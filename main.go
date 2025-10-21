package main

import (
	"Step_game/database"
	"Step_game/migrations"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"

	"Step_game/tg_bot"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

func main() {
	// переменные из .env файла
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Инициализация логгера
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Инициализация БД
	db := database.InitDBMust(os.Getenv("DB_PATH"))
	defer db.Close()

	// Применение миграций (паникует при ошибках)
	migrations.RunMigrations(db)
	
	// Конфигурация бота
	cfg := tg_bot.Config{
		Token: os.Getenv("TELEGRAM_BOT_TOKEN"),
		Debug: os.Getenv("DEBUG") == "true",
	}

	// Зависимости
	deps := tg_bot.Dependencies{
		DB:     db,
		Logger: logger,
	}

	// Создание бота
	bot, err := tg_bot.NewBot(cfg, deps)
	if err != nil {
		logger.Fatal("Failed to create bot", zap.Error(err))
	}
	defer func(bot *tg_bot.Bot) {
		err := bot.Close()
		if err != nil {
			logger.Fatal("Failed to close bot", zap.Error(err))
		}
	}(bot)

	// Запуск бота
	logger.Info("Starting bot...")

	// Обработка graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	bot.Run()

	<-quit
	logger.Info("Shutting down...")
}
