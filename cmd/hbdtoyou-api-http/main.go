package main

import (
	"flag"
	"hbdtoyou/cmd/hbdtoyou-api-http/server"
	"os"
)

func main() {
	p := flag.String("secret-path", "", "secret file path")
	t := flag.Bool("config-test", false, "run config test")

	flag.Parse()

	os.Exit(server.Run(server.Option{
		SecretPath: *p,
		ConfigTest: *t,
	}))
}
