package main

import (
	"Step_game/database"
	"Step_game/migrations"

	conf "Step_game/config"
	bot "Step_game/tg_bot"
)

func main() {
	//todo добавить явную инициализацию конфига
	//todo инициализация конфига, запуск бота(восстановление бота после паники, мягкое завершение)

	db := database.InitDBMust("step_game.db")
	defer db.Close()

	// Применение миграций (паникует при ошибках)
	migrations.RunMigrations(db)

	bot := bot.InitBot(conf.TgKey, db, false)
	bot.Run()

	// Используй .env файл для хранения ключей. и не забудь .env.template
	// viper библиотека для env файла
}
