package grpcfx

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func registerHealthCheckGrpcServer(grpcServer *grpc.Server, healthCheck *health.Server) {
	grpc_health_v1.RegisterHealthServer(grpcServer, healthCheck)
}
