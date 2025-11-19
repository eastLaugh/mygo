# 我的GO中心仓库

# MYGO
🚀 使用简短别名安装 Go 包（如 `gorm` 而不是 `gorm.io/gorm`）

🌐 包映射配置托管在 GitHub，无需本地配置文件

## 安装

### 使用 go install

```bash
go install github.com/eastLaugh/mygo/cmd/my@latest
```

### 克隆

```bash
git clone https://github.com/eastLaugh/mygo.git
cd mygo
make install
```

## 直接使用

```bash
# 安装单个包
my go gorm

# 安装多个包
my go gorm gin fiber

# 只显示命令，不执行（预览模式）
my go -n gorm
my go --no gorm gin
```

### 命令格式

```
my go [-n|--no] <package-name>...
```

## 包映射配置

包映射配置存储在 GitHub 仓库的 `my.toml` 文件中，以简单 KV 格式存储。

```toml
gorm = "gorm.io/gorm"
gin = "github.com/gin-gonic/gin"
fiber = "github.com/gofiber/fiber/v2"
```

## 项目结构

```
mygo/
├── cmd/
│   └── my/
│       └── main.go      # 主程序
├── my.toml              # 包映射配置
├── Makefile             
├── go.mod               
└── README.md            
```
