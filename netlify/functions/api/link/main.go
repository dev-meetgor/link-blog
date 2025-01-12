package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mr-destructive/link-blog/embedsql"
	"github.com/mr-destructive/link-blog/models"
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

	switch req.HTTPMethod {
	case "GET":
		return getLinks(req)
	case "POST":
		return createLink(req)
	case "PUT":
		return updateLink(req)
	case "DELETE":
		return deleteLink(req)
	default:
		return events.APIGatewayProxyResponse{StatusCode: 405}, nil
	}
}

func getLinks(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	postId := req.QueryStringParameters["id"]
	if postId != "" {

	} else {
		queries.ListPostsByAuthor()
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
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
