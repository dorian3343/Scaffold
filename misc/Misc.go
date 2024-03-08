package misc

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net"
	"net/http"
	"strconv"
	"strings"
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
func StartHttp(port int, static string) {
	mux := http.NewServeMux()

	// Handler to serve HTML files without .html extension
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			path = "/index.html" // Serve root
		} else {
			if strings.HasSuffix(path, "/") {
				path = strings.TrimSuffix(path, "/")
			}
			if !strings.HasSuffix(path, ".html") {
				path += ".html" // Append .html extension if not present
			}
		}
		http.ServeFile(w, r, static+path)
	})
	ip := GetLocalIP()
	url := fmt.Sprintf("http://%s:%d", ip, port)
	log.Info().Msgf("Server is running at: %s", url)

	err := http.ListenAndServe(":"+strconv.Itoa(port), mux)
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
