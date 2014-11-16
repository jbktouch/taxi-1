package taxi

import "testing"
import "strings"
import "os/exec"

func VerifyCommand(t *testing.T, cmd *exec.Cmd, docker Docker, params []string) {
	path := docker.Executable
	if cmd.Path != path {
		t.Errorf("Command path was %s but expected %s", cmd.Path, path)
	}

	expected_args := []string{path}
	expected_args = append(expected_args, params...)

	if strings.Join(cmd.Args, " ") != strings.Join(expected_args, " ") {
		t.Errorf("Command args were %q but expected %q", cmd.Args, expected_args)
	}
}

func ExampleDocker() Docker {
	return Docker{"/path/to/docker", "tcp://foobar:123", true}
}

func ExampleContainer() DockerContainer {
	return ExampleDocker().Container("my-name", "/here")
}

func TestRunEnvironment(t *testing.T) {
	d := ExampleDocker()

	has_host := false
	has_tls := false

	for _, v := range d.RunEnvironment() {
		if v == "DOCKER_HOST=tcp://foobar:123" {
			has_host = true
		}

		if v == "DOCKER_TLS_VERIFY=1" {
			has_tls = true
		}
	}

	if !has_host {
		t.Errorf("DOCKER_HOST=tcp://foobar:123 not in environment")
	}
	if !has_tls {
		t.Errorf("DOCKER_TLS_VERIFY=1 not in environment")
	}
}

func TestCommand(t *testing.T) {
	d := ExampleDocker()

	intended_args := []string{"ps", "-a"}

	cmd := d.Command(intended_args...)
	VerifyCommand(t, cmd, d, intended_args)
}

func TestContainer(t *testing.T) {
	d := ExampleDocker()

	name := "example-container"
	path := "/here/it/is"
	container := d.Container(name, path)

	if container.Name != name {
		t.Errorf("Expected container.Name to be %s, but was %s", name, container.Name)
	}

	if container.Path != path {
		t.Errorf("Expected container.Path to be %s, but was %s", path, container.Path)
	}
}

func TestBuildCommand(t *testing.T) {
	d := ExampleDocker()
	c := d.Container("example", "/path/to/it")

	cmd := c.BuildCommand()

	expected_args := []string{"build", "-t", "example", "/path/to/it"}

	VerifyCommand(t, cmd, c.DockerHost, expected_args)
}

func TestRunCommand(t *testing.T) {
	c := ExampleContainer()

	expected_args := []string{"run", "--name=my-name-run", "--rm", c.Name,
		"id", "-p"}
	cmd := c.RunCommand("id", "-p")

	VerifyCommand(t, cmd, c.DockerHost, expected_args)
}

func TestRunScriptCommand(t *testing.T) {
	c := ExampleContainer()

	script := "pip install flake8; flake8 /test"

	expected_args := []string{"run", "--name=my-name-run", "--rm", c.Name,
		"/bin/bash", "-c", script}
	cmd := c.RunScriptCommand(script)

	VerifyCommand(t, cmd, c.DockerHost, expected_args)
}

func TestDestroyCommand(t *testing.T) {
	c := ExampleContainer()

	expected_args := []string{"rmi", "--force", c.Name}
	cmd := c.DestroyCommand()

	VerifyCommand(t, cmd, c.DockerHost, expected_args)
}

func TestUntagCommand(t *testing.T) {
	c := ExampleContainer()

	expected_args := []string{"rmi", "--no-prune", "--force", c.Name}
	cmd := c.UntagCommand()

	VerifyCommand(t, cmd, c.DockerHost, expected_args)
}
