package main

import (
	"github.com/tkanos/gonfig"
	"log"
	"net/http"
	"time"
)

type Config struct {
	Email []string `json:"email"`
	Sleep int `json:"sleep"`
	Links []string `json:"link"`
	From string `json:"email"`
}

func main() {

	config:= Config{}
	err := gonfig.GetConf("config.json", &config)
	if err != nil {
		panic(err)
	}

	c := make(chan string) // channel for main routine and child routines communication

	for _, link := range config.Links {
        go checkLink(link, c,config) // go uses one CPU core by default
	}

	for l := range c { // wait for the channel then run the for loop
		go func(link string) { // main routine to child routine
			time.Sleep(time.Duration(config.Sleep) * time.Second) // 5 seconds sleep
			go checkLink(link,c,config) // instead of go checkLink(<-c,c)
		}(l) // passing l as an argument
	}
}

func checkLink(link string, c chan string,config Config){
	_, err := http.Get(link)
	if err != nil {
		log.Println(link, "might be down")
		go sendEmail(link,config)
		c <- link // sending data to channel
		return
	}

	log.Print(link, " is up and runnig!")
	c <- link // sending data to channel
}

func sendEmail(link string, config Config ) {
	body := link + "is down!"
	sender := NewSender("<YOUR EMAIL ADDRESS>", "<YOUR EMAIL PASSWORD>")

	Receiver := config.Email

	Subject := "Warning From Status Checker App"
	message := `
	<!DOCTYPE HTML PULBLIC "-//W3C//DTD HTML 4.01 Transitional//EN">
	<html>
	<head>
	<meta http-equiv="content-type" content="text/html"; charset=ISO-8859-1">
	</head>
	<body>` +body+`<br>
	<div class="moz-signature"><i><br>
	<br>
	Regards<br>
	Ogguz<br>
	<i></div>
	</body>
	</html>
	`
	bodyMessage := sender.WriteHTMLEmail(Receiver, Subject, message)

	sender.SendMail(Receiver, Subject, bodyMessage)
}