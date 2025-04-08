package services

import "github.com/gin-gonic/gin"

type StorageService interface {
	GetAllProducts(c *gin.Context)
}
