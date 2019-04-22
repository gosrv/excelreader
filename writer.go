package main

import (
	"fmt"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"strconv"
	"strings"
)

type ParseFunc map[string]func(string)interface{}

var ParseFuncs = ParseFunc {
	"int":func (val string) interface{} {
		if len(val) == 0 {
			return int32(0)
		}
		pv, err := strconv.ParseInt(val, 10, 32)
		verifyNoError(err)
		return int32(pv)
	},
	"str":func (val string) interface{} {
		return val
	},
	"float":func (val string) interface{} {
		if len(val) == 0 {
			return float32(0)
		}
		pv, err := strconv.ParseFloat(val, 32)
		verifyNoError(err)
		return float32(pv)
	},
}

func setArrayValue(msg *dynamic.Message, idx int, tp string, arrval []string) {
	for _, val := range arrval {
		pfunc, ok := ParseFuncs[tp]
		if !ok {
			panic("not support type " + tp)
		}
		pval := pfunc(val)
		msg.AddRepeatedFieldByNumber(idx, pval)
	}
}

func setValue(msg *dynamic.Message, idx int, tp string, val string) {
	isArr := isArrayType(tp)
	if isArr {
		if len(val) == 0 {
			return
		}
		tp = rawType(tp)
		setArrayValue(msg, idx, tp, strings.Split(val, "|"))
		return
	}

	pfunc, ok := ParseFuncs[tp]
	if !ok {
		panic("not support type " + tp)
	}
	pval := pfunc(val)
	msg.SetFieldByNumber(idx, pval)
}

func WriteProtoData(descriptor *desc.FileDescriptor, data *ExcelData, rowName int, rowType int, rowData int) []byte {
	rowDesc := descriptor.FindSymbol(xyl(data.name)).(*desc.MessageDescriptor)
	arrDesc := descriptor.FindSymbol(xyl(data.name) + "Array").(*desc.MessageDescriptor)

	dmArr := dynamic.NewMessage(arrDesc)

	for i := rowData; i < len(data.data); i++ {
		dmRow := dynamic.NewMessage(rowDesc)
		var key int64
		var err error
		for idx, name := range data.data[rowName] {
			if len(data.data[i][0]) == 0 {
				continue
			}
			cellVal := ""
			if idx < len(data.data[i]) {
				cellVal = data.data[i][idx]
			}
			setValue(dmRow, idx+1, data.data[rowType][idx], cellVal)
			if idx == 0 {
				key, err = strconv.ParseInt(cellVal, 10, 32)
				if err != nil || len(cellVal) == 0 {
					panic(fmt.Sprintf("data parse error %s:%s:%s", data.name, name, err))
				}
			}
			if err != nil {
				panic(fmt.Sprintf("data parse error %s:%s:%s", data.name, name, err))
			}
		}
		if key > 0 {
			dmArr.AddRepeatedFieldByNumber(1, int32(key))
			dmArr.AddRepeatedFieldByNumber(2, dmRow)
		}
	}

	protoData, err := dmArr.Marshal()
	verifyNoError(err)
	return protoData
}
