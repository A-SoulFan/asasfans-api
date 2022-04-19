package smsclient

import (
	"fmt"
	"math/rand"
	"strings"
)

type CodeGenerator interface {
	Generate(length int) string
}

var randomNumberCandidate = [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

type RandomNumberCodeGenerator struct {
	randomCandidate       [10]byte
	randomCandidateLength int
}

func NewRandomNumberCodeGenerator() CodeGenerator {
	temp := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	return &RandomNumberCodeGenerator{
		randomCandidate:       randomNumberCandidate,
		randomCandidateLength: len(temp),
	}
}

func (m *RandomNumberCodeGenerator) Generate(length int) string {
	var sb strings.Builder

	for i := 0; i < length; i++ {
		_, _ = fmt.Fprint(&sb, m.randomCandidate[rand.Intn(len(m.randomCandidate))])
	}

	return sb.String()
}
