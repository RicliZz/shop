package serviceStorage

import (
	"github.com/gin-gonic/gin"
)

type GetFullItemsInStorage struct {
}

func NewGetFullItemsInStorage() *GetFullItemsInStorage {
	return &GetFullItemsInStorage{}
}

func (s *GetFullItemsInStorage) GetAllProducts(c *gin.Context) {

}
