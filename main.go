package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func main() {
	if err := Main(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Main() error {
	if err := useStdIO(); err != nil {
		return err
	}
	return nil
}

func useStdIO() error {
	cmd := exec.Command("/bin/bash", "-c", "tr a-z A-Z | sort; echo ...done...")

	wc, _ := cmd.StdinPipe()
	rc, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		return err
	}

	go stdIn(wc)
	go stdOut(rc)

	return cmd.Wait()
}

func stdIn(wc io.WriteCloser) {
	defer wc.Close()

	io.WriteString(wc, "python\n")
	io.WriteString(wc, "csharp\n")
	io.WriteString(wc, "golang\n")
	io.WriteString(wc, "java\n")
}

func stdOut(rc io.ReadCloser) {
	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
