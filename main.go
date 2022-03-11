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
	fmt.Println("----runWithStdIO----")
	if err := runWithStdIO(); err != nil {
		return err
	}

	fmt.Println("----runWithCommands----")
	if err := runWithCommands(); err != nil {
		return err
	}
	return nil
}

func runWithStdIO() error {
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

func runWithCommands() error {
	cmds := []*exec.Cmd{
		exec.Command("ls", "-1", "-a"),
		exec.Command("grep", "-v", "-E", "^[.].*"),
	}

	if err := build(cmds); err != nil {
		return err
	}

	if err := start(cmds); err != nil {
		return err
	}

	if err := wait(cmds); err != nil {
		return err
	}

	return nil
}

func build(cmds []*exec.Cmd) error {
	cmds[0].Stdin = os.Stdin
	cmds[len(cmds)-1].Stdout = os.Stdout

	for i := 0; i < len(cmds)-1; i++ {
		current := cmds[i]
		next := cmds[i+1]

		out, err := current.StdoutPipe()
		if err != nil {
			return err
		}

		next.Stdin = out
	}

	return nil
}

func start(cmds []*exec.Cmd) error {
	for _, c := range cmds {
		if err := c.Start(); err != nil {
			return err
		}
	}

	return nil
}

func wait(cmds []*exec.Cmd) error {
	for _, c := range cmds {
		if err := c.Wait(); err != nil {
			return err
		}
	}

	return nil
}
