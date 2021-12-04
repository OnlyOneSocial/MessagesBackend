package model

import (
	"github.com/google/uuid"
	"github.com/katelinlis/goment"
)

//Dialogs ...
type Dialogs struct {
	Userid      int    `json:"userid"`
	WithIm      int    `json:"sendto"`
	LastMessage string `json:"text"`
	Timestamp   int64  `json:"timestamp"`

	Time     string    `json:"time"`
	RandomID uuid.UUID `json:"random_id"`
	Username string    `json:"username"`
	Avatar   string    `json:"avatar"`
}

/*//Validate ...
func (w *Dialogs) Validate() error {
	return validation.ValidateStruct(
		w,
		validation.Field(&w.Text, validation.Required, validation.Length(1, 400)),
		validation.Field(&w.Userid, validation.Required),
	)
}*/

//Proccessing ...
func (w *Dialogs) Proccessing() error {

	goment.SetLocale("ru")
	time, err := goment.Unix(w.Timestamp)
	if err != nil {
		return err
	}
	w.Time = time.FromNow()

	return nil
}

//GenerateUUID ...
func (w *Dialogs) GenerateUUID() error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	w.RandomID = uuid
	return nil
}
