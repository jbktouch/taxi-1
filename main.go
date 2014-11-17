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
	project.SkipFileTime = os.Getenv("TAXI_SKIP_FILETIME") == "1"

	if len(os.Args) == 1 {
		print("please pass install or test or cleanup...\n")
		print("prefferrably all of them, in seperate steps, and in that order.\n")

		print("Info about your install:\n")
		project.Describe()

		os.Exit(1)
	}

	cmd := os.Args[1]

	if cmd == "install" {
		project.Install()
	}

	if os.Args[1] == "test" {
		if len(os.Args) != 3 {
			print("taxi test 'your test script here'\n")
			os.Exit(1)
		}

		project.TestContainer(os.Args[2])
	}

	if os.Args[1] == "cleanup" {
		project.DestroyContainer()
	}
}
