package shell

type ContextDef struct {
	Name     string
	FileID   string
	FilePath string
}

type ShellContextList []ContextDef

// make ShellContextList sortable
func (a ShellContextList) Len() int           { return len(a) }
func (a ShellContextList) Less(i, j int) bool { return a[i].Name < a[j].Name }
func (a ShellContextList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
