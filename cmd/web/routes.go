package main

import (
	"net/http"

	"github.com/bmizerany/pat" // New import
	"github.com/justinas/alice"
	"golang.org/x/time/rate"
)

// NewRateLimiter creates a new RateLimiter
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		limiter: rate.NewLimiter(r, b),
	}
}

// RateLimiter struct to hold the rate limiter
type RateLimiter struct {
	limiter *rate.Limiter
}

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// Use the nosurf middleware on all our 'dynamic' routes.
	dynamicMiddleware := alice.New(app.session.Enable, noSurf)

	rateLimiter := NewRateLimiter(1, 15)
	// Middleware to apply rate limiting
	rateLimiterMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if rateLimiter.limiter.Allow() {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			}
		})
	}

	mux := pat.New()

	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/about", dynamicMiddleware.ThenFunc(app.about))

	mux.Get("/reviews", dynamicMiddleware.ThenFunc(app.reviews))
	mux.Get("/reviews/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createReviewForm))
	mux.Post("/reviews/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createReview))
	mux.Get("/reviews/delete", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.deleteReviewForm))
	mux.Post("/reviews/delete", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.deleteReview))

	mux.Get("/services", dynamicMiddleware.ThenFunc(app.servicess))
	mux.Get("/services/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createServiceForm))
	mux.Post("/services/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createService))
	mux.Get("/services/update", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.updateServiceForm))
	mux.Post("/services/update", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.updateService))
	mux.Get("/services/delete", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.deleteServiceForm))
	mux.Post("/services/delete", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.deleteService))
	mux.Get("/services/:id", dynamicMiddleware.ThenFunc(app.showService))
	mux.Get("/services/:sort", dynamicMiddleware.ThenFunc(app.showService))

	mux.Get("/appointments", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.appointmentss))
	mux.Get("/appointments/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createAppointmentForm))
	mux.Post("/appointments/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createAppointment))
	mux.Get("/appointments/update", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.updateAppointmentForm))
	mux.Post("/appointments/update", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.updateAppointment))
	mux.Get("/appointments/delete", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.deleteAppointmentForm))
	mux.Post("/appointments/delete", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.deleteAppointment))
	mux.Get("/appointments/:filter", dynamicMiddleware.ThenFunc(app.showService))

	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/confirm", dynamicMiddleware.ThenFunc(app.confirmUserForm))
	mux.Post("/user/confirm", dynamicMiddleware.ThenFunc(app.confirmUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	// Apply the rate limiter middleware to all routes
	return standardMiddleware.Append(rateLimiterMiddleware).Then(mux)
}
