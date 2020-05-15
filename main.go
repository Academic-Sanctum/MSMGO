package main

import (
	"MSMGO/backupmanagement"
	"MSMGO/msm"
	"fmt"
	"io"
	"os"
	"sync"
)

func main() {
	var strout string
	var stdin io.WriteCloser
	var sw sync.WaitGroup
	var restart sync.WaitGroup
	var issave sync.WaitGroup
	var player string
	var message string
	go func() {
		err := msm.Run(&sw, &strout, &stdin, "./start.sh")
		if err != nil {
			fmt.Println("run", err)
			os.Exit(1)
		}
	}()
	go func() {
		for {
			sw.Add(1)
			sw.Wait()
			if len(strout) > 15 {
				fmt.Println((strout)[11:47])
				if (strout)[11:47] == "[Server thread/INFO]: Saved the game" {
					fmt.Println("true")
					issave.Done()
				}
			}
			go func() {
				var err error
				player, message, err = msm.FormatStringPlayerMessage(strout)
				if err != nil {
					return
				}
				go backupmanagement.Backup(&sw, &restart, &issave, player, message, &stdin, &strout)
			}()

		}
	}()

	restart.Add(1)
	restart.Wait()
}
