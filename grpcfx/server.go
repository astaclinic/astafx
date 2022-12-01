package grpcfx

import (
	"context"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/reflection"
)

var Module = fx.Module("grpc",
	fx.Provide(NewGrpcServer),
	fx.Provide(health.NewServer), // Add health check
	fx.Invoke(RunGrpcServer),
	fx.Invoke(registerHealthCheckGrpcServer),
)

type GrpcConfig struct {
	ListenAddr string `mapstructure:"listen_addr" yaml:"listen_addr"  validate:"required,hostname_port"`
}

func NewGrpcServer() *grpc.Server {
	ser := grpc.NewServer()
	reflection.Register(ser) // Enable reflection
	return ser
}

type RunGrpcServerParams struct {
	fx.In
	Lifecycle  fx.Lifecycle
	GrpcServer *grpc.Server
	Config     *GrpcConfig
}

func RunGrpcServer(p RunGrpcServerParams) {
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			lis, err := net.Listen("tcp", p.Config.ListenAddr)
			if err != nil {
				return err
			}
			go p.GrpcServer.Serve(lis)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			p.GrpcServer.Stop()
			return nil
		},
	})
}
