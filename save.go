package main

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

func init() {
	// We register Entity types so that gob can encode them.
	gob.Register(&Player{})
	gob.Register(&Enemy{})
	gob.Register(&HealthPotion{})
	gob.Register(&MagicArrowScroll{})
	gob.Register(&ExplodeScroll{})
}

func Encode(g *Game) (encodedData []byte, err error) {
	data := bytes.Buffer{}
	enc := gob.NewEncoder(&data)
	err = enc.Encode(g)
	if err != nil {
		return 
	}
	var buf bytes.Buffer
	writeGzip := gzip.NewWriter(&buf)
	defer func () {
		writeGzip.Close()
	}()

	writeGzip.Write(data.Bytes())
	encodedData = buf.Bytes()
	return 
}

func EncodeNoGzip(g *Game) (data []byte, err error) {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err = enc.Encode(g)
	if err != nil {
		return 
	}
	data = buf.Bytes()
	return 
}

func DecodeNoGzip(data []byte) (g *Game, err error) {
	dec := gob.NewDecoder(bytes.NewReader(data))
	g = &Game{}
	err = dec.Decode(g)
	return
}

func Decode(data []byte) (g *Game, err error) {
	g = &Game{}
	buf := bytes.NewReader(data)
	readGzip, err := gzip.NewReader(buf)
	if err != nil {
		return
	}
	defer func ()  {
		readGzip.Close()
	}()

	dec := gob.NewDecoder(readGzip)
	err = dec.Decode(g)
	return 
}

// DataDir ... returns path to directory contains data file if there is not, make directory 
func DataDir () (path string, err error) {
	var xdg string 
	if runtime.GOOS == "windows" {
		xdg = os.Getenv("LOCALAPPDATA")
	} else { // linux, BSD, etc..
		xdg = os.Getenv("XDG_DATA_HOME")
	}

	if xdg == "" {
		xdg = filepath.Join(os.Getenv("HOME"), ".local", "share")
	}

	path = filepath.Join(xdg, "rt")
	_, err = os.Stat(path)
	if err != nil {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			err = fmt.Errorf("building data directory: %s\n", err.Error())
			return
		}
	}
	return 
}

func SaveFile(fileName string, data []byte) (err error) {
	dataDir, err := DataDir()
	if err != nil {
		return
	}
	tempFileName := filepath.Join(dataDir, "temp-" + fileName)
	tempFile, err := os.OpenFile(tempFileName, os.O_WRONLY | os.O_CREATE |os.O_TRUNC, 0644)
	if err != nil {
		return 
	}
	_, err = tempFile.Write(data)
	if err != nil {
		return 
	}

	// closing tempFile for renaming 
	if err = tempFile.Sync(); err != nil {
		return 
	}
	if err = tempFile.Close(); err != nil {
		return 
	}

	savedFileName := filepath.Join(dataDir, fileName)
	if err = os.Rename(tempFileName, savedFileName); err != nil {
		return 
	}
	return
} 

func LoadFile(fileName string) (data []byte, err error) {
	dataDir, err := DataDir()
	if err != nil {
		err = fmt.Errorf("could not read data directory: %s", err.Error())
		return 
	}
	filePath := filepath.Join(dataDir, fileName)
	if _, err = os.Stat(filePath); err != nil {
		err = fmt.Errorf("no such file: %s", err.Error())
		return 
	}
	data, err = ioutil.ReadFile(filePath)
	return 

}

func RemoveDataFile(fileName string) (err error) {
	dataDir, err := DataDir()
	if err != nil {
		err = fmt.Errorf("could not read data directory: %s", err.Error())
		return 
	}

	filePath := filepath.Join(dataDir, fileName)
	if _, err = os.Stat(filePath); err != nil {
		err = fmt.Errorf("no such file: %s", err.Error())
		return 
	}
	err = os.Remove(filePath)
	return
}