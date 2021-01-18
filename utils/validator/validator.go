package validator

import (
	"GopherBlog/utils"
	"errors"
	"github.com/go-playground/locales/zh_Hans_CN"
	unTrans "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
)

// 自定义验证器
func Validate(data interface{}) (string, error) {
	// 1.声明一个通用的翻译器（针对各国语言都可以，所以是通用的）
	// 2.unTrans.New()第一个参数是一个备用的翻译器，第二个参数是一个不定参数，表示设置该通用翻译器所支持的语言翻译器
	//
	// 例如：unTrans.New(en.New(), zh_Hans_CN.New(), zh_Hans_HK.New(), en_US.New())
	// 表示设置通用翻译器支持汉语(中国)、汉语（香港）、英语（美国）三种语言，然后备用（fallback）语言为英语（en.New()）
	//
	// 3.设置备用语言的用处：当调用uni.GetTranslator("zh_Hans_TW")时，首先会去通用翻译器中查找是否有“zh_Hans_TW”所对应的翻译器，
	// 找到了则返回，如果没有找到，则使用设置的备用翻译器（即en.New()）

	// 向通用翻译器设置一个备用翻译器为：汉语（中国）的翻译器，因为没有设置所支持的翻译器，所以调用GetTranslator("xxxx")时会一直返回备用翻译器
	uni := unTrans.New(zh_Hans_CN.New())
	// 1.从通用翻译器中获取zh_Hans_CN翻译器，这里可以根据实际情况修改，例如传入一个local参数，当local为"en"时返回英文翻译器，
	// 当local为"zh_Hans_CN"则返回中文翻译器
	//
	// 2.第二个参数返回一个bool值，如果在通用翻译器中找到了传入参数所对应的翻译器，则返回true，否则返回false
	trans, _ := uni.GetTranslator("zh_Hans_CN")

	// 创建验证器实例
	validate := validator.New()
	// TODO:注册默认翻译器, 如果设置为英文会怎么样？？
	err := zhTrans.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		utils.Logger.Error("验证器设置翻译器失败：", err)
		return "", err
	}
	// TODO:注册一个获取label tag的方法???
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		return label
	})

	err = validate.Struct(data)
	if err != nil {
		// 获取validator.ValidationErrors类型的errors
		for _, v := range err.(validator.ValidationErrors) {
			// 将错误信息翻译成对应的语言, 这里应该是只返回字段错误中的其中一个？？？？
			return v.Translate(trans), errors.New(v.Error())
		}
	}
	return "", nil
}
