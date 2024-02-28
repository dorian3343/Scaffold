package misc

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

// A simple message to display on startup
func WelcomeMessage() {
	text := "Thank you for using Scaffold! Happy dev'ing <3 !"
	fmt.Println(text)
}

// Startup a http server on a port
func StartHttp(port int) {
	log.Info().Msg("Starting HTTP server...")
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Fatal Error with http server")
	}

}

// Capitalize a string
func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-'a'+'A') + s[1:]
}
