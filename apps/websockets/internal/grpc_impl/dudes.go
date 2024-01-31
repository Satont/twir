package grpc_impl

import (
	"context"

	"github.com/satont/twir/apps/websockets/internal/protoutils"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *GrpcImpl) DudesJump(_ context.Context, req *websockets.DudesJumpRequest) (
	*emptypb.Empty,
	error,
) {
	json, err := protoutils.CreateJsonWithProto(req, nil)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, c.dudesServer.SendEvent(
		req.GetChannelId(),
		"jump",
		json,
	)
}
