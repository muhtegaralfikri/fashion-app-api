package routes


import (
    "appFashion/backend/controllers"
    "net/http"


    "github.com/gin-gonic/gin"
)


func SetupRoutes(r *gin.Engine) {
    // Grup rute untuk API v1
    apiV1 := r.Group("/api/v1")
    {
        // Endpoint sederhana untuk tes koneksi
        apiV1.GET("/ping", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{"message": "pong"})
        })


        // Grup untuk rute otentikasi
        auth := apiV1.Group("/auth")
        {
            auth.POST("/register", controllers.Register)
            // Tambahkan rute login dan forgot-password di sini nanti
            auth.POST("/login", controllers.Login)
            auth.POST("/forgot-password", controllers.ForgotPassword)
            auth.POST("/reset-password", controllers.ResetPassword)
        }


        // Grup untuk rute produk
        products := apiV1.Group("/products")
        {
        products.GET("/", controllers.GetProducts)
        products.GET("/:id", controllers.GetProductByID)
        }
    }
}

