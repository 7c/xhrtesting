package shared

import (
	"math/rand"
	"os"
	"path"
	"strings"
)

// Generates a random string of a given length
func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

// Generates a random slice of ints with a given length
func RandomArray(len int) []int {
	array := make([]int, len)
	for i := range array {
		array[i] = rand.Intn(100)
	}
	return array
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func LECerts(domain string) (cert string, key string) {
	certFile := path.Join("/etc/letsencrypt/live/", strings.ToLower(domain), "/fullchain.pem")
	keyFile := path.Join("/etc/letsencrypt/live/", strings.ToLower(domain), "/privkey.pem")
	if FileExists(certFile) && FileExists(keyFile) {
		return certFile, keyFile
	}
	return "", ""
}
