package main

import (
	"config"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//var ch = make(chan error)

func main() {

	var isZip string
	fmt.Printf("\n********************************************************\n\n")
	fmt.Printf("      目录设置如有变更请在程序根目录下修改.conf文件               \n")
	fmt.Printf("                  %s for  LFX                           \n", "V1.0")
	fmt.Printf("                                                          \n")
	fmt.Printf("********************************************************\n\n")

	fmt.Printf("\n\n-->获取配置文件信息......\n\n")
	fmt.Printf("-->获取成功、按回车开始收集备份")
	fmt.Scanln(&isZip)
	fmt.Printf("\n\n")
	//获取文件路径。
	myconfig := new(conf.Config)
	myconfig.InitConfig("./BackupPath.conf")
	//获取源文件所在目录集合
	allpath := myconfig.Read("default", "path")
	pathArray := strings.Split(allpath, ";")
	//获取处理完后存放的位置。
	toPathDirName := myconfig.Read("default", "toPathDirName")
	toDir := myconfig.Read("default", "topath")
	toDir, err4 := CreateDir(toDir, toPathDirName)
	if err4 != nil {
		fmt.Println("-->目标目录获取失败,或者目录已存在请删除他。  ！按任意键退出程序！ err=", err4)
		fmt.Scanln(&isZip)
		return
	}
	fmt.Println("-->目标目录是  ", toDir)
	//——————————————————————————————————————————————————————————————————————————
	for index := 0; index < len(pathArray); index++ {
		//获取文件完整路径
		path, err3 := GetFullPath(pathArray[index])
		if err3 != nil {
			fmt.Println("-->获取源文件失败 err=", err3)
			return
		}
		fmt.Printf("-->正在复制第%d个文件 ，共%d个文件\n", index+1, len(pathArray))
		CopyFile(toDir, path)
		//err0 := <-ch
		//if err0 != nil {
		fmt.Printf("-->第%d个文件完成复制\n", index+1)
		//}
	}

	fmt.Printf("\n\n-->是否要需要压缩打包？ 请选择 yes/no\n")
	fmt.Print("-->:")
	//var isZip string
	fmt.Scanln(&isZip)
	fmt.Println("-->这个功能还没做 。。。input=", isZip)
	fmt.Scanln(&isZip)

}

func GetFullPath(path string) (string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {

		return "", err
	}
	var fileName string
	var fileTime time.Time
	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			if fileTime.IsZero() {
				fileTime = file.ModTime()
				fileName = file.Name()
			} else {
				if fileTime.Before(file.ModTime()) {
					fileTime = file.ModTime()
					fileName = file.Name()
				}
			}

		}
	}
	path += "\\"
	path += fileName
	return path, nil
}

func CreateDir(path, name string) (string, error) {
	path += "\\"
	timeStr := strings.Replace(time.Now().String(), "-", "", -1)

	dirName := []byte(timeStr)

	if strings.Contains(name, "timeNow") {
		path += "TD"
		path += string(dirName[:8])
	} else {
		path += name
	}

	err := os.Mkdir(path, os.ModePerm)
	if err != nil {
		return "", err
	}
	return path, nil

}

func CopyFile(desPath, srcPath string) {
	//获取文件名
	_, fileName := filepath.Split(srcPath)
	desPath += "\\"
	desPath += fileName
	//打开源文件
	srcFile, err1 := os.Open(srcPath)
	defer srcFile.Close()
	if err1 != nil {
		fmt.Println("os.Open(srcPath) err=", err1)
		return
	}

	//创建目标文件。
	desFile, err2 := os.Create(desPath)
	defer desFile.Close()
	if err2 != nil {
		fmt.Println("os.Create(desPath) err=", err2)
		return
	}

	//开始复制文件
	buf := make([]byte, 4*1024)
	for {
		n, err3 := srcFile.Read(buf)
		if err3 != nil {
			if err3 == io.EOF {
				//ch <- err3
				return
			}
			return
			fmt.Println("srcFile.Read(buf) err=", err3)
		}
		desFile.Write(buf[:n])
	}
	return

}
