package utilities

import (
	"fmt"
	"os"
)

func GeneratePasswordString(password string, username string, id string) string {
	bcryptSalt := os.Getenv("BCRYPT_SALT")
	passwordString := fmt.Sprintf(
		"%v:%v:%v:%v",
		password,
		username,
		id,
		bcryptSalt,
	)

	return passwordString
}
