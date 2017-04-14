package resulthandler

import (
	"github.com/rssllyn/sms/commandsender"
	"github.com/axgle/mahonia"
	"github.com/tarm/serial"
	"log"
	"regexp"
	"time"
)

const (
	CARRIAGE_RETURN = 0x0D
	LINE_FEED       = 0x0A
)

var (
	RESULT           = regexp.MustCompile(`^\+(\w+): `)
	RESULT_DATA_CMT  = regexp.MustCompile(".*\r\n([0-9A-F]+)")
	RESULT_DATA_CMTI = regexp.MustCompile(`"(\w+)",(\d+)`)
	RESULT_DATA_CMGR = regexp.MustCompile(".*\r\n([0-9A-F]+)")
	RESULT_CODE_CMT  = "CMT"
	RESULT_CODE_CMTI = "CMTI"
	RESULT_CODE_CMGR = "CMGR"
	decoder          = mahonia.NewDecoder("gbk")
	partial_frame    = make(chan struct{})
	frame_buf        = make([]byte, 8024)
	frame_end        = 0
)

func Start(p *serial.Port) {
	buf := make([]byte, 1024)

	detect := frameBoundaryDetector()

	bytesRead := 0
	var lastLine *line_position = nil
	currentLine := &line_position{-1, -1, false}

	var lastFrameCanceller, lastFrameConfirm chan struct{}

	for {
		n, err := p.Read(buf[bytesRead:])
		if err != nil {
			log.Fatal(err)
		}
		// log.Printf("%d bytes arrived, %X", n, buf[:bytesRead+n])
		for i := bytesRead; i < bytesRead+n; i++ {
			if detect(buf[i]) {
				if currentLine.Start == -1 {
					// 找到帧头
					currentLine.Start = i - 1
					currentLine.StartWithNewLine = true
				} else {
					// 找到帧尾
					lastLine = &line_position{currentLine.Start, i, currentLine.StartWithNewLine}
					lastFrameCanceller, lastFrameConfirm = handleLine(buf[lastLine.GetDataStart():lastLine.GetDataEnd()])
					currentLine = &line_position{-1, -1, false}
				}
			}
			if lastLine != nil && i-lastLine.End == 2 {
				// 上一行之后又接收了两个字节
				if currentLine.Start == -1 {
					// 两个字节不是换行符
					// log.Println("potential frame cancelled")
					lastFrameCanceller <- struct{}{}
					currentLine.Start = lastLine.End + 1
					currentLine.StartWithNewLine = false
				} else {
					// 两个字节是换行符
					// log.Println("potential frame confirmed", currentLine.Start, lastLine.End)
					lastFrameConfirm <- struct{}{}
				}
				// last line只保留到接收到新的两个字节
				lastLine = nil
			}
		}
		bytesRead += n

		// 如果last line已经过期，则从buf中清除已经处理过的数据
		if lastLine == nil && currentLine.Start > 0 {
			copy(buf, buf[currentLine.Start:bytesRead])
			bytesRead -= currentLine.Start
			currentLine.Start = 0
		}
	}
}

func handleLine(potentialFrame []byte) (chan struct{}, chan struct{}) {
	log.Printf("potential frame %X", potentialFrame)
	canceller := make(chan struct{}, 2)
	confirm := make(chan struct{}, 2)
	if frame_end > 0 {
		copy(frame_buf[frame_end:], []byte{CARRIAGE_RETURN, LINE_FEED})
		frame_end += 2
	}
	copy(frame_buf[frame_end:], potentialFrame)
	frame_end += len(potentialFrame)
	go func() {
		select {
		case <-time.After(100 * time.Millisecond):
			// 等待一定的时间多接收两个字节，根据那两个字节是否是换行符判断potentialFrame是否是完整的命令帧
			// log.Println("potential frame confirmed because of timeout")
			handleFrame(frame_buf[:frame_end])
			frame_end = 0
		case <-confirm:
			// 接下来的两个字节是换行符，表明potentialFrame是一条完整的命令帧
			handleFrame(frame_buf[:frame_end])
			frame_end = 0
		case <-canceller:
			// 接下来的两个字节不是换行符，表明接下来的数据与potentialFrame是同一命令帧
		}
	}()
	return canceller, confirm
}

func handleFrame(frame []byte) {
	if len(frame) == 0 {
		log.Println("empty frame")
		return
	}
	// data := decoder.ConvertString(string(frame))
	data := string(frame)
	log.Printf("received raw: %X", frame)
	log.Println("received", string(data))
	// log.Println("pushing result")
	commandsender.PushFrame(data)

	matches := RESULT.FindStringSubmatch(data)
	if matches == nil {
		// log.Println("does not match", data)
		return
	}
	unsolicitedResultCode := matches[1]
	subdata := data[len(matches[0]):]
	// log.Println("sub data", subdata)
	switch unsolicitedResultCode {
	case RESULT_CODE_CMT:
		cmt(subdata)
	case RESULT_CODE_CMTI:
		cmti(subdata)
	case RESULT_CODE_CMGR:
		cmgr(subdata)
	}
}

func frameBoundaryDetector() func(byte) bool {
	lastByteIsCarriageReturn := false

	return func(b byte) (boundaryFound bool) {
		boundaryFound = lastByteIsCarriageReturn && b == LINE_FEED
		lastByteIsCarriageReturn = b == CARRIAGE_RETURN
		return
	}
}
