package Controllers

import (
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	domain "github.com/segnig/task-manager/Domains"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	UserUsecase domain.UserUsecase
	Password    domain.PasswordServiceProvider
	UserToken   domain.IUserToken
}

func (uc *UserController) Register(c *gin.Context) {
	var user domain.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	userID := primitive.NewObjectID()
	user.ID = userID
	user.UserID = userID.Hex()

	hashedPassword := uc.Password.HashPassword(user.Password)
	user.Password = hashedPassword

	var validUsername = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*$`)
	if !validUsername.MatchString(user.Username) {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "username should start with alphabet and only contain alpha numeric and underscore"})
		return
	}

	token, refreshToken, err := uc.UserToken.GenerateAllTokens(user.Username, user.UserType, user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	user.Token = token
	user.RefreshToken = refreshToken

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if err := uc.UserUsecase.Create(c, &user); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, domain.ErrorResponse{Message: "user created successfully"})
}

func (uc *UserController) Login(c *gin.Context) {
	var user domain.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	foundUser, err := uc.UserUsecase.GetUserByUsername(c, user.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	passwordIsValid, msg := uc.Password.VerifyPassword(foundUser.Password, user.Password)
	if !passwordIsValid {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: msg})
		return
	}
	token, refreshToken, err := uc.UserToken.GenerateAllTokens(foundUser.Username, foundUser.UserType, foundUser.UserID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
	}
	uc.UserUsecase.UpdateAllToken(c, token, refreshToken, foundUser.UserID)
	foundUser.Token = token
	foundUser.RefreshToken = refreshToken
	c.JSON(http.StatusOK, gin.H{
		"user_id":    foundUser.UserID,
		"username":   foundUser.Username,
		"first_name": foundUser.FirstName,
		"last_name":  foundUser.LastName,
		"token":      foundUser.Token,
		"user_type":  foundUser.UserType,
	})

}

func (uc *UserController) FetchAll(c *gin.Context) {
	users, err := uc.UserUsecase.FetchAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (uc *UserController) Fetch(c *gin.Context) {
	userID := c.Param("user_id")
	user, err := uc.UserUsecase.FetchById(c, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (uc *UserController) Update(c *gin.Context) {
	userID := c.GetString("user_id")
	var user domain.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	if err := uc.UserUsecase.UpdateById(c, userID, &user); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "User Update successfully"})
}

func (uc *UserController) Delete(c *gin.Context) {
	userID := c.GetString("user_id")

	if err := uc.UserUsecase.DeleteById(c, userID); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "user deleted successfully"})
}
