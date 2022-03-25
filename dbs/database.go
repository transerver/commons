package dbs

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gookit/color"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/transerver/commons/configs"
	"github.com/transerver/commons/logger"
	"strings"
	"sync"
)

var (
	dbs   = make(map[string]*Database)
	mutex sync.Mutex
)

type Database struct {
	*sqlx.DB
	config *configs.DBConfig
	Logger *logger.Logger
	Hook   DatabaseHook
}

type Option func(db *Database)

func WithLogger(logger *logger.Logger) Option {
	return func(db *Database) {
		db.Logger = logger
	}
}

func WithHook(hook DatabaseHook) Option {
	return func(db *Database) {
		db.Hook = hook
	}
}

func WithConfig(cfg *configs.DBConfig) Option {
	return func(db *Database) {
		db.config = cfg
	}
}

func NewDatabase(opts ...Option) *Database {
	db := &Database{}
	db.getOpts(opts...)
	return db
}

func (db *Database) getOpts(opts ...Option) {
	for _, o := range opts {
		o(db)
	}
}

// NamedQuery using this DB.
// Any named placeholder parameters are replaced with fields from arg.
func (db *Database) NamedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	querySQL, args, err := db.DB.BindNamed(query, arg)
	if err != nil {
		return nil, err
	}
	return db.Queryx(querySQL, args...)
}

// NamedExec using this DB.
// Any named placeholder parameters are replaced with fields from arg.
func (db *Database) NamedExec(query string, arg interface{}) (sql.Result, error) {
	querySQL, args, err := db.DB.BindNamed(query, arg)
	if err != nil {
		return nil, err
	}
	return db.Exec(querySQL, args...)
}

func (db *Database) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	var err error
	ctx = db.before(ctx, query, args...)
	err = db.DB.SelectContext(ctx, dest, query, args...)
	err = db.after(ctx, err, query, args...)
	return err
}

func (db *Database) Select(dest interface{}, query string, args ...interface{}) error {
	return db.SelectContext(context.Background(), dest, query, args...)
}

func (db *Database) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	var err error
	ctx = db.before(ctx, query, args...)
	err = db.DB.GetContext(ctx, dest, query, args...)
	err = db.after(ctx, err, query, args...)
	return err
}

func (db *Database) Get(dest interface{}, query string, args ...interface{}) error {
	return db.GetContext(context.Background(), dest, query, args...)
}

func (db *Database) ExecContext(ctx context.Context, executeSql string, args ...interface{}) (sql.Result, error) {
	var err error
	ctx = db.before(ctx, executeSql, args...)
	result, err := db.DB.ExecContext(ctx, executeSql, args...)
	err = db.after(ctx, err, executeSql, args...)
	return result, err
}

func (db *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.ExecContext(context.Background(), query, args...)
}

func (db *Database) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	ctx = db.before(ctx, query, args...)
	rows, err := db.DB.QueryContext(ctx, query, args...)
	err = db.after(ctx, err, query, args...)
	return rows, nil
}

func (db *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.QueryContext(context.Background(), query, args...)
}

func (db *Database) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	ctx = db.before(ctx, query, args...)
	rows, err := db.DB.QueryxContext(ctx, query, args...)
	err = db.after(ctx, err, query, args...)
	return rows, err
}

func (db *Database) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	return db.QueryxContext(context.Background(), query, args...)
}

func (db *Database) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	ctx = db.before(ctx, query, args...)
	rows := db.DB.QueryRowxContext(ctx, query, args...)
	_ = db.after(ctx, rows.Err(), query, args...)
	return rows
}

func (db *Database) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	return db.QueryRowxContext(context.Background(), query, args...)
}

func (db *Database) Close() error {
	err := db.DB.Close()
	if err != nil {
		db.Logger.Errorf("close connection fail: %+v", err)
		return err
	}
	return nil
}

func (db *Database) SetDatabaseHook(hook DatabaseHook) {
	db.Hook = hook
}

func (db *Database) SetDatabaseLoggerHook() {
	db.Hook = &DatabaseLoggerHook{db.Logger}
}

func FetchDB(dbName string) *Database {
	db, ok := dbs[dbName]
	if ok {
		return db
	}

	mutex.Lock()
	defer mutex.Unlock()

	if db, ok := dbs[dbName]; ok {
		return db
	}

	config := configs.FetchDBConfigWithName(dbName)
	db = NewDatabase(WithConfig(config))
	err := db.connect()
	if err != nil {
		return nil
	}

	if db != nil {
		return db
	}

	logger.Errorf("can't fetch the db with Alias:[%s]", dbName)
	return nil
}

func (db *Database) connect() error {
	config := db.config
	if config == nil {
		return errors.New(fmt.Sprintf("can't find database config with alias[%s]", config.DBName))
	}

	if db.Logger == nil {
		db.Logger = logger.NewLogger(logger.WithPrefix("DB.%s", strings.ToUpper(config.DBName)))
	}

	sdb, err := sqlx.Open(config.Driver, config.DSN)
	if err != nil {
		db.Logger.Errorf("connect fail: %+v", config.DBName, err)
		return err
	}
	if db.config.Options.MaxOpenConns > 0 {
		sdb.SetMaxOpenConns(db.config.Options.MaxOpenConns)
	}
	if db.config.Options.MaxIdleConns > 0 {
		sdb.SetMaxIdleConns(db.config.Options.MaxIdleConns)
	}
	if db.config.Options.ConnMaxIdleTime.Nanoseconds() > 0 {
		sdb.SetConnMaxIdleTime(db.config.Options.ConnMaxIdleTime)
	}
	if db.config.Options.ConnMaxLifetime.Nanoseconds() > 0 {
		sdb.SetConnMaxLifetime(db.config.Options.ConnMaxLifetime)
	}

	err = sdb.Ping()
	if err != nil {
		db.Logger.Errorf("ping fail: %+v", err)
		return err
	}

	db.DB = sdb
	db.Logger.Debugf(color.New(color.Bold, color.OpUnderscore, color.FgGreen).Sprintf("database connect successfully: [%s]", config.DesensitiseDSN))
	if db.Hook != nil {
		db.Hook.SetLogger(db.Logger)
	}
	dbs[config.DBName] = db
	return nil
}
