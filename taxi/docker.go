package taxi

import "fmt"
import "os/exec"

type Docker struct {
	Executable string
	host       string
	verify_tls bool
}

type DockerContainer struct {
	DockerHost Docker
	Name       string
	Path       string
}

func DockerHost(executable string, host string, verify_tls bool) Docker {
	return Docker{executable, host, verify_tls}
}

func (d Docker) RunEnvironment() []string {
	env := []string{fmt.Sprintf("DOCKER_HOST=%s", d.host)}

	if d.verify_tls {
		env = append(env, "DOCKER_TLS_VERIFY=1")
	}

	return env
}

func (d Docker) Command(args ...string) *exec.Cmd {
	cmd := exec.Command(d.Executable, args...)
	cmd.Env = d.RunEnvironment()

	return cmd
}

func (d Docker) Container(name string, path string) DockerContainer {
	return DockerContainer{d, name, path}
}

func (c DockerContainer) BuildCommand() *exec.Cmd {
	return c.DockerHost.Command("build", "-t", c.Name, c.Path)
}

func (c DockerContainer) RunCommand(command ...string) *exec.Cmd {
	arguments := []string{"run", "-i", fmt.Sprintf("--name=%s-run", c.Name),
		"--rm", c.Name}
	arguments = append(arguments, command...)
	cmd := c.DockerHost.Command(arguments...)

	return cmd
}

func (c DockerContainer) RunScriptCommand(command string) *exec.Cmd {
	return c.RunCommand("/bin/bash", "-c", command)
}

func (c DockerContainer) DestroyCommand() *exec.Cmd {
	return c.DockerHost.Command("rmi", "--force", c.Name)
}

func (c DockerContainer) TagCommand(tagname string) *exec.Cmd {
	return c.DockerHost.Command("tag", c.Name, tagname)
}

func (c DockerContainer) PushCommand(tagname string) *exec.Cmd {
	return c.DockerHost.Command("push", tagname)
}

func (c DockerContainer) UntagCommand() *exec.Cmd {
	return c.DockerHost.Command("rmi", "--no-prune", "--force", c.Name)
}

func (context Context) Docker() Docker {
	return DockerHost(context.DockerCli(), context.DockerHost, true)
}

func (context Context) Container() DockerContainer {
	return context.Docker().Container(context.ProjectName, context.RootPath)
}
