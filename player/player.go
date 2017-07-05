package player

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/lhz/sidpicker/hvsc"
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
var CurrentSong int
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
				StartTime = time.Now()
				playCmd = exec.Command("sidplayfp", "-ol", "-t0", "-o"+msg.Args[1], msg.Args[0])
				playCmd.Stdout = os.Stdout
				if err := playCmd.Start(); err != nil {
					log.Print("Failed to start player process: ", err)
				}
			case STOP_COMMAND:
				stopCommand(playCmd)
			case QUIT_COMMAND:
				stopCommand(playCmd)
				return
			}
		}
	}
}

func PlaySub(subTune int) {
	Stop()
	if subTune < 1 {
		subTune = CurrentTune.Header.StartSong
	}
	Playing = true
	CurrentSong = subTune
	MsgChan <- PlayerMsg{Command: PLAY_COMMAND, Args: []string{CurrentTune.FullPath(), strconv.Itoa(subTune)}}
}

func Play(index, subTune int) {
	Stop()
	tune := hvsc.FilteredTunes[index]
	CurrentTune = &tune
	if subTune < 1 {
		subTune = tune.Header.StartSong
	}
	CurrentSong = subTune
	Playing = true
	MsgChan <- PlayerMsg{Command: PLAY_COMMAND, Args: []string{tune.FullPath(), strconv.Itoa(subTune)}}
}

func PrevSong() {
	if Playing && CurrentSong > 1 {
		PlaySub(CurrentSong - 1)
	}
}

func NextSong() {
	if Playing && CurrentSong < CurrentTune.Header.Songs {
		PlaySub(CurrentSong + 1)
	}
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
	return TimeFormat(time.Since(StartTime))
}

func SongLength() string {
	if CurrentTune == nil || CurrentSong < 1 {
		return ""
	}
	return TimeFormat(CurrentTune.SongLengths[CurrentSong-1])
}

func TimeFormat(duration time.Duration) string {
	seconds := int(duration.Seconds())
	return fmt.Sprintf("%02d:%02d", seconds/60, seconds%60)
}

func stopCommand(cmd *exec.Cmd) {
	if cmd != nil {
		if runtime.GOOS == "windows" {
			cmd.Process.Kill()
		} else {
			cmd.Process.Signal(os.Interrupt)
			cmd.Wait()
		}
	}
}
