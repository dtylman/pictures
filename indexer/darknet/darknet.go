package darknet

import (
	"encoding/json"
	"io"
	"log"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
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
	log.Printf("Started %v with PID %v", p.cmd.Path, p.cmd.Process.Pid)
	p.decoder = json.NewDecoder(p.stdoutReader)
	return p, nil
}

func (p *Process) Detect(path string, timeout time.Duration) (*Result, error) {
	log.Printf("Darknet detecting '%v'...", path)
	_, err := p.stdInWriter.Write([]byte(path + "\n"))
	if err != nil {
		return nil, err
	}
	var result Result
	var done = false
	defer func() {
		done = true
	}()
	go func() {
		time.Sleep(timeout)
		log.Printf("Timeout reached while waiting for object detection on %v", path)
		if !done {
			p.Close()
		}

	}()
	err = p.decoder.Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *Process) Close() {
	p.stdInWriter.Close()
	p.stdoutReader.Close()
	if p.cmd.Process != nil {
		p.cmd.Process.Kill()
	}
	log.Printf("Closing process %v", p.cmd.Process.Pid)
	_, err := p.cmd.Process.Wait()
	if err != nil {
		log.Println(err)
	}
}

func getDarknetBinary() string {
	darknet := "darknet"
	if runtime.GOOS == "windows" {
		darknet += ".exe"
	}
	return filepath.Join(DarknetHome, darknet)

}
