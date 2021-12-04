package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/katelinlis/goment"
)

//Message ...
type Message struct {
	Userid    int    `json:"userid"`
	SendTo    int    `json:"sendto"`
	Text      string `json:"text"`
	Timestamp int64  `json:"timestamp"`

	Time     string    `json:"time"`
	RandomID uuid.UUID `json:"random_id"`
	AnswerTO uuid.UUID `json:"answerto"`
	Username string    `json:"username"`
	Avatar   string    `json:"avatar"`
}

//Validate ...
func (w *Message) Validate() error {
	return validation.ValidateStruct(
		w,
		validation.Field(&w.Text, validation.Required, validation.Length(1, 400)),
		validation.Field(&w.Userid, validation.Required),
	)
}

//Proccessing ...
func (w *Message) Proccessing() error {

	goment.SetLocale("ru")
	time, err := goment.Unix(w.Timestamp)
	if err != nil {
		return err
	}
	w.Time = time.FromNow()

	return nil
}

//GenerateUUID ...
func (w *Message) GenerateUUID() error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	w.RandomID = uuid
	return nil
}
