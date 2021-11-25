package sshx

import (
	"encoding/hex"
	"fmt"
	"github.com/lixianmin/got/convert"
	"math/rand"
	"os"
	"os/exec"
)

/********************************************************************
created:    2021-11-24
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type SSH struct {
	hostname string
	command  string
}

// NewSSH
// 1. format用于格式化ssh的参数部分，args可以有也可以没有
// 2. 因为是使用在目标机器上直接运行脚本的方式，因此只要命令在linux目标机上试验成功，直接copy命令就可以，不需要调整命令格式，不需要考虑 $ " 之类的特殊符号
func NewSSH(hostname, format string, args ...interface{}) *SSH {
	if hostname == "" {
		panic("hostname is empty")
	}

	if format == "" {
		panic("format is empty")
	}

	// 获取命令
	var command = format
	if len(args) > 0 {
		command = fmt.Sprintf(format, args...)
	}

	var my = &SSH{
		hostname: hostname,
		command:  command,
	}

	return my
}

func (my *SSH) Run() ([]byte, error) {
	var scriptName, err = writeScript(my.command)
	if err != nil {
		return nil, err
	}

	defer os.Remove(scriptName)

	// 复制脚本文件到远程主机
	output, err := exec.Command("scp", scriptName, my.hostname+":/tmp").CombinedOutput()
	if err != nil {
		return nil, err
	}

	// -tt 可避免 Pseudo-terminal will not be allocated because stdin is not a terminal
	// -o StrictHostKeyChecking=no  可避免 The authenticity of host 'xxx' can't be established. RSA key fingerprint is xxx. Are you sure you want to continue connecting (yes/no)
	var destFileName = "/tmp/" + scriptName
	var cmd = exec.Command("ssh", my.hostname, "-tt", "-o", "StrictHostKeyChecking=no", "/bin/bash", destFileName)
	output, err = cmd.CombinedOutput()

	// 删除目标机器上的临时文件
	// 这个最初试过在script中写入 rm self.xxx.sh的命令，后来发现有可能删除不掉
	_, _ = exec.Command("ssh", my.hostname, "-tt", "-o", "StrictHostKeyChecking=no", "/bin/rm", destFileName).CombinedOutput()
	return output, err
}

func (my *SSH) String() string {
	return my.command
}

func randomFileName(prefix, suffix string) string {
	var randBytes = make([]byte, 16)
	rand.Read(randBytes)
	return prefix + hex.EncodeToString(randBytes) + suffix
}

func writeScript(command string) (string, error) {
	var filename = randomFileName("ssh.", ".sh")
	fout, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return filename, err
	}

	defer fout.Close()

	_, err = fout.Write(convert.Bytes(command))
	return filename, err
}