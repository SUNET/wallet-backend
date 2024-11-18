package apiv1

import (
	"context"
	apiv1_mockas "vc/internal/mockas/apiv1"
)

type EventPublisher interface {
	MockNext(mockNextRequest *apiv1_mockas.MockNextRequest) error
	Close(ctx context.Context) error
}
