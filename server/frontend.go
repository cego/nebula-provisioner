package server

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	mrand "math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/slyngdk/nebula-provisioner/server/graph/model"
	"github.com/slyngdk/nebula-provisioner/webapp"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/slyngdk/nebula-provisioner/server/graph"
	"github.com/slyngdk/nebula-provisioner/server/graph/generated"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"github.com/slackhq/nebula"
	"github.com/slyngdk/nebula-provisioner/server/store"
	"golang.org/x/oauth2"
)

const CookieName = "session"

type frontend struct {
	sessions     sessions.Store
	oauth2Config oauth2.Config
	store        *store.Store
	ipManager    *store.IPManager

	config *nebula.Config
	l      *logrus.Logger
}

type userInfo struct {
	sub   string
	name  string
	email string
}

func NewFrontend(config *nebula.Config, logger *logrus.Logger, store *store.Store, ipManager *store.IPManager) (*frontend, error) {

	// Hash keys should be at least 32 bytes long
	authenticationKey := make([]byte, 64)
	_, err := rand.Read(authenticationKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate authenticationKey: %s", err)
	}
	encryptionKey := make([]byte, 32)
	_, err = rand.Read(encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate encryptionKey: %s", err)
	}

	cookieStore := sessions.NewCookieStore(authenticationKey, encryptionKey)
	cookieStore.MaxAge(int((10 * time.Minute).Seconds()))

	u, err := url.Parse(config.GetString("url", "https://localhost:51150/"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %s", err)
	}
	u.Path = "/oauth2"

	oauth2Config := oauth2.Config{
		ClientID:     config.GetString("oauth2.clientId", ""),
		ClientSecret: config.GetString("oauth2.clientSecret", ""),
		Scopes:       config.GetStringSlice("oauth2.scopes", []string{"openid", "profile", "email"}),
		RedirectURL:  u.String(),
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.GetString("oauth2.authUrl", ""),
			TokenURL: config.GetString("oauth2.tokenUrl", ""),
		},
	}

	return &frontend{cookieStore, oauth2Config, store, ipManager, config, logger}, nil
}

func (f *frontend) ServeHTTP() http.Handler {

	router := mux.NewRouter()
	router.Use(f.sessionMiddleware())
	router.HandleFunc("/login", f.login)
	router.HandleFunc("/oauth2", f.authorize)
	router.HandleFunc("/oauth2/access-denied", f.accessDenied)

	graphqlSrv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(f.store, f.ipManager, f.l)}))

	router.Handle("/graphql", graphqlSrv)

	w := webapp.WebHandler(f.l)
	router.PathPrefix("/").HandlerFunc(w)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		router.ServeHTTP(w, r)
	})
}

func (f *frontend) sessionMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/oauth2") {
				next.ServeHTTP(w, r)
				return
			}

			session, _ := f.sessions.Get(r, CookieName)

			if v, _ := session.Values["loggedIn"]; v != "true" {
				if session.Values["sub"] != nil && session.Values["sub"] != "" {
					if _, ok := f.store.IsUserApproved(session.Values["sub"].(string)); ok {
						session.Values["loggedIn"] = "true"
					}
				}
			}

			if strings.HasPrefix(r.URL.Path, "/graphql") {
				if v, ok := session.Values["loggedIn"]; !ok && v != "true" {
					w.WriteHeader(401)
					return
				}
			} else if f.authRedirect(w, r, session) {
				return
			}

			id := session.Values["sub"].(string)
			if id != "" {
				user, err := f.store.GetUserByID(id)
				if err == nil {
					ctx := graph.WithUser(r.Context(), model.User{
						ID:    user.Id,
						Name:  user.Name,
						Email: user.Email,
					})
					r = r.WithContext(ctx)
				}
			}

			session.Save(r, w)

			next.ServeHTTP(w, r)
		})
	}
}

func (f *frontend) authRedirect(w http.ResponseWriter, r *http.Request, session *sessions.Session) bool {
	if v, ok := session.Values["loggedIn"]; ok && v == "true" {
		return false
	}

	state := randomString(20)

	u := f.oauth2Config.AuthCodeURL(state)

	session.Values["state"] = state

	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}

	http.Redirect(w, r, u, http.StatusFound)
	return true
}

func (f *frontend) authorize(w http.ResponseWriter, r *http.Request) {
	session, _ := f.sessions.Get(r, CookieName)

	if cState, ok := session.Values["state"]; ok && len(cState.(string)) > 10 {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Unable to parse form data", http.StatusBadRequest)
			return
		}

		state := r.Form.Get("state")
		if state != cState {
			http.Error(w, "State invalid", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "Missing state on session", http.StatusBadRequest)
		return
	}

	code := r.Form.Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	token, err := f.oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	delete(session.Values, "state")

	u, err := f.getUserInfo(token)
	if err != nil {
		http.Error(w, "Failed to get userInfo", http.StatusInternalServerError)
		f.l.WithError(err).Error("Failed to get userInfo")
		return
	}

	session.Values["sub"] = u.sub
	session.Values["name"] = u.name
	session.Values["email"] = u.email

	if user, ok := f.store.IsUserApproved(u.sub); !ok {
		if user == nil {
			_, err = f.store.AddUser(&store.User{Id: u.sub, Name: u.name, Email: u.email})
			if err != nil {
				f.l.WithError(err).Error("Failed to add user")
			}
		}

		session.Values["loggedIn"] = "false"
		err = session.Save(r, w)
		http.Redirect(w, r, "/oauth2/access-denied", 302)
		return
	}

	session.Values["loggedIn"] = "true"

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (f *frontend) login(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("You are already logged in!"))
}

func (f *frontend) accessDenied(w http.ResponseWriter, r *http.Request) {
	session, _ := f.sessions.Get(r, CookieName)

	if session.Values["sub"] != nil && session.Values["sub"] != "" {
		if _, ok := f.store.IsUserApproved(session.Values["sub"].(string)); ok {
			session.Values["loggedIn"] = "true"
			err := session.Save(r, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/", 302)
			return
		}
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "Your are not allowed to login <a href=\"/\">Retry</a>")
}

func (f *frontend) getUserInfo(token *oauth2.Token) (*userInfo, error) {
	client := f.oauth2Config.Client(context.Background(), token)

	resp, err := client.Get(f.config.GetString("oauth2.userInfoUrl", ""))
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %s", err)
	} else {
		if resp.StatusCode == 200 {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to read body: %s\n", err)
			} else {
				f.l.WithField("userInfoBody", body).Trace("UserInfo body")
				j := make(map[string]interface{})

				err = json.Unmarshal(body, &j)
				if err != nil {
					return nil, fmt.Errorf("failed to parse userInfo as json: %s", err)
				}

				return &userInfo{j["sub"].(string), j["name"].(string), j["email"].(string)}, nil
			}
		} else {
			return nil, fmt.Errorf("failed to get user info with status : %d", resp.StatusCode)
		}
	}
}

func randomString(n int) string {
	const letterBytes = ("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[mrand.Intn(len(letterBytes))]
	}
	return string(b)
}
