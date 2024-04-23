package core

import "context"

func (s *pickUpPointCoreService) DeletePoint(ctx context.Context, id uint64) error {
	err := s.pointService.DeletePoint(ctx, id)
	if err != nil {
		return err
	}
	s.cache.DeletePoint(id)
	return nil
}
