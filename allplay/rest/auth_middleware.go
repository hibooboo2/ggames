package rest

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"
	"text/template"

	"github.com/hibooboo2/ggames/allplay/logger"
	"github.com/hibooboo2/ggames/allplay/pollen/constants"
	"github.com/hibooboo2/ggames/allplay/pollen/db"
)

const (
	userCtxKey = "USER_CTX_KEY"
)

func GetUsername(r *http.Request) string {
	user := r.Context().Value(userCtxKey)
	if user == nil {
		return ""
	}
	return user.(string)
}

func BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, c, isLoggedIn := db.IsLoggedIn(r, false)
		if !isLoggedIn {
			logger.Authln("Must login first")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), userCtxKey, username))

		http.SetCookie(w, c)
		next.ServeHTTP(w, r)
	})
}

var loginPage = template.Must(template.New("login").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta name="viewport" content="width=device-width, initial-scale=1">
<link rel="stylesheet" href="/static/css/login.css">
<script src="/static/js/functions.js"></script>
</head>
<body>

<h2>Login To Allplay Games</h2>

<div class="center">
	<form action="/login" method="post">
	<div class="imgcontainer">
		<img src="/static/images/Back_Green.png" alt="Pollen card" class="avatar">
	</div>

	<div class="container">
		{{ if .IsLoggedIn }}
			<button type="button" class="logout" onclick='window.location.href="/logout"'>
				Logout
			</button>
		{{ else }}
			<label for="uname"><b>Username</b></label>
			<input type="text" placeholder="Enter Username" name="uname" required>

			<label for="psw"><b>Password</b></label>
			<input type="password" placeholder="Enter Password" name="psw" required>
				
			{{ if .Invalid }}
				<label class="warning"><b>Invalid username or password</b></label>
			{{end}}
			<button type="submit">Login</button>
			<label>
			<input type="checkbox" checked="checked" name="remember"> Remember me
			</label>
		{{ end }}
	</div>

	<div class="container" style="background-color:#f1f1f1">
		<button type="button" class="cancelbtn">Cancel</button>
		<span class="register">No account? <a href="/register">Signup</a></span>
		<span class="psw">Forgot <a href="#">password?</a></span>
	</div>
	</form>
</div>

</body>
</html>

`))

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch strings.ToUpper(r.Method) {
	case http.MethodGet:
		_, _, isLoggedIn := db.IsLoggedIn(r, false)
		loginPage.Execute(w, struct {
			Invalid    bool
			IsLoggedIn bool
		}{
			r.URL.Query().Get("invalid") == "true",
			isLoggedIn,
		})
	case http.MethodPost:
		username := r.FormValue("uname")
		password := r.FormValue("psw")
		if username == "" || password == "" {
			http.Redirect(w, r, "/login?invalid=true", http.StatusSeeOther)
			return
		}

		us, err := db.Login(username, password, constants.SessionTimeout, true)
		if err != nil {
			http.Redirect(w, r, "/login?invalid=true", http.StatusSeeOther)
			return
		}

		http.SetCookie(w, us.Cookie())
		logger.Authln("Logged in", username)
		http.Redirect(w, r, "/", http.StatusFound)
	}

}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   constants.SessionCookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/login", http.StatusFound)
}

func basicAuthEnc(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func Register(w http.ResponseWriter, r *http.Request) {
	switch strings.ToUpper(r.Method) {
	case http.MethodPost:
		username := r.FormValue("username")
		password := r.FormValue("psw")
		passwordRepeat := r.FormValue("psw-repeat")
		email := r.FormValue("email")
		if email == "" {
			http.Error(w, "Email is required", http.StatusBadRequest)
			return
		}
		if username == "" {
			http.Error(w, "Username is required", http.StatusBadRequest)
			return
		}
		if password == "" {
			http.Error(w, "Password is required", http.StatusBadRequest)
			return
		}
		if password != passwordRepeat {
			http.Error(w, "Passwords do not match", http.StatusBadRequest)
			return
		}

		err := db.RegisterUser(email, username, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/login", http.StatusFound)
	case http.MethodGet:
		w.Write([]byte(`
		<!DOCTYPE html>
			<html>
			<head>
				<meta name="viewport" content="width=device-width, initial-scale=1">
				<link rel="stylesheet" href="/static/css/register.css">
			</head>
			<body>

			<form action="/register" method="post">
				<div class="container">
					<h1>Register</h1>
					<p>Please fill in this form to create an account.</p>
					<hr>

					<label for="email"><b>Email</b></label>
					<input type="text" placeholder="Enter Email" name="email" id="email" required>

					<label for="username"><b>Username</b></label>
					<input type="text" placeholder="Enter Username" name="username" id="username" required>

					<label for="psw"><b>Password</b></label>
					<input type="password" placeholder="Enter Password" name="psw" id="psw" required>

					<label for="psw-repeat"><b>Repeat Password</b></label>
					<input type="password" placeholder="Repeat Password" name="psw-repeat" id="psw-repeat" required>
					<hr>
					<p>By creating an account you agree to our <a href="#">Terms & Privacy</a>. Implement This</p>

					<button type="submit" class="registerbtn">Register</button>
				</div>
			
				<div class="container signin">
					<p>Already have an account? <a href="/login">Sign in</a>.</p>
				</div>
			</form>

			</body>
			</html>

		`))
	}
}

func TempID(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(db.GetTempID()))
}
