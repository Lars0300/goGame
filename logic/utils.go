package logic

import (
	"math/rand"
	"time"
)

const (
	minDuration = 10 * time.Second
	maxDuration = 60 * time.Second
)

var (
	HashToGame map[string]*Game   = make(map[string]*Game)
	OpenGames  map[*Game]struct{} = make(map[*Game]struct{})
	hashLetters []rune= []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func generateHash() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	hash := make([]rune, 5)
	for i := range hash {
		hash[i] = hashLetters[r.Intn(len(hashLetters))]
	}
	return string(hash)
}

func generateUsername() string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz1234567890")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	name := make([]rune, 8)
	for i := range name {
		name[i] = letters[r.Intn(len(letters))]
	}
	return string(name)
}

func isValidHash(suggestedHash string) bool{
	if len(suggestedHash) != 5{
		return false
	}
	for _, letter := range suggestedHash{
		if letter < 'A' || letter > 'Z'{
			return false
		}
	}
	return true
}
