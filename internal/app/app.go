package app

import (
	"context"

	"github.com/iKOPKACtraxa/otus-hw/project_rotation/internal/logger"
	"github.com/iKOPKACtraxa/otus-hw/project_rotation/internal/storage"
)

type App struct {
	Logger  logger.Logger
	Storage storage.Storage
}

// New returns a new App main object.
func New(logger logger.Logger, storage storage.Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

// AddBanner adds banner into slot.
func (a *App) AddBanner(ctx context.Context, BannerID storage.ID, SlotID storage.ID) error {
	return a.Storage.AddBanner(ctx, BannerID, SlotID)
}

// DeleteEvent deletes event from storage.
func (a *App) DeleteBanner(ctx context.Context, BannerID storage.ID, SlotID storage.ID) error {
	return a.Storage.DeleteBanner(ctx, BannerID, SlotID)
}

// DeleteEvent deletes event from storage.
func (a *App) ClicksIncreasing(ctx context.Context, SlotID, BannerID, SocGroupID storage.ID) error {
	return a.Storage.ClicksIncreasing(ctx, SlotID, BannerID, SocGroupID)
}

// BannerSelection ...todo
func (a *App) BannerSelection(ctx context.Context, SlotID, SocGroupID storage.ID) (storage.ID, error) {
	return a.Storage.BannerSelection(ctx, SlotID, SocGroupID)
}

// Debug make a message in logger at Debug-level.
func (a *App) Debug(msg string) {
	a.Logger.Debug(msg)
}

// Info make a message in logger at Info-level.
func (a *App) Info(msg string) {
	a.Logger.Info(msg)
}

// Warn make a message in logger at Warn-level.
func (a *App) Warn(msg string) {
	a.Logger.Warn(msg)
}

// Error make a message in logger at Error-level.
func (a *App) Error(msg string) {
	a.Logger.Error(msg)
}
