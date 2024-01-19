# 使用 Go 1.16 作为基础镜像
FROM golang:1.21.6

# 安装 supervisord
RUN apt-get update && apt-get install -y supervisor

# 设置工作目录
WORKDIR /app

# 将 go.mod 和 go.sum 文件复制到工作目录
COPY go.mod go.sum ./

# 下载所有依赖项
RUN go mod download

# 将源代码复制到工作目录
COPY . .

# 将 conf/config.yml 配置文件复制到工作目录
COPY conf/config.yml ./conf/

# 编译程序
RUN go build -o main .

# 将 supervisord 配置文件复制到 Docker 容器
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# 暴露端口
EXPOSE 8080

# 使用 supervisord 启动应用程序
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]