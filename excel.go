package main

import (
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"strings"
)

func FineExcelFiles(dirPth string) []string {
	var files []string = nil
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		panic(err)
	}

	for _, fi := range dir {
		if !fi.IsDir() && strings.HasSuffix(fi.Name(), "xlsx") &&
			!strings.HasPrefix(fi.Name(), "~$") {
			files = append(files, fi.Name())
		}
	}
	return files
}

func ReadExcel(name string) [][]string {
	xlFile, err := xlsx.OpenFile(name)
	if err != nil {
		panic(err)
	}
	var vals [][]string = nil
	sheet := xlFile.Sheets[0]
	for _, row := range sheet.Rows {
		rowVals := make([]string, len(row.Cells))
		vals = append(vals, rowVals)
		for idx, cell := range row.Cells {
			text := cell.Value
			rowVals[idx] = text
		}
	}

	return vals
}
