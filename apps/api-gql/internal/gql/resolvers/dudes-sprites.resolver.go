package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"fmt"

	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
)

// DudesSelectSprite is the resolver for the dudesSelectSprite field.
func (r *mutationResolver) DudesSelectSprite(ctx context.Context, channelID string, spriteID string) (bool, error) {
	panic(fmt.Errorf("not implemented: DudesSelectSprite - dudesSelectSprite"))
}

// DudesUnselectSprite is the resolver for the dudesUnselectSprite field.
func (r *mutationResolver) DudesUnselectSprite(ctx context.Context, channelID string, spriteID string) (bool, error) {
	panic(fmt.Errorf("not implemented: DudesUnselectSprite - dudesUnselectSprite"))
}

// DudesForkSprite is the resolver for the dudesForkSprite field.
func (r *mutationResolver) DudesForkSprite(ctx context.Context, spriteID string) (*gqlmodel.DudeSprite, error) {
	panic(fmt.Errorf("not implemented: DudesForkSprite - dudesForkSprite"))
}

// DudesCreateSprite is the resolver for the dudesCreateSprite field.
func (r *mutationResolver) DudesCreateSprite(ctx context.Context, input gqlmodel.CreateSpriteInput) (*gqlmodel.DudeSprite, error) {
	panic(fmt.Errorf("not implemented: DudesCreateSprite - dudesCreateSprite"))
}

// DudesDeleteSprite is the resolver for the dudesDeleteSprite field.
func (r *mutationResolver) DudesDeleteSprite(ctx context.Context, spriteID string) (bool, error) {
	panic(fmt.Errorf("not implemented: DudesDeleteSprite - dudesDeleteSprite"))
}

// DudesUpdateSprite is the resolver for the dudesUpdateSprite field.
func (r *mutationResolver) DudesUpdateSprite(ctx context.Context, spriteID string, input gqlmodel.UpdateSpriteInput) (*gqlmodel.DudeSprite, error) {
	panic(fmt.Errorf("not implemented: DudesUpdateSprite - dudesUpdateSprite"))
}

// DudesCreateLayer is the resolver for the dudesCreateLayer field.
func (r *mutationResolver) DudesCreateLayer(ctx context.Context, input gqlmodel.CreateLayerInput) (*gqlmodel.DudeSpriteLayer, error) {
	panic(fmt.Errorf("not implemented: DudesCreateLayer - dudesCreateLayer"))
}

// DudesDeleteLayer is the resolver for the dudesDeleteLayer field.
func (r *mutationResolver) DudesDeleteLayer(ctx context.Context, layerID string) (bool, error) {
	panic(fmt.Errorf("not implemented: DudesDeleteLayer - dudesDeleteLayer"))
}

// DudesChannelSelectSprite is the resolver for the dudesChannelSelectSprite field.
func (r *mutationResolver) DudesChannelSelectSprite(ctx context.Context, userID string, spriteID string) (bool, error) {
	panic(fmt.Errorf("not implemented: DudesChannelSelectSprite - dudesChannelSelectSprite"))
}

// DudesChannelUnselectSprite is the resolver for the dudesChannelUnselectSprite field.
func (r *mutationResolver) DudesChannelUnselectSprite(ctx context.Context, userID string, spriteID string) (bool, error) {
	panic(fmt.Errorf("not implemented: DudesChannelUnselectSprite - dudesChannelUnselectSprite"))
}

// DudesCatalogSprite is the resolver for the dudesCatalogSprite field.
func (r *queryResolver) DudesCatalogSprite(ctx context.Context, id string) (*gqlmodel.DudeSprite, error) {
	panic(fmt.Errorf("not implemented: DudesCatalogSprite - dudesCatalogSprite"))
}

// DudesCatalogSprites is the resolver for the dudesCatalogSprites field.
func (r *queryResolver) DudesCatalogSprites(ctx context.Context, input gqlmodel.DudesCatalogSpritesInput) ([]gqlmodel.DudeSprite, error) {
	panic(fmt.Errorf("not implemented: DudesCatalogSprites - dudesCatalogSprites"))
}

// DudesCatalogLayers is the resolver for the dudesCatalogLayers field.
func (r *queryResolver) DudesCatalogLayers(ctx context.Context, search *string, approved *bool, layersTypes []gqlmodel.DudeSpriteLayerType) ([]gqlmodel.DudeSpriteLayer, error) {
	panic(fmt.Errorf("not implemented: DudesCatalogLayers - dudesCatalogLayers"))
}
