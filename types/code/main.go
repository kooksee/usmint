package code

var (
	Ok = errs(00000, "ok")

	// Generic error definition, 10000-10999
	ErrInternal = errs(10000, "发生未知错误，请稍后再试～")
	ErrInMaintain = errs(10001, "服务器在维护中，请稍后再试～")
	ErrRequestFormat = errs(10002, "网络请求错误")
	ErrApiDeprecated = errs(10003, "当前版本过低，需要升级咯～")
	ErrSessionExpired = errs(10004, "用户已从其他设备登录")
	ErrInBlacklist = errs(10005, "当前用户已无法访问服务")
	ErrTempBlock = errs(10006, "您已暂时被禁止访问，请稍后再试")


	// app error definition, 12000-12999
	ErrTransactionDecode = errs(12001, "transaction decode error")
	ErrSignAccountNotExist = errs(12002, "ok")
	ErrUnknownMathod = errs(12003, "未知方法调用")
	ErrJsonDecode = errs(12004, "ok")
	ErrHexDecode = errs(12005, "convert hex string to hex bytes error")
	ErrDbDecode = errs(12006, "db error")
	ErrVerify = errs(12007, "签名验证失败")
	ErrLicenseSave = errs(12008, "license save error")
	ErrTransactionVerify = errs(12009, "transaction verify error")
	ErrValidatorAdd = errs(12010, "validator add error")



	// User/Auth error definition, 11000-11099
	ErrNicknameExists = errs(11000, "昵称已存在")
	ErrUserNotFound = errs(11001, "用户不存在")

	ErrUserExists = errs(11002, "用户已存在")
	ErrPasswordWrong = errs(11003, "密码错误")
	ErrUserBlacked = errs(11004, "当前用户已进入黑名单")
	ErrMobileUnmodified = errs(11005, "手机号码无法被更改")
	ErrMobileFormat = errs(11006, "手机号码格式错误")
	ErrMobileNotSet = errs(11007, "未绑定手机号码")
	ErrPasswordLength = errs(11008, "密码长度错误")
	ErrPasswordFormat = errs(11009, "密码格式错误")
	ErrNicknameStopWord = errs(11010, "您的昵称包含敏感词")

	ErrTokenNotFound = errs(11011, "网络请求错误")
	ErrTokenFormat = errs(11012, "网络请求错误")
	ErrTokenExpired = errs(11013, "网络请求错误")
	ErrTokenRefreshExpired = errs(11014, "网络请求错误")

	ErrRoleFormat = errs(11015, "请求权限错误")
	ErrRoleWrong = errs(11016, "请求权限错误")
)

