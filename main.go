package main

import (
	"./config"
	"./robot"
	"./utils/env"
	"./utils/log"
	"flag"
	"fmt"
	"io/ioutil"
)

var err error

func main() {
	//init log
	logger, err := log.NewLogger(env.LogPath, "debug")
	env.ErrExit(err)
	log.SetDefault(logger)

	// init config
	filename := flag.String("f", "config.json", "[file name] default: config.json")
	flag.Parse()
	fmt.Println("config filename:", *filename)
	data, err := ioutil.ReadFile(*filename)
	if err != nil {
		fmt.Println("initconfig failed! %v", err)
	}
	//set data
	env.ErrExit(config.InitConfig(data))

	revealrobot.Init()

	//robot run
	revealrobot.RevealRobot()
}
