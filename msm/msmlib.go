package msm

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

func Run(sw *sync.WaitGroup, strOut *string, stdIn *io.WriteCloser, CommandName string) (err error) {
	cmd := exec.Command(CommandName)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stdin, err := cmd.StdinPipe()

	if err != nil {
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
		time.Sleep(time.Second)
		*strOut = line
		if errr != nil || io.EOF == errr {
			break
		}
		fmt.Println(line)
		sw.Done()
		// sw2.Add(1)
		// sw2.Wait()
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

// func IsSave(strout string, issave *sync.WaitGroup) {

// 	fmt.Println("flase")
// }

//发送stop指令
func Stop(write *io.WriteCloser) {
	(*write).Write([]byte("stop\n"))
	(*write).Write(nil)
}

//发送save指令
func Save(write *io.WriteCloser) {
	(*write).Write([]byte("save-all\n"))
	(*write).Write(nil)
}

//发送信息全体
func SayAll(write *io.WriteCloser, data string) {
	(*write).Write([]byte("tellraw @a \"" + data + "\"\n"))
	(*write).Write(nil)

}

//发送信息制定玩家
func SayPlayer(write *io.WriteCloser, player string, data string) {
	(*write).Write([]byte("tellraw " + player + " \"" + data + "\"\n"))
	(*write).Write(nil)
}

//格式化字符串;并返回玩家名和玩家发送的信息;否则返回错误
func FormatStringPlayerMessage(line string) (player string, message string, err error) {
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
func FormatStringMessage(line string) (message string, err error) {
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
