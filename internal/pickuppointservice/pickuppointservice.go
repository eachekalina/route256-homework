package pickuppointservice

import (
	"context"
	"homework/internal/logger"
	"homework/internal/model"
	"sync"
)

type storage interface {
	Close() error
	Create(point model.PickUpPoint) error
	List() []model.PickUpPoint
	Get(id uint64) (model.PickUpPoint, error)
	Update(point model.PickUpPoint) error
	Delete(id uint64) error
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

// PickUpPointService allows concurrent working on pick-up points.
type PickUpPointService struct {
	stor  storage
	log   *logger.Logger
	read  chan request
	write chan request
	ctx   context.Context
	close context.CancelFunc
	wg    sync.WaitGroup
}

// NewPickUpPointService creates a new PickUpPointService.
func NewPickUpPointService(ctx context.Context, stor storage, log *logger.Logger) *PickUpPointService {
	ctx, cancel := context.WithCancel(ctx)
	s := &PickUpPointService{
		stor:  stor,
		log:   log,
		read:  make(chan request),
		write: make(chan request),
		ctx:   ctx,
		close: cancel,
	}
	s.run(ctx)
	return s
}

func (s *PickUpPointService) run(ctx context.Context) {
	s.wg = sync.WaitGroup{}
	s.wg.Add(2)

	go func() {
		defer s.wg.Done()
		for {
			s.log.Log("write thread: waiting for request")
			select {
			case <-ctx.Done():
				s.log.Log("write thread: closing")
				return
			case r := <-s.write:
				switch r.reqType {
				case createReq:
					s.log.Log("write thread: requested create")
					err := s.stor.Create(r.point)
					if err != nil {
						s.log.Log("write thread: error: %v", err)
					}
				case updateReq:
					s.log.Log("write thread: requested update")
					err := s.stor.Update(r.point)
					if err != nil {
						s.log.Log("write thread: error: %v", err)
					}
				case deleteReq:
					s.log.Log("write thread: requested delete")
					err := s.stor.Delete(r.id)
					if err != nil {
						s.log.Log("write thread: error: %v", err)
					}
				default:
					s.log.Log("write thread: invalid request")
				}
			}
		}
	}()

	go func() {
		defer s.wg.Done()
		for {
			s.log.Log("read thread: waiting for request")
			select {
			case <-ctx.Done():
				s.log.Log("read thread: closing")
				return
			case r := <-s.read:
				switch r.reqType {
				case listReq:
					s.log.Log("read thread: requested list")
					list := s.stor.List()
					s.log.PrintPoints(list)
				case getReq:
					s.log.Log("read thread: requested get")
					point, err := s.stor.Get(r.id)
					if err != nil {
						s.log.Log("read thread: error: %v", err)
					} else {
						s.log.PrintPoints([]model.PickUpPoint{point})
					}
				}
			}
		}
	}()
}

// Close stops all goroutines.
func (s *PickUpPointService) Close() {
	s.close()
	s.wg.Wait()
}

// CreatePoint creates a pick-up point.
func (s *PickUpPointService) CreatePoint(point model.PickUpPoint) {
	s.write <- request{reqType: createReq, point: point}
}

// ListPoints prints a slice of all pick-up points.
func (s *PickUpPointService) ListPoints() {
	s.read <- request{reqType: listReq}
}

// GetPoint prints a specified pick-up point.
func (s *PickUpPointService) GetPoint(id uint64) {
	s.read <- request{reqType: getReq, id: id}
}

// UpdatePoint updates a pick-up point info.
func (s *PickUpPointService) UpdatePoint(point model.PickUpPoint) {
	s.write <- request{reqType: updateReq, point: point}
}

// DeletePoint deletes a pick-up point.
func (s *PickUpPointService) DeletePoint(id uint64) {
	s.write <- request{reqType: deleteReq, id: id}
}
