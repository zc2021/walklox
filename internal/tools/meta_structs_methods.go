package tools

import "strings"

type FuncStr struct {
	Receiver string
	Name     string
	Params   []string
	Body     string
	Return   []string
}

func (f *FuncStr) ListParams() string {
	return strings.Join(f.Params, ", ")
}

func (f *FuncStr) ListReturns() string {
	return strings.Join(f.Return, ", ")
}

type InterfaceStr struct {
	Name string
	Sigs []FuncStr
}

type StructStr struct {
	Name   string
	Param  string
	Fields []FieldStr
	CnstBd []string
}

type FieldStr struct {
	Name  string
	Param string
	Type  string
	SetBd []string
}
