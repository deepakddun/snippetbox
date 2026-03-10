package main

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"

	"github.com/deepakddun/snippetbox/internal/models"
	"github.com/deepakddun/snippetbox/internal/validator"
)

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

type UsersCreateForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type UserLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippets.Latest(r.Context())

	if err != nil {

		//app.logger.Error(err.Error())
		app.serverError(w, r, err)
		return
	}
	//flash := app.sessionManager.PopString(r.Context(), "flash")
	data := app.newTemplateData(r)
	data.Snippets = snippets
	app.render(w, r, "home.tmpl", http.StatusOK, data)
	// err = ts.ExecuteTemplate(w, "base", data)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }

	//json.NewEncoder(w).Encode(snippets)

	// files := []string{
	// 	"./../../ui/html/pages/base.tmpl",
	// 	"./../../ui/html/pages/home.tmpl",
	// 	"./../../ui/html/pages/nav.tmpl",
	// }

	// tp, err := template.ParseFiles(files...)
	// // if tp.Lookup("base") == nil {
	// // 	log.Println(`Lookup("base") = nil (NOT defined)`)
	// // } else {
	// // 	log.Println(`Lookup("base") OK`)
	// // }

	// // for _, t := range tp.Templates() {
	// // 	log.Println("template:", t.Name())
	// // }

	// // for _, f := range files {
	// // 	abs, _ := filepath.Abs(f)
	// // 	log.Println("using file:", abs)
	// // }

	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }

	// err = tp.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }
	//w.WriteHeader()
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Starting view ...")
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil || id < 1 {
		app.serverError(w, r, err)
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippets.Get(r.Context(), id)

	if err != nil {
		if errors.Is(err, models.ErrorNoRecord) {

			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)

		}
		return
	}
	flash := app.sessionManager.PopString(r.Context(), "flash")
	data := templateData{

		Snippet:     snippet,
		CurrentYear: app.getCurrentYear(),
	}
	data.Flash = flash
	app.render(w, r, "view.tmpl", http.StatusOK, data)

}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	data := templateData{

		Snippet:     models.Snippet{},
		CurrentYear: app.getCurrentYear(),
	}
	data.Form = snippetCreateForm{
		Expires: 7,
	}
	data.IsAuthenticated = app.isAutheticated(r)
	app.render(w, r, "create.tmpl", http.StatusOK, data)

}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	// 	w.Header().Add("Server", "Go")
	// 	//w.WriteHeader(http.StatusCreated)
	// 	//w.WriteHeader(200)

	// 	title := "O snail"
	// 	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	// 	expires := 7
	// 	id, err := app.snippets.Insert(r.Context(), title, content, expires)

	// 	app.logger.Info("Created Record With ID", "Id ", id)

	// 	if err != nil {
	// 		app.serverError(w, r, err)
	// 		return
	// 	}
	// 	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	// 	//fmt.Fprintf(w, "Creating SnippetBox .Post..")
	// }

	err := r.ParseForm()

	if err != nil {

		app.ClientError(w, r, http.StatusBadRequest)
		return
	}
	var form snippetCreateForm
	// title := r.PostForm.Get("title")
	// content := r.PostForm.Get("content")
	// expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	err = app.formDecoder.Decode(&form, r.PostForm)
	if err != nil {

		app.ClientError(w, r, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")
	fmt.Println("%v", form)
	if !form.Valid() {
		data := templateData{

			CurrentYear: app.getCurrentYear(),
		}
		data.Form = form
		app.render(w, r, "create.tmpl", http.StatusUnprocessableEntity, data)
		return
	}
	id, err := app.snippets.Insert(r.Context(), form.Title, form.Content, form.Expires)

	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintln(w, "Display a form for signing up a new user...")

	data := templateData{
		CurrentYear: app.getCurrentYear(),
	}
	data.Form = UsersCreateForm{}
	app.render(w, r, "signup.tmpl", http.StatusOK, data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintln(w, "Create a new user...")
	err := r.ParseForm()

	if err != nil {
		app.ClientError(w, r, http.StatusBadRequest)
		return
	}

	var form UsersCreateForm

	err = app.formDecoder.Decode(&form, r.PostForm)

	if err != nil {
		app.ClientError(w, r, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "Name cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "Email cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "Password cannot be blank")
	form.CheckField(validator.EmailCheck(form.Email, validator.EmailRX), "email", "Email format is not valid")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Password needs to be atleast 8 characters long")

	if !form.Valid() {

		data := templateData{
			CurrentYear: app.getCurrentYear(),
		}
		data.Form = form
		app.render(w, r, "signup.tmpl", http.StatusUnprocessableEntity, data)
		return
	}

	err = app.users.Insert(r.Context(), form.Name, form.Email, form.Password)
	//fmt.Println(err)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {

			form.AddFieldError("email", "Email already used ")
			data := templateData{
				CurrentYear: app.getCurrentYear(),
			}
			data.Form = form
			app.render(w, r, "signup.tmpl", http.StatusUnprocessableEntity, data)

		} else {
			app.serverError(w, r, err)

		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "User created successfully")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)

}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintln(w, "Display a form for logging in a user...")

	data := templateData{
		CurrentYear: app.getCurrentYear(),
	}
	data.Form = UserLoginForm{}
	fmt.Println("User Signup")
	app.render(w, r, "login.tmpl", http.StatusOK, data)

}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "Authenticate and login the user...")

	var form UserLoginForm

	err := r.ParseForm()

	if err != nil {
		app.ClientError(w, r, http.StatusBadRequest)
		return
	}

	err = app.formDecoder.Decode(&form, r.PostForm)

	if err != nil {
		app.ClientError(w, r, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.EmailCheck(form.Email, validator.EmailRX), "email", "This field must be valid email address")

	if !form.Valid() {
		data := templateData{}
		data.CurrentYear = app.getCurrentYear()
		data.Form = form
		app.render(w, r, "login.tmpl", http.StatusUnprocessableEntity, data)
		return
	}
	fmt.Printf("email=%q type=%T\n", form.Email, form.Email)
	fmt.Printf("password=%q type=%T\n", form.Password, form.Password)
	id, err := app.users.Authenticate(r.Context(), form.Email, form.Password)

	if err != nil {

		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := templateData{}
			data.CurrentYear = app.getCurrentYear()
			data.Form = form
			data.IsAuthenticated = app.isAutheticated(r)
			app.render(w, r, "login.tmpl", http.StatusUnprocessableEntity, data)

		} else {
			app.serverError(w, r, err)
		}

	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintln(w, "Logout the user...")

	err := app.sessionManager.RenewToken(r.Context())

	if err != nil {
		debug.PrintStack()
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	app.sessionManager.Put(r.Context(), "flash", "You've have been logged out successfully")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
