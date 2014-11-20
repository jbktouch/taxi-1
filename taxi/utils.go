package taxi

import "os/exec"
import "fmt"
import "log"
import "os"
import "github.com/elwinar/logwriter"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func RunCommand(cmd *exec.Cmd) {
	fmt.Printf("Executing %q in %q\n", cmd.Args, cmd.Env)

	cmd.Stdin = os.Stdin
	cmd.Stdout = logwriter.New(log.New(os.Stdout, "stdout: ", log.Lshortfile))
	cmd.Stderr = logwriter.New(log.New(os.Stdout, "stderr: ", log.Lshortfile))

	check(cmd.Run())
}
