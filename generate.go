package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"text/template"
)

type TableField struct {
	Name  string
	Arr   bool
	Desc  string
	Idx   int
	Meta  string
	Type string
}

type TableClass struct {
	Name   string
	Fields []*TableField
}

func isArrayType(extp string) bool {
	return strings.HasSuffix(extp, "|arr")
}

func rawType(tp string) string {
	Arr := isArrayType(tp)
	if Arr {
		tp = tp[:len(tp)-len("|arr")]
	}
	return tp
}

func xyl(raw string) string {
	return strings.ToUpper(raw[0:1]) + raw[1:]
}

func GenerateProto(data []*ExcelData, cfgTmpl *CfgTemplate, rowName int, rowType int, rowDesc int) []byte {
	var tclzs []*TableClass = nil
	for _, edata := range data {
		tclass := &TableClass{}
		tclzs = append(tclzs, tclass)

		tclass.Name = xyl(edata.name)
		for idx, fname := range edata.data[rowName] {
			ftp := edata.data[rowType][idx]
			Arr := isArrayType(ftp)
			ftp = rawType(ftp)

			tfield := &TableField{
				Name:  xyl(fname),
				Arr: Arr,
				Desc:  edata.data[rowDesc][idx],
				Idx:   idx + 1,
				Type: cfgTmpl.MapType[ftp],
			}
			if len(tfield.Name) == 0 {
				panic(fmt.Sprintf("unknown field name [%v : %v]", tclass.Name, fname))
			}
			tclass.Fields = append(tclass.Fields, tfield)
		}
	}

	bwriter := bytes.NewBuffer(nil)
	fdata, err := ioutil.ReadFile(cfgTmpl.Tmpl)
	verifyNoError(err)
	t := template.New("")
	t, err = t.Parse(string(fdata))
	verifyNoError(err)
	err = t.Execute(bwriter, tclzs)
	if err != nil {
		panic(err)
	}
	return bwriter.Bytes()
}