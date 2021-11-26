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

type Script struct {
	hostname string
	script   string
	filename string
	options  scriptOptions
}

// NewScript script的内容会被当作shell脚本传输到目标机
func NewScript(hostname, script string, opts ...ScriptOption) *Script {
	if hostname == "" {
		panic("hostname is empty")
	}

	if script == "" {
		panic("script is empty")
	}

	// 默认值
	var options = scriptOptions{
		prefix: defaultPrefix,
	}

	// 初始化
	for _, opt := range opts {
		opt(&options)
	}

	var filename = options.prefix + sumSHA256(script) + ".sh"
	var my = &Script{
		hostname: hostname,
		script:   script,
		filename: filename,
		options:  options,
	}

	return my
}

func (my *Script) Run(args ...string) ([]byte, error) {
	// 这个方案有点redis的evalsha
	var output, err = my.runScript(args...)
	if err != nil { // 主要的目的是『如果没有则创建』
		var filename = my.filename
		err = filex.WriteAllText(filename, my.script)
		if err != nil {
			return nil, err
		}

		// remove本地的脚本文件
		defer os.Remove(filename)

		// scp脚本文件到远程主机
		output, err = exec.Command("scp", filename, my.hostname+":/tmp").CombinedOutput()
		if err != nil {
			return nil, err
		}

		output, err = my.runScript(args...)
	}

	return output, err
}

func (my *Script) runScript(args ...string) ([]byte, error) {
	var remotePath = "/tmp/" + my.filename

	// -tt 可避免 Pseudo-terminal will not be allocated because stdin is not a terminal
	// -o StrictHostKeyChecking=no  可避免 The authenticity of host 'xxx' can't be established. RSA key fingerprint is xxx. Are you sure you want to continue connecting (yes/no)
	var list = append([]string{my.hostname, "-tt", "-o", "StrictHostKeyChecking=no", "/bin/bash", remotePath}, args...)
	var cmd = exec.Command("ssh", list...)
	var output, err = cmd.CombinedOutput()
	return output, err
}

func (my *Script) String() string {
	return fmt.Sprintf("command=\"ssh %s /bin/bash /tmp/%s\", script=%q", my.hostname, my.filename, my.script)
}

func sumSHA256(input string) string {
	var data = convert.Bytes(input)
	var code = sha256.Sum256(data)
	return hex.EncodeToString(code[0:])
}
