package main

type ExcelData struct {
	name string
	data [][]string
}

type CfgTemplate struct {
	Package string
	MapType map[string]string
	Tmpl string
}

type CfgMeta struct {
	Meta struct{
		LineName int
		LineType int
		LineDes int
		LineData int
	}
	Templates map[string]*CfgTemplate
}