package grpchandler

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"homework/internal/app/core"
	"homework/internal/app/logger"
	"homework/internal/app/order"
	"homework/internal/app/pb"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	svc *core.OrderCoreService
	log logger.Logger
}

func NewOrderService(svc *core.OrderCoreService, log logger.Logger) *OrderService {
	return &OrderService{svc: svc, log: log}
}

func (s *OrderService) AcceptOrder(ctx context.Context, pbReq *pb.AcceptOrderRequest) (*pb.AcceptOrderResult, error) {
	req := core.AcceptOrderRequest{
		OrderId:       pbReq.OrderId,
		CustomerId:    pbReq.CustomerId,
		KeepDate:      pbReq.KeepDate.AsTime(),
		PriceRub:      pbReq.Price,
		WeightKg:      pbReq.WeightKg,
		PackagingType: pbReq.Packaging,
	}
	err := s.svc.AcceptOrder(req)
	if err != nil {
		if errors.As(err, &core.ValidationError{}) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, order.ErrIdAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		s.log.Log("%v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.AcceptOrderResult{}, nil
}

func (s *OrderService) ReturnOrder(ctx context.Context, pbReq *pb.ReturnOrderRequest) (*pb.ReturnOrderResult, error) {
	err := s.svc.ReturnOrder(pbReq.OrderId)
	if err != nil {
		if errors.As(err, &core.ValidationError{}) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, order.ErrNoItemFound) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		s.log.Log("%v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.ReturnOrderResult{}, nil
}

func (s *OrderService) GiveOrders(ctx context.Context, pbReq *pb.GiveOrdersRequest) (*pb.GiveOrdersResult, error) {
	err := s.svc.GiveOrders(pbReq.OrderIds)
	if err != nil {
		if errors.As(err, &core.ValidationError{}) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, order.ErrNoItemFound) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		s.log.Log("%v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.GiveOrdersResult{}, nil
}

func (s *OrderService) ListOrders(ctx context.Context, pbReq *pb.ListOrdersRequest) (*pb.ListOrdersResult, error) {
	req := core.ListOrdersRequest{
		CustomerId:   pbReq.CustomerId,
		DisplayCount: 0,
		FilterGiven:  pbReq.StoredOnly,
	}
	orders, err := s.svc.ListOrders(req)
	if err != nil {
		if errors.As(err, &core.ValidationError{}) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, order.ErrNoItemFound) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		s.log.Log("%v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resultOrders := make([]*pb.Order, len(orders))
	for i, o := range orders {
		var giveDate, returnDate *timestamppb.Timestamp
		if o.IsGiven {
			giveDate = timestamppb.New(o.GiveDate)
		}
		if o.IsReturned {
			returnDate = timestamppb.New(o.ReturnDate)
		}
		resultOrders[i] = &pb.Order{
			GiveDate:   giveDate,
			ReturnDate: returnDate,
			KeepDate:   timestamppb.New(o.KeepDate),
			AddDate:    timestamppb.New(o.AddDate),
			Id:         o.Id,
			CustomerId: o.CustomerId,
			PriceRub:   o.PriceRub,
			WeightKg:   o.WeightKg,
			IsGiven:    o.IsGiven,
			IsReturned: o.IsReturned,
		}
	}
	return &pb.ListOrdersResult{Orders: resultOrders}, nil
}

func (s *OrderService) AcceptReturn(ctx context.Context, pbReq *pb.AcceptReturnRequest) (*pb.AcceptReturnResult, error) {
	req := core.AcceptReturnRequest{
		OrderId:    pbReq.OrderId,
		CustomerId: pbReq.CustomerId,
	}
	err := s.svc.AcceptReturn(req)
	if err != nil {
		if errors.As(err, &core.ValidationError{}) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, order.ErrNoItemFound) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		s.log.Log("%v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.AcceptReturnResult{}, nil
}

func (s *OrderService) ListReturns(ctx context.Context, pbReq *pb.ListReturnsRequest) (*pb.ListReturnsResponse, error) {
	req := core.ListReturnsRequest{
		Count:   0,
		PageNum: 0,
	}
	returns, err := s.svc.ListReturns(req)
	if err != nil {
		if errors.As(err, &core.ValidationError{}) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, order.ErrNoItemFound) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		s.log.Log("%v", err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resultReturns := make([]*pb.Order, len(returns))
	for i, o := range returns {
		var giveDate, returnDate *timestamppb.Timestamp
		if o.IsGiven {
			giveDate = timestamppb.New(o.GiveDate)
		}
		if o.IsReturned {
			returnDate = timestamppb.New(o.ReturnDate)
		}
		resultReturns[i] = &pb.Order{
			GiveDate:   giveDate,
			ReturnDate: returnDate,
			KeepDate:   timestamppb.New(o.KeepDate),
			AddDate:    timestamppb.New(o.AddDate),
			Id:         o.Id,
			CustomerId: o.CustomerId,
			PriceRub:   o.PriceRub,
			WeightKg:   o.WeightKg,
			IsGiven:    o.IsGiven,
			IsReturned: o.IsReturned,
		}
	}
	return &pb.ListReturnsResponse{Returns: resultReturns}, nil
}
