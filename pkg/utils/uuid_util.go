package utils

import "github.com/google/uuid"

func GenerateUUID() string {
	return uuid.New().String()
}

func GetUUIDFromString(s string) uuid.UUID {
	id, _ := uuid.Parse(s)
	return id
}

func UUIDToString(id uuid.UUID) string {
	return id.String()
}
