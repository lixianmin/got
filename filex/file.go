package filex

import (
	"bufio"
	"fmt"
	"github.com/lixianmin/got/convert"
	"io"
	"os"
)

/********************************************************************
created:    2020-08-28
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func ReadLines(path string, handler func(line string) bool) error {
	var fin, err = os.Open(path)
	if err != nil {
		return err
	}

	defer fin.Close()

	return ForEachLine(fin, handler)
}

// ForEachLine 有些文件格式非常特殊，可能会把内存撑爆，这时可以考虑传入io.LimitReader(r, limit)，限制读入的最大字节数
func ForEachLine(fin io.Reader, handler func(line string) bool) error {
	if handler == nil {
		return fmt.Errorf("handler is nil")
	}

	var reader = bufio.NewReader(fin)
	for {
		var buffer, err = reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				// 最后一行了，没有结尾的'\n'，所以也不需要移除了
				var line = string(buffer)
				handler(line)
			}
			return err
		}

		var length = len(buffer)
		var lastIndex = length - 1 // \n的位置

		// windows下也需要处理 \r
		if length >= 2 && buffer[lastIndex-1] == '\r' {
			lastIndex--
		}

		var line = string(buffer[:lastIndex])
		if ok := handler(line); !ok { // 只要handler()返回false，就中止
			return nil
		}
	}
}

func ReadAllLines(path string) ([]string, error) {
	var fin, err = os.Open(path)
	if err != nil {
		return nil, err
	}

	defer fin.Close()

	var reader = bufio.NewReader(fin)
	var lines = make([]string, 0, 32)
	for {
		var buffer, err = reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				// 最后一行了，没有结尾的'\n'，所以也不需要移除了
				var line = string(buffer)
				lines = append(lines, line)
			}
			return lines, err
		}

		var line = string(buffer[:len(buffer)-1])
		lines = append(lines, line)
	}
}

func WriteAllText(path string, text string) error {
	var fout, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer fout.Close()

	var data = convert.Bytes(text)
	_, err = fout.Write(data)
	return err
}
