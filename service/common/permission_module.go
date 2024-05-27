package common

import (
	"contactcenter-api/common/model"
)

func HandleCreatePermissionMain() model.Permission {
	defaultData := model.ModuleMain{
		Status: model.ModuleDetail{
			Value:    true,
			Disabled: false,
		},
		Create: model.ModuleDetail{
			Value:    false,
			Disabled: false,
		},
		View: model.ModuleDetail{
			Value:    false,
			Disabled: false,
		},
		Edit: model.ModuleDetail{
			Value:    false,
			Disabled: false,
		},
		Delete: model.ModuleDetail{
			Value:    false,
			Disabled: false,
		},
		Import: model.ModuleDetail{
			Value:    false,
			Disabled: false,
		},
		Export: model.ModuleDetail{
			Value:    false,
			Disabled: false,
		},
		Search: model.ModuleDetail{
			Value:    false,
			Disabled: false,
		},
		Assign: model.ModuleDetail{
			Value:    false,
			Disabled: false,
		},
	}
	dashboard := defaultData
	permission := model.Permission{
		ModuleMain: defaultData,
		DashboardMain: model.ModuleMain{
			Status: model.ModuleDetail{
				Value:    true,
				Disabled: false,
			},
			Create: dashboard.Create,
			View: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Edit:   dashboard.Edit,
			Delete: dashboard.Delete,
			Import: dashboard.Import,
			Export: dashboard.Export,
			Assign: dashboard.Assign,
		},
		TicketMain:         defaultData,
		TicketCategoryMain: defaultData,
		ContactMain:        defaultData,
		MarketingMain:      defaultData,
		QCMain: model.ModuleMain{
			Status: model.ModuleDetail{
				Value:    true,
				Disabled: false,
			},
			Create: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			View: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Edit: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Delete: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Import: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Export: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Search: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Assign: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
		},
		CampaignMain:    defaultData,
		OmnichannelMain: defaultData,
		UserMain:        defaultData,
		ExtensionMain:   defaultData,
		SocialMain:      defaultData,
		ChatMain: model.ModuleMain{
			Status: model.ModuleDetail{
				Value:    true,
				Disabled: false,
			},
			View: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Edit: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Delete: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Import: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Export: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Search: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Assign: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
		},
		SettingChatMain: model.ModuleMain{
			Status: model.ModuleDetail{
				Value:    true,
				Disabled: false,
			},
			Create: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			View: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Edit: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Delete: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Import: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Export: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Search: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Assign: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
		},
		SmsZnsMain: defaultData,

		LogstashMain: model.ModuleMain{
			Status: model.ModuleDetail{
				Value:    true,
				Disabled: false,
			},
			Create: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			View: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Edit: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Delete: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Import: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Export: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Search: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Assign: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
		},
		CdrMain: model.ModuleMain{
			Status: model.ModuleDetail{
				Value:    true,
				Disabled: false,
			},
			Create: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			View: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Edit: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Delete: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Import: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Export: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Search: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Assign: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
		},
		ReportMain: model.ModuleMain{
			Status: model.ModuleDetail{
				Value:    true,
				Disabled: false,
			},
			Create: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			View: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Edit: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Delete: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Import: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Export: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Search: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Assign: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
		},
		RoleGroupMain: defaultData,
		UnitTreeMain:  defaultData,
		UnitMain:      defaultData,
		WorkDayMain: model.ModuleMain{
			Status: model.ModuleDetail{
				Value:    true,
				Disabled: false,
			},
			Create: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			View: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Edit: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Delete: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Import: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Export: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Search: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Assign: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
		},
		EmailMain:     defaultData,
		SolutionMain:  defaultData,
		QcSettingMain: defaultData,
		FileExport: model.ModuleMain{
			Status: model.ModuleDetail{
				Value:    true,
				Disabled: false,
			},
			Create: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			View: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Edit: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Delete: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Import: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Export: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Search: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Assign: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
		},
		ClassifyMain: defaultData,
		ToolMain: model.ModuleMain{
			Status: model.ModuleDetail{
				Value:    true,
				Disabled: false,
			},
			Create: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			View: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Edit: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Delete: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Import: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Export: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Search: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Assign: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
		},
		SupportMain: model.ModuleMain{
			Status: model.ModuleDetail{
				Value:    true,
				Disabled: false,
			},
			Create: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			View: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Edit: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Delete: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Import: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Export: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Search: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Assign: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
		},
		CarrierMain: defaultData,
	}

	return permission
}

func HandleCreatePermissionAdvance() model.PermissionAdvance {
	defaultData := model.ModuleMain{
		Status: model.ModuleDetail{
			Value:    true,
			Disabled: false,
		},
		Create: model.ModuleDetail{
			Value:    false,
			Disabled: false,
		},
		View: model.ModuleDetail{
			Value:    false,
			Disabled: false,
		},
		Edit: model.ModuleDetail{
			Value:    false,
			Disabled: false,
		},
		Delete: model.ModuleDetail{
			Value:    false,
			Disabled: false,
		},
		Import: model.ModuleDetail{
			Value:    false,
			Disabled: false,
		},
		Export: model.ModuleDetail{
			Value:    false,
			Disabled: false,
		},
		Search: model.ModuleDetail{
			Value:    false,
			Disabled: false,
		},
		Assign: model.ModuleDetail{
			Value:    false,
			Disabled: false,
		},
	}
	permission := model.PermissionAdvance{
		InfoEnterprise: model.ModuleMain{
			Status: model.ModuleDetail{
				Value:    true,
				Disabled: false,
			},
			Create: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			View: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Edit: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Delete: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Import: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Export: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Search: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Assign: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
		},
		AccountBalance: model.ModuleMain{
			Status: model.ModuleDetail{
				Value:    true,
				Disabled: false,
			},
			Create: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			View: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Edit: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Delete: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Import: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Export: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Search: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Assign: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
		},
		ServicePack: model.ModuleMain{
			Status: model.ModuleDetail{
				Value:    true,
				Disabled: false,
			},
			Create: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			View: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Edit: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Delete: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Import: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Export: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Search: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Assign: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
		},
		PaymentHistory: model.ModuleMain{
			Status: model.ModuleDetail{
				Value:    true,
				Disabled: false,
			},
			Create: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			View: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Edit: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Delete: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Import: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Export: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Search: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Assign: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
		},
		LoginConfig: model.ModuleMain{
			Status: model.ModuleDetail{
				Value:    true,
				Disabled: false,
			},
			Create: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			View: model.ModuleDetail{
				Value:    false,
				Disabled: false,
			},
			Edit: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Delete: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Import: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Export: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Search: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
			Assign: model.ModuleDetail{
				Value:    false,
				Disabled: true,
			},
		},
		Pbx: defaultData,
	}

	return permission
}

func HandleCreatePermissionUser() model.PermissionUser {
	return model.PermissionUser{}
}
