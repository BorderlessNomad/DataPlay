package main

import (
	bcrypt "code.google.com/p/go.crypto/bcrypt"
	"fmt"
	"net/http"
)

// Shortcut function that is used to redirect users to the login page if they are not logged in.
func checkAuth(res http.ResponseWriter, req *http.Request) {
	if !(IsUserLoggedIn(res, req)) {
		http.Redirect(res, req, "/login", http.StatusTemporaryRedirect)
		return
	}
}

func HandleLogin(res http.ResponseWriter, req *http.Request) {
	database := GetDB()
	defer database.Close()
	username := req.FormValue("username")
	password := req.FormValue("password")

	rows, e := database.Query("SELECT `password` FROM priv_users where email = ? LIMIT 1", username)
	check(e) // Check if the thing error's out
	rows.Next()

	var usrpassword string
	e = rows.Scan(&usrpassword)

	if usrpassword != "" && bcrypt.CompareHashAndPassword([]byte(usrpassword), []byte(password)) == nil { // Check the password with bcrypt
		var uid int
		e := database.QueryRow("SELECT uid FROM priv_users where email = ? LIMIT 1", username).Scan(&uid)
		check(e)
		e = SetSession(res, req, uid)
		if e != nil {
			http.Error(res, "Could not setup session.", http.StatusInternalServerError)
			return
		}

		http.Redirect(res, req, "/", http.StatusFound)
	} else {
		// Just in the case that the user is on a really old MD5 password (useful for admins resetting passwords too) check
		var md5test int
		e := database.QueryRow("SELECT count(*) FROM priv_users where email = ? AND password = MD5( ? ) LIMIT 1", username, password).Scan(&md5test)

		if e == nil {
			if md5test != 0 {
				// Ooooh, We need to upgrade this password!
				pwd, e := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
				if e == nil {
					database.Exec("UPDATE `DataCon`.`priv_users` SET `password`= ? WHERE `email`=?", pwd, username)

					var uid int
					e := database.QueryRow("SELECT uid FROM priv_users where email = ? LIMIT 1", username).Scan(&uid)
					check(e)
					e = SetSession(res, req, uid)
					if e != nil {
						http.Error(res, "Could not setup session.", http.StatusInternalServerError)
						return
					}

					http.Redirect(res, req, "/", http.StatusFound)
				}

				http.Redirect(res, req, fmt.Sprintf("/login?failed=3&r=%s", e), http.StatusFound)
			} else {
				http.Redirect(res, req, "/login?failed=1", http.StatusNotFound) // The user has failed this test as well :sad tuba:
			}
		} else {
			http.Redirect(res, req, "/login?failed=1", http.StatusNotFound) // Ditto to the above
		}
	}
}

func HandleLogout(res http.ResponseWriter, req *http.Request) {
	ClearSession(res, req)
	http.Redirect(res, req, "/login", http.StatusTemporaryRedirect)
	return
}

func HandleRegister(res http.ResponseWriter, req *http.Request) string {
	database := GetDB()
	defer database.Close()
	username := req.FormValue("username")
	password := req.FormValue("password")

	rows, e := database.Query("SELECT COUNT(*) FROM priv_users where email = ? LIMIT 1", username)
	check(e)
	rows.Next()

	var doesusrexist int
	e = rows.Scan(&doesusrexist)

	if doesusrexist == 0 {
		pwd, e := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if e != nil {
			return "The password you entered is invalid."
		}

		r, e := database.Exec("INSERT INTO `DataCon`.`priv_users` (`email`, `password`) VALUES (?, ?);", username, pwd)
		if e != nil {
			return "Could not make the user you requested."
		}

		newid, _ := r.LastInsertId()
		SetSession(res, req, int(newid))

		http.Redirect(res, req, "/", http.StatusFound)
		return ""
	} else {
		http.Error(res, "That username is already registered.", http.StatusConflict)
		return "That username is already registered."
	}
}
