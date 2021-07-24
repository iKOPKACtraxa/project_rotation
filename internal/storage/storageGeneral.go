package storage

import (
	"context"
)

type (
	ID int
)

type Storage interface {
	AddBanner(ctx context.Context, BannerID ID, SlotID ID) error
	DeleteBanner(ctx context.Context, BannerID ID, SlotID ID) error
	ClicksIncreasing(ctx context.Context, SlotID ID, BannerID ID, SocGroupID ID) error
	BannerSelection(ctx context.Context, SlotID, SocGroupID ID) (ID, error)
}
