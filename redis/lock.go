package redis

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/transerver/commons/logger"
	"io"
	"strconv"
	"sync"
	"time"
)

var (
	luaRefresh = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("pexpire", KEYS[1], ARGV[2]) else return 0 end`)
	luaRelease = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("del", KEYS[1]) else return 0 end`)
	luaPTTL    = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("pttl", KEYS[1]) else return -3 end`)

	// ErrNotObtained is returned when a lock cannot be obtained.
	ErrNotObtained = errors.New("redislock: not obtained")

	// ErrLockNotHeld is returned when trying to release an inactive lock.
	ErrLockNotHeld = errors.New("redislock: lock not held")

	tmp   []byte
	tmpMu sync.Mutex
)

type Lock struct {
	Key    string
	value  string
	Locked bool
}

func randomToken() (string, error) {
	tmpMu.Lock()
	defer tmpMu.Unlock()

	if tmp == nil {
		tmp = make([]byte, 16)
	}

	if _, err := io.ReadFull(rand.Reader, tmp); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(tmp), nil
}

func (c *redisClient) obtain(key string, ttl time.Duration) (*Lock, error) {
	token, err := randomToken()
	if err != nil {
		return nil, err
	}

	locked, err := c.SetNX(key, token, ttl).Result()
	if err != nil {
		return nil, err
	}
	return &Lock{Key: key, value: token, Locked: locked}, nil
}

func Obtain(ttl time.Duration, format string, v ...interface{}) (*Lock, error) {
	if len(v) > 0 {
		format = fmt.Sprintf(format, v...)
	}
	lock, err := Client().obtain(format, ttl)
	if err != nil {
		logger.Errorf("redislock:obtain %+v", err)
	}
	return lock, err
}

// TTL returns the remaining time-to-live. Returns 0 if the lock has expired.
func (l *Lock) TTL() (time.Duration, error) {
	res, err := luaPTTL.Run(client, []string{l.Key}, l.value).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	if num := res.(int64); num > 0 {
		return time.Duration(num) * time.Millisecond, nil
	}
	return 0, nil
}

// Refresh extends the lock with a new TTL.
// May return ErrNotObtained if refresh is unsuccessful.
func (l *Lock) Refresh(ttl time.Duration) error {
	ttlVal := strconv.FormatInt(int64(ttl/time.Millisecond), 10)
	status, err := luaRefresh.Run(client, []string{l.Key}, l.value, ttlVal).Result()
	if err != nil {
		return err
	} else if status == int64(1) {
		return nil
	}
	return ErrNotObtained
}

// Release manually releases the lock.
// May return ErrLockNotHeld.
func (l *Lock) Release() error {
	res, err := luaRelease.Run(client, []string{l.Key}, l.value).Result()
	if err == redis.Nil {
		return ErrLockNotHeld
	} else if err != nil {
		return err
	}

	if i, ok := res.(int64); !ok || i != 1 {
		return ErrLockNotHeld
	}
	return nil
}

func (l *Lock) LoggedRelease() {
	err := l.Release()
	if err != nil {
		logger.Errorf("%s, key: %s", err.Error(), l.Key)
	}
}
