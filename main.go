package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func main() {
	err := run("./start.sh")
	if err != nil {
		fmt.Println("run", err)
	}
}

func run(CommandName string) error {
	cmd := exec.Command(CommandName)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("cmd.StdoutPipe: ", err)
		return err
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("cmd.StdinPipe: ", err)
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	read := bufio.NewReader(stdout)
	for {

		line, errr := read.ReadString('\n')
		if errr != nil || io.EOF == errr {
			break
		}
		fmt.Println(line)
		go func() {
			write := bufio.NewReader(os.Stdin)
			writestr, err := write.ReadString('\n')
			if err != nil {
				fmt.Println("write.ReadString: ", err)
			}
			stdin.Write([]byte(writestr))
		}()

	}
	err = cmd.Wait()
	return err
}
