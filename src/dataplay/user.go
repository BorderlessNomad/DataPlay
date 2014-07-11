package main

import (
	bcrypt "code.google.com/p/go.crypto/bcrypt"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/codegangsta/martini"
	"github.com/jinzhu/gorm"
	"net/http"
)

type UserForm struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Shortcut function that is used to redirect users to the login page if they are not logged in.
func CheckAuthRedirect(res http.ResponseWriter, req *http.Request) {
	if !(IsUserLoggedIn(res, req)) {
		http.Redirect(res, req, "/login", http.StatusTemporaryRedirect)
		return
	}
}

func HandleLogin(res http.ResponseWriter, req *http.Request, login UserForm) string {
	if login.Username == "" || login.Password == "" {
		http.Error(res, "Username/Password missing.", http.StatusBadRequest)
		return ""
	}

	user := User{}
	var err error
	err = DB.Where("email = ?", login.Username).Find(&user).Error
	if err == gorm.RecordNotFound {
		http.Error(res, "Could not find a user.", http.StatusNotFound)
		return ""
	} else if err != nil {
		http.Error(res, "Could not find a user.", http.StatusInternalServerError)
		return ""
	}

	// Check the password with bcrypt
	if len(user.Password) > 0 && bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)) == nil {
		// Do nothing and continue :)
	} else {
		// Just in the case that the user is on a really old MD5 password (useful for admins resetting passwords too) check
		count := 0
		err := DB.Model(&user).Where("password = ?", GetMD5Hash(login.Password)).Count(&count).Error

		if err != nil && err != gorm.RecordNotFound {
			check(err)
			http.Error(res, "Could not find a user.", http.StatusInternalServerError)
			return ""
		}

		if count == 0 {
			http.Error(res, "Invalid username/password.", http.StatusBadRequest)
			return ""
		}

		// Ooooh, We need to upgrade this password!
		hashedPassword, e := bcrypt.GenerateFromPassword([]byte(login.Password), bcrypt.DefaultCost)
		if e != nil {
			http.Error(res, "Unable to upgrade the password.", http.StatusInternalServerError)
			return ""
		}

		err = DB.Model(&user).Update("password", string(hashedPassword)).Error
		if err != nil {
			check(err)
			http.Error(res, "Unable to update the password.", http.StatusInternalServerError)
			return ""
		}
	}

	session, e := SetSession(user.Uid)
	if e != nil {
		http.Error(res, e.Message, e.Code)
		return ""
	}

	u := map[string]interface{}{
		"username": user.Email,
		"session": map[string]interface{}{
			"name":   session.Name,
			"value":  session.Value,
			"expiry": 60 * 60 * 24 * 365, // Year (seconds)
		},
	}
	usr, _ := json.Marshal(u)

	return string(usr)
}

func HandleLogout(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	sid := params["session"]
	if len(sid) <= 0 {
		http.Error(res, "Missing session parameter.", http.StatusBadRequest)
		return ""
	}

	_, e := ClearSession(sid)
	if e != nil {
		http.Error(res, e.Message, e.Code)
		return ""
	}

	return ""
}

func HandleRegister(res http.ResponseWriter, req *http.Request, register UserForm) string {
	if register.Username == "" || register.Password == "" {
		http.Error(res, "Username/Password missing.", http.StatusBadRequest)
		return ""
	}

	user := User{}
	err := DB.Where("email = ?", register.Username).First(&user).Error
	if err != gorm.RecordNotFound {
		http.Error(res, "Username already exists.", http.StatusConflict)
		return ""
	}

	hashedPassword, err1 := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err1 != nil {
		http.Error(res, "Invalid Username/Password.", http.StatusBadRequest)
		return ""
	}

	user.Email = register.Username
	user.Password = string(hashedPassword)
	err2 := DB.Save(&user).Error
	if err2 != nil {
		check(err2)
		http.Error(res, "Unable to create user.", http.StatusInternalServerError)
		return ""
	}

	var session *http.Cookie
	var e *appError
	session, e = SetSession(user.Uid)
	if e != nil {
		http.Error(res, e.Message, e.Code)
		return ""
	}

	u := map[string]interface{}{
		"username": user.Email,
		"session": map[string]interface{}{
			"name":   session.Name,
			"value":  session.Value,
			"expiry": 60 * 60 * 24 * 365, // Year (seconds)
		},
	}
	usr, _ := json.Marshal(u)

	return string(usr)
}
