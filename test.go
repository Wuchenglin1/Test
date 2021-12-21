package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"time"
)

var secretKey = []byte("SecretKey") //模拟私钥

func main() {
	engine := gin.Default()

	routerGroup := engine.Group("")
	{
		routerGroup.POST("/register", Register)
		routerGroup.POST("/login", Login)
		routerGroup.GET("", Author, Hello) //查看信息需要jwt鉴权
	}

	err := engine.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}

var user struct {
	UserName string `json:"userName"`
	Pwd      string `json:"pwd"`
}

var userMap = make(map[string]string)

type Claims struct {
	UserName string `json:"userName"`
	jwt.StandardClaims
}

func Hello(c *gin.Context) {
	c.JSON(200, "你好！")
}

//Register 注册用户
func Register(c *gin.Context) {
	m := userMap
	u := user
	u.UserName = c.PostForm("userName")
	u.Pwd = c.PostForm("password")
	c.JSON(200, "注册成功！")
	m[u.UserName] = u.Pwd
	fmt.Println(m) //打印注册的注册信息
}

//Login 登录用户
func Login(c *gin.Context) {
	u := user
	u.UserName = c.PostForm("userName")
	u.Pwd = c.PostForm("password")
	if u.UserName == "" || u.Pwd == "" {
		c.JSON(403, "error:账号或密码不能为空")
		c.Abort()
		return
	}
	if u.Pwd == userMap[u.UserName] {
		fmt.Println(u)
		j, err := SetAdminJWT(u.UserName)
		c.SetCookie("userName", u.UserName, 3600, "/", "", false, true)
		c.SetCookie("jwt", j, 3600, "/", "", false, true)
		c.JSON(200, "恭喜您登录成功")
		if err != nil {
			c.JSON(403, gin.H{
				"error:": err,
			})
			return
		}
	} else {
		c.JSON(403, "账号不存在或密码错误！")
		return
	}
}

func SetAdminJWT(userName string) (string, error) {
	c := Claims{
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "wuchenglin",
			ExpiresAt: time.Now().Add(time.Second * 10).Unix(), //为了方便测试，签证有效期改为10s
		},
	}
	//创建一个token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	//用secretKey对token进行签名返回一个s字符串
	s, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Println(err)
		return s, err
	}
	return s, nil
}

//Author jwt鉴权
func Author(c *gin.Context) {
	userName, _ := c.Cookie("userName")
	//完善payload信息，将userName信息填充进去
	claim := Claims{
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "wuchenglin",
			ExpiresAt: time.Now().Add(time.Second * 10).Unix(),
		},
	}
	//登录之后获取存储在cookie中的jwt信息
	t, err := c.Cookie("jwt")
	if err != nil {
		c.JSON(403, "您还没有登录！")
		c.Abort()
		return
	}
	token, err1 := jwt.ParseWithClaims(t, &claim, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	//token.Valid(bool)检验签名是否有效
	if token.Valid {
		c.JSON(200, gin.H{
			"您好！": claim.UserName,
		})
	} else if ve, ok := err1.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			//ValidationErrorMalformed验证token是否为畸形
			//token的格式错误,
			c.JSON(403, "请输个像样的token")
			c.Abort()
			return
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			//jwt.ValidationErrorExpired是验证jwt签名是否过期
			//jwt.ValidationErrorNotValidYet是验证用户操作是否活跃
			c.JSON(403, "您不活跃或者验证已过期！")
			c.Abort()
			return
		} else {
			c.JSON(403, gin.H{
				"error:": err1,
			})
		}
	} else {
		c.JSON(403, gin.H{
			"不能识别此token": err1,
		})
	}
}
