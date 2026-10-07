package api

// Helper kecil bersama untuk handler API.

// pad3 — padStart(3, "0") untuk kode aksi ACT-001.
func pad3(n int) string {
	s := itos(n)
	for len(s) < 3 {
		s = "0" + s
	}
	return s
}

// srcRefSuffix — " (ref)" atau kosong.
func srcRefSuffix(ref string) string {
	if ref == "" {
		return ""
	}
	return " (" + ref + ")"
}
