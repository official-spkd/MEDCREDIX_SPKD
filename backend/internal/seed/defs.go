// Package seed — data demo MedCredix (port dari prisma/seed/*.ts).
// Semua data fiktif dan berlabel DEMO — tidak ada data pasien nyata.
package seed

import (
        "crypto/sha256"
        "encoding/hex"
        "time"
)

// sha256Hex — hashing kata sandi demo (SHA-256+salt; produksi: bcrypt/argon2).
func sha256Hex(s string) string {
        sum := sha256.Sum256([]byte(s))
        return hex.EncodeToString(sum[:])
}

// WIB — zona waktu acuan data demo (+07:00).
var WIB = time.FixedZone("WIB", 7*60*60)

// Now — titik nol data demo (deterministik, sama dengan seed asli).
var Now = time.Date(2026, 10, 3, 8, 30, 0, 0, WIB)

// D — tanggal relatif terhadap Now: offset hari + jam (default 09:00).
func D(offsetDays int, hour ...int) time.Time {
        h := 9
        if len(hour) > 0 {
                h = hour[0]
        }
        day := Now.Day() + offsetDays
        return time.Date(Now.Year(), Now.Month(), day, h, 0, 0, 0, WIB)
}

// MonthStr — string YYYY-MM relatif terhadap Now.
func MonthStr(offsetMonths int) string {
        m := int(Now.Month()) - 1 + offsetMonths
        y := Now.Year() + m/12
        m = ((m % 12) + 12) % 12
        months := []string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12"}
        return itoa(y) + "-" + months[m]
}

const (
        HNusantara = "h_nusantara"
        HMedika    = "h_medika"
        PW         = "medcredix2026"
)

func itoa(n int) string {
        if n == 0 {
                return "0"
        }
        neg := n < 0
        if neg {
                n = -n
        }
        var buf [12]byte
        i := len(buf)
        for n > 0 {
                i--
                buf[i] = byte('0' + n%10)
                n /= 10
        }
        if neg {
                i--
                buf[i] = '-'
        }
        return string(buf[i:])
}

// pad2 — padStart(2, "0").
func pad2(n int) string {
        if n < 10 {
                return "0" + itoa(n)
        }
        return itoa(n)
}
