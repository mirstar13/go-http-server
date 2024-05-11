package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/mirstar13/go-http-server/server"
	"github.com/mirstar13/go-http-server/server/config"
)

var (
	dirMatch = regexp.MustCompile(`--directory\s+((?:[^/]*/)*)(.*)`)
)

func main() {
	fmt.Println("The logs of your app will appear here")

	cmdOptions := setServerOptions()

	cfg := config.NewConfig(cmdOptions...)
	server := server.NewServer(cfg)

	if err := server.Start(); err != nil {
		log.Println("failed to start server: %w\n", err)
		os.Exit(1)
	}

}

func setServerOptions() []config.Option {
	var options []config.Option
	cmdOptions := strings.Join(os.Args[1:], " ")

	if params := dirMatch.FindStringSubmatch(cmdOptions); len(params) > 0 {
		dir := params[1]
		options = append(options, config.WithDir(dir))
	}

	return options
}
