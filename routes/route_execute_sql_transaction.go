package routes

import (
	"database/sql"
	"go-api-db/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RouteExecuteSQLTransaction(router *gin.Engine, db *sql.DB) {

	router.POST("/execute-sql-transaction", authorize(config.ConfigVar.Server.Authorization), func(c *gin.Context) {

		var requestBody config.RequestBodySQL
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao iniciar a transação"})
			return
		}

		_, err = tx.Exec(requestBody.SQL)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = tx.Commit()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao confirmar a transação"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "success"})

	})

}