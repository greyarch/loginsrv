package caddy

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/tarent/loginsrv/login"
	"github.com/tarent/loginsrv/model"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_ServeHTTP_200(t *testing.T) {	
	//Set the ServeHTTP *http.Request
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Unable to create request: %v", err)
	}
	
	//Set the ServeHTTP http.ResponseWriter
	w := httptest.NewRecorder()
	
	//Set the CaddyHandler config
	configh := login.DefaultConfig()
	configh.Backends = login.Options{"simple": {"bob": "secret"}}
	loginh, err := login.NewHandler(configh)
	if err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	}
	
	//Set the CaddyHandler that will use ServeHTTP
	h := &CaddyHandler{
	next: httpserver.HandlerFunc(func(w http.ResponseWriter, r *http.Request) (int, error) {
		return http.StatusOK, nil // not t.Fatalf, or we will not see what other methods yield
		}),
	config:       login.DefaultConfig(),
	loginHandler: loginh,
	}
	
	//Set user token
	userInfo := model.UserInfo{Sub: "bob", Expiry: time.Now().Add(time.Second).Unix()}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, userInfo)
	validToken, err := token.SignedString([]byte(h.config.JwtSecret))
	if err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	}
	
	//Set cookie for user token on the ServeHTTP http.ResponseWriter
	cookie := http.Cookie{Name: "jwt_token",Value: validToken}
    http.SetCookie(w, &cookie)
	
	status, err := h.ServeHTTP(w, r)
	
	if err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	}
	
	if status != 200 {
		t.Errorf("Expected returned status code to be %d, got %d", 0, status)
	}
}

func Test_ServeHTTP_login(t *testing.T) {
	//Set the ServeHTTP *http.Request
	r, err := http.NewRequest("GET", "/login", nil)
	if err != nil {
		t.Fatalf("Unable to create request: %v", err)
	}
	
	//Set the ServeHTTP http.ResponseWriter
	w := httptest.NewRecorder()
	
	//Set the CaddyHandler config
	configh := login.DefaultConfig()
	configh.Backends = login.Options{"simple": {"bob": "secret"}}
	loginh, err := login.NewHandler(configh)
	if err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	}
	
	//Set the CaddyHandler that will use ServeHTTP
	h := &CaddyHandler{
	next: httpserver.HandlerFunc(func(w http.ResponseWriter, r *http.Request) (int, error) {
		return http.StatusOK, nil // not t.Fatalf, or we will not see what other methods yield
		}),
	config:       login.DefaultConfig(),
	loginHandler: loginh,
	}
	
	status, err := h.ServeHTTP(w, r)
	
	if err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	}
	
	if status != 0 {
		t.Errorf("Expected returned status code to be %d, got %d", 0, status)
	}
}