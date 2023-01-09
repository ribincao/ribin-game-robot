package utils

import "fmt"

func GeneWebsocketURL(ip string, port int32) string {
	return fmt.Sprintf("ws://%s:%v", ip, port)
}
