package main

import (
	"flag"
	"github.com/stvp/go-toml-config" 
	"github.com/ermine/bullxmpp/client"
	"github.com/ermine/bullxmpp/extensions/iqversion"
	"fmt"
	"os"
)

var (
	myjid = config.String("jid", "user@server")
	password = config.String("password", "secret")
)

func loadConfig() {
	var path = flag.String("config", "bot.conf", "path to configuration file")
	flag.Parse()
	err := config.Parse(*path)
	if err != nil {
		panic(err)
	}
}

func main() {
	loadConfig()
	cli, err := client.New(*myjid, *password, "ru", true, true)
	checkError(err)
	for {
		p, err := cli.Read()
		checkError(err)
		if p == nil { break }
		switch p := p.(type) {
		case *client.Iq:
			if _, ok := p.Payload.(*iqversion.Version); ok {
				if p.Type != nil && *p.Type == client.IqTypeGet {
					version := iqversion.New("bull", "0.0001", "Linux")
					typ := client.IqTypeResult
					iq := client.Iq{To:p.From, Id:p.Id, Payload:version, Type: &typ}
					err := cli.Send(&iq)
					checkError(err)
				}
			}
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}



















