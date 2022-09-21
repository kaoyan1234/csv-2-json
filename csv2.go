package main

import (
	"fmt"
	"io"
	"os"

	"github.com/jszwec/csvutil"
)

type User struct {
	Index   string `csv:"编号"`
	Year    string `csv:"年份"`
	School  string `csv:"院校"`
	DesCode string `csv:"调剂专业（代码）"`
	DesName string `csv:"调剂专业（名称）"`
	AMD     string `csv:"学位类型"`
	Same    string `csv:"同校调剂？"`
	Subject string `csv:"统考科目要求"`
	HC      string `csv:"招生计划"`
	SrcCode string `csv:"报考专业要求（代码）"`
	SrcName string `csv:"报考专业要求（学科名称）"`
	LineSum string `csv:"初试分数要求（总分）"`
	LineBig string `csv:"初试分数要求（大分）"`
	LineSma string `csv:"初试分数要求（小分）"`
	PS      string `csv:"备注"`
	Link    string `csv:"来源网址"`
}

func main() {
	var csvInput = ReadCsv()

	var users []User

	err := csvutil.Unmarshal(csvInput, &users)
	if err != nil {
		fmt.Println("error:", err)
	}

	for _, u := range users {
		fmt.Printf("%+v\n", u)
	}
}

func ReadCsv() []byte {
	f, err := os.Open("表格视图.csv")
	if err != nil {
		fmt.Println("read f fail", err)
		return []byte{0}
	}
	defer f.Close()

	fd, err := io.ReadAll(f)
	if err != nil {
		fmt.Println("read fd fail", err)
		return []byte{0}
	}

	return fd
}
