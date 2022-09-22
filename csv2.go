package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

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

type Meili struct {
	Index   string   // 索引编号
	Year    string   // 年份 	【人工审核位】
	School  string   // 院校
	Des     string   // 调剂专业
	AMD     bool     // 是学硕吗？
	Same    bool     // 同校调剂？
	Src     []string // 接收专业
	Subject []string //统考科目要求
	Line    struct { // 分数线
		Sum int
		Big int
		Sma int
	}
	HC   int    // 招生人数
	Link string // 来源网址
	PS   string //备注
}

func main() {
	var err error
	var csvInput = ReadCsv()
	var users []User
	err = csvutil.Unmarshal(csvInput, &users)
	if err != nil {
		fmt.Println("csv error:", err)
	}
	//	{Index:000001 Year:2022 School:东北大学 DesCode:0801 DesName:力学 AMD:学硕 Same:否 Subject:数一,英一
	//	HC:4  SrcCode:0801,0814,0815 SrcName:力学,土木工程,水利工程 LineSum:320
	//	LineBig:75  LineSma:50  PS: Link:http://www.zitu.neu.edu.cn/2022/0402/c1036a211931/page.htm}

	meilis := ToJson(users)
	// for _, v := range meilis {
	// 	fmt.Printf("%+v\n", v)
	// }

	b, err := json.Marshal(meilis)
	if err != nil {
		fmt.Println("json error:", err)
	}
	// os.Stdout.Write(b)

	err = os.WriteFile("output.json", b, 0666)
	if err != nil {
		fmt.Println("write error", err)
	}
	fmt.Println("Done!")
}

func ReadCsv() []byte {
	f, err := os.Open("表格视图.csv")
	if err != nil {
		fmt.Println("read f error", err)
		return []byte{0}
	}
	defer f.Close()

	fd, err := io.ReadAll(f)
	if err != nil {
		fmt.Println("read fd error", err)
		return []byte{0}
	}

	return fd
}

func ToJson(users []User) []Meili {
	var meilis []Meili
	var mei Meili
	var err error
	for _, user := range users {
		// fmt.Printf("%+v\n", user)
		if user.Year == "" {
			continue
		}

		mei.Index = user.Index
		mei.Year = user.Year
		mei.School = user.School
		mei.Des = user.DesCode + " " + user.DesName
		mei.AMD = (user.AMD == "学硕")
		mei.Same = (user.Same != "否")
		src_code := strings.Split(user.SrcCode, ",")
		src_name := strings.Split(user.SrcName, ",")
		for k := range src_code {
			mei.Src = append(mei.Src, src_code[k]+" "+src_name[k])
		}
		mei.Subject = strings.Split(user.Subject, ",")

		// 可能增加国家线选项
		mei.Line.Sum, err = strconv.Atoi(strings.TrimSpace(user.LineSum))
		mei.Line.Big, _ = strconv.Atoi(strings.TrimSpace(user.LineBig))
		mei.Line.Sma, _ = strconv.Atoi(strings.TrimSpace(user.LineSma))
		if err != nil {
			fmt.Println("error:", err)
			break
		}

		mei.HC, _ = strconv.Atoi(strings.TrimSpace(user.HC))
		mei.Link = user.Link
		mei.PS = user.PS
		meilis = append(meilis, mei)
	}
	return meilis
}
