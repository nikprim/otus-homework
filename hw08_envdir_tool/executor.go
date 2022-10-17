package main

import (
	"log"
	"os"
	"os/exec"
)

func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Env = append(os.Environ(), env.toStrings()...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Start(); err != nil {
		log.Print(err)

		return 1
	}

	_ = command.Wait()

	return command.ProcessState.ExitCode()

}
