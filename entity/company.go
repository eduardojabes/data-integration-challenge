package entity

import "github.com/google/uuid"

type Companies struct {
	ID      uuid.UUID `json:"_id"`
	Name    string    `json:"name"`
	Zip     string    `json:"zipCode"`
	Website string    `json:"website"`
}
