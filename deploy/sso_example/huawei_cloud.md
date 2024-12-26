# 华为云单点登录
支持的单点登录方式：SAML2
## 配置方法
1. **创建身份提供商**：登录华为云，进入【统一身份认证】-【身份提供商】-【创建身份提供商】如下所示：<br><br>
![img.png](img/huawei-create-idp.jpg)<br><br>
   * 名称：指定一个名称，便于区分。
   * 协议：选择 `SAML`。
   * 类型：建议使用 `IAM用户SSO`，需要创建一个 IAM 实体用户。
   * 状态：启用。<br><br>
2. **身份提供商配置**：创建完身份提供商后点击【修改】进行配置，如下所示：<br><br>
![img.png](img/huawei-idp-config.jpg)<br><br>
   IDP的元数据文件可以访问 IDSphere 统一认证平台获取，地址为：`<externalUrl>/api/v1/sso/saml/metadata`，将网页中的内容保存至 `xml` 格式的文件中上传即可，另外还需要保存登录地址，后续在平台注册站点时使用。<br><br>
3. **创建IAM用户**：进入【统一身份认证】-【用户】-【创建用户】，如下图所示：<br><br>
![img.png](img/huawei-create-user.jpg)<br><br>
   * 用户名：登录用户名，与平台中用户的名保持一至。
   * 外部身份 ID：与用户名保持一致。  

   > **说明**：由于身份提供商的类型是 IAM 用户 SSO，所以需要创建一个 IAM 实体用户，否则无法登录。

4. **站点注册**：登录到平台，点击【资产管理】-【站点管理】-【新增】将华为云站点信息注册到平台，配置如下所示：<br><br>
![img.png](img/huawei-site.jpg)<br><br>
   * 站点名称：指定一个名称，便于用户区分。
   * 登录地址：华为云的登录地址，与华为云身份提供商配置界面保持一致。
   * SSO 认证：启用。
   * 认证类型：选择 `SAML2`。
   * 站点描述：描述信息。
   * SP Metadata URL：填写 [华为云元数据访问地址](https://auth.huaweicloud.com/authui/saml/metadata.xml "华为云元数据访问地址") 相关信息，点击【获取】按钮自动从华为云元数据中加载 `SP EntityID` 和 `SP 证书` 相关信息。<br><br>
5. **站点修改**：登录到 IDSphere 统一认证平台数据库，在 `site` 表中找到刚注册的站点信息，需要修改 `domain_id`、`redirect_url` 和 `idp_name` 这三个字段。<br><br>
   * domain_id：与华为云身份提供商配置界面的登录连接中的 `domain_id` 保持一致。
   * redirect_url：用户登录成功后的跳转地址，如：`https://console.huaweicloud.com/console/?region=cn-east-3` 。
   * idp_name：与华为云创建的身份提供商名称保持一致。