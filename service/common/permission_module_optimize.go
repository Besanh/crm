package common

import "contactcenter-api/common/model"

func HandlePermissionMainOptimize(data model.Permission) model.PermissionMainOptimize {
	permissionMainOptimize := model.PermissionMainOptimize{
		ModuleMain: model.ModuleMainOptimze{
			Create: data.ModuleMain.Create.Value,
			View:   data.ModuleMain.View.Value,
			Edit:   data.ModuleMain.Edit.Value,
			Delete: data.ModuleMain.Delete.Value,
			Import: data.ModuleMain.Import.Value,
			Export: data.ModuleMain.Export.Value,
			Search: data.ModuleMain.Search.Value,
			Assign: data.ModuleMain.Assign.Value,
		},
		DashboardMain: model.ModuleMainOptimze{
			Create: data.DashboardMain.Create.Value,
			View:   data.DashboardMain.View.Value,
			Edit:   data.DashboardMain.Edit.Value,
			Delete: data.DashboardMain.Delete.Value,
			Import: data.DashboardMain.Import.Value,
			Export: data.DashboardMain.Export.Value,
			Search: data.DashboardMain.Search.Value,
			Assign: data.DashboardMain.Assign.Value,
		},
		TicketMain: model.ModuleMainOptimze{
			Create: data.TicketMain.Create.Value,
			View:   data.TicketMain.View.Value,
			Edit:   data.TicketMain.Edit.Value,
			Delete: data.TicketMain.Delete.Value,
			Import: data.TicketMain.Import.Value,
			Export: data.TicketMain.Export.Value,
			Search: data.TicketMain.Search.Value,
			Assign: data.TicketMain.Assign.Value,
		},
		TicketCategoryMain: model.ModuleMainOptimze{
			Create: data.TicketCategoryMain.Create.Value,
			View:   data.TicketCategoryMain.View.Value,
			Edit:   data.TicketCategoryMain.Edit.Value,
			Delete: data.TicketCategoryMain.Delete.Value,
			Import: data.TicketCategoryMain.Import.Value,
			Export: data.TicketCategoryMain.Export.Value,
			Search: data.TicketCategoryMain.Search.Value,
			Assign: data.TicketCategoryMain.Assign.Value,
		},
		ContactMain: model.ModuleMainOptimze{
			Create: data.ContactMain.Create.Value,
			View:   data.ContactMain.View.Value,
			Edit:   data.ContactMain.Edit.Value,
			Delete: data.ContactMain.Delete.Value,
			Import: data.ContactMain.Import.Value,
			Export: data.ContactMain.Export.Value,
			Search: data.ContactMain.Search.Value,
			Assign: data.ContactMain.Assign.Value,
		},
		FormMain: model.ModuleMainOptimze{
			Create: data.FormMain.Create.Value,
			View:   data.FormMain.View.Value,
			Edit:   data.FormMain.Edit.Value,
			Delete: data.FormMain.Delete.Value,
			Import: data.FormMain.Import.Value,
			Export: data.FormMain.Export.Value,
			Search: data.FormMain.Search.Value,
			Assign: data.FormMain.Assign.Value,
		},
		MarketingMain: model.ModuleMainOptimze{
			Create: data.MarketingMain.Create.Value,
			View:   data.MarketingMain.View.Value,
			Edit:   data.MarketingMain.Edit.Value,
			Delete: data.MarketingMain.Delete.Value,
			Import: data.MarketingMain.Import.Value,
			Export: data.MarketingMain.Export.Value,
			Search: data.MarketingMain.Search.Value,
			Assign: data.MarketingMain.Assign.Value,
		},
		QCMain: model.ModuleMainOptimze{
			Create: data.QCMain.Create.Value,
			View:   data.QCMain.View.Value,
			Edit:   data.QCMain.Edit.Value,
			Delete: data.QCMain.Delete.Value,
			Import: data.QCMain.Import.Value,
			Export: data.QCMain.Export.Value,
			Search: data.QCMain.Search.Value,
			Assign: data.QCMain.Assign.Value,
		},
		CampaignMain: model.ModuleMainOptimze{
			Create: data.CampaignMain.Create.Value,
			View:   data.CampaignMain.View.Value,
			Edit:   data.CampaignMain.Edit.Value,
			Delete: data.CampaignMain.Delete.Value,
			Import: data.CampaignMain.Import.Value,
			Export: data.CampaignMain.Export.Value,
			Search: data.CampaignMain.Search.Value,
			Assign: data.CampaignMain.Assign.Value,
		},
		OmnichannelMain: model.ModuleMainOptimze{
			Create: data.OmnichannelMain.Create.Value,
			View:   data.OmnichannelMain.View.Value,
			Edit:   data.OmnichannelMain.Edit.Value,
			Delete: data.OmnichannelMain.Delete.Value,
			Import: data.OmnichannelMain.Import.Value,
			Export: data.OmnichannelMain.Export.Value,
			Search: data.OmnichannelMain.Search.Value,
			Assign: data.OmnichannelMain.Assign.Value,
		},
		UserMain: model.ModuleMainOptimze{
			Create: data.UserMain.Create.Value,
			View:   data.UserMain.View.Value,
			Edit:   data.UserMain.Edit.Value,
			Delete: data.UserMain.Delete.Value,
			Import: data.UserMain.Import.Value,
			Export: data.UserMain.Export.Value,
			Search: data.UserMain.Search.Value,
			Assign: data.UserMain.Assign.Value,
		},
		ExtensionMain: model.ModuleMainOptimze{
			Create: data.ExtensionMain.Create.Value,
			View:   data.ExtensionMain.View.Value,
			Edit:   data.ExtensionMain.Edit.Value,
			Delete: data.ExtensionMain.Delete.Value,
			Import: data.ExtensionMain.Import.Value,
			Export: data.ExtensionMain.Export.Value,
			Search: data.ExtensionMain.Search.Value,
			Assign: data.ExtensionMain.Assign.Value,
		},
		SocialMain: model.ModuleMainOptimze{
			Create: data.SocialMain.Create.Value,
			View:   data.SocialMain.View.Value,
			Edit:   data.SocialMain.Edit.Value,
			Delete: data.SocialMain.Delete.Value,
			Import: data.SocialMain.Import.Value,
			Export: data.SocialMain.Export.Value,
			Search: data.SocialMain.Search.Value,
			Assign: data.SocialMain.Assign.Value,
		},
		ChatMain: model.ModuleMainOptimze{
			Create: data.ChatMain.Create.Value,
			View:   data.ChatMain.View.Value,
			Edit:   data.ChatMain.Edit.Value,
			Delete: data.ChatMain.Delete.Value,
			Import: data.ChatMain.Import.Value,
			Export: data.ChatMain.Export.Value,
			Search: data.ChatMain.Search.Value,
			Assign: data.ChatMain.Assign.Value,
		},
		SettingChatMain: model.ModuleMainOptimze{
			Create: data.SettingChatMain.Create.Value,
			View:   data.SettingChatMain.View.Value,
			Edit:   data.SettingChatMain.Edit.Value,
			Delete: data.SettingChatMain.Delete.Value,
			Import: data.SettingChatMain.Import.Value,
			Export: data.SettingChatMain.Export.Value,
			Search: data.SettingChatMain.Search.Value,
			Assign: data.SettingChatMain.Assign.Value,
		},
		SmsZnsMain: model.ModuleMainOptimze{
			Create: data.SmsZnsMain.Create.Value,
			View:   data.SmsZnsMain.View.Value,
			Edit:   data.SmsZnsMain.Edit.Value,
			Delete: data.SmsZnsMain.Delete.Value,
			Import: data.SmsZnsMain.Import.Value,
			Export: data.SmsZnsMain.Export.Value,
			Search: data.SmsZnsMain.Search.Value,
			Assign: data.SmsZnsMain.Assign.Value,
		},
		LogstashMain: model.ModuleMainOptimze{
			Create: data.LogstashMain.Create.Value,
			View:   data.LogstashMain.View.Value,
			Edit:   data.LogstashMain.Edit.Value,
			Delete: data.LogstashMain.Delete.Value,
			Import: data.LogstashMain.Import.Value,
			Export: data.LogstashMain.Export.Value,
			Search: data.LogstashMain.Search.Value,
			Assign: data.LogstashMain.Assign.Value,
		},
		CdrMain: model.ModuleMainOptimze{
			Create: data.CdrMain.Create.Value,
			View:   data.CdrMain.View.Value,
			Edit:   data.CdrMain.Edit.Value,
			Delete: data.CdrMain.Delete.Value,
			Import: data.CdrMain.Import.Value,
			Export: data.CdrMain.Export.Value,
			Search: data.CdrMain.Search.Value,
			Assign: data.CdrMain.Assign.Value,
		},
		ReportMain: model.ModuleMainOptimze{
			Create: data.ReportMain.Create.Value,
			View:   data.ReportMain.View.Value,
			Edit:   data.ReportMain.Edit.Value,
			Delete: data.ReportMain.Delete.Value,
			Import: data.ReportMain.Import.Value,
			Export: data.ReportMain.Export.Value,
			Search: data.ReportMain.Search.Value,
			Assign: data.ReportMain.Assign.Value,
		},
		RoleGroupMain: model.ModuleMainOptimze{
			Create: data.RoleGroupMain.Create.Value,
			View:   data.RoleGroupMain.View.Value,
			Edit:   data.RoleGroupMain.Edit.Value,
			Delete: data.RoleGroupMain.Delete.Value,
			Import: data.RoleGroupMain.Import.Value,
			Export: data.RoleGroupMain.Export.Value,
			Search: data.RoleGroupMain.Search.Value,
			Assign: data.RoleGroupMain.Assign.Value,
		},
		UnitTreeMain: model.ModuleMainOptimze{
			Create: data.UnitTreeMain.Create.Value,
			View:   data.UnitTreeMain.View.Value,
			Edit:   data.UnitTreeMain.Edit.Value,
			Delete: data.UnitTreeMain.Delete.Value,
			Import: data.UnitTreeMain.Import.Value,
			Export: data.UnitTreeMain.Export.Value,
			Search: data.UnitTreeMain.Search.Value,
			Assign: data.UnitTreeMain.Assign.Value,
		},
		UnitMain: model.ModuleMainOptimze{
			Create: data.UnitMain.Create.Value,
			View:   data.UnitMain.View.Value,
			Edit:   data.UnitMain.Edit.Value,
			Delete: data.UnitMain.Delete.Value,
			Import: data.UnitMain.Import.Value,
			Export: data.UnitMain.Export.Value,
			Search: data.UnitMain.Search.Value,
			Assign: data.UnitMain.Assign.Value,
		},
		WorkDayMain: model.ModuleMainOptimze{
			Create: data.WorkDayMain.Create.Value,
			View:   data.WorkDayMain.View.Value,
			Edit:   data.WorkDayMain.Edit.Value,
			Delete: data.WorkDayMain.Delete.Value,
			Import: data.WorkDayMain.Import.Value,
			Export: data.WorkDayMain.Export.Value,
			Search: data.WorkDayMain.Search.Value,
			Assign: data.WorkDayMain.Assign.Value,
		},
		EmailMain: model.ModuleMainOptimze{
			Create: data.EmailMain.Create.Value,
			View:   data.EmailMain.View.Value,
			Edit:   data.EmailMain.Edit.Value,
			Delete: data.EmailMain.Delete.Value,
			Import: data.EmailMain.Import.Value,
			Export: data.EmailMain.Export.Value,
			Search: data.EmailMain.Search.Value,
			Assign: data.EmailMain.Assign.Value,
		},
		SolutionMain: model.ModuleMainOptimze{
			Create: data.SolutionMain.Create.Value,
			View:   data.SolutionMain.View.Value,
			Edit:   data.SolutionMain.Edit.Value,
			Delete: data.SolutionMain.Delete.Value,
			Import: data.SolutionMain.Import.Value,
			Export: data.SolutionMain.Export.Value,
			Search: data.SolutionMain.Search.Value,
			Assign: data.SolutionMain.Assign.Value,
		},
		QcSettingMain: model.ModuleMainOptimze{
			Create: data.QcSettingMain.Create.Value,
			View:   data.QcSettingMain.View.Value,
			Edit:   data.QcSettingMain.Edit.Value,
			Delete: data.QcSettingMain.Delete.Value,
			Import: data.QcSettingMain.Import.Value,
			Export: data.QcSettingMain.Export.Value,
			Search: data.QcSettingMain.Search.Value,
			Assign: data.QcSettingMain.Assign.Value,
		},
		FileExport: model.ModuleMainOptimze{
			Create: data.FileExport.Create.Value,
			View:   data.FileExport.View.Value,
			Edit:   data.FileExport.Edit.Value,
			Delete: data.FileExport.Delete.Value,
			Import: data.FileExport.Import.Value,
			Export: data.FileExport.Export.Value,
			Search: data.FileExport.Search.Value,
			Assign: data.FileExport.Assign.Value,
		},
		ClassifyMain: model.ModuleMainOptimze{
			Create: data.ClassifyMain.Create.Value,
			View:   data.ClassifyMain.View.Value,
			Edit:   data.ClassifyMain.Edit.Value,
			Delete: data.ClassifyMain.Delete.Value,
			Import: data.ClassifyMain.Import.Value,
			Export: data.ClassifyMain.Export.Value,
			Search: data.ClassifyMain.Search.Value,
			Assign: data.ClassifyMain.Assign.Value,
		},
		ToolMain: model.ModuleMainOptimze{
			Create: data.ToolMain.Create.Value,
			View:   data.ToolMain.View.Value,
			Edit:   data.ToolMain.Edit.Value,
			Delete: data.ToolMain.Delete.Value,
			Import: data.ToolMain.Import.Value,
			Export: data.ToolMain.Export.Value,
			Search: data.ToolMain.Search.Value,
			Assign: data.ToolMain.Assign.Value,
		},
		SupportMain: model.ModuleMainOptimze{
			Create: data.SupportMain.Create.Value,
			View:   data.SupportMain.View.Value,
			Edit:   data.SupportMain.Edit.Value,
			Delete: data.SupportMain.Delete.Value,
			Import: data.SupportMain.Import.Value,
			Export: data.SupportMain.Export.Value,
			Search: data.SupportMain.Search.Value,
			Assign: data.SupportMain.Assign.Value,
		},
		CarrierMain: model.ModuleMainOptimze{
			Create: data.CarrierMain.Create.Value,
			View:   data.CarrierMain.View.Value,
			Edit:   data.CarrierMain.Edit.Value,
			Delete: data.CarrierMain.Delete.Value,
			Import: data.CarrierMain.Import.Value,
			Export: data.CarrierMain.Export.Value,
			Search: data.CarrierMain.Search.Value,
			Assign: data.CarrierMain.Assign.Value,
		},
	}

	return permissionMainOptimize
}

func HandlePermissionAdvanceOptimize(data model.PermissionAdvance) model.PermissionAdvanceOptimize {
	permissionAdvanceOptimize := model.PermissionAdvanceOptimize{
		InfoEnterprise: model.ModuleMainOptimze{
			Create: data.InfoEnterprise.Create.Value,
			View:   data.InfoEnterprise.View.Value,
			Edit:   data.InfoEnterprise.Edit.Value,
			Delete: data.InfoEnterprise.Delete.Value,
			Import: data.InfoEnterprise.Import.Value,
			Export: data.InfoEnterprise.Export.Value,
			Search: data.InfoEnterprise.Search.Value,
			Assign: data.InfoEnterprise.Assign.Value,
		},
		AccountBalance: model.ModuleMainOptimze{
			Create: data.AccountBalance.Create.Value,
			View:   data.AccountBalance.View.Value,
			Edit:   data.AccountBalance.Edit.Value,
			Delete: data.AccountBalance.Delete.Value,
			Import: data.AccountBalance.Import.Value,
			Export: data.AccountBalance.Export.Value,
			Search: data.AccountBalance.Search.Value,
			Assign: data.AccountBalance.Assign.Value,
		},
		ServicePack: model.ModuleMainOptimze{
			Create: data.ServicePack.Edit.Value,
			Delete: data.ServicePack.Delete.Value,
			Import: data.ServicePack.Import.Value,
			Export: data.ServicePack.Export.Value,
			Search: data.ServicePack.Search.Value,
			Assign: data.ServicePack.Assign.Value,
		},
		PaymentHistory: model.ModuleMainOptimze{
			Create: data.PaymentHistory.Edit.Value,
			Delete: data.PaymentHistory.Delete.Value,
			Import: data.PaymentHistory.Import.Value,
			Export: data.PaymentHistory.Export.Value,
			Search: data.PaymentHistory.Search.Value,
			Assign: data.PaymentHistory.Assign.Value,
		},
		LoginConfig: model.ModuleMainOptimze{
			Create: data.LoginConfig.Edit.Value,
			Delete: data.LoginConfig.Delete.Value,
			Import: data.LoginConfig.Import.Value,
			Export: data.LoginConfig.Export.Value,
			Search: data.LoginConfig.Search.Value,
			Assign: data.LoginConfig.Assign.Value,
		},
		Pbx: model.ModuleMainOptimze{
			Create: data.Pbx.Edit.Value,
			View:   data.Pbx.View.Value,
			Edit:   data.Pbx.Edit.Value,
			Delete: data.Pbx.Delete.Value,
			Import: data.Pbx.Import.Value,
			Export: data.Pbx.Export.Value,
			Search: data.Pbx.Search.Value,
			Assign: data.Pbx.Assign.Value,
		},
	}

	return permissionAdvanceOptimize
}