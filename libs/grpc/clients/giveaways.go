package clients

import (
	"log"

	"github.com/twirapp/twir/libs/grpc/constants"
	"github.com/twirapp/twir/libs/grpc/giveaways"
	"google.golang.org/grpc"
)

func NewGiveaways(env string) giveaways.GiveawaysClient {
	serverAddress := createClientAddr(env, "giveaways", constants.GIVEAWAYS_SERVER_PORT)

	conn, err := grpc.Dial(
		serverAddress,
		defaultClientsOptions...,
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := giveaways.NewGiveawaysClient(conn)
	return c
}
