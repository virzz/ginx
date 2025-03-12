package code

type SystemCode int

const (
	OK       SystemCode = iota
	Client              // 客户请求错误
	Internal            // 内部错误
	Remote              // 远程错误
)

type BizCode int

const (
	Auth    BizCode = iota + 1 // Auth
	Captcha                    // Captcha
	Record                     // Record(Database)
	File                       // File
)

var (
	Success = New(OK, 0, 0, "成功")
)

var (
	UnknownErr    = New(Internal, 0, 0, "系统错误")
	UnImplemented = New(Internal, 0, 1, "未实现")
	NotFound      = New(Internal, 0, 2, "资源未找到")

	RecordUnknown  = New(Internal, Record, 0, "未知数据库错误")
	RecordExists   = New(Internal, Record, 1, "记录已存在")
	RecordUpdate   = New(Internal, Record, 2, "记录更新失败")
	RecordDelete   = New(Internal, Record, 3, "记录删除失败")
	RecordInsert   = New(Internal, Record, 4, "记录插入失败")
	RecordTruncate = New(Internal, Record, 5, "记录清空失败")
)

// Request
var (
	BadRequest   = New(Client, 0, 0, "请求不合法")
	RequestLimit = New(Client, 0, 1, "请求过于频繁")
	ParamInvalid = New(Client, 0, 2, "参数不合法")
)

// Auth
var (
	Unauthorized     = New(Client, Auth, 1, "未授权")  // Token Not Found in Request
	Forbidden        = New(Client, Auth, 2, "权限不足") // Permission Denied || Banned
	TokenInvalid     = New(Client, Auth, 3, "令牌无效")
	TokenExpired     = New(Client, Auth, 4, "令牌过期") // Token Not Found in Store
	PasswordInvalid  = New(Client, Auth, 5, "密码错误")
	PasswordWeak     = New(Client, Auth, 6, "密码强度不足")
	PasswordMatch    = New(Client, Auth, 7, "密码不匹配")
	PasswordReset    = New(Client, Auth, 8, "密码重置失败")
	EmailDisable     = New(Client, Auth, 9, "邮箱未开启")
	EmailCodeInvalid = New(Client, Auth, 10, "邮箱验证码无效")

	TokenCreate     = New(Internal, Auth, 1, "令牌创建失败")
	TokenDestory    = New(Internal, Auth, 2, "令牌销毁失败")
	RegisterExists  = New(Internal, Auth, 3, "用户已存在")
	RegisterFailed  = New(Internal, Auth, 4, "用户注册失败")
	RegisterDisable = New(Internal, Auth, 5, "用户注册已关闭")

	EmailCodeSend = New(Remote, Auth, 1, "邮箱验证码发送失败")
)

// Captcha
var (
	CaptchaInvalid = New(Client, Captcha, 1, "验证码错误")

	CaptchaGenerate = New(Internal, Captcha, 1, "验证码生成失败")
)

// File
var (
	FileSize = New(Client, File, 1, "文件大小超出限制")
	FileType = New(Client, File, 2, "文件类型不支持")

	FileUpload = New(Internal, File, 1, "文件上传失败")
	FileOpen   = New(Internal, File, 2, "文件打开失败")
	FileParse  = New(Internal, File, 3, "文件解析失败")
	FileExport = New(Internal, File, 4, "文件导出失败")

	FileUploadS3 = New(Remote, File, 1, "文件上传失败")
)
