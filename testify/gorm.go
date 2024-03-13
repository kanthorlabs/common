package testify

import (
	"log"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/kanthorlabs/common/mocks/clock"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GormStart(t *testing.T) *gorm.DB {
	u, err := url.Parse(testdata.SqliteUri)
	require.NoError(t, err)

	db, err := gorm.Open(sqlite.Open(u.RawPath), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				LogLevel: logger.Silent, // Log level
			},
		),
	})
	require.NoError(t, err)

	err = db.AutoMigrate(&testdata.User{})
	require.NoError(t, err)

	return db
}

func GormEnd(t *testing.T, db *gorm.DB) {
	conn, err := db.DB()
	require.NoError(t, err)
	err = conn.Close()
	require.NoError(t, err)
}

func GormInsert(t *testing.T, db *gorm.DB, count int) (map[string]*testdata.User, []string) {
	now := time.Now()
	watch := clock.NewClock(t)

	var ids []string
	rows := make(map[string]*testdata.User)
	for i := 0; i < count; i++ {
		watch.EXPECT().Now().Return(now.Add(-time.Minute * time.Duration(i))).Once()

		row := testdata.NewUser(watch)
		tx := db.Create(&row)
		require.NoError(t, tx.Error)

		ids = append(ids, row.Id)
		rows[row.Id] = &row
	}

	return rows, ids
}
