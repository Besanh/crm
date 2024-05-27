package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type ISipRegistration interface {
	GetSipRegistrationsNotCallOfDomain(ctx context.Context, domainName string) (*[]model.SipRegistration, error)
	GetSipRegistrationOfExtension(ctx context.Context, domainName, extension string) (any, error)
	GetSipRegistrations(ctx context.Context, domainName string) (*[]model.SipRegistrationMonitorView, error)
}

var SipRegistrationRepo ISipRegistration
