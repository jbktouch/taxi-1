package taxi

type Context struct {
	RootPath        string
	ProjectName     string
	CacheDirectory  string
	CertificatePath string
	DockerHost      string
	DockerPass      string
	SkipFileTime    bool
}

func SetupContext(path string, name string) Context {
	return Context{RootPath: path, ProjectName: name}
}

func (context Context) Install() {
	context.DownloadDocker()
	context.SecureDocker()
	if !context.SkipFileTime {
		context.ForceFileTimes()
	}
	context.BuildContainer()
}

func (context Context) SecureDocker() {
	certs := context.DockerCert()

	certs.DecryptKey()
	certs.MoveCertificates()
}

func (context Context) ForceFileTimes() {
	ForceFileTimes(context.RootPath)
}

func (context Context) BuildContainer() {
	RunCommand(context.Container().BuildCommand())
}

func (context Context) TestContainer(script string) {
	RunCommand(context.Container().RunScriptCommand(script))
}

func (context Context) DestroyContainer() {
	RunCommand(context.Container().UntagCommand())
}
