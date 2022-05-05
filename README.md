# 1. 打包编译
```
docker build -t lvzhancheng/golang_study:1.0.0 .
```
# 2. 推送到dockerhub
```docker login
docker push lvzhancheng/golang_study:1.0.0 
```
# 3. 本地启动
```
docker run -d -p 80:80 --name http_server lvzhancheng/golang_study:1.0.0
```
# 4. 查看容器ip
+ 获取容器pid
```
docker inspect -f {{.State.Pid}} http_server
```
+ 进入容器网络命名空间
```
nsenter -n -t Pid
```
+ 查看容器ip
```
ip addr
```