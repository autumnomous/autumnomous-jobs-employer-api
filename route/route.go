package route

import (
	"net/http"

	"bit-jobs-api/controller/v1/employers"
	"bit-jobs-api/route/middleware/acl"
	"bit-jobs-api/route/middleware/cors"
	hr "bit-jobs-api/route/middleware/httprouterwrapper"
	"bit-jobs-api/route/middleware/logrequest"

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

	r.POST("/employer/signup", hr.Handler(alice.New(acl.AllowAPIKey).ThenFunc(employers.SignUp)))
	r.POST("/employer/login", hr.Handler(alice.New(acl.AllowAPIKey).ThenFunc(employers.Login)))
	r.POST("/employer/update-password", hr.Handler(alice.New(acl.DisallowAnon).ThenFunc(employers.UpdatePassword)))
	r.GET("/employer/get", hr.Handler(alice.New(acl.DisallowAnon).ThenFunc(employers.GetEmployer)))
	r.POST("/employer/create/job", hr.Handler(alice.New(acl.DisallowAnon).ThenFunc(employers.CreateJob)))
	r.POST("/employer/edit/job", hr.Handler(alice.New(acl.DisallowAnon).ThenFunc(employers.EditJob)))
	r.GET("/employer/get/jobs", hr.Handler(alice.New(acl.DisallowAnon).ThenFunc(employers.GetJobs)))
	r.POST("/employer/get/job", hr.Handler(alice.New(acl.DisallowAnon).ThenFunc(employers.GetJob)))
	r.DELETE("/employer/delete/job", hr.Handler(alice.New(acl.DisallowAnon).ThenFunc(employers.DeleteJob)))
	r.GET("/employer/get/jobpackages/active", hr.Handler(alice.New(acl.DisallowAnon).ThenFunc(employers.GetActiveJobPackages)))

	// r.POST("/get-user", hr.Handler(alice.New(acl.DisallowAnon).ThenFunc(users.GetUser)))

	// r.GET("/get/client/registration", hr.Handler(alice.New(acl.DisallowAnon).ThenFunc(clients.CheckRegistration)))
	// r.POST("/set/client/registration", hr.Handler(alice.New(acl.DisallowAnon).ThenFunc(clients.SetRegistration)))

	// r.GET("/customers/:id", hr.Handler(alice.New(acl.DisallowAnon).ThenFunc(clients.GetClientCustomers)))
	// Enable Pprof
	// r.GET("/debug/pprof/*pprof", hr.Handler(alice.
	// 	New(acl.DisallowAnon).
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
