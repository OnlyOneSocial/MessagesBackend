package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/katelinlis/MessagesBackend/internal/app/model"
)

func (s *server) ConfigureWallRouter() {

	router := s.router.PathPrefix("/api/message").Subrouter()
	router.HandleFunc("/send", s.HandleSendMessage()).Methods("POST")   // Получение всей стены
	router.HandleFunc("/get/{im}", s.HandleGetMessage()).Methods("GET") // Получение последних сообщений
	router.HandleFunc("/longpoll", s.HandleLongpoll()).Methods("GET")   // получение сообщений в фоне
	router.HandleFunc("/im", s.HandleGetIm()).Methods("GET")            // получение диалогов
}

//MessageSend ...
type MessageSend struct {
while(true){
	Text   string `json:"text"`
	SendTo int    `json:"to"`}
}

func (s *server) HandleSendMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		userid, err := s.GetDataFromToken(w, request)
		if err != nil {
			fmt.Println(err)
		}
		var messageSend MessageSend
		err = json.NewDecoder(request.Body).Decode(&messageSend)
		if err != nil {
			s.error(w, request, http.StatusBadRequest, err)
		}

		var message model.Message
		message.Userid = int(userid)
		message.Text = messageSend.Text

		message.SendTo = messageSend.SendTo
		err = s.store.Message().Send(&message)

		if err != nil {
			s.respond(w, request, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, request, http.StatusOK, message)
	}
}

//MessageDialog ...
type MessageDialog struct {
	ImWith int `json:"userid"`
}

func (s *server) HandleGetMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		offset, limit := s.UrlLimitOffset(request)

		userid, err := s.GetDataFromToken(w, request)
		if err != nil {
			fmt.Println(err)
		}
		vars := mux.Vars(request)
		ImWith, err := strconv.Atoi(vars["im"])
		if err != nil {
			fmt.Println(err)
		}

		messages, err := s.store.Message().GetMessages(int(userid), ImWith, limit, offset)
		if err != nil {
			fmt.Println(err)
		}
		for index, element := range messages {
			UserObj := s.HTTPstore.User().GetUsernameAndAvatar(element.Userid)
			messages[index].Username = UserObj.Username
			messages[index].Avatar = UserObj.Avatar
		}
		s.respond(w, request, http.StatusOK, messages)
	}
}

func (s *server) HandleGetIm() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		userid, err := s.GetDataFromToken(w, request)
		if err != nil {
			fmt.Println(err)
		}

		offset, limit := s.UrlLimitOffset(request)
		Dialogs, err := s.store.Message().GetDialogList(int(userid), limit, offset)
		if err != nil {
			fmt.Println(err)
		}

		for index, element := range Dialogs {
			UserObj := s.HTTPstore.User().GetUsernameAndAvatar(element.WithIm)
			Dialogs[index].Username = UserObj.Username
			Dialogs[index].Avatar = UserObj.Avatar
		}

		s.respond(w, request, http.StatusOK, Dialogs)
	}
}

func (s *server) HandleLongpoll() http.HandlerFunc {
	//MessageDialog ...
	type LongPollIM struct {
		ImWith               int `json:"userid"`
		LastMessageTimestamp int `json:"lastMessageTimestamp"`
	}
	return func(w http.ResponseWriter, request *http.Request) {
		userid, err := s.GetDataFromToken(w, request)
		if err != nil {
			fmt.Println(err)
		}

		var messageDialog LongPollIM
		err = json.NewDecoder(request.Body).Decode(&messageDialog)
		if err != nil {
			s.error(w, request, http.StatusBadRequest, err)
		}

		messages, err := s.store.Message().ScanForNewMessage(int(userid), messageDialog.ImWith, messageDialog.LastMessageTimestamp)
		if err != nil {
			fmt.Println(err)
		}

		s.respond(w, request, http.StatusOK, messages)
	}
}
