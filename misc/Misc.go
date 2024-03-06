package misc

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net"
	"net/http"
	"strconv"
)

// A simple message to display on startup
func WelcomeMessage() {
	text := "Thank you for using Scaffold! Happy dev'ing <3 !"
	fmt.Println(text)
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// Startup a http server on a port
func StartHttp(port int) {
	ip := GetLocalIP()
	url := fmt.Sprintf("http://%s:%d", ip, port)
	log.Info().Msgf("Server is running at: %s", url)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Fatal Error with http server")
	}
}

// Capitalize a string
func Capitalize(s string) string {
	if len(s) == 0 || (s[0] >= 'A' && s[0] <= 'Z') || !isLetter(s[0]) {
		return s
	}
	return string(s[0]-'a'+'A') + s[1:]
}

// Check if its a letter
func isLetter(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
}
