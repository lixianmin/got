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
			break
		}

		var line = string(buffer[:len(buffer)-1])
		var ok = handler(line)
		if !ok {
			break
		}
	}

	if err != io.EOF {
		return err
	}

	return nil
}
