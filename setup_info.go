package webinit

type SetupInfo struct {
	RootFolder    string //root folder that hold controllers,views and public 
	Addr          string
	HotReloadView bool
	DelimLeft     string
	DelimRight    string
}
