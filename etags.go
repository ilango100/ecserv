package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"time"
)

func genDeps(pth string) (map[string][]string,error) {
	depfilename := path.Join(pth, "deps.json")
	if inf, err := os.Stat(depfilename); err != nil && !inf.IsDir() {
		return nil,err
	}
	deps := make(map[string][]string)
	depfile, err := os.Open(depfilename)
	if err != nil {
		return nil,err
	}
	dec := json.NewDecoder(depfile)
	if err := dec.Decode(&deps); err != nil {
		return nil,err
	}
	return deps,nil
}

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
