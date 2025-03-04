package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var UserStates = make(map[int64]string) // Хранит, что ожидает бот от пользователя
var UserNotes = make(map[int64]string)  // Временное хранилище названий заметок

func HandleCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	chatID := update.Message.Chat.ID

	switch update.Message.Command() {
	case "start":
		SendMessage(bot, chatID, "Привет, здесь ты можешь сохранить все самое важное. Испоьзуй команду /help для того чтобы ознакомиться всеми моими командами")
	case "help":
		SendMessage(bot, chatID, "Доступные команды:\n"+
			"/addnote Название Текст - Добавить заметку\n"+
			"/notes - Посмотреть список\n"+
			"/note Название - Открыть заметку\n"+
			"/deletenote Название - Удалить заметку")
	case "addnote":
		UserStates[chatID] = "waiting_for_title"
		SendMessage(bot, chatID, "Введи название заметки")
	default:
		SendMessage(bot, chatID, "Такой команды нет, посмотри /help")
	}
}

func SendMessage(bot *tgbotapi.BotAPI, chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	bot.Send(msg)
}
