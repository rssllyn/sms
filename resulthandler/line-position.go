package resulthandler

type line_position struct {
	Start, End       int
	StartWithNewLine bool
}

func (l *line_position) GetDataStart() int {
	if l.StartWithNewLine {
		return l.Start + 2
	} else {
		return l.Start
	}
}

func (l *line_position) GetDataEnd() int {
	return l.End - 1
}
