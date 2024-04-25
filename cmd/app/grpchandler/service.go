package grpchandler

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"homework/internal/app/core"
	"homework/internal/app/logger"
	"homework/internal/app/pb"
	"homework/internal/app/pickuppoint"
)

type Service struct {
	pb.UnimplementedPickUpPointServiceServer
	svc core.PickUpPointCoreService
	log logger.Logger
}

func NewService(svc core.PickUpPointCoreService, log logger.Logger) *Service {
	return &Service{svc: svc, log: log}
}

func (s *Service) Create(ctx context.Context, pbReq *pb.PickUpPointCreateRequest) (*pb.PickUpPointCreateResult, error) {
	req := core.CreatePointRequest{
		Id:      pbReq.Point.Id,
		Name:    pbReq.Point.Name,
		Address: pbReq.Point.Address,
		Contact: pbReq.Point.Contact,
	}

	err := s.svc.CreatePoint(ctx, req)
	if err != nil {
		if errors.Is(err, pickuppoint.ErrIdAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		s.log.Log("%v", err)
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &pb.PickUpPointCreateResult{}, nil
}

func (s *Service) List(ctx context.Context, pbReq *pb.PickUpPointListRequest) (*pb.PickUpPointListResult, error) {
	list, err := s.svc.ListPoints(ctx)
	if err != nil {
		s.log.Log("%v", err)
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	res := &pb.PickUpPointListResult{Points: make([]*pb.PickUpPoint, len(list))}
	for i, point := range list {
		res.Points[i] = &pb.PickUpPoint{
			Id:      point.Id,
			Name:    point.Name,
			Address: point.Address,
			Contact: point.Contact,
		}
	}
	return res, nil
}

func (s *Service) Get(ctx context.Context, pbReq *pb.PickUpPointGetRequest) (*pb.PickUpPointGetResult, error) {
	point, err := s.svc.GetPoint(ctx, pbReq.Id)
	if err != nil {
		if errors.Is(err, pickuppoint.ErrNoItemFound) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		s.log.Log("%v", err)
		return nil, status.Errorf(codes.Internal, "internal error")
	}
	res := &pb.PickUpPointGetResult{Point: &pb.PickUpPoint{
		Id:      point.Id,
		Name:    point.Name,
		Address: point.Address,
		Contact: point.Contact,
	}}
	return res, nil
}

func (s *Service) Update(ctx context.Context, pbReq *pb.PickUpPointUpdateRequest) (*pb.PickUpPointUpdateResult, error) {
	req := core.UpdatePointRequest{
		Id:      pbReq.Point.Id,
		Name:    pbReq.Point.Name,
		Address: pbReq.Point.Address,
		Contact: pbReq.Point.Contact,
	}
	err := s.svc.UpdatePoint(ctx, req)
	if err != nil {
		if errors.Is(err, pickuppoint.ErrNoItemFound) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		s.log.Log("%v", err)
		return nil, status.Errorf(codes.Internal, "internal error")
	}
	return &pb.PickUpPointUpdateResult{}, nil
}

func (s *Service) Delete(ctx context.Context, pbReq *pb.PickUpPointDeleteRequest) (*pb.PickUpPointDeleteResult, error) {
	err := s.svc.DeletePoint(ctx, pbReq.Id)
	if err != nil {
		if errors.Is(err, pickuppoint.ErrNoItemFound) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		s.log.Log("%v", err)
		return nil, status.Errorf(codes.Internal, "internal error")
	}
	return &pb.PickUpPointDeleteResult{}, nil
}
