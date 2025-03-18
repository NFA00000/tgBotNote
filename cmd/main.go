package main

import (
	"log"
	"os"
	"tgBotNote/internal/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

func main() {

	// Инициализируем базу данных
	bot.InitDB()

	// Загружаем переменные окружения
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	// Получаем токен
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("Токен бота не найден в .env файле")
	}

	// Инициализируем бота
	botAPI, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	// test
	botAPI.Debug = true
	log.Println("Едем")

	// Настраиваем получение обновлений (сообщений)
	u := tgbotapi.NewUpdate(0)               // Начинаем с самого первого обновления
	u.Timeout = 60                           // Таймаут ожидания ответа (рекомендуется 60 сек)
	updates, err := botAPI.GetUpdatesChan(u) // Канал для получения обновлений)
	if err != nil {
		log.Fatal(err)
	}

	// Обрабатываем входящие сообщения
	for update := range updates {
		if update.Message == nil { // Если в апдейте нет сообщения — пропускаем
			continue
		}

		bot.HandleUpdate(update, botAPI) // Обра��атываем входящее сообщение
	}
}
