package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"net/http"
	"text/template"
)

var db *sqlx.DB

func initDB() (err error) {
	dbsource := "root:123456@tcp(127.0.0.1:3306)/test"
	db, err = sqlx.Connect("mysql", dbsource)
	if err != nil {
		return err
	}
	//连接成功
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(16)
	return
}

//创建用户的函数
func creatUser(username, password string) error {
	sqlStr := "insert into userinfo(username,password) values(?,?)"
	_, err := db.Exec(sqlStr, username, password)
	if err != nil {
		fmt.Println("插入用户数据失败!")
		return err
	}
	return nil
}

func queryUser(username, password string) error {
	sqlStr := "select id from userinfo where username=? and password=? limit 1"
	var id int64
	err := db.Get(&id, sqlStr, username, password)
	if err != nil {
		fmt.Println("查询用户数据失败!")
		return err
	}
	return nil
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	//根据请求方法不同做不同的处理
	t, err := template.ParseFiles("./register.html")
	//POST:提取用户提交的form表单数据,去数据库创建一条记录
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(500)
		}
		username := r.FormValue("username")
		password := r.FormValue("password")
		//往数据库写数据
		err = creatUser(username, password)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		http.Redirect(w, r, "/login", 301)
	} else {
		//GET:返回HTML页面

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		t.Execute(w, nil)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(500)
		}
		username := r.FormValue("username")
		password := r.FormValue("password")
		//去数据库校验
		err = queryUser(username, password)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}
		http.Redirect(w, r, "https://www.baidu.com", 301)
	} else {
		t, err := template.ParseFiles("./login.html")
		if err != nil {
			w.WriteHeader(500)
			return
		}
		t.Execute(w, nil)
	}
}

func main() {
	err := initDB()
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("启动http失败,err:", err)
		return
	}
}
