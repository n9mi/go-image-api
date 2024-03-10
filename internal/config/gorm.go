package config

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(viperConfig *viper.Viper, log *logrus.Logger) *gorm.DB {
	host := viperConfig.GetString("DB_HOST")
	user := viperConfig.GetString("DB_USER")
	password := viperConfig.GetString("DB_PASSWORD")
	name := viperConfig.GetString("DB_NAME")
	port := viperConfig.GetInt("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		host,
		user,
		password,
		name,
		port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(&logrusWriter{Logger: log}, logger.Config{
			SlowThreshold:             5 * time.Second,
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			LogLevel:                  logger.Info,
		}),
	})
	if err != nil {
		log.Fatalf("failed to connect into database: %+v", err)
	}

	connection, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get database connection : %+v", err)
	}

	iddleConnection := viperConfig.GetInt("DB_IDDLE_CONNECTION")
	maxConnection := viperConfig.GetInt("DB_MAX_CONNECTION")
	maxLifetime := viperConfig.GetInt("DB_MAX_CONNECTION")
	connection.SetMaxIdleConns(iddleConnection)
	connection.SetMaxOpenConns(maxConnection)
	connection.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Second)

	return db
}

type logrusWriter struct {
	Logger *logrus.Logger
}

func (w *logrusWriter) Printf(message string, args ...interface{}) {
	w.Logger.Tracef(message, args...)
}
