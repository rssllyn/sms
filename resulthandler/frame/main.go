package frame

var (
	frames = make(chan string)
	bytes  = make(chan byte)
	buf    = make([]byte, 1024)
)

func Start() {
	return frames
}
