package service

import (
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"golang.org/x/text/encoding/unicode"
	"ops-api/config"
	"ops-api/dao"
	"ops-api/model"
	"ops-api/utils"
	"strconv"
	"strings"
	"time"
)

var AD ad

type ad struct{}

type LDAPServer struct {
	Conn   *ldap.Conn
	Config LDAPConfig
}

type LDAPConfig struct {
	Addr             string
	BindUserDN       string
	BindUserPassword string
	SearchDN         string
}

// UserList 用户同步结构体，用于LDAP用户同步
type UserList struct {
	Name              string     `json:"name"`
	Username          string     `json:"username"`
	Password          string     `json:"password"`
	IsActive          bool       `json:"is_active"`
	PhoneNumber       string     `json:"phone_number"`
	Email             string     `json:"email"`
	UserFrom          string     `json:"user_from"`
	PasswordExpiredAt *time.Time `json:"password_expired_at"`
}

// Connect 建立LDAP连接
func (a *ad) Connect() (*LDAPServer, error) {

	var (
		host     = config.Conf.Settings["ldapAddress"].(string)
		userDn   = config.Conf.Settings["ldapBindDn"].(string)
		password = config.Conf.Settings["ldapBindPassword"].(string)
		searchDn = config.Conf.Settings["ldapSearchDn"].(string)
	)

	// 密码解密
	str, _ := utils.Decrypt(password)

	conf := LDAPConfig{
		Addr:             host,
		BindUserDN:       userDn,
		BindUserPassword: str,
		SearchDN:         searchDn,
	}

	conn, err := ldap.DialURL(conf.Addr, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
	if err != nil {
		return nil, err
	}
	_, err = conn.SimpleBind(&ldap.SimpleBindRequest{
		Username: conf.BindUserDN,
		Password: conf.BindUserPassword,
	})
	if err != nil {
		return nil, err
	}

	return &LDAPServer{Conn: conn, Config: conf}, nil
}

// LDAPUserSearch 根据用户名查找用户信息
func (a *ad) LDAPUserSearch(username string) (result *ldap.SearchResult, err error) {

	var (
		userAttribute = config.Conf.Settings["ldapFilterAttribute"].(string)
		searchDn      = config.Conf.Settings["ldapSearchDn"].(string)
	)

	// 建立LDAP连接
	l, err := a.Connect()
	if err != nil {
		return nil, err
	}

	// 查找用户
	searchDN := strings.Split(searchDn, "&")
	for _, dn := range searchDN {

		// 构建查找请求
		searchRequest := ldap.NewSearchRequest(
			dn,
			ldap.ScopeWholeSubtree,
			ldap.NeverDerefAliases,
			0,
			0,
			false,
			fmt.Sprintf("(&(objectClass=person)(%s=%s))", userAttribute, username),
			[]string{},
			nil,
		)

		// 执行查找
		searchResult, err := l.Conn.Search(searchRequest)
		if err != nil {
			return nil, err
		}

		// 如果没有找到用户，则继续在下一个DN查找
		if len(searchResult.Entries) == 0 {
			continue
		}

		// 返回用户信息
		return searchResult, nil
	}

	return nil, errors.New("用户不存在")
}

// LDAPUserAuthentication 用户认证
func (a *ad) LDAPUserAuthentication(username, password string) (result *ldap.SearchResult, err error) {

	// 建立LDAP连接
	l, err := a.Connect()
	if err != nil {
		return nil, err
	}

	// 获取用户信息
	searchResult, err := a.LDAPUserSearch(username)

	if err != nil {
		return nil, err
	}

	// 检查账号和密码是否过期
	userAttribute := config.Conf.Settings["ldapFilterAttribute"].(string)
	if userAttribute == "uid" {
		// 获取当前日期距离1970年1月1日之间的天数
		currentDays := daysSinceEpoch()

		// 检查账号是否过期
		shadowExpire := searchResult.Entries[0].GetAttributeValue("shadowExpire")
		if shadowExpire != "" {
			// 转换shadowExpire为整数
			shadowExpireInt, err := strconv.Atoi(shadowExpire)
			if err != nil {
				return nil, err
			}

			// 检查账户是否已过期
			if shadowExpireInt <= currentDays {
				return nil, errors.New("账号已过期，请联系管理员")
			}
		}

		// 检查密码是否过期
		shadowMax := searchResult.Entries[0].GetAttributeValue("shadowMax")
		shadowLastChange := searchResult.Entries[0].GetAttributeValue("shadowLastChange")
		shadowMaxInt, err := strconv.Atoi(shadowMax)
		if err != nil {
			return nil, err
		}
		shadowLastChangeInt, err := strconv.Atoi(shadowLastChange)
		if err != nil {
			return nil, err
		}

		// 计算密码过期日期
		passwordExpireDay := shadowLastChangeInt + shadowMaxInt
		if currentDays >= passwordExpireDay {
			return nil, errors.New("密码已过期，请重置密码")
		}
	}

	// 密码认证
	userDN := searchResult.Entries[0].DN
	err = l.Conn.Bind(userDN, password)
	if err != nil {
		return nil, errors.New("用户或密码错误")
	}

	// 返回用户信息
	return searchResult, nil
}

// LDAPUserResetPassword 重置用户密码
func (a *ad) LDAPUserResetPassword(username, password string) (err error) {
	// 建立LDAP连接
	l, err := a.Connect()
	if err != nil {
		return err
	}

	// 获取用户信息
	searchResult, err := a.LDAPUserSearch(username)
	if err != nil {
		return err
	}

	// 构建修改请求
	userDN := searchResult.Entries[0].DN
	req := ldap.NewModifyRequest(userDN, []ldap.Control{})

	// 密码修改
	var (
		passwordExpiredAt *time.Time
		userAttribute     = config.Conf.Settings["ldapFilterAttribute"].(string)
	)
	if userAttribute == "uid" {
		// 使用 SHA1 算法对密码进行哈希处理
		hash := sha1.New()
		hash.Write([]byte(password))
		digest := hash.Sum(nil)

		// 将哈希结果进行 Base64 编码
		encoded := base64.StdEncoding.EncodeToString(digest)

		// LDAP 用户修改密码
		req.Replace("userPassword", []string{fmt.Sprintf("{SHA}%s", encoded)})

		// 获取当前日期距离1970年1月1日之间的天数
		shadowLastChange := daysSinceEpoch()

		// 更新密码最后更改时间
		req.Replace("shadowLastChange", []string{strconv.Itoa(shadowLastChange)})

		// 获取密码过期时间并转换为数字
		passwordExpired := searchResult.Entries[0].GetAttributeValue("shadowMax")

		shadowLastChangeStr := strconv.Itoa(shadowLastChange)
		expiredAt, err := getPasswordExpiredAt(&shadowLastChangeStr, &passwordExpired)
		if err != nil {
			return err
		}
		passwordExpiredAt = expiredAt

		// 执行修改请求
		if err := l.Conn.Modify(req); err != nil {
			return err
		}

	} else {
		// 对密码进行utf16编码
		utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
		pwdEncoded, _ := utf16.NewEncoder().String("\"" + password + "\"")

		// Windows AD用户修改密码
		req.Replace("unicodePwd", []string{pwdEncoded})

		// 执行修改请求，注意：修改用户密码需要确保BindUserDN账号具备修改用户密码权限，以及需要使用ldaps方式连接，ldaps默认端口号为636，如：ldaps://192.168.200.13:636
		if err := l.Conn.Modify(req); err != nil {
			return err
		}

		// 获取账号过期时间并转换为数字
		userAccountControl := searchResult.Entries[0].GetAttributeValue("userAccountControl")
		userAccountControlInt, err := strconv.Atoi(userAccountControl)
		if err != nil {
			return err
		}

		if userAccountControlInt == 66050 || userAccountControlInt == 66048 {
			passwordExpiredAt = nil
		} else {
			// 获取最后修改密码时间，休眠1秒钟，确保可以从Windows AD中获取到最新的时间
			time.Sleep(1 * time.Second)
			pwdLastSet := searchResult.Entries[0].GetAttributeValue("pwdLastSet")

			// 获取密码过期日期
			expiredAt, err := getPasswordExpiredAt(&pwdLastSet, nil)
			if err != nil {
				return err
			}

			passwordExpiredAt = expiredAt
		}
	}

	// 指定用户查找条件
	conditions := map[string]interface{}{
		"username": username,
	}

	// 本地数据库查找对应用户信息
	user, err := dao.User.GetUser(conditions)
	if err != nil {
		return err
	}

	// 修改本地用户密码过期时间
	if err := dao.User.UpdateUserPasswordExpiredAt(user.ID, passwordExpiredAt); err != nil {
		return err
	}

	return nil
}

// LDAPUserSync 用户同步
func (a *ad) LDAPUserSync() (err error) {
	var (
		userList               []UserList
		createOrUpdateUserList []*model.AuthUser
		searchDn               = config.Conf.Settings["ldapSearchDn"].(string)
		userAttribute          = config.Conf.Settings["ldapFilterAttribute"].(string)
	)

	// 建立LDAP连接
	l, err := a.Connect()
	if err != nil {
		return err
	}

	// 获取所有用户
	searchDN := strings.Split(searchDn, "&")
	for _, dn := range searchDN {
		// 构建查找请求
		searchRequest := ldap.NewSearchRequest(
			dn, // 指定查找范围
			ldap.ScopeWholeSubtree,
			ldap.NeverDerefAliases,
			0,
			0,
			false,
			fmt.Sprintf("(objectClass=person)"), // 指定过滤条件：类型为用户
			[]string{},
			nil,
		)

		// 执行查找
		searchResult, err := l.Conn.Search(searchRequest)
		if err != nil {
			return err
		}

		// 获取查结果
		for _, value := range searchResult.Entries {

			// 判断用户状态，创建一个默认值
			var (
				isActive          bool
				passwordExpiredAt *time.Time
			)
			if userAttribute == "uid" {
				// 获取账号过期时间并转换为数字
				shadowExpire := value.GetAttributeValue("shadowExpire")
				shadowExpireInt, err := strconv.Atoi(shadowExpire)
				if err != nil {
					return err
				}

				// 判断是否是永不过期
				if shadowExpireInt == 99999 {
					isActive = true
				}

				// 获取当前日期距离1970年1月1日之间的天数，如果获取出的值大于当前天数则未过期
				currentDays := daysSinceEpoch()
				if shadowExpireInt > 0 && shadowExpireInt > currentDays {
					isActive = true
				}

				// 获取密码过期天数
				passwordExpired := value.GetAttributeValue("shadowMax")

				// 获取密码上次更改时间
				passwordLastChange := value.GetAttributeValue("shadowLastChange")

				// 获取密码到期时间
				expiredAt, err := getPasswordExpiredAt(&passwordLastChange, &passwordExpired)
				passwordExpiredAt = expiredAt
			} else {
				// 获取账号过期时间并转换为数字
				userAccountControl := value.GetAttributeValue("userAccountControl")
				userAccountControlInt, err := strconv.Atoi(userAccountControl)
				if err != nil {
					return err
				}

				// userAccountControl的值说明：514 禁用，512 启用，66050 禁用+密码永不过期，66048 启用+密码永不过期
				if userAccountControlInt == 512 || userAccountControlInt == 66048 {
					isActive = true
				}

				if userAccountControlInt == 66050 || userAccountControlInt == 66048 {
					passwordExpiredAt = nil
				} else {
					// 获取最后一次修改密码日期
					pwdLastSet := value.GetAttributeValue("pwdLastSet")

					// 获取密码过期日期
					expiredAt, err := getPasswordExpiredAt(&pwdLastSet, nil)
					if err != nil {
						return err
					}
					passwordExpiredAt = expiredAt
				}
			}

			// 获取用户信息
			userInfo := &UserList{
				Name:              value.GetAttributeValue("cn"),
				Username:          value.GetAttributeValue(userAttribute),
				Password:          "",
				IsActive:          isActive,
				PhoneNumber:       value.GetAttributeValue("mobile"),
				Email:             value.GetAttributeValue("mail"),
				UserFrom:          "LDAP",
				PasswordExpiredAt: passwordExpiredAt,
			}
			// 将用户信息追加到结构体
			userList = append(userList, *userInfo)
		}
	}

	// 同步所有用户
	for _, user := range userList {
		createOrUpdateUserList = append(createOrUpdateUserList, &model.AuthUser{
			Username:          user.Username,
			Name:              user.Name,
			Email:             user.Email,
			Password:          user.Password,
			IsActive:          user.IsActive,
			PhoneNumber:       user.PhoneNumber,
			UserFrom:          user.UserFrom,
			PasswordExpiredAt: user.PasswordExpiredAt,
		})
	}
	return dao.User.SyncUsers(createOrUpdateUserList)
}

// getPasswordExpiredAt 获取密码过期时间
func getPasswordExpiredAt(lastChangeString, passwordExpiredString *string) (passwordExpiredAt *time.Time, err error) {

	var (
		userAttribute = config.Conf.Settings["ldapFilterAttribute"].(string)
		passwordAge   = config.Conf.Settings["ldapUserPasswordExpireDays"].(int)
	)

	if userAttribute == "uid" {
		// OpenLDAP的获取方法

		// 获取密码最后更改时间并转换为数字
		lastChangeInt, err := strconv.Atoi(*lastChangeString)
		if err != nil {
			return nil, err
		}

		// 获取密码过期时间并转换为数字
		passwordExpiredInt, err := strconv.Atoi(*passwordExpiredString)
		if err != nil {
			return nil, err
		}

		// 如果为99999则表示密码永不过期
		if passwordExpiredInt == 99999 {
			return nil, nil
		}

		// 计算上次更改密码的日期
		lastChangeDate := time.Unix(int64(lastChangeInt)*24*60*60, 0)

		// 计算密码过期日期
		expirationDate := lastChangeDate.Add(time.Duration(passwordExpiredInt) * 24 * time.Hour)

		return &expirationDate, nil
	} else {
		// Windows AD的获取方法

		// 密码过期时间
		maxPasswordAge := passwordAge

		// 将文件时间转换为Unix时间
		lastChangeInt, err := strconv.ParseInt(*lastChangeString, 10, 64)
		if err != nil {
			return nil, err
		}
		lastPasswordChangeTime := NtToUnix(lastChangeInt)

		// 计算密码过期时间
		expirationDuration := time.Duration(maxPasswordAge) * 24 * time.Hour
		passwordExpirationTime := lastPasswordChangeTime.Add(expirationDuration)
		return &passwordExpirationTime, nil
	}
}

// 计算自1970年1月1日以来的天数
func daysSinceEpoch() int {
	// 获取当前时间
	now := time.Now()
	// 获取自1970年1月1日起的时间差（秒）
	elapsed := now.UTC().Sub(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))
	// 转换为天数并返回
	return int(elapsed.Hours() / 24)
}

// NtToUnix Window NT时间转换为Unix时间
func NtToUnix(ntTime int64) (unixTime time.Time) {
	ntTime = (ntTime - 1.1644473600125e+17) / 1e+7
	return time.Unix(ntTime, 0)
}
