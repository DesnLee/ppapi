-- 给 pgsql 的 users 表添加一个 phone 字段
ALTER TABLE users ADD COLUMN phone VARCHAR(50) NOT NULL UNIQUE DEFAULT '';
