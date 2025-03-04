package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"tgBotNote/internal/bot"
)

func main() {

	bot.InitDB()

	botAPI, err := tgbotapi.NewBotAPI("7944124426:AAF5QWjCA6jq8BzH6SSpU-52nxO1f0q7KKc")
	if err != nil {
		log.Fatal(err)
	}

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

		chatID := update.Message.Chat.ID

		if update.Message.IsCommand() { // Если это команда — обрабатываем
			bot.HandleCommand(update, botAPI)
			continue
		}

		switch bot.UserStates[chatID] { // Обращаемся к переменной из bot
		case "waiting_for_title":
			bot.UserNotes[chatID] = update.Message.Text
			bot.UserStates[chatID] = "waiting_for_text"
			bot.SendMessage(botAPI, chatID, "Введи текст заметки")
		case "waiting_for_text":
			title := bot.UserNotes[chatID]
			text := update.Message.Text

			bot.AddNote(chatID, title, text)

			log.Printf("Добавлена заметка: Название: %s, Текст: %s", title, text)

			delete(bot.UserStates, chatID)
			delete(bot.UserNotes, chatID)

			bot.SendMessage(botAPI, chatID, "Заметка добавлена")
		}
	}
}
