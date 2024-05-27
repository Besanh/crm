package service

import (
	"contactcenter-api/common/log"
	"contactcenter-api/common/model"
	"contactcenter-api/repository"
	"context"
	"fmt"

	"github.com/google/uuid"
)

var FileLua = "ring_groups"

func handleExtensionRingGroup(ctx context.Context, domain model.Domain, extension string, ringGroupConfig model.RingGroupConfig) error {
	dialplans := make([]model.Dialplan, 0)
	dialplanDetailsMap := make(map[string][]model.DialplanDetail, 0)
	ringGroups := make([]model.RingGroup, 0)
	ringGroupDestinations := make([]model.RingGroupDestination, 0)
	ringGroupExtensionMain := ringGroupConfig.Main
	ringGroupExtensionSub := ringGroupConfig.Sub
	scriptName := ringGroupConfig.Script
	ringGroupMainName := fmt.Sprintf("API_RG_MAIN_%s", ringGroupExtensionMain)
	ringGroupSubName := fmt.Sprintf("API_RG_SUB_%s", ringGroupExtensionSub)
	if len(scriptName) < 1 {
		scriptName = "api-masterise-cc-new.sh"
	}
	transferApp := ""
	transferData := ""
	// Destination
	if len(ringGroupConfig.Destination) > 0 {
		if err := handleRingGroupDestination(ctx, domain, ringGroupConfig.Destination, ringGroupExtensionMain); err != nil {
			log.Error(err)
			return err
		}
	}
	// SUB
	if len(ringGroupExtensionSub) > 0 {
		ringGroupSub, err := repository.RingGroupRepo.GetRingGroupByExtension(ctx, domain.DomainUuid, ringGroupExtensionSub)
		if err != nil {
			log.Error(err)
			return err
		}
		if ringGroupSub == nil {
			ringGroupSub = &model.RingGroup{
				DomainUuid:                  domain.DomainUuid,
				RingGroupUuid:               uuid.NewString(),
				RingGroupExtension:          ringGroupExtensionSub,
				RingGroupName:               ringGroupSubName,
				RingGroupGreeting:           "",
				RingGroupContext:            domain.DomainName,
				RingGroupCallTimeout:        0,
				RingGroupForwardDestination: "",
				RingGroupForwardEnabled:     "false",
				RingGroupCallerIdName:       "",
				RingGroupCallerIdNumber:     "",
				RingGroupCidNamePrefix:      "",
				RingGroupCidNumberPrefix:    "",
				RingGroupStrategy:           "enterprise",
				RingGroupTimeoutApp:         "",
				RingGroupTimeoutData:        "",
				RingGroupDistinctiveRing:    "",
				RingGroupRingback:           "${us-ring}",
				RingGroupMissedCallApp:      "",
				RingGroupMissedCallData:     "",
				RingGroupEnabled:            "true",
				RingGroupDescription:        "",
				DialplanUuid:                "",
				RingGroupForwardTollAllow:   "",
			}
		}

		// Add destination to ring group
		ringGroupDestinationSub, err := repository.RingGroupRepo.GetRingGroupDestinationOfExtension(ctx, domain.DomainUuid, ringGroupSub.RingGroupUuid, extension)
		if err != nil {
			log.Error(err)
			return err
		} else if ringGroupDestinationSub == nil {
			ringGroupDestinationSub = &model.RingGroupDestination{
				DomainUuid:               domain.DomainUuid,
				RingGroupDestinationUuid: uuid.NewString(),
				RingGroupUuid:            ringGroupSub.RingGroupUuid,
				DestinationNumber:        extension,
				DestinationDelay:         0,
				DestinationTimeout:       30,
				DestinationPrompt:        0,
			}
		}
		ringGroupDestinations = append(ringGroupDestinations, *ringGroupDestinationSub)
		transferApp = "transfer"
		transferData = fmt.Sprintf("%s XML %s", ringGroupExtensionSub, domain.DomainName)
		// Add Dialplan for sub ring group
		dialplanRingGroupSub, err := repository.DialplanRepo.GetDialplanByNumber(ctx, domain.DomainUuid, ringGroupExtensionSub)
		if err != nil {
			log.Error(err)
			return err
		} else if dialplanRingGroupSub == nil {
			dialplanRingGroupSub = &model.Dialplan{
				DomainUuid:          domain.DomainUuid,
				DialplanUuid:        uuid.NewString(),
				AppUuid:             "1d61fb65-1eec-bc73-a6ee-a6203b4fe6f2",
				Hostname:            "",
				DialplanContext:     domain.DomainName,
				DialplanName:        ringGroupSubName,
				DialplanNumber:      ringGroupExtensionSub,
				DialplanContinue:    "false",
				DialplanXml:         "",
				DialplanOrder:       100,
				DialplanEnabled:     "true",
				DialplanDescription: ringGroupSubName,
			}
		}
		order := 0
		dialplanDetailsRingGroupSub := []model.DialplanDetail{}
		dialplanDetail := model.DialplanDetail{
			DomainUuid:           domain.DomainUuid,
			DialplanDetailUuid:   uuid.NewString(),
			DialplanUuid:         dialplanRingGroupSub.DialplanUuid,
			DialplanDetailTag:    "condition",
			DialplanDetailType:   "destination_number",
			DialplanDetailData:   fmt.Sprintf("^%s$", ringGroupExtensionSub),
			DialplanDetailOrder:  order,
			DialplanDetailGroup:  0,
			DialplanDetailBreak:  "",
			DialplanDetailInline: "",
		}
		order += 25
		dialplanDetailsRingGroupSub = append(dialplanDetailsRingGroupSub, dialplanDetail)
		dialplanDetail = model.DialplanDetail{
			DomainUuid:           domain.DomainUuid,
			DialplanDetailUuid:   uuid.NewString(),
			DialplanUuid:         dialplanRingGroupSub.DialplanUuid,
			DialplanDetailTag:    "action",
			DialplanDetailType:   "playback",
			DialplanDetailData:   "${recordings_dir}/${domain_name}/ring3s.wav",
			DialplanDetailOrder:  order,
			DialplanDetailGroup:  0,
			DialplanDetailBreak:  "",
			DialplanDetailInline: "",
		}
		order += 25
		dialplanDetailsRingGroupSub = append(dialplanDetailsRingGroupSub, dialplanDetail)
		dialplanDetail = model.DialplanDetail{
			DomainUuid:           domain.DomainUuid,
			DialplanDetailUuid:   uuid.NewString(),
			DialplanUuid:         dialplanRingGroupSub.DialplanUuid,
			DialplanDetailTag:    "action",
			DialplanDetailType:   "set",
			DialplanDetailData:   fmt.Sprintf("ring_group_uuid=%s", ringGroupSub.RingGroupUuid),
			DialplanDetailOrder:  order,
			DialplanDetailGroup:  0,
			DialplanDetailBreak:  "",
			DialplanDetailInline: "",
		}
		order += 25
		dialplanDetailsRingGroupSub = append(dialplanDetailsRingGroupSub, dialplanDetail)
		dialplanDetail = model.DialplanDetail{
			DomainUuid:           domain.DomainUuid,
			DialplanDetailUuid:   uuid.NewString(),
			DialplanUuid:         dialplanRingGroupSub.DialplanUuid,
			DialplanDetailTag:    "action",
			DialplanDetailType:   "lua",
			DialplanDetailData:   "app.lua " + FileLua,
			DialplanDetailOrder:  order,
			DialplanDetailGroup:  0,
			DialplanDetailBreak:  "",
			DialplanDetailInline: "",
		}
		order += 25
		dialplanDetailsRingGroupSub = append(dialplanDetailsRingGroupSub, dialplanDetail)
		dialplanRingGroupSub.DialplanXml = ParseDialplanDetailsToXML(ctx, *dialplanRingGroupSub, dialplanDetailsRingGroupSub)
		dialplans = append(dialplans, *dialplanRingGroupSub)
		dialplanDetailsMap[dialplanRingGroupSub.DialplanUuid] = dialplanDetailsRingGroupSub
		ringGroupSub.DialplanUuid = dialplanRingGroupSub.DialplanUuid
		ringGroups = append(ringGroups, *ringGroupSub)
	}
	// MAIN
	ringGroupMain, err := repository.RingGroupRepo.GetRingGroupByExtension(ctx, domain.DomainUuid, ringGroupExtensionMain)
	if err != nil {
		log.Error(err)
		return err
	}
	if ringGroupMain == nil {
		ringGroupMain = &model.RingGroup{
			DomainUuid:                  domain.DomainUuid,
			RingGroupUuid:               uuid.NewString(),
			RingGroupExtension:          ringGroupExtensionMain,
			RingGroupName:               ringGroupMainName,
			RingGroupGreeting:           "",
			RingGroupContext:            domain.DomainName,
			RingGroupCallTimeout:        0,
			RingGroupForwardDestination: "",
			RingGroupForwardEnabled:     "false",
			RingGroupCallerIdName:       "",
			RingGroupCallerIdNumber:     "",
			RingGroupCidNamePrefix:      "",
			RingGroupCidNumberPrefix:    "",
			RingGroupStrategy:           "enterprise",
			RingGroupTimeoutApp:         transferApp,
			RingGroupTimeoutData:        transferData,
			RingGroupDistinctiveRing:    "",
			RingGroupRingback:           "${us-ring}",
			RingGroupMissedCallApp:      "",
			RingGroupMissedCallData:     "",
			RingGroupEnabled:            "true",
			RingGroupDescription:        "",
			DialplanUuid:                "",
			RingGroupForwardTollAllow:   "",
		}
	}
	ringGroupDestinationMain, err := repository.RingGroupRepo.GetRingGroupDestinationOfExtension(ctx, domain.DomainUuid, ringGroupMain.RingGroupUuid, extension)
	if err != nil {
		log.Error(err)
		return err
	} else if ringGroupDestinationMain == nil {
		ringGroupDestinationMain = &model.RingGroupDestination{
			DomainUuid:               domain.DomainUuid,
			RingGroupDestinationUuid: uuid.NewString(),
			RingGroupUuid:            ringGroupMain.RingGroupUuid,
			DestinationNumber:        extension,
			DestinationDelay:         0,
			DestinationTimeout:       30,
			DestinationPrompt:        0,
		}
	}
	ringGroupDestinations = append(ringGroupDestinations, *ringGroupDestinationMain)
	dialplanRingGroupMain, err := repository.DialplanRepo.GetDialplanByNumber(ctx, domain.DomainUuid, ringGroupExtensionMain)
	if err != nil {
		log.Error(err)
		return err
	} else if dialplanRingGroupMain == nil {
		dialplanRingGroupMain = &model.Dialplan{
			DomainUuid:          domain.DomainUuid,
			DialplanUuid:        uuid.NewString(),
			AppUuid:             "1d61fb65-1eec-bc73-a6ee-a6203b4fe6f2",
			Hostname:            "",
			DialplanContext:     domain.DomainName,
			DialplanName:        ringGroupMainName,
			DialplanNumber:      ringGroupExtensionMain,
			DialplanContinue:    "false",
			DialplanXml:         "",
			DialplanOrder:       100,
			DialplanEnabled:     "true",
			DialplanDescription: ringGroupMainName,
		}
	}
	// Add Dialplan for main ring group
	order := 0
	dialplanDetailsRingGroupMain := []model.DialplanDetail{}
	dialplanDetail := model.DialplanDetail{
		DomainUuid:           domain.DomainUuid,
		DialplanDetailUuid:   uuid.NewString(),
		DialplanUuid:         dialplanRingGroupMain.DialplanUuid,
		DialplanDetailTag:    "condition",
		DialplanDetailType:   "destination_number",
		DialplanDetailData:   fmt.Sprintf("^%s$", ringGroupExtensionMain),
		DialplanDetailOrder:  order,
		DialplanDetailGroup:  0,
		DialplanDetailBreak:  "",
		DialplanDetailInline: "",
	}
	if len(scriptName) > 0 {
		order += 25
		dialplanDetailsRingGroupMain = append(dialplanDetailsRingGroupMain, dialplanDetail)
		dialplanDetail = model.DialplanDetail{
			DomainUuid:           domain.DomainUuid,
			DialplanDetailUuid:   uuid.NewString(),
			DialplanUuid:         dialplanRingGroupMain.DialplanUuid,
			DialplanDetailTag:    "action",
			DialplanDetailType:   "system",
			DialplanDetailData:   fmt.Sprintf("/bin/bash %s ${caller_id_number} %s ${call_uuid}", scriptName, extension),
			DialplanDetailOrder:  order,
			DialplanDetailGroup:  0,
			DialplanDetailBreak:  "",
			DialplanDetailInline: "",
		}
	}
	order += 25
	dialplanDetailsRingGroupMain = append(dialplanDetailsRingGroupMain, dialplanDetail)
	dialplanDetail = model.DialplanDetail{
		DomainUuid:           domain.DomainUuid,
		DialplanDetailUuid:   uuid.NewString(),
		DialplanUuid:         dialplanRingGroupMain.DialplanUuid,
		DialplanDetailTag:    "action",
		DialplanDetailType:   "playback",
		DialplanDetailData:   "${recordings_dir}/${domain_name}/ring3s.wav",
		DialplanDetailOrder:  order,
		DialplanDetailGroup:  0,
		DialplanDetailBreak:  "",
		DialplanDetailInline: "",
	}
	order += 25
	dialplanDetailsRingGroupMain = append(dialplanDetailsRingGroupMain, dialplanDetail)
	dialplanDetail = model.DialplanDetail{
		DomainUuid:           domain.DomainUuid,
		DialplanDetailUuid:   uuid.NewString(),
		DialplanUuid:         dialplanRingGroupMain.DialplanUuid,
		DialplanDetailTag:    "action",
		DialplanDetailType:   "set",
		DialplanDetailData:   fmt.Sprintf("ring_group_uuid=%s", ringGroupMain.RingGroupUuid),
		DialplanDetailOrder:  order,
		DialplanDetailGroup:  0,
		DialplanDetailBreak:  "",
		DialplanDetailInline: "",
	}
	order += 25
	dialplanDetailsRingGroupMain = append(dialplanDetailsRingGroupMain, dialplanDetail)
	dialplanDetail = model.DialplanDetail{
		DomainUuid:           domain.DomainUuid,
		DialplanDetailUuid:   uuid.NewString(),
		DialplanUuid:         dialplanRingGroupMain.DialplanUuid,
		DialplanDetailTag:    "action",
		DialplanDetailType:   "lua",
		DialplanDetailData:   "app.lua " + FileLua,
		DialplanDetailOrder:  order,
		DialplanDetailGroup:  0,
		DialplanDetailBreak:  "",
		DialplanDetailInline: "",
	}
	order += 25
	dialplanDetailsRingGroupMain = append(dialplanDetailsRingGroupMain, dialplanDetail)
	dialplanRingGroupMain.DialplanXml = ParseDialplanDetailsToXML(ctx, *dialplanRingGroupMain, dialplanDetailsRingGroupMain)
	// End Add
	// Add to dialplan list
	dialplans = append(dialplans, *dialplanRingGroupMain)
	dialplanDetailsMap[dialplanRingGroupMain.DialplanUuid] = dialplanDetailsRingGroupMain
	ringGroupMain.DialplanUuid = dialplanRingGroupMain.DialplanUuid
	ringGroups = append(ringGroups, *ringGroupMain)
	if err := repository.ExtensionRepo.UpdateExtensionRingGroupTransaction(ctx, dialplans, dialplanDetailsMap, ringGroups, ringGroupDestinations); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func ParseDialplanDetailsToXML(ctx context.Context, dialplan model.Dialplan, dialplanDetails []model.DialplanDetail) string {
	detailGroups := make(map[int][]model.DialplanDetail)
	for _, dialplanDetail := range dialplanDetails {
		detailGroups[dialplanDetail.DialplanDetailGroup] = append(detailGroups[dialplanDetail.DialplanDetailGroup], dialplanDetail)
	}
	xml := ""
	xml += fmt.Sprintf("<extension name=\"%s\" continue=\"%s\" uuid=\"%s\">\n", dialplan.DialplanName, dialplan.DialplanContinue, dialplan.DialplanUuid)
	for _, detailGroup := range detailGroups {
		isCloseCondition := false
		for _, detail := range detailGroup {
			if detail.DialplanDetailTag == "condition" {
				isCloseCondition = true
				xml += fmt.Sprintf("	<%s field=\"%s\" expression=\"%s\">\n", detail.DialplanDetailTag, detail.DialplanDetailType, detail.DialplanDetailData)
			} else {
				xml += fmt.Sprintf("		<%s application=\"%s\" data=\"%s\"/>\n", detail.DialplanDetailTag, detail.DialplanDetailType, detail.DialplanDetailData)
			}
		}
		if isCloseCondition {
			xml += "	</condition>\n"
		}
	}
	xml += "</extension>\n"
	return xml
}

func handleRingGroupDestination(ctx context.Context, domain model.Domain, destinationNumber, ringGroupExtension string) error {
	destination, err := repository.DestinationRepo.GetDestinationByNumber(ctx, domain.DomainUuid, destinationNumber)
	if err != nil {
		log.Error(err)
		return err
	} else if destination != nil {
		return nil
	}
	dialplan := model.Dialplan{
		DomainUuid:          domain.DomainUuid,
		DialplanUuid:        uuid.NewString(),
		AppUuid:             "c03b422e-13a8-bd1b-e42b-b6b9b4d27ce4",
		Hostname:            "",
		DialplanContext:     "public",
		DialplanName:        destinationNumber,
		DialplanNumber:      destinationNumber,
		DialplanContinue:    "false",
		DialplanXml:         "",
		DialplanOrder:       100,
		DialplanEnabled:     "true",
		DialplanDescription: "",
	}
	dialplanDetails := []model.DialplanDetail{}
	order := 20
	dialplanDetail := model.DialplanDetail{
		DomainUuid:           domain.DomainUuid,
		DialplanDetailUuid:   uuid.NewString(),
		DialplanUuid:         dialplan.DialplanUuid,
		DialplanDetailTag:    "condition",
		DialplanDetailType:   "destination_number",
		DialplanDetailData:   fmt.Sprintf("^(%s)$", destinationNumber),
		DialplanDetailOrder:  order,
		DialplanDetailGroup:  0,
		DialplanDetailBreak:  "",
		DialplanDetailInline: "",
	}
	dialplanDetails = append(dialplanDetails, dialplanDetail)
	order += 20
	dialplanDetail = model.DialplanDetail{
		DomainUuid:           domain.DomainUuid,
		DialplanDetailUuid:   uuid.NewString(),
		DialplanUuid:         dialplan.DialplanUuid,
		DialplanDetailTag:    "action",
		DialplanDetailType:   "set",
		DialplanDetailData:   "hangup_after_bridge=true",
		DialplanDetailOrder:  order,
		DialplanDetailGroup:  0,
		DialplanDetailBreak:  "",
		DialplanDetailInline: "",
	}
	dialplanDetails = append(dialplanDetails, dialplanDetail)
	order += 10
	dialplanDetail = model.DialplanDetail{
		DomainUuid:           domain.DomainUuid,
		DialplanDetailUuid:   uuid.NewString(),
		DialplanUuid:         dialplan.DialplanUuid,
		DialplanDetailTag:    "action",
		DialplanDetailType:   "set",
		DialplanDetailData:   "continue_on_fail=true",
		DialplanDetailOrder:  order,
		DialplanDetailGroup:  0,
		DialplanDetailBreak:  "",
		DialplanDetailInline: "",
	}
	dialplanDetails = append(dialplanDetails, dialplanDetail)
	order += 10
	dialplanDetail = model.DialplanDetail{
		DomainUuid:           domain.DomainUuid,
		DialplanDetailUuid:   uuid.NewString(),
		DialplanUuid:         dialplan.DialplanUuid,
		DialplanDetailTag:    "action",
		DialplanDetailType:   "set",
		DialplanDetailData:   "record_path=${recordings_dir}/${domain_name}/archive/${strftime(%Y)}/${strftime(%b)}/${strftime(%d)}",
		DialplanDetailOrder:  order,
		DialplanDetailGroup:  0,
		DialplanDetailBreak:  "",
		DialplanDetailInline: "true",
	}
	dialplanDetails = append(dialplanDetails, dialplanDetail)
	order += 10
	dialplanDetail = model.DialplanDetail{
		DomainUuid:           domain.DomainUuid,
		DialplanDetailUuid:   uuid.NewString(),
		DialplanUuid:         dialplan.DialplanUuid,
		DialplanDetailTag:    "action",
		DialplanDetailType:   "set",
		DialplanDetailData:   "record_name=${uuid}.${record_ext}",
		DialplanDetailOrder:  order,
		DialplanDetailGroup:  0,
		DialplanDetailBreak:  "",
		DialplanDetailInline: "true",
	}
	dialplanDetails = append(dialplanDetails, dialplanDetail)
	order += 10
	dialplanDetail = model.DialplanDetail{
		DomainUuid:           domain.DomainUuid,
		DialplanDetailUuid:   uuid.NewString(),
		DialplanUuid:         dialplan.DialplanUuid,
		DialplanDetailTag:    "action",
		DialplanDetailType:   "set",
		DialplanDetailData:   "record_append=true",
		DialplanDetailOrder:  order,
		DialplanDetailGroup:  0,
		DialplanDetailBreak:  "",
		DialplanDetailInline: "true",
	}
	dialplanDetails = append(dialplanDetails, dialplanDetail)
	order += 10
	dialplanDetail = model.DialplanDetail{
		DomainUuid:           domain.DomainUuid,
		DialplanDetailUuid:   uuid.NewString(),
		DialplanUuid:         dialplan.DialplanUuid,
		DialplanDetailTag:    "action",
		DialplanDetailType:   "set",
		DialplanDetailData:   "record_in_progress=true",
		DialplanDetailOrder:  order,
		DialplanDetailGroup:  0,
		DialplanDetailBreak:  "",
		DialplanDetailInline: "true",
	}
	dialplanDetails = append(dialplanDetails, dialplanDetail)
	order += 10
	dialplanDetail = model.DialplanDetail{
		DomainUuid:           domain.DomainUuid,
		DialplanDetailUuid:   uuid.NewString(),
		DialplanUuid:         dialplan.DialplanUuid,
		DialplanDetailTag:    "action",
		DialplanDetailType:   "set",
		DialplanDetailData:   "recording_follow_transfer=true",
		DialplanDetailOrder:  order,
		DialplanDetailGroup:  0,
		DialplanDetailBreak:  "",
		DialplanDetailInline: "true",
	}
	dialplanDetails = append(dialplanDetails, dialplanDetail)
	order += 10
	dialplanDetail = model.DialplanDetail{
		DomainUuid:           domain.DomainUuid,
		DialplanDetailUuid:   uuid.NewString(),
		DialplanUuid:         dialplan.DialplanUuid,
		DialplanDetailTag:    "action",
		DialplanDetailType:   "record_session",
		DialplanDetailData:   "${record_path}/${record_name}",
		DialplanDetailOrder:  order,
		DialplanDetailGroup:  0,
		DialplanDetailBreak:  "",
		DialplanDetailInline: "true",
	}
	dialplanDetails = append(dialplanDetails, dialplanDetail)
	order += 10
	dialplanDetail = model.DialplanDetail{
		DomainUuid:           domain.DomainUuid,
		DialplanDetailUuid:   uuid.NewString(),
		DialplanUuid:         dialplan.DialplanUuid,
		DialplanDetailTag:    "action",
		DialplanDetailType:   "transfer",
		DialplanDetailData:   fmt.Sprintf("%s XML %s", ringGroupExtension, domain.DomainName),
		DialplanDetailOrder:  order,
		DialplanDetailGroup:  0,
		DialplanDetailBreak:  "",
		DialplanDetailInline: "true",
	}
	dialplanDetails = append(dialplanDetails, dialplanDetail)
	destination = &model.Destination{
		DomainUuid:                domain.DomainUuid,
		DestinationUuid:           uuid.NewString(),
		DialplanUuid:              dialplan.DialplanUuid,
		FaxUuid:                   "",
		DestinationType:           "inbound",
		DestinationNumber:         destinationNumber,
		DestinationNumberRegex:    fmt.Sprintf("^(%s)$", destinationNumber),
		DestinationCallerIdName:   "",
		DestinationCallerIdNumber: "",
		DestinationCidNamePrefix:  "",
		DestinationContext:        "public",
		DestinationRecord:         "true",
		DestinationAccountcode:    "",
		DestinationApp:            "transfer",
		DestinationData:           fmt.Sprintf("%s XML %s", ringGroupExtension, domain.DomainName),
		DestinationEnabled:        "true",
		DestinationDescription:    "",
	}
	xml := ""
	xml += "<extension name=\"" + destinationNumber + "\" continue=\"false\" uuid=\"" + dialplan.DialplanUuid + "\">\n"
	xml += "	<condition field=\"destination_number\" expression=\"" + fmt.Sprintf("^(%s)$", destinationNumber) + "\">\n"
	xml += "		<action application=\"export\" data=\"call_direction=inbound\" inline=\"true\"/>\n"
	xml += "		<action application=\"set\" data=\"domain_uuid=" + domain.DomainUuid + "\" inline=\"true\"/>\n"
	xml += "		<action application=\"set\" data=\"domain_name=" + domain.DomainName + "\" inline=\"true\"/>\n"
	xml += "		<action application=\"set\" data=\"hangup_after_bridge=true\" inline=\"true\"/>\n"
	xml += "		<action application=\"set\" data=\"continue_on_fail=true\" inline=\"true\"/>\n"
	xml += "		<action application=\"set\" data=\"record_path=${recordings_dir}/${domain_name}/archive/${strftime(%Y)}/${strftime(%b)}/${strftime(%d)}\" inline=\"true\"/>\n"
	xml += "		<action application=\"set\" data=\"record_name=${uuid}.${record_ext}\" inline=\"true\"/>\n"
	xml += "		<action application=\"set\" data=\"record_append=true\" inline=\"true\"/>\n"
	xml += "		<action application=\"set\" data=\"record_in_progress=true\" inline=\"true\"/>\n"
	xml += "		<action application=\"set\" data=\"recording_follow_transfer=true\" inline=\"true\"/>\n"
	xml += "		<action application=\"record_session\" data=\"${record_path}/${record_name}\" inline=\"false\"/>\n"
	xml += "		<action application=\"transfer\" data=\"" + fmt.Sprintf("%s XML %s", ringGroupExtension, domain.DomainName) + "\"/>\n"
	xml += "	</condition>\n"
	xml += "</extension>\n"
	dialplan.DialplanXml = xml
	if err := repository.DestinationRepo.InsertDestination(ctx, destination); err != nil {
		log.Error(err)
		return err
	}
	if err := repository.DialplanRepo.UpdateDialplanTransaction(ctx, dialplan, dialplanDetails); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
