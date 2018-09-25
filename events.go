package main

import (
	"encoding/json"
)

type EventHandler func(*Event)

type Event struct {
	Name string      `json:"event"`
	Data interface{} `json:"data"`
}

func NewEventFromRaw(rawData []byte) (*Event, error) {
	eve := &Event{}
	err := json.Unmarshal(rawData, eve)
	return eve, err
}

func (e *Event) Raw() []byte {
	raw, _ := json.Marshal(e)
	return raw
}