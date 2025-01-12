package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mr-destructive/link-blog/embedsql"
	"github.com/mr-destructive/link-blog/models"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"golang.org/x/crypto/bcrypt"
)

var (
	queries *models.Queries
	sqlDB   *sql.DB
)

func TokenValidationMiddleware(next func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		var tokenString string
		cookies := getHeader(req.Headers, "Cookie")
		log.Printf("cookies: %v", cookies)
		for _, cookie := range strings.Split(cookies, ";") {
			c := string(cookie)
			log.Printf("c: %v", c)
			parts := strings.SplitN(c, "=", 2)
			log.Printf("parts: %v", parts)
			if len(parts) == 2 && parts[0] == "auth_token" {
				tokenString = parts[1]
				log.Printf("tokenString: %v", tokenString)
				break
			}
		}

		if tokenString == "" {
			return next(req)
		}

		log.Printf("tokenString C: %v", tokenString)
		claims, err := ValidateToken(tokenString, os.Getenv("JWT_SECRET"))
		log.Printf("tokenString D: %v", tokenString)
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				return errorResponse(http.StatusUnauthorized, "Invalid token signature"), nil
			}
			log.Printf("Error validating token: %v", err)
			return errorResponse(http.StatusInternalServerError, "Error validating token"), nil
		}

		userID, ok := claims.Claims.(jwt.MapClaims)["user_id"].(float64)
		if !ok {
			return errorResponse(http.StatusInternalServerError, "Invalid token claims"), nil
		}

		req.RequestContext.Authorizer = map[string]interface{}{"user_id": int(userID)}
		return next(req)
	}
}

func main() {
	authHandler := TokenValidationMiddleware(handler)
	lambda.Start(authHandler)
}

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	ctx := context.Background()
	dbName := os.Getenv("DB_NAME")
	dbToken := os.Getenv("DB_TOKEN")

	var err error
	dbString := fmt.Sprintf("libsql://%s?authToken=%s", dbName, dbToken)
	db, err := sql.Open("libsql", dbString)
	if err != nil {
		return errorResponse(http.StatusInternalServerError, "Database connection failed"), nil
	}
	defer db.Close()

	queries = models.New(db)
	if _, err := db.ExecContext(ctx, embedsql.DDL); err != nil {
		log.Printf("error creating tables: %v", err)
		return errorResponse(http.StatusInternalServerError, "Database connection failed"), nil
	}

	formData, err := url.ParseQuery(req.Body)
	if err != nil {
		return errorResponse(http.StatusBadRequest, "Invalid form data"), nil
	}
	userName := formData.Get("username")
	password := formData.Get("password")

	if userName == "" || password == "" {
		return errorResponse(http.StatusBadRequest, "Invalid form data"), nil
	}
	user, err := queries.GetUserByUsername(ctx, userName)
	if err != nil {
		return errorResponse(http.StatusInternalServerError, "Database connection failed"), nil
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return errorResponse(http.StatusUnauthorized, "Invalid username or password"), nil
	}

	token, err := CreateToken(user, os.Getenv("JWT_SECRET"))
	if err != nil {
		return errorResponse(http.StatusInternalServerError, "Failed to create token"), nil
	}
	tokenCookie := http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	}
	headers := map[string]string{
		"Content-Type": "text/plain",
		"Set-Cookie":   tokenCookie.String(),
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    headers,
		Body:       "Login successful",
	}, nil
}

func jsonResponse(statusCode int, data interface{}) events.APIGatewayProxyResponse {
	body, _ := json.Marshal(data)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}
}

func errorResponse(statusCode int, message string) events.APIGatewayProxyResponse {
	return jsonResponse(statusCode, map[string]string{"error": message})
}

func CreateToken(user models.User, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString string, secret string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
}

func getHeader(headers map[string]string, key string) string {
	if val, ok := headers[key]; ok {
		return val
	}

	lowerKey := strings.ToLower(key)
	for k, v := range headers {
		if strings.ToLower(k) == lowerKey {
			return v
		}
	}
	return ""
}
