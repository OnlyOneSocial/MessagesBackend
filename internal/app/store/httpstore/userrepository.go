package httpstore

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/katelinlis/MessagesBackend/internal/app/model"
)

//UserRepository ...
type UserRepository struct {
	store *Store
}

//GetUsernameAndAvatar ...
func (r *UserRepository) GetUsernameAndAvatar(AuthorID int) model.UserObj {
	if val, ok := r.store.userCache[AuthorID]; ok {
		return val
	}
	userID := strconv.Itoa(AuthorID)

	client := http.Client{}
	resp, err := client.Get(`http://localhost:3044/api/user/get/` + userID)
	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	fmt.Println(result["user"]["username"])
	fmt.Println(result["user"]["avatar"])

	usrObj := model.UserObj{
		Username: result["user"]["username"],
		Avatar:   result["user"]["avatar"],
	}

	r.store.userCache[AuthorID] = usrObj

	return usrObj
}

//GetFriends ...
func (r *UserRepository) GetFriends(AuthorID int) []int {
	if val, ok := r.store.friendsCache[AuthorID]; ok {
		return val
	}
	userID := strconv.Itoa(AuthorID)

	client := http.Client{}
	resp, err := client.Get(`http://localhost:3044/api/user/array_friends/` + userID)
	if err != nil {
		log.Fatalln(err)
	}

	var result []int
	json.NewDecoder(resp.Body).Decode(&result)

	fmt.Println(result)
	r.store.friendsCache[AuthorID] = result
	return result
}
