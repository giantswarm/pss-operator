package test

import (
	"context"
)

func (r *Handler) EnsureDeleted(ctx context.Context, obj interface{}) error {
	r.logger.Debugf(ctx, "KUBA: EnsureCreated: %+v", obj)
	return nil
}
