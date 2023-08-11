package database

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB       *gorm.DB
	host     = "localhost"
	user     = "root"
	password = "123456"
	dbname   = "pp_dev"
	port     = "5432"
)

// sqlc
func Connect() {
	// dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	// db, err := sql.Open("postgres", dsn)

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Shanghai", host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	} else {
		DB = db
		log.Println("数据库连接成功！")
	}
}

func Close() {
	if sqlDB, err := DB.DB(); err != nil {
		log.Fatalln(err)
	} else {
		err = sqlDB.Close()
		if err != nil {
			log.Println("E: 数据库连接关闭失败！")
			log.Fatalln(err)
		} else {
			log.Println("数据库连接已关闭！")
		}
	}
}

type User struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email     string
	Phone     string `gorm:"varchar(20);uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func CreateTables() {
	if err := DB.Migrator().CreateTable(&User{}); err != nil {
		log.Fatalln(err)
	} else {
		log.Println("创建表成功！")
	}
}

func Crud() {
	user := User{Email: "test1@qq.com"}
	tx := DB.Create(&user)
	if tx.Error != nil {
		log.Fatalln(tx.Error)
	} else {
		log.Println("创建成功！")
	}
	fmt.Println(user.ID)
	user.Phone = "123456789"
	tx = DB.Save(&user)
	if tx.Error != nil {
		log.Fatalln(tx.Error)
	} else {
		log.Println("更新成功！")
	}
}

func Migrate() {
	if err := DB.AutoMigrate(&User{}); err != nil {
		log.Fatalln(err)
	} else {
		log.Println("迁移表成功！")
	}
}

// gorm
// func Connect() {
// 	// dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
// 	// db, err := sql.Open("postgres", dsn)
//
// 	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Shanghai", host, user, password, dbname, port)
// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
//
// 	if err != nil {
// 		log.Fatalln(err)
// 	} else {
// 		DB = db
// 		log.Println("数据库连接成功！")
// 	}
// }
//
// func Close() {
// 	if sqlDB, err := DB.DB(); err != nil {
// 		log.Fatalln(err)
// 	} else {
// 		err = sqlDB.Close()
// 		if err != nil {
// 			log.Println("E: 数据库连接关闭失败！")
// 			log.Fatalln(err)
// 		} else {
// 			log.Println("数据库连接已关闭！")
// 		}
// 	}
// }
//
// type User struct {
// 	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
// 	Email     string
// 	Phone     string `gorm:"varchar(20);uniqueIndex"`
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// }
//
// func CreateTables() {
// 	if err := DB.Migrator().CreateTable(&User{}); err != nil {
// 		log.Fatalln(err)
// 	} else {
// 		log.Println("创建表成功！")
// 	}
// }
//
// func Crud() {
// 	user := User{Email: "test1@qq.com"}
// 	tx := DB.Create(&user)
// 	if tx.Error != nil {
// 		log.Fatalln(tx.Error)
// 	} else {
// 		log.Println("创建成功！")
// 	}
// 	fmt.Println(user.ID)
// 	user.Phone = "123456789"
// 	tx = DB.Save(&user)
// 	if tx.Error != nil {
// 		log.Fatalln(tx.Error)
// 	} else {
// 		log.Println("更新成功！")
// 	}
// }
//
// func Migrate() {
// 	if err := DB.AutoMigrate(&User{}); err != nil {
// 		log.Fatalln(err)
// 	} else {
// 		log.Println("迁移表成功！")
// 	}
// }
