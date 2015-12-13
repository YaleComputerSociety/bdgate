package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

var form = `
<html>
<head>
<title>Sign Up!</title>
</head>
<body>
<form method="POST" action="/signup/post" accept-charset="UTF-8">
<input type="text" name="name">
<input type="text" name="email">
<!--
    The default template tag used by the CSRF middleware .
    This will be replaced with a hidden <input> field containing the
    masked CSRF token.
-->
{{ .csrfField }}
<input type="submit" value="Sign up!">
</form>
</body>
</html>
`

var t = template.Must(template.New("signup_form.tmpl").Parse(form))

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/signup", ShowSignupForm)
	// All POST requests without a valid token will return HTTP 403 Forbidden.
	r.HandleFunc("/signup/post", SubmitSignupForm)

	// Add the middleware to your router by wrapping it.
	http.ListenAndServe(":8000",
		csrf.Protect([]byte("3byte-long-auth-key-asdfasdfsdaf"))(r))
}

func ShowSignupForm(w http.ResponseWriter, r *http.Request) {
	// signup_form.tmpl just needs a {{ .csrfField }} template tag for
	// csrf.TemplateField to inject the CSRF token into. Easy!
	t.ExecuteTemplate(w, "signup_form.tmpl", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

func SubmitSignupForm(w http.ResponseWriter, r *http.Request) {
	// We can trust that requests making it this far have satisfied
	// our CSRF protection requirements.
	fmt.Fprintf(w, "%v\n", r.PostForm)
}
