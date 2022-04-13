package entity

import "github.com/google/uuid"

type Company struct {
	ID   uuid.UUID
	Name string
	Zip  string
}
