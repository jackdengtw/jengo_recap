package action

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/golang/glog"
	_ "github.com/qetuantuan/jengo_recap/context"
)

const (
	SHEBANG = "#!/bin/bash"
	SET_E   = "set -e"
)

// Bash action is for *nix only
type Bash struct {
	Noop

	ScriptContent string
	Env           []string
	Timeout       time.Duration

	cmd        *exec.Cmd
	scriptPath string
}

func (b *Bash) Do() {
	file, err := ioutil.TempFile(b.ctx.LogDir(), "bash_"+b.Id)
	if err != nil {
		glog.Errorf("Create temp script failed: %v", err)
		b.state = FAILED
		return
	}

	b.scriptPath, err = filepath.Abs(file.Name())
	if err != nil {
		glog.Errorf("Get temp file path failed: %v", err)
		b.state = FAILED
		return
	}

	_, err = file.WriteString(
		fmt.Sprintf(
			"%v\n%v\n%v",
			SHEBANG,
			SET_E,
			b.ScriptContent))
	if err != nil {
		glog.Errorf("Write to script file failed: %v", err)
		b.state = FAILED
		return
	}
	file.Close()

	b.cmd = exec.Command("bash", b.scriptPath)
	b.cmd.Env = b.Env

	file, err = ioutil.TempFile(b.ctx.LogDir(), "bash_"+b.Id+"_out")
	if err != nil {
		glog.Errorf("Create temp out failed: %v", err)
		b.state = FAILED
		return
	}

	b.cmd.Stdout = file
	defer file.Close()

	file, err = ioutil.TempFile(b.ctx.LogDir(), "bash_"+b.Id+"_err")
	if err != nil {
		glog.Errorf("Create temp out failed: %v", err)
		b.state = FAILED
		return
	}
	b.cmd.Stderr = file
	defer file.Close()

	err = b.cmd.Start()
	if err != nil {
		glog.Errorf("Start Script err: %v", err)
		b.state = FAILED
		return
	}

	timer := time.NewTimer(time.Second * b.Timeout)
	go func() {
		<-timer.C
		// a chance for script to terminate by itself
		b.cmd.Process.Signal(syscall.SIGTERM)
		time.Sleep(time.Second * 5)
		if !b.cmd.ProcessState.Exited() {
			b.cmd.Process.Signal(syscall.SIGKILL)
		}
	}()

	err = b.cmd.Wait()
	if err != nil {
		glog.Errorf("Wait process failed: %v", err)
		b.state = FAILED
		return
	}

	if b.cmd.ProcessState.Success() {
		b.state = SUCCESS
	} else {
		b.state = FAILED
	}
}
