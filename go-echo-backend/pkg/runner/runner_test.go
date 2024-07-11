package runner

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/logger"
)

func TestRunner(t *testing.T) {
	logger.Init()
	runner := New(5)
	defer runner.Release()

	for i := 0; i < 100; i++ {
		i := i
		runner.Submit(func() {
			fmt.Println("Job: ", i)
			time.Sleep(time.Second * 1)
		})
	}

	runner.Wait()
}

func TestRunnerInsideWorker(t *testing.T) {
	logger.Init()
	var wg = sync.WaitGroup{}
	wg.Add(1)

	go func() {
		runner := New(5)
		defer runner.Release()

		for i := 0; i < 100; i++ {
			i := i
			_ = runner.Submit(func() {
				fmt.Println("Job: ", i)
				time.Sleep(time.Second * 1)
			})
		}

		runner.Wait()
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("Done")

}
func TestRunnerZeroWorker(t *testing.T) {
	logger.Init()
	runner := New(0)
	defer runner.Release()

	for i := 0; i < 100; i++ {
		i := i
		_ = runner.Submit(func() {
			fmt.Println("Job: ", i)
			time.Sleep(time.Second * 1)
		})
	}

	runner.Wait()
}
