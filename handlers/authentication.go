package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// RedirectGoogleLogin redirects user to google login page with redirect information
func RedirectGoogleLogin(clientID, redirectURI string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scope := "profile email"
		URL := fmt.Sprintf(
			"https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&scope=%s&response_type=code",
			clientID,
			redirectURI,
			scope,
		)
		http.Redirect(w, r, URL, http.StatusTemporaryRedirect)
	})
}

type googleExchangeToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}
type googleUserInfoDTO struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	Email   string `json:"email"`
}

// HandleGoogleCallback takes token from request query and exchange for user information
func HandleGoogleCallback(clientID, clientSecret, redirectURI string, emailWhitelist []string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(
				w,
				"Failed to get exchange code from Google redirect.",
				http.StatusBadRequest,
			)
			return
		}
		URL := "https://www.googleapis.com/oauth2/v4/token"
		resp, err := http.PostForm(
			URL,
			url.Values{
				"code":          {code},
				"client_id":     {clientID},
				"client_secret": {clientSecret},
				"redirect_uri":  {redirectURI},
				"grant_type":    {"authorization_code"},
			},
		)
		if err != nil {
			http.Error(
				w,
				"Failed to exchange code for access_token from Google",
				http.StatusBadRequest,
			)
			return
		}
		defer resp.Body.Close()
		var target googleExchangeToken

		json.NewDecoder(resp.Body).Decode(&target)

		// continue using access_token to retrieve user email information
		userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo"
		client := &http.Client{}
		req, err := http.NewRequest("GET", userInfoURL, nil)
		if err != nil {
			http.Error(
				w,
				"Failed to get user information from Google API",
				http.StatusBadRequest,
			)
			fmt.Println("Failed to create request to send to Google for grabbing user information", err)
			return
		}
		req.Header.Set("Authorization", fmt.Sprintf(
			"Bearer %s",
			target.AccessToken,
		))
		res, err := client.Do(req)
		if err != nil {
			http.Error(
				w,
				"Failed to get user information from Google API",
				http.StatusBadRequest,
			)
			fmt.Println("Failed to send request to Google for grabbing user information", err)
			return
		}
		defer res.Body.Close()
		var userInfo googleUserInfoDTO

		json.NewDecoder(res.Body).Decode(&userInfo)

		fmt.Printf(
			"User %s logged in. Checking user email is withing white list...",
			userInfo.Email,
		)
		for _, whitelistedEmail := range emailWhitelist {
			if userInfo.Email == whitelistedEmail {
				// success
				expiration := time.Now().Add(24 * time.Hour)
				userJSON, _ := json.Marshal(userInfo)
				cookie := http.Cookie{Name: "user", Value: string(userJSON), Expires: expiration}
				http.SetCookie(w, &cookie)
				http.Redirect(
					w, r,
					"/",
					http.StatusTemporaryRedirect,
				)
			}
		}
		http.Error(
			w,
			fmt.Sprintf("User %s is not authorized to login.", userInfo.Email),
			http.StatusUnauthorized,
		)
	})
}

// HandleLogout invalidates the login session
func HandleLogout() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: handle login and store user info into session?
	})
}
