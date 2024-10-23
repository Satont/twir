package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"log/slog"

	model "github.com/satont/twir/libs/gomodels"
	dataloader "github.com/twirapp/twir/apps/api-gql/internal/gql/data-loader"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/graph"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/mappers"
)

// User is the resolver for the user field.
func (r *auditLogResolver) User(ctx context.Context, obj *gqlmodel.AuditLog) (*gqlmodel.TwirUserTwitchInfo, error) {
	if obj.UserID == nil {
		return nil, nil
	}

	user, err := dataloader.GetHelixUserById(ctx, *obj.UserID)
	if err != nil {
		r.logger.Error(
			"failed to get helix user for audit log",
			slog.String("user_id", *obj.UserID),
			slog.String("audit_log_id", obj.ID.String()),
		)

		return nil, err
	}

	return user, nil
}

// AuditLog is the resolver for the auditLog field.
func (r *queryResolver) AuditLog(ctx context.Context) ([]gqlmodel.AuditLog, error) {
	dashboardID, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	var logs []model.AuditLog
	if err := r.gorm.
		WithContext(ctx).
		Limit(100).
		Order("created_at DESC").
		Where("channel_id = ?", dashboardID).
		Find(&logs).Error; err != nil {
		r.logger.Error("error in fetching audit logs", slog.Any("err", err))
		return nil, err
	}

	gqllogs := make([]gqlmodel.AuditLog, 0, len(logs))
	for _, l := range logs {
		gqllogs = append(
			gqllogs,
			gqlmodel.AuditLog{
				ID:            l.ID,
				System:        mappers.AuditTableNameToGqlSystem(l.Table),
				OperationType: mappers.AuditTypeModelToGql(l.OperationType),
				OldValue:      l.OldValue.Ptr(),
				NewValue:      l.NewValue.Ptr(),
				ObjectID:      l.ObjectID.Ptr(),
				UserID:        l.UserID.Ptr(),
				CreatedAt:     l.CreatedAt,
			},
		)
	}

	return gqllogs, nil
}

// AuditLog is the resolver for the auditLog field.
func (r *subscriptionResolver) AuditLog(ctx context.Context) (<-chan *gqlmodel.AuditLog, error) {
	dashboardID, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	auditLogs, err := r.auditLogsPubSub.Subscribe(ctx, dashboardID)
	if err != nil {
		return nil, err
	}

	channel := make(chan *gqlmodel.AuditLog)

	go func() {
		defer func() {
			_ = auditLogs.Close()
			close(channel)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case auditLog := <-auditLogs.Channel():
				channel <- mappers.AuditLogToGql(auditLog)
			}
		}
	}()

	return channel, nil
}

// AuditLog returns graph.AuditLogResolver implementation.
func (r *Resolver) AuditLog() graph.AuditLogResolver { return &auditLogResolver{r} }

type auditLogResolver struct{ *Resolver }
