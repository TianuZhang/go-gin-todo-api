# Todo API 学习项目

一个用 Go 和 Gin 编写的 RESTful API 练习项目。当前阶段使用内存保存用户数据；后续将接入 MySQL 和 Redis。

## 运行

```powershell
go mod tidy
go run .
```

服务默认监听 `http://localhost:8080`。

## 接口

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/health` | 服务健康检查 |
| GET | `/users` | 查询用户列表 |
| GET | `/users/:id` | 查询一个用户 |
| POST | `/users` | 创建用户 |
| PATCH | `/users/:id` | 部分修改用户 |
| DELETE | `/users/:id` | 删除用户 |

创建用户请求示例：

```json
{
  "name": "Charlie",
  "email": "charlie@example.com"
}
```

## 当前限制

数据只保存在内存中，因此服务重启后会恢复为初始数据。接入 MySQL 后，用户数据将持久化保存。
