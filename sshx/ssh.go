package sshx

import (
	"encoding/hex"
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
	script   string
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
		debug:  false,
	}

	// 初始化
	for _, opt := range opts {
		opt(&options)
	}

	var my = &SSH{
		hostname: hostname,
		script:   script,
		options:  options,
	}

	return my
}

func (my *SSH) Run() ([]byte, error) {
	var scriptName, err = writeScript(my.script, my.options.prefix)
	if err != nil {
		return nil, err
	}

	var isDebug = my.options.debug
	defer func() {
		if !isDebug {
			_ = os.Remove(scriptName)
		}
	}()

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
	// 这个最初试过在script中写入 rm self.xxx.sh的命令，后来发现有可能删除不掉 (也许是因为脚本执行出错了？)
	if !isDebug {
		_, _ = exec.Command("ssh", my.hostname, "-tt", "-o", "StrictHostKeyChecking=no", "/bin/rm", destFileName).CombinedOutput()
	}

	return output, err
}

func (my *SSH) String() string {
	return my.script
}

func randomFileName(prefix, suffix string) string {
	var randBytes = make([]byte, 16)
	rand.Read(randBytes)
	return prefix + hex.EncodeToString(randBytes) + suffix
}

func writeScript(script string, prefix string) (string, error) {
	var filename = randomFileName(prefix, ".sh")
	fout, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return filename, err
	}

	defer fout.Close()

	_, err = fout.Write(convert.Bytes(script))
	return filename, err
}
