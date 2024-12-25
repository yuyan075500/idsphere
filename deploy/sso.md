# IDSphere 统一认证平台 SSO 功能介绍
IDSphere 统一认证平台支持与使用 `CAS 3.0`、`OAuth 2.0`、`OIDC`和`SAML2` 协议的客户端进行对接。本指南提供了标准客户端的对接方法，另外也提供了 [已通过测试的客户端列表](#已测试通过的客户端) 以供参考。<br><br>
**注意：目前所有对接协议都不支持单点注销，不支持Token刷新。**
# 自定义协议接入
以下方法将提供自定义接入所必须的接口及使用方法，适用于企业自研应用接入。在开始接入前需要在 IDSphere 统一认证平台上创建好对应的站点信息，站点创建好后包含应用接入所必须地相关配置信息，如：`client_id`、`client_secret`等。
# CAS3.0 客户端接入指南

| 接口名称 | 接口地址                  | 请求方法 | <div style="width:400px;">请求参数</div>                  | 请求参数类型 | 返回参数                                                                                  | 
|------|-----------------------|------|-------------------------------------------------------|--------|---------------------------------------------------------------------------------------|
| 登录地址 | `/login`              | GET  | `service`：客户端应用回调地址，必选。                               | Query  | 在 `URL` 中携带 `ticket` 信息，重定向至客户端应用指定的回调地址，如：<br>`https://a.com/callback?ticket=xxxxxx` |
| 票据授权 | `/p3/serviceValidate` | GET  | `service`：客户端应用回调地址，必选。<br>`ticket`：登录成功后获取到的票据信息，必选。 | Query  | `XML` 编码后的用户信息。                                                                       |

认证成功后返回的用户信息包含：`id`、`email`、`name`、`phone_number`和`username`。
# OAuth2.0 客户端接入指南

| 接口名称    | 接口地址                         | 请求方法        | <div style="width:440px;">请求参数</div>                                                                                                                                                                         | 请求参数类型 | 返回参数                                                                                                                                                                                                                              |
|---------|------------------------------|-------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| 登录地址    | `/login`                     | GET         | `response_type`：授权类型，固定值 `code`，必选。<br>`client_id`：客户端ID，从 IDSphere 统一认证平台获取，必选。<br>`redirect_uri`：客户端应用回调地址，必选。<br>`state`：状态码，由客户端指定任意值，请求成功后原样返回，必选。<br>`scope`：授权范围，固定值 `openid`，必选。                     | Query  | 在 `URL` 中携带 `code` 和 `state` 信息，重定向至客户端应用指定的回调地址，如：<br>`https://a.com/callback?code=xxxxxx&state=xxxxxx`                                                                                                                          |
| 获取Token | `/api/v1/sso/oauth/token`    | POST        | `grant_type`：授权类型，固定值 `authorization_code`，必选。<br>`client_id`：客户端ID，从 IDSphere 统一认证平台获取，必选。<br>`redirect_uri`：授权成功后的回调地址，必选。<br>`client_secret`：客户端Secret，从 IDSphere 统一认证平台获取，必选。<br>`code`：登录成功后获取到的授权码，必选。 | Body   | 返回 `Token` 信息，如：<br><pre><code>{<br>  "id_token": "",<br>  "access_token": "",<br>  "token_type": "bearer",<br>  "expires_in": 3600,<br>  "scope": "openid"<br>}</code></pre>`id_token` 和 `access_token` 都是采用 `JWT` 格式生成的`Token`。 |
| 获取用户信息  | `/api/v1/sso/oauth/userinfo` | GET<br>POST | 在请求头中携带 `Authorization` 字段，值为`Bearer <access_token>`，必选。                                                                                                                                                     | Header | JSON格式用户信息。                                                                                                                                                                                                                       |

认证成功后返回的用户信息包含：`id`、`email`、`name`、`phone_number`和`username`。
# OIDC 客户端接入指南
`OIDC` 客户端接入参考 `OAuth2.0` 接入指南即可，另 IDSphere 统一认证平台提供了专属 `OIDC` 配置信息接口，接口地址为：`/.well-known/openid-configuration`。<br><br>
使用 `OIDC` 获取的 JWT Token 的 Claims 信息如下：
```shell
{
    "id": <用户ID>,
    "name": "<用户姓名>",
    "username": "<用户名>",
    "azp": "<ClientID>",
    "policy": "readwrite",
    "nonce": "<授权时传入的nonce参数>",
    "iss": "<externalUrl>",
    "sub": "user-<用户ID>",
    "aud": [
        "<ClientID>"
    ],
    "exp": <过期时间>,
    "nbf": <生效时间>,
    "iat": <签发时间>
}
```
# SAML2 客户端接入指南
`SAML2` 客户端接入所需要信息可以通过 `IDP` 元数据接口地址获取，接口地址为：`/api/v1/sso/saml/metadata`，认证成功后返回的用户信息包含：

| 属性值                              | 属性名称                            |
|----------------------------------|---------------------------------|
| name                             | 姓名                              |
| username                         | 用户名                             |
| email                            | 邮箱地址                            |
| phone_number                     | 电话号码                            |
| IAM_SAML_Attributes_xUserId      | 华为云专属，对应用户名                     |
| IAM_SAML_Attributes_redirect_url | 华为云专属，登录成功后的跳地址                 |
| IAM_SAML_Attributes_domain_id    | 华为云专属，在华为云配置时，自动生成的 `domain_id` |
| IAM_SAML_Attributes_idp_id       | 华为云专属，在华为云配置时指定的的身份提供商名称        |

# Nginx 代理鉴权
此功能可以针对 `Nginx` 代理的路径进行鉴权，以实现基于 `Cookie` 的单点登录。如：Kibana Dashboard、Consul Server UI等其它所有不需要鉴权就能访问的页面，为了确保安全性，那么可以使用Nginx作为代理，使其需要认证才能访问。<br><br>
由于该功能是基=于基于 `Cookie` 实现，所以在使用上有一定限制，要求如下：
1. 客户端的所在域必须与 IDSphere 统一认证平台域一致，假如IDSphere 统一认证平台的访问域为：`test.idsphere.cn`，则客户端的所在域必须为`xxx.idsphere.cn`。
2. 不支持使用 `IP` 访问的客户端。如：`127.0.0.1`、`localhost`、`192.168.1.10`。 
3. 不支持非HTTPS应用。
## 认证流程
![img.png](sso_example/img/nginx.jpg)
## 认证规则
假设有A、B、C三个客户端应用都使用Nginx进行代理鉴权，鉴权规则如下：
* 如果用户未登录，当访问A、B、C三个客户端应用中的任何一个都将跳转至IDSphere 统一认证平台登录界面。
* 如果用户已经登录，当访问A、B、C三个客户端应用中的任何一个都能直接访问应用。
## Nginx 配置
在开始配置前请确保 `Nginx` 支持 `auth_request` 模块，可以使用命令 `nginx -V` 查看，具体配置如下：
```nginx
server {
	listen 80;

	location / {
		auth_request /auth;
		error_page 401 500 = @error401;
		# 后面是被代理应用的相关配置
	}
	
	location = /auth {
		internal;
		proxy_pass_request_body off;
		proxy_set_header Content-Length "";
		proxy_set_header X-Original-URI $request_uri;
		proxy_set_header Cookie $http_cookie;
		proxy_pass <externalUrl>/api/v1/sso/cookie/auth;
	}

	# 认证失败后的处理
	location @error401 {
		# 跳转至登录页
		return 302 <externalUrl>/login?nginx_redirect_uri=$scheme://$host$request_uri;
	}
}
```
在上面的示例中对 `/` 根路径进行了代理鉴权，如果需要对其它路径进行鉴权可以添加其它 `location`，在 `location` 里面添加对应的 `auth_request` 和 `error_page`。
## Kubernetes Ingress 配置
此功能同样支持在 `Kubernetes` 中使用 `Ingress` 进行配置，示例如下：
```yaml
kind: Ingress
apiVersion: networking.k8s.io/v1
metadata:
  name: my-app
  annotations:
    nginx.ingress.kubernetes.io/auth-url: <externalUrl>/api/v1/sso/cookie/auth
    nginx.ingress.kubernetes.io/server-snippet: |
      error_page 401 500 = @login;
      proxy_set_header Cookie $http_cookie;
      location @login {
        return 302 <externalUrl>/login?nginx_redirect_uri=$scheme://$host$request_uri;
      }
spec:
  ingressClassName: nginx
  tls:
    - hosts:
        - my-app.idsphere.cn
      secretName: idsphere.cn
  rules:
    - host: my-app.idsphere.cn
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: my-app
                port:
                  number: 80
```
# 已测试通过的客户端
针对 `CAS3.0`、`OAuth2.0`、`SAML2`、`OIDC` 协议，目前经测试对接成功的SSO客户端如下：

| 客户端名称      | 对接协议名称   | 参考文档                                                                                                |
|:-----------|:---------|-----------------------------------------------------------------------------------------------------|
| Grafana    | OAuth2.0 | [参考文档](https://github.com/yuyan075500/idsphere/blob/main/deploy/sso_example/grafana.md "参考文档")      |
| Jenkins    | CAS3.0   | [参考文档](https://github.com/yuyan075500/idsphere/blob/main/deploy/sso_example/jenkins.md "参考文档")      |
| Zabbix     | SAML2    | [参考文档](https://github.com/yuyan075500/idsphere/blob/main/deploy/sso_example/zabbix.md "参考文档")       |
| 华为云        | SAML2    | [参考文档](https://github.com/yuyan075500/idsphere/blob/main/deploy/sso_example/huawei_cloud.md "参考文档") |
| JumpServer | OAuth2.0 | [参考文档](https://github.com/yuyan075500/idsphere/blob/main/deploy/sso_example/jumpserver.md "参考文档")   |
| Jira       | OAuth2.0 | [参考文档](https://github.com/yuyan075500/idsphere/blob/main/deploy/sso_example/jira.md "参考文档")         |
| Confluence | OAuth2.0 | [参考文档](https://github.com/yuyan075500/idsphere/blob/main/deploy/sso_example/confluence.md "参考文档")   |
| KubePi     | OIDC     | [参考文档](https://github.com/yuyan075500/idsphere/blob/main/deploy/sso_example/kubepi.md "参考文档")       |
| 阿里云        | SAML2    | [参考文档](https://github.com/yuyan075500/idsphere/blob/main/deploy/sso_example/aliyun.md "参考文档")       |
| 腾讯云        | SAML2    | [参考文档](https://github.com/yuyan075500/idsphere/blob/main/deploy/sso_example/tencent.md "参考文档")      |
| Minio      | OIDC     | [参考文档](https://github.com/yuyan075500/idsphere/blob/main/deploy/sso_example/minio.md "参考文档")        |
| GitLab     | OIDC     | [参考文档](https://github.com/yuyan075500/idsphere/blob/main/deploy/sso_example/gitlab.md "参考文档")       |
| 天翼云        |          | 待测试                                                                                                 |
| Rancher    |          | 待测试                                                                                                 |
| 禅道         |          | 待测试                                                                                                 |
| AWS        |          | 待测试                                                                                                 |

PS：如果你有其它第三方系统需要对接可以提交 `Issue` 请求。