package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/knavdeep152002/stream-sense/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (a *Auth) createUserInDB(authInput *authInput) (user models.User, err error) {
	var userFound models.User
	a.DB.Where("username=?", authInput.Username).Find(&userFound)

	if userFound.ID != 0 {
		err = errors.New("username already used")
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(authInput.Password), bcrypt.DefaultCost)
	if err != nil {
		err = errors.New("error hashing password")
		return
	}

	user = models.User{
		Username: authInput.Username,
		Password: string(passwordHash),
	}

	a.DB.Create(&user)
	return
}

func (a *Auth) getSignedToken(authInput *authInput) (token string, err error) {
	var userFound models.User
	a.DB.Where("username=?", authInput.Username).Find(&userFound)
	if userFound.ID == 0 {
		err = errors.New("user not found")
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(authInput.Password)); err != nil {
		err = errors.New("invalid password")
		return
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   userFound.ID,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
		"iss":  "cohort",
		"name": userFound.Username,
	})

	token, err = generateToken.SignedString([]byte(os.Getenv("SECRET")))
	return
}
