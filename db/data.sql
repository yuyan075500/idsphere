SET NAMES 'utf8mb4';

# 一级菜单
INSERT INTO `system_menu` VALUES (1, '用户管理', 'User', 'menu-user', '/user', 'Layout', 2, null);
INSERT INTO `system_menu` VALUES (2, '资产管理', 'Asset', 'menu-asset', '/asset', 'Layout', 3, null);
INSERT INTO `system_menu` VALUES (3, '日志审计', 'Audit', 'menu-audit', '/audit', 'Layout', 4, null);
INSERT INTO `system_menu` VALUES (4, '系统设置', 'System', 'menu-system', '/system', 'Layout', 5, null);

# 二级菜单
INSERT INTO `system_sub_menu` VALUES (1, '用户管理', 'UserManagement', 'sub-menu-user', 'user', 'user/user/index', 1, null, 1);
INSERT INTO `system_sub_menu` VALUES (2, '分组管理', 'GroupManagement', 'sub-menu-group', 'group', 'user/group/index', 2, null, 1);
INSERT INTO `system_sub_menu` VALUES (3, '账号管理', 'AccountManagement', 'sub-menu-account', 'account', 'asset/account/index', 1, null, 2);
INSERT INTO `system_sub_menu` VALUES (4, '站点管理', 'SiteManagement', 'sub-menu-site', 'site', 'asset/site/index', 2, null, 2);
INSERT INTO `system_sub_menu` VALUES (5, '登录日志', 'AuditLoginRecord', 'sub-menu-login-record', 'login', 'audit/login/index', 1, null, 3);
INSERT INTO `system_sub_menu` VALUES (6, '短信记录', 'AuditSMSRecord', 'sub-menu-sms-record', 'sms', 'audit/sms/index', 2, null, 3);
INSERT INTO `system_sub_menu` VALUES (7, '操作日志', 'AuditOplog', 'sub-menu-oplog', 'oplog', 'audit/oplog/index', 3, null, 3);
INSERT INTO `system_sub_menu` VALUES (8, '菜单管理', 'MenuManagement', 'sub-menu-menu', 'menu', 'system/menu/index', 1, null, 4);
INSERT INTO `system_sub_menu` VALUES (9, '定时任务', 'CornManagement', 'sub-menu-corn', 'corn', 'system/corn/index', 2, null, 4);
INSERT INTO `system_sub_menu` VALUES (10, '系统配置', 'ConfManagement', 'sub-menu-conf', 'conf', 'system/settings/index', 3, null, 4);

# API接口
INSERT INTO `system_path` VALUES (1, 'AddUser', '/api/v1/user', 'POST', 'UserManagement', '新增用户');
INSERT INTO `system_path` VALUES (2, 'UpdateUser', '/api/v1/user', 'PUT', 'UserManagement', '修改用户');
INSERT INTO `system_path` VALUES (3, 'UpdateUserPassword', '/api/v1/user/reset_password', 'PUT', 'UserManagement', '密码重置');
INSERT INTO `system_path` VALUES (4, 'ResetUserMFA', '/api/v1/user/reset_mfa/:id', 'PUT', 'UserManagement', 'MAF重置');
INSERT INTO `system_path` VALUES (5, 'DeleteUser', '/api/v1/user/:id', 'DELETE', 'UserManagement', '删除用户');
INSERT INTO `system_path` VALUES (6, 'GetUserList', '/api/v1/users', 'GET', 'UserManagement', '获取用户列表（表格）');
INSERT INTO `system_path` VALUES (7, 'UserSyncAd', '/api/v1/user/sync/ad', 'POST', 'UserManagement', '用户同步');
INSERT INTO `system_path` VALUES (8, 'GetUserListAll', '/api/v1/user/list', 'GET', 'UserManagement', '获取用户列表（所有）');
INSERT INTO `system_path` VALUES (9, 'AddGroup', '/api/v1/group', 'POST', 'GroupManagement', '新增分组');
INSERT INTO `system_path` VALUES (10, 'UpdateGroup', '/api/v1/group', 'PUT', 'GroupManagement', '修改分组');
INSERT INTO `system_path` VALUES (11, 'UpdateGroupUser', '/api/v1/group/users', 'PUT', 'GroupManagement', '更改分组用户');
INSERT INTO `system_path` VALUES (12, 'UpdateGroupPermission', '/api/v1/group/permissions', 'PUT', 'GroupManagement', '更改分组权限');
INSERT INTO `system_path` VALUES (13, 'DeleteGroup', '/api/v1/group/:id', 'DELETE', 'GroupManagement', '删除分组');
INSERT INTO `system_path` VALUES (14, 'GetGroupList', '/api/v1/groups', 'GET', 'GroupManagement', '获取分组列表');
INSERT INTO `system_path` VALUES (15, 'GetMenuListAll', '/api/v1/menu/list', 'GET', 'GroupManagement', '获取菜单列表');
INSERT INTO `system_path` VALUES (16, 'GetPathListAll', '/api/v1/path/list', 'GET', 'GroupManagement', '获取接口列表');
INSERT INTO `system_path` VALUES (17, 'GetSiteList', '/api/v1/sites', 'GET', 'SiteManagement', '获取站点列表');
INSERT INTO `system_path` VALUES (18, 'AddSite', '/api/v1/site', 'POST', 'SiteManagement', '新增站点');
INSERT INTO `system_path` VALUES (19, 'UpdateSite', '/api/v1/site', 'PUT', 'SiteManagement', '修改站点');
INSERT INTO `system_path` VALUES (20, 'DeleteSite', '/api/v1/site/:id', 'DELETE', 'SiteManagement', '删除站点');
INSERT INTO `system_path` VALUES (21, 'AddSiteGroup', '/api/v1/site/group', 'POST', 'SiteManagement', '新增站点分组');
INSERT INTO `system_path` VALUES (22, 'UpdateSiteGroup', '/api/v1/site/group', 'PUT', 'SiteManagement', '修改站点分组');
INSERT INTO `system_path` VALUES (23, 'DeleteSiteGroup', '/api/v1/site/group/:id', 'DELETE', 'SiteManagement', '删除站点分组');
INSERT INTO `system_path` VALUES (24, 'UpdateSiteUser', '/api/v1/site/users', 'PUT', 'SiteManagement', '更改站点用户');
INSERT INTO `system_path` VALUES (25, 'UpdateSiteTag', '/api/v1/site/tags', 'PUT', 'SiteManagement', '更改站点标签');
INSERT INTO `system_path` VALUES (26, 'GetSMSRecordList', '/api/v1/audit/sms', 'GET', 'AuditSMSRecord', '获取短信发送记录');
INSERT INTO `system_path` VALUES (27, 'GetLoginRecordList', '/api/v1/audit/login', 'GET', 'AuditLoginRecord', '获取用户登录记录');
INSERT INTO `system_path` VALUES (28, 'GetOplogList', '/api/v1/audit/oplog', 'GET', 'AuditOplog', '获取用户操作日志');
INSERT INTO `system_path` VALUES (29, 'GetMenuList', '/api/v1/menus', 'GET', 'MenuManagement', '获取菜单列表');
INSERT INTO `system_path` VALUES (30, 'GetPathList', '/api/v1/paths', 'GET', 'MenuManagement', '获取菜单接口');
INSERT INTO `system_path` VALUES (31, 'GetTaskList', '/api/v1/tasks', 'GET', 'CornManagement', '获取定时任务列表');
INSERT INTO `system_path` VALUES (32, 'AddTask', '/api/v1/task', 'POST', 'CornManagement', '新增定时任务');
INSERT INTO `system_path` VALUES (33, 'UpdateTask', '/api/v1/task', 'PUT', 'CornManagement', '修改定时任务');
INSERT INTO `system_path` VALUES (34, 'DeleteTask', '/api/v1/task/:id', 'DELETE', 'CornManagement', '删除定时任务');
INSERT INTO `system_path` VALUES (35, 'GetTaskLogList', '/api/v1/task/logs', 'GET', 'CornManagement', '获取定时任务执行日志列表');
INSERT INTO `system_path` VALUES (36, 'GetSettings', '/api/v1/settings', 'GET', 'ConfManagement', '获取配置信息');
INSERT INTO `system_path` VALUES (37, 'UpdateLogo', '/api/v1/settings/logoUpload', 'POST', 'ConfManagement', '修改 Logo');
INSERT INTO `system_path` VALUES (38, 'UpdateSettings', '/api/v1/settings', 'PUT', 'ConfManagement', '修改配置信息');

# 系统默认配置
INSERT INTO `settings` VALUES (1, 'externalUrl', 'https://example.idsphere.cn', 'string');
INSERT INTO `settings` VALUES (2, 'logo', null, 'string');
INSERT INTO `settings` VALUES (3, 'mfa', 'false', 'boolean');
INSERT INTO `settings` VALUES (4, 'issuer', 'IDSphere 统一认证中心', 'string');
INSERT INTO `settings` VALUES (5, 'secret', 'swfqezjzoqssvjck', 'string');
INSERT INTO `settings` VALUES (6, 'ldapAddress', null, 'string');
INSERT INTO `settings` VALUES (7, 'ldapBindDn', null, 'string');
INSERT INTO `settings` VALUES (8, 'ldapBindPassword', null, 'string');
INSERT INTO `settings` VALUES (9, 'ldapSearchDn', null, 'string');
INSERT INTO `settings` VALUES (10, 'ldapFilterAttribute', 'uid', 'string');
INSERT INTO `settings` VALUES (11, 'ldapUserPasswordExpireDays', '90', 'int');
INSERT INTO `settings` VALUES (12, 'passwordExpireDays', '90', 'int');
INSERT INTO `settings` VALUES (13, 'passwordLength', '8', 'int');
INSERT INTO `settings` VALUES (14, 'passwordComplexity', '["numbers","uppercase","lowercase"]', 'list');
INSERT INTO `settings` VALUES (15, 'passwordExpiryReminderDays', '7', 'int');
INSERT INTO `settings` VALUES (16, 'certificate', '-----BEGIN CERTIFICATE-----
MIIDazCCAlOgAwIBAgIUTmM3pO+/sy8prxjOo5s3RlGcWYswDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yNDA4MTIwOTAzMzNaFw0zNDA4
MTAwOTAzMzNaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQDBZtVdydM3KSU83ePj4DqhExqV4taIS1h4n84ODa7M
wGjjxQq0R10mtxmJQH4NCqqa7Z0crjoqjbM9eh605Rk/naX3a8NU5OvVkB88qJMq
44wezdgsIoRZMAiCc1HQCR8H3WkHg+SsJFSFa4K2EsP9+MBGARpTerRBhs+ZsIBr
Gt1lVzG6B36UdOHiAV5JVfw4SHexX8oGJ3T0a48WTVOYycxK5+8INQNCtMkoTiDU
IvRwAhgkAt/vBYJfOjIEWXl4sUBz7ZqKEQJjeQ+cjSLn1bZ0jSojOpySBlo1RhaL
8ZzyJYRk0cxE4ta+LX5OHvk8OdvgMj25sYHAMkUodmTRAgMBAAGjUzBRMB0GA1Ud
DgQWBBTVSwUcHGKg8zpb7rwzdS7GQeoMaDAfBgNVHSMEGDAWgBTVSwUcHGKg8zpb
7rwzdS7GQeoMaDAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQBW
19bRn4k4aTOWhbb1tqIw5tEMCBsT+UfS8qbD8i1172RMn08BOHfdZF1k39EyqsXE
2B4b9eYParduYSqpPlx7PlLAA6aJilJkqKRa/y8m9Le8iWT89NpwDeXbkmrd2f0Q
4+vqj/OSkIoK4N49mk9I3J0EKvcND5bCONudGXIOV9VowMa9/nGQMuwcUXTQEZk4
4vbnXoE/ctmFqMYPDADAmDXl6YDztz2xbXvQA3voEaATFvhvyYFn5VjenupYpoG2
2GxGZNlI2l6PxBKxri3qwlblVbIee6Q4SNJ51f1A9Gzh9bwBd/mJGtNZXk8qBTDP
ezxXCSubKEeDFdRHiPIF
-----END CERTIFICATE-----', 'string');
INSERT INTO `settings` VALUES (17, 'publicKey', '-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwWbVXcnTNyklPN3j4+A6
oRMaleLWiEtYeJ/ODg2uzMBo48UKtEddJrcZiUB+DQqqmu2dHK46Ko2zPXoetOUZ
P52l92vDVOTr1ZAfPKiTKuOMHs3YLCKEWTAIgnNR0AkfB91pB4PkrCRUhWuCthLD
/fjARgEaU3q0QYbPmbCAaxrdZVcxugd+lHTh4gFeSVX8OEh3sV/KBid09GuPFk1T
mMnMSufvCDUDQrTJKE4g1CL0cAIYJALf7wWCXzoyBFl5eLFAc+2aihECY3kPnI0i
59W2dI0qIzqckgZaNUYWi/Gc8iWEZNHMROLWvi1+Th75PDnb4DI9ubGBwDJFKHZk
0QIDAQAB
-----END PUBLIC KEY-----', 'string');
INSERT INTO `settings` VALUES (18, 'privateKey', '-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDBZtVdydM3KSU8
3ePj4DqhExqV4taIS1h4n84ODa7MwGjjxQq0R10mtxmJQH4NCqqa7Z0crjoqjbM9
eh605Rk/naX3a8NU5OvVkB88qJMq44wezdgsIoRZMAiCc1HQCR8H3WkHg+SsJFSF
a4K2EsP9+MBGARpTerRBhs+ZsIBrGt1lVzG6B36UdOHiAV5JVfw4SHexX8oGJ3T0
a48WTVOYycxK5+8INQNCtMkoTiDUIvRwAhgkAt/vBYJfOjIEWXl4sUBz7ZqKEQJj
eQ+cjSLn1bZ0jSojOpySBlo1RhaL8ZzyJYRk0cxE4ta+LX5OHvk8OdvgMj25sYHA
MkUodmTRAgMBAAECggEAA3ZI9SwC5Ra498J3F7dbWqfMuYh9KPpu+Sskf/mRuqY/
xRgGuDHGLUZo0B1V7vt/u1SWk+gF2t4nf78QY/ztO0d3VvBYZ6F4xp1qWbvLdUtB
Jl9ZR8z5tokBa1MPe1xhGv3FKWEci6eo1ayWLevfE0xL0NT943xWVA1K6RLmRu8v
qCCxD/GS4FvqH7fx2YXdYvyW8UE/sXTB5T4SWJMPYkbBMk3CJSBtF6YO6/VyOcgP
pHOnNtxVECIXyJuE7VSga+Is69auu5P5SsZCOJl9ofvcJ35oIwiYEVpYa0LpGsfX
es4Hfq5zNF72RDAWXRCXRjfQQWEJ10ACzVcnjZJh4wKBgQD6YdcZ5wY5IME5nU9+
5KzyBsrtw58Tjz9xVldkknPawZnySuX50noHwuAdT6Ju8AlXVbiMzxcH54ZUcm1h
jpGfGLZtF25zjzWo7CCdFc+OMWMcg56QP62buuqYtnxYIWiXJzg1ZHB3G5o7Vrpu
xACwSllM8GtaAb3rBtU2ZNMHUwKBgQDFvbSRcsehMFHUNYzoM5ae3XJ/hlnWNMMh
1wob5jYTeuOXk+3uCLC3CXQrY1Zeqnc1DHmyWjiP9C/zLpkNz4bz1LPXCA/icMrL
LHN4aaK59bsggFImgPDi4krnOQiU/JqkHEgzmy2IMHwTR4NhC5CXQI9jlFFoHFyc
SNkP6j1SywKBgCl+f3RWegyLphoTPPJtmU++2nCO49UM/1mcEn2jW7ncLdQen4BI
ZlrU6+lPoj66XwHvPddwFoQD1Zo3IHNzeiSgptLreC2EhUMKZtlBnRUWkDNQiL2l
H/NYBbrrOy4r0zaBlGoczBqhI21EET98EhUlLSl4CoJvGXdSuZD7IpHrAoGAHnAd
I2ZvpDgz4F592iBWxw1/WnHr0jU89DCNtc2x9T2tWt/CeCmOSh6Ca0tXOCs1Pk01
Tmbk3gPQfbZmiOGw/Ed5h1gOWeTS0oN9IsPf8JAKxe36t0KR0drTfNQipgxcIXbZ
BliUoaoU70LKzl1hXGbrq4BhJ412E/iCsRh1aBECgYEApEFLxRvhe19OrjnXVWTG
Abn2baPW6YbbXEyFbERGfS12lmsADYqQO1eVlODV9fLuxTez+BIwBnt3Q+sj3TLm
116HYANFTiU6+doZ2JA+Xbng0OL6DsJ2yV+nHM6sROskcRJMAnYxylEBycSG9Y9z
TJTvjLxCzTsUCGhK6oen3NI=
-----END PRIVATE KEY-----', 'string');
INSERT INTO `settings` VALUES (19, 'mailAddress', null, 'string');
INSERT INTO `settings` VALUES (20, 'mailPort', '25', 'int');
INSERT INTO `settings` VALUES (21, 'mailForm', null, 'string');
INSERT INTO `settings` VALUES (22, 'mailPassword', null, 'string');
INSERT INTO `settings` VALUES (23, 'smsProvider', 'huawei', 'string');
INSERT INTO `settings` VALUES (24, 'smsSignature', null, 'string');
INSERT INTO `settings` VALUES (25, 'smsEndpoint', null, 'string');
INSERT INTO `settings` VALUES (26, 'smsSender', null, 'string');
INSERT INTO `settings` VALUES (27, 'smsAppKey', null, 'string');
INSERT INTO `settings` VALUES (28, 'smsAppSecret', null, 'string');
INSERT INTO `settings` VALUES (29, 'smsCallbackUrl', null, 'string');
INSERT INTO `settings` VALUES (30, 'smsTemplateId', null, 'string');
INSERT INTO `settings` VALUES (31, 'dingdingAppKey', null, 'string');
INSERT INTO `settings` VALUES (32, 'dingdingAppSecret', null, 'string');
INSERT INTO `settings` VALUES (33, 'feishuAppId', null, 'string');
INSERT INTO `settings` VALUES (34, 'feishuAppSecret', null, 'string');
INSERT INTO `settings` VALUES (35, 'wechatCorpId', null, 'string');
INSERT INTO `settings` VALUES (36, 'wechatAgentId', null, 'int');
INSERT INTO `settings` VALUES (37, 'wechatSecret', null, 'string');
INSERT INTO `settings` VALUES (38, 'tokenExpiresTime', 12, 'int');
INSERT INTO `settings` VALUES (39, 'swagger','false', 'boolean');