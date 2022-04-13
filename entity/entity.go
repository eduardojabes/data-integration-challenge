package entity

import "github.com/google/uuid"

type Companies struct {
	ID      uuid.UUID
	Name    string
	Zip     string
	Website string
}
