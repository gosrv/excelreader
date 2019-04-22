package main

// excell 数据
type ExcelData struct {
	name string
	data [][]string
}

// 模板配置
type CfgTemplate struct {
	Package string
	MapType map[string]string
	Tmpl string
}

// meta文件
type CfgMeta struct {
	Meta struct{
		// 行号以0起始
		LineName int	// 名字的行号
		LineType int	// 类型的行号
		LineDes int		// 注释的行号
		LineData int	// 数据起始行
	}
	Templates map[string]*CfgTemplate
}