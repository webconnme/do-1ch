package main

import (
	"github.com/webconnme/go-webconn"
	"github.com/webconnme/go-webconn-gpio"
	"log"
)

var client webconn.Webconn
var g *gpio.Gpio

func D1_OUT(buf []byte) error{

	data := string(buf)
	log.Println(">>>out data : ",data)

	if data == "high" {
		if err := g.Out(gpio.HIGH); err != nil {
			log.Println(err)
			return err
		}
	} else if data == "low" {
		if err := g.Out(gpio.LOW); err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func main() {

	g = &gpio.Gpio{248, gpio.OUT}
	err := g.Open()
	if err != nil {
		log.Println(err)
	}
	defer g.Close()

	client = webconn.NewClient("http://192.168.4.180:3006/v01/do1ch/80")
	client.AddHandler("do",D1_OUT)

	client.Run()
}
