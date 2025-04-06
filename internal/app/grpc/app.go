package grpc

import (
	"fmt"
	"github.com/bwjson/kolesa_auth/internal/grpc/auth"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, authService auth.Auth, port int) *App {
	grpcServer := grpc.NewServer()

	auth.Register(grpcServer, authService)

	return &App{
		log:        log,
		gRPCServer: grpcServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("Failed to run gRPC server: %v", err)
	}

	a.log.Info("gRPC server listening", slog.String("address:", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("Failed to run gRPC server: %v", err)
	}

	return nil
}

func (a *App) Stop() {
	a.log.Info("gRPC server stopping")

	a.gRPCServer.GracefulStop()
}
