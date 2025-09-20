package main

import (
	conf "Step_game/config"
	bot "Step_game/tg_bot"
)

func main() {
	//todo добавить явную инициализацию конфига
	//todo инициализация конфига, инициализация бота, запуск бота(восстановление бота после паники, мягкое завершение)

	bot.InitBot(conf.TgKey, conf.DbPath, false)

	// Используй .env файл для хранения ключей. и не забудь .env.template
	// viper библиотека для env файла
}
