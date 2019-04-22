package main

import (
	"errors"
	"fmt"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"strconv"
	"strings"
)

func setArrayValue(msg *dynamic.Message, idx int, tp string, arrval []string) error {
	for _, val := range arrval {
		switch tp {
		case "int":
			if len(val) > 0 {
				pv, err := strconv.ParseInt(val, 10, 32)
				if err != nil {
					return err
				}
				msg.AddRepeatedFieldByNumber(idx, int32(pv))
			} else {
				msg.AddRepeatedFieldByNumber(idx, int32(0))
			}
		case "str":
			msg.AddRepeatedFieldByNumber(idx, val)
		case "float":
			if len(val) > 0 {
				pv, err := strconv.ParseFloat(val, 32)
				if err != nil {
					return err
				}
				msg.AddRepeatedFieldByNumber(idx, float32(pv))
			} else {
				msg.AddRepeatedFieldByNumber(idx, float32(0))
			}
		default:
			return errors.New("unknown type " + tp)
		}
	}
	return nil
}

func setValue(msg *dynamic.Message, idx int, tp string, val string) error {
	isArr := isArrayType(tp)
	if isArr {
		if len(val) == 0 {
			return nil
		}
		tp = rawType(tp)
		return setArrayValue(msg, idx, tp, strings.Split(val, "|"))
	}

	switch tp {
	case "int":
		if len(val) > 0 {
			pv, err := strconv.ParseInt(val, 10, 32)
			if err != nil {
				return err
			}
			msg.SetFieldByNumber(idx, int32(pv))
		} else {
			msg.SetFieldByNumber(idx, int32(0))
		}
	case "str":
		msg.SetFieldByNumber(idx, val)
	case "float":
		if len(val) > 0 {
			pv, err := strconv.ParseFloat(val, 32)
			if err != nil {
				return err
			}
			msg.SetFieldByNumber(idx, float32(pv))
		} else {
			msg.SetFieldByNumber(idx, float32(0))
		}
	default:
		return errors.New("unknown type " + tp)
	}
	return nil
}

func WriteProtoData(descriptor *desc.FileDescriptor, data *ExcelData, rowName int, rowType int, rowData int) []byte {
	rowDesc := descriptor.FindSymbol(xyl(data.name)).(*desc.MessageDescriptor)
	arrDesc := descriptor.FindSymbol(xyl(data.name) + "Array").(*desc.MessageDescriptor)

	dmArr := dynamic.NewMessage(arrDesc)

	for i := rowData; i < len(data.data); i++ {
		dmRow := dynamic.NewMessage(rowDesc)
		var key int64
		for idx, name := range data.data[rowName] {
			if len(data.data[i][0]) == 0 {
				continue
			}
			cellVal := ""
			if idx < len(data.data[i]) {
				cellVal = data.data[i][idx]
			}
			err := setValue(dmRow, idx+1, data.data[rowType][idx], cellVal)
			if idx == 0 {
				key, err = strconv.ParseInt(cellVal, 10, 64)
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
