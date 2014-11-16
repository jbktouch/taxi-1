package taxi

import "io"
import "os"
import "net/http"
import "path/filepath"

func LatestDockerSource() io.ReadCloser {
	src := "https://get.docker.io/builds/Linux/x86_64/docker-latest"
	resp, err := http.Get(src)
	if err != nil {
		panic(err)
	}

	return resp.Body
}

func LocalDockerTarget(target string) io.WriteCloser {
	out, err := os.Create(target)
	if err != nil {
		panic(err)
	}

	return out
}

func PerformCopy(input io.ReadCloser, output io.WriteCloser) {
	defer input.Close()
	defer output.Close()

	_, err := io.Copy(output, input)
	if err != nil {
		panic(err)
	}
}

func (context Context) DockerCli() string {
	return filepath.Join(context.CacheDirectory, "docker131")
}

func (context Context) DownloadDocker() {
	filename := context.DockerCli()

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		PerformCopy(LatestDockerSource(), LocalDockerTarget(filename))
	}

	check(os.Chmod(filename, 0500))
}
