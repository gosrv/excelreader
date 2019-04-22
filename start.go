package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jhump/protoreflect/desc/protoparse"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// 生成模板文件
func generateTemplateFile(data []*ExcelData, name string, fname string, meta *CfgMeta) {
	f, err := os.Create(fname)
	verifyNoError(err)
	protoData := GenerateProto(data, meta.Templates[name],
		meta.Meta.LineName, meta.Meta.LineType, meta.Meta.LineDes)
	_, err = f.Write(protoData)
	verifyNoError(err)
	err = f.Close()
	verifyNoError(err)
}

func main() {
	metaFile := flag.String("meta", "", "meta file")
	excelDir := flag.String("excel", "", "excel dir")
	outputFile := flag.String("out", "", "output dir")
	tmplName := flag.String("tmpl", "", "template name")
	generateData := flag.String("data", "", "template name")
	flag.Parse()
	if len(*excelDir) == 0 || len(*outputFile) == 0 || (len(*tmplName) == 0 && len(*generateData) == 0) {
		flag.PrintDefaults()
		return
	}
	err := os.MkdirAll(path.Dir(*outputFile), os.ModePerm)
	verifyNoError(err)

	excelFiles := FindExcelFiles(*excelDir)
	var edatas []*ExcelData = nil
	for _, excelFile := range excelFiles {
		fmt.Println("parse excel file ", excelFile)
		edata := &ExcelData{
			name: strings.Split(excelFile, ".")[0],
			data: ReadExcel(path.Join(*excelDir, excelFile)),
		}
		edatas = append(edatas, edata)
	}

	f, err := os.Open(*metaFile)
	verifyNoError(err)

	metaData, err := ioutil.ReadAll(f)
	verifyNoError(err)

	cfgMeta := &CfgMeta{}
	err = json.Unmarshal(metaData, cfgMeta)
	verifyNoError(err)

	if len(*generateData) > 0 {
		// 生成一个proto临时文件，反射需要用
		tmpProtoFile := path.Join(*outputFile, "__table.tmp.proto")
		generateTemplateFile(edatas, "proto", tmpProtoFile, cfgMeta)
		defer os.Remove(tmpProtoFile)
		// 解析proto
		parser := &protoparse.Parser{
			ImportPaths: []string{},
		}
		fileDesc, err := parser.ParseFiles(tmpProtoFile)
		verifyNoError(err)
		// 生成proto数据
		for _, edata := range edatas {
			fmt.Println("begin write data ", edata.name)
			protoData := WriteProtoData(fileDesc[0], edata, cfgMeta.Meta.LineName, cfgMeta.Meta.LineType, cfgMeta.Meta.LineData)
			protoDataFile, err := os.Create(path.Join(*outputFile, strings.ToLower(edata.name) + ".bytes"))
			verifyNoError(err)

			_, err = protoDataFile.Write(protoData)
			verifyNoError(err)
			err = protoDataFile.Close()
			verifyNoError(err)
		}
		return
	} else {
		// 生成模板文件
		generateTemplateFile(edatas, *tmplName, *outputFile, cfgMeta)
	}

	return
}
