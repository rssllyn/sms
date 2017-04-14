package commandsender

import (
	"github.com/tarm/serial"
	"log"
	"regexp"
	"time"
)

var (
	commands           = make(chan *Command, 1)
	results            = make(chan *CommandResult, 1)
	receivedFrames     = make(chan string, 50)
	inCommand          = false
	COMMAND_RESULT_OK  = regexp.MustCompile(`\r\nOK`)
	COMMAND_RESULT_ERR = regexp.MustCompile(`((\r\n)|^)\+CME ERROR`)
)

func Start(p *serial.Port) {
	for {
		select {
		case command := <-commands:
			inCommand = true

			log.Printf("sending: %q", command.CommandString)
			_, err := p.Write([]byte(command.CommandString))
			if err != nil {
				log.Fatalf("failed to send AT command: %s, %s", command.CommandString, err)
			}

			result := CommandResult{Success: true}
			for i := 0; i < command.ExpectedResultFrames; i++ {
				frame := <-receivedFrames
				result.Frames = append(result.Frames, frame)

				if COMMAND_RESULT_ERR.MatchString(frame) {
					// 命令执行失败
					result.Success = false
					break
				} else if COMMAND_RESULT_OK.MatchString(frame) {
					// 命令执行成功但没有返回信息
					result.Success = true
					break
				}
			}
			log.Printf("result of command %q: %q", command, result)
			results <- &result

			inCommand = false
		}
		// AT命令发送时间间隔
		time.Sleep(300 * time.Millisecond)
	}
}

func Raw(command string, resultFramesExpected int) *CommandResult {

	commands <- &Command{
		CommandString:        command,
		ExpectedResultFrames: resultFramesExpected,
	}
	return <-results
}

func PushFrame(frame string) {
	if inCommand {
		receivedFrames <- frame
	}
}
