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
	if err := run(); err != nil {
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

func run() error {
	var (
		cmds = []*exec.Cmd{
			exec.Command("ls", "-1", "-a", homeDir()),
			exec.Command("grep", "-v", "-E", "^[.].*"),
			exec.Command("wc", "-l"),
		}
		err error
	)

	if err = prepare(cmds); err != nil {
		return err
	}

	if err = build(cmds); err != nil {
		return err
	}

	if err = start(cmds); err != nil {
		return err
	}

	if err = wait(cmds); err != nil {
		return err
	}

	return nil
}

func homeDir() string {
	var (
		dir string
		err error
	)

	if dir, err = os.UserHomeDir(); err != nil {
		panic(err)
	}

	return dir
}

func prepare(cmds []*exec.Cmd) error {
	cmds[0].Stdin = os.Stdin
	cmds[len(cmds)-1].Stdout = os.Stdout

	return nil
}

func build(cmds []*exec.Cmd) error {
	for i := 0; i < len(cmds)-1; i++ {
		var (
			curr = cmds[i]
			next = cmds[i+1]
			out  io.ReadCloser
			err  error
		)

		if out, err = curr.StdoutPipe(); err != nil {
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
