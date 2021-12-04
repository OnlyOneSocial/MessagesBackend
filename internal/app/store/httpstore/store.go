package httpstore

import (
	"github.com/katelinlis/MessagesBackend/internal/app/model"
	"github.com/katelinlis/MessagesBackend/internal/app/store"
)

//Store ...

//Store ...
type Store struct {
	userCache      map[int]model.UserObj
	friendsCache   map[int][]int
	userRepository *UserRepository
}

/*
Todo поставить ограничение на карту

 и удалять некоторую информацию в случае достяжении лимита

 for k := range userCache {
    delete(userCache, k)
}

а лучше подключить это https://github.com/patrickmn/go-cache
*/
//New ...
func New() *Store {
	return &Store{
		userCache:    make(map[int]model.UserObj),
		friendsCache: make(map[int][]int),
	}
}

//User ...
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}
