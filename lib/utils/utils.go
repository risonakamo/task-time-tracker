package utils

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/rs/zerolog"
)

// set zerolog global logger default options
func ConfigureDefaultZeroLogger() {
    log.Logger=log.Output(zerolog.ConsoleWriter{
        Out:os.Stdout,
        TimeFormat: "2006/01/02 15:04:05",
    })
}

// give folder location of the exe that calls this func
func GetHereDirExe() string {
    var exePath string
    var e error
    exePath,e=os.Executable()

    if e!=nil {
        panic(e)
    }

    return filepath.Dir(exePath)
}

// try to open web url or file with default program.
// essentially runs program like it was double clicked
func OpenTargetWithDefaultProgram(url string) error {
    var cmd *exec.Cmd=exec.Command("cmd","/c","start",url)
    var e error=cmd.Run()

    if e!=nil {
        return e
    }

    return nil
}

// overwrite target json file with a new file
func WriteJson(filename string,data any) error {
	var wfile *os.File
	var e error
	wfile,e=os.Create(filename)

	if e!=nil {
		panic(e)
	}

	defer wfile.Close()

	var jsondata []byte
	jsondata,e=json.Marshal(data)

	if e!=nil {
		panic(e)
	}

	wfile.Write(jsondata)
	return nil
}

// read an deserialise json file
func ReadJson[DataT any](filename string) (DataT,error) {
    var data []byte
	var e error
	data,e=os.ReadFile(filename)

	if errors.Is(e,fs.ErrNotExist) {
		log.Info().Msgf("file not found: %s",filename)
		var def DataT
		return def,e
	}

	if e!=nil {
		var def DataT
		return def,e
	}

	var parsedData DataT
    json.Unmarshal(data,&parsedData)

	return parsedData,nil
}
