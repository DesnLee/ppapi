package controller_helper

import (
	"fmt"

	"ppapi.desnlee.com/internal/model"
)

// ValidateCreateTagRequestBody 验证创建标签请求体
func ValidateCreateTagRequestBody(b *model.CreateTagRequestBody) error {
	if err := validateKind(string(b.Kind)); err != nil {
		return err
	}
	if b.Name == "" {
		return fmt.Errorf("标签名称不能为空")
	}
	if b.Sign == "" {
		return fmt.Errorf("标签图标不能为空")
	}
	return nil
}
