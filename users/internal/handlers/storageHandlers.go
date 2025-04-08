package handlers

import (
	"github.com/RicliZz/shop/users/internal/services"
	"github.com/gin-gonic/gin"
)

type Storage struct {
	StorageService services.StorageService
}

func NewStorageHandler(storage services.StorageService) *Storage {
	return &Storage{StorageService: storage}
}

func (h *Storage) InitRoutes(router *gin.RouterGroup) {
	storageRouter := router.Group("storage/")
	{
		storageRouter.GET("get-all/", h.StorageService.GetAllProducts)
	}
}
