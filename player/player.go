package player

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/lhz/considerate/hvsc"
	//"io/ioutil"
)

const (
	PLAY_COMMAND = iota
	STOP_COMMAND
	QUIT_COMMAND
)

type PlayerMsg struct {
	Command int
	Args    []string
}

var CurrentTune *hvsc.SidTune
var StartTime time.Time
var MsgChan chan PlayerMsg
var Playing = false

func Run() {
	MsgChan = make(chan PlayerMsg)

	var playCmd *exec.Cmd

	for {
		select {
		case msg := <-MsgChan:
			switch msg.Command {
			case PLAY_COMMAND:
				playCmd = exec.Command("/usr/bin/sidplay2", msg.Args[0])
				playCmd.Stdout = os.Stdout
				if err := playCmd.Start(); err != nil {
					log.Print("Failed to start player process: ", err)
				}
				StartTime = time.Now()
			case STOP_COMMAND:
				stopCommand(playCmd)
			case QUIT_COMMAND:
				stopCommand(playCmd)
				return
			}
		}
	}
}

func Play(index, subTune int) {
	Stop()
	tune := hvsc.FilteredTunes[index]
	CurrentTune = &tune
	Playing = true
	MsgChan <- PlayerMsg{Command: PLAY_COMMAND, Args: []string{tune.FullPath(), strconv.Itoa(subTune)}}
}

func Stop() {
	Playing = false
	MsgChan <- PlayerMsg{Command: STOP_COMMAND, Args: []string{}}
}

func Quit() {
	Playing = false
	MsgChan <- PlayerMsg{Command: QUIT_COMMAND, Args: []string{}}
}

func Elapsed() string {
	duration := time.Since(StartTime)
	seconds := int(duration.Seconds())
	s := seconds % 60
	m := seconds / 60
	return fmt.Sprintf("%02d:%02d", m, s)
}

func stopCommand(cmd *exec.Cmd) {
	if cmd != nil {
		cmd.Process.Signal(os.Interrupt)
		cmd.Wait()
	}
}
