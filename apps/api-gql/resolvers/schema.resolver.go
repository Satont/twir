package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"fmt"

	"github.com/twirapp/twir/apps/api-gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/graph"
)

// Empty is the resolver for the empty field.
func (r *mutationResolver) Empty(ctx context.Context) (*gqlmodel.Empty, error) {
	panic(fmt.Errorf("not implemented: Empty - empty"))
}

// Empty is the resolver for the empty field.
func (r *queryResolver) Empty(ctx context.Context) (*gqlmodel.Empty, error) {
	panic(fmt.Errorf("not implemented: Empty - empty"))
}

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
