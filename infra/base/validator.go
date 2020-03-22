package base

import (
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/sirupsen/logrus"
	"github.com/ybinhome/envelope/infra"
	"gopkg.in/go-playground/validator.v9"
	vtzh "gopkg.in/go-playground/validator.v9/translations/zh"
)

var validate *validator.Validate
var translator ut.Translator

func Validate() *validator.Validate {
	Check(validate)
	return validate
}

func Translator() ut.Translator {
	Check(validate)
	return translator
}

type ValidatorStarter struct {
	infra.BaseStarter
}

func (v *ValidatorStarter) Init(ctx infra.StarterContext) {
	// 实例化 validate 对象
	validate = validator.New()

	// 创建 validate 消息输出中文翻译器
	//   通过 locales 的中文包创建中文翻译器
	cn := zh.New()
	//   通过 universal-translator 包的通用翻译器来转换，并获取一个中文翻译器
	uni := ut.New(cn, cn)
	var found bool
	translator, found = uni.GetTranslator("zh")
	if found {
		//  将翻译器注册到 validate 中，通过 validate 框架中的 Translations 包下的语言包来实现，中文包和 locales 中的中文包名冲突，此处需要手动导入重命名
		err := vtzh.RegisterDefaultTranslations(validate, translator)
		if err != nil {
			logrus.Error(err)
		}
	} else {
		logrus.Errorf("Not found translator: zh")
	}
}

// 验证输入参数是否合法
func ValidateStruct(s interface{}) (err error) {
	err = Validate().Struct(s)
	if err != nil {
		// 使用类型断言来判断 err 的类型，如果是 *validator.InvalidValidationError 类型，则传入的参数 &dto 为空，将自定义的错误信息和 err 输出到日志
		_, ok := err.(*validator.InvalidValidationError)
		if ok {
			logrus.Error("验证错误", err)
		}

		// 使用类型断言来判断 err 的类型，如果是 validator.ValidationErrors 类型，则传入的参数 &dto 中字段有误，同时使用循环来便利 errs 切片，通过转换器中文输出那些字段有问题到日志
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			for _, e := range errs {
				logrus.Error(e.Translate(Translator()))
			}
		}

		// 最后向用户范围一个空的 *services.AccountDTO 和错误信息
		return err
	}
	return nil
}
