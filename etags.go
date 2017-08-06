package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"time"
)

func genEtags(pth string) (map[string]string, error) {
	info, err := os.Stat(pth)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		files, _ := ioutil.ReadDir(pth)
		etags := make(map[string]string)
		for _, file := range files {
			if file.IsDir() {
				var et map[string]string
				var er error
				if et, er = genEtags(path.Join(pth, file.Name())); er != nil {
					return nil, er
				}
				for n, e := range et {
					etags[path.Join(file.Name(), n)] = e
				}
			} else {
				etags[file.Name()] = tTag(file.ModTime())
			}
		}
		return etags, nil
	}
	return nil, errors.New("Not a directory")
}

func tTag(t time.Time) string {
	buf := make([]byte, 3, 3)
	buf[0] = 64 + byte(t.Day())
	buf[1] = 97 + byte(t.Hour())
	buf[2] = 65 + byte(t.Minute())
	return string(buf)
}
