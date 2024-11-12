package osutil

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const windows = "windows"

func IsWindows() bool {
	return windows == runtime.GOOS
}

func toEnvVarsList(envVarsAsMap map[string]string) []string {
	var envVarsAsList = make([]string, len(envVarsAsMap))
	var i int
	for key, value := range envVarsAsMap {
		envVarsAsList[i] = fmt.Sprintf("%s=%s", key, value)
		i += 1
	}
	return envVarsAsList
}

var ErrorForceKill = errors.New("force killed failed")

type CmdOutput struct {
	Stdout *bytes.Buffer
	Stderr *bytes.Buffer
}

type Options struct {
	CancelCtx  context.Context
	Command    string
	CliArgs    []string
	BinPath    string
	Env        map[string]string
	ErrWriter  io.Writer
	Stdin      io.Reader
	WorkingDir string
}

const EnvVarTMP = "TMP"
const EnvVarHome = "HOME"

func Exec(opt *Options) (*CmdOutput, error) {
	logger.Debug(fmt.Sprintf("Running command: %s %s", opt.Command, strings.Join(opt.CliArgs, " ")))
	out := &CmdOutput{Stdout: &bytes.Buffer{}, Stderr: &bytes.Buffer{}}
	// 设置执行目录
	if opt.BinPath == "" {
		binary, err := exec.LookPath(opt.Command)
		if err != nil {
			return nil, err
		}
		opt.BinPath = binary
	}

	opt.Env[EnvVarTMP] = filepath.Clean(os.Getenv(EnvVarTMP))
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	opt.Env[EnvVarHome] = home

	cmd := exec.Command(opt.BinPath, opt.CliArgs...) //nolint

	if opt.Stdin == nil {
		opt.Stdin = os.Stdin
	}
	cmd.Env = toEnvVarsList(opt.Env)

	if opt.WorkingDir == "" {
		opt.WorkingDir = home
	}

	cmd.Path = opt.BinPath
	cmd.Dir = opt.WorkingDir
	cmd.Dir = filepath.ToSlash(cmd.Dir)

	if opt.ErrWriter != nil {
		cmd.Stderr = io.MultiWriter(opt.ErrWriter, out.Stderr)
	} else {
		cmd.Stderr = out.Stderr
	}
	cmd.Stdout = out.Stdout

	logger.Debug(fmt.Sprintf("Running command: %s %s at %s", opt.Command, strings.Join(opt.CliArgs, " "), opt.WorkingDir))
	logger.Debug(fmt.Sprintf("Running command: %s env is %v", opt.Command, cmd.Env))
	if err = cmd.Start(); err != nil {
		return out, err
	}
	// 完成通道
	finishCh := make(chan error, 1)
	go func() {
		cmdErr := cmd.Wait()
		if cmdErr != nil {
			finishCh <- cmdErr
		}
		finishCh <- nil
	}()
	select {
	case err = <-finishCh:
	case <-opt.CancelCtx.Done():
		if IsWindows() && cmd.Process.Kill() != nil {
			logger.Error("force kill app failed")
			err = ErrorForceKill
			return out, err
		}
		if cmd.Process.Signal(os.Interrupt) != nil {
			logger.Error(fmt.Sprint("force kill app failed, in ", runtime.GOOS))
			err = ErrorForceKill
		}
	}
	return out, err
}

var logger *slog.Logger

func init() {
	logger = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey && len(groups) == 0 {
					return slog.Attr{}
				}
				return a
			},
		}),
	).WithGroup("Gone Script")
}
