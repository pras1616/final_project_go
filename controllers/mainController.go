package controllers

import (
	"errors"
	"final_project/helpers"
	"final_project/models"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Controllers interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type ControllersStruct struct {
	DB_Comment     *gorm.DB
	DB_User        *gorm.DB
	DB_Photo       *gorm.DB
	DB_SocialMedia *gorm.DB
}

func NewCarsController(db1 *gorm.DB, db2 *gorm.DB, db3 *gorm.DB, db4 *gorm.DB) Controllers {
	return &ControllersStruct{
		DB_Comment:     db1,
		DB_User:        db2,
		DB_Photo:       db3,
		DB_SocialMedia: db4,
	}
}

func validation(val interface{}) string {
	validate := validator.New()
	errVal := validate.Struct(val)
	if errVal != nil {
		if _, ok := errVal.(*validator.InvalidValidationError); ok {
			fmt.Println(errVal)
			return "true"
		}

		msgErr := ""
		for _, err := range errVal.(validator.ValidationErrors) {
			msgErr = err.StructField() + " " + err.ActualTag() + " " + err.Param()
			if err.ActualTag() == "email" {
				msgErr = "email not valid"
			}

		}
		return msgErr
	}
	return "true"
}

func (g *ControllersStruct) Register(c *gin.Context) {
	var payload models.User
	if err := c.ShouldBindJSON(&payload); err != nil {
		fmt.Println("error found: ", err)
		c.JSON(400, gin.H{
			"code":    400,
			"data":    "-",
			"message": strings.Split(err.Error(), "User.")[1],
		})
		return
	}

	errVal := validation(&payload)
	if errVal != "true" {
		c.JSON(409, gin.H{
			"code":    409,
			"data":    "-",
			"message": errVal,
		})
		return
	}

	errFind := g.DB_User.Where("email = ? AND username = ?", &payload.Email, &payload.Username).First(&payload).Error
	if errFind != nil { // if not found
		// func CompareHashAndPassword(hashedPassword, password []byte) error
		hash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}

		register := models.User{
			Username:   payload.Username,
			Email:      payload.Email,
			Password:   string(hash),
			Age:        payload.Age,
			Created_at: time.Now(),
			Updated_at: time.Now(),
		}

		errUser := g.DB_User.Create(&register).Error
		if err != nil {
			fmt.Println("error found: ", errUser)
			c.JSON(500, gin.H{
				"message": "internal server error",
			})
			return
		}

		type Response struct {
			Age      int    `form:"age" json:"age"`
			Email    string `form:"email" json:"email"`
			ID       uint   `form:"id" json:"id"`
			Username string `form:"username" json:"username"`
		}

		response := Response{
			Age:      register.Age,
			Email:    register.Email,
			ID:       register.ID,
			Username: register.Username,
		}

		c.JSON(201, gin.H{
			"code":    201,
			"data":    response,
			"message": "succesful to resgister new user",
		})
	} else {
		c.JSON(400, gin.H{
			"code":    400,
			"data":    "-",
			"message": "username or email already exist",
		})
		return
	}
}

func (g *ControllersStruct) Login(c *gin.Context) {
	type Login struct {
		Email    string `form:"email" json:"email" xml:"email" binding:"required" validate:"required,email"`
		Password string `form:"password" json:"password" xml:"password" binding:"required" validate:"required,min=8"`
	}
	var payload Login
	if err := c.ShouldBindJSON(&payload); err != nil {
		fmt.Println("error found: ", err)
		c.JSON(400, gin.H{
			"code":    400,
			"data":    "-",
			"message": "email or password not valid",
		})
		return
	}

	var user models.User
	errFind := g.DB_User.Where("email = ?", &payload.Email).First(&user).Error
	if errFind != nil { // if not found
		c.JSON(400, gin.H{
			"code":    400,
			"data":    "-",
			"message": "email or password not valid",
		})
	} else {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
			c.JSON(400, gin.H{
				"code":    400,
				"data":    "-",
				"message": "email or password not valid value",
			})
			return
		}
		var token = helpers.GenerateToken(user.ID, user.Email)
		c.JSON(200, gin.H{
			"code":    200,
			"token":   token,
			"message": "success login",
		})
		return
	}
}

func (g *ControllersStruct) UpdateUser(c *gin.Context) {
	type UserBody struct {
		Email    string `form:"email" json:"email" xml:"email" binding:"required" validate:"required,email"`
		Username string `form:"username" json:"username" xml:"username" binding:"required" validate:"required,min=4,max=50"`
	}

	var payload UserBody
	if err := c.ShouldBindJSON(&payload); err != nil {
		fmt.Println("error found: ", err)
		c.JSON(400, gin.H{
			"code":    400,
			"data":    "-",
			"message": "err binding",
		})
		return
	}
	idStr, _ := c.Params.Get("userId")
	id, err := strconv.ParseFloat(idStr, 64)
	if err != nil {
		fmt.Println("error found: ", err)
		c.JSON(400, gin.H{
			"code":    409,
			"data":    "-",
			"message": err,
		})
		return
	}
	errVal := validation(&payload)
	if errVal != "true" {
		c.JSON(409, gin.H{
			"code":    409,
			"data":    "-",
			"message": errVal,
		})
		return
	}

	userData := c.MustGet("userData").(jwt.MapClaims)
	userId := float64(userData["id"].(float64))
	if userId != id {
		c.JSON(400, gin.H{
			"code":    400,
			"data":    "-",
			"message": "not allowed, id not yours",
		})
		return
	}
	req := models.User{
		ID:       uint(id),
		Email:    payload.Email,
		Username: payload.Username,
		// Age: ,
		Updated_at: time.Now(),
	}
	errUpdate := g.DB_User.Model(&req).Where("id = ?", id).Updates(models.User{
		Email:      req.Email,
		Username:   req.Username,
		Updated_at: req.Updated_at,
	}).Error
	if errors.Is(errUpdate, gorm.ErrRecordNotFound) {
		fmt.Println("error found: ", errUpdate)
		c.JSON(404, gin.H{
			"code":    404,
			"data":    "-",
			"message": "data not found",
		})
		return
	}
	if errUpdate != nil {
		fmt.Println("error found: ", errUpdate)
		c.JSON(500, gin.H{
			"code":    500,
			"data":    "-",
			"message": "internal server error",
		})
		return
	}
	type Res struct {
		ID         uint      `form:"id" json:"id"`
		Username   string    `form:"username" json:"username"`
		Email      string    `form:"email" json:"email"`
		Age        int       `form:"age" json:"age"`
		Updated_at time.Time `form:"updated_at" json:"updated_at"`
	}

	var mUser models.User
	errorFInd := g.DB_User.Where("id = ?", req.ID).Find(&mUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("error found: ", errorFInd)
		c.JSON(404, gin.H{
			"message": "data empty",
		})
		return
	}
	if errorFInd != nil {
		fmt.Println("error found: ", errorFInd)
		c.JSON(500, gin.H{
			"message": "internal server error",
		})
		return
	}
	var res = Res{
		ID:         mUser.ID,
		Username:   mUser.Username,
		Email:      mUser.Email,
		Age:        mUser.Age,
		Updated_at: mUser.Updated_at,
	}
	c.JSON(200, gin.H{
		"code":    200,
		"data":    res,
		"message": "success update data user",
	})
}

func (g *ControllersStruct) DeleteUser(c *gin.Context) {
	idStr, _ := c.Params.Get("userId")
	id, err := strconv.ParseFloat(idStr, 64)
	if err != nil {
		fmt.Println("error found: ", err)
		c.JSON(400, gin.H{
			"code":    409,
			"data":    "-",
			"message": err,
		})
		return
	}

	userData := c.MustGet("userData").(jwt.MapClaims)
	userId := float64(userData["id"].(float64))
	if userId != id {
		c.JSON(400, gin.H{
			"code":    400,
			"data":    "-",
			"message": "not allowed, id not yours",
		})
		return
	}

	req := models.User{}

	errDel := g.DB_User.Delete(&req, "id = ?", uint(id)).Error
	if errors.Is(errDel, gorm.ErrRecordNotFound) {
		fmt.Println("error found: ", errDel)
		c.JSON(404, gin.H{
			"code":    404,
			"data":    "-",
			"message": "data not found",
		})
		return
	}
	if errDel != nil {
		fmt.Println("error found: ", errDel)
		c.JSON(500, gin.H{
			"code":    500,
			"data":    "-",
			"message": "internal server error",
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"data":    "-",
		"message": "your account successfully deleted",
	})
}
