package utils

import "math/rand"

func GenerateName() string {
	stringPool := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	nameLength := 5
	name := make([]byte, nameLength)

	for i := range name {
		name[i] = stringPool[rand.Intn(len(stringPool))]
	}

	return string(name)

}
