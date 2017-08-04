package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func genEtags(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		files, _ := ioutil.ReadDir(path)
		etags := make(map[string]string)
		for _, file := range files {
			if file.IsDir() {
				if er := genEtags(filepath.Join(path, file.Name())); er != nil {
					return er
				}
			} else if file.Name() == "etags.json" {
				continue
			} else {
				etags[file.Name()] = tTag(file.ModTime())
			}
		}
		etf, err := os.Create(filepath.Join(path, "etags.json"))
		if err != nil {
			return err
		}
		enc := json.NewEncoder(etf)
		enc.SetIndent("", " ")
		enc.Encode(etags)
		etf.Close()
		return nil
	}
	return errors.New("Not a directory")
}

func tTag(t time.Time) string {
	buf := make([]byte, 3, 3)
	buf[0] = 64 + byte(t.Day())
	buf[1] = 97 + byte(t.Hour())
	buf[2] = 65 + byte(t.Minute())
	return string(buf)
}
