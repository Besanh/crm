package model

import "github.com/uptrace/bun"

type VContactEmail struct {
	bun.BaseModel    `bun:"v_contact_emails,alias:ce"`
	ContactEmailUuid string `json:"contact_email_uuid" bun:"contact_email_uuid"`
	ContactUuid      string `json:"contact_uuid"  bun:"contact_uuid"`
	DomainUuid       string `json:"domain_uuid"   bun:"domain_uuid"`
	EmailAddress     string `json:"email_address" bun:"email_address"`
	EmailPrimary     int32  `json:"email_primary" bun:"email_primary"`
}
