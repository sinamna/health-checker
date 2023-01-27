package utils

import (
	"crypto/md5"
	"fmt"
)

func HashString(value string) string {
	md5Hash := md5.Sum([]byte(value))
	return fmt.Sprintln(md5Hash)[:40]
}
