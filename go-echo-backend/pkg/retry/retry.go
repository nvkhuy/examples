package retry

import "github.com/cenkalti/backoff/v4"

type ExecFunc = func() error

func Exec(maxAttempts uint64, fn ExecFunc) error {
	var b = backoff.WithMaxRetries(backoff.NewExponentialBackOff(), maxAttempts)
	var err = backoff.Retry(fn, b)
	return err
}
