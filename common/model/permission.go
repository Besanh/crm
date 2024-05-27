package model

type Permission struct {
	ModuleMain    ModuleMain `json:"permision_module_main"`
	DashboardMain ModuleMain `json:"dashboard_main"`

	// TicketMain      TicketChild `json:"ticket_main"`
	TicketMain         ModuleMain `json:"ticket_main"`
	TicketCategoryMain ModuleMain `json:"ticket_category_main"`

	ContactMain     ModuleMain `json:"contact_main"`
	FormMain        ModuleMain `json:"form_main"`
	MarketingMain   ModuleMain `json:"marketing_main"`
	QCMain          ModuleMain `json:"qc_main"`
	CampaignMain    ModuleMain `json:"campaign_main"`
	OmnichannelMain ModuleMain `json:"omnichannel_main"`
	UserMain        ModuleMain `json:"user_main"`
	ExtensionMain   ModuleMain `json:"extension_main"`
	SocialMain      ModuleMain `json:"social_main"`
	ChatMain        ModuleMain `json:"chat_main"`
	SmsZnsMain      ModuleMain `json:"sms_zns_main"`
	SettingChatMain ModuleMain `json:"setting_chat_main"`
	LogstashMain    ModuleMain `json:"logstash_main"`
	CdrMain         ModuleMain `json:"cdr_main"`
	ReportMain      ModuleMain `json:"report_main"`

	// ConfigMain      ConfigChild `json:"config_main"`
	RoleGroupMain ModuleMain `json:"role_group_main"`
	UnitTreeMain  ModuleMain `json:"unit_tree_main"`
	UnitMain      ModuleMain `json:"unit_main"`
	WorkDayMain   ModuleMain `json:"work_day_main"`
	EmailMain     ModuleMain `json:"email_main"`
	SolutionMain  ModuleMain `json:"solution_main"`
	QcSettingMain ModuleMain `json:"qc_config_main"`
	FileExport    ModuleMain `json:"file_export_main"`
	ClassifyMain  ModuleMain `json:"classify_main"`

	ToolMain    ModuleMain `json:"tool_main"`
	SupportMain ModuleMain `json:"support_main"`
	CarrierMain ModuleMain `json:"carrier_main"`
}

type PermissionView struct {
	ModuleMain         ModuleMainView `json:"permision_module_main"`
	DashboardMain      ModuleMainView `json:"dashboard_main"`
	TicketMain         ModuleMainView `json:"ticket_main"`
	TicketCategoryMain ModuleMainView `json:"ticket_category_main"`
	ContactMain        ModuleMainView `json:"contact_main"`
	FormMain           ModuleMainView `json:"form_main"`
	MarketingMain      ModuleMainView `json:"marketing_main"`
	QCMain             ModuleMainView `json:"qc_main"`
	CampaignMain       ModuleMainView `json:"campaign_main"`
	OmnichannelMain    ModuleMainView `json:"omnichannel_main"`
	UserMain           ModuleMainView `json:"user_main"`
	ExtensionMain      ModuleMain     `json:"extension_main"`
	SocialMain         ModuleMain     `json:"social_main"`
	ChatMain           ModuleMainView `json:"chat_main"`
	SettingChatMain    ModuleMainView `json:"setting_chat_main"`
	SmsZnsMain         ModuleMainView `json:"sms_zns_main"`
	LogstashMain       ModuleMainView `json:"logstash_main"`
	CdrMain            ModuleMainView `json:"cdr_main"`
	ReportMain         ModuleMainView `json:"report_main"`
	RoleGroupMain      ModuleMainView `json:"role_group_main"`
	UnitTreeMain       ModuleMainView `json:"unit_tree_main"`
	UnitMain           ModuleMainView `json:"unit_main"`
	WorkDayMain        ModuleMainView `json:"work_day_main"`
	EmailMain          ModuleMainView `json:"email_main"`
	SolutionMain       ModuleMainView `json:"solution_main"`
	QcSettingMain      ModuleMainView `json:"qc_config_main"`
	FileExportMain     ModuleMainView `json:"file_export_main"`
	ClassifyMain       ModuleMain     `json:"classify_main"`

	ToolMain    ModuleMainView `json:"tool_main"`
	SupportMain ModuleMainView `json:"support_main"`
	CarrierMain ModuleMainView `json:"carrier_main"`
}

type TicketChild struct {
	Status         ModuleDetail `json:"status"`
	Ticket         ModuleMain   `json:"ticket"`
	TicketCategory ModuleMain   `json:"ticket_category"`
}

type ConfigChild struct {
	Status    ModuleDetail `json:"status"`
	RoleGroup ModuleMain   `json:"role_group"`
	UnitTree  ModuleMain   `json:"unit_tree"`
	Unit      ModuleMain   `json:"unit"`
	WorkDay   ModuleMain   `json:"work_day"`
	Email     ModuleMain   `json:"email"`
	Solution  ModuleMain   `json:"solution"`
	Qc        ModuleMain   `json:"qc"`
}

type ModuleDetail struct {
	Value    bool `json:"value"`
	Disabled bool `json:"disabled"`
}

type ModuleMain struct {
	Status ModuleDetail `json:"status"`
	Create ModuleDetail `json:"create"`
	View   ModuleDetail `json:"view"`
	Edit   ModuleDetail `json:"edit"`
	Delete ModuleDetail `json:"delete"`
	Import ModuleDetail `json:"import"`
	Export ModuleDetail `json:"export"`
	Search ModuleDetail `json:"search"`
	Assign ModuleDetail `json:"assign"`
}

type ModuleMainView struct {
	Status bool `json:"status"`
	Create bool `json:"create"`
	View   bool `json:"view"`
	Edit   bool `json:"edit"`
	Delete bool `json:"delete"`
	Import bool `json:"import"`
	Export bool `json:"export"`
	Search bool `json:"search"`
	Assign bool `json:"assign"`
}

type PermissionAdvance struct {
	InfoEnterprise ModuleMain `json:"info_enterprise"`
	AccountBalance ModuleMain `json:"account_balance"`
	ServicePack    ModuleMain `json:"service_pack"`
	PaymentHistory ModuleMain `json:"payment_history"`
	LoginConfig    ModuleMain `json:"login_config"`
	Pbx            ModuleMain `json:"pbx"`
}

type PermissionAdvanceView struct {
	InfoEnterprise ModuleMainView `json:"info_enterprise"`
	AccountBalance ModuleMainView `json:"account_balance"`
	ServicePack    ModuleMainView `json:"service_pack"`
	PaymentHistory ModuleMainView `json:"payment_history"`
	LoginConfig    ModuleMainView `json:"login_config"`
	Pbx            ModuleMainView `json:"pbx"`
}

type PermissionUser struct {
	UserUuid string         `json:"user_uuid"`
	Module   PermissionView `json:"module"`
}
