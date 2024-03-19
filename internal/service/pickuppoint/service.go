package pickuppoint

import (
	"context"
	"golang.org/x/sync/errgroup"
	"homework/internal/logger"
	"homework/internal/model"
)

type storage interface {
	Create(ctx context.Context, point model.PickUpPoint) error
	List(ctx context.Context) ([]model.PickUpPoint, error)
	Get(ctx context.Context, id uint64) (model.PickUpPoint, error)
	Update(ctx context.Context, point model.PickUpPoint) error
	Delete(ctx context.Context, id uint64) error
}

const (
	createReq = iota
	listReq
	getReq
	updateReq
	deleteReq
)

type request struct {
	reqType int
	id      uint64
	point   model.PickUpPoint
}

// Service allows concurrent working on pick-up points.
type Service struct {
	stor  storage
	log   *logger.Logger
	read  chan request
	write chan request
}

// NewPickUpPointService creates a new Service.
func NewPickUpPointService(stor storage, log *logger.Logger) *Service {
	s := &Service{
		stor: stor,
		log:  log,
	}
	return s
}

func (s *Service) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return s.writeThread(ctx)
	})

	eg.Go(func() error {
		return s.readThread(ctx)
	})

	return eg.Wait()
}

func (s *Service) writeThread(ctx context.Context) error {
	s.write = make(chan request)
	for {
		s.log.Log("write thread: waiting for request")
		select {
		case <-ctx.Done():
			s.log.Log("write thread: closing")
			close(s.write)
			return ctx.Err()
		case r := <-s.write:
			switch r.reqType {
			case createReq:
				s.log.Log("write thread: requested create")
				err := s.stor.Create(ctx, r.point)
				if err != nil {
					s.log.Log("write thread: error: %v", err)
				}
			case updateReq:
				s.log.Log("write thread: requested update")
				err := s.stor.Update(ctx, r.point)
				if err != nil {
					s.log.Log("write thread: error: %v", err)
				}
			case deleteReq:
				s.log.Log("write thread: requested delete")
				err := s.stor.Delete(ctx, r.id)
				if err != nil {
					s.log.Log("write thread: error: %v", err)
				}
			default:
				s.log.Log("write thread: invalid request")
			}
		}
	}
}

func (s *Service) readThread(ctx context.Context) error {
	s.read = make(chan request)

	for {
		s.log.Log("read thread: waiting for request")
		select {
		case <-ctx.Done():
			s.log.Log("read thread: closing")
			close(s.read)
			return ctx.Err()
		case r := <-s.read:
			switch r.reqType {
			case listReq:
				s.log.Log("read thread: requested list")
				list, err := s.stor.List(ctx)
				if err != nil {
					s.log.Log("read thread: error: %v", err)
				} else {
					s.log.PrintPoints(list)
				}
			case getReq:
				s.log.Log("read thread: requested get")
				point, err := s.stor.Get(ctx, r.id)
				if err != nil {
					s.log.Log("read thread: error: %v", err)
				} else {
					s.log.PrintPoints([]model.PickUpPoint{point})
				}
			}
		}
	}
}

// CreatePoint creates a pick-up point.
func (s *Service) CreatePoint(point model.PickUpPoint) {
	s.write <- request{reqType: createReq, point: point}
}

// ListPoints prints a slice of all pick-up points.
func (s *Service) ListPoints() {
	s.read <- request{reqType: listReq}
}

// GetPoint prints a specified pick-up point.
func (s *Service) GetPoint(id uint64) {
	s.read <- request{reqType: getReq, id: id}
}

// UpdatePoint updates a pick-up point info.
func (s *Service) UpdatePoint(point model.PickUpPoint) {
	s.write <- request{reqType: updateReq, point: point}
}

// DeletePoint deletes a pick-up point.
func (s *Service) DeletePoint(id uint64) {
	s.write <- request{reqType: deleteReq, id: id}
}
