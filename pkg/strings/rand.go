package strings

import "math/rand"

// RandomString generates a pseudo random string of the given length using the given characters.
// The distribution is not perfect, but works for general purposes
func RandomString(source string, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = source[rand.Int()%len(source)]
	}
	return string(b)
}

const passwordLower = "abcdefghijkmnopqrstuvwxyz"
const passwordUpper = "ABCDEFGHJKLMNPQRSTUVWXYZ"
const passwordNum = "23456789"
const passwordSym = "!@#%?+=_"
const passwordBytes = passwordLower + passwordUpper + passwordNum + passwordSym

// PasswordString generates a pseudo random password with the given length using characters that are common in passwords.
// We attempt to leave out letters that are easily visually confused.
// Specific letters excluded are lowercase l, upper case I and the number 1, upper case O and the number 0
// Also, only easily identifiable and describable symbols are used.
// It also tries to protect against accidentally creating an easily guessed value by making sure the password has at
// least one lower-case letter, one upper-case letter, one number, and one symbol.
// n must be at least 4
func PasswordString(n int) string {
	if n < 4 {
		panic ("n must be at least 4")
	}
	b := make([]byte, n)
	b[0] = passwordLower[rand.Int()%len(passwordLower)]
	b[1] = passwordUpper[rand.Int()%len(passwordUpper)]
	b[2] = passwordNum[rand.Int()%len(passwordNum)]
	b[3] = passwordSym[rand.Int()%len(passwordSym)]
	for i := 4; i < len(b); i++ {
		b[i] = passwordBytes[rand.Int()%len(passwordBytes)]
	}
	rand.Shuffle(n, func(i,j int) {
		temp := b[i]
		b[i] = b[j]
		b[j] = temp
	})

	return string(b)
}