package routes

import (
	// Import your API layer (Ensure module name matches go.mod)
	"my-course-backend/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter initializes the Gin engine and defines all routes
func SetupRouter() *gin.Engine {
	// Initialize Gin with default middleware (Logger and Recovery)
	r := gin.Default()

	// 1. Configure CORS (Cross-Origin Resource Sharing)
	// This logic was moved here from main.go to keep the entry point clean.
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true // Warning: Allow all origins for development; restrict in production.
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	// Apply CORS middleware globally
	r.Use(cors.New(config))

	// 2. Register Route Groups

	// Auth Route Group
	// Prefix: /auth
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", api.Register)
		// CHANGED: SuperManager creates manager invite codes
		authRoutes.POST("/manager/invite-codes", api.CreateManagerInviteCode)
		authRoutes.POST("/manager/register", api.ManagerRegister)
		authRoutes.POST("/login", api.Login)

		// New Profile Endpoints
		authRoutes.GET("/profile", api.GetProfile)
		authRoutes.PUT("/profile", api.UpdateProfile)

		authRoutes.POST("/roles/assign", api.AssignUserRole)
	}
	// User Route Group
	// Prefix: /users
	userRoutes := r.Group("/users")
	{
		// DELETE /users/:id (e.g., /users/1)
		userRoutes.DELETE("/:id", api.DeleteUser)
		// GET /users/:id/enrollments (authenticated user can only get their own enrolled courses)
		userRoutes.GET("/:id/enrollments", api.GetUserEnrolledClasses)
		// GET /users/:id/enrollment-summary?days=30 (authenticated user can only get their own summary)
		//userRoutes.GET("/:id/enrollment-summary", api.GetUserEnrollmentSummary)
		// GET /users/:id/enrollment-summary/{7days|1mon|3mon}
		userRoutes.GET("/:id/analytics", api.GetUserAnalytics)
	}

	// Class Route Group
	// Prefix: /classes
	classRoutes := r.Group("/classes")
	{
		// public
		classRoutes.GET("", api.ListClasses)
		classRoutes.GET("/:id", api.GetClass)
		// manager-only
		classRoutes.GET("/:id/enrollments", api.ListClassEnrollments)

		// enrollment actions
		classRoutes.POST("/register", api.RegisterClass)
		classRoutes.POST("/drop", api.DropClass)

		// manager-only
		classRoutes.POST("", api.ManagerCreateClass)
		classRoutes.PUT("/:id", api.ManagerUpdateClass)
		classRoutes.DELETE("/:id", api.ManagerDeleteClass)
	}

	instructorRoutes := r.Group("/instructor")
	{
		instructorRoutes.GET("/courses", api.InstructorListCourses)
		instructorRoutes.GET("/courses/:id/enrollments", api.InstructorListCourseEnrollments)
		instructorRoutes.PATCH("/courses/:id/enrollments", api.InstructorUpdateEnrollmentStatus)
	}

		// ✅ Manager Route Group
	// Prefix: /manager
	managerRoutes := r.Group("/manager")
	{
		managerRoutes.GET("/users", api.ManagerListUsers)
		managerRoutes.GET("/users/:id/enrollments", api.ManagerListUserEnrollments)
		managerRoutes.POST("/users/:id/enrollments", api.ManagerAddUserEnrollment)
		managerRoutes.DELETE("/users/:id/enrollments/:course_id", api.ManagerDeleteUserEnrollment)
	}


	return r
}
