# badExtraction

一个用于清理错误解压文件的Go工具。

## 为什么需要这个项目？

在日常使用中，我们经常会遇到这样的情况：在一个包含大量文件的目录中解压了一个压缩包，结果压缩包内的文件和文件夹散落在当前目录中，与原本的文件混在一起。这时候很难区分哪些是原有的文件，哪些是刚解压出来的文件，手动清理非常麻烦且容易出错。

`cleanup` 工具通过分析压缩包的第一级内容，自动识别当前目录中可能来自该压缩包的文件和文件夹，并提供安全的删除功能。

## 功能特性

- **多格式支持**: 支持 ZIP、TAR、TAR.GZ、TAR.BZ2、TAR.XZ 等常见压缩格式
- **智能检测**: 自动分析压缩包第一级内容，匹配当前目录中的对应文件/文件夹
- **安全确认**: 删除前会列出所有待删除项目，需要用户确认后才执行
- **错误处理**: 完善的错误提示和异常处理机制

## 安装方法

### 使用 Makefile 安装

```bash
# 克隆项目
git clone <repository-url>
cd badExtraction

# 编译并安装到 $HOME/bin
make install

# 或者仅编译
make build

# 清理编译文件
make clean

# 卸载程序
make uninstall
```

### 手动编译安装

```bash
# 编译
go build -o cleanup cleanup.go

# 手动安装到系统路径
sudo cp cleanup /usr/local/bin/
# 或安装到用户路径
cp cleanup $HOME/bin/
```

## 使用方法

```bash
cleanup <压缩包文件名>
```

### 使用示例

```bash
# 处理 ZIP 文件
cleanup archive.zip

# 处理 TAR.XZ 文件
cleanup software-1.0.tar.xz

# 处理 TAR.GZ 文件
cleanup project.tar.gz
```

### 使用流程

1. 在包含错误解压文件的目录中运行 `cleanup` 命令
2. 程序会分析指定压缩包的第一级内容
3. 显示当前目录中匹配的文件和文件夹列表
4. 确认删除操作 (输入 `y` 或 `yes` 确认，其他任何输入取消)
5. 程序执行删除并显示结果

### 示例输出

```
$ cleanup example.tar.xz
发现以下文件/文件夹可能来自压缩包解压:
  - src/
  - README.md
  - Makefile
确认删除这些文件和文件夹吗? (y/N): y
已删除: src/
已删除: README.md
已删除: Makefile
清理完成
```

## 支持的压缩格式

- `.zip` - ZIP 压缩文件
- `.tar` - TAR 归档文件
- `.tar.gz` / `.tgz` - TAR + GZIP 压缩文件
- `.tar.bz2` - TAR + BZIP2 压缩文件
- `.tar.xz` - TAR + XZ 压缩文件

## 安全说明

- 程序只会删除与压缩包第一级内容匹配的文件和文件夹
- 删除前会显示完整列表并要求用户确认
- 如果压缩包不存在或格式不支持，程序会给出明确提示
- 删除操作不可逆，请谨慎确认

## 系统要求

- Go 1.16 或更高版本 (编译时)
- Linux/macOS/Windows 系统
- 对于 `.tar.xz` 格式，需要系统安装 `xz` 工具

## 开发

```bash
# 获取帮助
make help

# 编译
make build

# 安装
make install

# 清理
make clean
```

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。