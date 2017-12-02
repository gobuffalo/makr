package makr

import (
	"os/exec"

	"github.com/gobuffalo/envy"
)

// GoInstall compiles and installs packages and dependencies
func GoInstall(pkg string, opts ...string) *exec.Cmd {
	args := append([]string{"install"}, opts...)
	args = append(args, pkg)
	return exec.Command(envy.Get("GO_BIN", "go"), args...)
}

// GoGet downloads and installs packages and dependencies
func GoGet(pkg string, opts ...string) *exec.Cmd {
	args := append([]string{"get"}, opts...)
	args = append(args, pkg)
	return exec.Command(envy.Get("GO_BIN", "go"), args...)
}

// GoFmt is command that will use `goimports` if available,
// or fail back to `gofmt` otherwise.
func GoFmt(files ...string) *exec.Cmd {
	if len(files) == 0 {
		files = []string{"."}
	}
	c := "goimports"
	_, err := exec.LookPath(c)
	if err != nil {
		c = "gofmt"
		_, err = exec.LookPath(c)
		if err != nil {
			return []string{"could not find gofmt or goimports"}
		}
	}
	args := append([]string{"-w"}, files...)
	return exec.Command(c, args...)
}
