package encryption

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// bcrypt stuff

func HashPassword(password []byte) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	return hashedPassword, err
}

func CompareHashes(dbhashword []byte, inputpassword []byte) bool {

	err := bcrypt.CompareHashAndPassword(dbhashword, inputpassword)

	if err != nil {
		return false
	}

	return true

}

func GeneratePassword(length int) []byte {
	// from https://yourbasic.org/golang/generate-random-string/
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	specials := "~=+%^*/()[]{}/!@#$?|"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		digits + specials

	if length == 0 {
		length = 8
	}

	buf := make([]byte, length)
	buf[0] = digits[rand.Intn(len(digits))]
	buf[1] = specials[rand.Intn(len(specials))]
	for i := 2; i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	str := string(buf) // E.g. "3i[g0|)z"

	return []byte(str)
}
