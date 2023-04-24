package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/google/renameio/v2"
)

func mainE() error {
	for _, filename := range []string{"go.mod", "go.sum"} {
		fmt.Printf("Parsing %s... ", filename)

		b, err := os.ReadFile(filename)
		if err != nil {
			return err
		}

		w, err := renameio.TempFile("", filename)
		if err != nil {
			return err
		}
		defer w.Cleanup()

		lines := strings.Split(string(b), "\n")
		found := 0
		for _, l := range lines {
			if strings.HasPrefix(l, "<<<<<<<") {
				found++
				continue
			}
			if strings.HasPrefix(l, ">>>>>>>") || strings.HasPrefix(l, "=======") {
				continue
			}
			fmt.Fprintf(w, "%s\n", l)
		}

		if err := w.CloseAtomicallyReplace(); err != nil {
			return fmt.Errorf("error while trying to overwrite %q: %w", filename, err)
		}
		fmt.Printf("%d conflicts found\n", found)
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
