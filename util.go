package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var Debug = true

const (
	OR = 1
	HQ = 2
	FA = 3
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func GetSite() int {
	h, _ := os.Hostname()
	if Debug {
		log.Printf("finding site for hostname: %s\n", h)
	}
	if strings.Contains(strings.ToLower(h), "oriole") {
		return OR
	} else if strings.Contains(strings.ToLower(h), "hqts") {
		return HQ
	} else if strings.Contains(strings.ToLower(h), "falcon") {
		return FA
	} else {
		return -1
	}
}

func IsRootUser() bool {
	return os.Geteuid() == 0
}

func ExecuteCmdRaw(cmd string) string {
	out, err := exec.Command(cmd).Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

func ListFiles(path string, ext string) map[string]os.FileInfo {
	ret := make(map[string]os.FileInfo)
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	fmt.Println(path, ext)
	files, err := filepath.Glob(path + ext)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(files)
	for _, file := range files {
		fi := FileInfo(file)
		fmt.Println(fi.Name(), fi.Size(), fi.ModTime())
		ret[fi.Name()] = fi
	}

	return ret
}

func FileInfo(f string) os.FileInfo {
	fi, e := os.Stat(f)
	if e != nil {
		log.Println(e)
		return nil
	}
	return fi
}

func DeleteFile(f string) int {
	err := os.Remove(f)
	if err != nil {
		log.Println(err)
		return 1
	}
	return 0 // success
}
