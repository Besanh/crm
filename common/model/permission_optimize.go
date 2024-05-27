package model

type PermissionMainOptimize struct {
	ModuleMain    ModuleMainOptimze `json:"permision_module_main"`
	DashboardMain ModuleMainOptimze `json:"dashboard_main"`

	TicketMain         ModuleMainOptimze `json:"ticket_main"`
	TicketCategoryMain ModuleMainOptimze `json:"ticket_category_main"`

	ContactMain     ModuleMainOptimze `json:"contact_main"`
	FormMain        ModuleMainOptimze `json:"form_main"`
	MarketingMain   ModuleMainOptimze `json:"marketing_main"`
	QCMain          ModuleMainOptimze `json:"qc_main"`
	CampaignMain    ModuleMainOptimze `json:"campaign_main"`
	OmnichannelMain ModuleMainOptimze `json:"omnichannel_main"`
	UserMain        ModuleMainOptimze `json:"user_main"`
	ExtensionMain   ModuleMainOptimze `json:"extension_main"`
	SocialMain      ModuleMainOptimze `json:"social_main"`
	ChatMain        ModuleMainOptimze `json:"chat_main"`
	SettingChatMain ModuleMainOptimze `json:"setting_chat_main"`
	SmsZnsMain      ModuleMainOptimze `json:"sms_zns_main"`
	LogstashMain    ModuleMainOptimze `json:"logstash_main"`
	CdrMain         ModuleMainOptimze `json:"cdr_main"`
	ReportMain      ModuleMainOptimze `json:"report_main"`

	RoleGroupMain ModuleMainOptimze `json:"role_group_main"`
	UnitTreeMain  ModuleMainOptimze `json:"unit_tree_main"`
	UnitMain      ModuleMainOptimze `json:"unit_main"`
	WorkDayMain   ModuleMainOptimze `json:"work_day_main"`
	EmailMain     ModuleMainOptimze `json:"email_main"`
	SolutionMain  ModuleMainOptimze `json:"solution_main"`
	QcSettingMain ModuleMainOptimze `json:"qc_config_main"`
	FileExport    ModuleMainOptimze `json:"file_export_main"`
	ClassifyMain  ModuleMainOptimze `json:"classify_main"`

	ToolMain    ModuleMainOptimze `json:"tool_main"`
	SupportMain ModuleMainOptimze `json:"support_main"`
	CarrierMain ModuleMainOptimze `json:"carrier_main"`
}

type PermissionAdvanceOptimize struct {
	InfoEnterprise ModuleMainOptimze `json:"info_enterprise"`
	AccountBalance ModuleMainOptimze `json:"account_balance"`
	ServicePack    ModuleMainOptimze `json:"service_pack"`
	PaymentHistory ModuleMainOptimze `json:"payment_history"`
	LoginConfig    ModuleMainOptimze `json:"login_config"`
	Pbx            ModuleMainOptimze `json:"pbx"`
}

type ModuleMainOptimze struct {
	Create bool `json:"create"`
	View   bool `json:"view"`
	Edit   bool `json:"edit"`
	Delete bool `json:"delete"`
	Import bool `json:"import"`
	Export bool `json:"export"`
	Search bool `json:"search"`
	Assign bool `json:"assign"`
}
