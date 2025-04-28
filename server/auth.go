package server

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"teamdev/internal/models"
)

type fullFormData struct {
	Email          string `form:"InputEmail"`
	Password       string `form:"InputPassword"`
	PasswordRepeat string `form:"InputPassword2"`
	Name           string `form:"InputName"`
	Surname        string `form:"InputSurname"`
	Phone          string `form:"InputPhone"`
	Address        string `form:"InputAddress"`
}

func (s *Services) authenticatedUser(c *gin.Context) *models.User {
	session := sessions.Default(c)
	sessionID := session.Get("userID")

	if sessionID != nil {
		strUserID, ok := sessionID.(string)
		if ok {
			userId, err := uuid.Parse(strUserID)
			if err == nil {
				user, err := s.Services.UserService.GetUserByID(userId)
				if err == nil {
					return user
				}
			}
		}
	}
	return nil
}

func (s *Services) signupGet(c *gin.Context) {
	c.HTML(200, "signup", gin.H{
		"title": "Регистрация",
	})
}

func (s *Services) signupPost(c *gin.Context) {
	var data fullFormData
	if err := c.Bind(&data); err != nil {
		c.HTML(http.StatusBadRequest, "signup", gin.H{
			"title":    "Регистрация",
			"error":    err.Error(),
			"formData": data,
		})
		return
	}

	if data.Password != data.PasswordRepeat {
		c.HTML(http.StatusBadRequest, "signup", gin.H{
			"title":    "Регистрация",
			"error":    "Пароли не совпадают",
			"formData": data,
		})
		return
	}

	// Check if the user exists already
	_, err := s.Services.UserService.GetUserByEmail(data.Email)
	if err == nil {
		c.HTML(http.StatusBadRequest, "signup", gin.H{
			"title":    "Регистрация",
			"error":    "Пользователь с таким email уже существует",
			"formData": data,
		})
		return
	}

	// Create the user
	user, err := s.Services.UserService.Register(&models.User{
		Email:       data.Email,
		Name:        data.Name,
		Surname:     data.Surname,
		PhoneNumber: data.Phone,
		Address:     data.Address,
	}, data.Password)

	if err != nil {
		c.HTML(http.StatusBadRequest, "signup", gin.H{
			"title":    "Регистрация",
			"error":    err.Error(),
			"formData": data,
		})
		return
	}

	// Set the session.
	session := sessions.Default(c)
	session.Set("userID", user.ID.String())
	session.Save()

	c.Redirect(http.StatusFound, "/")
}

type loginFormData struct {
	Email    string `form:"InputEmail"`
	Password string `form:"InputPassword"`
}

func (s *Services) signinGet(c *gin.Context) {
	c.HTML(200, "signin", gin.H{
		"title": "Вход",
	})
}

func (s *Services) signinPost(c *gin.Context) {
	var data loginFormData
	if err := c.Bind(&data); err != nil {
		c.HTML(http.StatusBadRequest, "signin", gin.H{
			"title": "Вход",
			"error": err.Error(),
		})
		return
	}

	// try to login
	user, err := s.Services.UserService.Login(data.Email, data.Password)
	if err != nil {
		c.HTML(http.StatusBadRequest, "signin", gin.H{
			"title":    "Вход",
			"error":    "Неверный пароль или пользователь с таким email не существует",
			"formData": data,
		})
		return
	}

	// Set the session.
	session := sessions.Default(c)
	session.Set("userID", user.ID.String())
	ok := session.Save()
	if ok != nil {
		c.HTML(http.StatusBadRequest, "signin", gin.H{
			"title":    "Вход",
			"error":    "Не удалось сохранить сессию",
			"formData": data,
		})
		return
	}

	c.Redirect(http.StatusFound, "/")

}

func (s *Services) logout(c *gin.Context) {
	// Delete the session
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Status(http.StatusAccepted)
	c.Redirect(http.StatusFound, "/")
}
