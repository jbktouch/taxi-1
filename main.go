package main

import "os"
import "fmt"
import "math/rand"
import "path/filepath"
import "strings"
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

func TaxiBuildRepo() string {
	name := os.Getenv("TAXI_REPO")

	if name == "" {
		name = os.Getenv("TRAVIS_REPO_SLUG")
	}

	if name == "" {
		name = "test_build"
	}

	return name
}

func TaxiBuildId() string {
	id := os.Getenv("TAXI_BUILD_ID")

	if id == "" {
		id = os.Getenv("TRAVIS_JOB_ID")
	}

	if id == "" {
		id = fmt.Sprintf("%d", rand.Intn(1000))
	}

	return id
}

func TaxiBuildName() string {
	return fmt.Sprintf("%s:%s", TaxiBuildRepo(), TaxiBuildId())
}

func TaxiContainerName() string {
	name := TaxiBuildName()
	name = strings.Replace(name, "/", ".", -1)
	name = strings.Replace(name, ":", ".", -1)
	name = strings.ToLower(name)

	return name
}

func main() {
	if TaxiBuildRepo() == "test_build" {
		println("Please define TAXI_REPO prior to running taxi.")
		println("Using the name 'test_build' for now.")
	}

	project := taxi.SetupContext(".", TaxiContainerName())
	project.CacheDirectory = MkAbsDir(".taxi-cache")
	project.CertificatePath = MkAbsDir(".docker")
	project.DockerHost = os.Getenv("DOCKER_HOST")
	project.DockerPass = os.Getenv("DOCKER_PASS")
	project.SkipFileTime = os.Getenv("TAXI_SKIP_FILETIME") == "1"

	if len(os.Args) == 1 {
		println("please pass install or test or cleanup...")
		println("prefferrably all of them, in seperate steps, and in that order.")
		println("Info about your install:")
		project.Describe()

		os.Exit(1)
	}

	cmd := os.Args[1]

	if cmd == "info" {
		project.Describe()
	}

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

	if os.Args[1] == "pushtag" {
		if len(os.Args) != 3 {
			println("taxi pushtag internalserver:5000/repo/name:tag\n")
			os.Exit(1)
		}

		project.PushAndTag(os.Args[2])
	}

	if os.Args[1] == "cleanup" {
		project.DestroyContainer()
	}
}
