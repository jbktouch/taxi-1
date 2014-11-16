package main

import "os"
import "path/filepath"

import "github.com/grahamc/taxi/taxi"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func MkAbsDir(path string) string {
	path, err := filepath.Abs(path)
	check(err)

	err = os.Mkdir(path, 0700)
	if !os.IsExist(err) {
		check(err)
	}

	return path
}

func main() {

	project := taxi.SetupContext(".", "test")
	project.CacheDirectory = MkAbsDir(".taxi-cache")
	project.CertificatePath = MkAbsDir(".docker")
	project.DockerHost = os.Getenv("DOCKER_HOST")
	project.DockerPass = os.Getenv("DOCKER_PASS")
	project.SkipFileTime = true

	if len(os.Args) == 1 {
		print("please pass install or test or cleanup...\n")
		print("prefferrably all of them, in seperate steps, and in that order.\n")
		print("or pass all and it'll do all of them.\n")

		os.Exit(1)
	}

	cmd := os.Args[1]

	if cmd == "install" || cmd == "all" {
		project.Install()
	}

	if len(os.Args) > 2 && (os.Args[1] == "test" || cmd == "all") {
		project.TestContainer(os.Args[2])
	}

	if os.Args[1] == "cleanup" || cmd == "all" {
		project.DestroyContainer()
	}
}
