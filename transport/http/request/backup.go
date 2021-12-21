package request

// Backup request DTO.
type Backup struct {
	LicenseID string `json:"license_id"`
	FileName  string `json:"file_name"`
}
