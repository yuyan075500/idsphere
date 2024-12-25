# 项目部署
支持使用 [Docker Compose一键部署](#docker-compose部署) 和 [Kubernetes部署](#Kubernetes部署)，生产环境推荐使用 Kubernetes 部署。
## Docker Compose部署
如果你想快速拥有一个简易的环境用于测试、演示，且对性能、稳定性和安全性没有任何求的，那么推荐使用该部署方式。  
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
   配置文件位于 `conf/config.yaml`，修改方法参考 [配置文件说明](#配置文件说明)，以下配置必修改项：
   * `externalUrl` 需要更改为 IDSphere 统一认证平台在浏览器实际的访问地址，否则导致单点功能等相关功能无法正常使用。
   * `oss.accessKey` 和 `oss.secretKey` 中指定的 `AK` 和 `SK` 需要在 Minio 启动完成后登录到后台手动创建。
   * `oss.endpoint` 配置的地址必须确保使用 IDSphere 统一认证平台的客户端电脑可以访问，如果实际的地址协议为 `HTTPS` 则需要将 `oss.ssl` 更改为 `true`。<br><br>
6. **创建证书**：<br><br>
   参考 [创建项目证书](#项目证书)，将生成的新证书保存至`certs`目录中并覆盖目标文件，测试环境可以跳过此步骤。<br><br>
7. **创建 Minio 数据目录**：<br><br>
   需要手动创建 Minio 数据目录，并更改权限为 `1001:1001`。
   ```shell
   mkdir -p data/minio
   chown -R 1001:1001 data/minio
   ```
8. **执行部署**：<br><br>
    ```shell
    docker-compose up -d
    ```
9. **系统登录**：<br><br>
   部署完成后，会自动创建一个超级用户，此用户不受 Casbin 权限控制，默认用户名为：`admin`，密码为：`admin@123...`。<br><br>
10. **密码更改**：<br><br>
   为确保系统安全请务必更改 `admin` 账号的初始密码。
## Kubernetes部署
生产环境推荐使用此种部署方法，你需要准备以下相关资源：
* [x] [Kubernetes](https://kubernetes.io "Kubernetes") 软件运行必要环境。
* [x] [Helm](https://helm.sh "Helm") 部署客户端工具，此工具需要能访问到 Kubernetes 集群。
* [x] MySQL 8.0。
* [x] Redis 5.x。
* [x] MinIO 或华为云 OBS 对象存储。<br><br>
**注意**：需要确保 Minio API 地址使用的 `scheme` 和 IDSphere 统一认证平台使用的 `scheme` 一致，否则有可能导致上传到 Minio 的图片无法展示。
### 部署
1. **克隆项目**：
    ```shell
    git clone https://github.com/yuyan075500/idsphere.git
    或
    git clone https://gitee.com/yybluestorm/idsphere
    ```
2. **切换工作目录**：
    ```shell
    cd idsphere/deploy/kubernetes
    ```
3. **创建证书**：<br><br>
   创建 [项目证书](#项目证书)，证书创建完成后需要使用新的证书替换 `templates/configmap.yaml` 文件中对应的内容。<br><br>
4. **修改项目配置**：<br><br>
   配置文件位于 `templates/configmap.yaml`，修改方法参考 [配置文件说明](#配置文件说明)。<br><br>
5. **部署**：
   ```shell
   helm install <APP_NAME> --namespace <NAMESPACE_NAME> .
   ```
6. **系统登录**：<br><br>
   部署完成后，会自动创建一个超级用户，此用户不受 Casbin 权限控制，默认用户名为：`admin`，密码为：`admin@123...`。<br><br>

   **PS：如果需要高可以自行调整应用的副本数。**

# 配置文件说明
```yaml
server: "0.0.0.0:8000"
externalUrl: ""
secret: "swfqezjzoqssvjck"
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
jwt:
  secret: "swfqezjzoqssvjck"
  expires: 6
mfa:
  enable: false
  issuer: "IDSphere 统一认证中心"
oss:
  endpoint: ""
  accessKey: ""
  secretKey: ""
  bucketName: ""
  ssl: true
ldap:
  host: ""
  bindUserDN: ""
  bindUserPassword: ""
  searchDN: ""
  userAttribute: ""
  maxPasswordAge: 90
sms:
  provider: ""
  url: ""
  appKey: ""
  appSecret: ""
  callbackUrl: ""
  resetPassword:
    sender: ""
    templateId: ""
    signature: ""
mail:
  smtpHost: ""
  smtpPort: 587
  from: ""
  password: ""
dingTalk:
   appKey: ""
   appSecret: ""
wechat:
   corpId: ""
   agentId: ""
   secret: ""
feishu:
   appId: ""
   appSecret: ""
swagger: true
```
* [x] server：后端服务监听的地址和端口。
* [x] externalUrl：IDSphere 统一认证平台对外提供的访问地址，格式为：`http[s]://<address>[:<port>]`。
* [x] secret: 签名字符串，生产环境请务必修改。
* [x] mysql：`MySQL` 数据库相关配置。
* [x] redis：`Redis` 相关配置。
* [x] jwt：`JWT` 相关配置。
* [x] mfa：双因素认证相关配置，`issuer` 是 APP 扫码后显示的名称。
* [x] oss：对象存储相关配置，支持 MinIO 和华为云 OBS。
* [ ] ldap：参考 [LDAP配置](#LDAP配置)。
* [ ] sms：参考 [短信配置](#短信配置)。
* [ ] mail：邮件相关配置。
* [ ] dingTalk：钉钉扫码登录相关配置，参考 [钉钉配置](https://github.com/yuyan075500/idsphere/blob/main/deploy/dingtalk.md "钉钉配置")。
* [ ] wechat：企业微信扫码登录相关配置，参考 [企业微信配置](https://github.com/yuyan075500/idsphere/blob/main/deploy/wechat.md "企业微信配置")。
* [ ] feishu：飞书扫码登录相关配置，参考 [飞书配置](https://github.com/yuyan075500/idsphere/blob/main/deploy/feishu.md "飞书配置")。
* [x] swagger：Swagger 接口配置，生产环境建议关闭。<br><br>

**注意：`externalUrl` 地址一经固定，切忽随意更改，更改后影响 SSO 的相关功能，如果更改后 SSO 客户端无法登录，那么你需要重置进行相关客户端配置。**
## LDAP配置
该配置项支持与 Windows AD 或 OpenLDAP 进行对接，实现用户认证，使配置说明如下：
* [x] host：服务器地址，格式为：`ldap[s]://<host>:<port>`。
* [x] bindUserDN：绑定的用户DN，格式为：`cn=admin,dc=idsphere,dc=cn`。
* [x] bindUserPassword：绑定的用户密码。
* [x] searchDN：允许登录用户的范围，格式为：`ou=IT,dc=idsphere,dc=cn`，支持配置多个，之间使用 `&` 分割，如：`ou=IT,dc=idsphere,dc=cn&ou=HR,dc=idsphere,dc=cn`。
* [x] userAttribute：用户属性，如果是 OpenLDAP 则为 `uid`，如果是 Windows AD 则为 `sAMAccountName`。
* [ ] maxPasswordAge：密码最大有效期，此参数仅针对 `Windows AD`，需要与实际的域控用户密码有效期保持一致。<br><br>

配置完成后还需要将用户同步到本地数据库，否则用户无法登录。可以在系统中手动执行【用户同步】或通过【定时任务】功能创建自动同步任务。<br><br>
**注意：如果需要更改 Windows AD 或 OpenLDAP 的用户密码功能，需要绑定的用户有足够的权限。如果是 Windows AD 还要求使用 `ldaps` 协议进行连接，`ldaps` 协议的默认端口为 `636`。**
## 短信配置
如果有使用到短信功能则需要配置，如果没有则无需要配置,支持使用华为云和阿里云短信，具体配置如下所示：
* [x] provider：指定短信服务商，固定值，`aliyun` 或 `huawei`。
* [x] url：短信服务地址，不同服务商的地址不同，阿里云参考 [短信服务接入点](https://help.aliyun.com/zh/sms/developer-reference/api-dysmsapi-2017-05-25-endpoint "阿里云短信服务接入点")，华为云参考 [API请求地址](https://support.huaweicloud.com/api-msgsms/sms_05_0000.html#section1 "API请求地址")。
* [x] appKey: 华为云参考 [开发数据准备](https://support.huaweicloud.com/devg-msgsms/sms_04_0006.html "开发数据准备")，阿里云参考 [创建AccessKey](https://help.aliyun.com/zh/ram/user-guide/create-an-accesskey-pair "创建AccessKey")。
* [x] appSecret: 华为云参考 [开发数据准备](https://support.huaweicloud.com/devg-msgsms/sms_04_0006.html "开发数据准备")，阿里云参考 [创建AccessKey](https://help.aliyun.com/zh/ram/user-guide/create-an-accesskey-pair "创建AccessKey")。
* [ ] callbackUrl：短信回调地址，用于接收短信发送状态，仅华为云需要配置，回调地址为 `<externalUrl>/api/v1/sms/huawei/callback`，请确保该地址公网可以访问到。
* [ ] resetPassword.sender：短信通道号，仅华为云需要配置。
* [x] resetPassword.templateId：短信模板ID。
* [x] resetPassword.signature：短信签名名称。
### 短信模板
阿里云模板如下所示：
```
# 短信模板
您的验证码为：${code}，验证码在5分钟内有效，请勿泄漏他人！
```
华为云模板如下所示：
```
# 短信模板
您的验证码为：${1}，验证码在5分钟内有效，请勿泄漏他人！
```
# 项目证书
为确保重要信息不会泄露，在项目部署时建议生成一套全新的证书，推荐使用 [证书在线生成工具](https://www.qvdv.net/tools/qvdv-csrpfx.html "在线生成工具") 创建。建议将证书有效期设置为10年，证书生成完成后需要下载CRT证书文件、证书公钥和证书私钥并严格按以下名称命名：
* private.key：私钥
* public.key：公钥
* certificate.crt：证书  

你也可以使用`openssl`工具生成自签证书，参考命令如下所示：
```shell
# 生成私钥
openssl genpkey -algorithm RSA -out private.key -pkeyopt rsa_keygen_bits:2048 -outform PEM
# 创建证书
openssl req -new -x509 -key private.key -out certificate.crt -days 3650
# 从证书中提取公钥
openssl rsa -in private.key -pubout -out public.key
```
# IP地址库
记录用户登录信息中的源IP来源于离线库文件，该文件位于项目 `config/GeoLite2-City.mmdb` 目录，最后更新日志为 `2024-11-8`，最新库文件可从官方获取并替换即可。
