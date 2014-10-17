package main

import (
	bcrypt "code.google.com/p/go.crypto/bcrypt"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"github.com/codegangsta/martini"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

// REPUTATION POINTS
const obsCredit int = 5      // observation is voted up
const discCredit int = 15    // discovery is valdiated
const discObs int = 2        // discovery receives an observation
const rankUp int = 10        // reach new rank
const topRank int = 100      // reach top 10 Experts rank
const discHot int = 50       // discovery is hot
const obsDiscredit int = -1  // observation is voted down
const discDiscredit int = -2 // discovery is voted down
const obsSpam int = -100     // observation deleted after being flagged

/// USER TYPES
const UserTypeNormal int = 0
const UserTypeAdmin int = 1

type UserForm struct {
	Username string `json:"username"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"`
}

type UserNameForm struct {
	Username string `json:"username" binding:"required"`
}

type UserSocialForm struct {
	Network   string `json:"network" binding:"required"`
	Id        string `json:"id" binding:"required"`
	Email     string `json:"email" binding:"required"`
	FullName  string `json:"full_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Image     string `json:"image"`
}

type UserDetailsForm struct {
	Username string `json:"username"`
	Email    string `json:"email"`
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
	if login.Username == "" && login.Email == "" || login.Password == "" {
		http.Error(res, "Username/Password missing.", http.StatusBadRequest)
		return ""
	}

	user := User{}
	var err error

	if login.Username != "" {
		err = DB.Where("username = ?", login.Username).Find(&user).Error
		if err == gorm.RecordNotFound {
			http.Error(res, "No such user found!", http.StatusNotFound)
			return ""
		} else if err != nil {
			http.Error(res, "No such user found!", http.StatusInternalServerError)
			return ""
		}
	} else {
		err = DB.Where("email = ?", login.Email).Find(&user).Error
		if err == gorm.RecordNotFound {
			http.Error(res, "No such user found!", http.StatusNotFound)
			return ""
		} else if err != nil {
			http.Error(res, "No such user found!", http.StatusInternalServerError)
			return ""
		}
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
			http.Error(res, "Unable to find user with MD5.", http.StatusInternalServerError)
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
		"user":     user.Username,
		"session":  session.Value,
		"usertype": user.Usertype,
	}
	usr, _ := json.Marshal(u)

	return string(usr)
}

func HandleSocialLogin(res http.ResponseWriter, req *http.Request, login UserSocialForm) string {
	if login.Network == "" || login.Id == "" || login.Email == "" {
		http.Error(res, "Network/Id/Email missing.", http.StatusBadRequest)
		return ""
	}

	user := User{}
	var err error
	err = DB.Where("email = ?", login.Email).Find(&user).Error

	if err == nil { // if email is found then automatically log user in, no password required for social login
		session, e := SetSession(user.Uid)
		if e != nil {
			http.Error(res, e.Message, e.Code)
			return ""
		}

		u := map[string]interface{}{
			"user":     user.Email,
			"session":  session.Value,
			"usertype": user.Usertype,
		}
		usr, _ := json.Marshal(u)

		return string(usr)

	} else if err != nil && err == gorm.RecordNotFound { // if user does not exist then create social user
		user.Email = login.Email
		user.Password = GetMD5Hash(user.Email + time.Now().String()) // not important here, just used hashed version of their email plus the time
		user.Avatar = login.Image
		user.Username = login.FirstName + login.LastName
		user.Usertype = 1

		err = DB.Save(&user).Error
		if err != nil {
			http.Error(res, err.Error()+" - database query failed - Unable to save user's standard details.", http.StatusInternalServerError)
			return ""
		}

		social := Social{}
		social.FirstName = login.FirstName
		social.Network = login.Network
		social.LastName = login.LastName
		social.FullName = login.FullName
		social.NetworkUserId = login.Id
		newUser := User{}
		err = DB.Where("email = ?", login.Email).Find(&newUser).Error // find newly created user to get generated uid

		if err != nil {
			http.Error(res, err.Error()+" - could not recall user after creation", http.StatusInternalServerError)
			return ""
		}

		social.Uid = newUser.Uid

		err = DB.Save(&social).Error
		if err != nil {
			http.Error(res, err.Error()+" - database query failed - Unable to save user's social details.", http.StatusInternalServerError)
			return ""
		}

		session, e := SetSession(social.Uid)
		if e != nil {
			http.Error(res, e.Message, e.Code)
			return ""
		}

		u := map[string]interface{}{
			"user":     newUser.Email,
			"session":  session.Value,
			"usertype": newUser.Usertype,
		}
		usr, _ := json.Marshal(u)

		return string(usr)
	}

	http.Error(res, err.Error()+" - problem finding registered user", http.StatusInternalServerError)
	return ""
}

func HandleLogout(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	sid := req.Header.Get("X-API-SESSION")
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
	if register.Username == "" || register.Password == "" || register.Email == "" {
		http.Error(res, "Username, Password or Email missing.", http.StatusBadRequest)
		return ""
	}

	user := User{}
	err := DB.Where("username = ?", register.Username).First(&user).Error
	if err != gorm.RecordNotFound {
		http.Error(res, "Username already exists.", http.StatusConflict)
		return ""
	}

	err = DB.Where("email = ?", register.Email).First(&user).Error
	if err != gorm.RecordNotFound {
		http.Error(res, "Email already exists.", http.StatusConflict)
		return ""
	}

	hashedPassword, err1 := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err1 != nil {
		http.Error(res, "Unable to generate password hash.", http.StatusInternalServerError)
		return ""
	}

	user.Email = register.Email
	user.Username = register.Username
	user.Password = string(hashedPassword)
	err2 := DB.Save(&user).Error
	if err2 != nil {
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
		"user":     user.Username,
		"session":  session.Value,
		"usertype": user.Usertype,
	}
	usr, _ := json.Marshal(u)

	return string(usr)
}

func HandleCheckUsername(res http.ResponseWriter, req *http.Request, user UserNameForm) string {
	if user.Username == "" {
		http.Error(res, "Invalid or empty username.", http.StatusBadRequest)
		return ""
	}

	validUser := User{}
	err := DB.Where("username = ?", user.Username).First(&validUser).Error
	if err != nil && err != gorm.RecordNotFound {
		http.Error(res, "Unable to find that User.", http.StatusInternalServerError)
		return ""
	} else if err == gorm.RecordNotFound {
		http.Error(res, "We couldn't find an account associated with "+user.Username, http.StatusNotFound)
		return ""
	}

	u := map[string]interface{}{
		"user":   validUser.Username,
		"exists": true,
	}
	usr, _ := json.Marshal(u)

	return string(usr)
}

func HandleForgotPassword(res http.ResponseWriter, req *http.Request, user UserNameForm) string {
	if user.Username == "" {
		http.Error(res, "Invalid or empty username.", http.StatusBadRequest)
		return ""
	}

	validUser := User{}
	err := DB.Where("username = ?", user.Username).First(&validUser).Error
	if err != nil && err != gorm.RecordNotFound {
		http.Error(res, "Database query failed (User)", http.StatusInternalServerError)
		return ""
	} else if err == gorm.RecordNotFound {
		http.Error(res, "We couldn't find an account associated with "+user.Username, http.StatusNotFound)
		return ""
	}

	randomString := make([]byte, 64)
	_, e := rand.Read(randomString)
	if e != nil {
		http.Error(res, "Unable to generate hash.", http.StatusInternalServerError)
		return ""
	}
	hashString := sha1.New()
	hashString.Write(randomString)
	hash := hex.EncodeToString(hashString.Sum(nil))

	token := UserTokens{
		Uid:     validUser.Uid,
		Hash:    hash,
		Created: time.Now(),
	}

	dbError := DB.Save(&token).Error
	if dbError != nil {
		http.Error(res, "Database query failed (Token)", http.StatusInternalServerError)
		return ""
	}

	u := map[string]interface{}{
		"user":  validUser.Username,
		"token": hash,
	}
	usr, _ := json.Marshal(u)

	return string(usr)
}

func ResetPassword(hash, username, password string) *appError {
	if username == "" || hash == "" {
		return &appError{nil, "No username/token found!", http.StatusBadRequest}
	}

	user := User{}
	err := DB.Where("username = ?", username).First(&user).Error
	if err != nil && err != gorm.RecordNotFound {
		return &appError{nil, "Database query failed (User).", http.StatusInternalServerError}
	} else if err == gorm.RecordNotFound {
		return &appError{nil, "Invalid email!", http.StatusNotFound}
	}

	token := UserTokens{}
	dbError := DB.Where("uid = ?", user.Uid).Where("hash = ?", hash).Where("used = ?", false).Last(&token).Error
	if dbError != nil && dbError != gorm.RecordNotFound {
		return &appError{nil, "Database query failed (Token).", http.StatusInternalServerError}
	} else if dbError == gorm.RecordNotFound {
		return &appError{nil, "Invalid token!", http.StatusNotFound}
	}

	if password == "" {
		return nil
	}

	hashedPassword, errH := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if errH != nil {
		return &appError{errH, "Unable to generate password hash.", http.StatusInternalServerError}
	}

	user.Password = string(hashedPassword)
	err = DB.Save(&user).Error
	if err != nil {
		return &appError{err, "Database query failed (Password).", http.StatusInternalServerError}
	}

	token.Used = true
	err = DB.Save(&token).Error
	if err != nil {
		return &appError{err, "Database query failed (Token).", http.StatusInternalServerError}
	}

	return nil
}

func HandleResetPasswordCheck(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	err := ResetPassword(params["token"], params["username"], "")
	if err != nil {
		http.Error(res, err.Message, err.Code)
		return ""
	}

	return "OK"
}

func HandleResetPassword(res http.ResponseWriter, req *http.Request, params martini.Params, user UserForm) string {
	err := ResetPassword(params["token"], user.Username, user.Password)
	if err != nil {
		http.Error(res, err.Message, err.Code)
		return ""
	}

	return "OK"
}

func GetUserDetails(res http.ResponseWriter, req *http.Request) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter.", http.StatusBadRequest)
		return ""
	}

	uid, err := GetUserID(session)
	if err != nil {
		http.Error(res, err.Message, err.Code)
		return ""
	}

	user := User{}
	err1 := DB.Where("uid = ?", uid).First(&user).Error
	if err1 != nil && err1 != gorm.RecordNotFound {
		http.Error(res, "Database query failed (User).", http.StatusInternalServerError)
		return ""
	} else if err1 == gorm.RecordNotFound {
		http.Error(res, "No such user found!", http.StatusNotFound)
		return ""
	}

	u := map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
		"usertype": user.Usertype,
	}
	usr, _ := json.Marshal(u)

	return string(usr)
}

func UpdateUserDetails(res http.ResponseWriter, req *http.Request, data UserDetailsForm) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter.", http.StatusBadRequest)
		return ""
	}

	if data.Email == "" && data.Username == "" {
		http.Error(res, "Missing request parameter.", http.StatusBadRequest)
		return ""
	}

	uid, err := GetUserID(session)
	if err != nil {
		http.Error(res, err.Message, err.Code)
		return ""
	}

	user := User{}
	err1 := DB.Where("uid = ?", uid).First(&user).Error
	if err1 != nil && err1 != gorm.RecordNotFound {
		http.Error(res, "Database query failed (Select).", http.StatusInternalServerError)
		return ""
	} else if err1 == gorm.RecordNotFound {
		http.Error(res, "No such user found!", http.StatusNotFound)
		return ""
	}

	if data.Username != "" {
		user.Username = data.Username
	}

	if data.Email != "" {
		user.Email = data.Email
	}

	err2 := DB.Save(&user).Error
	if err2 != nil {
		http.Error(res, "Database query failed (Update).", http.StatusInternalServerError)
		return ""
	}

	u := map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
	}
	usr, _ := json.Marshal(u)

	return string(usr)
}

func Reputation(uid int, points int) string {
	usr := User{}

	err := DB.Where("uid = ?", uid).Find(&usr).Error
	if err != nil {
		return "failed to find user"
	}

	r := usr.Reputation + points

	err = DB.Model(usr).Where("uid = ?", uid).Update("reputation", r).Error
	if err != nil {
		return "failed to update reputation"
	}

	return ""
}
