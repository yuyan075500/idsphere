# Docker Compose 部署
此方式适用于测试、演示环境，对性能、稳定性和安全性没有任何求的。
1. **部署环境准备**：<br><br>
   你需要准备一台 Linux 服务器，并安装以下组件。
   * [x] Docker。
   * [x] Docker Compose。
     `Docker`和`Docker Compose`是部署必须准备的，其它组件在 `docker-compose.yaml` 配置清单中已指定。<br><br>
2. **克隆项目**：<br><br>
    ```shell
    git clone https://github.com/yuyan075500/idsphere.git
    或
    git clone https://gitee.com/yybluestorm/idsphere
    ```
3. **切换工作目录**：<br><br>
    ```shell
    cd idsphere/deploy/docker-compose
    ```
4. **配置环境变量**：<br><br>
   配置文件位于 `.env`，此配置文件中主要指定了 MySQL 数据库、Redis 缓存、MinIO 的初始化配置和项目启动的版本，该步骤可以跳过。<br><br>
5. **修改项目配置**：<br><br>
   配置文件位于 `conf/config.yaml`，修改方法参考 [配置文件说明](#配置文件说明)，注意事项如下：
   * `oss.accessKey` 和 `oss.secretKey` 中指定的 `AK` 和 `SK` 需要在 Minio 启动完成后登录到 Minio 后台手动创建。
   * `oss.endpoint` 配置的地址必须确保使用 IDSphere 统一认证平台的客户端电脑可以访问，如果实际的地址协议为 `HTTPS` 则需要将 `oss.ssl` 更改为 `true`。<br><br>
6. **创建 Minio 数据目录**：<br><br>
   需要手动创建 Minio 数据目录，并更改权限为 `1001:1001`。
   ```shell
   mkdir -p data/minio
   chown -R 1001:1001 data/minio
   ```
7. **执行部署**：<br><br>
    ```shell
    docker-compose up -d
    ```
8. **系统登录**：<br><br>
   系统会自动创建一个超级用户，此用户不受系统权限控制，默认用户名为：`admin`，密码为：`admin@123...`。<br><br>
9. **密码更改**：<br><br>
   为确保系统安全请务必更改 `admin` 账号的初始密码。
# Kubernetes 部署
生产环境推荐使用此种部署方法，你需要准备以下相关资源：
* [x] [Kubernetes](https://kubernetes.io "Kubernetes") 软件运行必要环境。
* [x] [Helm](https://helm.sh "Helm") 部署客户端工具，此工具需要能访问到 Kubernetes 集群。
* [x] MySQL 8.0。
* [x] Redis 5.x。
* [x] MinIO 或华为云 OBS 对象存储。<br><br>
  **注意**：需要确保 Minio API 地址使用的 `scheme` 和 IDSphere 统一认证平台使用的 `scheme` 一致，否则有可能导致上传到 Minio 的图片无法展示。
1. **克隆项目**：<br><br>
    ```shell
    git clone https://github.com/yuyan075500/idsphere.git
    或
    git clone https://gitee.com/yybluestorm/idsphere
    ```
2. **切换工作目录**：<br><br>
    ```shell
    cd idsphere/deploy/kubernetes
    ```
4. **修改项目配置**：<br><br>
   配置文件位于 `templates/configmap.yaml`，修改方法参考 [配置文件说明](#配置文件说明)。<br><br>
5. **部署**：<br><br>
   ```shell
   helm install <APP_NAME> --namespace <NAMESPACE_NAME> .
   ```
6. **系统登录**：<br><br>
   系统会自动创建一个超级用户，此用户不受系统权限控制，默认用户名为：`admin`，密码为：`admin@123...`。<br><br>
7. **密码更改**：<br><br>
   为确保系统安全请务必更改 `admin` 账号的初始密码。
# 配置文件说明
```yaml
server: "0.0.0.0:8000"
mysql:
   host: "127.0.0.1"
   port: 3306
   db: "ops"
   user: "root"
   password: ""
   maxIdleConns: 10
   maxOpenConns: 100
   maxLifeTime: 30
redis:
   host: "127.0.0.1:6379"
   password: ""
   db: 0
oss:
   endpoint: ""
   accessKey: ""
   secretKey: ""
   bucketName: ""
   ssl: true
```
* [x] server：后端服务监听的地址和端口。
* [x] mysql：`MySQL` 数据库相关配置。
* [x] redis：`Redis` 相关配置。
* [x] oss：对象存储相关配置，支持 MinIO 和华为云 OBS。<br><br>