/**
 * The MIT License (MIT)
 *
 * Copyright (c) 2015 Jane Lee <jane@webconn.me>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

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

	g = &gpio.Gpio{28, gpio.OUT}
	err := g.Open()
	if err != nil {
		log.Println(err)
	}
	defer g.Close()

	client = webconn.NewClient("http://192.168.4.180:3006/v01/do1ch/80")
	client.AddHandler("do",D1_OUT)

	client.Run()
}
