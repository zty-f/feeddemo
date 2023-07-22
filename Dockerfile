FROM public-env-mirror-service-registry.cn-beijing.cr.aliyuncs.com/dist/golang:1.16

WORKDIR /tmp/source
# 准备工作
#RUN export
COPY go.mod .
COPY go.sum .
COPY . .

#加入git访问权限
#COPY .netrc /root/.netrc

# 编译
RUN go env -w GOPRIVATE=codeup.aliyun.com
RUN go env
#WORKDIR /tmp/source/cmd/server
RUN GOPROXY="https://goproxy.cn" GO111MODULE=on go build -o /tmp/source/bin/server .
RUN chmod +x /tmp/source/bin/server

WORKDIR /tmp/source
#ARG envType=test
#COPY configs/config_${envType}.toml conf/env/env.toml
# 执行编译生成的二进制文件
CMD ["/tmp/source/bin/server","-conf","conf/env/env.toml"]
# 暴露端口
EXPOSE 80