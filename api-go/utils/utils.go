package utils

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"regexp"
	"time"
)

const IdLength = 32

func GenUuid() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x", b[0:4], b[4:6], (b[6]&0x0f)|0x40, (b[8]&0x3f)|0x80, b[10:])
}

func IsUuid(id string) bool {
	re := regexp.MustCompile(fmt.Sprintf("^[a-z0-9]{%d}$", IdLength))
	return re.MatchString(id)
}

func IsValidString(str string) bool {
	re := regexp.MustCompile(`^[0-9a-zA-Z_!?:=+\\-,.\sА-Яа-я]{1,64}$`)
	return re.MatchString(str)
}

func RandomName() string {
	adjectives := []string{"Amiable", "Blissful", "Cheerful", "Delightful", "Enchanting", "Friendly", "Gracious", "Harmonious", "Invigorating", "Jovial", "Kindhearted"}
	animals := []string{"Dragonfly", "Alpaca", "Rabbit", "Vulture", "Jackrabbit", "Bunny", "Butterfly"}

	adj, _ := rand.Int(rand.Reader, big.NewInt(int64(len(adjectives))))
	anim, _ := rand.Int(rand.Reader, big.NewInt(int64(len(animals))))

	return adjectives[adj.Int64()] + " " + animals[anim.Int64()]
}

var startTime time.Time

func Timer() float64 {
	if startTime.IsZero() {
		startTime = time.Now()
		return 0
	}
	return time.Since(startTime).Seconds()
}

func Log(str string) {
	log.Printf("%s (%0.3f): %s\n", time.Now().Format("2006-01-02 15:04:05.000"), Timer(), str)
}

func OhMySimpleHash(str string) uint32 {
	var hash uint32
	for i := 0; i < len(str); i++ {
		hash = (hash << 5) - hash + uint32(str[i])
	}
	return hash
}
