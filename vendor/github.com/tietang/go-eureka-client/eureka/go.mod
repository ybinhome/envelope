module github.com/tietang/go-eureka-client/eureka

go 1.12

replace (
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190325154230-a5d413f7728c
	golang.org/x/net => github.com/golang/net v0.0.0-20190327025741-74e053c68e29
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190227155943-e225da77a7e6
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190322080309-f49334f85ddc
	golang.org/x/text => github.com/golang/text v0.3.0
)

require (
	github.com/go-ini/ini v1.42.0 // indirect
	github.com/sirupsen/logrus v1.4.0
	github.com/smartystreets/goconvey v0.0.0-20190306220146-200a235640ff
	github.com/tietang/go-utils v0.0.0-20190308094824-9e17fa5e3788
	github.com/tietang/props v2.1.0+incompatible
	github.com/valyala/fasttemplate v1.0.1 // indirect
	gopkg.in/ini.v1 v1.42.0 // indirect
	gopkg.in/yaml.v2 v2.2.2
)
