package tools

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func UpperString(s string) string {
	return cases.Title(language.Und, cases.NoLower).String(s)
}

func ConstructorFunc(ss *StructStr) FuncStr {
	casedStruct := UpperString(ss.Name)
	newName := fmt.Sprintf("New%s", casedStruct)
	pms := make([]string, len(ss.Fields))
	var bld strings.Builder
	for _, str := range ss.CnstBd {
		bld.WriteString(fmt.Sprintf("%s\n", str))
	}
	bld.WriteString(fmt.Sprintf("return &%s{\n", ss.Name))
	for i, field := range ss.Fields {
		pms[i] = fmt.Sprintf("%s %s", field.Param, field.Type)
		bld.WriteString(fmt.Sprintf("%s: %s,\n", field.Name, field.Param))
	}
	bld.WriteRune('}')
	return FuncStr{
		Name:   newName,
		Params: pms,
		Return: []string{fmt.Sprintf("*%s", ss.Name)},
		Body:   bld.String(),
	}
}

func Setters(ss *StructStr) []FuncStr {
	strs := make([]FuncStr, len(ss.Fields))
	for i, field := range ss.Fields {
		var bld strings.Builder
		for _, str := range field.SetBd {
			bld.WriteString(fmt.Sprintf("%s\n", str))
		}
		bld.WriteString(fmt.Sprintf("%s.%s = %s", ss.Param, field.Name, field.Param))
		strs[i] = FuncStr{
			Name:     fmt.Sprintf("Set%s", UpperString(field.Name)),
			Receiver: fmt.Sprintf("%s *%s", ss.Param, ss.Name),
			Params:   []string{fmt.Sprintf("%s %s", field.Param, field.Type)},
			Body:     bld.String(),
		}
	}
	return strs
}

func Getters(ss *StructStr) []FuncStr {
	gtrs := make([]FuncStr, len(ss.Fields))
	for i, field := range ss.Fields {
		gtrs[i] = FuncStr{
			Name:     UpperString(field.Name),
			Receiver: fmt.Sprintf("%s *%s", ss.Param, ss.Name),
			Return:   []string{field.Type},
			Body:     fmt.Sprintf("return %s.%s", ss.Param, field.Name),
		}
	}
	return gtrs
}

func VisitSig(ss *StructStr, void bool) FuncStr {
	casedAcceptor := UpperString(ss.Name)
	name := fmt.Sprintf("Visit%s", casedAcceptor)
	retstr := "interface{}"
	if void {
		retstr = ""
	}
	return FuncStr{
		Name: name,
		Params: []string{
			fmt.Sprintf("%s *%s", ss.Param, ss.Name),
		},
		Return: []string{retstr},
	}
}

func AcceptMethod(ss *StructStr, void bool) FuncStr {
	bdfmt := "return v.%s(%s)"
	retstr := "interface{}"
	if void {
		bdfmt = "v.%s(%s)"
		retstr = ""
	}
	casedAcceptor := UpperString(ss.Name)
	visitName := fmt.Sprintf("Visit%s", casedAcceptor)
	return FuncStr{
		Name:     "Accept",
		Receiver: fmt.Sprintf("%s *%s", ss.Param, ss.Name),
		Params:   []string{"v Visitor"},
		Body:     fmt.Sprintf(bdfmt, visitName, ss.Param),
		Return:   []string{retstr},
	}
}
