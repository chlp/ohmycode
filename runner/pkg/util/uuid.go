package util

import (
	"crypto/rand"
	"fmt"
	"log"
	"regexp"
)

const IdLength = 32

var uuidRegex = regexp.MustCompile(fmt.Sprintf("^[a-z0-9]{%d}$", IdLength))

func GenUuid() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal("uuid: cannot generate uuid", err)
	}
	return fmt.Sprintf("%08x%04x%04x%04x%12x", b[0:4], b[4:6], (b[6]&0x0f)|0x40, (b[8]&0x3f)|0x80, b[10:])
}

func IsUuid(id string) bool {
	return uuidRegex.MatchString(id)
}
