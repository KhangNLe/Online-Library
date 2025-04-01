package mybook

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func getFavoriteBook(c *gin.Context, query *sqlx.Tx, libId int) {

}
