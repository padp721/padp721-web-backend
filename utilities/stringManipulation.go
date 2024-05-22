package utilities

import (
	"fmt"
	"os"
)

func GeneratePasswordString(password string, username string, name string) string {
	bcryptSalt := os.Getenv("BCRYPT_SALT")
	passwordString := fmt.Sprintf(
		"%v:%v:%v:%v",
		password,
		username,
		name,
		bcryptSalt,
	)

	return passwordString
}
