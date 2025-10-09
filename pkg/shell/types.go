package shell

type ContextDef struct {
	Name     string
	FileID   string
	FilePath string
}

type ShellContextList []ContextDef
