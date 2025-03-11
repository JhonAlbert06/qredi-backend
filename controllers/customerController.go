package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"prestamosbackend/initializers"
	"prestamosbackend/models"
	"prestamosbackend/responses"
	"prestamosbackend/utils"
)

func CreateCustomer(c *gin.Context) {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get program directory",
		})
		return
	}

	var imageDir string
	if os.Getenv("GO_ENV") == "production" {
		imageDir = filepath.Join(dir, "files/customer/images")
	} else {
		imageDir = filepath.Join("files/customer/images")
	}

	// imageDir := filepath.Join("files/customer/images")
	if err := os.MkdirAll(imageDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create image directory",
		})
		return
	}

	// Create a new customer
	customerID := uuid.New()

	//Get the body
	var body struct {
		CompanyId   string  `form:"companyId" json:"companyId"`
		Cedula      string  `form:"cedula" json:"cedula"`
		Names       string  `form:"names" json:"names"`
		LastNames   string  `form:"lastNames" json:"lastNames"`
		Address     string  `form:"address" json:"address"`
		Phone       string  `form:"phone" json:"phone"`
		CivilStatus string  `form:"civilStatus" json:"civilStatus"`
		Reference   string  `form:"reference" json:"reference"`
		Latitude    *string `form:"latitude" json:"latitude"`
		Longitude   *string `form:"longitude" json:"longitude"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read body",
		})
		fmt.Println(err)
		return
	}

	if body.Cedula == "" || body.Names == "" || body.LastNames == "" || body.Address == "" || body.Phone == "" || body.CivilStatus == "" || body.Reference == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Empty fields",
		})
		return
	}

	// Process and save image with the recipe ID as the filename
	imageFile, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read image file",
		})
		return
	}
	defer func(imageFile multipart.File) {
		err := imageFile.Close()
		if err != nil {

		}
	}(imageFile)

	companyId := uuid.MustParse(body.CompanyId)

	// Create Customer with file paths
	customer := models.Customer{
		CompanyId:   companyId,
		ID:          customerID,
		Cedula:      body.Cedula,
		Names:       body.Names,
		LastNames:   body.LastNames,
		Address:     body.Address,
		Phone:       body.Phone,
		CivilStatus: body.CivilStatus,
		Reference:   body.Reference,
		Latitude:    body.Latitude,
		Longitude:   body.Longitude,
	}

	result := initializers.DB.Create(&customer)

	if result.Error != nil {
		fmt.Println(result.Error)
		initializers.DB.Unscoped().Delete(&customer)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create customer",
		})
		return
	}

	imageFileName := customerID.String() + ".jpg"
	imagePath := filepath.Join(imageDir, imageFileName)
	imageDestination, err := os.Create(imagePath)
	if err != nil {
		initializers.DB.Unscoped().Delete(&customer)
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": err,
		})
		return
	}
	defer func(imageDestination *os.File) {
		err := imageDestination.Close()
		if err != nil {

		}
	}(imageDestination)

	_, err = io.Copy(imageDestination, imageFile)
	if err != nil {
		initializers.DB.Unscoped().Delete(&customer)
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Failed to save image",
		})
		return
	}

	// Location images
	if body.Latitude != nil && body.Longitude != nil {

		apiKey := os.Getenv("MAP_API_KEY")
		if apiKey == "" {
			initializers.DB.Unscoped().Delete(&customer)
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": "Google Maps API key not set",
			})
			return
		}

		mapURL := fmt.Sprintf("https://maps.googleapis.com/maps/api/staticmap?center=%s,%s&zoom=15&size=512x256&key=%s", *body.Latitude, *body.Longitude, apiKey)
		mapImagePath := filepath.Join(imageDir, customerID.String()+"_map.jpg")

		resp, err := http.Get(mapURL)
		if err != nil {
			initializers.DB.Unscoped().Delete(&customer)
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": "Failed to download map image",
			})
			return
		}
		defer resp.Body.Close()

		mapImageFile, err := os.Create(mapImagePath)
		if err != nil {
			initializers.DB.Unscoped().Delete(&customer)
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": "Failed to create map image file",
			})
			return
		}
		defer mapImageFile.Close()

		_, err = io.Copy(mapImageFile, resp.Body)
		if err != nil {
			initializers.DB.Unscoped().Delete(&customer)
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": "Failed to save map image",
			})
			return
		}
	}
	customerRes := responses.NewCustomerResponse(customer)

	// Respond
	c.JSON(http.StatusCreated, customerRes)
}

func UpdateCustomer(c *gin.Context) {

	var body struct {
		Id          string  `form:"id" json:"id"`
		Cedula      string  `form:"cedula" json:"cedula"`
		Names       string  `form:"names" json:"names"`
		LastNames   string  `form:"lastNames" json:"lastNames"`
		Address     string  `form:"address" json:"address"`
		Phone       string  `form:"phone" json:"phone"`
		CivilStatus string  `form:"civilStatusId" json:"civilStatus"`
		Reference   string  `form:"reference" json:"reference"`
		Latitude    *string `form:"latitude" json:"latitude"`
		Longitude   *string `form:"longitude" json:"longitude"`
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read body",
		})
		fmt.Println(err)
		return
	}

	var existingCustomer models.Customer
	if err := initializers.DB.Where("id = ?", body.Id).First(&existingCustomer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Cliente no encontrado",
		})
		return
	}

	if body.Cedula != "" {
		existingCustomer.Cedula = body.Cedula
	}

	if body.Names != "" {
		existingCustomer.Names = body.Names
	}

	if body.LastNames != "" {
		existingCustomer.LastNames = body.LastNames
	}

	if body.Address != "" {
		existingCustomer.Address = body.Address
	}

	if body.Phone != "" {
		existingCustomer.Phone = body.Phone
	}

	if body.CivilStatus != "" {
		existingCustomer.CivilStatus = body.CivilStatus
	}

	if body.Reference != "" {
		existingCustomer.Reference = body.Reference
	}

	if body.Latitude != nil {
		existingCustomer.Latitude = body.Latitude
	}

	if body.Longitude != nil {
		existingCustomer.Longitude = body.Longitude
	}

	// Revisa si la imagen fue enviada en la solicitud
	imageFile, _, err := c.Request.FormFile("image")
	var imagePath string
	if err != nil {
		// Si no se envió una imagen, simplemente deja imagePath como una cadena vacía
		imagePath = ""
	} else {
		defer func(imageFile multipart.File) {
			err := imageFile.Close()
			if err != nil {
				// Manejar el error al cerrar el archivo si es necesario
			}
		}(imageFile)

		// Si se envió una imagen, guarda la imagen
		imageFileName := existingCustomer.ID.String() + ".jpg"
		imagePath = filepath.Join("files/customer/images", imageFileName)
		imageDestination, err := os.Create(imagePath)
		if err != nil {
			initializers.DB.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": "Failed to save image",
			})
			return
		}
		defer func(imageDestination *os.File) {
			err := imageDestination.Close()
			if err != nil {
				// Manejar el error al cerrar el archivo si es necesario
			}
		}(imageDestination)

		// Si se proporcionó una imagen, copiarla al destino
		_, err = io.Copy(imageDestination, imageFile)
		if err != nil {
			initializers.DB.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": "Failed to save image",
			})
			return
		}
	}

	result := initializers.DB.Save(&existingCustomer)
	if result.Error != nil {
		fmt.Println(result.Error)
		initializers.DB.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update customer",
		})
		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{})
}

func SearchCustomerByParameter(c *gin.Context) {
	query := c.Query("query")
	searchField := c.Query("field")

	db := initializers.DB

	user := c.MustGet("user").(models.User)
	db = db.Where("company_id = ?", user.CompanyId)

	var customers []models.Customer

	switch searchField {
	case "names":
		db = db.Where("names LIKE ?", "%"+query+"%").Limit(15)
	case "last_names":
		db = db.Where("last_names LIKE ?", "%"+query+"%").Limit(15)
	case "cedula":
		query = utils.EliminarGuiones(query)
		query = utils.FormatearCedula(query)
		db = db.Where("cedula LIKE ?", "%"+query+"%").Limit(15)
	case "phone":
		query = utils.EliminarGuiones(query)
		query = utils.FormatearTelefono(query)
		db = db.Where("phone LIKE ?", "%"+query+"%").Limit(15)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Campo de búsqueda no válido",
		})
		return
	}

	if err := db.Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error to load customers",
		})
		return
	}

	var customerResponse []responses.CustomerResponse

	for _, customer := range customers {
		customerRes := responses.NewCustomerResponse1(customer)
		customerResponse = append(customerResponse, customerRes)
	}

	c.JSON(http.StatusOK, customerResponse)
}

func SearchCustomerById(c *gin.Context) {
	id := c.Param("id")

	db := initializers.DB

	user := c.MustGet("user").(models.User)
	db = db.Where("company_id = ?", user.CompanyId)

	var customer models.Customer
	if err := db.Where("id = ?", id).Find(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load customer",
		})
		return
	}

	customerResponse := responses.NewCustomerResponse(customer)

	c.JSON(http.StatusOK, customerResponse)
}

func GetCustomerImage(c *gin.Context) {

	id := c.Param("id")

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get program directory",
		})
		return
	}

	var imageDir string
	if os.Getenv("GO_ENV") == "production" {
		imageDir = filepath.Join(dir, "files/customer/images")
	} else {
		imageDir = filepath.Join("files/customer/images")
	}

	imagePath := filepath.Join(imageDir, id+".jpg")
	if _, err := os.Stat(imagePath); err == nil {
		c.File(imagePath)
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Image not found",
		})
	}
}

func UpdateCustomerImage(c *gin.Context) {

	id := c.Param("id")

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get program directory",
		})
		return
	}

	var imageDir string
	if os.Getenv("GO_ENV") == "production" {
		imageDir = filepath.Join(dir, "files/customer/images")
	} else {
		imageDir = filepath.Join("files/customer/images")
	}

	// Crear las carpetas si no existen
	if err := os.MkdirAll(imageDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create image directory",
		})
		return
	}

	// Verificar si el cliente existe en la base de datos
	var existingCustomer models.Customer
	if err := initializers.DB.Where("id = ?", id).First(&existingCustomer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User no encontrado",
		})
		return
	}

	// Procesar y guardar la nueva imagen
	newImageFile, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fallo al leer el archivo de imagen",
		})
		return
	}
	defer func(newImageFile multipart.File) {
		err := newImageFile.Close()
		if err != nil {
			// Manejar el error de cierre si es necesario
		}
	}(newImageFile)

	// Eliminar la imagen anterior si existe
	oldImageFilePath := filepath.Join(imageDir, existingCustomer.ID.String()+".jpg")
	if err := os.Remove(oldImageFilePath); err != nil && !os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Fallo al eliminar la imagen anterior",
		})
		return
	}

	// Guardar la nueva imagen
	newImageFileName := existingCustomer.ID.String() + ".jpg"
	newImagePath := filepath.Join(imageDir, newImageFileName)
	newImageDestination, err := os.Create(newImagePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Fallo al guardar la nueva imagen",
		})
		return
	}
	defer func(newImageDestination *os.File) {
		err := newImageDestination.Close()
		if err != nil {
			// Manejar el error de cierre si es necesario
		}
	}(newImageDestination)

	_, err = io.Copy(newImageDestination, newImageFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Fallo al guardar la nueva imagen",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Imagen del cliente actualizada con éxito",
	})
}

func GetCustomerImageMap(c *gin.Context) {
	id := c.Param("id")

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get program directory",
		})
		return
	}

	var imageDir string
	if os.Getenv("GO_ENV") == "production" {
		imageDir = filepath.Join(dir, "files/customer/images")
	} else {
		imageDir = filepath.Join("files/customer/images")
	}

	var imagePath string
	imagePath = filepath.Join(imageDir, id+"_map.jpg")

	if _, err := os.Stat(imagePath); err == nil {
		c.File(imagePath)
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Image not found",
		})
	}
}
