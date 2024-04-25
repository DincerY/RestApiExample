package models

import (
	u "RestApiExample/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strings"
	"time"
)

type Token struct {
	UserId   uint
	Username string
	jwt.StandardClaims
}

type Account struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token";sql:"-"`
}

func (account *Account) Validate() (map[string]interface{}, bool) {

	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "Email address is not valid"), false
	}

	if len(account.Password) < 8 {
		return u.Message(false, "Password must be at least 8 characters"), false
	}

	temp := &Account{}

	err := GetDB().Table("accounts").Where("email = ?", account.Email).First(temp).Error

	if err != nil && err == gorm.ErrRecordNotFound {
		return u.Message(false, "Bağlantı hatası oluştu. Lütfen tekrar deneyiniz!"), false
	}

	if temp.Email != "" {
		return u.Message(false, "Email address already exist"), false
	}

	return u.Message(false, "Requirement passed"), true
}

func (account *Account) Create() map[string]interface{} {

	if resp, ok := account.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	GetDB().Create(account)

	if account.ID <= 0 {
		return u.Message(false, "Bağlantı hatası oluştu. Kullanıcı yaratılamadı!")
	}

	tk := &Token{UserId: account.ID, Username: strings.Split(account.Email, "@")[0]}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	account.Password = " "

	response := u.Message(true, "Account created")

	response["account"] = account
	return response
}

func Login(w http.ResponseWriter, email, password string) map[string]interface{} {
	account := &Account{}

	err := GetDB().Table("accounts").Where("email = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email not found")
		}
		return u.Message(false, "There are connection error")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return u.Message(false, "Email or password is not match")
	}

	account.Password = ""

	tk := &Token{UserId: account.ID, Username: strings.Split(account.Email, "@")[0]}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	resp := u.Message(true, "Logging successfully")
	resp["account"] = account

	c := http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		Domain:   "localhost",
		Expires:  time.Time{},
		HttpOnly: false,
		SameSite: 0,
	}
	http.SetCookie(w, &c)

	return resp
}

func GetUser(id string) map[string]interface{} {
	acc := &Account{}
	GetDB().Table("accounts").Where("id = ?", id).First(acc)
	if acc.Email == "" {
		return u.Message(false, "Girilen id de bir hesap bulunamadı")
	}
	acc.Password = ""
	resp := u.Message(true, "Account found")
	resp["account"] = acc
	return resp
}
