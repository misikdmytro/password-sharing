package helper

import (
	"errors"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type RandomGenerator interface {
	RandomString(int) (string, error)
}

type RandomGeneratorFactory interface {
	NewRandomGenerator() RandomGenerator
}

type randomGenerator struct {
}

type randomGeneratorFactory struct {
}

func NewRandomFactory() RandomGeneratorFactory {
	return &randomGeneratorFactory{}
}

func (f *randomGeneratorFactory) NewRandomGenerator() RandomGenerator {
	return &randomGenerator{}
}

var symbols = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func (r *randomGenerator) RandomString(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("requested string should have positive length")
	}

	chars := make([]rune, length)
	for i := range chars {
		chars[i] = symbols[rand.Intn(len(symbols))]
	}

	return string(chars), nil
}
