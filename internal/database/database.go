package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var (
	DB       *sql.DB
	host     = "localhost"
	user     = "pp"
	password = "ljk950805"
	dbname   = "pp_dev"
	port     = "5432"
)

func Connect() {
	connStr := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Shanghai", host, user, password, dbname, port)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}
	DB = db
	if err = db.Ping(); err != nil {
		log.Fatalln(err)
	}
	log.Println("数据库连接成功！")
}

func Close() {
	if err := DB.Close(); err != nil {
		log.Fatalln(err)
	}
	log.Println("数据库连接已关闭！")
}

func CreateTables() {
	// 创建 Users 表
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS users (  
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        age INT NOT NULL,
        email VARCHAR(255) NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    )`)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("创建 users 表成功！")
}
