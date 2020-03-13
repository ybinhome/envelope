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
	return validate
}

func Translator() ut.Translator {
	return translator
}

type ValidatorStarter struct {
	infra.BaseStarter
}

func (v *ValidatorStarter) Init(ctx infra.StarterContext) {
	// 实例化 validate 对象
	validate := validator.New()

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
