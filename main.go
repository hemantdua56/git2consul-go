package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/cleung2010/go-git2consul/config"
	"github.com/cleung2010/go-git2consul/consul"
	"github.com/cleung2010/go-git2consul/exit"
	"github.com/cleung2010/go-git2consul/repository"
)

var quit = make(chan bool)

func init() {
	exit.Listen(func(os.Signal) {
		log.Info("Caught interrupt signal, terminating git2consul")
		close(quit)
	})
}

func main() {
	var filename string
	var v bool
	var d bool

	flag.StringVar(&filename, "config", "", "path to config file")
	flag.BoolVar(&v, "v", false, "show version")
	flag.BoolVar(&d, "d", false, "enable debugging mode")
	flag.Parse()

	if d {
		log.SetLevel(log.DebugLevel)
	}

	if v {
		fmt.Println(Version)
		return
	}

	log.Infof("Starting git2consul version: %s", Version)

	if len(filename) == 0 {
		log.Fatal("No configuration file provided")
	}

	cfg, err := config.Load(filename)
	if err != nil {
		log.Error(err)
		close(quit)
	}

	repos, err := repository.LoadRepos(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = repos.WatchRepos()
	if err != nil {
		log.Fatal(err)
	}

	// Create Consul client
	client, err := consul.NewClient(cfg)
	client.WatchChanges(repos)

	//Wait for shutdown signal
	<-quit
}
