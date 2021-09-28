package response

var CodeMessage = map[int]string{
	SUCCESS:         "成功",
	FAIL:            "请求失败",
	AuthError:       "用户认证失败",
	ERROR:           "系统错误",
	UserExist:       "用户已存在",
	UserNotExist:    "用户不存在",
	PasswordError:   "密码错误",
	ParamError:      "参数提交错误",
	PermissionError: "没有权限",
}

func GetMessage(code int) string {
	msg, ok := CodeMessage[code]
	if ok {
		return msg
	}
	return CodeMessage[ERROR]
}
