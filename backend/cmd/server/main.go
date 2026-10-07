// SPKD MedCredix API Server — backend Golang.
// "From Hospital Data to Action." — PT Sistem Pelayanan Kesehatan dan Data (SPKD)
//
// Menjalankan:
//   go run ./cmd/server           (dev)
//   go build -o medcredix-api ./cmd/server   (build)
//
// Environment:
//   MEDCREDIX_PORT   port (default 8081)
//   MEDCREDIX_DB     path SQLite (default data/medcredix.db)
//   MEDCREDIX_STATIC direktori build frontend (opsional, mode single-binary)
package main

import (
	"log"
	"os"

	"medcredix/backend/internal/api"
	"medcredix/backend/internal/db"
)

func main() {
	port := os.Getenv("MEDCREDIX_PORT")
	if port == "" {
		port = "8081"
	}
	gdb := db.Open("")
	if err := api.SeedIfEmpty(gdb); err != nil {
		log.Fatalf("[main] seed gagal: %v", err)
	}
	r := api.NewRouter()
	log.Printf("🚀 SPKD MedCredix API (Golang) berjalan di :%s — From Hospital Data to Action.", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
