package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/caching"
	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/locker"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/validation"
	"github.com/jackc/pgx/v5/pgconn"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

type CallbackInterface interface {
	Register(db *gorm.DB)
	RegisterGenerateXID(db *gorm.DB)
	RegisterGenerateBlurhash(db *gorm.DB)
}

var instance *DB

// DB instance
type DB struct {
	*gorm.DB
	Configuration *config.Configuration
	CustomLogger  *logger.Logger
	Validator     *validation.Validator
	Cache         *caching.Client
	Locker        *locker.Locker
}

// New new db
func New(config *config.Configuration, callback CallbackInterface, caching *caching.Client) *DB {
	var dbLogger = DefaultLogger("db")
	if !config.IsProd() {
		dbLogger.LogLevel = gormLogger.Info
	}

	var dbConfig = &gorm.Config{
		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		Logger:                                   dbLogger,
	}

	db, err := gorm.Open(postgres.Open(config.GetDatabaseURI()), dbConfig)
	if err != nil {
		dbLogger.ZapLogger.Errorf("Connecting %s err=%+v", config.GetDatabaseURI(), err)
		panic(err)
	}

	if config.DBReplicaHost != "" {
		err = db.Use(dbresolver.Register(dbresolver.Config{
			TraceResolverMode: true,
			Replicas:          []gorm.Dialector{postgres.Open(config.GetDatabaseReplicaURI())},
		}))
		if err != nil {
			panic(err)
		}
		dbLogger.ZapLogger.Debugf("Setup replica %s success\n", config.GetDatabaseReplicaURI())
	}

	sqlDB, err := db.DB()
	if err != nil {
		dbLogger.ZapLogger.Errorf("Connecting %s err=%+v", config.GetDatabaseURI(), err)
		panic(err)
	}

	for i := 0; i < 10; i++ {
		err = sqlDB.Ping()
		if err == nil {
			dbLogger.ZapLogger.Debugf("Database %s connected success", config.GetDatabaseURI())
			break
		}
		time.Sleep(time.Second * 2)
		dbLogger.ZapLogger.Debugf("Retry connect to %s %d/%d", config.GetDatabaseURI(), i, 10)
	}

	if callback != nil {
		callback.Register(db)
	}

	sqlDB.SetMaxIdleConns(50)
	sqlDB.SetMaxOpenConns(200)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	instance = &DB{
		DB:            db,
		Locker:        locker.GetInstance(),
		CustomLogger:  dbLogger.ZapLogger,
		Configuration: config,
		Validator:     validation.RegisterValidation(),
		Cache:         caching,
	}

	instance.setupExtensions()
	instance.setupFunctions()
	return instance
}

func NewAnalytic(config *config.Configuration, callback CallbackInterface, caching *caching.Client) *DB {
	var dbLogger = DefaultLogger("analytic_db")
	if !config.IsProd() {
		dbLogger.LogLevel = gormLogger.Info
	}

	var dbConfig = &gorm.Config{
		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		Logger:                                   dbLogger,
	}

	db, err := gorm.Open(postgres.Open(config.GetAnalyticDatabaseURI()), dbConfig)
	if err != nil {
		dbLogger.ZapLogger.Errorf("Connecting %s err=%+v", config.GetAnalyticDatabaseURI(), err)
		panic(err)
	}

	if config.ADBReplicaHost != "" {
		err = db.Use(dbresolver.Register(dbresolver.Config{
			TraceResolverMode: true,
			Replicas:          []gorm.Dialector{postgres.Open(config.GetAnalyticDatabaseReplicaURI())},
		}))
		if err != nil {
			panic(err)
		}
		dbLogger.ZapLogger.Debugf("Setup replica %s success\n", config.GetAnalyticDatabaseReplicaURI())
	}

	sqlDB, err := db.DB()
	if err != nil {
		dbLogger.ZapLogger.Errorf("Connecting %s err=%+v", config.GetAnalyticDatabaseURI(), err)
		panic(err)
	}

	for i := 0; i < 10; i++ {
		err = sqlDB.Ping()
		if err == nil {
			dbLogger.ZapLogger.Debugf("Database %s connected success", config.GetAnalyticDatabaseURI())
			break
		}
		time.Sleep(time.Second * 2)
		dbLogger.ZapLogger.Debugf("Retry connect to %s %d/%d", config.GetAnalyticDatabaseURI(), i, 10)
	}

	if callback != nil {
		callback.RegisterGenerateXID(db)
	}

	sqlDB.SetMaxIdleConns(50)
	sqlDB.SetMaxOpenConns(200)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	instance = &DB{
		DB:            db,
		Locker:        locker.GetInstance(),
		CustomLogger:  dbLogger.ZapLogger,
		Configuration: config,
		Validator:     validation.RegisterValidation(),
		Cache:         caching,
	}

	instance.setupExtensions()
	instance.setupFunctions()
	return instance
}

// GetInstance get instance
func GetInstance() *DB {
	if instance == nil {
		panic("Must be call New() first")
	}
	return instance
}

// WithGorm enable debug log
func (db *DB) WithGorm(gdb *gorm.DB) *DB {
	return &DB{
		DB:            gdb,
		Cache:         db.Cache,
		Configuration: db.Configuration,
		CustomLogger:  db.CustomLogger,
		Validator:     db.Validator,
		Locker:        db.Locker,
	}
}

// WithLogger set logger
func (db *DB) WithLogger(logger *logger.Logger) *DB {
	db.CustomLogger = logger
	return db
}

// WithContext with conntext
func (db *DB) WithContext(ctx context.Context) *DB {
	db.DB = db.DB.WithContext(ctx)
	return db

}

// IsRecordNotFoundError check record not found
func (db *DB) IsRecordNotFoundError(err error) bool {
	return errors.Is(err, sql.ErrNoRows) || errors.Is(err, gorm.ErrRecordNotFound)

}

// DisableDebug check record not found
func (db *DB) SetDebug(debug bool) *gorm.DB {
	if debug {
		return db.Session(&gorm.Session{Logger: gormLogger.Default.LogMode(gormLogger.Silent)})
	}
	return db.DB
}

// Validate validate
func (db *DB) Validate(dest interface{}) error {
	return db.Validator.Validate(dest)
}

func (db *DB) IsDuplicateConstraint(err error) (bool, *pgconn.PgError) {
	if e, ok := err.(*pgconn.PgError); ok {
		if e.Code == "23505" {
			return true, e
		}
	}

	return false, nil
}

func (db *DB) CheckUserDuplicateConstraint(err error) error {
	if db.IsDuplicateUserEmailConstraint(err) {
		return errs.ErrEmailTaken
	}

	if db.IsDuplicateUserPhoneConstraint(err) {
		return errs.ErrPhoneTaken
	}

	if db.IsDuplicateUserNameConstraint(err) {
		return errs.ErrUserNameTaken
	}

	return nil
}

func (db *DB) IsDuplicateUserEmailConstraint(err error) bool {
	if duplicated, e := db.IsDuplicateConstraint(err); duplicated {
		if e != nil && (e.ColumnName == "email" || e.ConstraintName == "users_email_key" || e.ConstraintName == "idx_users_email") {
			return true
		}
	}

	return false
}

func (db *DB) IsDuplicateUserPhoneConstraint(err error) bool {
	if duplicated, e := db.IsDuplicateConstraint(err); duplicated {
		if e.ColumnName == "phone_number" || e.ConstraintName == "users_phone_number_key" || e.ConstraintName == "idx_users_phone_number" {
			return true
		}
	}

	return false
}

func (db *DB) IsDuplicateUserNameConstraint(err error) bool {
	if duplicated, e := db.IsDuplicateConstraint(err); duplicated {
		if e.ColumnName == "user_name" || e.ConstraintName == "users_user_name_key" {
			return true
		}
	}

	return false
}
