//@File     mysql_test.go
//@Time     2024/7/18
//@Author   #Suyghur,

package gormx

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/yyxxgame/gopkg/internal/dbtest"
	"gorm.io/gorm"
)

type (
	trackedDb struct {
		//redis *redis.Redis
		//execValue      bool
		//queryRowsValue bool
		//transactValue  bool
	}
)

func mockGormDB(t *testing.T, conn *sql.DB) *gorm.DB {
	return NewGormDBFromConn(conn)
}

func TestRawSqlQuery(t *testing.T) {
	dbtest.RunTest(t, func(conn *sql.DB, mock sqlmock.Sqlmock) {
		db := mockGormDB(t, conn)

		resp := make(map[string]any)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT `name` FROM `t_users` WHERE `id` = ? LIMIT 1")).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "suyghur"))

		err := db.Raw("SELECT `name` FROM `t_users` WHERE `id` = ? LIMIT 1", 1).Find(&resp).Error

		assert.Nil(t, err)
		assert.Equal(t, "suyghur", resp["name"])
	})
}

func TestRawSqlInsert(t *testing.T) {
	dbtest.RunTest(t, func(conn *sql.DB, mock sqlmock.Sqlmock) {
		db := mockGormDB(t, conn)

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `t_users` (`username`.`email`) VALUES (?, ?)")).
			WithArgs("suyghur", "suyghurmjp@outlook.com").
			WillReturnResult(sqlmock.NewResult(1, 1))

		resp := db.Exec("INSERT INTO `t_users` (`username`.`email`) VALUES (?, ?)", "suyghur", "suyghurmjp@outlook.com")
		assert.Equal(t, int64(1), resp.RowsAffected)
	})
}

func TestRawSqlUpdate(t *testing.T) {
	dbtest.RunTest(t, func(conn *sql.DB, mock sqlmock.Sqlmock) {
		db := mockGormDB(t, conn)

		mock.ExpectExec(regexp.QuoteMeta("UPDATE `t_users` SET `username` = ? WHERE `id` = ?")).
			WithArgs("suyghur", 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		resp := db.Exec("UPDATE `t_users` SET `username` = ? WHERE `id` = ?", "suyghur", 1)
		assert.Equal(t, int64(1), resp.RowsAffected)
	})
}

func TestRawSqlDelete(t *testing.T) {
	dbtest.RunTest(t, func(conn *sql.DB, mock sqlmock.Sqlmock) {
		db := mockGormDB(t, conn)

		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `t_users` WHERE `username` = ?")).
			WithArgs("suyghur").
			WillReturnResult(sqlmock.NewResult(1, 1))

		resp := db.Exec("DELETE FROM `t_users` WHERE `username` = ?", "suyghur")
		assert.Equal(t, int64(1), resp.RowsAffected)
	})
}

func TestRawSqlError(t *testing.T) {
	dbtest.RunTest(t, func(conn *sql.DB, mock sqlmock.Sqlmock) {
		db := mockGormDB(t, conn)

		resp := make(map[string]any)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT `name` FROM `t_users` WHERE `id` = ? LIMIT 1")).
			WithArgs(1).
			WillReturnError(errors.New("test error"))

		err := db.Raw("SELECT `name` FROM `t_users` WHERE `id` = ? LIMIT 1", 1).Find(&resp).Error

		assert.NotNil(t, err)
	})
}
