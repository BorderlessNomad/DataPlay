package main

import (
	bcrypt "code.google.com/p/go.crypto/bcrypt"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
)

type User struct {
	Uid      int `primaryKey:"yes"`
	Email    string
	Password string
}

func (u User) TableName() string {
	return "priv_users"
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Shortcut function that is used to redirect users to the login page if they are not logged in.
func checkAuth(res http.ResponseWriter, req *http.Request) {
	if !(IsUserLoggedIn(res, req)) {
		http.Redirect(res, req, "/login", http.StatusTemporaryRedirect)
		return
	}
}

func HandleLogin(res http.ResponseWriter, req *http.Request) {
	username := req.FormValue("username")
	password := req.FormValue("password")

	user := User{}
	err := DB.Select("uid, email, password").Where("email = ?", username).Find(&user).Error
	check(err)

	if user.Password != "" && bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil { // Check the password with bcrypt
		err := DB.Select("uid").Where("email = ?", username).Find(&user).Error
		check(err)

		if SetSession(res, req, user.Uid) != nil {
			http.Error(res, "Could not setup session.", http.StatusInternalServerError)
			return
		}

		http.Redirect(res, req, "/", http.StatusFound)
	} else {
		// Just in the case that the user is on a really old MD5 password (useful for admins resetting passwords too) check
		count := 0
		err := DB.Model(&user).Where("password = ?", GetMD5Hash(password)).Count(&count).Error
		check(err)

		if count != 0 {
			// Ooooh, We need to upgrade this password!
			hashedPassword, e := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if e == nil {
				err := DB.Model(&user).Update("password", string(hashedPassword)).Error
				check(err)

				if SetSession(res, req, user.Uid) != nil {
					http.Error(res, "Could not setup session.", http.StatusInternalServerError)
					return
				}

				http.Redirect(res, req, "/", http.StatusFound)
			}

			http.Redirect(res, req, fmt.Sprintf("/login?failed=3&r=%s", e), http.StatusFound)
		} else {
			http.Redirect(res, req, "/login?failed=1", http.StatusNotFound) // The user has failed this test as well :sad tuba:
		}
	}
}

func HandleLogout(res http.ResponseWriter, req *http.Request) {
	ClearSession(res, req)

	http.Redirect(res, req, "/login", http.StatusTemporaryRedirect)
	return
}

func HandleRegister(res http.ResponseWriter, req *http.Request) string {
	username := req.FormValue("username")
	password := req.FormValue("password")

	user := User{}
	err := DB.Where("email = ?", username).First(&user).Error
	if err != gorm.RecordNotFound {
		http.Error(res, "That username is already registered.", http.StatusConflict)
		return "That username is already registered."
	}

	hashedPassword, err1 := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err1 != nil {
		return "The password you entered is invalid."
	}

	user.Email = username
	user.Password = string(hashedPassword)
	err2 := DB.Save(&user).Error
	if err2 != nil {
		return "Could not make the user you requested."
	}

	SetSession(res, req, user.Uid)

	http.Redirect(res, req, "/", http.StatusFound)
	return ""
}
