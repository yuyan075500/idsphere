# IDSphere 统一认证平台项目介绍
仅需一次认证，即可访问所有授权访问的应用系统，为企业办公人员提供高效、便捷的访问体验。
## 演示站点
演示站点访问地址为：`https://sso.idsphere.cn`，相关账号密码如下：

| 账号名称 | 用户名      | 密码                | 用户来源     | 系统权限 | 可访问站点                |
|------|----------|-------------------|----------|------|----------------------|
| 王五   | `wangw`  | `Wang-W@123...`   | OpenLDAP | 管理员  | 全部                   |
| 张三   | `zhangs` | `Zahang-S@123...` | OpenLDAP | 管理员  | 全部                   |
| 李四   | `lisi`   | `Li-Si@123...`    | 本地       | 普通用户 | `Jenkins` 和 `GitLab` |
| 审计用户 | `audit`  | `AAudit@123...`   | 本地       | 审计员  | `GitLab`             |

演示环境系统部分功能不支持，列表如下：
1. 双因素认证。
2. 短信验证码。
3. 钉钉、企业微信、飞书扫码登录。
4. 系统配置修改。

如果有其它需求可 [联系作者](#项目交流)。
## 架构设计
项目采用前后端分离架构设计，项目地址如下：

| 项目  | 项目地址                                        |
|:----|:--------------------------------------------|
| 前端  | https://github.com/yuyan075500/idsphere-web |                                                                                                              |
| 后端  | https://github.com/yuyan075500/idsphere     |

以上两个仓库与 `Gitee` 保持同步，仓库地址如下：

| 项目  | 项目地址                                       |
|:----|:-------------------------------------------|
| 前端  | https://gitee.com/yybluestorm/idsphere-web |                                                                                                              |
| 后端  | https://gitee.com/yybluestorm/idsphere     |
## 后端目录说明
* config：全局配置。
* controller：路由规则配置和接口的入参与响应。
* service：接口处理逻辑。
* dao：数据库操作。
* model：数据库模型定义。
* db：数据库、缓存以及文件存储客户端初始化。
* middleware：全局中间件层，如跨域、JWT认证、权限校验等。
* utils：全局工具层，如Token解析、文件操作、字符串操作以及加解密等。
## 后端 Code 状态码说明
* 0：请求成功。
* 90400：请求参数错误。
* 90401：认证失败。
* 90403：拒绝访问。
* 90404：访问的对象或资源不存在。
* 90514：Token过期或无效。
* 90500：其它服务器错误。
# 项目功能介绍
## 认证相关
* **SSO 单点登录**：支持 `CAS 3.0`、`OAuth 2.0`、`OIDC`和`SAML2` 协议，客户端对接对接方法可以参考 [客户端配置指南](https://github.com/yuyan075500/idsphere/wiki/6%E3%80%81%E5%8D%95%E7%82%B9%E7%99%BB%E5%BD%95%EF%BC%88SSO%EF%BC%89%E5%AE%A2%E6%88%B7%E7%AB%AF%E6%8E%A5%E5%85%A5%E6%8C%87%E5%8D%97 "SSO 客户端对接") 和 [已测试客户端列表](https://github.com/yuyan075500/idsphere/wiki/6%E3%80%81%E5%8D%95%E7%82%B9%E7%99%BB%E5%BD%95%EF%BC%88SSO%EF%BC%89%E5%AE%A2%E6%88%B7%E7%AB%AF%E6%8E%A5%E5%85%A5%E6%8C%87%E5%8D%97#%E5%B7%B2%E9%80%9A%E8%BF%87%E6%B5%8B%E8%AF%95%E7%9A%84%E5%AE%A2%E6%88%B7%E7%AB%AF%E5%88%97%E8%A1%A8 "已测试客户端列表")。
* **用户认证**：支持使用 [钉钉扫码登录](https://github.com/yuyan075500/idsphere/wiki/5%E3%80%81%E7%94%A8%E6%88%B7%E8%AE%A4%E8%AF%81#%E9%92%89%E9%92%89 "钉钉扫码配置")、[企业微信扫码登录](https://github.com/yuyan075500/idsphere/wiki/5%E3%80%81%E7%94%A8%E6%88%B7%E8%AE%A4%E8%AF%81#%E4%BC%81%E4%B8%9A%E5%BE%AE%E4%BF%A1 "企业微信扫码配置")、[飞书扫码登录](https://github.com/yuyan075500/idsphere/wiki/5%E3%80%81%E7%94%A8%E6%88%B7%E8%AE%A4%E8%AF%81#%E9%A3%9E%E4%B9%A6 "飞书扫码配置")、[OpenLDAP 账号密码登录](https://github.com/yuyan075500/idsphere/wiki/5%E3%80%81%E7%94%A8%E6%88%B7%E8%AE%A4%E8%AF%81#openldap "OpenLDAP 配置")和[Windows AD 账号密码登录](https://github.com/yuyan075500/idsphere/wiki/5%E3%80%81%E7%94%A8%E6%88%B7%E8%AE%A4%E8%AF%81#windows-ad "Windows AD配置") 登录。登录页面支持个性化配置，隐藏或显示必要的登录选项，可以参考 [前端配置指南](https://github.com/yuyan075500/ops-web "前端配置")。
* **双因素认证**：支持使用 Google Authenticator、阿里云和华为云手机 APP 进行双因素认证，双因素认证仅在使用账号密码认证时生效。

    <br>
    <img src="deploy/image/login-1.gif" alt="img" width="350" height="200"/>
    <img src="deploy/image/login-mfa.gif" alt="img" width="350" height="200"/>
    <br>

## 企业级账号管理
支持账号资产管理功能，面向企业人员提供安全、便捷的企业信息化账号管理。
* 加密保护：账号密码采用密钥对加密存储，避免本地存储带来的泄露风险。
* 共享机制：支持用户间的账号共享，无需依赖邮件或聊天工具，分享更加高效安全。
* 所有权变更机制：人员变动时账号可直接移交给新负责人，避免资产流失。

具体可以参考 [账号管理相关文档](https://github.com/yuyan075500/idsphere/wiki/4%E3%80%81%E7%B3%BB%E7%BB%9F%E4%BD%BF%E7%94%A8#%E8%B4%A6%E5%8F%B7%E7%AE%A1%E7%90%86 "账号管理")。
## 域名及证书管理
计划于 3.0 版本上线
## 其它
* 支持`Swagger`接口文档：部署成功后访问地址为：`/swagger/index.html`，无需要登录。
* 支持用户密码自助更改：部署成功后访问地址：`/reset_password`，无需要登录。
* 支持企业网站导航：部署成功后访问地址：`/sites`，无需要登录。
# 项目部署
参考 [Docker Compose部署](https://github.com/yuyan075500/idsphere/wiki/%E5%AE%89%E8%A3%85%E9%83%A8%E7%BD%B2#docker-compose-%E9%83%A8%E7%BD%B2 "docker-compose部署") 和 [Kubernetes部署](https://github.com/yuyan075500/idsphere/wiki/%E5%AE%89%E8%A3%85%E9%83%A8%E7%BD%B2#kubernetes-%E9%83%A8%E7%BD%B2 "Kubernetes部署")。
# 开发环境搭建
参考 [开发环境搭建](https://github.com/yuyan075500/idsphere/wiki/%E5%BC%80%E5%8F%91%E7%8E%AF%E5%A2%83%E6%90%AD%E5%BB%BA "开发环境搭建")。
# 项目交流
如果你对此项目感兴趣，可添加作者联系方式
WeChat：270142877。  
Email：270142877@qq.com。  
<br>
联系时请注名来意。