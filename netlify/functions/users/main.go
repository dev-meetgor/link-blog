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

func main() {
	lambda.Start(handler)
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
	email := formData.Get("email")
	password := formData.Get("password")

	switch req.HTTPMethod {

	case "POST":
		if userName == "" || password == "" || email == "" {
			return errorResponse(http.StatusBadRequest, "Invalid form data"), nil
		}
		user, err := queries.GetUserByEmail(ctx, email)
		if err != nil {
			//create user
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				return errorResponse(http.StatusInternalServerError, "Failed to hash password"), nil
			}
			_, err = queries.CreateUser(ctx, models.CreateUserParams{
				Username:     userName,
				Email:        email,
				PasswordHash: string(hashedPassword),
			})
			if err != nil {
				return errorResponse(http.StatusInternalServerError, "Failed to create user"), nil
			}
			user, err = queries.GetUserByEmail(ctx, email)
			if err != nil {
				return errorResponse(http.StatusInternalServerError, "Failed to get user"), nil
			}
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
	default:
		return errorResponse(http.StatusMethodNotAllowed, "Method not allowed"), nil
	}

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
