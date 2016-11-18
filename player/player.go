package player

import "log"

var MsgChan chan string

func Setup() {
	MsgChan = make(chan string)
	go Run()
}

func Run() {
	for {
		select {
		case message := <-MsgChan:
			log.Printf("Player got message %q.", message)
			if message == "quit" {
				log.Print("Player says goodbye.")
				return
			}
		}
	}
}
