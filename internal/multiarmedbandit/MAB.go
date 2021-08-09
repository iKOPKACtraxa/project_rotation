package multiarmedbandit

import (
	"context"
	"database/sql"
	"fmt"
	"math"

	"github.com/iKOPKACtraxa/project_rotation/internal/storage"
)

// Select selects 1 banner (ID) for socGroupID.
func Select(ctx context.Context, tx *sql.Tx, socGroupID storage.ID) (storage.ID, error) {
	var rowsCount int
	row := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM Statistic WHERE SocGroupID=$1", socGroupID)
	err := row.Scan(&rowsCount)
	if err != nil {
		return 0, fmt.Errorf("counting of rows from Statistic with SocGroupID has got an error: %w", err)
	}

	rows, err := tx.QueryContext(ctx, "SELECT BannerID, Impressions, Clicks FROM Statistic WHERE SocGroupID=$1", socGroupID)
	if err != nil {
		return 0, fmt.Errorf("selecting of dataForAnalysis from Statistic has got an error: %w", err)
	}
	defer rows.Close()
	impressionsOfBanner := make(map[storage.ID]float64, rowsCount)
	clicksOfBanner := make(map[storage.ID]float64, rowsCount)
	var allImpressions float64
	for rows.Next() {
		var bannerID storage.ID
		var impressions float64
		var clicks float64
		err = rows.Scan(&bannerID, &impressions, &clicks)
		if err != nil {
			return 0, fmt.Errorf("scanning of IDs with SocGroupID from Statistic has got an error: %w", err)
		}
		impressionsOfBanner[bannerID] += impressions
		clicksOfBanner[bannerID] += clicks
		allImpressions += impressions
	}
	if err = rows.Err(); err != nil {
		return 0, fmt.Errorf("scanning of dataForAnalysis from Statistic has got an error: %w", err)
	}

	var selectedBanner storage.ID
	var maxTotalScoreOfBanner float64
	for k := range impressionsOfBanner {
		scoreOfBanner := clicksOfBanner[k]/impressionsOfBanner[k] + math.Sqrt(2*math.Log(allImpressions)/impressionsOfBanner[k])
		if scoreOfBanner > maxTotalScoreOfBanner {
			selectedBanner = k
			maxTotalScoreOfBanner = scoreOfBanner
		}
	}

	return selectedBanner, nil
}
