package app

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/iKOPKACtraxa/project_rotation/internal/logger"
	"github.com/iKOPKACtraxa/project_rotation/internal/storage"
	"github.com/streadway/amqp"
)

type App struct {
	Logger     logger.Logger
	Storage    storage.Storage
	ChannelMQ  *amqp.Channel
	Confirms   chan amqp.Confirmation
	ExchangeMQ string
	RoutingKey string
}

// New returns a new App main object.
func New(ctx context.Context, logger logger.Logger, storage storage.Storage, uri, exchange, routingKey string, reliable bool) *App {
	var channelMQ *amqp.Channel
	var confirms chan amqp.Confirmation
	for channelMQ == nil {
		channelMQ, confirms = mqStart(ctx, uri, exchange, reliable, logger)
		if channelMQ == nil {
			logger.Info("connection awaiting...")
			time.Sleep(7 * time.Second)
		}
	}
	return &App{
		Logger:     logger,
		Storage:    storage,
		ChannelMQ:  channelMQ,
		Confirms:   confirms,
		ExchangeMQ: exchange,
		RoutingKey: routingKey,
	}
}

type eventToMQ struct {
	SlotID     storage.ID `json:"slotID"`
	BannerID   storage.ID `json:"bannerID"`
	SocGroupID storage.ID `json:"socGroupID"`
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
	err := a.Storage.ClicksIncreasing(ctx, slotID, bannerID, socGroupID)
	if err != nil {
		return fmt.Errorf("ClicksIncreasing err: %w", err)
	}
	event := eventToMQ{
		SlotID:     slotID,
		BannerID:   bannerID,
		SocGroupID: socGroupID,
	}
	err = a.sendToMQ(event)
	if err != nil {
		return fmt.Errorf("sendToMQ err: %w", err)
	}
	return nil
}

// BannerSelection selects banner using multiarmedbandit logic and increases impressions.
func (a *App) BannerSelection(ctx context.Context, slotID, socGroupID storage.ID) (storage.ID, error) {
	id, err := a.Storage.BannerSelection(ctx, slotID, socGroupID)
	if err != nil {
		return 0, fmt.Errorf("BannerSelection err: %w", err)
	}
	event := eventToMQ{
		SlotID:     slotID,
		BannerID:   id,
		SocGroupID: socGroupID,
	}
	err = a.sendToMQ(event)
	if err != nil {
		return 0, fmt.Errorf("sendToMQ err: %w", err)
	}
	return id, nil
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

func mqStart(ctx context.Context, uri, exchange string, reliable bool, logger logger.Logger) (*amqp.Channel, chan amqp.Confirmation) {
	logger.Info("MQ: dialing ", uri)
	connection, err := amqp.Dial(uri)
	if err != nil {
		logger.Error("dial err: ", err)
		return nil, nil
	}
	go func() {
		<-ctx.Done()
		connection.Close()
	}()
	logger.Info("MQ: got Connection, getting Channel")
	channel, err := connection.Channel()
	if err != nil {
		logger.Error("channel err: ", err)
	}
	logger.Info("MQ: got Channel, declaring direct Exchange: ", exchange)
	if err := channel.ExchangeDeclare(
		exchange, // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // noWait
		nil,      // arguments
	); err != nil {
		logger.Error("exchange Declare: ", err)
	}
	var confirms chan amqp.Confirmation
	if reliable {
		logger.Info("MQ: enabling publishing confirms")
		if err := channel.Confirm(false); err != nil {
			logger.Error("Channel could not be put into confirm mode: ", err)
		}
		confirms = channel.NotifyPublish(make(chan amqp.Confirmation, 1))
	}
	return channel, confirms
}

func (a *App) confirmOne() {
	a.Logger.Info("MQ: waiting for confirmation of one publishing")
	if confirmed := <-a.Confirms; confirmed.Ack {
		a.Logger.Info("MQ: confirmed delivery with delivery tag:", confirmed.DeliveryTag)
	} else {
		a.Logger.Info("MQ: failed delivery of delivery tag:", confirmed.DeliveryTag)
	}
}

func (a *App) sendToMQ(event eventToMQ) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal err: %w", err)
	}
	a.Logger.Debug("MQ: declared Exchange, publishing ", len(body), "B body (", body, ")")
	if err = a.ChannelMQ.Publish(
		a.ExchangeMQ,
		a.RoutingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            body,
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
		},
	); err != nil {
		return fmt.Errorf("exchange publish err: %w", err)
	}
	if a.Confirms != nil {
		a.confirmOne()
	}
	return nil
}
