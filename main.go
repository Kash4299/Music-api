package main

import (
	"io"
	"os"
	"path/filepath"
	"sample/api"
	"sample/common/cache"
	"sample/docs"
	"sample/internal/mongodb"
	"sample/internal/redis"
	"sample/repository"
	"sample/repository/db"
	"sample/service"

	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Dir      string `env:"CONFIG_DIR" envDefault:"config/config.json"`
	Port     string
	Redis    string
	LogType  string
	LogFile  string
	LogLevel string
	Mongodb  string
}

var config Config

func init() {
	if err := env.Parse(&config); err != nil {
		log.Panic("Get environment values fail")
		log.Fatal(err)
	}
	viper.SetConfigFile(config.Dir)
	if err := viper.ReadInConfig(); err != nil {
		log.Println(err.Error())
		panic(err)
	}
	cfg := Config{
		Port:     viper.GetString(`main.port`),
		Redis:    viper.GetString(`main.redis`),
		LogFile:  viper.GetString(`main.log_file`),
		LogType:  viper.GetString(`main.log_type`),
		LogLevel: viper.GetString(`main.log_level`),
		Mongodb:  viper.GetString(`main.mongodb`),
	}
	if cfg.Redis == "enabled" {
		var err error
		redis.Redis, err = redis.NewRedis(redis.Config{
			Addr:         viper.GetString(`redis.address`),
			Password:     viper.GetString(`redis.password`),
			DB:           viper.GetInt(`redis.database`),
			PoolSize:     30,
			PoolTimeout:  20,
			IdleTimeout:  10,
			ReadTimeout:  20,
			WriteTimeout: 15,
		})
		if err != nil {
			panic(err)
		}
	}

	config = cfg
}

// @title Music API
// @version 1.0
// @description This is a sample music API.
// @license.name hobaduy
// @host localhost:8000
// @BasePath /v1
func main() {
	_ = os.Mkdir(filepath.Dir(config.LogFile), 0755)
	file, _ := os.OpenFile(config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer file.Close()
	setAppLogger(config, file)

	cache.MCache = cache.NewMemCache()
	defer cache.MCache.Close()

	if config.Redis == "enabled" {
		cache.RCache = cache.NewRedisCache(redis.Redis.GetClient())
		defer cache.RCache.Close()
	}

	if config.Mongodb == "enabled" {
		client, err := mongodb.InitMongodb(viper.GetString(`mongodb_uri`))
		if err != nil {
			panic(err)
		}

		repository.TrackRepo = db.NewTrack(client)
		repository.PlaylistRepo = db.NewPlaylist(client)

		defer mongodb.CloseDB()
	}

	server := api.NewServer()
	musicTrackService := service.NewTrack()
	api.APIMusicTrackHandler(server.Engine, musicTrackService)

	playlistService := service.NewPlaylist()
	api.APIPlaylistHandler(server.Engine, playlistService)

	docs.SwaggerInfo.BasePath = "/v1"
	api.APISwaggerHandler(server.Engine)

	server.Start(config.Port)
}

func setAppLogger(cfg Config, file *os.File) {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.999Z",
	})

	switch cfg.LogLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	switch cfg.LogType {
	case "DEFAULT":
		log.SetOutput(os.Stdout)
	case "FILE":
		if file != nil {
			log.SetOutput(io.MultiWriter(os.Stdout, file))
		} else {
			log.SetOutput(os.Stdout)
		}
	default:
		log.SetOutput(os.Stdout)
	}
}
