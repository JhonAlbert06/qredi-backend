package controllers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"prestamosbackend/initializers"
	"prestamosbackend/models"
	"prestamosbackend/responses"
	"prestamosbackend/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {

	//Get the phone and password
	var body struct {
		CompanyID string `json:"companyId" form:"companyId"`
		FirstName string `json:"firstName" form:"firstName"`
		LastName  string `json:"lastName" form:"lastName"`
		Username  string `json:"userName" form:"userName"`
		Password  string `json:"password" form:"password"`
		IsAdmin   bool   `json:"isAdmin" form:"isAdmin"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// IsStringEmpty checks if a string is null or empty.
	if utils.IsStringEmpty(body.FirstName) || utils.IsStringEmpty(body.LastName) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "FirstName and LastName cannot be empty",
		})

		return
	}

	// IsStringEmpty checks if a string is null or empty.
	if utils.IsStringEmpty(body.Username) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "The username cannot be empty",
		})

		return
	}

	// Validate password length
	if !utils.IsPasswordValid(body.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password must be between 1 and 40 characters",
		})

		return
	}

	//Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})

		return
	}

	var role models.Role
	initializers.DB.First(&role, "name = ?", "User")

	//Create User
	user := models.User{
		ID:              uuid.New(),
		CompanyId:       uuid.MustParse(body.CompanyID),
		RoleId:          role.ID,
		IsNew:           true,
		FirstName:       body.FirstName,
		LastName:        body.LastName,
		UserName:        body.Username,
		Password:        string(hash),
		PasswordVersion: uuid.New(),
	}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	//Respond
	c.JSON(http.StatusCreated, gin.H{})
}

func Login(c *gin.Context) {

	var body struct {
		Username string `json:"userName" form:"userName"`
		Password string `json:"password" form:"password"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	// Validate password length
	if !utils.IsPasswordValid(body.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password must be between 1 and 40 characters",
		})

		return
	}

	// look up requested user
	var user models.User
	initializers.DB.First(&user, "user_name = ?", body.Username)

	if user.ID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid User Name",
		})

		return
	}

	//Compare sent in pass with saved user pass hash

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid phone number or password",
		})

		return
	}

	//Generate jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":         user.ID,
		"exp":         time.Now().Add(time.Hour * 24 * 30).Unix(),
		"pwd_version": user.PasswordVersion,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})

		return
	}

	//Sent it back
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

func ChangePassword(c *gin.Context) {

	var body struct {
		CurrentPassword string `json:"currentPassword" form:"currentPassword"`
		NewPassword     string `json:"newPassword" form:"newPassword"`
	}

	if c.Bind(&body) != nil {
		fmt.Println(body)
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read body",
		})
		return
	}

	// Validate password length
	if !utils.IsPasswordValid(body.CurrentPassword) {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "The Current Password must be between 1 and 40 characters",
		})

		return
	}

	// Validate password length
	if !utils.IsPasswordValid(body.NewPassword) {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "The New Password must be between 1 and 40 characters",
		})

		return
	}

	user, _ := c.Get("user")
	if user.(models.User).ID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user",
		})
		return
	}

	// Hash the new password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Compare the current password with the stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.(models.User).Password), []byte(body.CurrentPassword)); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Incorrect current password",
		})
		return
	}

	// Change the password and save the user in the database
	if u, ok := user.(models.User); ok {
		u.Password = string(hash)
		u.PasswordVersion = uuid.New()
		result := initializers.DB.Save(&u)

		if result.Error != nil {
			fmt.Println(result.Error)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to change the password",
			})
			return
		}

		// Respond with a success message
		c.JSON(http.StatusOK, gin.H{})
	} else {
		fmt.Println(ok)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to load user",
		})
	}
}

func LoadUser(c *gin.Context) {

	user, _ := c.Get("user")

	if u, ok := user.(models.User); ok {

		userResp := responses.NewUserResponse(u)

		c.JSON(http.StatusOK, userResp)
	} else {
		fmt.Println(ok)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to load user",
		})
	}
}

func UpdateUserImage(c *gin.Context) {

	user, _ := c.Get("user")
	userID := user.(models.User).ID

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get program directory",
		})
		return
	}

	var imageDir string
	if os.Getenv("GO_ENV") == "production" {
		imageDir = filepath.Join(dir, "files/user/images")
	} else {
		imageDir = filepath.Join("files/user/images")
	}

	// imageDir := filepath.Join("files/customer/images")
	// Crear las carpetas si no existen
	if err := os.MkdirAll(imageDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create image directory",
		})
		return
	}

	// Verificar si el cliente existe en la base de datos
	var existingUser models.User
	if err := initializers.DB.Where("id = ?", userID).First(&existingUser).Error; err != nil {
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
	oldImageFilePath := filepath.Join(imageDir, userID.String()+".jpg")
	if err := os.Remove(oldImageFilePath); err != nil && !os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Fallo al eliminar la imagen anterior",
		})
		return
	}

	// Guardar la nueva imagen
	newImageFileName := userID.String() + ".jpg"
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

func GetUserImage(c *gin.Context) {

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
		imageDir = filepath.Join(dir, "files/user/images")
	} else {
		imageDir = filepath.Join("files/user/images")
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
