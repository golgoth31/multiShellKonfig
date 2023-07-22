package shell

import "regexp"

type ContextDef struct {
	Name     string
	FileID   string
	FilePath string
}

type ShellContextList []ContextDef

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))" //nolint

var re = regexp.MustCompile(ansi)
