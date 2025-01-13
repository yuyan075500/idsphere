package db

import (
	"fmt"
	"github.com/wonderivan/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
	"ops-api/config"
	"ops-api/global"
	"ops-api/model"
	"os"
	"strings"
	"time"
)

// MySQLInit MySQL初始化
func MySQLInit() error {

	// 组装数据库连接配置
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Conf.MySQL.User,
		config.Conf.MySQL.Password,
		config.Conf.MySQL.Host,
		config.Conf.MySQL.Port,
		config.Conf.MySQL.DB,
	)

	// 建立数据库连接，并生成*gorm.DB对象
	client, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger2.Default.LogMode(logger2.Silent),
	})
	if err != nil {
		return err
	}

	// 表迁移
	_ = client.AutoMigrate(
		&model.AuthUser{},
		&model.AuthGroup{},
		&model.SiteGroup{},
		&model.Tag{},
		&model.Site{},
		&model.Menu{},
		&model.SubMenu{},
		&model.SystemPath{},
		&model.LogSMS{},
		&model.LogLogin{},
		&model.LogOplog{},
		&model.SsoOAuthTicket{},
		&model.SsoCASTicket{},
		&model.ScheduledTask{},
		&model.ScheduledTaskExecLog{},
		&model.Account{},
		&model.Settings{},
	)

	// 设置数据库连接池
	DB, _ := client.DB()
	DB.SetMaxIdleConns(config.Conf.MySQL.MaxIdleConns)
	DB.SetMaxOpenConns(config.Conf.MySQL.MaxOpenConns)
	DB.SetConnMaxLifetime(time.Duration(config.Conf.MySQL.MaxLifeTime) * time.Second)

	global.MySQLClient = client

	// 初始化站点信息
	if err := initializeSites(client); err != nil {
		return err
	}

	// 初始化基础数据
	if err := initializeSQL(client); err != nil {
		return err
	}

	// 初始化配置信息
	if err := InitConfig(client); err != nil {
		return err
	}

	// 创建超级用户
	if err := createSuperUser(client); err != nil {
		return err
	}

	// 初始化定时任务
	if err := InitializeScheduledTask(client); err != nil {
		return err
	}

	logger.Info("MySQL客户端及数据初始化成功.")

	return nil
}

// 创建超级用户
func createSuperUser(client *gorm.DB) error {
	user := model.AuthUser{
		Name:     "管理员",
		Username: "admin",
		IsActive: true,
		Password: "admin@123...",
	}

	result := client.FirstOrCreate(&user, model.AuthUser{Username: "admin"})

	if result.RowsAffected == 0 {
		logger.Warn("admin用户已存在.")
	} else {
		logger.Info("admin用户创建成功.")
	}
	return nil
}

// 初始化站点信息
func initializeSites(client *gorm.DB) error {
	var count int64
	if err := client.Model(&model.SiteGroup{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	var siteGroup model.SiteGroup
	if err := client.FirstOrCreate(&siteGroup, model.SiteGroup{Name: "系统默认（可以删除，删除前请确保存在至少1个分组，否则系统在下次重启时又将自动创建）"}).Error; err != nil {
		return err
	}

	sites := []model.Site{
		{
			Name:        "密码重置",
			Description: "密码自助更改平台，支持本地、Windows AD 和 OpenLDAP 用户密码修改",
			Address:     fmt.Sprintf("/reset_password"),
			SSO:         false,
			SiteGroupID: siteGroup.ID,
		},
		{
			Name:        "接口文档",
			Description: "IDSphere 统一认证平台 Swagger 接口文档",
			Address:     fmt.Sprintf("/swagger/index.html"),
			SSO:         false,
			SiteGroupID: siteGroup.ID,
		},
		{
			Name:        "站点导航",
			Description: "IDSphere 统一认证平台站点导航页",
			Address:     fmt.Sprintf("/sites"),
			SSO:         false,
			SiteGroupID: siteGroup.ID,
		},
		{
			Name:        "IDP 元数据",
			Description: "IDP 元数据配置文件接口，接入 SAML2 SP 使用",
			Address:     fmt.Sprintf("/api/v1/sso/saml/metadata"),
			SSO:         false,
			SiteGroupID: siteGroup.ID,
		},
		{
			Name:        "OIDC 配置",
			Description: "OIDC 配置信息接口，接入 OIDC 客户端使用",
			Address:     fmt.Sprintf("/.well-known/openid-configuration"),
			SSO:         false,
			SiteGroupID: siteGroup.ID,
		},
	}

	for _, site := range sites {
		client.FirstOrCreate(&site, model.Site{Address: site.Address})
	}

	return nil
}

// 初始化SQL
func initializeSQL(client *gorm.DB) error {

	var count int64
	if err := client.Model(&model.Menu{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	// 读取初始化SQL
	content, err := os.ReadFile("db/data.sql")
	if err != nil {
		return err
	}

	// 按分号分割SQL语句
	queries := strings.Split(string(content), ";")
	for _, query := range queries {
		// 去除前后空格
		query = strings.TrimSpace(query)
		if query == "" {
			// 跳过空语句
			continue
		}

		// 执行SQL语句
		if err := client.Exec(query).Error; err != nil {
			return err
		}
	}

	return nil
}

// InitializeScheduledTask 初始化定时任务
func InitializeScheduledTask(client *gorm.DB) error {

	var count int64
	if err := client.Model(&model.ScheduledTask{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	tasks := []model.ScheduledTask{
		{
			Name:          "密码过期通知",
			Type:          2,
			CronExpr:      "0 8 * * *",
			BuiltInMethod: "password_expire_notify",
			Enabled:       false,
		},
		{
			Name:          "用户同步",
			Type:          2,
			CronExpr:      "* 5 * * *",
			BuiltInMethod: "user_sync",
			Enabled:       false,
		},
	}

	for _, task := range tasks {
		client.FirstOrCreate(&task, model.ScheduledTask{BuiltInMethod: task.BuiltInMethod})
	}

	return nil
}

// InitConfig 初始化配置
func InitConfig(client *gorm.DB) error {

	// 加载所有设置项
	settings := make(map[string]interface{})
	var configItems []model.Settings
	if err := client.Find(&configItems).Error; err != nil {
		return err
	}

	// 保存配置项
	for _, item := range configItems {
		if err := item.ParseValue(); err != nil {
			continue
		}
		settings[item.Key] = item.ParsedValue
	}

	config.Conf.Settings = settings
	return nil
}
