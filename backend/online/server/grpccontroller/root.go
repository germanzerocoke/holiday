package grpccontroller

import (
	"backend/online/server/controller"
	pb "backend/proto"
)

type GrpcController struct {
	pb.UnimplementedSignalServiceServer
	controller *controller.Controller
}

func NewGrpcController(controller *controller.Controller) *GrpcController {
	gc := GrpcController{controller: controller}

	return &gc
}
