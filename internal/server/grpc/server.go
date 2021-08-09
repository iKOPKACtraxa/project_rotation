package grpcserver

import (
	"context"
	"fmt"
	"net"

	"github.com/iKOPKACtraxa/project_rotation/internal/app"
	"github.com/iKOPKACtraxa/project_rotation/internal/pb"
	"github.com/iKOPKACtraxa/project_rotation/internal/storage"
	"google.golang.org/grpc"
)

type Service struct {
	rotation *app.App
	pb.UnimplementedRotationServer
}

func newService(rotation *app.App) *Service {
	return &Service{
		rotation: rotation,
	}
}

//
func (s *Service) AddBanner(ctx context.Context, req *pb.AddBannerRequest) (*pb.AddBannerResponse, error) {
	err := s.rotation.AddBanner(ctx, storage.ID(req.BannerID), storage.ID(req.SlotID)) // todo это другой контекст? GRPC его предоставил? А ctx из main уже тут не нужен?
	if err != nil {
		return nil, err
	}
	return &pb.AddBannerResponse{}, nil
}

//
func (s *Service) DeleteBanner(ctx context.Context, req *pb.DeleteBannerRequest) (*pb.DeleteBannerResponse, error) {
	err := s.rotation.DeleteBanner(ctx, storage.ID(req.BannerID), storage.ID(req.SlotID))
	if err != nil {
		return nil, err
	}
	return &pb.DeleteBannerResponse{}, nil
}

//
func (s *Service) ClicksIncreasing(ctx context.Context, req *pb.ClicksIncreasingRequest) (*pb.ClicksIncreasingResponse, error) {
	err := s.rotation.ClicksIncreasing(ctx, storage.ID(req.SlotID), storage.ID(req.BannerID), storage.ID(req.SocGroupID))
	if err != nil {
		return nil, err
	}
	return &pb.ClicksIncreasingResponse{}, nil
}

//
func (s *Service) BannerSelection(ctx context.Context, req *pb.BannerSelectionRequest) (*pb.BannerSelectionResponse, error) {
	id, err := s.rotation.BannerSelection(ctx, storage.ID(req.SlotID), storage.ID(req.SocGroupID))
	if err != nil {
		return nil, err
	}
	return &pb.BannerSelectionResponse{
		BannerID: uint64(id),
	}, nil
}

//
func Serve(ctx context.Context, rotation *app.App, port string) error {
	lsn, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("listening of tcp port %v has got an error: %w", port, err)
	}

	server := grpc.NewServer()
	go func() {
		<-ctx.Done()
		rotation.Logger.Info("server GRPC is stopping gracefully...")
		server.GracefulStop()
		rotation.Logger.Info("server GRPC is stopped gracefully...")
	}()
	pb.RegisterRotationServer(server, newService(rotation))

	rotation.Logger.Info("starting GRPC server on", lsn.Addr().String())
	if err := server.Serve(lsn); err != nil {
		return fmt.Errorf("serving has got an error: %w", err)
	}
	return nil
}
