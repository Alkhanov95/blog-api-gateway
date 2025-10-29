package model

import "time"

type Item struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Quantity  int       `json:"quantity"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
