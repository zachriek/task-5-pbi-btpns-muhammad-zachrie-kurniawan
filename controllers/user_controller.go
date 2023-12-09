package controllers

import (
	"net/http"
	"strconv"
	"task-5-pbi-btpns-muhammad-zachrie-kurniawan/app"
	"task-5-pbi-btpns-muhammad-zachrie-kurniawan/database"
	"task-5-pbi-btpns-muhammad-zachrie-kurniawan/helpers"
	"task-5-pbi-btpns-muhammad-zachrie-kurniawan/models"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

func Register(context *gin.Context) {
	var userFormRegister app.UserFormRegister
	if err := context.ShouldBindJSON(&userFormRegister); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if _, err := govalidator.ValidateStruct(userFormRegister); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	if len(userFormRegister.Password) < 6 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Password minimal 6 karakter!"})
		context.Abort()
		return
	}

	if err := database.Instance.Where("email = ?", userFormRegister.Email).First(&user).Error; err == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah terdaftar!"})
		context.Abort()
		return
	}

	if err := database.Instance.Where("username = ?", userFormRegister.Username).First(&user).Error; err == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Username sudah terdaftar!"})
		context.Abort()
		return
	}

	user = models.User{
		Username: userFormRegister.Username,
		Email:    userFormRegister.Email,
		Password: userFormRegister.Password,
	}

	if err := user.HashPassword(user.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	record := database.Instance.Create(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Berhasil membuat akun!"})
}

func Login(context *gin.Context) {
	var userFormLogin app.UserFormLogin
	if err := context.ShouldBindJSON(&userFormLogin); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := govalidator.ValidateStruct(userFormLogin); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.Instance.Where("email = ?", userFormLogin.Email).First(&user).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Email tidak sesuai!"})
		return
	}

	if err := user.CheckPassword(userFormLogin.Password); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Password tidak sesuai!"})
		return
	}

	token, err := helpers.GenerateJWT(user.ID, user.Email, user.Username)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Berhasil masuk ke akun!", "token": token})
}

func GetUserByID(context *gin.Context) {
	userID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "ID pengguna tidak valid!"})
		return
	}

	tokenString := context.GetHeader("Authorization")
	claims, err := helpers.ParseToken(tokenString)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	if userID != int(claims.ID) {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak diizinkan!"})
		context.Abort()
		return
	}

	var user models.User
	if err := database.Instance.Where("id = ?", claims.ID).First(&user).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Pengguna tidak ditemukan!"})
		return
	}

	if err := database.Instance.First(&user, userID).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan!"})
		return
	}

	var userResult app.UserResult
	userResult.ID = user.ID
	userResult.Username = user.Username
	userResult.Email = user.Email
	userResult.CreatedAt = user.CreatedAt.String()
	userResult.UpdatedAt = user.UpdatedAt.String()

	context.JSON(http.StatusOK, gin.H{"data": userResult})
}

func UpdateUserByID(context *gin.Context) {
	userID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "ID pengguna tidak valid!"})
		return
	}
	var userFormUpdate app.UserFormUpdate
	if err := context.ShouldBindJSON(&userFormUpdate); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := govalidator.ValidateStruct(userFormUpdate); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	if len(userFormUpdate.Password) < 6 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Password minimal 6 karakter!"})
		context.Abort()
		return
	}

	if err := database.Instance.Where("email = ? AND id != ?", userFormUpdate.Email, userID).First(&user).Error; err == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah terdaftar!"})
		context.Abort()
		return
	}

	if err := database.Instance.Where("username = ? AND id != ?", userFormUpdate.Username, userID).First(&user).Error; err == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Username sudah terdaftar!"})
		context.Abort()
		return
	}

	tokenString := context.GetHeader("Authorization")
	claims, err := helpers.ParseToken(tokenString)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if userID != int(claims.ID) {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak diizinkan!"})
		context.Abort()
		return
	}

	if err := database.Instance.Where("id = ?", claims.ID).First(&user).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Pengguna tidak ditemukan!"})
		return
	}

	if err := database.Instance.First(&user, userID).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan!"})
		return
	}

	user.Username = userFormUpdate.Username
	user.Email = userFormUpdate.Email
	if userFormUpdate.Password != "" {
		if err := user.HashPassword(userFormUpdate.Password); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			context.Abort()
			return
		}
	}

	if err := database.Instance.Save(&user).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Berhasil memperbarui data pengguna!"})
}

func DeleteUserByID(context *gin.Context) {
	var user models.User
	userID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "ID pengguna tidak valid!"})
		return
	}

	tokenString := context.GetHeader("Authorization")
	claims, err := helpers.ParseToken(tokenString)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if userID != int(claims.ID) {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Tidak diizinkan!"})
		context.Abort()
		return
	}

	if err := database.Instance.Where("id = ?", claims.ID).First(&user).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Pengguna tidak ditemukan!"})
		return
	}

	if err := database.Instance.First(&user, userID).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan!"})
		return
	}

	if err := database.Instance.Delete(&user).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Berhasil menghapus akun pengguna!"})
}
