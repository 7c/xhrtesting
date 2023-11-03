package shared

import (
	"os"
	"path"
	"strings"
)

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
