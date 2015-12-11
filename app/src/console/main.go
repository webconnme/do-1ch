package main

import (
	"encoding/json"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"github.com/webconnme/go-webconn"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

/*
#include <stdio.h>
#include <termio.h>

struct termios save;
int saved = 0;

void saveTerm(void) {
    saved = 1;
    tcgetattr(0,&save);
}

void restoreTerm(void) {
    if (saved == 1) {
        tcsetattr(0, TCSAFLUSH, &save);
    }
}

int getch(void) {
    char ch;
    struct termios buf;

    saveTerm();
    buf = save;
    buf.c_lflag &= ~(ICANON|ECHO);
    buf.c_cc[VMIN] = 1;
    buf.c_cc[VTIME] = 0;
    tcsetattr(0, TCSAFLUSH, &buf);

    ch = getchar();

    restoreTerm();
    return ch;
}
*/
import "C"

func HandleSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	raiseCount := 0

	for {
		// Block until a signal is received.
		select {
		case s := <-c:
		//fmt.Println("Signal Raised: ", s)
			switch s {
			case syscall.SIGINT:
				raiseCount++
				//fmt.Printf("Raise count: %d\n", raiseCount)
				if raiseCount >= 0 {
					C.restoreTerm()
					os.Exit(0)
				}
			case syscall.SIGKILL:
				fallthrough
			case syscall.SIGTERM:
				C.restoreTerm()
				os.Exit(0)
			}
		case <-time.NewTicker(100 * time.Millisecond).C:
			if raiseCount > 0 {

			}
			raiseCount = 0
			runtime.Gosched()
		}
	}
}

var context *zmq.Context
var sock *zmq.Socket

func OnReceive(s *zmq.Socket) error {
	buf, err := s.RecvBytes(0)
	if err != nil {
		return err
	}

	var messages []webconn.Message
	err = json.Unmarshal(buf, &messages)
	if err != nil {
		return err
	}

	for _, m := range messages {
		if m.Command == "do" {
			fmt.Printf(string(m.Data))
		}
	}
	return nil
}

func SendDo(b bool) error {

	var messages []webconn.Message
	if b {
		messages = append(messages, webconn.Message{"do", "high"})
	} else {
		messages = append(messages, webconn.Message{"do", "low"})
	}

	j, err := json.Marshal(messages)
	if err != nil {
		return err
	}

	_, err = sock.SendBytes(j, 0)
	if err != nil {
		return err
	}

	return nil
}

func HandleNetwork(url string) {
	var err error
	context, err = zmq.NewContext()
	if err != nil {
		log.Panic(err)
	}
	defer context.Term()

	sock, err = context.NewSocket(zmq.PAIR)
	if err != nil {
		log.Panic(err)
	}
	defer sock.Close()

	sock.Connect(url)

	reactor := zmq.NewReactor()
	reactor.AddSocket(sock, zmq.POLLIN, func(state zmq.State) error { return OnReceive(sock) })

	err = reactor.Run(time.Second)

	if err != nil {
		log.Panic(err)
	}
}

func HandleKeyboard() {
	for {
		ch := byte(C.getch())

		if ch == 'h' || ch == 'H' {
			SendDo(true)
			fmt.Println("Sent a digital output High")
		} else if ch == 'l' || ch == 'L' {
			SendDo(false)
			fmt.Println("\nSent a digital output Low")
		}
		fmt.Printf("[h] DO High , [l] DO Low : ")
	}
}

func main() {

	done := make(chan bool)

	fmt.Printf("[h] Digital Output High , [l] Digital Output Low : ")
	go HandleNetwork("tcp://192.168.4.180:3007")
	go HandleSignal()
	go HandleKeyboard()

	<-done
}
