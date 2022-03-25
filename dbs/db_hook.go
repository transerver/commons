package dbs

import (
	"context"
	"fmt"
	"github.com/transerver/commons/logger"
	"strconv"
	"strings"
	"time"
)

type DatabaseHook interface {
	Before(ctx context.Context, query string, args ...interface{}) (context.Context, error)
	After(ctx context.Context, query string, args ...interface{}) (context.Context, error)
	OnError(ctx context.Context, err error, query string, args ...interface{})
	SetLogger(logger *logger.Logger)
}

type DatabaseLoggerHook struct {
	Logger *logger.Logger
}

func (h *DatabaseLoggerHook) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	return context.WithValue(ctx, "startTime", time.Now()), nil
}

func (h *DatabaseLoggerHook) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	if h.Logger.Level > logger.DebugLevel {
		return ctx, nil
	}

	if startTime, ok := ctx.Value("startTime").(time.Time); ok {
		took := time.Since(startTime)
		h.Logger.Debugf("SQL: %s, Took: %s", expansionArgs(query, args...), took)
	} else {
		h.Logger.Debugf("SQL: %s", expansionArgs(query, args...))
	}
	return ctx, nil
}

func (h *DatabaseLoggerHook) OnError(ctx context.Context, err error, query string, args ...interface{}) {
	if startTime, ok := ctx.Value("startTime").(time.Time); ok {
		took := time.Since(startTime)
		h.Logger.Errorf("SQL: %s, Took: %s, Error: %+v", expansionArgs(query, args...), took, err)
	} else {
		h.Logger.Errorf("SQL: %s, Error: %+v", expansionArgs(query, args...), err)
	}
}

func (h *DatabaseLoggerHook) SetLogger(logger *logger.Logger) {
	h.Logger = logger
}

func expansionArgs(query string, args ...interface{}) string {
	for i, arg := range args {
		query = strings.Replace(query, "$"+strconv.Itoa(i+1), fmt.Sprintf("`%+v`", arg), 1)
	}
	return query
}

func (db *Database) before(ctx context.Context, query string, args ...interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if db.Hook == nil {
		return ctx
	}

	bctx, err := db.Hook.Before(ctx, query, args...)
	if err != nil {
		db.Hook.OnError(bctx, err, query, args...)
	}
	return bctx
}

func (db *Database) after(ctx context.Context, err error, query string, args ...interface{}) error {
	if db.Hook == nil {
		return err
	}

	if err != nil {
		db.Hook.OnError(ctx, err, query, args...)
		return err
	}
	_, err = db.Hook.After(ctx, query, args...)
	if err != nil {
		db.Hook.OnError(ctx, err, query, args...)
	}
	return err
}
