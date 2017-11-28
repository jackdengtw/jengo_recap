package util

import (
	"fmt"
)

func GetHostPort4Client(host string, port int) string {
	if host == "0.0.0.0" {
		host = "127.0.0.1"
	}
	return fmt.Sprintf("%v:%v", host, port)
}
