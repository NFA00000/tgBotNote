package bot

import (
	"context"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Глобальные переменные
var (
	UserStates = make(map[int64]string)             // Хранит, что ожидает бот от пользователя
	UserNotes  = make(map[int64]string)             // Временное хранилище названий заметок
)

// Обработчик обновлений
func HandleUpdate(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	chatID := update.Message.Chat.ID

	if update.Message.IsCommand() {
		HandleCommand(update, bot) // Обрабатываем команды
		return
	}

	// Обрабатываем состояние пользователя
	HandleUserState(chatID, update.Message.Text, bot)
}

// Обработчик команд
func HandleCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	chatID := update.Message.Chat.ID

	switch update.Message.Command() {
	case "start":
		SendMessage(bot, chatID, "Привет! Используй команду /help, чтобы узнать мои команды.")
	case "help":
		SendMessage(bot, chatID, "Доступные команды:\n"+
			"/addnote - Добавить заметку\n"+
			"/notes - Посмотреть список\n"+
			"/note Название - Открыть заметку\n"+
			"/deletenote Название - Удалить заметку")
	case "addnote":
		UserStates[chatID] = "waiting_for_title"
		SendMessage(bot, chatID, "Введи название заметки")

	case "notes":
		notes, err := GetNotesByUserID(chatID)
		if err != nil {
			SendMessage(bot, chatID, "Щшибка при получении списка заметок")
			return
		}

		if len(notes) == 0 {
			SendMessage(bot, chatID, "У тебя пока нет заметок")
			return
		}

		message := "Твои заметки:\n"
		for _, note := range notes {
			message += "• " + note + "\n"
		}
		SendMessage(bot, chatID, message)
	default:
		SendMessage(bot, chatID, "Неизвестная команда. Используй /help.")
	}
}

// Обработчик состояний пользователя
func HandleUserState(chatID int64, text string, bot *tgbotapi.BotAPI) {
	switch UserStates[chatID] {
	case "waiting_for_title":
		UserNotes[chatID] = text
		UserStates[chatID] = "waiting_for_text"
		SendMessage(bot, chatID, "Введи текст заметки")

	case "waiting_for_text":
		title := UserNotes[chatID]
		noteText := text

		// Сохраняем в БД
		AddNote(chatID, title, noteText)

		log.Printf("Добавлена заметка: %s - %s", title, noteText)

		ResetUserStatus(chatID) // Сбрасываем состояния

		SendMessage(bot, chatID, "Заметка добавлена")
	}
}

// Отправка сообщения
func SendMessage(bot *tgbotapi.BotAPI, chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	}
}

// Сброс состояния
func ResetUserStatus(chatID int64) {
	delete(UserStates, chatID)
	delete(UserNotes, chatID)
	log.Printf("Статус пользователя %d сброшен", chatID)
}
