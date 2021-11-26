package sshx

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/lixianmin/got/convert"
	"github.com/lixianmin/got/filex"
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
	script   string
	sha256   string
	options  sshOptions
}

// NewSSH script的内容会被当作shell脚本传输到目标机
func NewSSH(hostname, script string, opts ...SSHOption) *SSH {
	if hostname == "" {
		panic("hostname is empty")
	}

	if script == "" {
		panic("script is empty")
	}

	// 默认值
	var options = sshOptions{
		prefix: defaultPrefix,
	}

	// 初始化
	for _, opt := range opts {
		opt(&options)
	}

	var my = &SSH{
		hostname: hostname,
		script:   script,
		sha256:   sumSHA256(script),
		options:  options,
	}

	return my
}

func (my *SSH) Run(args ...string) ([]byte, error) {
	var scriptName = my.options.prefix + my.sha256 + ".sh"
	// 这个方案有点redis的evalsha
	var output, err = my.runScript(scriptName, args...)
	if err != nil { // 主要的目的是『如果没有则创建』
		err = filex.WriteAllText(scriptName, my.script)
		if err != nil {
			return nil, err
		}

		// remove本地的脚本文件
		defer os.Remove(scriptName)

		// scp脚本文件到远程主机
		output, err = exec.Command("scp", scriptName, my.hostname+":/tmp").CombinedOutput()
		if err != nil {
			return nil, err
		}

		output, err = my.runScript(scriptName, args...)
	}

	return output, err
}

func (my *SSH) runScript(scriptName string, args ...string) ([]byte, error) {
	var remotePath = "/tmp/" + scriptName

	// -tt 可避免 Pseudo-terminal will not be allocated because stdin is not a terminal
	// -o StrictHostKeyChecking=no  可避免 The authenticity of host 'xxx' can't be established. RSA key fingerprint is xxx. Are you sure you want to continue connecting (yes/no)
	var list = append([]string{my.hostname, "-tt", "-o", "StrictHostKeyChecking=no", "/bin/bash", remotePath}, args...)
	var cmd = exec.Command("ssh", list...)
	var output, err = cmd.CombinedOutput()
	return output, err
}

func (my *SSH) String() string {
	return fmt.Sprintf("hostname=%q, sha256=%q, script=%q", my.hostname, my.sha256, my.script)
}

func sumSHA256(input string) string {
	var data = convert.Bytes(input)
	var code = sha256.Sum256(data)
	return hex.EncodeToString(code[0:])
}
