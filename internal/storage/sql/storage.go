package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/iKOPKACtraxa/project_rotation/internal/logger"
	"github.com/iKOPKACtraxa/project_rotation/internal/multiarmedbandit"
	"github.com/iKOPKACtraxa/project_rotation/internal/storage"
	"github.com/lib/pq"
)

const (
	isExistErrCode = "23505"
)

var ErrBannerinSlotIsNotExist = errors.New("in slot there is no banners")

type StorageInDB struct {
	db   *sql.DB
	logg logger.Logger
}

// New returns a StorageInDB.
func New(ctx context.Context, connStr string, logg logger.Logger) (*StorageInDB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	s := &StorageInDB{
		db:   db,
		logg: logg,
	}
	go s.Close(ctx, logg)
	return s, nil
}

// Close closes connection to DB.
func (s *StorageInDB) Close(ctx context.Context, logg logger.Logger) {
	<-ctx.Done()
	logg.Info("connection to DB is closing gracefully...")
	if err := s.db.Close(); err != nil {
		logg.Error("failed to stop DB:", err)
	} else {
		logg.Info("connection to DB is closed gracefully...")
	}
}

// AddBanner adds banner into slot by adding a row in BannersInSlots.
// If it already exists in slot transaction is ends.
// Also AddBanner creates in Statistic rows for every socGroup:
// with 1 Impressions and 0 Clicks.
func (s *StorageInDB) AddBanner(ctx context.Context, bannerID, slotID storage.ID) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning of AddBanner transaction has got an error: %w", err)
	}
	_, err = tx.ExecContext(ctx, "INSERT INTO BannersInSlots (BannerID, SlotID) VALUES ($1, $2)",
		bannerID,
		slotID,
	)
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case isExistErrCode:
			s.logg.Info("Banner", bannerID, "already added to slot", slotID)
			err = tx.Rollback()
			if err != nil {
				s.logg.Error("rollback is not complete wit err:", err)
			}
			return nil
		default:
			err = tx.Rollback()
			if err != nil {
				s.logg.Error("rollback is not complete wit err:", err)
			}
			return fmt.Errorf("addition of banner in slot has got an error: %w", err)
		}
	}

	var socGroups []storage.ID
	rows, err := tx.QueryContext(ctx, "SELECT ID FROM SocGroups")
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			s.logg.Error("rollback is not complete wit err:", err)
		}
		return fmt.Errorf("selecting of IDs from SocGroups has got an error: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id storage.ID
		err = rows.Scan(&id)
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				s.logg.Error("rollback is not complete wit err:", err)
			}
			return fmt.Errorf("scanning of IDs from SocGroups has got an error: %w", err)
		}
		socGroups = append(socGroups, id)
	}
	if err = rows.Err(); err != nil {
		err = tx.Rollback()
		if err != nil {
			s.logg.Error("rollback is not complete wit err:", err)
		}
		return fmt.Errorf("scanning of IDs from SocGroups has got an error: %w", err)
	}

	for _, socGroup := range socGroups {
		_, err = tx.ExecContext(ctx, "INSERT INTO Statistic (SlotID, BannerID, SocGroupID, Impressions, Clicks) VALUES ($1, $2, $3, $4, $5)",
			slotID,
			bannerID,
			socGroup,
			1,
			0,
		)
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				s.logg.Error("rollback is not complete wit err:", err)
			}
			return fmt.Errorf("new banner addition to Statistic has got an error: %w", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit of AddBanner transaction has got an error: %w", err)
	}
	s.logg.Info("Banner", bannerID, "added to slot", slotID)
	return nil
}

// DeleteBanner deletes banner from slot by removing a row from BannersInSlots.
// If it not exists in slot ErrBannerinSlotIsNotExist returned.
// Also DeleteBanner deletes in Statistic rows for every socGroup.
func (s *StorageInDB) DeleteBanner(ctx context.Context, bannerID, slotID storage.ID) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin of DeleteBanner transaction has got an error: %w", err)
	}
	_, err = tx.ExecContext(ctx,
		"DELETE FROM BannersInSlots WHERE BannerID=$1 AND SlotID=$2", bannerID, slotID)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			s.logg.Error("rollback is not complete wit err:", err)
		}
		return fmt.Errorf("deleting of banner in slot has got an error: %w", err)
	}
	_, err = tx.ExecContext(ctx,
		"DELETE FROM Statistic WHERE BannerID=$1 AND SlotID=$2", bannerID, slotID)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			s.logg.Error("rollback is not complete wit err:", err)
		}
		return fmt.Errorf("deleting of rows from Statistic has got an error: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit of DeleteBanner transaction has got an error: %w", err)
	}
	s.logg.Info("Banner", bannerID, "deleted from slot", slotID)
	return nil
}

// ClicksIncreasing adds 1 to Clicks in Statistic table.
func (s *StorageInDB) ClicksIncreasing(ctx context.Context, slotID, bannerID, socGroupID storage.ID) error {
	_, err := s.db.ExecContext(ctx,
		"UPDATE Statistic SET Clicks = Clicks + 1 WHERE SlotID=$1 AND BannerID=$2 AND SocGroupID=$3", slotID, bannerID, socGroupID)
	if err != nil {
		return fmt.Errorf("ClicksIncreasing has got an error: %w", err)
	}
	s.logg.Info("Clicks for slot:", slotID, ", banner:", bannerID, ", socGroup:", socGroupID, " increased")
	return nil
}

// BannerSelection selects 1 banner (ID) for slotID and socGroupID.
// Banner must be added in Slots in BannersInSlots, otherwise there is no banners to choose.
func (s *StorageInDB) BannerSelection(ctx context.Context, slotID, socGroupID storage.ID) (storage.ID, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin of BannerSelection transaction has got an error: %w", err)
	}

	err = s.isExistInBannersInSlotsCheck(ctx, slotID)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			s.logg.Error("rollback is not complete wit err:", err)
		}
		return 0, err
	}

	selectedBanner, err := multiarmedbandit.Select(ctx, tx, socGroupID)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			s.logg.Error("rollback is not complete wit err:", err)
		}
		return 0, err
	}

	_, err = s.db.ExecContext(ctx,
		"UPDATE Statistic SET Impressions = Impressions + 1 WHERE SlotID=$1 AND BannerID=$2 AND SocGroupID=$3", slotID, selectedBanner, socGroupID)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			s.logg.Error("rollback is not complete wit err:", err)
		}
		return 0, fmt.Errorf("ClicksIncreasing has got an error: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("commit of BannerSelection transaction has got an error: %w", err)
	}
	s.logg.Info("Selected banner:", selectedBanner, " for slot:", slotID, " and socGroup:", socGroupID)
	return selectedBanner, nil
}

// To think: Method BannerSelection works in transaction. If it must work faster it might use
// some kind of precalculated pool of ScoreOfBanner for every bannerID slotID and socGroupID.
// This pool must be updated for needed period of time.

// isExistInBannersInSlotsCheck checks wether banner in slot is exist.
func (s *StorageInDB) isExistInBannersInSlotsCheck(ctx context.Context, slotID storage.ID) error {
	row := s.db.QueryRowContext(ctx, "SELECT BannerID FROM BannersInSlots WHERE BannerID=ANY(SELECT ID FROM Banners) AND SlotID=$1 LIMIT 1", slotID)
	var bannerInSlot string
	err := row.Scan(&bannerInSlot)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrBannerinSlotIsNotExist
	}
	if err != nil {
		return err
	}
	return nil
}

// To think: social groups are known at start, and new groups is not adding. If so, when new social group is added it needs to add new rows in Statistic with all SlotID BannerID variations and Impressions=1 Clicks=0
