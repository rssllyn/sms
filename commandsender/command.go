package commandsender

type Command struct {
	CommandString        string
	ExpectedResultFrames int
}

type CommandResult struct {
	Frames  []string
	Success bool
}
