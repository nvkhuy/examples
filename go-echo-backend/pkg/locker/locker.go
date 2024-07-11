package locker

import (
	"fmt"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/go-redsync/redsync/v4"
	redsyncgoredis "github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"github.com/rotisserie/eris"
	"github.com/rs/xid"
)

var instance *Locker

type LockReleaseFunc = func() error

type Locker struct {
	workspace     string
	redisSync     *redsync.Redsync
	logger        *logger.Logger
	defaultExpiry time.Duration
}

func New(config *config.Configuration) *Locker {
	var redisClient = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress[0],
		Password: config.RedisPassword,
		DB:       config.RedisDBMode,
	})
	var pool = redsyncgoredis.NewPool(redisClient)
	var redisSync = redsync.New(pool)

	instance = &Locker{
		workspace:     "inflow_locker",
		redisSync:     redisSync,
		logger:        logger.New("utils/locker"),
		defaultExpiry: time.Minute,
	}
	return instance
}

func (locker *Locker) AcquireLock(key string, timeoutDuration time.Duration) (LockReleaseFunc, error) {
	if key == "" {
		key = xid.New().String()
	}

	var timeout = time.NewTimer(helper.GetTimeout(timeoutDuration, time.Minute))
	var start = time.Now()
	for {
		select {
		case <-timeout.C:
			var fn = func() error {
				return nil
			}
			return fn, eris.Errorf("timeout reached after %v", time.Since(start))
		default:
			var lockKey = fmt.Sprintf("%s_%s", locker.workspace, key)
			var mutex = locker.redisSync.NewMutex(lockKey, redsync.WithExpiry(locker.defaultExpiry))
			var err = mutex.Lock()
			if err != nil {
				locker.logger.Errorf("Acquire lock for key %s error: %+v", lockKey, err)
				time.Sleep(time.Millisecond * 200)
			} else {
				locker.logger.WithSkipCaller(1).Infof("Acquire lock for key = %s expiry = %0.2fs success", lockKey, locker.defaultExpiry.Seconds())
				var fn = func() error {
					if ok, err := mutex.Unlock(); !ok || err != nil {
						var e = eris.Errorf("Release lock for key %s error: %+v", mutex.Name(), err)
						locker.logger.ErrorAny(e)
						return eris.Wrap(e, "")
					}
					locker.logger.WithSkipCaller(1).Infof("Release lock for key = %s success after %0.2fs", lockKey, time.Since(start).Seconds())
					return nil
				}
				return fn, nil
			}

		}
	}

}
func GetInstance() *Locker {
	return instance
}
