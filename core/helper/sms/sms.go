package sms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ohwin/micro/core/config"
	"github.com/ohwin/micro/core/errx"
	"github.com/ohwin/micro/core/helper/binding"
	"github.com/ohwin/micro/core/log"
	"github.com/ohwin/micro/core/rest/req"
	"github.com/ohwin/micro/core/rest/response"
	"github.com/ohwin/micro/core/store"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	api "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type Template struct {
	Code string `json:"code"`
}

func CreateClient() (client *api.Client, err error) {
	aliyun := config.App.Aliyun

	c := &openapi.Config{
		AccessKeyId:     tea.String(aliyun.AccessKeyId),
		AccessKeySecret: tea.String(aliyun.AccessKeySecret),
	}

	c.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	return api.NewClient(c)
}

func Send(phone string, template interface{}, prefixes ...string) (err error) {
	client, err := CreateClient()

	if err != nil {
		return err
	}

	sms := config.App.Aliyun.Sms
	param, err := json.Marshal(template)
	runtime := &util.RuntimeOptions{}
	sendSmsRequest := &api.SendSmsRequest{
		PhoneNumbers:  tea.String(phone),
		SignName:      tea.String(sms.SignName),
		TemplateCode:  tea.String(sms.TemplateCode),
		TemplateParam: tea.String(string(param)),
	}

	resp, err := client.SendSmsWithOptions(sendSmsRequest, runtime)
	if err != nil || *resp.Body.Code != "OK" {
		log.Error(fmt.Sprintf("send sms error:[code: %d][messages:%s]", *resp.StatusCode, *resp.Body.Message))
		return err
	}

	key := generateKey(phone, prefixes...)
	_ = save(key, param)

	return nil
}

func VerifyMiddleware(ctx *gin.Context) {
	req := new(req.SmsVerifyReq)

	err := binding.Auto(ctx, req, binding.JSON)
	if err != nil {
		return
	}

	template := Template{Code: req.Code}
	if !verify(req.Phone, template) {
		ctx.Abort()
		response.Fail(ctx, errx.ErrVerify)
		return
	}

	ctx.Next()
}

func verify(phone string, template interface{}, prefixes ...string) bool {
	param, _ := json.Marshal(template)
	key := generateKey(phone, prefixes...)

	val, err := get(key)
	if err != nil {
		return false
	}

	if bytes.Equal(param, val) {
		del(key)
		return true
	}

	return false
}

func save(key string, code []byte) (err error) {
	cache := store.RDB
	expire := config.App.Aliyun.Sms.Expire
	_, err = cache.Set(context.Background(), key, code, time.Second*expire).Result()
	return err
}

// 获取key
func get(key string) (val []byte, err error) {
	cache := store.RDB

	result, err := cache.Get(context.Background(), key).Bytes()
	if err != nil {
		return nil, err
	}

	return result, nil
}

// 删除key
func del(key string) {
	cache := store.RDB
	_ = cache.Del(context.Background(), key)
}

// 生成Key
func generateKey(phone string, prefixes ...string) (key string) {
	prefixes = append([]string{"sms"}, prefixes...)
	prefix := strings.Join(prefixes, ":")
	return fmt.Sprintf("%s:phone:%s", prefix, phone)
}
