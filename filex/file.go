package filex

import (
	"bufio"
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

		var line = string(buffer[:len(buffer)-1])
		var ok = handler(line)
		if !ok {
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
