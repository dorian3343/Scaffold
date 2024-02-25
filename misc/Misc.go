package misc

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

func WelcomeMessage() {
	bold := "\033[1m"
	reset := "\033[0m"
	text := "Thank you for using Scaffold! Happy dev'ing <3 !"
	fmt.Println(bold + text + reset + "\n")
}

func StartHttp(port int) {
	log.Info().Msg("Starting HTTP server...")
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Fatal Error with http server")
	}

}
