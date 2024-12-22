package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"

	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/dataloader"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/graph"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	audit_logs "github.com/twirapp/twir/apps/api-gql/internal/services/audit-logs"
)

// User is the resolver for the user field.
func (r *auditLogResolver) User(ctx context.Context, obj *gqlmodel.AuditLog) (*gqlmodel.TwirUserTwitchInfo, error) {
	if obj.UserID == nil {
		return nil, nil
	}

	return dataloader.GetHelixUserById(ctx, *obj.UserID)
}

// AuditLog is the resolver for the auditLog field.
func (r *queryResolver) AuditLog(ctx context.Context) ([]gqlmodel.AuditLog, error) {
	dashboardID, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	logs, err := r.deps.AuditLogsService.GetMany(
		ctx, audit_logs.GetManyInput{
			ChannelID: &dashboardID,
			Page:      0,
			Limit:     100,
		},
	)
	if err != nil {
		return nil, err
	}

	result := make([]gqlmodel.AuditLog, 0, len(logs))
	for _, l := range logs {
		result = append(result, mappers.AuditLogToGql(l))
	}

	return result, nil
}

// AuditLog is the resolver for the auditLog field.
func (r *subscriptionResolver) AuditLog(ctx context.Context) (<-chan *gqlmodel.AuditLog, error) {
	dashboardID, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	logsChannel, err := r.deps.AuditLogsService.Subscribe(ctx, dashboardID)
	if err != nil {
		return nil, err
	}

	channel := make(chan *gqlmodel.AuditLog)

	go func() {
		defer func() {
			close(channel)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case auditLog := <-logsChannel:
				val := mappers.AuditLogToGql(auditLog)
				channel <- &val
			}
		}
	}()

	return channel, nil
}

// AuditLog returns graph.AuditLogResolver implementation.
func (r *Resolver) AuditLog() graph.AuditLogResolver { return &auditLogResolver{r} }

type auditLogResolver struct{ *Resolver }
