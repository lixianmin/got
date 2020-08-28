package filex

import (
	"bufio"
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
			break
		}

		var line = string(buffer[:len(buffer)-1])
		var ok = handler(line)
		if !ok {
			break
		}
	}

	return err
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
			break
		}

		var line = string(buffer[:len(buffer)-1])
		lines = append(lines, line)
	}

	return lines, err
}
