package sqlstore

import (
	//"database/sql"

	"fmt"
	"time"

	"github.com/katelinlis/MessagesBackend/internal/app/model"
)

//MessageRepository ...
type MessageRepository struct {
	store *Store
}

//Send ...
func (r *MessageRepository) Send(p *model.Message) error {

	if err := p.Validate(); err != nil {
		return err
	}
	p.GenerateUUID()

	var err2 = r.store.db.QueryRow(
		"INSERT INTO messages (userid,sendto,text,timestamp,random_id) VALUES ($1,$2,$3,$4,$5) RETURNING random_id",
		p.Userid,
		p.SendTo,
		p.Text,
		time.Now().Unix(),
		p.RandomID.String(),
	).Scan(&p.RandomID)

	return err2
}

//Update ...
func (r *MessageRepository) Update(p *model.Message) error {

	if err := p.Validate(); err != nil {
		return err
	}

	fmt.Println(p.RandomID.String())

	var err2 = r.store.db.QueryRow(
		"UPDATE wall set answer_to = $1 where random_id=$2 RETURNING random_id",
		p.AnswerTO.String(),
		p.RandomID.String(),
	).Scan(&p.RandomID)

	return err2
}

// GetMessages ...
func (r *MessageRepository) GetMessages(userid int, ImWith int, limit int, offset int) ([]model.Message, error) {
	messages := []model.Message{}
	// SELECT status,user2 from friends where (user1 = $1 and user2= $2) OR (user1 = $2 and user2= $1)
	var rows, err2 = r.store.db.Query("select userid,sendto,text,timestamp,random_id from messages where (userid = $1 and sendto= $2) OR (userid = $2 and sendto= $1) ORDER BY timestamp ASC limit $3 OFFSET $4", userid, ImWith, limit, offset)

	if err2 != nil {
		return messages, err2
	}

	for rows.Next() {
		message := model.Message{}
		err := rows.Scan(&message.Userid, &message.SendTo, &message.Text, &message.Timestamp, &message.RandomID)
		if err != nil {
			return messages, err
		}

		message.Proccessing()

		messages = append(messages, message)
	}

	return messages, err2
}

//GetDialogList ...
func (r *MessageRepository) GetDialogList(userid int, limit int, offset int) ([]model.Dialogs, error) {
	Dialogs := []model.Dialogs{}

	rows, err := r.store.db.Query("select * from (select DISTINCT on (sendto) sendto,userid,text,timestamp from((select DISTINCT on (sendto) sendto,userid,text,timestamp from messages where userid = $1 order by sendto, timestamp DESC )union(select DISTINCT on (userid) userid,sendto,text,timestamp from messages where sendto = $1 order by userid, timestamp DESC )) as b order by sendto, timestamp DESC) as bb order by timestamp desc",
		userid,
	)
	if err != nil {
		return Dialogs, err
	}

	for rows.Next() {
		dialog := model.Dialogs{}
		err := rows.Scan(&dialog.WithIm, &dialog.Userid, &dialog.LastMessage, &dialog.Timestamp)
		if err != nil {
			return Dialogs, err
		}

		dialog.Proccessing()

		Dialogs = append(Dialogs, dialog)

	}
	return Dialogs, err
}

//GetAnswersCount ...
func (r *MessageRepository) GetAnswersCount(PostID string) (int, error) {
	var count int
	err := r.store.db.QueryRow(
		"select COUNT(answer_to) from wall where answer_to = $1",
		PostID,
	).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, err
}

//ScanForNewMessage ...
func (r *MessageRepository) ScanForNewMessage(userid int, ImWith int, lastMessageTimestamp int) ([]model.Message, error) {
	newMessages := []model.Message{}

	var rows, err = r.store.db.Query("select userid,sendto,text,timestamp,random_id from messages where (userid = $1 and sendto= $2) OR (userid = $2 and sendto= $1)AND timestamp > $3  ORDER BY timestamp ASC limit $4", userid, ImWith, lastMessageTimestamp, 50)

	if err != nil {
		return newMessages, err
	}

	for rows.Next() {
		message := model.Message{}
		err := rows.Scan(&message.Userid, &message.SendTo, &message.Text, &message.Timestamp, &message.RandomID)
		if err != nil {
			return newMessages, err
		}

		message.Proccessing()

		newMessages = append(newMessages, message)
	}

	//post.RandomID = uuid
	return newMessages, err
}
