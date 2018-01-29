package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	args := os.Args

	if len(args) != 2 || args[0] != "gitget" {
		fmt.Println(" illegal command line !!")
		return
	}

	url := args[1]
	if strings.Contains(url, "http://") || strings.Contains(url, "https://") {
		fmt.Println(" the url cannot start with http:// or https:// ..")
		return
	}

	downloadUrl := "https://" + url + "/archive/master.zip"

	res, err := http.Get(downloadUrl)
	if err != nil {
		panic(err)
	}

	exist, _ := pathExists(url)
	if !exist {
		if err := os.MkdirAll(url, 0777); err != nil {
			panic(err)
		}
	}

	f, err := os.Create(url + "/package.zip")
	if err != nil {
		panic(err)
	}
	io.Copy(f, res.Body)
	f.Close()

	if err := deCompress(url+"/package.zip", url); err != nil {
		panic(err)
	}

	if err := os.Remove(url + "/package.zip"); err != nil {
		panic(err)
	}
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

func deCompress(zipFile, dest string) error {

	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		filename := dest + filterRootDir(file.Name)
		err = os.MkdirAll(getDir(filename), 0777)
		if err != nil {
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			//return err
			w.Close()
			continue
		}
		defer w.Close()
		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
		w.Close()
		rc.Close()
	}
	return nil
}

func filterRootDir(name string) string {
	path := ""
	infos := strings.Split(name, "/")
	length := len(infos)
	if length <= 1 {
		return ""
	}
	for i := 1; i < length; i++ {
		path += "/"
		path += infos[i]
	}
	return path
}

func getDir(path string) string {
	return subString(path, 0, strings.LastIndex(path, "/"))
}

func subString(str string, start, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < start || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}
