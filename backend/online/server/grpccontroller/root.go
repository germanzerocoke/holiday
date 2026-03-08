package grpccontroller

import (
	"backend/online/server/controller"
	pb "backend/proto"
	"context"
)

type GrpcController struct {
	pb.UnimplementedSignalServiceServer
	controller *controller.Controller
}

func NewGrpcController(controller *controller.Controller) *GrpcController {
	gc := GrpcController{controller: controller}

	return &gc
}

func (gc *GrpcController) RelaySignal(ctx context.Context, in *pb.RelaySignalRequest) (*pb.RelaySignalResponse, error) {
	err := gc.controller.RelaySignal(ctx, in.FromId, in.ToId, in.Signal)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
