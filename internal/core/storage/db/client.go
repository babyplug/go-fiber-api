package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm/clause"

	"go-fiber-api/internal/core/config"
	"go-fiber-api/internal/core/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbCon    *gorm.DB
	dbOnce   sync.Once
	entities = []interface{}{model.User{}}
)

const (
	DSNFormat = "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s"
)

type Client interface {
	AddError(err error) error
	Assign(attrs ...interface{}) (tx *gorm.DB)
	Association(column string) *gorm.Association
	Attrs(attrs ...interface{}) (tx *gorm.DB)
	AutoMigrate(dst ...interface{}) error
	Begin(opts ...*sql.TxOptions) *gorm.DB
	Clauses(conds ...clause.Expression) (tx *gorm.DB)
	Commit() *gorm.DB
	Connection(fc func(tx *gorm.DB) error) (err error)
	Count(count *int64) (tx *gorm.DB)
	Create(value interface{}) (tx *gorm.DB)
	CreateInBatches(value interface{}, batchSize int) (tx *gorm.DB)
	DB() (*sql.DB, error)
	Debug() (tx *gorm.DB)
	Delete(value interface{}, conds ...interface{}) (tx *gorm.DB)
	Distinct(args ...interface{}) (tx *gorm.DB)
	Exec(sql string, values ...interface{}) (tx *gorm.DB)
	Find(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	FindInBatches(dest interface{}, batchSize int, fc func(tx *gorm.DB, batch int) error) *gorm.DB
	First(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	FirstOrCreate(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	FirstOrInit(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	Get(key string) (interface{}, bool)
	Group(name string) (tx *gorm.DB)
	Having(query interface{}, args ...interface{}) (tx *gorm.DB)
	InnerJoins(query string, args ...interface{}) (tx *gorm.DB)
	InstanceGet(key string) (interface{}, bool)
	InstanceSet(key string, value interface{}) *gorm.DB
	Joins(query string, args ...interface{}) (tx *gorm.DB)
	Last(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	Limit(limit int) (tx *gorm.DB)
	MapColumns(m map[string]string) (tx *gorm.DB)
	Migrator() gorm.Migrator
	Model(value interface{}) (tx *gorm.DB)
	Not(query interface{}, args ...interface{}) (tx *gorm.DB)
	Offset(offset int) (tx *gorm.DB)
	Omit(columns ...string) (tx *gorm.DB)
	Or(query interface{}, args ...interface{}) (tx *gorm.DB)
	Order(value interface{}) (tx *gorm.DB)
	Pluck(column string, dest interface{}) (tx *gorm.DB)
	Preload(query string, args ...interface{}) (tx *gorm.DB)
	Raw(sql string, values ...interface{}) (tx *gorm.DB)
	Rollback() *gorm.DB
	RollbackTo(name string) *gorm.DB
	Row() *sql.Row
	Rows() (*sql.Rows, error)
	Save(value interface{}) (tx *gorm.DB)
	SavePoint(name string) *gorm.DB
	Scan(dest interface{}) (tx *gorm.DB)
	ScanRows(rows *sql.Rows, dest interface{}) error
	Scopes(funcs ...func(*gorm.DB) *gorm.DB) (tx *gorm.DB)
	Select(query interface{}, args ...interface{}) (tx *gorm.DB)
	Session(config *gorm.Session) *gorm.DB
	Set(key string, value interface{}) *gorm.DB
	SetupJoinTable(model interface{}, field string, joinTable interface{}) error
	Table(name string, args ...interface{}) (tx *gorm.DB)
	Take(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	ToSQL(queryFn func(tx *gorm.DB) *gorm.DB) string
	Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error)
	Unscoped() (tx *gorm.DB)
	Update(column string, value interface{}) (tx *gorm.DB)
	UpdateColumn(column string, value interface{}) (tx *gorm.DB)
	UpdateColumns(values interface{}) (tx *gorm.DB)
	Updates(values interface{}) (tx *gorm.DB)
	Use(plugin gorm.Plugin) error
	Where(query interface{}, args ...interface{}) (tx *gorm.DB)
	WithContext(ctx context.Context) *gorm.DB
}

func ProvideDB(cfg *config.Configuration) Client {
	dbOnce.Do(func() {
		fmt.Println(cfg.DBHost)
		dbURI := fmt.Sprintf(DSNFormat, cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)

		pgConfig := postgres.Config{
			DSN: dbURI,
		}

		var err error
		dbCon, err = gorm.Open(postgres.New(pgConfig))

		if err != nil {
			log.Fatalf("e: %v", err)
		}

		if err := autoMigrate(dbCon); err != nil {
			log.Fatalf("migrate schema failed: %v", err)
		}
	})

	return dbCon
}

func GetDbTestMode() (*gorm.DB, error) {
	name := uuid.New().String()
	dbCon, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%v?mode=memory&cache=shared", name)), &gorm.Config{})

	if err := dbCon.Migrator().DropTable(entities...); err != nil {
		return nil, err
	}

	if err := autoMigrate(dbCon); err != nil {
		return nil, err
	}

	return dbCon, err
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(entities...)
}
