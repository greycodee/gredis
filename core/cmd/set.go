package cmd

func set() func(c string) []byte {
	return func(c string) []byte {
		cmdList := readCmdString(c)
		return buildCommand(cmdList...)
	}
}

func get() func(c string) []byte {
	return func(c string) []byte {
		cmdList := readCmdString(c)
		return buildCommand(cmdList...)
	}
}
