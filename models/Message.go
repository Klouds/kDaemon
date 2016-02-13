package models

type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}
