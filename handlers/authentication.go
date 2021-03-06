package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/BentoBoxSchools/web"
	"github.com/gorilla/sessions"
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

// HandleGoogleCallback takes token from request query and exchange for user information
func HandleGoogleCallback(store sessions.Store, clientID, clientSecret, redirectURI string, emailWhitelist []string) http.HandlerFunc {
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
		var userInfo web.User

		json.NewDecoder(res.Body).Decode(&userInfo)

		fmt.Printf(
			"User %s logged in. Checking user email is withing white list...\n",
			userInfo.Email,
		)
		for _, whitelistedEmail := range emailWhitelist {
			if userInfo.Email == whitelistedEmail {
				// success
				session, err := store.Get(r, "user")
				if err != nil {
					fmt.Println("failed to get session store to store authenticated user information.", err)
					break
				}
				session.Values["user"] = userInfo
				if e := session.Save(r, w); e != nil {
					fmt.Println("Failed to save session", e)
				}
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
func HandleLogout(store sessions.Store) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "user")
		delete(session.Values, "user")
		session.Options.MaxAge = -1
		_ = session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	})
}
