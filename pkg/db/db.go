package db

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/natefinch/lumberjack"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBConnections stores multiple database connections
var dbConnections = make(map[string]*gorm.DB)

// InitDB initializes a database connection and stores it in DBConnections
func InitDB(dbName, dbType, dsn string, customLogger logger.Interface) {
	var err error

	var db *gorm.DB
	switch dbType {
	case "postgres":
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: customLogger,
		})
	case "mysql":
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: customLogger,
		})
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
			Logger: customLogger,
		})
	default:
		log.Fatalf("Unsupported database type: %s", dbType)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Set connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database from GORM: %v", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Store the connection in the map
	dbConnections[dbName] = db

	log.Printf("Database connection '%s' established", dbName)
}

// GetDB retrieves a database connection by name
func getDB(dbName string) *gorm.DB {
	db, exists := dbConnections[dbName]
	if !exists {
		log.Fatalf("Database connection '%s' not found", dbName)
	}
	return db
}

// Create inserts a new record into the specified database
func Create(dbName string, value interface{}) error {
	db := getDB(dbName)
	return db.Create(value).Error
}

// Update modifies an existing record in the specified database
func Update(dbName string, value interface{}) error {
	db := getDB(dbName)
	return db.Save(value).Error
}

// Delete removes a record from the specified database
func Delete(dbName string, value interface{}) error {
	db := getDB(dbName)
	return db.Delete(value).Error
}

// Find retrieves records from the specified database based on conditions
func Find(dbName string, out interface{}, where ...interface{}) error {
	db := getDB(dbName)
	return db.Find(out, where...).Error
}

// ExecSQL executes a raw SQL query on the specified database
func ExecSQL(dbName, sql string, values ...interface{}) error {
	db := getDB(dbName)
	return db.Exec(sql, values...).Error
}

// QuerySQL executes a raw SQL query on the specified database and scans the result into the provided destination
func QuerySQL(dbName string, dest interface{}, sql string, values ...interface{}) error {
	db := getDB(dbName)
	return db.Raw(sql, values...).Scan(dest).Error
}

// ExecuteInTransaction executes the given function within a database transaction
func ExecuteInTransaction(dbName string, txFunc func(tx *gorm.DB) error) error {
	db := getDB(dbName)
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // re-throw the panic after Rollback
		}
	}()

	if err := txFunc(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// PaginateQuery performs a paginated query on the specified table
func PaginateQuery(dbName, tableName string, pageSize int, processFunc func(records []map[string]interface{}) error, conds ...interface{}) error {
	db := getDB(dbName)
	var (
		page       = 1
		totalCount int64
	)

	query := db.Table(tableName)

	// 如果有条件，添加到查询中
	if len(conds) > 0 {
		query = query.Where(conds[0], conds[1:]...)
	}

	// 获取总记录数
	if err := query.Count(&totalCount).Error; err != nil {
		return err
	}

	log.Printf("Total records in %s: %d", tableName, totalCount)

	// 循环分页查询
	for {
		var records []map[string]interface{}

		// 查询当前页的数据
		if err := query.Limit(pageSize).
			Offset((page - 1) * pageSize).
			Find(&records).Error; err != nil {
			return err
		}

		// 如果没有更多数据，退出循环
		if len(records) == 0 {
			break
		}

		// 处理查询到的数据
		if err := processFunc(records); err != nil {
			return err
		}

		log.Printf("Processed page %d with %d records", page, len(records))

		// 如果已经查询完所有数据，退出循环
		if page*pageSize >= int(totalCount) {
			break
		}

		page++
	}

	return nil
}

// Paginate retrieves records from the specified database based on conditions and supports pagination
func Paginate(dbName string, out interface{}, page, pageSize int, where ...interface{}) error {
	db := getDB(dbName)
	offset := (page - 1) * pageSize
	return db.Limit(pageSize).Offset(offset).Find(out, where...).Error
}

// DBLoggerConfig defines the configuration for the custom logger
type DBLoggerConfig struct {
	LogFilePath   string          // 日志文件路径（为空时输出到控制台）
	MaxSize       int             // 单个日志文件的最大大小（单位：MB）, 当日志文件路径不为空时生效
	MaxBackups    int             // 保留的旧日志文件的最大数量, , 当日志文件路径不为空时生效
	MaxAge        int             // 日志文件的最大保存天数 , , 当日志文件路径不为空时生效
	Compress      bool            // 是否压缩旧日志文件, 当日志文件路径不为空时生效
	LogLevel      logger.LogLevel // 日志等级  	Silent Error Warn Info
	SlowThreshold time.Duration   // 慢查询阈值
}

// NewDBCustomLogger creates a custom logger for GORM with log rotation or stdout
func NewDBCustomLogger(config DBLoggerConfig) logger.Interface {
	var (
		logOutput   *log.Logger
		logFilePath string
	)

	if config.LogFilePath == "" {
		// 如果日志文件路径为空，则输出到控制台
		logOutput = log.New(os.Stdout, "\r\n", log.LstdFlags)
	} else {

		if filepath.IsAbs(config.LogFilePath) {
			logFilePath = config.LogFilePath // 绝对路径
		} else {
			baseDIR, err := os.Getwd()
			if err != nil {
				log.Fatalf("Failed to get current working directory: %v", err)
			}

			// 相对路径基于 logs 目录, 相对路径
			logFilePath = filepath.Join(baseDIR, config.LogFilePath)
		}

		logDir := filepath.Dir(logFilePath)
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			log.Fatalf("Failed to create log directory: %v", err)
		}

		// 配置日志轮转
		rotatingLogger := &lumberjack.Logger{
			Filename:   logFilePath,       // 日志文件路径
			MaxSize:    config.MaxSize,    // 单个日志文件的最大大小（MB）
			MaxBackups: config.MaxBackups, // 保留的旧日志文件的最大数量
			MaxAge:     config.MaxAge,     // 日志文件的最大保存天数
			Compress:   config.Compress,   // 是否压缩旧日志文件
		}
		logOutput = log.New(rotatingLogger, "\r\n", log.LstdFlags)
	}

	// 创建自定义日志器
	return logger.New(
		logOutput,
		logger.Config{
			SlowThreshold:             config.SlowThreshold * time.Millisecond, // 慢查询阈值
			LogLevel:                  config.LogLevel,                         // 日志等级
			IgnoreRecordNotFoundError: true,                                    // 忽略记录未找到的错误
			Colorful:                  config.LogFilePath == "",                // 控制台日志启用彩色输出
		},
	)
}
