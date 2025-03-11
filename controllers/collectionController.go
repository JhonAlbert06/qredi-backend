package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"prestamosbackend/initializers"
	"prestamosbackend/models"
	"prestamosbackend/utils"
	"time"
)

type CollectionDetail struct {
	FeeId      uuid.UUID `json:"feeId"`
	PaidAmount float32   `json:"paidAmount"`
}

func CreateCollection(c *gin.Context) {

	collectionId := uuid.New()
	db := initializers.DB

	var body struct {
		RouteId          string  `form:"routeId"`
		Amount           float32 `form:"amount"`
		CollectionDetail string  `form:"collectionDetail"`
	}

	var route models.Route
	if err := db.Where("id = ?", body.RouteId).First(&route).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error al parsear el route",
		})
		return
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read body",
		})
		fmt.Println(err)
		return
	}

	//Get the user
	u, _ := c.Get("user")
	if u.(models.User).ID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user",
		})
		return
	}

	// Parse the user
	var user models.User
	if u, ok := u.(models.User); ok {
		user = u
	} else {
		fmt.Println(ok)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to load user",
		})
	}

	jsonDetailCollection := body.CollectionDetail
	var DetailCollectionList []CollectionDetail

	detailParseErro := json.Unmarshal([]byte(jsonDetailCollection), &DetailCollectionList)
	if detailParseErro != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error al parsear el detalle",
		})

		return
	}

	for _, item := range DetailCollectionList {

		fmt.Println(item.FeeId)

		var fee models.Fee
		if err := db.Where("id = ?", item.FeeId).First(&fee).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Error al parsear el fee",
			})
			fmt.Println(err)
			return
		}

		var loan models.Loan
		if err := db.First(&loan, fee.LoanId).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Error al parsear el loan",
			})
			return
		}

		var customer models.Customer
		if err := db.First(&customer, loan.CustomerId).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Error al parsear el customer",
			})
			return
		}

		interestAmount := float32(loan.Interest) / 100 * loan.Amount

		var payment = models.Payment{
			ID:         uuid.New(),
			UserId:     user.ID,
			FeeId:      fee.ID,
			PaidDate:   time.Now(),
			PaidAmount: item.PaidAmount,
		}

		if result := db.Create(&payment).Error; result != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "could not save payment",
			})
			return
		}

		if utils.MarkLoanAsPaidIfAllFeesPaid(db, fee.LoanId) {
			c.JSON(http.StatusOK, gin.H{
				"error": "El prestates esta pago",
			})
			return
		}

		if result := db.Create(&models.CollectionDetail{
			ID:               uuid.New(),
			CollectionId:     collectionId,
			LoanId:           fee.LoanId,
			FeeId:            fee.ID,
			FeeNumber:        fee.Number,
			FeeQuantity:      loan.FeesQuantity,
			ExpectedAmount:   interestAmount,
			PaidAmount:       item.PaidAmount,
			CustomerId:       loan.CustomerId,
			CustomerFullName: customer.Names + " " + customer.LastNames,
		}).Error; result != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Error al guardar el detalle",
			})
			return
		}
	}

	if result := db.Create(&models.Collection{
		ID:           collectionId,
		RouteId:      uuid.MustParse(body.RouteId),
		RouteName:    route.Name, // Ojo
		UserId:       user.ID,
		UserName:     user.UserName,
		UserFullName: user.FirstName + " " + user.LastName,
		Date:         time.Now(),
		Amount:       body.Amount,
	}).Error; result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error al guardar la collection",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}
