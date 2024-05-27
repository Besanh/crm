package service

import (
	"contactcenter-api/common/log"
	"contactcenter-api/common/response"
	"contactcenter-api/common/util"
	"contactcenter-api/repository"
	"context"
	"errors"
	"os"
)

type (
	IAttachment interface {
		DownloadAttachment(ctx context.Context, domainUuid, folder, fileName string) (string, error)
		DeleteAttachment(ctx context.Context, domainUuid, folder, fileName, entityUuid string) (int, any)
	}
	Attachment struct{}
)

func NewAttachment() IAttachment {
	return &Attachment{}
}

func (s *Attachment) DownloadAttachment(ctx context.Context, domainUuid, folder, fileName string) (string, error) {
	// TODO
	pathFile := util.PUBLIC_DIR + domainUuid + "/" + folder + "/" + fileName
	_, err := os.Stat(pathFile)
	if os.IsNotExist(err) {
		return "", err
	}
	return pathFile, nil
}

func (s *Attachment) DeleteAttachment(ctx context.Context, domainUuid, folder, fileName, entityUuid string) (int, any) {
	// TODO
	pathFile := util.PUBLIC_DIR + domainUuid + "/" + folder + "/" + fileName
	err := os.Remove(pathFile)
	if err != nil {
		return response.ServiceUnavailableMsg(err.Error())
	}
	if len(entityUuid) > 0 {
		if folder == "ticket" {
			ticket, err := repository.TicketRepo.GetTicketById(ctx, domainUuid, entityUuid)
			if err != nil {
				log.Error(err)
				return response.ServiceUnavailableMsg(err.Error())
			} else if ticket == nil {
				return response.ServiceUnavailableMsg(errors.New("ticket is not found"))
			}
			var attachment []string
			if len(ticket.Attachment) > 0 {
				for _, val := range ticket.Attachment {
					if fileName != val {
						attachment = append(attachment, val)
					}
				}
			}
			ticket.Attachment = attachment
			if err := repository.TicketRepo.PatchTicketAttachment(ctx, domainUuid, *ticket); err != nil {
				log.Error(err)
				return response.ServiceUnavailableMsg(err.Error())
			}
		}
	}
	return response.OK(map[string]any{})
}
