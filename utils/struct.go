package utils

// excludedFields 不记录日志的字段
var excludedFields = map[string]struct{}{
	"password":          {},
	"re_password":       {},
	"client_id":         {},
	"client_secret":     {},
	"certificate":       {},
	"mfa_code":          {},
	"DeletedAt":         {},
	"ldapBindPassword":  {},
	"mailPassword":      {},
	"smsAppSecret":      {},
	"dingdingAppSecret": {},
	"feishuAppSecret":   {},
	"wechatSecret":      {},
	"access_key":        {},
	"secret_key":        {},
	"iam_password":      {},
}

// FilterFields 递归过滤敏感字段
func FilterFields(data map[string]interface{}) {
	for key, value := range data {
		// 检查是否为敏感字段
		if _, ok := excludedFields[key]; ok {
			// 删除敏感字段
			delete(data, key)
		} else if nestedMap, ok := value.(map[string]interface{}); ok {
			// 递归处理嵌套 map
			FilterFields(nestedMap)
		} else if nestedArray, ok := value.([]interface{}); ok {
			// 递归处理数组中的map
			for _, item := range nestedArray {
				if itemMap, ok := item.(map[string]interface{}); ok {
					FilterFields(itemMap)
				}
			}
		}
	}
}
