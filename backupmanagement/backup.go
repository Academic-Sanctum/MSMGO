package backupmanagement

import (
	"MSMGO/msm"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"time"
)

type List struct {
	Name string
	Time string
}

func Backup(sw *sync.WaitGroup, restart *sync.WaitGroup, issave *sync.WaitGroup, player string, message string, stdin *io.WriteCloser, strout *string) {

	if len(message) >= 4 {
		if message[0:4] == ">>bm" {
			if message == ">>bm" {
				msm.SayPlayer(stdin, player, "[BackUp] >>bm save 注释   备份")
				msm.SayPlayer(stdin, player, "[BackUp] >>bm back 序号   恢复")
				msm.SayPlayer(stdin, player, "[BackUp] >>bm list        列表")
			} else if len(message) >= 10 {

				if message[0:10] == ">>bm save " {

					msm.SayAll(stdin, "[BackUp] 备份开始")
					timeStart := time.Now()
					issave.Add(1)
					msm.Save(stdin)
					issave.Wait()
					timeS := getTime()

					go func() {
						jsonPath := "./backupmanagement/backup/backupList.json"
						backupPath := "./backupmanagement/backup"
						pExists, err := pathExists(backupPath)
						if err != nil {
							fmt.Println("获取目录错误", err)
							return
						}
						if !pExists {
							err = os.Mkdir(backupPath, os.ModePerm)
							if err != nil {
								return
							}
						}
						jpExists, err := pathExists(jsonPath)
						if err != nil {
							return
						}
						if jpExists {
							jsonMap, err := readJSON(jsonPath)
							if err != nil {
								return
							}
							jsonMapLen := len(jsonMap)
							jsonMap[jsonMapLen] = &List{message[10:], timeS}
							writeJSON(jsonMap, jsonPath)
						} else {
							f, err := os.Create(jsonPath)
							if err != nil {
								fmt.Println(err.Error())
							}
							f.Close()
							fmt.Println("创建完成")
							var jsonMap = make(map[int]*List)
							jsonMap[0] = &List{message[10:], timeS}
							writeJSON(jsonMap, jsonPath)
						}
					}()

					err := bmMake(message[10:])
					if err != nil {
						fmt.Println("创建失败", err)
						msm.SayAll(stdin, "[BackUp] 备份创建失败"+err.Error())
						return
					}
					timeStop := time.Now()
					timeUse := timeStop.Sub(timeStart)
					msm.SayAll(stdin, "[BackUp] 备份成功")
					msm.SayAll(stdin, "[BackUp] 日期: "+timeS)
					msm.SayAll(stdin, "[BackUp] 耗时:"+timeUse.String())
					msm.SayAll(stdin, "[BackUp] 注释:"+message[10:])
				} else if message[0:10] == ">>bm back " {
					jsonPath := "./backupmanagement/backup/backupList.json"
					pExists, err := pathExists(jsonPath)
					if err != nil {
						return
					}
					if !pExists {
						return
					}
					jsonMap, err := readJSON(jsonPath)
					if err != nil {
						return
					}
					id, err := strconv.Atoi(message[10:])
					if err != nil {
						return
					}
					msm.SayAll(stdin, "[BackUp] 准备还原备份 "+"序号:"+strconv.Itoa(id)+jsonMap[id].Name)
					for i := 5; i > 0; i-- {
						msm.SayAll(stdin, "[BackUp] "+strconv.Itoa(i)+"秒")
						time.Sleep(time.Second)
					}

					restart.Add(1)
					msm.Stop(stdin)
					time.Sleep(time.Second * 5)
					os.RemoveAll("./server/world")
					DeCompress("./backupmanagement/backup/"+jsonMap[id].Name+".zip", "./server")
					go func() {
						err := msm.Run(sw, strout, stdin, "./start.sh")
						if err != nil {
							fmt.Println("run", err)
							os.Exit(1)
						}
					}()
					restart.Done()

				}
			} else if message == ">>bm list" {
				jsonPath := "./backupmanagement/backup/backupList.json"
				pExists, err := pathExists(jsonPath)
				if err != nil {
					return
				}
				if !pExists {
					return
				}
				jsonMap, err := readJSON(jsonPath)
				if err != nil {
					return
				}
				jsonMapLen := len(jsonMap)
				for i := 0; i < jsonMapLen; i++ {
					msm.SayPlayer(stdin, player, "[BackUp] 序号"+strconv.Itoa(i)+"    名字："+jsonMap[i].Name+"    日期"+jsonMap[i].Time)
				}
			} else {
				msm.SayPlayer(stdin, player, "[BackUp] bm命令不正确")
			}
		}
	}

}

func bmMake(backupName string) (err error) {
	backupPath := "./backupmanagement/backup"
	pExists, err := pathExists(backupPath)
	if err != nil {
		fmt.Println("获取目录错误", err)
		return err
	}
	if !pExists {
		err = os.Mkdir(backupPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	file, err := os.Open("./server/world")
	if err != nil {
		return err
	}
	filess := file
	files := []*os.File{filess}
	dest := backupPath + "/" + backupName + ".zip"
	Compress(files, dest)
	return err
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func getTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func readFile(filePath string) (string, error) {
	byteRead, err := ioutil.ReadFile(filePath)
	str := string(byteRead)
	return str, err
}

func writeFile(filePath string, write string) error {
	writeByte := []byte(write)
	err := ioutil.WriteFile(filePath, writeByte, 0644)
	return err
}

func readJSON(jsonPath string) (map[int]*List, error) {
	jsonStr, err := readFile(jsonPath)

	var jsonMap = make(map[int]*List)

	err = json.Unmarshal([]byte(jsonStr), &jsonMap)
	if err != nil {
		fmt.Println("解码失败", err)
	} else {
		fmt.Println("解析成功", jsonMap)
	}
	return jsonMap, err
}

func writeJSON(jsonMap map[int]*List, jsonPath string) error {
	bytes, err := json.Marshal(jsonMap)
	if err != nil {
		fmt.Println("编码错误", err)
	} else {
		writeFile(jsonPath, string(bytes))
	}
	return err
}
