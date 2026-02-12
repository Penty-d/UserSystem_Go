# Go 用户管理系统

[![Go Version](https://img.shields.io/badge/go-1.25.5-blue)](https://golang.org/)
[![Gin Framework](https://img.shields.io/badge/gin-1.11.0-blue)](https://gin-gonic.com/)

基于Golang和MySQL构建的RESTful用户管理系统，提供完整的用户认证和权限管理功能。

## 功能特性

- ✅ 用户注册与登录
- ✅ Token令牌认证
- ✅ 密码加密存储（bcrypt）
- ✅ 基于角色的访问控制（admin/user）
- ✅ 用户信息管理
- ✅ 自动数据库初始化
- ✅ 请求日志记录
- ✅ 服务端错误恢复

## 技术栈

- **语言**: Go 1.25.5
- **框架**: Gin 1.11.0
- **数据库**: MySQL
- **安全**: 
  - Bcrypt密码哈希
  - Token令牌认证
- **工具**: 
  - go-sql-driver/mysql
  - golang.org/x/crypto

## 快速开始

### 前置需求
- Go 1.25+
- MySQL 5.7+

### 安装步骤

# 克隆仓库
git clone https://github.com/yourusername/user-system-go.git

# 进入目录
cd user-system-go

# 安装依赖
go mod download


### 配置环境
复制.env.example并配置数据库：
cp config/.env.example config/.env

环境变量配置示例：

DB_USER=root
DB_PASSWORD=123456
DB_HOST=localhost
DB_PORT=3306
DB_NAME=usersystem


### 启动服务

go run main.go


## API文档

### 公开端点
| 方法 | 路径         | 描述       |
|------|--------------|------------|
| POST | /api/register | 用户注册   |
| POST | /api/login    | 用户登录   |

### 受保护端点
| 方法 | 路径               | 描述         |
|------|--------------------|--------------|
| POST   | /api/delete         | 删除用户     |
| POST   | /api/change_password| 修改密码     |
| GET    | /api/users          | 获取用户信息 |

**认证要求**：在Authorization Header中添加Bearer Token

## 项目结构

usersystem_go/
├── config/            # 配置管理
├── database/          # 数据库连接
├── middleware/        # 中间件
│   └── middleware.go  # 认证/日志/恢复中间件
├── models/            # 数据模型
├── repositories/      # 数据访问层
├── userhandler/       # 控制器
├── utils/             # 工具函数
│   ├── auth.go        # JWT认证
│   └── password.go    # 密码加密
├── go.mod
└── main.go            # 入口文件

## 许可证
MIT License