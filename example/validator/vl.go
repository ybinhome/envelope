package main

import (
	"fmt"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	vtzh "gopkg.in/go-playground/validator.v9/translations/zh"
)

// 创建带 validate 限制条件的结构体，字段 tag 必须以 validate 开头，required 表示不能为空，gte 表示 >，lte 表示 <，email 表示必须为邮箱格式
type User struct {
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
	Age       uint8  `validate:"gte=0,lte=130"`
	Email     string `validate:"required,email"`
}

func main() {
	// 实例化 User 对象
	user := &User{
		FirstName: "firstName",
		LastName:  "lastName",
		Age:       136,
		Email:     "vl163.com",
	}

	// 实例化一个 validate 对象
	validate := validator.New()

	// 默认 validate 的报错输出不友好，创建国际化的消息翻译器
	//  创建中文翻译器
	cn := zh.New()
	//  通过通用翻译器来转换，并获取一个中文翻译器
	uni := ut.New(cn, cn)
	transtator, found := uni.GetTranslator("zh")
	if found {
		// 将翻译器注册到 validate 中，通过 validate 框架中的 Translations 包下的语言包来实现，中文包和 locales 中的中文包名冲突，此处需要手动导入重命名
		err := vtzh.RegisterDefaultTranslations(validate, transtator)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("not found")
	}

	// 使用 validate 验证结构体实例化对象 user 是否合法
	err := validate.Struct(user)
	if err != nil {
		_, ok := err.(*validator.InvalidValidationError)
		if ok {
			fmt.Println(err)
		}
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			// 默认的错误输出方法
			// fmt.Println(errs)

			// 中文翻译器
			for _, err := range errs {
				fmt.Println(err.Translate(transtator))
			}
		}
	}
}
