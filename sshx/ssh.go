package sshx

import (
	"fmt"
	"os/exec"
	"strings"
)

/********************************************************************
created:    2021-11-24
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type SSH struct {
	args []string
}

// NewSSH 创建一条ssh命令。注：1. format只用于格式化参数部分，不带ssh 2. format命令中间以空格分隔
func NewSSH(address, format string, a ...interface{}) *SSH {
	if address == "" {
		panic("address is empty")
	}

	if format == "" {
		panic("format is empty")
	}

	// -tt 可避免 Pseudo-terminal will not be allocated because stdin is not a terminal
	// -o StrictHostKeyChecking=no  可避免 The authenticity of host 'xxx' can't be established. RSA key fingerprint is xxx. Are you sure you want to continue connecting (yes/no)
	var text = address + " -tt -o StrictHostKeyChecking=no " + fmt.Sprintf(format, a...)
	var args = strings.Split(text, " ")
	var my = &SSH{args: args}
	return my
}

func (my *SSH) Run() ([]byte, error) {
	var cmd = exec.Command("ssh", my.args...)
	var output, err = cmd.CombinedOutput()
	return output, err
}

func (my *SSH) String() string {
	return strings.Join(my.args, " ")
}
