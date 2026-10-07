// Package scoring — MedCredix Analytics Engine (port dari lib/scoring.ts).
// Formula version: medcredix-scorer/1.0
// Prinsip: SKOR → MENGAPA → BUKTI → RISIKO → REKOMENDASI → AKSI
package scoring

import (
        "math"
        "strconv"
        "strings"
        "time"

        "medcredix/backend/internal/model"
)

const FormulaVersion = "medcredix-scorer/1.0"

// StageThresholds — ambang Stage 0–7.
var StageThresholds = [][2]float64{{90, 7}, {80, 6}, {66, 5}, {60, 4}, {50, 3}, {35, 2}, {20, 1}, {0, 0}}

// StageOf — konversi skor ke Stage.
func StageOf(score float64) int {
        for _, t := range StageThresholds {
                if score >= t[0] {
                        return int(t[1])
                }
        }
        return 0
}

// StageRequirement — deskripsi ambang per stage.
func StageRequirement(stage int) string {
        reqs := []string{
                "Prasyarat wajib belum lengkap / skor < 20",
                "Skor 20–34: fondasi data mulai terbangun",
                "Skor 35–49: digitalisasi parsial",
                "Skor 50–59: mayoritas indikator aktif",
                "Skor 60–65: mendekati ambang kelulusan",
                "Skor 66–79: LULUS — kredensial digital terpenuhi",
                "Skor 80–89: LULUS — kematangan tinggi",
                "Skor ≥90: LULUS — excellence, layak jadi benchmark",
        }
        if stage > 7 {
                stage = 7
        }
        if stage < 0 {
                stage = 0
        }
        return reqs[stage]
}

// TrustMultiplier — pengali skor berdasarkan Trust (transparan, bukan tersembunyi).
func TrustMultiplier(trust float64) float64 {
        switch {
        case trust >= 90:
                return 1.0
        case trust >= 70:
                return 0.95
        case trust >= 50:
                return 0.85
        default:
                return 0.7
        }
}

// TrustBand — rentang Trust Multiplier untuk tampilan transparansi.
type TrustBand struct {
        MinTrust   float64 `json:"minTrust"`
        Multiplier float64 `json:"multiplier"`
        Label      string  `json:"label"`
}

var TrustBands = []TrustBand{
        {90, 1.0, "×1,00 — penuh"},
        {70, 0.95, "×0,95"},
        {50, 0.85, "×0,85"},
        {0, 0.7, "×0,70"},
}

// CcsGrade — grade klaim.
func CcsGrade(ccs float64) string {
        switch {
        case ccs >= 90:
                return "EXCELLENT"
        case ccs >= 75:
                return "GOOD"
        case ccs >= 60:
                return "FAIR"
        default:
                return "POOR"
        }
}

// ScoreStatus — label kualitas skor.
func ScoreStatus(s float64) string {
        switch {
        case s >= 90:
                return "excellent"
        case s >= 75:
                return "good"
        case s >= 60:
                return "fair"
        default:
                return "poor"
        }
}

// Indicator — indikator penyusun domain.
type Indicator struct {
        Label     string   `json:"label"`
        Value     float64  `json:"value"`
        Unit      string   `json:"unit"`
        Threshold *float64 `json:"threshold"`
        Status    string   `json:"status"`
        Weight    float64  `json:"weight"`
        Source    string   `json:"source"`
        Fhir      *string  `json:"fhir,omitempty"`
        Note      *string  `json:"note,omitempty"`
}

// Domain — domain skor.
type Domain struct {
        Key         string       `json:"key"`
        Label       string       `json:"label"`
        Score       float64      `json:"score"`
        Weight      float64      `json:"weight"`
        Group       string       `json:"group"`
        Formula     string       `json:"formula"`
        Indicators  []Indicator  `json:"indicators"`
        Positives   []string     `json:"positives"`
        Negatives   []string     `json:"negatives"`
}

type specialistCovRow struct {
        Specialty string  `json:"specialty"`
        Required  int     `json:"required"`
        Available int     `json:"available"`
        Coverage  float64 `json:"coverage"`
        Note      *string `json:"note"`
}

type ccsSubRow struct {
        Key    string  `json:"key"`
        Label  string  `json:"label"`
        Weight float64 `json:"weight"`
        Ok     int     `json:"ok"`
        Warn   int     `json:"warn"`
        Bad    int     `json:"bad"`
        Score  float64 `json:"score"`
}

type HospitalBrief struct {
        ID       string `json:"id"`
        Name     string `json:"name"`
        Class    string `json:"class"`
        City     string `json:"city"`
        Province string `json:"province"`
}

type BcdsInfo struct {
        Raw       float64 `json:"raw"`
        Trust     float64 `json:"trust"`
        Multiplier float64 `json:"multiplier"`
        Final     float64 `json:"final"`
        Stage     int     `json:"stage"`
        Lulus     bool    `json:"lulus"`
        PrereqOk  bool    `json:"prereqOk"`
}

type KelompokInfo struct {
        A float64 `json:"A"`
        B float64 `json:"B"`
        C float64 `json:"C"`
        D float64 `json:"D"`
}

type CcsInfo struct {
        Score         float64     `json:"score"`
        Grade         string      `json:"grade"`
        Subs          []ccsSubRow `json:"subs"`
        FlaggedClaims int         `json:"flaggedClaims"`
        TotalClaims   int         `json:"totalClaims"`
}

type CpcsInfo struct {
        Score         float64 `json:"score"`
        Coverage      float64 `json:"coverage"`
        Adherence     float64 `json:"adherence"`
        Documentation float64 `json:"documentation"`
}

type UrInfo struct {
        Alos              float64 `json:"alos"`
        Readmission       float64 `json:"readmission"`
        Bor               float64 `json:"bor"`
        Score             float64 `json:"score"`
        AlosScore         float64 `json:"alosScore"`
        ReadmScore        float64 `json:"readmScore"`
        BorScore          float64 `json:"borScore"`
        TargetAlos        float64 `json:"targetAlos"`
        TargetReadmission float64 `json:"targetReadmission"`
        TargetBorMin      float64 `json:"targetBorMin"`
        TargetBorMax      float64 `json:"targetBorMax"`
}

type SafetyInfo struct {
        Score      float64 `json:"score"`
        Allergy    float64 `json:"allergy"`
        ClosedLoop float64 `json:"closedLoop"`
        Discharge  float64 `json:"discharge"`
}

type TrustInfo struct {
        Score        float64             `json:"score"`
        Multiplier   float64             `json:"multiplier"`
        Factors      []model.TrustFactor `json:"factors"`
        AnomalyFlags int                 `json:"anomalyFlags"`
        UnderReview  bool                `json:"underReview"`
}

type WorkforceInfo struct {
        SpecialistCov []specialistCovRow `json:"specialistCov"`
        Coverage      float64            `json:"coverage"`
        Practitioners int                `json:"practitioners"`
}

type CredentialingInfo struct {
        Completeness   float64 `json:"completeness"`
        LicenseActive  float64 `json:"licenseActive"`
        Expired        int     `json:"expired"`
        Within30       int     `json:"within30"`
        Within90       int     `json:"within90"`
        PipelineHealth float64 `json:"pipelineHealth"`
}

type SaranaInfo struct {
        Score float64 `json:"score"`
        Total int     `json:"total"`
        Wajib int     `json:"wajib"`
}

type TrendItem struct {
        Period string  `json:"period"`
        Health float64 `json:"health"`
        Bcds   float64 `json:"bcds"`
        Trust  float64 `json:"trust"`
        Ccs    float64 `json:"ccs"`
        Cpcs   float64 `json:"cpcs"`
        Safety float64 `json:"safety"`
}

type PrereqItem struct {
        Code     string     `json:"code"`
        Name     string     `json:"name"`
        Status   string     `json:"status"`
        DocExpiry *time.Time `json:"docExpiry"`
        DaysLeft *int       `json:"daysLeft"`
}

// Intelligence — keluaran lengkap ComputeIntelligence.
type Intelligence struct {
        FormulaVersion string           `json:"formulaVersion"`
        ComputedAt     time.Time        `json:"computedAt"`
        DataAt         time.Time        `json:"dataAt"`
        Hospital       *HospitalBrief   `json:"hospital"`
        Domains        []Domain         `json:"domains"`
        HealthScore    float64          `json:"healthScore"`
        HealthStatus   string           `json:"healthStatus"`
        Bcds           BcdsInfo         `json:"bcds"`
        Kelompok       KelompokInfo     `json:"kelompok"`
        Ccs            CcsInfo          `json:"ccs"`
        Cpcs           CpcsInfo         `json:"cpcs"`
        Ur             UrInfo           `json:"ur"`
        Safety         SafetyInfo       `json:"safety"`
        Trust          TrustInfo        `json:"trust"`
        Workforce      WorkforceInfo    `json:"workforce"`
        Credentialing  CredentialingInfo `json:"credentialing"`
        Sarana         SaranaInfo       `json:"sarana"`
        Trend          []TrendItem      `json:"trend"`
        Prereq         []PrereqItem     `json:"prereq"`
}

// ---- helper matematis ----

func clamp(v, lo, hi float64) float64 { return math.Max(lo, math.Min(hi, v)) }

func clamp100(v float64) float64 { return clamp(v, 0, 100) }

func round1(v float64) float64 { return math.Round(v*10) / 10 }

type weightedVal struct{ value, weight float64 }

func wavg(rows []weightedVal) float64 {
        var sum, wsum float64
        for _, r := range rows {
                sum += r.value * r.weight
                wsum += r.weight
        }
        if wsum == 0 {
                return 0
        }
        return sum / wsum
}

func strPtr(s string) *string { return &s }

func fltPtr(f float64) *float64 { return &f }

func ind(label string, value float64, unit string, threshold *float64, weight float64, source string, fhir *string, note *string) Indicator {
        status := "—"
        if threshold != nil {
                if unit == "ms" || unit == "jam" {
                        if value <= *threshold {
                                status = "Baik"
                        } else if value <= *threshold*1.25 {
                                status = "Perlu Perhatian"
                        } else {
                                status = "Kritis"
                        }
                } else {
                        if value >= *threshold {
                                status = "Baik"
                        } else if value >= *threshold*0.85 {
                                status = "Perlu Perhatian"
                        } else {
                                status = "Kritis"
                        }
                }
        }
        return Indicator{
                Label: label, Value: round1(value), Unit: unit, Threshold: threshold,
                Status: status, Weight: weight, Source: source, Fhir: fhir, Note: note,
        }
}

// indicatorScore — skor indikator digital (arah diperhitungkan).
func indicatorScore(value, threshold float64, direction, unit string) float64 {
        if direction == "down" {
                if value <= threshold {
                        return 100
                }
                return clamp100(100 - ((value-threshold)/threshold)*100*1.2)
        }
        return clamp100(value)
}

var sipBad = map[string]bool{"Kedaluwarsa": true, "Ditangguhkan": true, "Dicabut": true}

// IntelData — kumpulan data mentah yang dibutuhkan engine skor.
type IntelData struct {
        Practitioners []model.Practitioner
        Credentials   []model.PractitionerCredential
        Stds          []model.SpecialtyStandard
        Facilities    []model.Facility
        Beds          []model.BedStat
        DigInd        []model.DigitalIndicator
        Claims        []model.Claim
        Pathways      []model.ClinicalPathway
        Episodes      []model.PathwayEpisode
        Utilizations  []model.UtilizationStat
        Complaints    []model.ComplaintStat
        Trusts        []model.TrustFactor
        DqIssues      []model.DataQualityIssue
        Requirements  []model.Requirement
        Hospital      *model.Hospital
        Snapshots     []model.ScoreSnapshot
}

// Loader — injeksi dependensi agar dapat diuji & dipakai lintas handler.
type Loader interface {
        Load(hospitalID string) (IntelData, error)
}

// ComputeIntelligence — port computeIntelligence() dari lib/scoring.ts.
func ComputeIntelligence(l Loader, hospitalID string) (*Intelligence, error) {
        d, err := l.Load(hospitalID)
        if err != nil {
                return nil, err
        }
        now := time.Now()
        days := func(dt *time.Time) int {
                if dt == nil {
                        return 9999
                }
                return int(math.Ceil(float64(dt.Sub(now)) / float64(24*time.Hour)))
        }
        findDig := func(contains string) (model.DigitalIndicator, bool) {
                for _, x := range d.DigInd {
                        if strings.Contains(x.Indicator, contains) {
                                return x, true
                        }
                }
                return model.DigitalIndicator{}, false
        }

        // ===== 1. CLINICAL WORKFORCE =====
        activePractitioners := make([]model.Practitioner, 0, len(d.Practitioners))
        for _, p := range d.Practitioners {
                if p.Status == "Aktif" {
                        activePractitioners = append(activePractitioners, p)
                }
        }
        specialistCov := make([]specialistCovRow, 0, len(d.Stds))
        for _, s := range d.Stds {
                available := 0
                for _, p := range activePractitioners {
                        if (p.Profession == "Dokter Spesialis" || p.Profession == "Dokter Subspesialis") && p.Specialty != nil && *p.Specialty == s.Specialty {
                                available++
                        }
                }
                specialistCov = append(specialistCov, specialistCovRow{
                        Specialty: s.Specialty, Required: s.Required, Available: available,
                        Coverage: clamp100(float64(available) / float64(s.Required) * 100), Note: s.Note,
                })
        }
        var coverage float64
        if len(specialistCov) > 0 {
                rows := make([]weightedVal, len(specialistCov))
                for i, s := range specialistCov {
                        rows[i] = weightedVal{s.Coverage, 1}
                }
                coverage = wavg(rows)
        }
        sipActive := 0
        for _, p := range d.Practitioners {
                if p.SipExpiry != nil && days(p.SipExpiry) >= 0 {
                        sipActive++
                }
        }
        sipActivePct := clamp100(float64(sipActive) / float64(max(1, len(d.Practitioners))) * 100)
        activePct := clamp100(float64(len(activePractitioners)) / float64(max(1, len(d.Practitioners))) * 100)
        metCount := 0
        for _, s := range specialistCov {
                if s.Coverage >= 100 {
                        metCount++
                }
        }
        subSpec := 0
        for _, p := range d.Practitioners {
                if strings.Contains(p.Profession, "Subspesialis") {
                        subSpec++
                }
        }
        workforce := Domain{
                Key: "workforce", Label: "Clinical Workforce", Weight: 0.12, Group: "A. Kapabilitas Klinis & SDM",
                Formula: "Cakupan spesialis (50%) + SIP aktif (30%) + status tenaga aktif (20%) — cakupan = rata-rata min(tersedia/wajib) per spesialis (Permenkes 56/2014)",
                Score: round1(wavg([]weightedVal{{coverage, 0.5}, {sipActivePct, 0.3}, {activePct, 0.2}})),
                Indicators: []Indicator{
                        ind("Cakupan spesialis wajib", coverage, "%", fltPtr(100), 0.5, "HRIS + Permenkes 56/2014", strPtr("Practitioner"),
                                strPtr(itoa(metCount)+"/"+itoa(len(specialistCov))+" spesialis terpenuhi")),
                        ind("Praktisi dengan SIP masih berlaku", sipActivePct, "%", fltPtr(95), 0.3, "HRIS", strPtr("Practitioner"), nil),
                        ind("Status tenaga aktif", activePct, "%", fltPtr(90), 0.2, "HRIS", strPtr("Practitioner"), nil),
                },
                Positives: []string{
                        itoa(metCount) + " dari " + itoa(len(specialistCov)) + " spesialis wajib terpenuhi",
                        itoa(subSpec) + " subspesialis tersedia (di atas standar minimal)",
                },
        }
        for _, s := range specialistCov {
                if s.Coverage < 100 {
                        note := "perlu rekrutmen"
                        if s.Note != nil {
                                note = *s.Note
                        }
                        workforce.Negatives = append(workforce.Negatives, s.Specialty+": "+itoa(s.Available)+"/"+itoa(s.Required)+" — "+note)
                }
        }

        // ===== 2. CREDENTIALING READINESS =====
        byPractitioner := map[string][]model.PractitionerCredential{}
        for _, c := range d.Credentials {
                byPractitioner[c.PractitionerID] = append(byPractitioner[c.PractitionerID], c)
        }
        firstByType := func(creds []model.PractitionerCredential, t string) *model.PractitionerCredential {
                for i := range creds {
                        if creds[i].Type == t {
                                return &creds[i]
                        }
                }
                return nil
        }
        completeCount, licenseActiveCount := 0, 0
        for _, p := range d.Practitioners {
                creds := byPractitioner[p.ID]
                str := firstByType(creds, "STR")
                sip := firstByType(creds, "SIP")
                okStr := str != nil && !sipBad[str.Status]
                okSip := sip != nil && !sipBad[sip.Status]
                if okStr && okSip {
                        completeCount++
                        licenseActiveCount++
                }
        }
        completeness := clamp100(float64(completeCount) / float64(max(1, len(d.Practitioners))) * 100)
        licenseActive := clamp100(float64(licenseActiveCount) / float64(max(1, len(d.Practitioners))) * 100)
        activeCredStates := map[string]bool{"Aktif": true, "Segera Berakhir": true, "Disetujui": true, "Kedaluwarsa": true, "Dicabut": true}
        finalStateCount := 0
        for _, c := range d.Credentials {
                if activeCredStates[c.Status] {
                        finalStateCount++
                }
        }
        pipelineHealth := clamp100(float64(finalStateCount) / float64(max(1, len(d.Credentials))) * 100)
        expired, within30, within90 := 0, 0, 0
        for _, c := range d.Credentials {
                if c.ExpiryDate != nil {
                        dd := days(c.ExpiryDate)
                        if dd < 0 {
                                expired++
                        } else if dd <= 30 {
                                within30++
                        } else if dd <= 90 {
                                within90++
                        }
                }
        }
        expiryScore := clamp100(100 - float64(expired*10+within30*4+within90*2))
        credentialing := Domain{
                Key: "credentialing", Label: "Credentialing Readiness", Weight: 0.14, Group: "A. Kapabilitas Klinis & SDM",
                Formula: "Kelengkapan STR+SIP praktisi (40%) + lisensi aktif (25%) + kesehatan pipeline (15%) + skor expiry (20%) — expiry = 100 − (kedaluwarsa×10 + ≤30h×4 + ≤90h×2)",
                Score: round1(wavg([]weightedVal{{completeness, 0.4}, {licenseActive, 0.25}, {pipelineHealth, 0.15}, {expiryScore, 0.2}})),
                Indicators: []Indicator{
                        ind("Kelengkapan kredensial inti (STR+SIP)", completeness, "%", fltPtr(95), 0.4, "HRIS + Credentialing", strPtr("Practitioner + PractitionerRole"),
                                strPtr(itoa(completeCount)+"/"+itoa(len(d.Practitioners))+" praktisi lengkap")),
                        ind("Lisensi aktif (STR & SIP tidak bermasalah)", licenseActive, "%", fltPtr(95), 0.25, "HRIS", strPtr("Practitioner"), nil),
                        ind("Kredensial pada status final", pipelineHealth, "%", fltPtr(90), 0.15, "Credentialing", strPtr("PractitionerCredential"), nil),
                        ind("Skor kesehatan expiry", expiryScore, "/100", fltPtr(80), 0.2, "Credentialing", strPtr("PractitionerCredential"),
                                strPtr(itoa(expired)+" kedaluwarsa; "+itoa(within30)+" ≤30 hari; "+itoa(within90)+" ≤90 hari")),
                },
                Positives: []string{itoa(completeCount) + "/" + itoa(len(d.Practitioners)) + " praktisi memiliki STR & SIP lengkap"},
        }
        if within30 > 0 {
                credentialing.Negatives = append(credentialing.Negatives, itoa(within30)+" dokumen berakhir dalam 30 hari")
        }
        if expired > 0 {
                credentialing.Negatives = append(credentialing.Negatives, itoa(expired)+" dokumen telah kedaluwarsa")
        }
        if within90 > 0 {
                credentialing.Negatives = append(credentialing.Negatives, itoa(within90)+" dokumen berakhir dalam 31–90 hari")
        }

        // ===== Dimensi digital =====
        dimScore := func(dim string) (float64, []Indicator) {
                var rows []model.DigitalIndicator
                for _, r := range d.DigInd {
                        if r.Dimension == dim {
                                rows = append(rows, r)
                        }
                }
                if len(rows) == 0 {
                        return 0, []Indicator{}
                }
                indicators := make([]Indicator, 0, len(rows))
                var wrows []weightedVal
                for _, r := range rows {
                        th := r.Threshold
                        indicators = append(indicators, ind(r.Indicator, r.Value, r.Unit, &th, r.Weight, r.SourceSystem, r.FhirResource, nil))
                        indicators[len(indicators)-1].Status = r.Status
                        wrows = append(wrows, weightedVal{indicatorScore(r.Value, r.Threshold, r.Direction, r.Unit), r.Weight})
                }
                return round1(wavg(wrows)), indicators
        }
        _, emrInd := dimScore("EMR")
        infraScore, infraInd := dimScore("Infrastruktur")
        medScore, medInd := dimScore("Medikasi")
        diagScore, diagInd := dimScore("Diagnostik")
        interopScore, interopInd := dimScore("Interoperabilitas")
        _, safetyInd := dimScore("Keselamatan")
        emrRi, emrRj := 0.0, 0.0
        if x, ok := findDig("Rawat Inap"); ok && x.Dimension == "EMR" {
                emrRi = x.Value
        }
        if x, ok := findDig("Rawat Jalan"); ok && x.Dimension == "EMR" {
                emrRj = x.Value
        }
        emrScore := round1(emrRi*0.55 + emrRj*0.45) // formula legacy TransDCA dipertahankan

        uptimeVal := 0.0
        if x, ok := findDig("Uptime"); ok {
                uptimeVal = x.Value
        }
        msVal := 0.0
        for _, x := range d.DigInd {
                if x.Unit == "ms" {
                        msVal = x.Value
                        break
                }
        }
        digital := Domain{
                Key: "digital", Label: "Digital Maturity", Weight: 0.10, Group: "B. Digitalisasi & Sistem Informasi",
                Formula: "Skor EMR = (Rawat Inap × 55% + Rawat Jalan × 45%) — formula legacy dipertahankan; skor dimensi = rata-rata tertimbang indikator vs ambang",
                Score:        round1(wavg([]weightedVal{{emrScore, 0.5}, {infraScore, 0.5}})),
                Indicators:   append(append([]Indicator{}, emrInd...), infraInd...),
                Positives:    []string{"Adopsi RME rawat inap " + ftoa(emrRi) + "% (ambang 85%)", "Uptime API " + ftoa(uptimeVal) + "%"},
                Negatives:    []string{"Adopsi RME rawat jalan " + ftoa(emrRj) + "% masih di bawah ambang 85%", "Latensi API " + ftoa(msVal) + " ms di atas ambang 300 ms"},
        }
        cpoeVal := 0.0
        for _, x := range d.DigInd {
                if strings.HasPrefix(x.Indicator, "CPOE") {
                        cpoeVal = x.Value
                        break
                }
        }
        adminVal := 0.0
        if x, ok := findDig("Administrasi"); ok {
                adminVal = x.Value
        }
        medication := Domain{
                Key: "medication", Label: "Medication Management", Weight: 0.08, Group: "B. Digitalisasi & Sistem Informasi",
                Formula: "Rata-rata tertimbang CPOE, dispensing, administrasi, closed-loop, kelengkapan alergi (ambang masing-masing indikator)",
                Score: medScore, Indicators: medInd,
                Positives: []string{"CPOE " + ftoa(cpoeVal) + "%"},
                Negatives: []string{"Administrasi obat digital " + ftoa(adminVal) + "% — di bawah ambang 80%"},
        }
        labVal := 0.0
        if x, ok := findDig("laboratorium"); ok {
                labVal = x.Value
        }
        critVal := 0.0
        if x, ok := findDig("Critical"); ok {
                critVal = x.Value
        }
        diagnostic := Domain{
                Key: "diagnostic", Label: "Diagnostic Readiness", Weight: 0.08, Group: "B. Digitalisasi & Sistem Informasi",
                Formula: "Rata-rata tertimbang digitalisasi lab/radiologi, kelengkapan laporan, TAT, critical value flagging",
                Score: diagScore, Indicators: diagInd,
                Positives: []string{"Digitalisasi laboratorium " + ftoa(labVal) + "%"},
                Negatives: []string{"Critical value flagging " + ftoa(critVal) + "% di bawah ambang 95%"},
        }
        valFhir := 0.0
        if x, ok := findDig("Validasi"); ok {
                valFhir = x.Value
        }
        interopDomain := Domain{
                Key: "interop", Label: "Interoperability", Weight: 0.08, Group: "B. Digitalisasi & Sistem Informasi",
                Formula: "Rata-rata tertimbang submission SATUSEHAT, validasi FHIR, sukses integrasi BPJS (TDS v2.0.0 — 24 resource)",
                Score: interopScore, Indicators: interopInd,
                Positives: []string{"Coverage resource FHIR 24/24 tersedia"},
                Negatives: []string{"Validasi FHIR " + ftoa(valFhir) + "% — sebagian resource ditolak (field wajib kosong)"},
        }

        // ===== 7. CLAIM INTEGRITY (CCS) =====
        cats := []struct {
                key    string
                label  string
                weight float64
        }{
                {"Diagnosis", "Diagnosis Consistency", 0.25},
                {"Medication", "Medication Consistency", 0.20},
                {"Procedure", "Procedure Consistency", 0.15},
                {"Diagnostic", "Diagnostic Consistency", 0.15},
                {"Stock", "Konsistensi Stok & Klaim (Triangle)", 0.15},
                {"Timeline", "Completeness & Timeline", 0.10},
        }
        claimCount := max(1, len(d.Claims))
        subs := make([]ccsSubRow, 0, len(cats))
        flagCounts := make([]map[string]string, len(d.Claims))
        flaggedClaims := 0
        for i, c := range d.Claims {
                flags := parseFlags(c.FlagsJson)
                flagCounts[i] = flags
                for _, v := range flags {
                        if v != "ok" {
                                flaggedClaims++
                                break
                        }
                }
        }
        for _, cat := range cats {
                ok, warn, bad := 0, 0, 0
                for _, flags := range flagCounts {
                        switch flags[cat.key] {
                        case "ok":
                                ok++
                        case "warn":
                                warn++
                        case "bad":
                                bad++
                        }
                }
                subs = append(subs, ccsSubRow{
                        Key: cat.key, Label: cat.label, Weight: cat.weight, Ok: ok, Warn: warn, Bad: bad,
                        Score: round1((float64(ok) + float64(warn)*0.5) / float64(claimCount) * 100),
                })
        }
        var ccsRows []weightedVal
        for _, s := range subs {
                ccsRows = append(ccsRows, weightedVal{s.Score, s.Weight})
        }
        ccs := round1(wavg(ccsRows))
        claimDomain := Domain{
                Key: "claim", Label: "Claim Integrity", Weight: 0.16, Group: "C. Kinerja JKN",
                Formula: "CCS = Σ(sub-skor × bobot): Diagnosis 25% + Medication 20% + Procedure 15% + Diagnostic 15% + Stock 15% + Timeline 10%; sub-skor = (ok + warn×0,5) / total klaim",
        }
        for _, s := range subs {
                note := itoa(s.Ok) + " ok · " + itoa(s.Warn) + " perlu tinjauan · " + itoa(s.Bad) + " tinggi"
                claimDomain.Indicators = append(claimDomain.Indicators, ind(s.Label, s.Score, "%", fltPtr(85), s.Weight, "BPJS VClaim × EMR/RME", strPtr("Claim"), strPtr(note)))
        }
        claimDomain.Score = ccs
        claimDomain.Positives = []string{itoa(len(d.Claims)-flaggedClaims) + " klaim bersih tanpa sinyal", "Konsistensi diagnosis " + ftoa(subs[0].Score) + "%"}
        claimDomain.Negatives = []string{itoa(flaggedClaims) + " klaim memiliki minimal satu sinyal inkonsistensi"}
        for _, s := range subs {
                if s.Score < 80 {
                        claimDomain.Negatives = append(claimDomain.Negatives, s.Label+": "+ftoa(s.Score)+"%")
                }
        }

        // ===== 8. CLINICAL PATHWAY (CPCS) =====
        sumAdh := 0.0
        for _, e := range d.Episodes {
                sumAdh += e.Adherence
        }
        avgAdh := 0.0
        if len(d.Episodes) > 0 {
                avgAdh = sumAdh / float64(len(d.Episodes))
        }
        devEpisodes := make([]model.PathwayEpisode, 0)
        for _, e := range d.Episodes {
                if e.DeviationType != "Tidak Ada" {
                        devEpisodes = append(devEpisodes, e)
                }
        }
        documentedPct := 100.0
        if len(devEpisodes) > 0 {
                doc := 0
                for _, e := range devEpisodes {
                        if e.Documented {
                                doc++
                        }
                }
                documentedPct = float64(doc) / float64(len(devEpisodes)) * 100
        }
        cpTotal, cpCovered := 14, 0
        if len(d.Pathways) >= 8 {
                cpTotal, cpCovered = 20, 16
        } else {
                cpCovered = min(len(d.Pathways), 11)
        }
        cpCoveragePct := float64(cpCovered) / float64(cpTotal) * 100
        cpcs := round1(cpCoveragePct*0.25 + avgAdh*0.5 + documentedPct*0.25)
        docCount := 0
        for _, e := range devEpisodes {
                if e.Documented {
                        docCount++
                }
        }
        pathwayDomain := Domain{
                Key: "pathway", Label: "Clinical Pathway", Weight: 0.10, Group: "C. Kinerja JKN",
                Formula: "CPCS = Coverage 25% (diagnosis teratas yang punya CP) + Adherence 50% (kepatuhan rata-rata episode) + Dokumentasi Deviasi 25%",
                Score: cpcs,
                Indicators: []Indicator{
                        ind("Ketersediaan panduan klinis (coverage)", cpCoveragePct, "%", fltPtr(85), 0.25, "EMR/RME", strPtr("CarePlan"),
                                strPtr(itoa(cpCovered)+"/"+itoa(cpTotal)+" diagnosis utama")),
                        ind("Kepatuhan klinisi (adherence)", avgAdh, "%", fltPtr(85), 0.5, "EMR/RME", strPtr("CarePlan"),
                                strPtr(itoa(len(d.Episodes))+" episode dianalisis")),
                        ind("Dokumentasi deviasi", documentedPct, "%", fltPtr(75), 0.25, "EMR/RME", strPtr("CarePlan"),
                                strPtr(itoa(docCount)+"/"+itoa(len(devEpisodes))+" deviasi terdokumentasi")),
                },
                Positives: []string{itoa(len(d.Pathways)) + " clinical pathway aktif"},
        }
        lowAdh := 0
        for _, p := range d.Pathways {
                if p.AvgAdherence < 75 && lowAdh < 2 {
                        pathwayDomain.Negatives = append(pathwayDomain.Negatives, p.Diagnosis+": adherence "+ftoa(p.AvgAdherence)+"%")
                        lowAdh++
                }
        }
        pathwayDomain.Negatives = append(pathwayDomain.Negatives, itoa(len(devEpisodes)-docCount)+" deviasi tanpa alasan klinis terdokumentasi")

        // ===== 9. PATIENT SAFETY =====
        allergy, closedLoop, discharge := 0.0, 0.0, 0.0
        if x, ok := findDig("Alergi"); ok {
                allergy = x.Value
        }
        if x, ok := findDig("High-Alert"); ok {
                closedLoop = x.Value
        }
        if x, ok := findDig("Pemulangan"); ok {
                discharge = x.Value
        }
        safetyScore := round1(allergy*0.4 + closedLoop*0.3 + discharge*0.3)
        safetyDomain := Domain{
                Key: "safety", Label: "Patient Safety", Weight: 0.14, Group: "D. Keselamatan Pasien",
                Formula: "Skor Keselamatan = Alergi Obat × 40% + High-Alert & Closed-Loop × 30% + Dokumentasi Pemulangan × 30% (akumulasi 3 domain — bukan rata-rata)",
                Score: safetyScore, Indicators: safetyInd,
                Positives: []string{"Keamanan alergi obat " + ftoa(allergy) + "% (ambang 95%)", "Dokumentasi pemulangan " + ftoa(discharge) + "% (ambang 80%)"},
                Negatives: []string{"High-alert & closed-loop " + ftoa(closedLoop) + "% di bawah ambang 80%"},
        }

        // ===== SARANA (untuk kelompok A) =====
        wajib := make([]model.Facility, 0)
        for _, f := range d.Facilities {
                if f.Notes == nil || !strings.Contains(strings.ToLower(*f.Notes), "opsional") {
                        wajib = append(wajib, f)
                }
        }
        saranaScore := 0.0
        if len(wajib) > 0 {
                rows := make([]weightedVal, len(wajib))
                for i, f := range wajib {
                        v := 40.0
                        if f.Availability == "Ada" && f.Functional == "Berfungsi" && f.Feasible == "Laik Pakai" && f.Calibration == "Valid" {
                                v = 100
                        } else if f.Feasible == "Perlu Perbaikan" {
                                v = 75
                        }
                        rows[i] = weightedVal{v, 1}
                }
                saranaScore = round1(wavg(rows))
        }

        domains := []Domain{workforce, credentialing, digital, medication, diagnostic, interopDomain, claimDomain, pathwayDomain, safetyDomain}
        var hRows []weightedVal
        for _, dm := range domains {
                hRows = append(hRows, weightedVal{dm.Score, dm.Weight})
        }
        healthScore := round1(wavg(hRows))

        // ===== BDCS (kelompok) =====
        kelompokA := round1(wavg([]weightedVal{{workforce.Score, 0.5}, {saranaScore, 0.5}}))
        kelompokB := round1(wavg([]weightedVal{
                {emrScore, 0.25}, {infraScore, 0.15}, {medScore, 0.2}, {diagScore, 0.15}, {interopScore, 0.25},
        }))
        var alosV, readmV, borV float64
        if len(d.Utilizations) > 0 {
                u := d.Utilizations[len(d.Utilizations)-1]
                alosV, readmV, borV = u.Alos, u.Readmission, u.Bor
        }
        alosScore := clamp100(100 - math.Max(0, (alosV-5)*20))
        readmScore := clamp100(100 - math.Max(0, (readmV-5)*15))
        borScore := 100.0
        if borV < 60 || borV > 85 {
                borScore = clamp100(100 - math.Abs(map[bool]float64{true: 60 - borV, false: borV - 85}[borV < 60])*2)
        }
        urScore := round1((alosScore + readmScore + borScore) / 3)
        kelompokC := round1(wavg([]weightedVal{{ccs, 0.5}, {cpcs, 0.35}, {urScore, 0.15}}))
        kelompokD := safetyScore
        bcdsRaw := round1(kelompokA*0.25 + kelompokB*0.3 + kelompokC*0.35 + kelompokD*0.1)

        // ===== TRUST =====
        trust := 100.0
        if len(d.Trusts) > 0 {
                rows := make([]weightedVal, len(d.Trusts))
                for i, t := range d.Trusts {
                        rows[i] = weightedVal{t.Score, t.Weight}
                }
                trust = round1(wavg(rows))
        }
        multiplier := TrustMultiplier(trust)
        bcds := round1(bcdsRaw * multiplier)
        stage := StageOf(bcds)
        anomalyFlags := 0
        for _, q := range d.DqIssues {
                if q.Anomaly && q.Status != "Diperbaiki" {
                        anomalyFlags++
                }
        }
        underReview := anomalyFlags >= 3
        prereqOk := true
        for _, r := range d.Requirements {
                if r.Code == "P1" || r.Code == "P2" || r.Code == "P3" {
                        if r.Status != "Valid" {
                                prereqOk = false
                        }
                }
        }
        effectiveStage := stage
        if !prereqOk {
                effectiveStage = 0
        }

        // ===== Trend & Prasyarat =====
        trend := make([]TrendItem, 0, len(d.Snapshots))
        for _, s := range d.Snapshots {
                trend = append(trend, TrendItem{Period: s.Period, Health: s.HealthScore, Bcds: s.Bcds, Trust: s.Trust, Ccs: s.Ccs, Cpcs: s.Cpcs, Safety: s.SafetyScore})
        }
        prereq := make([]PrereqItem, 0, len(d.Requirements))
        for _, r := range d.Requirements {
                item := PrereqItem{Code: r.Code, Name: r.Name, Status: r.Status, DocExpiry: r.DocExpiry}
                if r.DocExpiry != nil {
                        dl := days(r.DocExpiry)
                        item.DaysLeft = &dl
                }
                prereq = append(prereq, item)
        }

        intel := &Intelligence{
                FormulaVersion: FormulaVersion,
                ComputedAt:     now,
                DataAt:         now,
                Domains:        domains,
                HealthScore:    healthScore,
                HealthStatus:   ScoreStatus(healthScore),
                Bcds: BcdsInfo{Raw: bcdsRaw, Trust: trust, Multiplier: multiplier, Final: bcds, Stage: effectiveStage, Lulus: effectiveStage >= 5, PrereqOk: prereqOk},
                Kelompok:       KelompokInfo{A: kelompokA, B: kelompokB, C: kelompokC, D: kelompokD},
                Ccs:            CcsInfo{Score: ccs, Grade: CcsGrade(ccs), Subs: subs, FlaggedClaims: flaggedClaims, TotalClaims: len(d.Claims)},
                Cpcs:           CpcsInfo{Score: cpcs, Coverage: round1(cpCoveragePct), Adherence: round1(avgAdh), Documentation: round1(documentedPct)},
                Ur: UrInfo{Alos: alosV, Readmission: readmV, Bor: borV, Score: urScore, AlosScore: alosScore, ReadmScore: readmScore, BorScore: borScore,
                        TargetAlos: 5, TargetReadmission: 5, TargetBorMin: 60, TargetBorMax: 85},
                Safety:      SafetyInfo{Score: safetyScore, Allergy: allergy, ClosedLoop: closedLoop, Discharge: discharge},
                Trust:       TrustInfo{Score: trust, Multiplier: multiplier, Factors: d.Trusts, AnomalyFlags: anomalyFlags, UnderReview: underReview},
                Workforce:   WorkforceInfo{SpecialistCov: specialistCov, Coverage: round1(coverage), Practitioners: len(d.Practitioners)},
                Credentialing: CredentialingInfo{Completeness: completeness, LicenseActive: licenseActive, Expired: expired,
                        Within30: within30, Within90: within90, PipelineHealth: pipelineHealth},
                Sarana: SaranaInfo{Score: saranaScore, Total: len(d.Facilities), Wajib: len(wajib)},
                Trend:  trend,
                Prereq: prereq,
        }
        if d.Hospital != nil {
                intel.Hospital = &HospitalBrief{ID: d.Hospital.ID, Name: d.Hospital.Name, Class: d.Hospital.Class, City: d.Hospital.City, Province: d.Hospital.Province}
        }
        return intel, nil
}

// Recommendation — rekomendasi rule-based per domain (arsitektur siap AI assistive).
type Recommendation struct {
        Title          string `json:"title"`
        Priority       string `json:"priority"`
        SuggestedOwner string `json:"suggestedOwner"`
}

// DomainRecommendations — port domainRecommendations() dari lib/scoring.ts.
func DomainRecommendations(d *Domain, intel *Intelligence) []Recommendation {
        var recs []Recommendation
        switch d.Key {
        case "credentialing":
                if intel.Credentialing.Expired > 0 {
                        recs = append(recs, Recommendation{
                                Title: "Segera perpanjang " + itoa(intel.Credentialing.Expired) + " dokumen kedaluwarsa & proses " + itoa(intel.Credentialing.Within30) + " dokumen ≤30 hari",
                                Priority: "Tinggi", SuggestedOwner: "Petugas Kredensial"})
                }
                recs = append(recs, Recommendation{Title: "Lengkapi dokumen praktisi yang belum lengkap melalui alur verifikasi", Priority: "Sedang", SuggestedOwner: "HR/SDM"})
        case "claim":
                recs = append(recs, Recommendation{Title: "Tinjau " + itoa(intel.Ccs.FlaggedClaims) + " klaim dengan sinyal inkonsistensi (prioritas kategori terlemah)", Priority: "Tinggi", SuggestedOwner: "Tim Klaim"})
                recs = append(recs, Recommendation{Title: "Rekonsiliasi stok-klaim (triangle) untuk kategori dengan sinyal tertinggi", Priority: "Sedang", SuggestedOwner: "Farmasi Klinik"})
        case "pathway":
                recs = append(recs, Recommendation{Title: "Audit episode dengan adherence <70% & lengkapi dokumentasi alasan deviasi", Priority: "Sedang", SuggestedOwner: "Komite Medis"})
                recs = append(recs, Recommendation{Title: "Susun clinical pathway untuk 4 diagnosis teratas yang belum tercakup", Priority: "Sedang", SuggestedOwner: "Komite Medis"})
        case "medication":
                recs = append(recs, Recommendation{Title: "Deploy e-MAR + barcode administrasi obat hingga ≥85%", Priority: "Tinggi", SuggestedOwner: "Ka. Keperawatan"})
        case "interop":
                recs = append(recs, Recommendation{Title: "Perbaiki field wajib FHIR yang kosong (Composition, MedicationDispense) & jalankan Test Koneksi", Priority: "Sedang", SuggestedOwner: "IT/SIMRS"})
        case "safety":
                recs = append(recs, Recommendation{Title: "Tingkatkan closed-loop medication ke ambang 80%", Priority: "Tinggi", SuggestedOwner: "Farmasi Klinik"})
        case "workforce":
                for _, s := range intel.Workforce.SpecialistCov {
                        if s.Coverage < 100 {
                                recs = append(recs, Recommendation{Title: "Rekrut spesialis yang belum terpenuhi (lihat rincian cakupan)", Priority: "Sedang", SuggestedOwner: "HR/SDM"})
                                break
                        }
                }
        case "digital":
                recs = append(recs, Recommendation{Title: "Tingkatkan adopsi RME rawat jalan menuju ambang 85%", Priority: "Sedang", SuggestedOwner: "IT/SIMRS"})
        case "diagnostic":
                recs = append(recs, Recommendation{Title: "Aktifkan critical value flagging otomatis dari LIS ke RME", Priority: "Sedang", SuggestedOwner: "IT/SIMRS"})
        }
        return recs
}

// parseFlags — parse flagsJson klaim menjadi map kategori → flag.
func parseFlags(js *string) map[string]string {
        out := map[string]string{}
        if js == nil {
                return out
        }
        // format: {"Diagnosis":"ok",...} — parse sederhana tanpa dependensi
        s := *js
        s = strings.TrimPrefix(strings.TrimSuffix(s, "}"), "{")
        for _, part := range strings.Split(s, ",") {
                kv := strings.SplitN(part, ":", 2)
                if len(kv) != 2 {
                        continue
                }
                k := strings.Trim(strings.TrimSpace(kv[0]), `"`)
                v := strings.Trim(strings.TrimSpace(kv[1]), `"`)
                out[k] = v
        }
        return out
}

func itoa(n int) string { return strconv.Itoa(n) }

func fmtInt(n int) string { return strconv.Itoa(n) }

func fmtF(v float64) string { return strconv.FormatFloat(v, 'f', 1, 64) }

func ftoa(v float64) string {
        // angka tampil tanpa desimal bila bulat (menyerupai template literal TS)
        if v == math.Trunc(v) {
                return fmtInt(int(v))
        }
        return fmtF(v)
}
