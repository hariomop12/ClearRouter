package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hariomop12/clearrouter/apps/backend/internal/models"
	"github.com/hariomop12/clearrouter/apps/backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type OAuthHandler struct {
	db *gorm.DB
}

func NewOAuthHandler(db *gorm.DB) *OAuthHandler {
	return &OAuthHandler{db: db}
}

func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (h *OAuthHandler) setStateCookie(c *gin.Context, state string) {
	c.SetCookie("oauth_state", state, 300, "/", "", false, true)
}

func (h *OAuthHandler) verifyState(c *gin.Context, state string) bool {
	cookie, err := c.Cookie("oauth_state")
	if err != nil || cookie == "" || cookie != state {
		return false
	}
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)
	return true
}

func frontendURL() string {
	if u := os.Getenv("FRONTEND_URL"); u != "" {
		return u
	}
	return "http://localhost:5173"
}

func (h *OAuthHandler) authRedirect(c *gin.Context, token string, user *models.User) {
	userJSON, _ := json.Marshal(gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
	loc := fmt.Sprintf("%s/oauth/callback?token=%s&user=%s",
		frontendURL(),
		url.QueryEscape(token),
		url.QueryEscape(string(userJSON)),
	)
	c.Redirect(http.StatusTemporaryRedirect, loc)
}

func (h *OAuthHandler) findOrCreateOAuthUser(email, name, providerID, provider string) (*models.User, error) {
	var user models.User
	email = strings.ToLower(strings.TrimSpace(email))

	err := h.db.Where("email = ?", email).First(&user).Error
	if err == nil {
		needsUpdate := false
		if provider == "google" && user.GoogleID == nil {
			needsUpdate = true
		} else if provider == "github" && user.GitHubID == nil {
			needsUpdate = true
		}
		if needsUpdate {
			if err := h.db.Model(&user).Update(provider+"_id", providerID).Error; err != nil {
				fmt.Printf("[OAUTH] Failed to link %s account: %v\n", provider, err)
				return nil, err
			}
		}
		return &user, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Printf("[OAUTH] DB error looking up user: %v\n", err)
		return nil, err
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(uuid.New().String()), bcrypt.DefaultCost)
	newUser := models.User{
		Name:          name,
		Email:         email,
		PasswordHash:  string(hash),
		EmailVerified: true,
	}
	if provider == "google" {
		newUser.GoogleID = &providerID
	} else {
		newUser.GitHubID = &providerID
	}

	if err := h.db.Create(&newUser).Error; err != nil {
		// If duplicate email, try to find and link instead
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			if findErr := h.db.Where("email = ?", email).First(&user).Error; findErr == nil {
				if err := h.db.Model(&user).Update(provider+"_id", providerID).Error; err == nil {
					return &user, nil
				}
			}
		}
		fmt.Printf("[OAUTH] Failed to create user via %s: email=%s err=%v\n", provider, email, err)
		return nil, err
	}
	fmt.Printf("[OAUTH] Created user via %s: email=%s id=%s\n", provider, email, newUser.ID)
	return &newUser, nil
}

// GoogleLogin redirects to Google OAuth
func (h *OAuthHandler) GoogleLogin(c *gin.Context) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	if clientID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Google OAuth not configured"})
		return
	}

	state, err := generateState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state"})
		return
	}
	h.setStateCookie(c, state)

	redirectURI := os.Getenv("GOOGLE_REDIRECT_URI")
	if redirectURI == "" {
		redirectURI = "http://localhost:8080/auth/google/callback"
	}

	u := fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=email%%20profile&state=%s",
		url.QueryEscape(clientID),
		url.QueryEscape(redirectURI),
		url.QueryEscape(state),
	)
	c.Redirect(http.StatusTemporaryRedirect, u)
}

// GoogleCallback handles the Google OAuth callback
func (h *OAuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")

	if errorParam != "" {
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=%s", frontendURL(), url.QueryEscape(errorParam)))
		return
	}

	if !h.verifyState(c, state) {
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=invalid_state", frontendURL()))
		return
	}

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURI := os.Getenv("GOOGLE_REDIRECT_URI")
	if redirectURI == "" {
		redirectURI = "http://localhost:8080/auth/google/callback"
	}

	tokenResp, err := http.PostForm("https://oauth2.googleapis.com/token", url.Values{
		"code":          {code},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"redirect_uri":  {redirectURI},
		"grant_type":    {"authorization_code"},
	})
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=token_exchange_failed", frontendURL()))
		return
	}
	defer tokenResp.Body.Close()

	body, _ := io.ReadAll(tokenResp.Body)
	if tokenResp.StatusCode != http.StatusOK {
		fmt.Printf("[OAUTH] Google token exchange failed: %s — %s\n", tokenResp.Status, string(body))
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=token_exchange_failed", frontendURL()))
		return
	}

	var tokenData struct {
		AccessToken string `json:"access_token"`
		IDToken     string `json:"id_token"`
	}
	if err := json.Unmarshal(body, &tokenData); err != nil {
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=invalid_token_response", frontendURL()))
		return
	}

	if tokenData.AccessToken == "" {
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=empty_access_token", frontendURL()))
		return
	}

	userInfoResp, err := http.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", url.QueryEscape(tokenData.AccessToken)))
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=userinfo_failed", frontendURL()))
		return
	}
	defer userInfoResp.Body.Close()

	body, err = io.ReadAll(userInfoResp.Body)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=userinfo_failed", frontendURL()))
		return
	}
	if userInfoResp.StatusCode != http.StatusOK {
		fmt.Printf("[OAUTH] Google userinfo failed: %s — %s\n", userInfoResp.Status, string(body))
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=userinfo_failed", frontendURL()))
		return
	}

	var googleUser struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.Unmarshal(body, &googleUser); err != nil {
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=invalid_userinfo", frontendURL()))
		return
	}

	if googleUser.Email == "" {
		fmt.Printf("[OAUTH] Google returned empty email for user %+v\n", googleUser)
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=no_email", frontendURL()))
		return
	}

	if googleUser.Name == "" {
		googleUser.Name = strings.Split(googleUser.Email, "@")[0]
	}

	user, err := h.findOrCreateOAuthUser(googleUser.Email, googleUser.Name, googleUser.ID, "google")
	if err != nil {
		fmt.Printf("[OAUTH] findOrCreateOAuthUser failed: %v\n", err)
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=auth_failed", frontendURL()))
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=token_generation_failed", frontendURL()))
		return
	}

	h.authRedirect(c, token, user)
}

// GitHubLogin redirects to GitHub OAuth
func (h *OAuthHandler) GitHubLogin(c *gin.Context) {
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	if clientID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "GitHub OAuth not configured"})
		return
	}

	state, err := generateState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state"})
		return
	}
	h.setStateCookie(c, state)

	redirectURI := os.Getenv("GITHUB_REDIRECT_URI")
	if redirectURI == "" {
		redirectURI = "http://localhost:8080/auth/github/callback"
	}

	u := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user:email&state=%s",
		url.QueryEscape(clientID),
		url.QueryEscape(redirectURI),
		url.QueryEscape(state),
	)
	c.Redirect(http.StatusTemporaryRedirect, u)
}

// GitHubCallback handles the GitHub OAuth callback
func (h *OAuthHandler) GitHubCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")

	if errorParam != "" {
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=%s", frontendURL(), url.QueryEscape(errorParam)))
		return
	}

	if !h.verifyState(c, state) {
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=invalid_state", frontendURL()))
		return
	}

	clientID := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	redirectURI := os.Getenv("GITHUB_REDIRECT_URI")
	if redirectURI == "" {
		redirectURI = "http://localhost:8080/auth/github/callback"
	}

	tokenResp, err := http.PostForm("https://github.com/login/oauth/access_token", url.Values{
		"code":          {code},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"redirect_uri":  {redirectURI},
	})
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=token_exchange_failed", frontendURL()))
		return
	}
	defer tokenResp.Body.Close()

	body, _ := io.ReadAll(tokenResp.Body)
	if tokenResp.StatusCode != http.StatusOK {
		fmt.Printf("[OAUTH] GitHub token exchange failed: %s — %s\n", tokenResp.Status, string(body))
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=token_exchange_failed", frontendURL()))
		return
	}

	vals, _ := url.ParseQuery(string(body))
	accessToken := vals.Get("access_token")

	if accessToken == "" {
		fmt.Printf("[OAUTH] GitHub empty access token — response: %s\n", string(body))
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=invalid_token", frontendURL()))
		return
	}

	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")
	userInfoResp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=userinfo_failed", frontendURL()))
		return
	}
	defer userInfoResp.Body.Close()

	body, err = io.ReadAll(userInfoResp.Body)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=userinfo_failed", frontendURL()))
		return
	}
	if userInfoResp.StatusCode != http.StatusOK {
		fmt.Printf("[OAUTH] GitHub userinfo failed: %s — %s\n", userInfoResp.Status, string(body))
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=userinfo_failed", frontendURL()))
		return
	}

	var githubUser struct {
		ID    int    `json:"id"`
		Login string `json:"login"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.Unmarshal(body, &githubUser); err != nil {
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=invalid_userinfo", frontendURL()))
		return
	}

	if githubUser.Email == "" {
		emailReq, _ := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
		emailReq.Header.Set("Authorization", "Bearer "+accessToken)
		emailReq.Header.Set("Accept", "application/json")
		emailResp, err := http.DefaultClient.Do(emailReq)
		if err == nil {
			defer emailResp.Body.Close()
			emailBody, _ := io.ReadAll(emailResp.Body)
			if emailResp.StatusCode == http.StatusOK {
				var emails []struct {
					Email   string `json:"email"`
					Primary bool   `json:"primary"`
					Verified bool  `json:"verified"`
				}
				if err := json.Unmarshal(emailBody, &emails); err == nil {
					for _, e := range emails {
						if e.Primary && e.Verified {
							githubUser.Email = e.Email
							break
						}
					}
					if githubUser.Email == "" && len(emails) > 0 {
						githubUser.Email = emails[0].Email
					}
				}
			}
		}
	}

	if githubUser.Email == "" {
		fmt.Printf("[OAUTH] GitHub returned no email for user %+v\n", githubUser)
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=no_email", frontendURL()))
		return
	}

	if githubUser.Name == "" {
		githubUser.Name = strings.Split(githubUser.Email, "@")[0]
	}

	user, err := h.findOrCreateOAuthUser(githubUser.Email, githubUser.Name, fmt.Sprintf("%d", githubUser.ID), "github")
	if err != nil {
		fmt.Printf("[OAUTH] findOrCreateOAuthUser failed: %v\n", err)
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=auth_failed", frontendURL()))
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/login?error=token_generation_failed", frontendURL()))
		return
	}

	h.authRedirect(c, token, user)
}

// OAuthStatus returns which providers are configured
func (h *OAuthHandler) OAuthStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"google": os.Getenv("GOOGLE_CLIENT_ID") != "",
		"github": os.Getenv("GITHUB_CLIENT_ID") != "",
	})
}
