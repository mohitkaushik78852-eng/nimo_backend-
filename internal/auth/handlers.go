package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"nimo-backend/internal/models"
)

type Handler struct {
	DB             *gorm.DB
	JWTSecret      string
	GoogleClientID string
}

func NewHandler(db *gorm.DB, jwtSecret string, googleClientID string) *Handler {
	return &Handler{DB: db, JWTSecret: jwtSecret, GoogleClientID: googleClientID}
}

type registerRequest struct {
	Name     string `json:"name" binding:"required,max=100"`
	Email    string `json:"email" binding:"required,email,max=254"`
	Password string `json:"password" binding:"required,min=6,max=128"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email,max=254"`
	Password string `json:"password" binding:"required,max=128"`
}

type googleRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}

func (h *Handler) authResponse(c *gin.Context, user models.User) {
	token, err := GenerateToken(user.ID, user.Email, h.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":            user.ID,
			"name":          user.Name,
			"email":         user.Email,
			"auth_provider": user.AuthProvider,
		},
	})
}

func (h *Handler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	var existing models.User
	if err := h.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "an account with this email already exists"})
		return
	}

	hash, err := HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not process password"})
		return
	}

	user := models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hash,
		AuthProvider: "password",
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create account"})
		return
	}

	h.authResponse(c, user)
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	var user models.User
	if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	if user.AuthProvider != "password" || user.PasswordHash == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "this account uses Google sign-in, log in with Google instead"})
		return
	}

	if !CheckPassword(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	h.authResponse(c, user)
}

func (h *Handler) GoogleLogin(c *gin.Context) {
	var req googleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	info, err := VerifyGoogleIDToken(req.IDToken, h.GoogleClientID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err = h.DB.Where("email = ?", info.Email).First(&user).Error
	if err != nil {
		user = models.User{
			Name:         info.Name,
			Email:        info.Email,
			GoogleID:     info.Sub,
			AuthProvider: "google",
		}
		if err := h.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create account"})
			return
		}
	} else if user.GoogleID == "" {
		user.GoogleID = info.Sub
		h.DB.Save(&user)
	}

	h.authResponse(c, user)
}

func (h *Handler) Me(c *gin.Context) {
	userID := c.GetUint("user_id")
	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":            user.ID,
		"name":          user.Name,
		"email":         user.Email,
		"auth_provider": user.AuthProvider,
	})
}
