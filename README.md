# 1. 打包编译
```
docker build -t lvzhancheng/golang_study:1.3.0 .
```
# 2. 推送到dockerhub
```
docker login
输入账号密码,提示successful表示登录成功
docker push lvzhancheng/golang_study:1.3.0 
```
# 3. 本地启动
```
docker run -d -p 80:80 --name http_server lvzhancheng/golang_study:1.3.0
```
# 4. 查看容器ip
+ 获取容器pid
```
Pid=`docker inspect -f {{.State.Pid}} http_server`
```
+ 进入容器网络命名空间
```
nsenter -n -t $Pid
```
+ 查看容器ip
```
ip addr
```

# 5. 模块八第一部分作业
现在你对 Kubernetes 的控制面板的工作机制是否有了深入的了解呢？

是否对如何构建一个优雅的云上应用有了深刻的认识，那么接下来用最近学过的知识把你之前编写的 http 以优雅的方式部署起来吧，你可能需要审视之前代码是否能满足优雅上云的需求。

作业要求：编写 Kubernetes 部署脚本将 httpserver 部署到 Kubernetes 集群，以下是你可以思考的维度。

+ 优雅启动
+ 优雅终止
+ 资源需求和 QoS 保证
+ 探活
+ 日常运维需求，日志等级
+ 配置和代码分离