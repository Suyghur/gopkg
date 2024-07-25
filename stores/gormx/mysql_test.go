//@File     mysql_test.go
//@Time     2024/7/18
//@Author   #Suyghur,

package gormx

import (
	"database/sql"
	"regexp"
	"testing"

	"gorm.io/gorm"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/yyxxgame/gopkg/internal/dbtest"
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
	//c := &trackedGormc{}
	//
	//c.redis = dbtest.CreateRedis(t)
	//assert.NotNil(t, c.redis)
	//
	//db, err := gorm.Open(mysql.New(mysql.Config{
	//	SkipInitializeWithVersion: true,
	//	Conn:                      db,
	//}), &gorm.Config{})
	//if err != nil {
	//	t.Fatalf("an error '%s' was not expected when initialize a gorm instance", err)
	//}
	//c.IGormc = NewNodeConn(gormDb, c.redis, WithExpiry(time.Second*10))

	return NewGormDBFromConn(conn)
}

func TestRaw(t *testing.T) {
	dbtest.RunTest(t, func(conn *sql.DB, mock sqlmock.Sqlmock) {
		db := mockGormDB(t, conn)

		resp := map[string]any{}
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `t_users` WHERE `t_users`.`id` = ? LIMIT 1")).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "suyghur"))

		sql.Result()

		db.Raw("SELECT * FROM `t_users` WHERE `t_users`.`id` = ? LIMIT 1", 1).Find(&resp)

		db.Raw("SELECT * FROM `t_users` WHERE `t_users`.`id` = ? LIMIT 1", 1).Find(&resp)

		db.Query("cacheKey", func(db *gorm.DB) {
		})

		//db.Update("cacheKey")
		//db.Delete("cacheKey", "")
		//db.Create()

		db.Exec("cacheKey", true, func(db *grom.DB) {
			return db.Raw("SELECT * FROM `t_users` WHERE `t_users`.`id` = ? LIMIT 1", 1).Find(&resp)
		})
		//err := c.Raw(&resp, "any", func(conn *gorm.DB, value interface{}) error {
		//	c.queryRowsValue = true
		//	return conn.Table("`t_users`").Where("`t_users`.`id` = ?", 1).Take(&resp).Error
		//})
		//assert.Nil(t, err)
		//assert.True(t, c.queryRowsValue)

		//bResp, _ := json.Marshal(resp)
		//value, _ := db.redis.Get("any")
		//assert.Equal(t, bResp, []byte(value))
	})
}
