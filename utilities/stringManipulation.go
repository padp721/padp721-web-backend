package utilities

import (
	"fmt"
	"os"
)

func GeneratePasswordString(password string, username string) string {
	passwordString := fmt.Sprintf(
		"%v:%v:%v",
		password,
		username,
		os.Getenv("BCRYPT_SALT"),
	)

	return passwordString
}
