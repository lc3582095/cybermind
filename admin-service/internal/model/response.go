package model

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageResponse 分页响应结构
type PageResponse struct {
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}

// LoginResponse 登录响应结构
type LoginResponse struct {
	Token  string `json:"token"`
	Admin  Admin  `json:"admin"`
} 