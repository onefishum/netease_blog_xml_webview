package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 使用 https://crypot.51strive.com/jsontogo/ 自动转换

// Netease 网易blogxml结构
type Netease struct {
	XMLName xml.Name `xml:"root"`
	Text    string   `xml:",chardata"`
	Blog    []struct {
		Text         string `xml:",chardata"`
		ID           int    `xml:"id"`
		UserID       string `xml:"userId"`
		UserName     string `xml:"userName"`
		UserNickname string `xml:"userNickname"`
		Title        string `xml:"title"`
		PublishTime  int64  `xml:"publishTime"`
		Ispublished  string `xml:"ispublished"`
		ClassID      string `xml:"classId"`
		ClassName    string `xml:"className"`
		AllowView    string `xml:"allowView"`
		Content      string `xml:"content"`
		Valid        string `xml:"valid"`
		MoveForm     string `xml:"moveForm"`
	} `xml:"blog"`
}

var stNetease Netease
var htmlBlogList = `
<!DOCTYPE html>
    <html lang="zh-cn">

    <head>
        <meta charset="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <meta name="referrer" content="origin" />
        <title>blog 列表</title>
        <style>
            ul { margin: 0; padding: 0; list-style: none; } 
            .table { display: table; border-collapse: collapse; border: 1px solid #ccc; }
            .table-caption { display: table-caption; margin: 0; padding: 0; font-size: 16px; }
            .table-column-group { display: table-column-group; }
            .table-column { display: table-column; width: 100px; }
            .table-row-group { display: table-row-group; }
            .table-row { display: table-row; }
            .table-row-group .table-row:hover, .table-footer-group .table-row:hover { background: #f6f6f6; }
            .table-cell { display: table-cell; padding: 0 5px; border: 1px solid #ccc; }
            .table-header-group { display: table-header-group; background: #eee; font-weight: bold; }
            .table-footer-group { display: table-footer-group; }
        </style>
    </head>

    <body>
        <div class="table">
            <h2 class="table-caption">Blog列表</h2>
            <div class="table-column-group">
                <div class="table-column"></div>
                <div class="table-column"></div>
                <div class="table-column"></div>
                <div class="table-column"></div>
            </div>
            <div class="table-header-group">
                <ul class="table-row">
                    <li class="table-cell">id</li>
                    <li class="table-cell">建立日期</li>
                    <li class="table-cell">标题</li>
                    <li class="table-cell">标签</li>

                </ul>
            </div>
            <!--    
            <div class="table-footer-group">
                <ul class="table-row">
                    <li class="table-cell">footer</li>
                    <li class="table-cell">footer</li>
                    <li class="table-cell">footer</li>
                    <li class="table-cell">footer</li>
                </ul>
            </div>
            -->
            %s
        </div>
    </body>
</html>
`

var htmlContent = `<!DOCTYPE html>
<html lang="zh-cn">

<head>
	<meta charset="utf-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
	<meta name="referrer" content="origin" />
	<title>blog %s</title>
</head>
<body>
%s
</body>
</html>
`

func main() {
	var filename string

	flag.StringVar(&filename, "f", "", "网易博客xml文件")
	flag.Parse()
	if filename == "" {
		fmt.Println("Usage: \n\tnetease_blog -f xxx.xml")
		return
	}
	fmt.Println("读取数据")
	getNetese(filename)
	fmt.Println("数据读取完成")

	// 生成页面数据
	createBlogList()
	fmt.Println("列表数据生成完成")

	fmt.Println("启动WEB服务器")
	http.HandleFunc("/", listBlog)
	http.HandleFunc("/blog", showBlog)

	HServer := http.Server{
		Addr:    "0.0.0.0:80",
		Handler: http.DefaultServeMux,
	}
	fmt.Println(HServer.Addr)

	err := HServer.ListenAndServe() // 设置监听的端口
	ErrHandler(err)

	fmt.Println("服务已经停止")
}

// listBlog 博客列表页
func listBlog(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, htmlBlogList)
}

// showBlog 显示blog内容
func showBlog(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	ids, ok := vars["id"]
	if ok {
		id, err := strconv.Atoi(ids[0])
		if err == nil {
			fmt.Fprint(w, fmt.Sprintf(htmlContent, stNetease.Blog[id].Title, stNetease.Blog[id].Content))
		}
	}
}

// createBlogList 生成博客列表页
func createBlogList() {
	item := `<ul class="table-row">
	<li class="table-cell">%d</li>
	<li class="table-cell">%s</li>
	<li class="table-cell"><a href="/blog?id=%d">%s </a></li>
	<li class="table-cell">%s</li>
</ul>`
	var blogList string
	for index, blog := range stNetease.Blog {

		tm := time.Unix(blog.PublishTime/1000, 0)
		blogList += fmt.Sprintf(item, blog.ID, tm.Format("2006-01-02 03:04:05 PM"), index, blog.Title, blog.ClassName)
	}
	htmlBlogList = fmt.Sprintf(htmlBlogList, blogList)
}

func getNetese(filename string) {
	content, err := ioutil.ReadFile("test.xml")
	ErrHandler(err)

	xmls := strings.Replace(string(content), string("\u000F"), "/", -1)
	xmls = strings.Replace(xmls, string("\u000E"), ".", -1)
	err = xml.Unmarshal([]byte(xmls), &stNetease) //将文件转化成对象
	ErrHandler(err)
	fmt.Println()
}

// ErrHandler 错误处理函数
func ErrHandler(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}
