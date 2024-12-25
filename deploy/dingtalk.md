# 配置钉钉应用
1. **登录钉钉开放平台**：https://open.dingtalk.com ，并进入开发者后台。
2. **创建应用**：参考 [官方文档](https://open.dingtalk.com/document/orgapp/create-an-application "官方文档")。
3. **应用配置**：进入应用详情页，单击【开发配置】 > 【安全设置】，填写重定向 URL（回调域名）。回调域名为 IDSphere 统一认证平台的登录地址，为`<externalUrl>/login`，当使用钉钉手机 APP 扫码后浏览器默认会跳转至该地址。
4. **应用授权**：进入应用详情页，单击【开发配置】 > 【权限管理】，在权限搜索框中输入权限名称并申请权限。需要授与该应用`Contact.User.mobile`和`Contact.User.Read`权限，参考 [官方文档](https://open.dingtalk.com/document/orgapp/tutorial-obtaining-user-personal-information#c4647d84328mg "官方文档")。
5. **发布应用**：参考 [官方文档](https://open.dingtalk.com/document/orgapp/publish-dingtalk-application "官方文档")。
# 后端应用配置
需要在配置中添加钉钉应用的相关配置。
```yaml
dingtalk:
  appKey: ""
  appSecret: ""
```
* [x] appKey：在钉钉开放平台，应用详情页，左侧的【凭证与基础信息】中获取。
* [x] appSecret：在钉钉开放平台，应用详情页，左侧的【凭证与基础信息】中获取。
# 前端应用配置
参考 [前端配置](https://github.com/yuyan075500/idsphere-web "前端配置") 相关文档，修改配置文件中有关于钉钉相关的配置项，并手动构建打包项目，生成新的容器镜像。
# 创建本地用户
钉钉扫码登录需要事先在 IDSphere 统一认证平台中创建对应的用户，或者从 LDAP、Windows AD 中同步用户到本地。确保`用户姓名`和`手机号`与钉钉用户一致且本地用户状态（账号未禁用、密码未过期）正常，则可以登录成功。