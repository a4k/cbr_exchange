package model

import (
	"encoding/json"
)

type Popmail struct {
	Id       	string 	`json:"id" db:"id"`
	Uidl   		string 	`json:"uidl" db:"uidl"`
	Type   		string 	`json:"type" db:"type"`
	Date    	int64 	`json:"date" db:"date"`
	CreateAt 	int64  	`json:"create_at" db:"create_at"`
}

func (js *Popmail) ToJson() string {
	if b, err := json.Marshal(js); err != nil {
		return ""
	} else {
		return string(b)
	}
}

func (u *Popmail) PreSave() {
	if u.Id == "" {
		u.Id = NewId()
	}
}