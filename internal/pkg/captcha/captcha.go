package captcha

import (
	"fmt"
	"github.com/mojocn/base64Captcha"
)

var store = base64Captcha.DefaultMemStore

func GenerateCaptcha() (string, string) {
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, store)
	id, base64, err := cp.Generate()
	if err != nil {
		fmt.Printf("generate captcha fail, err: %v\n", err)
	}
	return id, base64
}

func ParseCaptcha(id, requestId string) bool {
	return store.Verify(id, requestId, true)
}
