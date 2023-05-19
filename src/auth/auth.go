package auth

import (
	"os"
	"strings"
)

func IsAuthenticated(id string) bool {
	authed := strings.Split(os.Getenv("AUTHORIZED"), ",")
	for _, a := range authed {
		if a == id {
			return true
		}
	}
	return false
}
