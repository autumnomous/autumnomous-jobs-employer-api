package route

import (
	"net/http"

	"autumnomous-jobs-employer-api/controller/v1/employers"
	"autumnomous-jobs-employer-api/controller/v1/utilities"
	"autumnomous-jobs-employer-api/route/middleware/acl"
	"autumnomous-jobs-employer-api/route/middleware/cors"
	hr "autumnomous-jobs-employer-api/route/middleware/httprouterwrapper"
	"autumnomous-jobs-employer-api/route/middleware/logrequest"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// LoadRoutes returns the routes and middleware
func LoadRoutes() http.Handler {
	//return routes()
	return middleware(routes())
}

// *****************************************************************************
// Routes
// *****************************************************************************

func routes() *httprouter.Router {
	r := httprouter.New()

	r.POST("/upload/image", hr.Handler(alice.New(acl.AllowAPIKey).ThenFunc(utilities.UploadImage)))

	r.POST("/employer/signup", hr.Handler(alice.New(acl.AllowAPIKey).ThenFunc(employers.SignUp)))
	r.POST("/employer/login", hr.Handler(alice.New(acl.AllowAPIKey).ThenFunc(employers.Login)))

	r.POST("/employer/update-password", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(employers.UpdatePassword)))
	r.POST("/employer/update-account", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(employers.UpdateAccount)))
	r.POST("/employer/update-company", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(employers.UpdateCompany)))
	r.POST("/employer/update-payment-method", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(employers.UpdatePaymentMethod)))
	r.POST("/employer/update-payment-details", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(employers.UpdatePaymentDetails)))

	r.GET("/employer/get", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(employers.GetEmployer)))
	r.GET("/employer/get/company", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(employers.GetEmployerCompany)))
	r.POST("/employer/create/job", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(employers.CreateJob)))
	r.POST("/employer/edit/job", hr.Handler(alice.New(acl.ValidateJWT, acl.ValidateJWT).ThenFunc(employers.EditJob)))
	r.GET("/employer/get/jobs", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(employers.GetJobs)))
	r.POST("/employer/get/job", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(employers.GetJob)))
	r.DELETE("/employer/delete/job", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(employers.DeleteJob)))
	r.GET("/employer/get/jobpackages/active", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(employers.GetActiveJobPackages)))
	r.POST("/employer/get/location/autocomplete", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(employers.GetAutocompleteLocationData)))

	r.POST("/employer/buy/job-package", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(employers.PurchaseJobPackage)))

	// r.POST("/get-user", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(users.GetUser)))

	// r.GET("/get/client/registration", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(clients.CheckRegistration)))
	// r.POST("/set/client/registration", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(clients.SetRegistration)))

	// r.GET("/customers/:id", hr.Handler(alice.New(acl.ValidateJWT).ThenFunc(clients.GetClientCustomers)))
	// Enable Pprof
	// r.GET("/debug/pprof/*pprof", hr.Handler(alice.
	// 	New(acl.ValidateJWT).
	// 	ThenFunc(pprofhandler.Handler)))

	return r
}

// *****************************************************************************
// Middleware
// *****************************************************************************

func middleware(h http.Handler) http.Handler {
	// Log every request
	h = logrequest.Handler(h)

	// Cors for swagger-ui
	h = cors.Handler(h)

	// Clear handler for Gorilla Context
	h = context.ClearHandler(h)

	return h
}
