package app

import (
	"context"

	"github.com/iKOPKACtraxa/project_rotation/internal/logger"
	"github.com/iKOPKACtraxa/project_rotation/internal/storage"
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
func (a *App) AddBanner(ctx context.Context, bannerID storage.ID, slotID storage.ID) error {
	return a.Storage.AddBanner(ctx, bannerID, slotID)
}

// DeleteEvent deletes event from storage.
func (a *App) DeleteBanner(ctx context.Context, bannerID storage.ID, slotID storage.ID) error {
	return a.Storage.DeleteBanner(ctx, bannerID, slotID)
}

// DeleteEvent deletes event from storage.
func (a *App) ClicksIncreasing(ctx context.Context, slotID, bannerID, socGroupID storage.ID) error {
	return a.Storage.ClicksIncreasing(ctx, slotID, bannerID, socGroupID)
}

// BannerSelection ...todo.
func (a *App) BannerSelection(ctx context.Context, slotID, socGroupID storage.ID) (storage.ID, error) {
	return a.Storage.BannerSelection(ctx, slotID, socGroupID)
}

// Debug make a message in logger at Debug-level.
func (a *App) Debug(args ...interface{}) {
	a.Logger.Debug(args)
}

// Info make a message in logger at Info-level.
func (a *App) Info(args ...interface{}) {
	a.Logger.Info(args)
}

// Warn make a message in logger at Warn-level.
func (a *App) Warn(args ...interface{}) {
	a.Logger.Warn(args)
}

// Error make a message in logger at Error-level.
func (a *App) Error(args ...interface{}) {
	a.Logger.Error(args)
}
