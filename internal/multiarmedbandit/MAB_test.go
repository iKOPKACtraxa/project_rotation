package multiarmedbandit

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/iKOPKACtraxa/project_rotation/internal/storage"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	dbtestname     string     = "rotationdbtest"
	connStrinit    string     = "user=postgres sslmode=disable"
	connStr        string     = "user=postgres sslmode=disable" + " dbname=" + dbtestname
	someSocGroup   storage.ID = 1
	someSlotID     storage.ID = 1
	scoreForWinner int        = 33 // at start has 10 clicks as all
	scoreForLooser int        = 15 // at start has 10 clicks as all
	winner         storage.ID = 1  // in statistic has 3 clicks
	looser         storage.ID = 3  // in statistic has 1 clicks
)

func TestSelect(t *testing.T) {
	db := start(t)
	defer unMigration(t)
	defer func() {
		err := db.Close()
		assert.NoErrorf(t, err, "db.Close: ", err)
	}()

	for i := 0; i < 50; i++ {
		tx, err := db.BeginTx(context.Background(), nil)
		require.NoErrorf(t, err, "db.BeginTx: ", err)
		id, err := Select(context.Background(), tx, someSocGroup)
		_, err = afterSelectLogic(tx, id, err)
		require.NoErrorf(t, err, "Select or afterSelectLogic has got an err: ", err)
	}
	var count int
	t.Run("check for best bunner", func(t *testing.T) {
		row := db.QueryRow("SELECT impressions FROM statistic WHERE bannerid=$1 AND slotid=1 AND socgroupid=1", winner)
		err := row.Scan(&count)
		require.NoErrorf(t, err, "row.Scan has got an err: ", err)
		require.Greaterf(t, count, scoreForWinner, "winner has not required amount, has:%v, need:%v", count, scoreForWinner)
	})
	t.Run("check for looser still got impressions", func(t *testing.T) {
		row := db.QueryRow("SELECT impressions FROM statistic WHERE bannerid=$1 AND slotid=1 AND socgroupid=1", looser)
		err := row.Scan(&count)
		require.NoErrorf(t, err, "row.Scan has got an err: ", err)
		require.Greaterf(t, count, scoreForLooser, "looser has not required amount, has:%v, need:%v", count, scoreForLooser)
	})
}

func start(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", connStrinit)
	require.NoErrorf(t, err, "open connection err: ", err)
	_, err = db.Exec("CREATE DATABASE " + dbtestname)
	require.NoErrorf(t, err, "CREATE DATABASE err: ", err)
	err = db.Close()
	require.NoErrorf(t, err, "close connection err: ", err)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		unMigration(t)
	}
	require.NoErrorf(t, err, "sql.Open: ", err)
	err = goose.Up(db, "../../migrations")
	require.NoErrorf(t, err, "first goose up: ", err)
	err = goose.Up(db, "../../migrations/fortest")
	require.NoErrorf(t, err, "second goose up: ", err)
	return db
}

func afterSelectLogic(tx *sql.Tx, idIn storage.ID, errIn error) (id storage.ID, err error) {
	if errIn != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return 0, fmt.Errorf("rollback err: %w", errRollback)
		}
		return 0, fmt.Errorf("select is not complete with err: %w", errIn)
	}
	_, err = tx.ExecContext(context.Background(),
		"UPDATE Statistic SET Impressions = Impressions + 1 WHERE SlotID=$1 AND BannerID=$2 AND SocGroupID=$3", someSlotID, idIn, someSocGroup)
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return 0, fmt.Errorf("rollback err: %w", errRollback)
		}
		return 0, fmt.Errorf("increasing of impressions has got an error: %w", errIn)
	}
	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("commit of BannerSelection transaction has got an error: %w", errIn)
	}
	return idIn, nil
}

func unMigration(t *testing.T) {
	db, err := sql.Open("postgres", connStrinit)
	require.NoErrorf(t, err, "open connection err: ", err)
	_, err = db.Exec("DROP DATABASE " + dbtestname)
	require.NoErrorf(t, err, "DROP DATABASE err: ", err)
	err = db.Close()
	require.NoErrorf(t, err, "close connection err: ", err)
}
