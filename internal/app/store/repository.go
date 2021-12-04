package store

import "github.com/katelinlis/MessagesBackend/internal/app/model"

//MessageRepository ...
type MessageRepository interface {
	Send(*model.Message) error                                                                   // Отправка сообщения
	GetMessages(userid int, ImWith int, limit int, offset int) ([]model.Message, error)          // Получение сообщений в диалоги
	GetDialogList(userid int, limit int, offset int) ([]model.Dialogs, error)                    // Получить список диалогов
	ScanForNewMessage(userid int, ImWith int, lastMessageTimestamp int) ([]model.Message, error) // Сканирование на новое сообщение
}

//UserRepository ...
type UserRepository interface {
	GetUsernameAndAvatar(int) model.UserObj //Получение имени пользователя
	GetFriends(int) []int                   //Получение списка друзей пользователя
}
