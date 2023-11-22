package api

import "github.com/google/uuid"

func GenerateCorrelatorID() string {
	id := uuid.New()
	return id.String()
}
