package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/app"
	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/middlewares"
	"github.com/engineeringinflow/inflow-backend/pkg/validation"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/rotisserie/eris"

	"github.com/hibiken/asynqmon"
)

var instance *Worker

var (
	queueLow      string = "low"
	queueDefault  string = "default"
	queueMedium   string = "medium"
	queueHigh     string = "high"
	queueCritical string = "critical"
)

var (
	QueueLow      = asynq.Queue(queueLow)
	QueueDefault  = asynq.Queue(queueDefault)
	QueueMedium   = asynq.Queue(queueMedium)
	QueueHigh     = asynq.Queue(queueHigh)
	QueueCritical = asynq.Queue(queueCritical)
)

// TaskHandler handler
type TaskHandler interface {
	TaskName() string
	Handler(context.Context, *asynq.Task) error
	GetPayload() []byte
	Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error)
}

type Handler func(context.Context, *asynq.Task) error

// Worker dispatcher
type Worker struct {
	*Config

	Server    *asynq.Server
	Validator *validation.Validator
	Client    *asynq.Client
	Inspector *asynq.Inspector

	Scheduler *asynq.Scheduler
	mux       *asynq.ServeMux

	tasks map[string]TaskHandler
}

// Config config
type Config struct {
	Namespace    string
	DefaultQueue string
	Logger       *logger.Logger
	App          *app.App
	Size         int
	IsConsumer   bool
}

// New dispatcher
func New(config *Config) *Worker {

	var client = asynq.NewClient(getRedisConnOpt(config.App.Config))
	instance = &Worker{
		Validator: validation.RegisterValidation(),
		Client:    client,
		Config:    config,
		Inspector: asynq.NewInspector(getRedisConnOpt(config.App.Config)),
	}

	if config.IsConsumer {
		instance.Scheduler = instance.createScheduler()
		instance.Server = instance.createServer()
	}

	return instance
}

// Start start async
func (worker *Worker) Run() {

	go func() {
		var err = worker.Server.Run(worker.mux)
		if err != nil {
			worker.Logger.Error("Run server error", zap.Error(err))
			return
		}
		worker.Logger.Debugf("Run server successfully")

	}()

	go func() {
		var err = worker.Scheduler.Run()
		if err != nil {
			worker.Logger.Error("Run scheduler server error", zap.Error(err))
		}

		worker.Logger.Debugf("Run scheduler successfully")
	}()

}

// SendTaskWithContext send task with context
func (worker *Worker) SendDynamicTask(task *asynq.Task, handler Handler) (*asynq.TaskInfo, error) {
	worker.mux.HandleFunc(task.Type(), handler)

	return worker.Client.EnqueueContext(context.Background(), task)

}

// SendTaskWithContext send task with context
func (worker *Worker) SendTaskWithContext(ctx context.Context, task TaskHandler, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	var options = []asynq.Option{
		asynq.MaxRetry(3),
		asynq.Retention(time.Hour * 24),
	}

	return worker.Client.EnqueueContext(ctx, asynq.NewTask(task.TaskName(), task.GetPayload(), append(options, opts...)...))

}

func (worker *Worker) SendTask(task TaskHandler, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return worker.SendTaskWithContext(context.Background(), task, opts...)
}

func (worker *Worker) SendTaskAt(task TaskHandler, time time.Time, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	var options = []asynq.Option{
		asynq.ProcessAt(time),
	}
	return worker.SendTaskWithContext(context.Background(), task, append(options, opts...)...)
}

func (worker *Worker) SendTaskIn(task TaskHandler, duration time.Duration, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	var options = []asynq.Option{
		asynq.ProcessIn(duration),
	}
	return worker.SendTaskWithContext(context.Background(), task, append(options, opts...)...)
}

func GetInstance() *Worker {
	if instance == nil {
		panic("Worker instance is nil, must be call New() first")
	}

	return instance
}

func (worker *Worker) createServer() *asynq.Server {
	return asynq.NewServer(
		getRedisConnOpt(worker.App.Config),
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: worker.Size,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				queueCritical: 6,
				queueDefault:  3,
				queueLow:      1,
			},
			Logger:   worker.Logger.Sugar(),
			LogLevel: asynq.DebugLevel,
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				worker.Logger.Debugf("Task error task=%s payload=%s err=%+v", task.Type(), string(task.Payload()), err)
			}),
			// See the godoc for other configuration options
		},
	)

}

func (worker *Worker) createScheduler() *asynq.Scheduler {
	return asynq.NewScheduler(
		getRedisConnOpt(worker.App.Config),
		&asynq.SchedulerOpts{
			Logger:   worker.Logger.Sugar(),
			LogLevel: asynq.DebugLevel,
			PreEnqueueFunc: func(task *asynq.Task, opts []asynq.Option) {
				worker.Logger.Debugf("Scheduler::Enqueueing task=%s payload=%s", task.Type(), string(task.Payload()))
			},
			PostEnqueueFunc: func(info *asynq.TaskInfo, err error) {
				worker.Logger.Debugf("Scheduler::Enqueued task=%s result=%s payload=%s", info.Type, info.State, string(info.Payload))
			},
		},
	)

}

func (worker *Worker) ServeMonitoringRoutes(e *echo.Echo) {
	var mon = asynqmon.New(asynqmon.Options{
		RootPath:     "/tasks",
		RedisConnOpt: getRedisConnOpt(worker.App.Config),
	})

	e.Any("/tasks/*", echo.WrapHandler(mon), middlewares.IsBasicAuth())
}

func (worker *Worker) BindAndValidate(rawValue []byte, task TaskHandler) error {
	var err = worker.Bind(rawValue, task)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = worker.Validator.Validate(task)
	if err != nil {
		worker.Logger.DebugAny(task, fmt.Sprintf("Validate task error=%s+v", err))
		return eris.Wrap(err, err.Error())
	}

	return nil
}

func (worker *Worker) Bind(rawValue []byte, task TaskHandler) error {
	var err = json.Unmarshal(rawValue, task)
	if err != nil {
		worker.Logger.ErrorAny(err, "Bind error")
		return eris.Wrap(err, err.Error())
	}
	return nil
}

func (worker *Worker) CreateTaskHandler(tasks ...TaskHandler) *asynq.ServeMux {
	var mux = asynq.NewServeMux()
	worker.tasks = map[string]TaskHandler{}

	for _, task := range tasks {
		mux.HandleFunc(task.TaskName(), task.Handler)

		worker.tasks[task.TaskName()] = task
	}

	worker.mux = mux

	return worker.mux
}

func (worker *Worker) ScheduleTask(conspec string, task TaskHandler, opts ...asynq.Option) (string, error) {
	return worker.Scheduler.Register(conspec, asynq.NewTask(task.TaskName(), task.GetPayload(), opts...))
}

func getRedisConnOpt(config *config.Configuration) asynq.RedisConnOpt {
	if len(config.RedisAddress) > 1 {
		return asynq.RedisClusterClientOpt{
			Addrs: config.RedisAddress,
		}
	}

	return asynq.RedisClientOpt{
		Addr:     config.RedisAddress[0],
		PoolSize: 100,
	}
}
