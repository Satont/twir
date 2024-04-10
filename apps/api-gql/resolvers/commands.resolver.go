package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/graph"
)

// CreateCommand is the resolver for the createCommand field.
func (r *mutationResolver) CreateCommand(ctx context.Context, opts gqlmodel.CreateCommandInput) (*gqlmodel.Command, error) {
	responses := make([]gqlmodel.CommandResponse, 0, len(opts.Responses.Value()))
	for _, response := range opts.Responses.Value() {
		responses = append(
			responses, gqlmodel.CommandResponse{
				ID:    uuid.NewString(),
				Text:  response.Text,
				Order: response.Order,
			},
		)
	}

	newCommand := gqlmodel.Command{
		ID:          uuid.NewString(),
		Name:        opts.Name,
		Description: opts.Description.Value(),
		Aliases:     opts.Aliases.Value(),
		Responses:   responses,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	r.NewCommandChann <- &newCommand

	commands = append(commands, newCommand)
	return &newCommand, nil
}

// UpdateCommand is the resolver for the updateCommand field.
func (r *mutationResolver) UpdateCommand(ctx context.Context, id string, opts gqlmodel.UpdateCommandOpts) (*gqlmodel.Command, error) {
	var cmd *gqlmodel.Command
	cmdIndex := 0

	for i, command := range commands {
		if command.ID == id {
			cmd = &command
			cmdIndex = i
			break
		}
	}

	if cmd == nil {
		return nil, fmt.Errorf("command with id %s not found", id)
	}

	if opts.Name.IsSet() {
		cmd.Name = *opts.Name.Value()
	}

	if opts.Description.IsSet() {
		cmd.Description = opts.Description.Value()
	}

	if opts.Aliases.IsSet() {
		cmd.Aliases = opts.Aliases.Value()
	}

	cmd.UpdatedAt = time.Now()

	commands[cmdIndex] = *cmd

	return cmd, nil
}

// Commands is the resolver for the commands field.
func (r *queryResolver) Commands(ctx context.Context) ([]gqlmodel.Command, error) {
	return commands, nil
}

// NewCommand is the resolver for the newCommand field.
func (r *subscriptionResolver) NewCommand(ctx context.Context) (<-chan *gqlmodel.Command, error) {
	ch := make(chan *gqlmodel.Command)

	fmt.Println("Subscription Started")

	go func() {
		defer close(ch)

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Subscription Closed")
				return

			case cmd := <-r.NewCommandChann:
				fmt.Println("New Command")
				ch <- cmd
			}
		}
	}()

	// We return the channel and no error.
	return ch, nil
}

// Subscription returns graph.SubscriptionResolver implementation.
func (r *Resolver) Subscription() graph.SubscriptionResolver { return &subscriptionResolver{r} }

type subscriptionResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
var commands []gqlmodel.Command
