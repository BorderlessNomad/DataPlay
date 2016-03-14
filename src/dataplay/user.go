package main

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"github.com/ahirmayur/gorm"
	"github.com/codegangsta/martini"
	bcrypt "golang.org/x/crypto/bcrypt"
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
	Email string `json:"email" binding:"required"`
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

	if login.Username != "" {
		gErr := DB.Where("username = ?", login.Username).Find(&user).Error
		if gErr == gorm.RecordNotFound {
			http.Error(res, "No such user found!", http.StatusNotFound)
			return ""
		} else if gErr != nil {
			http.Error(res, "No such user found!", http.StatusInternalServerError)
			return ""
		}
	} else {
		gErr := DB.Where("email = ?", login.Email).Find(&user).Error
		if gErr == gorm.RecordNotFound {
			http.Error(res, "No such user found!", http.StatusNotFound)
			return ""
		} else if gErr != nil {
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
		gErr := DB.Model(&user).Where("password = ?", GetMD5Hash(login.Password)).Count(&count).Error

		if gErr != nil && gErr != gorm.RecordNotFound {
			http.Error(res, "Unable to find user with MD5.", http.StatusInternalServerError)
			return ""
		}

		if count == 0 {
			http.Error(res, "Invalid username/password.", http.StatusBadRequest)
			return ""
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(login.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, "Unable to upgrade the password.", http.StatusInternalServerError)
			return ""
		}

		gErr = DB.Model(&user).Update("password", string(hashedPassword)).Error
		if gErr != nil {
			http.Error(res, "Unable to update the password.", http.StatusInternalServerError)
			return ""
		}
	}

	session, err := SetSession(user.Uid)
	if err != nil {
		http.Error(res, err.Message, err.Code)
		return ""
	}

	u := map[string]interface{}{
		"uid":      user.Uid,
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
	gErr := DB.Where("email = ?", login.Email).Find(&user).Error

	if gErr == nil { // if email is found then automatically log user in, no password required for social login
		session, err := SetSession(user.Uid)
		if err != nil {
			http.Error(res, err.Message, err.Code)
			return ""
		}

		u := map[string]interface{}{
			"uid":        user.Uid,
			"user":       user.Email,
			"email_hash": GetMD5Hash(user.Email),
			"session":    session.Value,
			"usertype":   user.Usertype,
		}

		usr, _ := json.Marshal(u)

		return string(usr)

	} else if gErr != nil && gErr == gorm.RecordNotFound { // if user does not exist then create social user
		user.Email = login.Email
		user.Password = GetMD5Hash(user.Email + time.Now().String()) // not important here, just used hashed version of their email plus the time
		user.Avatar = login.Image
		user.Username = login.FirstName + login.LastName
		user.Usertype = UserTypeNormal

		gErr = DB.Save(&user).Error
		if gErr != nil {
			http.Error(res, gErr.Error()+" - database query failed - Unable to save user's standard details.", http.StatusInternalServerError)
			return ""
		}

		social := Social{}
		social.FirstName = login.FirstName
		social.Network = login.Network
		social.LastName = login.LastName
		social.FullName = login.FullName
		social.NetworkUserId = login.Id
		newUser := User{}
		gErr = DB.Where("email = ?", login.Email).Find(&newUser).Error // find newly created user to get generated uid

		if gErr != nil {
			http.Error(res, gErr.Error()+" - could not recall user after creation", http.StatusInternalServerError)
			return ""
		}

		social.Uid = newUser.Uid

		gErr = DB.Save(&social).Error
		if gErr != nil {
			http.Error(res, gErr.Error()+" - database query failed - Unable to save user's social details.", http.StatusInternalServerError)
			return ""
		}

		session, err := SetSession(social.Uid)
		if err != nil {
			http.Error(res, err.Message, err.Code)
			return ""
		}

		u := map[string]interface{}{
			"uid":        user.Uid,
			"user":       newUser.Email,
			"email_hash": GetMD5Hash(newUser.Email),
			"session":    session.Value,
			"usertype":   newUser.Usertype,
		}

		usr, _ := json.Marshal(u)

		return string(usr)
	}

	http.Error(res, gErr.Error()+" - problem finding registered user", http.StatusInternalServerError)
	return ""
}

func HandleLogout(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	sid := req.Header.Get("X-API-SESSION")
	if len(sid) <= 0 {
		http.Error(res, "Missing session parameter.", http.StatusBadRequest)
		return ""
	}

	_, err := ClearSession(sid)
	if err != nil {
		http.Error(res, err.Message, err.Code)
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
	gErr := DB.Where("username = ?", register.Username).First(&user).Error
	if gErr != gorm.RecordNotFound {
		http.Error(res, "Username already exists.", http.StatusConflict)
		return ""
	}

	gErr = DB.Where("email = ?", register.Email).First(&user).Error
	if gErr != gorm.RecordNotFound {
		http.Error(res, "Email already exists.", http.StatusConflict)
		return ""
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(res, "Unable to generate password hash.", http.StatusInternalServerError)
		return ""
	}

	user.Email = register.Email
	user.Username = register.Username
	user.Password = string(hashedPassword)
	gErr = DB.Save(&user).Error
	if gErr != nil {
		http.Error(res, "Unable to create user.", http.StatusInternalServerError)
		return ""
	}

	var session *http.Cookie

	session, sErr := SetSession(user.Uid)
	if sErr != nil {
		http.Error(res, sErr.Message, sErr.Code)
		return ""
	}

	u := map[string]interface{}{
		"user":       user.Username,
		"email_hash": GetMD5Hash(user.Email),
		"session":    session.Value,
		"usertype":   user.Usertype,
	}

	usr, _ := json.Marshal(u)

	return string(usr)
}

func HandleCheckUsername(res http.ResponseWriter, req *http.Request, user UserNameForm) string {
	if user.Email == "" {
		http.Error(res, "Invalid or empty username.", http.StatusBadRequest)
		return ""
	}

	validUser := User{}
	gErr := DB.Where("email = ?", user.Email).First(&validUser).Error
	if gErr != nil && gErr != gorm.RecordNotFound {
		http.Error(res, "Unable to find that User.", http.StatusInternalServerError)
		return ""
	} else if gErr == gorm.RecordNotFound {
		http.Error(res, "We couldn't find an account associated with "+user.Email, http.StatusNotFound)
		return ""
	}

	u := map[string]interface{}{
		"user":   validUser.Email,
		"exists": true,
	}

	usr, _ := json.Marshal(u)

	return string(usr)
}

func HandleForgotPassword(res http.ResponseWriter, req *http.Request, user UserNameForm) string {
	if user.Email == "" {
		http.Error(res, "Invalid or empty username.", http.StatusBadRequest)
		return ""
	}

	validUser := User{}
	gErr := DB.Where("email = ?", user.Email).First(&validUser).Error
	if gErr != nil && gErr != gorm.RecordNotFound {
		http.Error(res, "Database query failed (User)", http.StatusInternalServerError)
		return ""
	} else if gErr == gorm.RecordNotFound {
		http.Error(res, "We couldn't find an account associated with "+user.Email, http.StatusNotFound)
		return ""
	}

	randomString := make([]byte, 64)
	_, err := rand.Read(randomString)

	if err != nil {
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

	gErr = DB.Save(&token).Error
	if gErr != nil {
		http.Error(res, "Database query failed (Token)", http.StatusInternalServerError)
		return ""
	}

	u := map[string]interface{}{
		"user":       validUser.Username,
		"email_hash": GetMD5Hash(validUser.Email),
		"token":      hash,
	}

	usr, _ := json.Marshal(u)

	return string(usr)
}

func ResetPassword(hash, email, password string) *appError {
	if email == "" || hash == "" {
		return &appError{nil, "No email/token found!", http.StatusBadRequest}
	}

	user := User{}
	gErr := DB.Where("email = ?", email).First(&user).Error
	if gErr != nil && gErr != gorm.RecordNotFound {
		return &appError{nil, "Database query failed (User).", http.StatusInternalServerError}
	} else if gErr == gorm.RecordNotFound {
		return &appError{nil, "Invalid email!", http.StatusNotFound}
	}

	token := UserTokens{}
	gErr = DB.Where("uid = ?", user.Uid).Where("hash = ?", hash).Where("used = ?", false).Last(&token).Error
	if gErr != nil && gErr != gorm.RecordNotFound {
		return &appError{nil, "Database query failed (Token).", http.StatusInternalServerError}
	} else if gErr == gorm.RecordNotFound {
		return &appError{nil, "Invalid token!", http.StatusNotFound}
	}

	if password == "" {
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return &appError{err, "Unable to generate password hash.", http.StatusInternalServerError}
	}

	user.Password = string(hashedPassword)
	gErr = DB.Save(&user).Error
	if gErr != nil {
		return &appError{gErr, "Database query failed (Password).", http.StatusInternalServerError}
	}

	token.Used = true
	gErr = DB.Save(&token).Error
	if gErr != nil {
		return &appError{gErr, "Database query failed (Token).", http.StatusInternalServerError}
	}

	return nil
}

func HandleResetPasswordCheck(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	err := ResetPassword(params["token"], params["email"], "")
	if err != nil {
		http.Error(res, err.Message, err.Code)
		return ""
	}

	return "OK"
}

func HandleResetPassword(res http.ResponseWriter, req *http.Request, params martini.Params, user UserForm) string {
	err := ResetPassword(params["token"], user.Email, user.Password)
	if err != nil {
		http.Error(res, err.Message, err.Code)
		return ""
	}

	return "OK"
}

func GetUserDetails(res http.ResponseWriter, req *http.Request, params martini.Params) string {
	session := req.Header.Get("X-API-SESSION")
	if len(session) <= 0 {
		http.Error(res, "Missing session parameter.", http.StatusBadRequest)
		return ""
	}

	user := User{}
	err := &appError{}
	if params["username"] != "" {
		user, err = GetUserDetailsByUsername(params["username"])
		if err != nil {
			http.Error(res, err.Message, err.Code)
		}
	} else {
		uid, err := GetUserID(session)
		if err != nil {
			http.Error(res, err.Message, err.Code)
			return ""
		}

		user, err = GetUserDetailsById(uid)
		if err != nil {
			http.Error(res, err.Message, err.Code)
		}
	}

	u := map[string]interface{}{
		"uid":        user.Uid,
		"username":   user.Username,
		"email":      user.Email,
		"email_hash": GetMD5Hash(user.Email),
		"avatar":     user.Avatar,
		"usertype":   user.Usertype,
	}

	usr, _ := json.Marshal(u)

	return string(usr)
}

func GetUserDetailsById(userid int) (User, *appError) {
	user := User{}

	if userid == 0 {
		return user, &appError{nil, "No username found!", http.StatusBadRequest}
	}

	gErr := DB.Where("uid = ?", userid).First(&user).Error
	if gErr != nil && gErr != gorm.RecordNotFound {
		return user, &appError{gErr, "Database query failed (User).", http.StatusInternalServerError}
	} else if gErr == gorm.RecordNotFound {
		return user, &appError{gErr, "No such user found!", http.StatusNotFound}
	}

	return user, nil
}

func GetUserDetailsByUsername(username string) (User, *appError) {
	user := User{}

	if username == "" {
		return user, &appError{nil, "No username found!", http.StatusBadRequest}
	}

	gErr := DB.Where("username = ?", username).First(&user).Error
	if gErr != nil && gErr != gorm.RecordNotFound {
		return user, &appError{gErr, "Database query failed (User).", http.StatusInternalServerError}
	} else if gErr == gorm.RecordNotFound {
		return user, &appError{gErr, "No such user found!", http.StatusNotFound}
	}

	return user, nil
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
	gErr := DB.Where("uid = ?", uid).First(&user).Error
	if gErr != nil && gErr != gorm.RecordNotFound {
		http.Error(res, "Database query failed (Select).", http.StatusInternalServerError)
		return ""
	} else if gErr == gorm.RecordNotFound {
		http.Error(res, "No such user found!", http.StatusNotFound)
		return ""
	}

	if data.Username != "" {
		user.Username = data.Username
	}

	if data.Email != "" {
		user.Email = data.Email
	}

	gErr = DB.Save(&user).Error
	if gErr != nil {
		http.Error(res, "Database query failed (Update).", http.StatusInternalServerError)
		return ""
	}

	u := map[string]interface{}{
		"username":   user.Username,
		"email":      user.Email,
		"email_hash": GetMD5Hash(user.Email),
	}

	usr, _ := json.Marshal(u)

	return string(usr)
}

func Reputation(uid int, points int) string {
	usr := User{}

	gErr := DB.Where("uid = ?", uid).Find(&usr).Error
	if gErr != nil {
		return "failed to find user"
	}

	r := usr.Reputation + points

	gErr = DB.Model(usr).Where("uid = ?", uid).Update("reputation", r).Error
	if gErr != nil {
		return "failed to update reputation"
	}

	return ""
}
