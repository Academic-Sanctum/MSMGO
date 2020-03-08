package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func main() {
	err := run("./start.sh")
	if err != nil {
		fmt.Println("run", err)
	}
}

//临时的方案
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
		//var player string
		var Fline string
		site0 := strings.Index(line, "<")
		if site0 != -1 {
			site1 := strings.Index(line[site0:], ">")
			if site1 != -1 {
				Lsite0 := len(line)
				//player = line[site0+1 : site0+site1]
				//go stdin.Write([]byte("say " + player + "\r"))
				Fline = line[site0+site1+2 : Lsite0-1]

				//go stdin.Write([]byte("tell " + player + " " + Fline + "\r"))
				go stdin.Write([]byte("tellraw @a \"" + Fline + "\"\r"))
				//go stdin.Write([]byte("say " + Fline + "\r"))
				fmt.Println(Fline)
			}

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

func say(write io.WriteCloser, data string) {
	go write.Write([]byte("tellraw @a {\"text\":" + data + "}"))
}
