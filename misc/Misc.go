package misc

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

func WelcomeMessage() {
	text := "Thank you for using Scaffold! Happy dev'ing <3 !"
	fmt.Println(text)
}

func StartHttp(port int) {
	log.Info().Msg("Starting HTTP server...")
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Fatal Error with http server")
	}

}
