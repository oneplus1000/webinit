package webinit

type SetupInfo struct {
	RootFolder            string //RootFolder, root folder that hold controllers,views and public
	Addr                  string //bind address  example  ":8080"
	DelimLeft, DelimRight string

	//true - load template where reload page (recommand when develop)
	//false - load template only one time when  app stating (recommand when deploy)
	HotReloadView bool
	JsVersion     string
}
