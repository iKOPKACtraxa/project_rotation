package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"

	"github.com/iKOPKACtraxa/otus-hw/project_rotation/internal/logger"
	"github.com/iKOPKACtraxa/otus-hw/project_rotation/internal/storage"
	"github.com/lib/pq"
)

var (
	ErrBannerInSlotIsExist      = errors.New("banner in slot is already exists")
	ErrBannerinSlotIsNotExist   = errors.New("banner in slot is not exists")
	ErrKeyinStatisticIsNotExist = errors.New("key in Statistic is not exists, it must have been created by BannerSelection")
)

type StorageInDB struct {
	DB *sql.DB
}

// New returns a StorageInDB.
func New(ctx context.Context, connStr string, logg logger.Logger) (*StorageInDB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	s := &StorageInDB{
		DB: db,
	}
	go s.Close(ctx, logg)
	return s, nil
}

// Close closes connection to DB.
func (s *StorageInDB) Close(ctx context.Context, logg logger.Logger) {
	<-ctx.Done()
	if err := s.DB.Close(); err != nil {
		logg.Error("failed to stop DB: " + err.Error())
	}
}

// AddBanner adds banner into slot by adding a row in BannersInSlots.
// If it already exists in slot ErrBannerInSlotIsExist returned.
// Also AddBanner creates in Statistic rows for every socGroup with 1 Impressions
// and 0 Clicks.
func (s *StorageInDB) AddBanner(ctx context.Context, BannerID, SlotID storage.ID) error {
	_, err := s.DB.ExecContext(ctx, "INSERT INTO BannersInSlots (BannerID, SlotID) VALUES ($1, $2)",
		BannerID,
		SlotID,
	)
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505":
			return ErrBannerInSlotIsExist
		default:
			return fmt.Errorf("creating of banner in slot has got an error: %w", err)
		}
	}

	var socGroups []storage.ID
	rows, err := s.DB.QueryContext(ctx, "SELECT ID FROM SocGroups")
	if err != nil {
		return fmt.Errorf("selecting of IDs from SocGroups has got an error: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id storage.ID
		err = rows.Scan(&id)
		socGroups = append(socGroups, id)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("scanning of IDs from SocGroups has got an error: %w", err)
	}

	for _, socGroup := range socGroups {
		_, err = s.DB.ExecContext(ctx, "INSERT INTO Statistic (SlotID, BannerID, SocGroupID, Impressions, Clicks) VALUES ($1, $2, $3, $4, $5)",
			SlotID,
			BannerID,
			socGroup,
			1,
			0,
		)
	}
	return nil
}

// DeleteBanner deletes banner from slot by removing a row from BannersInSlots.
// If it not exists in slot ErrBannerinSlotIsNotExist returned.
// Also DeleteBanner deletes in Statistic rows for every socGroup.
func (s *StorageInDB) DeleteBanner(ctx context.Context, BannerID, SlotID storage.ID) error {
	err := s.isExistInBannersInSlotsCheck(ctx, BannerID, SlotID)
	if err != nil {
		return err
	}
	_, err = s.DB.ExecContext(ctx,
		"DELETE FROM BannersInSlots WHERE BannerID=$1 AND SlotID=$2", BannerID, SlotID)
	if err != nil {
		return fmt.Errorf("deleting of banner in slot has got an error: %w", err)
	}
	_, err = s.DB.ExecContext(ctx,
		"DELETE FROM Statistic WHERE BannerID=$1 AND SlotID=$2", BannerID, SlotID)
	if err != nil {
		return fmt.Errorf("deleting of rows from Statistic has got an error: %w", err)
	}
	return nil
}

// isExistInBannersInSlotsCheck checks wether banner in slot is exist.
func (s *StorageInDB) isExistInBannersInSlotsCheck(ctx context.Context, BannerID, SlotID storage.ID) error {
	var bannerInSlot string
	row := s.DB.QueryRowContext(ctx, "SELECT * FROM BannersInSlots WHERE BannerID=$1 AND SlotID=$2", BannerID, SlotID)
	if errors.Is(row.Scan(&bannerInSlot), sql.ErrNoRows) {
		return ErrBannerinSlotIsNotExist
	}
	return nil
}

// ClicksIncreasing adds 1 to Clicks in Statistic table.
func (s *StorageInDB) ClicksIncreasing(ctx context.Context, SlotID, BannerID, SocGroupID storage.ID) error {
	err := s.isExistInStatisticCheck(ctx, SlotID, BannerID, SocGroupID)
	if err != nil {
		return err
	}
	_, err = s.DB.ExecContext(ctx,
		"UPDATE Statistic SET Clicks = Clicks + 1 WHERE SlotID=$1 AND BannerID=$2 AND SocGroupID=$3", SlotID, BannerID, SocGroupID)
	if err != nil {
		return fmt.Errorf("ClicksIncreasing has got an error: %w", err)
	}
	return nil
}

// isExistInStatisticCheck checks wether key of SlotID, BannerID, SocGroupID is exist.
func (s *StorageInDB) isExistInStatisticCheck(ctx context.Context, SlotID, BannerID, SocGroupID storage.ID) error {
	var keyInTable string
	row := s.DB.QueryRowContext(ctx, "SELECT * FROM Statistic WHERE SlotID=$1 AND BannerID=$2 AND SocGroupID=$3", SlotID, BannerID, SocGroupID)
	if errors.Is(row.Scan(&keyInTable), sql.ErrNoRows) {
		return ErrKeyinStatisticIsNotExist
	}
	return nil
}

// BannerSelection is ...todo
func (s *StorageInDB) BannerSelection(ctx context.Context, SlotID, SocGroupID storage.ID) (storage.ID, error) {
	// 1. Выбрать строки с соцгруппой
	var impressionsOfBanner map[storage.ID]float64
	var clicksOfBanner map[storage.ID]float64
	var allImpressions float64
	rows, err := s.DB.QueryContext(ctx, "SELECT BannerID Impressions Clicks FROM Statistic WHERE SocGroupID=$1", SocGroupID)
	if err != nil {
		return 0, fmt.Errorf("selecting of dataForAnalysis from Statistic has got an error: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var BannerID storage.ID
		var Impressions float64
		var Clicks float64
		err = rows.Scan(&BannerID, &Impressions, &Clicks)
		impressionsOfBanner[BannerID] = impressionsOfBanner[BannerID] + Impressions
		clicksOfBanner[BannerID] = clicksOfBanner[BannerID] + Clicks
		allImpressions = allImpressions + Impressions
	}
	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("scanning of IDs from SocGroups has got an error: %w", err)
	}

	// 2. Посчитать для каждого величину по формуле (обе составляющие в рамках сумм по всем слотам)
	var totalScoreOfBanner map[storage.ID]float64
	for k := range impressionsOfBanner {
		totalScoreOfBanner[k] = clicksOfBanner[k]/impressionsOfBanner[k] + math.Sqrt(2*math.Log(allImpressions)/impressionsOfBanner[k])
	}

	// 3. Выдать баннер с максимальной величиной totalScoreOfBanner

	// 4. Увеличить показы

	return 0, nil
}

// To think: social groups are known at start, and new groups is not adding. If so, when new social group is added it needs to add new rows in Statistic with all SlotID BannerID variations and Impressions=1 Clicks=0
