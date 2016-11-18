package player

import "log"

var Chan chan string

func Setup() {
	Chan = make(chan string)
}

func Run() {
	for {
		select {
		case message := <-Chan:
			log.Printf("Player got message %q.", message)
		}
	}
}
