package util

import (
	"crypto/rand"
	"math/big"
	"regexp"
)

func IsValidName(str string) bool {
	re := regexp.MustCompile(`^[0-9a-zA-Z_!?:=+\-,.\sА-Яа-яЁё]{1,64}$`)
	return re.MatchString(str)
}

func RandomName() string {
	adjectives := []string{"Amiable", "Blissful", "Cheerful", "Delightful", "Enchanting", "Friendly", "Gracious", "Harmonious", "Invigorating", "Jovial", "Kindhearted"}
	animals := []string{"Dragonfly", "Alpaca", "Rabbit", "Vulture", "Jackrabbit", "Bunny", "Butterfly"}

	adj, _ := rand.Int(rand.Reader, big.NewInt(int64(len(adjectives))))
	anim, _ := rand.Int(rand.Reader, big.NewInt(int64(len(animals))))

	return adjectives[adj.Int64()] + " " + animals[anim.Int64()]
}

func OhMySimpleHash(str string) uint32 {
	var hash uint32
	for i := 0; i < len(str); i++ {
		hash = (hash << 5) - hash + uint32(str[i])
	}
	return hash
}
