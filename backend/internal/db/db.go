// Package db — koneksi database (SQLite via GORM, driver pure-Go), AutoMigrate,
// dan generator ID bergaya cuid untuk primary key string.
package db

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"medcredix/backend/internal/model"
)

var (
	DB    *gorm.DB
	idSeq atomic.Uint64
)

// Open membuka (atau membuat) database SQLite dan menjalankan AutoMigrate.
func Open(dsn string) *gorm.DB {
	if dsn == "" {
		dsn = os.Getenv("MEDCREDIX_DB")
	}
	if dsn == "" {
		dsn = "data/medcredix.db"
	}
	gdb, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("[db] gagal membuka %s: %v", dsn, err)
	}
	gdb.Exec("PRAGMA journal_mode = WAL;")
	gdb.Exec("PRAGMA busy_timeout = 5000;")
	gdb.Exec("PRAGMA foreign_keys = ON;")

	// Callback global: isi primary key string kosong dengan ID baru.
	_ = gdb.Callback().Create().Before("gorm:create").Register("medcredix:auto_id", func(tx *gorm.DB) {
		if tx.Statement.Schema == nil {
			return
		}
		field, ok := tx.Statement.Schema.FieldsByDBName["id"]
		if !ok || field.FieldType.String() != "string" {
			return
		}
		val, zero := field.ValueOf(tx.Statement.Context, tx.Statement.ReflectValue)
		if zero || val == "" {
			_ = field.Set(tx.Statement.Context, tx.Statement.ReflectValue, NewID())
		}
	})

	if err := gdb.AutoMigrate(model.AllModels()...); err != nil {
		log.Fatalf("[db] AutoMigrate gagal: %v", err)
	}
	DB = gdb
	return gdb
}

// NewID menghasilkan ID string unik bergaya cuid (c + base36 timestamp + counter + random).
func NewID() string {
	n := idSeq.Add(1)
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return fmt.Sprintf("c%08x%04x%02x%02x%02x",
		uint32(time.Now().UnixMilli()), uint16(n&0xffff), b[0], b[1], b[2])
}
