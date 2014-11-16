package taxi

import "io/ioutil"
import "fmt"
import "os"
import "os/exec"
import "path/filepath"

type DockerCerts struct {
	Pass     string
	PassFile string

	KeyFileSource       string
	CertFileSource      string
	AuthorityFileSource string

	KeyFileTarget       string
	CertFileTarget      string
	AuthorityFileTarget string
}

func (context Context) DockerCert() DockerCerts {
	return DockerCert(context.DockerPass, context.RootPath, context.CertificatePath)
}

func DockerCert(pass string, source_dir string, target_dir string) DockerCerts {
	cert := DockerCerts{
		pass,
		filepath.Join(target_dir, "key.pass"),
		filepath.Join(source_dir, ".docker.pem"),
		filepath.Join(source_dir, ".docker.pem"),
		filepath.Join(source_dir, ".docker.ca"),
		filepath.Join(target_dir, "key.pem"),
		filepath.Join(target_dir, "cert.pem"),
		filepath.Join(target_dir, "ca.pem"),
	}

	return cert
}

func (certs DockerCerts) DecryptKey() {
	defer os.Remove(certs.PassFile)
	check(ioutil.WriteFile(certs.PassFile, []byte(certs.Pass), 0644))

	passfile := fmt.Sprintf("file:%s", certs.PassFile)
	RunCommand(exec.Command("openssl", "rsa", "-in", certs.KeyFileSource,
		"-passin", passfile, "-out", certs.KeyFileTarget))
}

func (certs DockerCerts) MoveCertificates() {
	check(os.Rename(certs.CertFileSource, certs.CertFileTarget))
	check(os.Rename(certs.AuthorityFileSource, certs.AuthorityFileTarget))
}
