package webinit

type SetupInfo struct {
	//RootFolder, root folder that hold controllers,views and public
	RootFolder, Addr, HotReloadView, DelimLeft, DelimRight string
}
