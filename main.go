package main

import (
	"github.com/tkanos/gonfig"
	"log"
	"net/http"
	"time"
)

type Config struct {
	Email string `json:"email"`
	Sleep int `json:"sleep"`
	Links []string `json:"link"`
}
func main() {

	config := Config{}
	err := gonfig.GetConf("config.json", &config)
	if err != nil {
		panic(err)
	}

	c := make(chan string) // channel for main routine and child routines communication

	for _, link := range config.Links {
        go checkLink(link, c) // go uses one CPU core by default
	}

	for l := range c { // wait for the channel then run the for loop
		go func(link string) { // main routine to child routine
			time.Sleep(time.Duration(config.Sleep) * time.Second) // 5 seconds sleep TODO: 5 should be parameter
			go checkLink(link,c) // instead of go checkLink(<-c,c)
		}(l) // passing l as an argument
	}
}

func checkLink(link string, c chan string){ // TODO: add sleep as a paramater
	_, err := http.Get(link)
	if err != nil {
		log.Println(link, "might be down") // TODO: warn via email
		c <- link // sending data to channel
		return
	}

	log.Print(link, " is up and runnig!")
	c <- link // sending data to channel
}

