package pickuppointservice

import "homework/internal/model"

type storage interface {
	Close() error
	Create(point model.PickUpPoint) error
	List() []model.PickUpPoint
	Get(id uint64) (model.PickUpPoint, error)
	Update(point model.PickUpPoint) error
	Delete(id uint64) error
}

type request[T any, R any] struct {
	value  T
	result chan R
}

func newRequest[T any, R any](value T) request[T, R] {
	return request[T, R]{
		value:  value,
		result: make(chan R),
	}
}

type getResult struct {
	point model.PickUpPoint
	err   error
}

// PickUpPointService allows concurrent working on pick-up points.
type PickUpPointService struct {
	stor   storage
	create chan request[model.PickUpPoint, error]
	list   chan request[any, []model.PickUpPoint]
	get    chan request[uint64, getResult]
	update chan request[model.PickUpPoint, error]
	delete chan request[uint64, error]
	close  chan any
	saved  chan any
}

// NewPickUpPointService creates a new PickUpPointService.
func NewPickUpPointService(stor storage) *PickUpPointService {
	s := &PickUpPointService{
		stor:   stor,
		create: make(chan request[model.PickUpPoint, error]),
		list:   make(chan request[any, []model.PickUpPoint]),
		get:    make(chan request[uint64, getResult]),
		update: make(chan request[model.PickUpPoint, error]),
		delete: make(chan request[uint64, error]),
		close:  make(chan any),
		saved:  make(chan any),
	}
	s.run()
	return s
}

func (s *PickUpPointService) run() {
	go func() {
		for {
			select {
			case <-s.close:
				s.stor.Close()
				close(s.saved)
				return
			case r := <-s.create:
				r.result <- s.stor.Create(r.value)
				close(r.result)
			case r := <-s.update:
				r.result <- s.stor.Update(r.value)
				close(r.result)
			case r := <-s.delete:
				r.result <- s.stor.Delete(r.value)
				close(r.result)
			}
		}
	}()

	go func() {
		for {
			select {
			case <-s.close:
				return
			case r := <-s.list:
				r.result <- s.stor.List()
				close(r.result)
			case r := <-s.get:
				point, err := s.stor.Get(r.value)
				r.result <- getResult{
					point: point,
					err:   err,
				}
				close(r.result)
			}
		}
	}()
}

// Close stops all goroutines and saves storage.
func (s *PickUpPointService) Close() {
	close(s.close)
	<-s.saved
}

// CreatePoint creates a pick-up point.
func (s *PickUpPointService) CreatePoint(point model.PickUpPoint) error {
	req := newRequest[model.PickUpPoint, error](point)
	s.create <- req
	return <-req.result
}

// ListPoints returns a slice of all pick-up points.
func (s *PickUpPointService) ListPoints() []model.PickUpPoint {
	req := newRequest[any, []model.PickUpPoint](nil)
	s.list <- req
	return <-req.result
}

// GetPoint returns a specified pick-up point.
func (s *PickUpPointService) GetPoint(id uint64) (model.PickUpPoint, error) {
	req := newRequest[uint64, getResult](id)
	s.get <- req
	res := <-req.result
	return res.point, res.err
}

// UpdatePoint updates a pick-up point info.
func (s *PickUpPointService) UpdatePoint(point model.PickUpPoint) error {
	req := newRequest[model.PickUpPoint, error](point)
	s.update <- req
	return <-req.result
}

// DeletePoint deletes a pick-up point.
func (s *PickUpPointService) DeletePoint(id uint64) error {
	req := newRequest[uint64, error](id)
	s.delete <- req
	return <-req.result
}
