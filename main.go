package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/google/renameio"
)

func mainE() error {
	goModFilename := "go.mod"

	b, err := ioutil.ReadFile(goModFilename)
	if err != nil {
		return err
	}

	w, err := renameio.TempFile("", goModFilename)
	if err != nil {
		return err
	}
	defer w.Cleanup()

	lines := strings.Split(string(b), "\n")
	for _, l := range lines {
		if strings.HasPrefix(l, "<<<<<<<") || strings.HasPrefix(l, ">>>>>>>") || strings.HasPrefix(l, "=======") {
			continue
		}
		fmt.Fprintf(w, "%s\n", l)
	}

	if err := w.CloseAtomicallyReplace(); err != nil {
		return fmt.Errorf("error while trying to overwrite %q: %w", goModFilename, err)
	}

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func main() {
	if err := mainE(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
