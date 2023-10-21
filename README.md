# 启动本地数据库

## psql

```bash
docker run -d --name pg-for-go-mangosteen -e POSTGRES_USER=mangosteen -e POSTGRES_PASSWORD=123456 -e POSTGRES_DB=mangosteen_dev -e PGDATA=/var/lib/postgresql/data/pgdata -v pg-go-mangosteen-data:/var/lib/postgresql/data --network=network1 postgres:14
```

## mysql

```bash
docker run -d --network=network1 --name mysql-for-go-mangosteen -e MYSQL_DATABASE=mangosteen_dev -e MYSQL_USER=mangosteen -e MYSQL_PASSWORD=123456 -e MYSQL_ROOT_PASSWORD=123456 -v mysql-go-mangosteen-data:/var/lib/mysql mysql:8 --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
```
# 数据库迁移

## 安装工具

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## 创建迁移文件

```bash
go build . && ./ppapi db migrate:new migration_name
# 或者
migrate create -ext sql -dir db/migrations -seq migration_name
```
## 运行迁移文件

```bash
go build . && ./ppapi db migrate:up
# 或者
migrate -database "postgres://root:123456@localhost:5432/pp_dev?sslmode=disable" -source "file://$(pwd)/db/migrations" up
```

## 回滚迁移文件

```bash
go build . && ./ppapi db migrate:down 1
# 或者
migrate -database "postgres://root:123456@localhost:5432/pp_dev?sslmode=disable" -source "file://$(pwd)/db/migrations" down 1
```

# 启动服务

```bash
sh ./scripts/run.sh
```
