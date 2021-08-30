package exception

import (
	"inochat/client/res"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	ErrSign         = &SignErr{Msg: "sign error."}
	ErrParam        = &ParamErr{Msg: "param error."}
	ErrUnAuthorized = &AuthorizeErr{Msg: "unauthrized."}
)

// 捕获所有的错误信息
func ExceptionHandle(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.Abort()
			// switch e := err.(type) {
			// case *ParamErr:
			// 	log.Errorf("%v", err)
			// 	c.JSON(res.PARAMETERERR, res.Result(res.PARAMETERERR, e.Error()))
			// case *AuthorizeErr:
			// 	log.Errorf("%v", err)
			// 	c.JSON(res.UNAUTHORIZATION, res.Result(res.UNAUTHORIZATION, e.Error()))
			// case *SignErr:
			// 	log.Errorf("%v", err)
			// 	c.JSON(res.PROHIBITED, res.Result(res.PROHIBITED, e.Error()))
			// case *SystemErr:
			// 	log.Errorf("%v", err)
			// 	c.JSON(res.FAILED, res.Err("网络异常"))
			// default:
			// 	log.Errorf("%v", err)
			// 	c.JSON(res.FAILED, res.Err("网络异常"))
			// }
			log.Errorf("%v", err)
			c.JSON(res.FAILED, res.Err("网络异常"))
		}
	}()
	c.Next()
}

func SysErr(msg string) error {
	return &SystemErr{Msg: msg}
}

type SystemErr struct {
	Msg string
}

func (ex *SystemErr) Error() string {
	return ex.Msg
}

type ParamErr struct {
	Msg string
}

func (ex *ParamErr) Error() string {
	return ex.Msg
}

type AuthorizeErr struct {
	Msg string
}

func (ex *AuthorizeErr) Error() string {
	return ex.Msg
}

type SignErr struct {
	Msg string
}

func (ex *SignErr) Error() string {
	return ex.Msg
}
