// Package model — SPKD MedCredix data models (port dari Prisma schema → GORM).
// Semua entitas bertenant (HospitalID) untuk isolasi multi-rumah sakit.
// Tag JSON mengikuti serialisasi camelCase Prisma agar kontrak API identik.
package model

import "time"

// ============ TENANT & IDENTITAS ============

type Hospital struct {
	ID               string    `gorm:"primaryKey;size:64" json:"id"`
	Name             string    `json:"name"`
	Code             string    `gorm:"uniqueIndex;size:32" json:"code"`
	Class            string    `gorm:"size:4" json:"class"`
	City             string    `json:"city"`
	Province         string    `json:"province"`
	BedLicensed      int       `json:"bedLicensed"`
	BedActive        int       `json:"bedActive"`
	BedNotOperating  int       `json:"bedNotOperating"`
	NotOperatingNote *string   `json:"notOperatingNote"`
	BpjsFacilityCode *string   `json:"bpjsFacilityCode"`
	SatusehatOrgID   *string   `gorm:"column:satusehat_org_id" json:"satusehatOrgId"`
	Tier             string    `gorm:"size:32" json:"tier"`
	Onboarding       bool      `json:"onboarding"`
	CreatedAt        time.Time `json:"createdAt"`
}

type User struct {
	ID           string     `gorm:"primaryKey;size:64" json:"id"`
	HospitalID   *string    `gorm:"column:hospital_id;index" json:"hospitalId"`
	Email        string     `gorm:"uniqueIndex;size:128" json:"email"`
	PasswordHash string     `gorm:"size:128" json:"-"`
	Name         string     `json:"name"`
	Role         string     `gorm:"size:32" json:"role"`
	Title        *string    `json:"title"`
	Phone        *string    `json:"phone"`
	Active       bool       `json:"active"`
	LastLogin    *time.Time `json:"lastLogin"`
	CreatedAt    time.Time  `json:"createdAt"`
}

// ============ CREDENTIALING ============

type Practitioner struct {
	ID             string     `gorm:"primaryKey;size:64" json:"id"`
	HospitalID     string     `gorm:"column:hospital_id;index:idx_practitioner_hosp_prof;index:idx_practitioner_hosp_spec" json:"hospitalId"`
	FullName       string     `json:"fullName"`
	NIKMasked      *string    `gorm:"column:nik_masked" json:"nikMasked"`
	Profession     string     `json:"profession"`
	Specialty      *string    `json:"specialty"`
	SubSpecialty   *string    `json:"subSpecialty"`
	Department     string     `json:"department"`
	EmploymentType string     `json:"employmentType"`
	Status         string     `json:"status"`
	StrNumber      *string    `json:"strNumber"`
	StrExpiry      *time.Time `json:"strExpiry"`
	SipNumber      *string    `json:"sipNumber"`
	SipExpiry      *time.Time `json:"sipExpiry"`
	Phone          *string    `json:"phone"`
	JoinDate       *time.Time `json:"joinDate"`
	FhirID         *string    `gorm:"column:fhir_id" json:"fhirId"`
	CreatedAt      time.Time  `json:"createdAt"`
}

type PractitionerCredential struct {
	ID             string     `gorm:"primaryKey;size:64" json:"id"`
	HospitalID     string     `gorm:"column:hospital_id;index:idx_cred_hosp_status" json:"hospitalId"`
	PractitionerID string     `gorm:"column:practitioner_id;index" json:"practitionerId"`
	Type           string     `gorm:"size:64" json:"type"`
	Name           string     `json:"name"`
	Number         *string    `json:"number"`
	IssueDate      *time.Time `json:"issueDate"`
	ExpiryDate     *time.Time `json:"expiryDate"`
	Status         string     `gorm:"size:32" json:"status"`
	Version        int        `json:"version"`
	VerifiedByID   *string    `gorm:"column:verified_by_id" json:"verifiedById"`
	VerifiedAt     *time.Time `json:"verifiedAt"`
	ApprovedByID   *string    `gorm:"column:approved_by_id" json:"approvedById"`
	ApprovedAt     *time.Time `json:"approvedAt"`
	Notes          *string    `json:"notes"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

type CredentialDocument struct {
	ID           string    `gorm:"primaryKey;size:64" json:"id"`
	HospitalID   string    `gorm:"column:hospital_id" json:"hospitalId"`
	CredentialID string    `gorm:"column:credential_id" json:"credentialId"`
	FileName     string    `json:"fileName"`
	FileType     string    `gorm:"size:8" json:"fileType"`
	SizeKb       int       `json:"sizeKb"`
	UploadedByID *string   `gorm:"column:uploaded_by_id" json:"uploadedById"`
	UploadedAt   time.Time `json:"uploadedAt"`
	Status       string    `gorm:"size:32" json:"status"`
	Notes        *string   `json:"notes"`
}

// ============ REGULASI & REQUIREMENT ENGINE ============

type Requirement struct {
	ID               string     `gorm:"primaryKey;size:64" json:"id"`
	HospitalID       string     `gorm:"column:hospital_id" json:"hospitalId"`
	Code             string     `gorm:"size:8" json:"code"`
	Name             string     `json:"name"`
	Category         string     `gorm:"size:64" json:"category"`
	RegulatorySource string     `json:"regulatorySource"`
	Version          string     `gorm:"size:16" json:"version"`
	Required         bool       `json:"required"`
	DocType          string     `json:"docType"`
	ValidationRule   *string    `json:"validationRule"`
	EvidenceRequired bool       `json:"evidenceRequired"`
	ApprovalRequired bool       `json:"approvalRequired"`
	ResponsibleRole  string     `json:"responsibleRole"`
	Status           string     `gorm:"size:32" json:"status"`
	DocNumber        *string    `json:"docNumber"`
	IssuedAt         *time.Time `json:"issuedAt"`
	DocExpiry        *time.Time `json:"docExpiry"`
	FileName         *string    `json:"fileName"`
	Notes            *string    `json:"notes"`
	UpdatedAt        time.Time  `json:"updatedAt"`
}

type Regulation struct {
	ID          string     `gorm:"primaryKey;size:64" json:"id"`
	Code        string     `gorm:"uniqueIndex;size:32" json:"code"`
	Title       string     `json:"title"`
	Number      string     `json:"number"`
	Authority   string     `json:"authority"`
	Year        int        `json:"year"`
	Category    string     `gorm:"size:48" json:"category"`
	Status      string     `gorm:"size:32" json:"status"`
	Summary     string     `json:"summary"`
	Source      string     `json:"source"`
	EffectiveAt *time.Time `json:"effectiveAt"`
}

type RegulationVersion struct {
	ID            string    `gorm:"primaryKey;size:64" json:"id"`
	RegulationID  string    `gorm:"column:regulation_id;index" json:"regulationId"`
	Version       string    `json:"version"`
	IssuedAt      time.Time `json:"issuedAt"`
	Status        string    `gorm:"size:32" json:"status"`
	ChangeSummary *string   `json:"changeSummary"`
	ImpactLevel   *string   `json:"impactLevel"`
	SupersededBy  *string   `json:"supersededBy"`
}

type ComplianceMapping struct {
	ID             string     `gorm:"primaryKey;size:64" json:"id"`
	RegulationID   string     `gorm:"column:regulation_id;index" json:"regulationId"`
	HospitalID     string     `gorm:"column:hospital_id" json:"hospitalId"`
	Article        string     `json:"article"`
	Requirement    string     `json:"requirement"`
	Module         string     `json:"module"`
	Feature        string     `json:"feature"`
	DataField      *string    `json:"dataField"`
	Validation     *string    `json:"validation"`
	EvidenceRef    *string    `json:"evidenceRef"`
	AuditRequired  bool       `json:"auditRequired"`
	Owner          string     `json:"owner"`
	Status         string     `gorm:"size:32" json:"status"`
	LastReviewedAt *time.Time `json:"lastReviewedAt"`
	Reviewer       *string    `json:"reviewer"`
}

// ============ SCORING ============

type ScoreSnapshot struct {
	ID             string    `gorm:"primaryKey;size:64" json:"id"`
	HospitalID     string    `gorm:"column:hospital_id;index" json:"hospitalId"`
	Period         string    `gorm:"size:8" json:"period"`
	HealthScore    float64   `json:"healthScore"`
	Bcds           float64   `json:"bcds"`
	Trust          float64   `json:"trust"`
	Ccs            float64   `json:"ccs"`
	Cpcs           float64   `json:"cpcs"`
	UrScore        float64   `json:"urScore"`
	SafetyScore    float64   `json:"safetyScore"`
	FormulaVersion string    `gorm:"size:48" json:"formulaVersion"`
	CalcAt         time.Time `json:"calcAt"`
	DataAt         time.Time `json:"dataAt"`
	DetailsJson    *string   `gorm:"column:details_json" json:"detailsJson"`
}

// ============ RISK & EARLY WARNING ============

type Risk struct {
	ID                string     `gorm:"primaryKey;size:64" json:"id"`
	HospitalID        string     `gorm:"column:hospital_id;index" json:"hospitalId"`
	Code              string     `gorm:"size:16" json:"code"`
	Category          string     `json:"category"`
	Level             string     `gorm:"size:16" json:"level"`
	Title             string     `json:"title"`
	Description       string     `json:"description"`
	DetectedAt        time.Time  `json:"detectedAt"`
	Domain            string     `json:"domain"`
	AffectedRecords   *string    `json:"affectedRecords"`
	EvidenceJson      *string    `gorm:"column:evidence_json" json:"evidenceJson"`
	BusinessImpact    string     `json:"businessImpact"`
	RecommendedAction string     `json:"recommendedAction"`
	SuggestedOwner    string     `json:"suggestedOwner"`
	Deadline          *time.Time `json:"deadline"`
	Status            string     `gorm:"size:32" json:"status"`
	CreatedAt         time.Time  `json:"createdAt"`
}

type EarlyWarning struct {
	ID             string     `gorm:"primaryKey;size:64" json:"id"`
	HospitalID     string     `gorm:"column:hospital_id" json:"hospitalId"`
	Code           string     `gorm:"size:16" json:"code"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Severity       string     `gorm:"size:16" json:"severity"`
	Domain         string     `json:"domain"`
	ThresholdRule  string     `json:"thresholdRule"`
	TriggerValue   string     `json:"triggerValue"`
	DetectedAt     time.Time  `json:"detectedAt"`
	Status         string     `gorm:"size:16" json:"status"`
	AcknowledgedBy *string    `json:"acknowledgedBy"`
	AcknowledgedAt *time.Time `json:"acknowledgedAt"`
	AssignedTo     *string    `json:"assignedTo"`
	DueDate        *time.Time `json:"dueDate"`
}

type AlertRule struct {
	ID              string     `gorm:"primaryKey;size:64" json:"id"`
	HospitalID      string     `gorm:"column:hospital_id" json:"hospitalId"`
	Name            string     `json:"name"`
	Domain          string     `json:"domain"`
	Metric          string     `json:"metric"`
	Condition       string     `gorm:"column:condition" json:"condition"`
	Severity        string     `gorm:"size:16" json:"severity"`
	Channels        string     `json:"channels"`
	Active          bool       `json:"active"`
	TriggeredCount  int        `json:"triggeredCount"`
	LastTriggeredAt *time.Time `json:"lastTriggeredAt"`
}

// ============ ACTION CENTER ============

type Action struct {
	ID              string     `gorm:"primaryKey;size:64" json:"id"`
	HospitalID      string     `gorm:"column:hospital_id" json:"hospitalId"`
	Code            string     `gorm:"size:16" json:"code"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	Priority        string     `gorm:"size:16" json:"priority"`
	Status          string     `gorm:"size:32" json:"status"`
	SourceType      string     `json:"sourceType"`
	SourceRef       *string    `json:"sourceRef"`
	Owner           string     `json:"owner"`
	OwnerRole       string     `json:"ownerRole"`
	DueDate         *time.Time `json:"dueDate"`
	ExpectedOutcome *string    `json:"expectedOutcome"`
	ProgressNote    *string    `json:"progressNote"`
	ResolutionNote  *string    `json:"resolutionNote"`
	ResolvedAt      *time.Time `json:"resolvedAt"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

type ActionEvent struct {
	ID        string    `gorm:"primaryKey;size:64" json:"id"`
	ActionID  string    `gorm:"column:action_id;index" json:"actionId"`
	Type      string    `gorm:"size:32" json:"type"`
	FromValue *string   `json:"fromValue"`
	ToValue   *string   `json:"toValue"`
	ByName    string    `json:"byName"`
	Note      *string   `json:"note"`
	CreatedAt time.Time `json:"createdAt"`
}

// ============ CLAIM INTELLIGENCE ============

type Claim struct {
	ID            string     `gorm:"primaryKey;size:64" json:"id"`
	HospitalID    string     `gorm:"column:hospital_id" json:"hospitalId"`
	ClaimNumber   string     `gorm:"index" json:"claimNumber"`
	SepNumber     string     `json:"sepNumber"`
	PatientInit   string     `json:"patientInit"`
	InacbgCode    string     `json:"inacbgCode"`
	InacbgDesc    string     `json:"inacbgDesc"`
	CareType      string     `json:"careType"`
	AdmissionDate time.Time  `json:"admissionDate"`
	DischargeDate *time.Time `json:"dischargeDate"`
	BilledAmount  int        `json:"billedAmount"`
	Status        string     `gorm:"size:32" json:"status"`
	FlagsJson     *string    `gorm:"column:flags_json" json:"flagsJson"`
	CreatedAt     time.Time  `json:"createdAt"`
}

type ClaimAnomaly struct {
	ID             string    `gorm:"primaryKey;size:64" json:"id"`
	ClaimID        string    `gorm:"column:claim_id;index" json:"claimId"`
	HospitalID     string    `gorm:"column:hospital_id" json:"hospitalId"`
	Category       string    `gorm:"size:32" json:"category"`
	Severity       string    `gorm:"size:16" json:"severity"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	EvidenceJson   *string   `gorm:"column:evidence_json" json:"evidenceJson"`
	Recommendation string    `json:"recommendation"`
	Status         string    `gorm:"size:32" json:"status"`
	DetectedAt     time.Time `json:"detectedAt"`
	ReviewedBy     *string   `json:"reviewedBy"`
}

// ============ CLINICAL INTELLIGENCE ============

type ClinicalPathway struct {
	ID                  string  `gorm:"primaryKey;size:64" json:"id"`
	HospitalID          string  `gorm:"column:hospital_id" json:"hospitalId"`
	Diagnosis           string  `json:"diagnosis"`
	Icd10               string  `gorm:"size:16" json:"icd10"`
	StandardLos         float64 `json:"standardLos"`
	TargetAdherence     float64 `json:"targetAdherence"`
	EpisodeCount        int     `json:"episodeCount"`
	AvgAdherence        float64 `json:"avgAdherence"`
	MedicationAdherence float64 `json:"medicationAdherence"`
	LabAdherence        float64 `json:"labAdherence"`
	ProcedureAdherence  float64 `json:"procedureAdherence"`
	LosAdherence        float64 `json:"losAdherence"`
	DeviationCount      int     `json:"deviationCount"`
	DocumentedDeviation float64 `json:"documentedDeviation"`
	CostRatio           float64 `json:"costRatio"`
	Status              string  `gorm:"size:16" json:"status"`
}

type PathwayEpisode struct {
	ID              string  `gorm:"primaryKey;size:64" json:"id"`
	PathwayID       string  `gorm:"column:pathway_id;index" json:"pathwayId"`
	HospitalID      string  `gorm:"column:hospital_id" json:"hospitalId"`
	EpisodeNumber   string  `json:"episodeNumber"`
	PatientInit     string  `json:"patientInit"`
	Dpjp            string  `json:"dpjp"`
	Adherence       float64 `json:"adherence"`
	LosActual       float64 `json:"losActual"`
	LosStandard     float64 `json:"losStandard"`
	DeviationType   string  `json:"deviationType"`
	DeviationReason *string `json:"deviationReason"`
	Documented      bool    `json:"documented"`
	CostRatio       float64 `json:"costRatio"`
	Period          string  `gorm:"size:8" json:"period"`
}

type UtilizationStat struct {
	ID          string  `gorm:"primaryKey;size:64" json:"id"`
	HospitalID  string  `gorm:"column:hospital_id" json:"hospitalId"`
	Month       string  `gorm:"size:8" json:"month"`
	MonthLabel  string  `gorm:"size:8" json:"monthLabel"`
	Alos        float64 `json:"alos"`
	Readmission float64 `json:"readmission"`
	Bor         float64 `json:"bor"`
	Score       float64 `json:"score"`
}

type ComplaintStat struct {
	ID           string  `gorm:"primaryKey;size:64" json:"id"`
	HospitalID   string  `gorm:"column:hospital_id" json:"hospitalId"`
	Month        string  `gorm:"size:8" json:"month"`
	Total        int     `json:"total"`
	Antrian      int     `json:"antrian"`
	Pelayanan    int     `json:"pelayanan"`
	Obat         int     `json:"obat"`
	Administrasi int     `json:"administrasi"`
	Fasilitas    int     `json:"fasilitas"`
	Lainnya      int     `json:"lainnya"`
	Satisfaction float64 `json:"satisfaction"`
	Per1000      float64 `json:"per1000"`
	PeerAvg      int     `json:"peerAvg"`
}

// ============ SARANA & DIGITAL ============

type Facility struct {
	ID           string  `gorm:"primaryKey;size:64" json:"id"`
	HospitalID   string  `gorm:"column:hospital_id" json:"hospitalId"`
	Name         string  `json:"name"`
	Category     string  `json:"category"`
	Availability string  `json:"availability"`
	Functional   string  `json:"functional"`
	Feasible     string  `json:"feasible"`
	Calibration  string  `json:"calibration"`
	Suitable     string  `json:"suitable"`
	StandardRef  string  `json:"standardRef"`
	Notes        *string `json:"notes"`
}

type BedStat struct {
	ID                 string  `gorm:"primaryKey;size:64" json:"id"`
	HospitalID         string  `gorm:"column:hospital_id" json:"hospitalId"`
	Ward               string  `json:"ward"`
	BedClass           string  `gorm:"column:bed_class" json:"bedClass"`
	Total              int     `json:"total"`
	Active             int     `json:"active"`
	Usable             int     `json:"usable"`
	NonOperatingReason *string `json:"nonOperatingReason"`
}

type SpecialtyStandard struct {
	ID         string  `gorm:"primaryKey;size:64" json:"id"`
	HospitalID string  `gorm:"column:hospital_id" json:"hospitalId"`
	RsClass    string  `gorm:"column:rs_class;size:4" json:"rsClass"`
	Specialty  string  `json:"specialty"`
	Required   int     `json:"required"`
	Note       *string `json:"note"`
}

type DigitalIndicator struct {
	ID           string    `gorm:"primaryKey;size:64" json:"id"`
	HospitalID   string    `gorm:"column:hospital_id" json:"hospitalId"`
	Dimension    string    `json:"dimension"`
	Indicator    string    `json:"indicator"`
	Value        float64   `json:"value"`
	Unit         string    `gorm:"size:8" json:"unit"`
	Threshold    float64   `json:"threshold"`
	Weight       float64   `json:"weight"`
	Status       string    `gorm:"size:32" json:"status"`
	FhirResource *string   `json:"fhirResource"`
	SourceSystem string    `json:"sourceSystem"`
	LastSyncAt   time.Time `json:"lastSyncAt"`
	Direction    string    `gorm:"size:8;default:up" json:"direction"`
}

// ============ DATA INTELLIGENCE ============

type DataSource struct {
	ID            string     `gorm:"primaryKey;size:64" json:"id"`
	HospitalID    string     `gorm:"column:hospital_id" json:"hospitalId"`
	Name          string     `json:"name"`
	Type          string     `json:"type"`
	Status        string     `json:"status"`
	AuthMethod    string     `json:"authMethod"`
	ConfigSummary *string    `json:"configSummary"`
	SuccessRate   float64    `json:"successRate"`
	AvgLatencyMs  int        `json:"avgLatencyMs"`
	LastSyncAt    *time.Time `json:"lastSyncAt"`
	FhirResources int        `json:"fhirResources"`
	CoveragePct   float64    `json:"coveragePct"`
	Uptime30d     float64    `json:"uptime30d"`
	Notes         *string    `json:"notes"`
}

type IntegrationJob struct {
	ID            string     `gorm:"primaryKey;size:64" json:"id"`
	HospitalID    string     `gorm:"column:hospital_id" json:"hospitalId"`
	DataSourceID  string     `gorm:"column:data_source_id" json:"dataSourceId"`
	JobType       string     `json:"jobType"`
	StartedAt     time.Time  `json:"startedAt"`
	FinishedAt    *time.Time `json:"finishedAt"`
	Status        string     `gorm:"size:16" json:"status"`
	RecordsPulled int        `json:"recordsPulled"`
	DurationMs    *int       `json:"durationMs"`
	Message       *string    `json:"message"`
}

type DataQualityIssue struct {
	ID              string     `gorm:"primaryKey;size:64" json:"id"`
	HospitalID      string     `gorm:"column:hospital_id" json:"hospitalId"`
	Dimension       string     `json:"dimension"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	AffectedRecords int        `json:"affectedRecords"`
	SourceSystem    string     `json:"sourceSystem"`
	Severity        string     `gorm:"size:16" json:"severity"`
	Anomaly         bool       `json:"anomaly"`
	Status          string     `gorm:"size:32" json:"status"`
	DetectedAt      time.Time  `json:"detectedAt"`
	ResolvedAt      *time.Time `json:"resolvedAt"`
}

type TrustFactor struct {
	ID         string  `gorm:"primaryKey;size:64" json:"id"`
	HospitalID string  `gorm:"column:hospital_id" json:"hospitalId"`
	Factor     string  `json:"factor"`
	Score      float64 `json:"score"`
	Weight     float64 `json:"weight"`
	Notes      *string `json:"notes"`
}

// ============ EVIDENCE & AUDIT ============

type EvidenceItem struct {
	ID               string     `gorm:"primaryKey;size:64" json:"id"`
	HospitalID       string     `gorm:"column:hospital_id;index" json:"hospitalId"`
	RefType          string     `gorm:"size:16" json:"refType"`
	RefID            string     `gorm:"column:ref_id" json:"refId"`
	Title            string     `json:"title"`
	Metric           *string    `json:"metric"`
	Value            *string    `json:"value"`
	Unit             *string    `json:"unit"`
	SourceSystem     string     `json:"sourceSystem"`
	FhirResource     *string    `json:"fhirResource"`
	RecordCount      *int       `json:"recordCount"`
	LastSyncAt       *time.Time `json:"lastSyncAt"`
	ValidationStatus string     `json:"validationStatus"`
	Timestamp        time.Time  `json:"timestamp"`
	DemoData         bool       `json:"demoData"`
}

type AuditLog struct {
	ID            string    `gorm:"primaryKey;size:64" json:"id"`
	HospitalID    *string   `gorm:"column:hospital_id;index:idx_audit_hosp_created" json:"hospitalId"`
	ActorName     string    `json:"actorName"`
	ActorRole     string    `json:"actorRole"`
	Action        string    `gorm:"size:24" json:"action"`
	EntityType    string    `json:"entityType"`
	EntityID      *string   `json:"entityId"`
	Description   string    `json:"description"`
	Before        *string   `json:"before"`
	After         *string   `json:"after"`
	IP            *string   `gorm:"column:ip" json:"ip"`
	CorrelationID *string   `json:"correlationId"`
	CreatedAt     time.Time `gorm:"index:idx_audit_hosp_created" json:"createdAt"`
}

// ============ REPORT & NOTIFIKASI ============

type ReportSnapshot struct {
	ID          string    `gorm:"primaryKey;size:64" json:"id"`
	HospitalID  string    `gorm:"column:hospital_id" json:"hospitalId"`
	Type        string    `gorm:"size:32" json:"type"`
	Title       string    `json:"title"`
	Period      string    `json:"period"`
	GeneratedAt time.Time `json:"generatedAt"`
	GeneratedBy string    `json:"generatedBy"`
	Methodology string    `json:"methodology"`
	Version     string    `gorm:"size:16" json:"version"`
	AuditRef    string    `json:"auditRef"`
	ContentJson string    `gorm:"column:content_json" json:"contentJson"`
}

type Notification struct {
	ID         string    `gorm:"primaryKey;size:64" json:"id"`
	HospitalID string    `gorm:"column:hospital_id" json:"hospitalId"`
	UserID     *string   `json:"userId"`
	Title      string    `json:"title"`
	Body       string    `json:"body"`
	Severity   string    `gorm:"size:16" json:"severity"`
	Channel    string    `gorm:"size:16" json:"channel"`
	Read       bool      `json:"read"`
	Link       *string   `json:"link"`
	CreatedAt  time.Time `json:"createdAt"`
}

// AllModels mengembalikan seluruh model untuk AutoMigrate.
func AllModels() []any {
	return []any{
		&Hospital{}, &User{}, &Practitioner{}, &PractitionerCredential{}, &CredentialDocument{},
		&Requirement{}, &Regulation{}, &RegulationVersion{}, &ComplianceMapping{},
		&ScoreSnapshot{}, &Risk{}, &EarlyWarning{}, &AlertRule{},
		&Action{}, &ActionEvent{}, &Claim{}, &ClaimAnomaly{},
		&ClinicalPathway{}, &PathwayEpisode{}, &UtilizationStat{}, &ComplaintStat{},
		&Facility{}, &BedStat{}, &SpecialtyStandard{}, &DigitalIndicator{},
		&DataSource{}, &IntegrationJob{}, &DataQualityIssue{}, &TrustFactor{},
		&EvidenceItem{}, &AuditLog{}, &ReportSnapshot{}, &Notification{},
	}
}
