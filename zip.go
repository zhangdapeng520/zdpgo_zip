package zdpgo_zip

import (
	"archive/zip"
	"github.com/zhangdapeng520/zdpgo_log"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

/*
@Time : 2022/5/30 14:13
@Author : 张大鹏
@File : zip.go
@Software: Goland2021.3.1
@Description:
*/

type Zip struct {
	Config *Config
	Log    *zdpgo_log.Log
}

func New() *Zip {
	return NewWithConfig(&Config{})
}

func NewWithConfig(config *Config) *Zip {
	z := &Zip{}

	// 日志
	if config.LogFilePath == "" {
		config.LogFilePath = "logs/zdpgo/zdpgo_zip.log"
	}
	z.Log = zdpgo_log.NewWithDebug(config.Debug, config.LogFilePath)

	// 配置
	z.Config = config

	// 返回
	return z
}

// Zip 压缩文件夹
func (z *Zip) Zip(dirPath, targetName string) error {

	// 预防：旧文件无法覆盖
	_ = os.RemoveAll(targetName)

	// 创建：zip文件
	zipFile, err := os.Create(targetName)
	if err != nil {
		z.Log.Error("创建压缩文件失败", "error", err, "targetName", targetName)
		return err
	}
	defer zipFile.Close()

	// 打开：zip文件
	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	// 遍历路径信息
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, _ error) error {

		// 如果是源路径，提前进行下一个遍历
		if path == dirPath {
			return nil
		}

		// 获取：文件头信息
		header, _ := zip.FileInfoHeader(info)
		header.Name = strings.TrimPrefix(path, dirPath+`/`)

		// 判断：文件是不是文件夹
		if info.IsDir() {
			header.Name += `/`
		} else {
			// 设置：zip的文件压缩算法
			header.Method = zip.Deflate
		}

		// 创建：压缩包头部信息
		writer, _ := archive.CreateHeader(header)
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				z.Log.Error("打开文件失败", "error", err, "path", path)
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
			if err != nil {
				z.Log.Error("复制文件失败", "error", err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		z.Log.Error("压缩文件失败", "error", err)
		return err
	}
	return nil
}

// ZipAndDelete 压缩并删除文件夹
func (z *Zip) ZipAndDelete(dirPath, targetName string) error {
	err := z.Zip(dirPath, targetName)
	if err != nil {
		z.Log.Error("压缩文件失败", "error", err, "dir", dirPath)
		return err
	}

	err = os.RemoveAll(dirPath)
	if err != nil {
		z.Log.Error("删除文件夹失败", "error", err, "dir", dirPath)
		return err
	}

	return nil
}

// Unzip 解压缩文件夹
func (z *Zip) Unzip(zipFileName, saveDir string) error {
	//打开并读取压缩文件中的内容
	fr, err := zip.OpenReader(zipFileName)
	if err != nil {
		z.Log.Error("打开压缩文件失败", "error", err, "zipFileName", zipFileName)
		return err
	}
	defer fr.Close()

	//r.reader.file 是一个集合，里面包括了压缩包里面的所有文件
	for _, file := range fr.Reader.File {
		//判断文件该目录文件是否为文件夹
		if file.FileInfo().IsDir() {
			err = os.MkdirAll(file.Name, 0644)
			if err != nil {
				z.Log.Error("创建文件夹失败", "error", err)
				return err
			}
			continue
		}

		// 为文件时，打开文件
		r, err := file.Open()

		//文件为空的时候，打印错误
		if err != nil {
			z.Log.Error("打开文件失败", "error", err)
			continue
		}

		// 创建保存文件夹
		dirPath, fileName := filepath.Split(file.Name)
		finalDir := path.Join(saveDir, dirPath)
		_ = os.MkdirAll(finalDir, os.ModePerm)

		// 创建文件
		finalName := path.Join(finalDir, fileName)
		newFile, err := os.Create(finalName)
		if err != nil {
			z.Log.Error("创建文件失败", "error", err)
			continue
		}

		// 将内容复制
		_, err = io.Copy(newFile, r)
		if err != nil {
			z.Log.Error("复制文件内容失败", "error", err)
			return err
		}
		newFile.Close()
		r.Close()
	}

	return nil
}

// UnzipToCurrentDir 解压到当前目录
func (z *Zip) UnzipToCurrentDir(zipFileName string) error {
	return z.Unzip(zipFileName, "./")
}

// UnzipToCurrentDirAndDelete 解压到当前目录并删除zip文件
func (z *Zip) UnzipToCurrentDirAndDelete(zipFileName string) error {
	err := z.Unzip(zipFileName, "./")
	if err != nil {
		z.Log.Error("解压文件失败", "error", err, "file", zipFileName)
		return err
	}

	err = os.RemoveAll(zipFileName)
	if err != nil {
		z.Log.Error("删除文件失败", "error", err, "file", zipFileName)
		return err
	}

	return nil
}
