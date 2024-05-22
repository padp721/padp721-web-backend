package utilities

import (
	"fmt"
	"os"
)

func GeneratePasswordString(password string, username string) string {
	bcryptSalt := os.Getenv("BCRYPT_SALT")
	passwordString := fmt.Sprintf(
		"%v:%v:%v",
		password,
		username,
		bcryptSalt,
	)

	return passwordString
}
