package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

/*
=== Взаимодействие с ОС ===

Необходимо реализовать собственный шелл

встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах

Реализовать утилиту netcat (nc) клиент
принимать данные из stdin и отправлять в соединение (tcp/udp)
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

var (
	userHomeDir string
)

var (
	errorShellParse = errors.New("shell: parse error near `|'")
	errorUnknownCmd = errors.New("shell: command not found")
	errorChangeDir  = errors.New("cd: string not in pwd")
	errorIllegalPID = errors.New("kill: illegal pid")
)

type command struct {
	name string
	args []string
}

func init() {
	var err error
	userHomeDir, err = os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		run(scanner.Bytes())
	}
}

func validateAndBatch(cmdsBytes []byte) ([]command, error) {
	_cmds := strings.Split(string(cmdsBytes), " | ")

	availableCmds := map[string]struct{}{
		"cd":     {},
		"pwd":    {},
		"echo":   {},
		"kill":   {},
		"ps":     {},
		"\\quit": {},
	}

	cmds := make([]command, len(_cmds))
	for i, cmd := range _cmds {
		words := strings.Fields(cmd)
		if len(words) == 0 {
			return nil, errorShellParse
		}

		if _, ok := availableCmds[words[0]]; !ok {
			return nil, fmt.Errorf("%w: %s", errorUnknownCmd, words[0])
		}

		cmds[i] = command{
			name: words[0],
			args: words[1:],
		}
	}

	return cmds, nil
}

func run(cmdsBytes []byte) {
	cmds, err := validateAndBatch(cmdsBytes)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	for _, cmd := range cmds {
		switch cmd.name {
		case "cd":
			cmd.changeDirectory()
		case "pwd":
			cmd.presentWorkingDirectory()
		case "echo":
			cmd.echo()
		case "kill":
			cmd.kill()
		case "ps":
			cmd.processStatus()
		case "\\quit":
			os.Exit(0)
		}
	}
}

func getCurrDir() (string, error) {
	return os.Getwd()
}

func (c *command) changeDirectory() {
	if len(c.args) > 1 {
		fmt.Fprintf(os.Stderr, "%s: %s\n", errorChangeDir.Error(), c.args[0])
		return
	}

	path := userHomeDir
	if len(c.args) == 1 {
		if c.args[0][0] == '/' {
			path = c.args[1]
		} else {
			var err error
			path, err = getCurrDir()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", c.name, err.Error())
				return
			}
			path = filepath.Join(path, c.args[0])
		}
	}

	if err := os.Chdir(path); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", c.name, err.Error())
	}
}

func (c *command) presentWorkingDirectory() {
	path, err := getCurrDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", c.name, err.Error())
		return
	}
	fmt.Fprintln(os.Stdout, path)
}

func (c *command) echo() {
	fmt.Fprintln(os.Stdout, strings.Join(c.args, " "))
}

func (c *command) kill() {
	pids := make([]int, len(c.args))
	for i, arg := range c.args {
		pid, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", errorIllegalPID.Error(), arg)
			return
		}
		pids[i] = pid
	}

	for _, pid := range pids {
		proc, err := os.FindProcess(pid)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", c.name, err.Error())
		}
		if err = proc.Kill(); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", c.name, err.Error())
		}
	}
}

func (c *command) processStatus() {
	cmd := exec.Command(c.name)

	stdout, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", c.name, err.Error())
		return
	}

	fmt.Fprint(os.Stdout, string(stdout))
}
