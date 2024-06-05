package main

import (
	"errors"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"strconv"
	"time"

	"aitunews.kz/snippetbox/pkg/forms"
	"aitunews.kz/snippetbox/pkg/models"
)

func (app *application) showService(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	s, err := app.services.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{
		Service: s,
	})

}

func (app *application) createServiceForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		// Pass a new empty forms.Form object to the template.
		Form: forms.New(nil),
	})
}

func (app *application) createService(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content", "master", "price")
	form.MaxLength("title", 100)
	form.MaxLength("master", 100)
	//form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	i, err1 := strconv.Atoi(form.Get("price"))
	if err1 != nil {
		return
	}

	id, err := app.services.Insert(form.Get("title"), form.Get("content"), form.Get("master"), i)
	if err != nil && id == 0 {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "flash", "Service successfully created!")

	http.Redirect(w, r, "/services", http.StatusSeeOther)

}

func (app *application) updateServiceForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "update.page.tmpl", &templateData{
		// Pass a new empty forms.Form object to the template.
		Form: forms.New(nil),
	})
}

func (app *application) updateService(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "price")

	if !form.Valid() {
		app.render(w, r, "update.page.tmpl", &templateData{Form: form})
		return
	}

	i, err1 := strconv.Atoi(form.Get("price"))
	if err1 != nil {
		return
	}

	err = app.services.Update(form.Get("title"), i)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "flash", "Service successfully updated!")

	http.Redirect(w, r, "/services", http.StatusSeeOther)

}

func (app *application) deleteServiceForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "delete.page.tmpl", &templateData{
		// Pass a new empty forms.Form object to the template.
		Form: forms.New(nil),
	})
}

func (app *application) deleteService(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title")

	if !form.Valid() {
		app.render(w, r, "delete.page.tmpl", &templateData{Form: form})
		return
	}

	err = app.services.Delete(form.Get("title"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "flash", "Service successfully deleted!")

	http.Redirect(w, r, "/services", http.StatusSeeOther)

}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})
}

func (app *application) servicess(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	sort_type := r.URL.Query().Get("sort_type")

	if r.URL.Path != "/services" {
		app.notFound(w)
		return
	}
	s, err := app.services.Latest("services", sort, sort_type)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "services.page.tmpl", &templateData{
		Services: s,
	})
}

func (app *application) appointmentss(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := 10 // Number of items per page
	offset := 0 // Offset for SQL query

	if p, err := strconv.Atoi(page); err == nil && p > 1 {
		offset = (p - 1) * limit
	}

	if r.URL.Path != "/appointments" {
		app.notFound(w)
		return
	}
	user_id := app.session.GetInt(r, "authenticatedUserID")

	s, err := app.appointments.Latest(user_id, limit, offset)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "appointments.page.tmpl", &templateData{
		Appointments: s,
	})
}

func (app *application) createAppointmentForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "createAppointment.page.tmpl", &templateData{
		// Pass a new empty forms.Form object to the template.
		Form: forms.New(nil),
	})
}

func (app *application) createAppointment(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("service_id", "time")

	if !form.Valid() {
		app.render(w, r, "createAppointment.page.tmpl", &templateData{Form: form})
		return
	}

	user_id := app.session.GetInt(r, "authenticatedUserID")

	id, err := app.appointments.Insert(user_id, form.Get("service_id"), form.Get("time"))
	if err != nil && id == 0 {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "flash", "Good!")

	http.Redirect(w, r, "/appointments", http.StatusSeeOther)

}

func (app *application) updateAppointmentForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "updateAppointment.page.tmpl", &templateData{
		// Pass a new empty forms.Form object to the template.
		Form: forms.New(nil),
	})
}

func (app *application) updateAppointment(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("appointment_id", "time")

	if !form.Valid() {
		app.render(w, r, "updateAppointment.page.tmpl", &templateData{Form: form})
		return
	}

	i, err1 := strconv.Atoi(form.Get("appointment_id"))
	if err1 != nil {
		return
	}

	err = app.appointments.Update(i, form.Get("time"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "flash", "Appointment successfully updated!")

	http.Redirect(w, r, "/appointments", http.StatusSeeOther)

}

func (app *application) deleteAppointmentForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "deleteAppointment.page.tmpl", &templateData{
		// Pass a new empty forms.Form object to the template.
		Form: forms.New(nil),
	})
}

func (app *application) deleteAppointment(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("appointment_id")

	if !form.Valid() {
		app.render(w, r, "deleteAppointment.page.tmpl", &templateData{Form: form})
		return
	}

	i, err1 := strconv.Atoi(form.Get("appointment_id"))
	if err1 != nil {
		return
	}

	err = app.appointments.Delete(i)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "flash", "Appointment successfully deleted!")

	http.Redirect(w, r, "/appointments", http.StatusSeeOther)

}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/about" {
		app.notFound(w)
		return
	}
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "about.page.tmpl", &templateData{
		Snippets: s,
	})
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {

	rand.Seed(time.Now().UnixNano())

	randomInt := rand.Intn(9000) + 1000 // Generates a random integer in the range [0, 100]
	confirmation := false

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("full_name", "email", "password", "phone", "role")
	form.MaxLength("full_name", 255)
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)
	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}
	err = app.users.Insert(form.Get("full_name"), form.Get("email"), form.Get("phone"), form.Get("password"), form.Get("role"), randomInt, confirmation)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Address is already in use")
			app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.session.Put(r, "flash", "Your signup was successful. Please confirm your email.")

	// Sender's email credentials
	from := "daniyar.0586@gmail.com"
	password := "fkmm flqo zpha zddd"

	// Receiver's email address
	to := []string{form.Get("email")}

	// SMTP server configuration
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message
	msg := []byte("To: " + form.Get("email") + "\r\n" +
		"Subject: Email Confirmation\r\n" +
		"\r\n" +
		strconv.Itoa(randomInt))

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Email sent successfully!")

	http.Redirect(w, r, "/user/confirm", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Check whether the credentials are valid. If they're not, add a generic error
	// message to the form failures map and re-display the login page.
	form := forms.New(r.PostForm)

	confirmation := app.users.ConfirmCheck(form.Get("email"))

	if !confirmation {
		http.Redirect(w, r, "/user/confirm", http.StatusSeeOther)
		return
	}

	code := app.users.CodeCheck(form.Get("email"))

	if code != "0" {
		http.Redirect(w, r, "/user/confirm", http.StatusSeeOther)
		return
	}

	id, role, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password is incorrect")
			app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}
	// Add the ID of the current user to the session, so that they are now 'logged in'.

	randomInt := rand.Intn(9000) + 1000

	from := "daniyar.0586@gmail.com"
	password := "fkmm flqo zpha zddd"

	// Receiver's email address
	to := []string{form.Get("email")}

	// SMTP server configuration
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message
	msg := []byte("To: " + form.Get("email") + "\r\n" +
		"Subject: Code Confirmation\r\n" +
		"\r\n" +
		strconv.Itoa(randomInt))

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Email sent successfully!")

	app.users.AddCode(form.Get("email"), strconv.Itoa(randomInt))

	http.Redirect(w, r, "/user/confirm", http.StatusSeeOther)

	app.session.Put(r, "authenticatedUserID", id)
	app.session.Put(r, "authenticatedUserRole", role)
	// Redirect the user to the create snippet page.
	// http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) confirmUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "confirmation.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) confirmUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)

	confirmed := app.users.ConfirmCheck(form.Get("email"))
	if !confirmed {
		err = app.users.Confirm(form.Get("email"), form.Get("code"))
		if err != nil {
			app.serverError(w, err)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		app.users.AddCode(form.Get("email"), "0")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	// Remove the authenticatedUserID from the session data so that the user is
	// 'logged out'.
	app.session.Remove(r, "authenticatedUserID")
	app.session.Remove(r, "authenticatedUserRole")
	// Add a flash message to the session to confirm to the user that they've been
	// logged out.
	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

/////////////

func (app *application) reviews(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter")

	if r.URL.Path != "/reviews" {
		app.notFound(w)
		return
	}
	user_id := app.session.GetInt(r, "authenticatedUserID")

	s, err := app.snippets.Latestt(user_id, filter)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "reviews.page.tmpl", &templateData{
		Snippets: s,
	})
}

func (app *application) createReviewForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "createReview.page.tmpl", &templateData{
		// Pass a new empty forms.Form object to the template.
		Form: forms.New(nil),
	})
}

func (app *application) createReview(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("content")

	if !form.Valid() {
		app.render(w, r, "createReview.page.tmpl", &templateData{Form: form})
		return
	}

	user_id := app.session.GetInt(r, "authenticatedUserID")
	currentTime := time.Now()
	created := currentTime.Format("02.01.2006 15:04:05")

	id, err := app.snippets.Insert(user_id, form.Get("content"), created)
	if err != nil && id == 0 {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "flash", "Service successfully created!")

	http.Redirect(w, r, "/reviews", http.StatusSeeOther)

}

func (app *application) deleteReviewForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "deleteReview.page.tmpl", &templateData{
		// Pass a new empty forms.Form object to the template.
		Form: forms.New(nil),
	})
}

func (app *application) deleteReview(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("review_id")

	if !form.Valid() {
		app.render(w, r, "deleteReview.page.tmpl", &templateData{Form: form})
		return
	}

	i, err1 := strconv.Atoi(form.Get("review_id"))
	if err1 != nil {
		return
	}

	err = app.snippets.Delete(i)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "flash", "Review successfully deleted!")

	http.Redirect(w, r, "/reviews", http.StatusSeeOther)

}
