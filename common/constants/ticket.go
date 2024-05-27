package constants

const (
	TICKET_STATUS_OPEN       string = "OPEN"
	TICKET_STATUS_PROCESSING string = "PROCESSING" //when assign to department
	TICKET_STATUS_WAITING    string = "WAITING"
	TICKET_STATUS_SOLVED     string = "SOLVED"
	TICKET_STATUS_PENDING    string = "PENDING"
	TICKET_STATUS_REOPEN     string = "REOPEN"

	NOTIFICATION_TICKET_USER         string = "NOTIFICATION_TICKET_USER"
	NOTIFICATION_TICKET_USER_HOUR    string = "NOTIFICATION_TICKET_USER_HOUR"
	NOTIFICATION_TICKET_USER_EXPIRED string = "NOTIFICATION_TICKET_USER_EXPIRED"
	exportListKey                    string = "export_"
	TOTAL_SLA_ONTIME                 string = "total_sla_ontime"
	TOTAL_TICKET                     string = "total_ticket"
	TICKET_CHAT                      string = "ticket_chat"
	TICKET_EMAIL                     string = "ticket_email"
)
