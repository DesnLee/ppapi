package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os/exec"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"ppapi.desnlee.com/db/sqlcExec"
)

var (
	DB    *sql.DB
	DBCtx = context.Background()

	host     = "localhost"
	user     = "root"
	password = "123456"
	dbname   = "pp_dev"
	port     = "5432"
)

// sqlc
func Connect() {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Shanghai", host, user, password, dbname, port)
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		log.Fatalln(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalln(err)
	}
	DB = db
	log.Println("数据库连接成功！")
}

func Close() {
	if err := DB.Close(); err != nil {
		log.Fatalln(err)
	} else {
		log.Println("数据库连接已关闭！")
	}
}

func MigrateNew(name string) {
	cmd := exec.Command("migrate", "create", "-ext", "sql", "-dir", "db/migrations", "-seq", name)
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}
}

func MigrateUp() {
	sourceUrl := fmt.Sprintf("file://db/migrations")
	dbUrl := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", user, password, host, port, dbname)
	m, err := migrate.New(sourceUrl, dbUrl)
	if err != nil {
		log.Fatalln(err)
	}

	err = m.Up()
	if err != nil {
		log.Fatalln(err)
	}
}

func MigrateDown(step int) {
	sourceUrl := fmt.Sprintf("file://db/migrations")
	dbUrl := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", user, password, host, port, dbname)
	m, err := migrate.New(sourceUrl, dbUrl)
	if err != nil {
		log.Fatalln(err)
	}

	err = m.Steps(step * -1)
	if err != nil {
		log.Fatalln(err)
	}
}

func Crud() {
	q := sqlcExec.New(DB)

	// num := rand.Int()
	// u, err := q.CreateUser(DBCtx, sqlcExec.CreateUserParams{Email: fmt.Sprintf("%v@qq.com", num), Phone: fmt.Sprintf("%v", num)})
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	id := uuid.MustParse("042ef509-32c6-492f-93b3-faea02dac1f1")
	newU, err := q.UpdateUser(DBCtx, sqlcExec.UpdateUserParams{ID: id, Phone: "123456789"})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(newU)
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
