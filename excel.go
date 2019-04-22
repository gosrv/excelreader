package main

import (
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"strings"
)

// 查找目录下的excel文件，不包含子目录
func FindExcelFiles(dirPth string) []string {
	var files []string = nil
	dir, err := ioutil.ReadDir(dirPth)
	verifyNoError(err)

	for _, fi := range dir {
		// 只搜索excel文件
		if !fi.IsDir() && strings.HasSuffix(fi.Name(), "xlsx") &&
			// 不包含临时文件
			!strings.HasPrefix(fi.Name(), "~$") {
			files = append(files, fi.Name())
		}
	}
	return files
}

// 读取excel数据，转为一个二维字符串数组
func ReadExcel(name string) [][]string {
	xlFile, err := xlsx.OpenFile(name)
	verifyNoError(err)

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
