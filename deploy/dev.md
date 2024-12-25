# 开发环境搭建
## 准备开发工具
非必须的工具为推荐使用，如果有可代替的可自行选择，本指南不提供工具的使用教程，需要自行查找相关文档。

| 工具名称    | 用途                                                                                                                                | 必须  | 版本       |
|:--------|:----------------------------------------------------------------------------------------------------------------------------------|:----|:---------|
| VS Code | 前端辅助开发工具                                                                                                                          | ❌   | 推荐最新版本   |                                                                                                              |
| Goland  | 后端辅助开发工具                                                                                                                          | ❌   | 推荐最新版本   |
| Nginx   | 前后端代理程序                                                                                                                           | ✅   | 推荐最新版本   |
| Golang  | 后端运行环境                                                                                                                            | ✅   | 1.23.1   |
| Node.js | 前端运行环境                                                                                                                            | ✅   | v20.12.0 |
| 域名      | [Nginx鉴权](https://github.com/yuyan075500/idsphere/blob/main/deploy/sso.md#nginx%E4%BB%A3%E7%90%86%E9%89%B4%E6%9D%83 "Nginx鉴权") 必须 | ❌   | 无要求      |
| 域名证书    | [Nginx鉴权](https://github.com/yuyan075500/idsphere/blob/main/deploy/sso.md#nginx%E4%BB%A3%E7%90%86%E9%89%B4%E6%9D%83 "Nginx鉴权") 必须 | ❌   | 无要求      |
## 前端项目配置
前端项目开发环境配置可以参考 [前端项目开发调试](https://github.com/yuyan075500/idsphere-web?tab=readme-ov-file#%E5%BC%80%E5%8F%91%E8%B0%83%E8%AF%95 "前端项目开发调试") 文档。
## 后端项目配置
1. **创建项目配置文件**：需要在项目的 `config` 目录下创建 `config.yaml` 配置文件，配置文件内容为：<br><br>
   ```yaml
   server: "0.0.0.0:8000"
   externalUrl: "http://192.168.200.21"
   secret: "swfqezjzoqssvjck"
   mysql:
     host: "mysql"
     port: 3306
     db: "ops-api"
     user: "root"
     password: "X3UhzF9F"
     maxIdleConns: 10
     maxOpenConns: 100
     maxLifeTime: 30
   redis:
     host: "redis:6379"
     password: "o0qYcTrt"
     db: 0
   jwt:
     expires: 6
   mfa:
     enable: false
     issuer: "IDSphere 统一认证中心"
   oss:
     endpoint: "minio:9000"
     accessKey: "mXBbXV8nhjmLs8Ho"
     secretKey: "Zicc4ifKsX8dGwZHwro1"
     bucketName: "ops-api"
     ssl: false
   ldap:
     host: ""
     bindUserDN: ""
     bindUserPassword: ""
     searchDN: ""
     userAttribute: "uid"
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
   配置项修改请参考 [项目配置说明](https://github.com/yuyan075500/idsphere/blob/main/deploy/deploy.md#%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6%E8%AF%B4%E6%98%8E "项目配置说明")。<br><br>
2. **安装项目依赖包**：<br><br>
   ```shell
   go mod tidy
   ```
3. 构建项目：<br><br>
   ```shell
   go build
   ```
4. 运行项目：<br><br>
   ```shell
   ./ops-api
   ```
   项目默认监听 `8000` 端口，更换端口请修改配置文件中 `server` 项即可。
## 代理配置
由于项目使用前后端分离开发，为适配单点登录（SSO）相关协议以需要使用代理进行前后端分离部署，Nginx的配置文件内容如下：
```shell
server {
    listen       80;
    server_name  localhost;

    # 超时配置
    proxy_connect_timeout 60s;
    proxy_send_timeout 60s;
    proxy_read_timeout 60s;

    # 缓冲相关配置
    proxy_buffering on;
    proxy_buffers 16 4k;
    proxy_buffer_size 8k;
    proxy_busy_buffers_size 16k;

    # 后端API相关接口
    location ~ ^/(swagger|api|p3|validate|\.well-known/openid-configuration) {
        proxy_pass <后端服务器地址>;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 前端
    location / {
        proxy_pass <前端服务器地址>;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```
**注意**：如果需要用于 [Nginx鉴权](https://github.com/yuyan075500/idsphere/blob/main/deploy/sso.md#nginx%E4%BB%A3%E7%90%86%E9%89%B4%E6%9D%83 "Nginx鉴权") ，请使用域名+证书的访问方式作为 IDSphere 统一认证平台代理。