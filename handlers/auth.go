package handlers

import (
	"html/template"
	"log"

	"github.com/gin-gonic/gin"

	"sprinto/models"
	"sprinto/service"
)

type AuthHandler struct {
	svc     service.AuthService
	teamSvc service.TeamMemberService
}

func NewAuthHandler(svc service.AuthService, teamSvc service.TeamMemberService) *AuthHandler {
	return &AuthHandler{svc: svc, teamSvc: teamSvc}
}

// renderAuth renders a standalone page (no layout).
func renderAuth(c *gin.Context, page string, data interface{}) {
	t, err := template.New("").Funcs(funcMap).ParseFiles("templates/" + page + ".html")
	if err != nil {
		c.String(500, "Template error: %s", err.Error())
		return
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Status(200)
	if err := t.ExecuteTemplate(c.Writer, page, data); err != nil {
		log.Printf("renderAuth %s: %v", page, err)
	}
}

func (h *AuthHandler) LoginPage(c *gin.Context) {
	renderAuth(c, "login", gin.H{})
}

func (h *AuthHandler) RegisterPage(c *gin.Context) {
	renderAuth(c, "register", gin.H{})
}

func (h *AuthHandler) Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	user, err := h.svc.Login(email, password)
	if err != nil {
		renderAuth(c, "login", gin.H{"Error": "Invalid email or password", "Email": email})
		return
	}
	setSession(c, user.ID)
	redirectTo(c, "/")
}

func (h *AuthHandler) Register(c *gin.Context) {
	fullName := c.PostForm("full_name")
	email := c.PostForm("email")
	password := c.PostForm("password")
	if err := h.svc.Register(fullName, email, password); err != nil {
		renderAuth(c, "register", gin.H{"Error": err.Error(), "FullName": fullName, "Email": email})
		return
	}
	user, _ := h.svc.Login(email, password)
	// Auto-add registered user as a team member
	h.teamSvc.CreateForUser(fullName, "", email, user.ID)
	setSession(c, user.ID)
	redirectTo(c, "/")
}

func (h *AuthHandler) Logout(c *gin.Context) {
	clearSession(c)
	redirectTo(c, "/login")
}

// LoadUserMiddleware reads the session cookie and injects *models.User into context.
// Soft — does not redirect if no session.
func LoadUserMiddleware(svc service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if id, ok := getSessionUserID(c); ok {
			if user, err := svc.UserByID(id); err == nil {
				c.Set("current_user", &user)
			}
		}
		c.Next()
	}
}

// AuthRequiredMiddleware redirects to /login if no user is loaded in context.
func AuthRequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exists := c.Get("current_user"); !exists {
			c.Redirect(302, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}

func currentUserFromCtx(c *gin.Context) *models.User {
	if u, exists := c.Get("current_user"); exists {
		if user, ok := u.(*models.User); ok {
			return user
		}
	}
	return nil
}
