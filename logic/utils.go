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
	hashToGame map[string]*Game = make(map[string]*Game)
	openGames map[*Game]struct{} = make(map[*Game]struct{})
)

func generateHash() string {
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	hash := make([]rune, 5)
	for i := range hash {
		hash[i] = letters[r.Intn(len(letters))]
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

