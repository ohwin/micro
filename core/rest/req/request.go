package req

type PageInfo struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"pageSize" json:"pageSize"`
}

type SmsVerifyReq struct {
	Phone string `json:"phone" form:"phone"`
	Code  string `json:"code" form:"code"`
}
