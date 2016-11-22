package player

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	//"io/ioutil"
)

const (
	PLAY_COMMAND = iota
	QUIT_COMMAND
)

type PlayerMsg struct {
	Command int
	Args    []string
}


var MsgChan chan PlayerMsg

func Setup() {
	MsgChan = make(chan PlayerMsg)
	go Run()
}

func Run() {
	var playCmd *exec.Cmd
	for {
		select {
		case msg := <-MsgChan:
			//log.Printf("Player got message %v.", msg)
			if msg.Command == QUIT_COMMAND {
				if playCmd != nil {
					playCmd.Process.Signal(os.Interrupt)
					playCmd.Wait()
					// if err := playCmd.Process.Kill(); err != nil {
					// 	log.Print("Failed to kill player process: ", err)
					// }
				}
			}
			if msg.Command == PLAY_COMMAND {
				playCmd = exec.Command("/usr/bin/sidplay2", msg.Args[0])
				playCmd.Stdout = os.Stdout
				if err := playCmd.Start(); err != nil {
					log.Print("Failed to start player process: ", err)
				}
			}
		}
	}
	playCmd.Process.Kill()
}

func Play(path string, subTune int) {
	Quit()
	MsgChan <- PlayerMsg{Command: PLAY_COMMAND, Args: []string{path, strconv.Itoa(subTune)}}
}

func Quit() {
	MsgChan <- PlayerMsg{Command: QUIT_COMMAND, Args: []string{}}
}
