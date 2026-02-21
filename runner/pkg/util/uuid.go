package util

import (
	"crypto/rand"
	"log"
	"regexp"
)

const IdLength = 22
const idAlphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var idRegex = regexp.MustCompile(`^([0-9A-Za-z]{22}|[a-z0-9]{32})$`)

func GenUuid() string {
	b := make([]byte, IdLength)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal("id: cannot generate id", err)
	}
	for i := range b {
		b[i] = idAlphabet[b[i]%62]
	}
	return string(b)
}

func IsUuid(id string) bool {
	return idRegex.MatchString(id)
}
