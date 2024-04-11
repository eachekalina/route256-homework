package core

import "context"

func (s *pickUpPointCoreService) DeletePoint(ctx context.Context, id uint64) error {
	return s.pointService.DeletePoint(ctx, id)
}
