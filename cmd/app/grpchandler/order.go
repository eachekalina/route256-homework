package grpchandler

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"homework/internal/app/core"
	"homework/internal/app/logger"
	"homework/internal/app/metrics"
	"homework/internal/app/order"
	"homework/internal/app/pb"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	svc    *core.OrderCoreService
	log    logger.Logger
	metric *metrics.Metrics
	tracer trace.Tracer
}

func NewOrderService(svc *core.OrderCoreService, log logger.Logger, metric *metrics.Metrics) *OrderService {
	return &OrderService{svc: svc, log: log, metric: metric, tracer: otel.Tracer("cmd/app/grpchandler/order")}
}

func (s *OrderService) AcceptOrder(ctx context.Context, pbReq *pb.AcceptOrderRequest) (*pb.AcceptOrderResult, error) {
	ctx, span := s.tracer.Start(ctx, "AcceptOrder")
	defer span.End()

	req := core.AcceptOrderRequest{
		OrderId:       pbReq.OrderId,
		CustomerId:    pbReq.CustomerId,
		KeepDate:      pbReq.KeepDate.AsTime(),
		PriceRub:      pbReq.Price,
		WeightKg:      pbReq.WeightKg,
		PackagingType: pbReq.Packaging,
	}
	err := s.svc.AcceptOrder(ctx, req)
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
	ctx, span := s.tracer.Start(ctx, "ReturnOrder")
	defer span.End()

	err := s.svc.ReturnOrder(ctx, pbReq.OrderId)
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
	ctx, span := s.tracer.Start(ctx, "GiveOrders")
	defer span.End()

	orders, err := s.svc.GiveOrders(ctx, pbReq.OrderIds)
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
	s.metric.OrdersGiven.Add(float64(len(pbReq.OrderIds)))
	for _, o := range orders {
		s.metric.TimeBeforeGiven.Observe(o.GiveDate.Sub(o.AddDate).Seconds())
	}
	return &pb.GiveOrdersResult{}, nil
}

func (s *OrderService) ListOrders(ctx context.Context, pbReq *pb.ListOrdersRequest) (*pb.ListOrdersResult, error) {
	ctx, span := s.tracer.Start(ctx, "ListOrders")
	defer span.End()

	req := core.ListOrdersRequest{
		CustomerId:   pbReq.CustomerId,
		DisplayCount: 0,
		FilterGiven:  pbReq.StoredOnly,
	}
	orders, err := s.svc.ListOrders(ctx, req)
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
	ctx, span := s.tracer.Start(ctx, "AcceptReturn")
	defer span.End()

	req := core.AcceptReturnRequest{
		OrderId:    pbReq.OrderId,
		CustomerId: pbReq.CustomerId,
	}
	o, err := s.svc.AcceptReturn(ctx, req)
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
	s.metric.OrdersReturned.Add(1)
	s.metric.TimeBeforeReturn.Observe(o.ReturnDate.Sub(o.GiveDate).Seconds())
	return &pb.AcceptReturnResult{}, nil
}

func (s *OrderService) ListReturns(ctx context.Context, pbReq *pb.ListReturnsRequest) (*pb.ListReturnsResponse, error) {
	ctx, span := s.tracer.Start(ctx, "ListReturns")
	defer span.End()

	req := core.ListReturnsRequest{
		Count:   0,
		PageNum: 0,
	}
	returns, err := s.svc.ListReturns(ctx, req)
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
