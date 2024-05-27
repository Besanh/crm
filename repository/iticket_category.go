package repository

import (
	"contactcenter-api/common/model"
	"context"
)

type ITicketCategory interface {
	InsertTicketCategory(ctx context.Context, ticketCategory *model.TicketCategory) error
	InsertTicketCategoryAndSLAPolicy(ctx context.Context, ticketCategory *model.TicketCategory, slaPolicies *[]model.SlaPolicy) error
	GetTicketCategories(ctx context.Context, domainUuid string, categoryCode string) (*[]model.TicketCategory, error)
	GetTicketCategoryByCode(ctx context.Context, domainUuid string, categoryCode string) (*model.TicketCategory, error)
	GetTicketCategoriesInfo(ctx context.Context, domainUuid string, limit, offset int, filter model.TicketCategoryFilter) (*[]model.TicketCategoryInfo, int, error)
	GetParentTicketCategoryById(ctx context.Context, domainUuid string, ticketCategoryUuid string) (*model.TicketCategory, error)
	GetParentTicketCategoriesById(ctx context.Context, domainUuid string, ticketCategoryUuid string) (*[]model.TicketCategory, error)
	GetTicketCategoryById(ctx context.Context, domainUuid string, ticketCategoryUuid string) (*model.TicketCategory, error)
	GetTicketCategoryInfoById(ctx context.Context, domainUuid string, ticketCategoryUuid string) (*model.TicketCategoryInfo, error)
	UpdateTicketCategory(ctx context.Context, ticketCategory *model.TicketCategory) error
	DeleteTicketCategory(ctx context.Context, domainUuid, ticketCategoryUuid string) error
	UpdateTicketCategoryAndSLAPolicy(ctx context.Context, ticketCategory *model.TicketCategory, slaPolicies *[]model.SlaPolicyInfo) error
}

var TicketCategoryRepo ITicketCategory
