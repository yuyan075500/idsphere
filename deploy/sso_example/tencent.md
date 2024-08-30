# 腾讯云单点登录
支持的单点登录方式：SAML2
## 配置方法
1. **创建身份提供商**：登录腾讯云，进入【访问管理】-【身份提供商】-【用户SSO】，默认SSO登录是未开启，点击另边的【编辑】开启SSO，如下图所示：
![img.png](img/tencent-config2.jpg)
   * 用户SSO：启用。
   * SSO协议：`SAML`。
   * SAML服务提供商元数据URL：此URL复制保存，后面在平台注册站点时需要。
   * 身份提供商元数据文档：IDP的元数据文件可以访问平台`<protocol>://<address>[:<port>]/api/v1/sso/saml/metadata`获取，将网页中的内容保存到本地`xml`格式的文件中上传即可。
2. **创建CAM子用户**：接上步，进入【用户列表】-【】-【新建用户】，如下图所示：
   * 用户名：登录用户名，与平台中用户的`username`保持一至。
   > **说明**：当启用SSO后，所有子用户都无法通过账号密码通过腾讯云控制台登录，都将统一跳转至IDP认证。
3. **站点注册**：登录到平台，点击【资产管理】-【站点管理】-【新增】将腾讯云站点信息注册到平台，配置如下所示：
![img.png](img/tencent-site.jpg)
   * 站点名称：指定一个名称，便于用户区分。
   * 登录地址：腾讯的登录地址，默认为：`https://cloud.tencent.com/login/subAccount/<您的账号ID>?type=subAccount`。
   * SSO认证：启用。
   * 认证类型：选择`SAML2`。
   * 站点描述：描述信息。
   * SP Metadata URL：填写第一步中SAML服务提供商元数据URL，点击【获取】可以自动从腾讯云元数据中加载`SP EntityID`和`SP 证书`相关信息。