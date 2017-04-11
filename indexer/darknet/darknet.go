package darknet

import (
	"os/exec"
	"io"
	"encoding/json"
	"path/filepath"
	"runtime"
)

type Process struct {
	cmd          *exec.Cmd
	stdoutReader io.ReadCloser
	stdInWriter  io.WriteCloser
	decoder      *json.Decoder
}

const Success = "success"

var DarknetHome = ""

func NewProcess() (*Process, error) {
	p := new(Process)
	var err error
	p.cmd = exec.Command(getDarknetBinary())
	p.cmd.Dir = DarknetHome
	p.stdoutReader, err = p.cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	p.stdInWriter, err = p.cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	err = p.cmd.Start()
	if err != nil {
		return nil, err
	}
	p.decoder = json.NewDecoder(p.stdoutReader)
	return p, nil
}

func (p*Process) Detect(path string) (*Result, error) {
	_, err := io.WriteString(p.stdInWriter, path + "\n")
	if err != nil {
		return nil, err
	}
	var result Result
	err = p.decoder.Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p*Process) Close() {
	p.stdInWriter.Close()
	p.stdoutReader.Close()
	if (p.cmd.Process != nil) {
		p.cmd.Process.Kill()
	}
	p.cmd.Process.Wait()
}

func getDarknetBinary() string {
	darknet := "darknet"
	if runtime.GOOS == "windows" {
		darknet += ".exe"
	}
	return filepath.Join(DarknetHome, darknet)

}