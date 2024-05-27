package model

type CallTransfer struct {
	CampaignUuid string `json:"campaign_uuid"`
	UserUuid     string `json:"user_uuid"`
	PhoneNumber  string `json:"phone_number"`
	CallId       string `json:"call_id"`
	Token        string `json:"token"`
}
