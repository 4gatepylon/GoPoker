package utils

import (
	"fmt"
	"math/rand"
	"strings"
)

// Creates a random string of length n from the alphabet including upper and lower case
// copied from `https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go`
// (optimized to utilize uint64 random numbers are batches)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

// Creates a random uint64
func RandInt64() uint64 {
	return rand.Uint64()
}

var animals = []string{
	"giraffe",
	"chimpanzee",
	"flamingo",
	"frog",
	"triceratops",
	"leopard",
	"camel",
	"human",
	"sparrow",
	"aardavark",
	"weasel",
	"goose",
	"puppy",
	"wombat",
	"spider",
	"zebra",
}

var adjectives = []string{
	"pink",
	"excited",
	"energetic",
	"green",
	"blue",
	"angry",
	"surprised",
	"humble",
	"soaring",
	"majestic",
	"emphatic",
	"inquisitive",
	"grey",
	"spooky",
	"lethargic",
	"dancing",
	"fencing",
	"acrobatic",
}

func RandAdjAnimal(separator *string) string {
	sep := "-"
	if separator != nil {
		sep = *separator
	}
	adjective := adjectives[rand.Intn(len(adjectives))]
	animal := animals[rand.Intn(len(animals))]
	return fmt.Sprintf("%s%s%s", adjective, sep, animal)
}

var verbs = []string{
	"tread",
	"jump",
	"run",
	"think",
	"be",
	"beep",
	"boop",
	"sleep",
	"read",
	"fly",
	"roll",
	"float",
}

var adverbs = []string{
	"lightly",
	"quickly",
	"slowly",
	"merrily",
	"tragically",
	"amazingly",
	"swimmingly",
	"auspiciously",
	"magnanimously",
	"surprisingly",
	"like a ravioli",
}

func RandVerbAdv(separator *string) string {
	sep := "-"
	if separator != nil {
		sep = *separator
	}
	verb := verbs[rand.Intn(len(verbs))]
	adverb := adverbs[rand.Intn(len(adverbs))]
	return fmt.Sprintf("%s%s%s", verb, sep, adverb)
}
