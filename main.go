package main

import (
	"database/sql"
	"embed"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/mereith/nav/goscraper"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
	// _ "github.com/mattn/go-sqlite3"
)

const INDEX = "index.html"

func getIcon(url string) string {
	s, err := goscraper.Scrape(url, 5)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	var result string = ""
	if strings.Contains(s.Preview.Icon, "http:") || strings.Contains(s.Preview.Icon, "https:") {
		result = s.Preview.Icon
	} else {
		//  如果 link 最后一个是 /
		var first string = s.Preview.Link
		var second string = s.Preview.Icon
		if !strings.Contains(s.Preview.Link[len(s.Preview.Link)-1:len(s.Preview.Link)], "/") {
			first = s.Preview.Link + "/"
		}
		// 如果 icon 第一个是 /
		if strings.Contains(s.Preview.Icon[0:1], "/") {
			second = s.Preview.Icon[1:len(s.Preview.Icon)]
		}
		result = first + second
	}
	fmt.Println(result)
	return result
}

func updateCatelog(data updateCatelogDto, db *sql.DB) {
	sql_update_catelog := `
		UPDATE nav_catelog
		SET name = ?
		WHERE id = ?;
		`
	stmt, err := db.Prepare(sql_update_catelog)
	checkErr(err)
	res, err := stmt.Exec(data.Name, data.Id)
	checkErr(err)
	_, err = res.RowsAffected()
	checkErr(err)
	// fmt.Println(affect)
}

func updateTool(data updateToolDto, db *sql.DB) {
	sql_update_tool := `
		UPDATE nav_table
		SET name = ?, url = ?, logo = ?, catelog = ?, desc = ?
		WHERE id = ?;
		`
	stmt, err := db.Prepare(sql_update_tool)
	checkErr(err)
	res, err := stmt.Exec(data.Name, data.Url, data.Logo, data.Catelog, data.Desc, data.Id)
	checkErr(err)
	_, err = res.RowsAffected()
	checkErr(err)
	// fmt.Println(affect)
}

func updateSetting(data Setting, db *sql.DB) {
	sql_update_setting := `
		UPDATE nav_setting
		SET favicon = ?, title = ?
		WHERE id = ?;
		`
	stmt, err := db.Prepare(sql_update_setting)
	checkErr(err)
	res, err := stmt.Exec(data.Favicon, data.Title, 0)
	checkErr(err)
	_, err = res.RowsAffected()
	checkErr(err)
	// fmt.Println(affect)
}

func addApiTokenInDB(data Token, db *sql.DB) {
	sql_add_api_token := `
		INSERT INTO nav_api_token (id,name,value,disabled)
		VALUES (?,?,?,?);
		`
	// fmt.Println("增加分类：",data)
	stmt, err := db.Prepare(sql_add_api_token)
	checkErr(err)

	res, err := stmt.Exec(data.Id, data.Name, data.Value, data.Disabled)
	checkErr(err)
	_, err = res.LastInsertId()
	checkErr(err)
}

func updateUser(data updateUserDto, db *sql.DB) {
	sql_update_user := `
		UPDATE nav_user
		SET name = ?, password = ?
		WHERE id = ?;
		`
	stmt, err := db.Prepare(sql_update_user)
	checkErr(err)
	res, err := stmt.Exec(data.Name, data.Password, data.Id)
	checkErr(err)
	_, err = res.RowsAffected()
	checkErr(err)
	// fmt.Println(affect)
}

func addCatelog(data addCatelogDto, db *sql.DB) {
	// 先检查重复不重复
	existCatelogs := getAllCatelog(db)
	var existCatelogsArr []string
	for _, catelogDto := range existCatelogs {
		existCatelogsArr = append(existCatelogsArr, catelogDto.Name)
	}
	if in(data.Name, existCatelogsArr) {
		return
	}
	sql_add_catelog := `
		INSERT INTO nav_catelog (id,name)
		VALUES (?,?);
		`
	// fmt.Println("增加分类：",data)
	stmt, err := db.Prepare(sql_add_catelog)
	checkErr(err)
	res, err := stmt.Exec(generateId(), data.Name)
	checkErr(err)
	_, err = res.LastInsertId()
	checkErr(err)
	// fmt.Println(id)
}

func addTool(data addToolDto, db *sql.DB) {
	sql_add_tool := `
		INSERT INTO nav_table (id,name, url, logo, catelog, desc)
		VALUES (?, ?, ?, ?, ?, ?);
		`
	stmt, err := db.Prepare(sql_add_tool)
	checkErr(err)
	res, err := stmt.Exec(generateId(), data.Name, data.Url, data.Logo, data.Catelog, data.Desc)
	checkErr(err)
	_, err = res.LastInsertId()
	checkErr(err)
	// fmt.Println(id)
}

func getAllTool(db *sql.DB) []Tool {
	sql_get_all := `
		SELECT * FROM nav_table;
		`
	results := make([]Tool, 0)
	rows, err := db.Query(sql_get_all)
	checkErr(err)
	for rows.Next() {
		var tool Tool
		err = rows.Scan(&tool.Id, &tool.Name, &tool.Url, &tool.Logo, &tool.Catelog, &tool.Desc)
		checkErr(err)
		results = append(results, tool)
	}
	defer rows.Close()
	return results
}

func getAllCatelog(db *sql.DB) []Catelog {
	sql_get_all := `
		SELECT * FROM nav_catelog;
		`
	results := make([]Catelog, 0)
	rows, err := db.Query(sql_get_all)
	checkErr(err)
	for rows.Next() {
		var catelog Catelog
		err = rows.Scan(&catelog.Id, &catelog.Name)
		checkErr(err)
		results = append(results, catelog)
	}
	defer rows.Close()
	return results
}

func generateId() int {
	// 生成一个随机 id
	return int(time.Now().Unix())
}

var db *sql.DB

func PathExistsOrCreate(path string) {
	_, err := os.Stat(path)
	if err == nil {
		return
	}
	os.Mkdir(path, os.ModePerm)
}

func initDB() {
	PathExistsOrCreate("./data")
	// 创建数据库
	db, _ = sql.Open("sqlite", "./data/nav.db")
	// user 表
	sql_create_table := `
		CREATE TABLE IF NOT EXISTS nav_user (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			password TEXT
		);
		`
	_, err := db.Exec(sql_create_table)
	checkErr(err)
	// setting 表
	sql_create_table = `
	CREATE TABLE IF NOT EXISTS nav_setting (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		favicon TEXT,
		title TEXT
	);
	`
	_, err = db.Exec(sql_create_table)
	checkErr(err)
	// 默认 tools 用的 表
	sql_create_table = `
		CREATE TABLE IF NOT EXISTS nav_table (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			url TEXT,
			logo TEXT,
			catelog TEXT,
			desc TEXT
		);
		`
	_, err = db.Exec(sql_create_table)
	checkErr(err)
	// 分类表
	sql_create_table = `
		CREATE TABLE IF NOT EXISTS nav_catelog (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT
		);
		`
	_, err = db.Exec(sql_create_table)
	checkErr(err)
	// api token 表
	sql_create_table = `
		CREATE TABLE IF NOT EXISTS nav_api_token (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			value TEXT,
			disabled INTEGER
		);
		`
	_, err = db.Exec(sql_create_table)
	checkErr(err)
	// 如果不存在，就初始化用户
	sql_get_user := `
		SELECT * FROM nav_user;
		`
	rows, err := db.Query(sql_get_user)
	checkErr(err)
	if !rows.Next() {
		sql_add_user := `
			INSERT INTO nav_user (id, name, password)
			VALUES (?, ?, ?);
			`
		stmt, err := db.Prepare(sql_add_user)
		checkErr(err)
		res, err := stmt.Exec(generateId(), "admin", "admin")
		checkErr(err)
		_, err = res.LastInsertId()
		checkErr(err)
		// fmt.Println(id)
	}
	rows.Close()
	// 如果不存在设置，就初始化
	sql_get_setting := `
		SELECT * FROM nav_setting;
		`
	rows, err = db.Query(sql_get_setting)
	checkErr(err)
	if !rows.Next() {
		sql_add_setting := `
			INSERT INTO nav_setting (id, favicon, title)
			VALUES (?, ?, ?);
			`
		stmt, err := db.Prepare(sql_add_setting)
		checkErr(err)
		res, err := stmt.Exec(0, "https://pic.mereith.com/img/male.svg", "Van Nav")
		checkErr(err)
		_, err = res.LastInsertId()
		checkErr(err)
		// fmt.Println(id)
	}
	rows.Close()
	fmt.Println("数据库初始化成功。。。")
}

//go:embed public
var fs embed.FS

type binaryFileSystem struct {
	fs   http.FileSystem
	root string
}

func (b *binaryFileSystem) Open(name string) (http.File, error) {
	// fmt.Println("打开文件",name)
	openPath := path.Join(b.root, name)
	return b.fs.Open(openPath)
}

func (b *binaryFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		var name string
		if p == "" {
			// fmt.Println("找 index")
			name = path.Join(b.root, p, INDEX)
		} else {
			name = path.Join(b.root, p)
		}
		// 判断
		// fmt.Println("文件是否存在？",name)
		if _, err := b.fs.Open(name); err != nil {
			return false
		}
		return true
	}
	return false
}
func BinaryFileSystem(data embed.FS, root string) *binaryFileSystem {
	fs := http.FS(data)
	return &binaryFileSystem{
		fs,
		root,
	}
}
func main() {
	initDB()
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// 嵌入文件夹

	// public,_ := fs.ReadDir("./public")
	// router.StaticFS("/",http.FS(fs))

	router.Use(static.Serve("/", BinaryFileSystem(fs, "public")))
	// router.Use(static.Serve("/", static.LocalFile("./public", true)))
	api := router.Group("/api")
	{
		// 获取数据的路由
		api.GET("/", GetAllHandler)
		// 获取用户信息

		api.POST("/login", LoginHandler)
		api.GET("/logout", LogoutHandler)

		// 管理员用的
		admin := api.Group("/admin")
		admin.Use(JWTMiddleware())
		{
			admin.POST("/apiToken", AddApiTokenHandler)
			admin.DELETE("/apiToken/:id", DeleteApiTokenHandler)
			admin.GET("/all", GetAdminAllDataHandler)

			admin.GET("/exportTools", ExportToolsHandler)

			admin.POST("/importTools", ImportToolsHandler)

			admin.PUT("/user", UpdateUserHandler)

			admin.PUT("/setting", UpdateSettingHandler)

			admin.POST("/tool", AddToolHandler)
			admin.DELETE("/tool/:id", DeleteToolHandler)
			admin.PUT("/tool/:id", UpdateToolHandler)

			admin.POST("/catelog", AddCatelogHandler)
			admin.DELETE("/catelog/:id", DeleteCatelogHandler)
			admin.PUT("/catelog/:id", UpdateCatelogHandler)
		}
	}
	fmt.Println("应用启动成功，网址:   http://localhost:6412")
	router.Run(":6412")
}

func importTools(data []Tool) {
	var catelogs []string
	for _, v := range data {
		// if ()
		if !in(v.Catelog, catelogs) {
			catelogs = append(catelogs, v.Catelog)
		}
		sql_add_tool := `
			INSERT INTO nav_table (id, name, catelog, url, logo, desc)
			VALUES (?, ?, ?, ?, ?, ?);
			`
		stmt, err := db.Prepare(sql_add_tool)
		checkErr(err)
		res, err := stmt.Exec(v.Id, v.Name, v.Catelog, v.Url, v.Logo, v.Desc)
		checkErr(err)
		_, err = res.LastInsertId()
		checkErr(err)
	}
	for _, catelog := range catelogs {
		var addCatelogDto addCatelogDto
		addCatelogDto.Name = catelog
		addCatelog(addCatelogDto, db)
	}
}

func getSetting(db *sql.DB) Setting {
	sql_get_user := `
		SELECT * FROM nav_setting WHERE id = ?;
		`
	var setting Setting
	row := db.QueryRow(sql_get_user, 0)
	err := row.Scan(&setting.Id, &setting.Favicon, &setting.Title)
	checkErr(err)
	return setting
}

func getApiTokens(db *sql.DB) []Token {
	sql_get_api_tokens := `
		SELECT * FROM nav_api_token WHERE disabled = 0;
		`
	results := make([]Token, 0)
	rows, err := db.Query(sql_get_api_tokens)
	checkErr(err)
	for rows.Next() {
		var token Token
		err = rows.Scan(&token.Id, &token.Name, &token.Value, &token.Disabled)
		checkErr(err)
		results = append(results, token)
	}
	defer rows.Close()
	return results
}

func getUser(name string, db *sql.DB) User {
	sql_get_user := `
		SELECT * FROM nav_user WHERE name = ?;
		`
	var user User
	row := db.QueryRow(sql_get_user, name)
	err := row.Scan(&user.Id, &user.Name, &user.Password)
	checkErr(err)
	return user
}
