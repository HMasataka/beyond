package infrastructure

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/HMasataka/beyond/config"
	"github.com/HMasataka/transactor/rdbms"
	"github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
)

var (
	db     rdbms.Conn
	dbOnce sync.Once
)

func newConnection(cfg *config.Config) rdbms.Conn {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}

	c := mysql.Config{
		DBName:               cfg.MySQL.DB,
		User:                 cfg.MySQL.User,
		Passwd:               cfg.MySQL.Password,
		Addr:                 fmt.Sprintf("%s:%d", cfg.MySQL.Host, cfg.MySQL.Port),
		Net:                  cfg.MySQL.Net,
		ParseTime:            true,
		Collation:            "utf8mb4_unicode_ci",
		Loc:                  jst,
		AllowNativePasswords: true,
	}

	log.Info().Str("DSN", c.FormatDSN()).Send()

	conn, err := sql.Open("mysql", c.FormatDSN())
	if err != nil {
		panic(err)
	}

	return conn
}

func NewConnectionOnce(conf *config.Config) rdbms.Conn {
	dbOnce.Do(func() {
		db = newConnection(conf)
	})

	return db
}
