package model

type ExportMap struct {
	Name             string `json:"name"`
	ExportTime       string `json:"export_time"`
	ExportTimeFinish string `json:"export_time_finish"`
	TotalRows        int    `json:"total_rows"`
	Status           string `json:"status"`
	DomainUuid       string `json:"domain_uuid"`
	UserUuid         string `json:"user_uuid"`
	Type             string `json:"type"`
	Folder           string `json:"folder"`
}