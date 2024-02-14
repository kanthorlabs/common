package testify

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/kanthorlabs/common/faker/timer"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GormStart(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(testdata.SqliteUri), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				LogLevel: logger.Silent, // Log level
			},
		),
	})
	require.Nil(t, err)

	err = db.AutoMigrate(&testdata.User{})
	require.Nil(t, err)

	return db
}

func GormEnd(t *testing.T, db *gorm.DB) {
	conn, err := db.DB()
	require.Nil(t, err)
	err = conn.Close()
	require.Nil(t, err)
}

func GormInsert(t *testing.T, db *gorm.DB, count int) (map[string]*testdata.User, []string) {
	now := time.Now()
	clock := timer.NewTimer(t)

	var ids []string
	rows := make(map[string]*testdata.User)
	for i := 0; i < count; i++ {
		clock.EXPECT().Now().Return(now.Add(-time.Minute * time.Duration(i))).Once()

		row := testdata.NewUser(clock)
		tx := db.Create(row)
		require.Nil(t, tx.Error)

		ids = append(ids, row.Id)
		rows[row.Id] = row
	}

	return rows, ids
}
