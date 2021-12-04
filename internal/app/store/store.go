package store

/*
Store репозитории данных
*/
type Store interface {
	Message() MessageRepository // интерфейс для сообщений
}

//HTTPStore ...
type HTTPStore interface {
	User() UserRepository //интерфейс для пользователей
}
