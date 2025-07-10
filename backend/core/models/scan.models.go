package models

type Scan struct {
	Id              int    `json:"id"`
	RepositoryId    int    `json:"repository_id"`
	ScanTime        string `json:"scan_time"`
	Result          string `json:"result"`
	Vulnerabilities int    `json:"vulnerabilities" `
	Status          string `json:"status" `
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type ScanFull struct {
	Scan                  Scan   `json:"scan"`
	RepositoryName        string `json:"repository_name"`
	RepositoryPath        string `json:"repository_path"`
	RepositoryDescription string `json:"repository_description"`
}

type ScanCreate struct {
	RepositoryId    int    `json:"repository_id" validate:"required"`
	UserId          int    `json:"user_id" validate:"required"`
	ScanTime        string `json:"scan_time"`
	Result          string `json:"result" validate:"required"`
	Vulnerabilities int    `json:"vulnerabilities"`            // Assuming this is an integer count of vulnerabilities
	Status          string `json:"status" validate:"required"` // Assuming status is a string like "completed", "in_progress", etc.
}
