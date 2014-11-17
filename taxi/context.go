package taxi

import "fmt"

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

func (context Context) Describe() {
	var passoverten string
	if len(context.DockerPass) > 10 {
		passoverten = "Yes"
	} else {
		passoverten = "No"
	}

	fmt.Printf("  Context:\n")
	fmt.Printf("     Root:\t%s\n", context.RootPath)
	fmt.Printf("     Name:\t%s\n", context.ProjectName)
	fmt.Printf("    Cache:\t%s\n", context.CacheDirectory)
	fmt.Printf("    Certs:\t%s\n", context.CertificatePath)
	fmt.Printf("     Host:\t%s\n", context.DockerHost)
	fmt.Printf("  Pass>10:\t%s\n", passoverten)
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
