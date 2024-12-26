# Jenkins 单点登录
Jenkins 支持的单点登录方式：CAS3.0
## 配置方法
1. **站点注册**：登录到 IDSphere 统一认证平台，点击【资产管理】-【站点管理】-【新增】将 Jenkins 站点信息注册到 IDSphere 统一认证平台，配置如下所示：<br><br>
![img.png](img/jenkins-site.jpg)<br><br>
   * 站点名称：指定一个名称，便于用户区分。
   * 登录地址：Jenkins 的登录地址。
   * SSO 认证：启用。
   * 认证类型：选择 `CAS3.0`。
   * 站点描述：描述信息。
   * 回调地址：Jenkins 的回调地址，默认为：`http[s]://<address>[:<port>]/securityRealm/finishLogin`。<br><br>
2. **插件安装**：登录到 Jenkins 并进入【系统管理】-【插件管理】，如下图所示：<br><br>
![img.png](img/jenkins-plugin.png)<br><br>
选中箭头指示的 `CAS` 插件，点击右上角的【安装】，安装完成重启 Jenkins 即可。<br><br>
3. **插件配置**：登录到 Jenkins 并进入【系统管理】-【全局安全配置】，如下图所示：<br><br>
![img.png](img/jenins-cas.jpg)<br><br>
   * 安全域：选择 `CAS`。
   * CAS Server URL：CAS认证地址： `<externalUrl>`。
   * CAS Protocol：选择 `CAS 3.0`。
   * Full Name Attribute：更改为 `name`。
   * Email Attribute：更改为 `email`。<br><br>

**注意**：
   * 请注意红色框中的配置项，确保配置正确。
   * 如果 IDSphere 统一认证平台的用户登录到 Jenkins 后，用户姓名显示为乱码，引起该问题有可能是 Jenkins 在启动时 `JAVA_OPTS` 要指定使用 `UTF-8` 编码。