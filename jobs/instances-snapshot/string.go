package main

import "crypto/rand"

// RandomString is used to generate a random string containing upper and lower characters
// + number, of size n.
func RandomString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	bytes := make([]byte, n)

	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}

	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}

	return string(bytes)
}
