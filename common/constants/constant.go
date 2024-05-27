package constants

import "time"

const (
	SERVICE_NAME string = "contactcenter-api"
	VERSION_NAME string = "v1.0"
	VERSION      string = "v1"

	// Priority
	LOW    string = "LOW"
	MEDIUM string = "MEDIUM"
	HIGH   string = "HIGH"

	// Channel
	WEB      string = "WEB"
	CHAT     string = "CHAT"
	MSMAIL   string = "MSMAIL"
	GMAIL    string = "GMAIL"
	SMS      string = "SMS"
	ZNS      string = "ZNS"
	AUTOCALL string = "AUTOCALL"

	// KPI
	RATIO_AFTER_CALL string = "RATIO_AFTER_CALL"
	COUNT_AFTER_CALL string = "COUNT_AFTER_CALL"
	KPI_SCORE        string = "KPI_SCORE"
	KPI_COUNT        string = "KPI_COUNT"
	AGENT_SCORE_LATE string = "AGENT_SCORE_LATE"
	AGENT_COUNT_LATE string = "AGENT_COUNT_LATE"

	// Notify
	NOTIFY_SYSTEM_USER         string = "NOTIFY_SYSTEM_USER"
	NOTIFIY_TICKET_USER_ASSIGN string = "NOTIFIY_TICKET_USER_ASSIGN"

	// Check duplicate
	DUPLICATE_CAMPAIGN = "dupcamp"
	DUPLICATE_SYSTEM   = "dupsys"

	// Log level
	ERROR   = "error"
	INFO    = "info"
	DEBUG   = "debug"
	WARNING = "warning"

	// Level
	ADMIN      = "admin"
	SUPERADMIN = "superadmin"
	LEADER     = "leader"
	MANAGER    = "manager"
	USER       = "user"
	AGENT      = "agent"

	// File
	FILE_SIZE_THRESHOLD = 512000

	// Upload contact
	MAX_CONTACT = 2000000
	DELIMITER   = ","

	// Omni
	OMNI_INFO = "OMNI_INFO"
	TTL_OMNI  = 1 * time.Minute

	// Export
	EXPORT_KEY = "export_crm_"
	EXPORT_DIR = "/root/go/src/exported/"
	SHEET1     = "Sheet1"

	ROLE_ADMIN   = "admin"
	ROLE_MANAGER = "manager"
	ROLE_USER    = "user"
	ROLE_LEADER  = "leader"

	PRIVILEGE_USER_UNIT = "privilege_user_unit_"
	CALL_TRANSFER       = "call_transfer"
)
