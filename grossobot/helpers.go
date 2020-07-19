package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func containsVal(s []string, e string) int {
	for i, a := range s {
		if a == e {
			return i
		}
	}
	return -1
}

func archiveJSON(fn string, ty interface{}) {
	f, err := os.Create(fn)
	if err != nil {
		return
	}

	defer f.Close()

	arch, err := json.Marshal(ty)
	if err != nil {
		return
	}
	c, err := f.Write(arch)
	if err != nil {
		return
	}
	fmt.Println("bytes: ", c)
}

func unarchiveJSON(fn string, ty interface{}) {
	if fileExists(fn) {
		dat, err := ioutil.ReadFile(fn)
		if err != nil {
			return
		}
		json.Unmarshal(dat, ty)
	}
}

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}
