package validator

import (
	"GopherBlog/constant"
	"GopherBlog/utils"
	"github.com/go-playground/locales/zh_Hans_CN"
	unTrans "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
)

// Validate 自定义验证器
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
	//trans, _ := uni.GetTranslator("zh_Hans_CN")
	trans, _ := uni.GetTranslator("zh_Hans_CN")

	// 创建验证器实例
	validate := validator.New()
	// 为内置标签注册一组默认的翻译器，可以理解zhTrans表示内置标签的中文翻译,
	// 而RegisterDefaultTranslations的作用是为翻译器trans加上内置标签的中文翻译
	err := zhTrans.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		utils.Logger.Error(constant.ConvertForLog(constant.SetValidatorError), err)
	}
	// 注册一个获取label tag的方法，这样返回的验证器校验后的字段错误信息中会使用字段的label名称
	// 例如： User模型中的username的label为“用户名”，则如果该字段设置为长度在4到6之间，则打印信息为：用户名长度必须至少为4个字符
	// 而不是username长度必须至少4个字符
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		return label
	})

	err = validate.Struct(data)
	if err != nil {
		// 获取validator.ValidationErrors类型的errors
		for _, v := range err.(validator.ValidationErrors) {
			// 将错误信息翻译成对应的语言, 这里应该是只返回字段错误中的其中一个？？？？
			return v.Translate(trans), err
		}
	}
	return "", nil
}
