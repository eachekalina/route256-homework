package core

import "context"

func (s *pickUpPointCoreService) DeletePoint(ctx context.Context, id uint64) error {
	ctx, span := s.tracer.Start(ctx, "DeletePoint")
	defer span.End()

	err := s.pointService.DeletePoint(ctx, id)
	if err != nil {
		return err
	}
	s.cache.DeletePoint(ctx, id)
	return nil
}
