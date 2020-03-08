package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

func main() {
	var strout string
	var stdin io.WriteCloser
	var sw sync.WaitGroup

	go func() {
		for {
			sw.Add(1)
			sw.Wait()
			message, err := formatStringMessage(strout)
			if err != nil {
				//fmt.Println(err)
			}
			if message == ">>bm" {
				sayAll(stdin, ">>bm make + \\\"备份名\\\"    备份")
				sayAll(stdin, ">>bm load + \\\"备份名\\\"    加载")
				sayAll(stdin, ">>bm list                    显示备份列表")
			}
		}
	}()
	err := run(&sw, &strout, &stdin, "./start.sh")
	if err != nil {
		fmt.Println("run", err)
		os.Exit(1)
	}

}

//整合
func run(sw *sync.WaitGroup, strOut *string, stdIn *io.WriteCloser, CommandName string) (err error) {
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
	*stdIn = stdin

	err = cmd.Start()
	if err != nil {
		return err
	}

	//读取信息
	read := bufio.NewReader(stdout)

	for {
		line, errr := read.ReadString('\n')

		*strOut = line
		if errr != nil || io.EOF == errr {
			break
		}
		fmt.Println(line)
		sw.Done()
		runtime.Gosched()
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

//发送信息全体
func sayAll(write io.WriteCloser, data string) {
	write.Write([]byte("tellraw @a \"" + data + "\"\n"))
	write.Write(nil)

}

//发送信息制定玩家
func sayPlayer(write io.WriteCloser, player string, data string) {
	write.Write([]byte("tellraw " + player + " \"" + data + "\"\n"))
	write.Write(nil)
}

//格式化字符串;并返回玩家名和玩家发送的信息;否则返回错误
func formatStringPlayerMessage(line string) (player string, message string, err error) {
	site0 := strings.Index(line, "<")
	if site0 != -1 {
		site1 := strings.Index(line[site0:], ">")
		if site1 != -1 {
			Lsite0 := len(line)
			player = line[site0+1 : site0+site1]
			message = line[site0+site1+2 : Lsite0-1]
			err = nil
			return player, message, err
		}
		err = errors.New("not found")
		return player, message, err
	}
	err = errors.New("not found")
	return player, message, err
}

//格式化字符串;并返回玩家发送的信息;否则返回错误
func formatStringMessage(line string) (message string, err error) {
	site0 := strings.Index(line, "<")
	if site0 != -1 {
		site1 := strings.Index(line[site0:], ">")
		if site1 != -1 {
			Lsite0 := len(line)
			message = line[site0+site1+2 : Lsite0-1]
			err = nil
			return message, err
		}
		err = errors.New("not found")
		return message, err
	}
	err = errors.New("not found")
	return message, err
}
