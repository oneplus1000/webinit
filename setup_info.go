package webinit

type SetupInfo struct {
	RootFolder                                 string //RootFolder, root folder that hold controllers,views and public
	Addr, HotReloadView, DelimLeft, DelimRight string
}
