package entity

import "time"

type ItemId int64

type Item struct {
	Id         ItemId    `json:"id"`
	Title      string    `json:"title"`
	Url        string    `json:"url"`
	ModifiedAt time.Time `json:"modifiedAt"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Items []*Item
