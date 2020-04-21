package model

import "time"

type JsonTime time.Time

func (j JsonTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(j).Format("2006-01-02 15:04") + `"`), nil
}

type Message struct {
	Id       int64    `json:"id" xorm:"pk autoincr"`
	From     string   `json:"from"`
	To       string   `json:"to"`
	Content  string   `json:"content"`
	SendTime JsonTime `json:"send_time" xorm:"created"`
	IsRead   int      `json:"is_read" xorm:"int"`
}
