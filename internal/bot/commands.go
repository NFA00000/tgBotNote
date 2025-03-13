package bot

import (
	"context"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var UserTimers = make(map[int64]context.CancelFunc) // Хранилище таймеров
var UserStates = make(map[int64]string)             // Хранит, что ожидает бот от пользователя
var UserNotes = make(map[int64]string)              // Временное хранилище названий заметок

// Команды бота
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

// Функция отправки сообщения
func SendMessage(bot *tgbotapi.BotAPI, chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка при отправке сообщения пользователю %d: %v", chatID, err)
	}
}

// Функция запуска тацмера
func StartTimer(chatID int64, duration time.Duration, onTimeout func()) {
	//Если есть активный таймер - отменяем  его
	if cancel, exists := UserTimers[chatID]; exists {
		cancel()
	}

	//Создаем новый контекст с таймером
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	UserTimers[chatID] = cancel

	go func() {
		<-ctx.Done()               // Проверяем, не закончился ли таймер
		onTimeout()                // Вызываем переданную функцию (сохранение или сброс)
		delete(UserTimers, chatID) // Удаляем таймер из хранилища
	}()
}

// Зачистка временных статусов пользователя
func ResetUserStatus(userID int64) {
	delete(UserStates, userID) // Удаляем статус пользователя
	delete(UserNotes, userID)  // Удаляем временные заметки
	delete(UserTimers, userID) // Удаляем таймер пользователя
	log.Printf("Статус пользователя %d сброшен", userID)
}
