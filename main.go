package main

import (
	"net/http"
	"os"
	"path/filepath"
	"prestamosbackend/controllers"
	"prestamosbackend/initializers"
	"prestamosbackend/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
	initializers.CreateData()
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/dir", func(c *gin.Context) {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"dir": dir,
			})
		}
	})

	r.GET("/dashboard/:id", middleware.RequireAuth, controllers.Dashboard)

	//Company routes
	r.POST("/company", controllers.CreateCompany)

	//User routes
	r.POST("/user/signup", controllers.SignUp)
	r.POST("/user/login", controllers.Login)

	r.POST("/user/changeUserPassword", middleware.RequireAuth, controllers.ChangePassword)
	r.GET("/user/loadUser", middleware.RequireAuth, controllers.LoadUser)
	r.PUT("/user/changeUserImage", middleware.RequireAuth, controllers.UpdateUserImage)

	r.GET("/image/user/:id", controllers.GetUserImage)

	//Customer routes
	r.POST("/customer", middleware.RequireAuth, controllers.CreateCustomer)
	r.PUT("/customer", middleware.RequireAuth, controllers.UpdateCustomer)
	r.GET("/customer", middleware.RequireAuth, controllers.SearchCustomerByParameter)
	r.GET("/customer/:id", middleware.RequireAuth, controllers.SearchCustomerById)

	//Route routes
	r.POST("/route", middleware.RequireAuth, controllers.CreateRoute)
	r.GET("/route", middleware.RequireAuth, controllers.SearchRouteByParameter)
	r.GET("/routes", middleware.RequireAuth, controllers.GetAllRoutes)
	r.GET("/route/:id", middleware.RequireAuth, controllers.SearchRouteById)
	r.PUT("/route", middleware.RequireAuth, controllers.EditRoute)

	//Loan
	r.POST("/loan/confirmation", middleware.RequireAuth, controllers.LoanConfirmation)
	r.POST("/loan", middleware.RequireAuth, controllers.CreateLoan)
	r.GET("/loan", middleware.RequireAuth, controllers.SearchLoanByParameter)
	r.GET("/loanDate", middleware.RequireAuth, controllers.GetLoansByDate)
	r.GET("/loanParameter", middleware.RequireAuth, controllers.GetLoansByParameter)
	r.GET("/loan/:id", middleware.RequireAuth, controllers.SearchLoanById)
	r.PUT("/loan/:id", middleware.RequireAuth, controllers.PayOffLoan)
	r.PUT("/loan/setLoanToPaid/:id", middleware.RequireAuth, controllers.SetLoanToPaid)
	r.POST("/loan/reNew/confirmation", middleware.RequireAuth, controllers.RenewLoanConfirmation)
	r.POST("/loan/reNew", middleware.RequireAuth, controllers.CreateRenewLoan)

	//Fee routes
	r.PUT("/fee/payOffFee", middleware.RequireAuth, controllers.PayOffFee)
	r.GET("/fee", middleware.RequireAuth, controllers.GetFeesByDate)

	// Spent routes
	r.POST("/spent", middleware.RequireAuth, controllers.CreateSpent)
	r.GET("/spent/:id", middleware.RequireAuth, controllers.GetSpent)
	r.GET("/spent", middleware.RequireAuth, controllers.GetAllSpents)

	// Spent Type routes
	r.GET("/spent/type", middleware.RequireAuth, controllers.GetAllTypesSpent)
	r.POST("/spent/type", middleware.RequireAuth, controllers.CreateTypeSpent)

	err := r.Run()
	if err != nil {
		return
	}
}
