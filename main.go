package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mackerelio/mackerel-client-go"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// 同じものをコピペ
// https://github.com/catatsuy/private-isu/blob/abba9a52a4b561650a898dedc27881b847a591a0/benchmarker/cli.go#L42-L48
type Output struct {
	Pass     bool     `json:"pass"`
	Score    int64    `json:"score"`
	Suceess  int64    `json:"success"`
	Fail     int64    `json:"fail"`
	Messages []string `json:"messages"`
}

type Message struct {
	Message string `json:"message"`
}

func getTeamNameFromSpreadSheets(clientIP string) (string, error) {
	spreadsheetID := os.Getenv("SPREADSHEETID")
	credential := option.WithCredentialsJSON([]byte(os.Getenv("SPREADSHEET_CREDENTIALS_JSON")))
	srv, err := sheets.NewService(context.TODO(), credential)
	if err != nil {
		log.Println(err)
		return "", err
	}

	readRange := os.Getenv("SPREADSHEET_RANGE")
	res, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Println(err)
		return "", err
	}
	if len(res.Values) == 0 {
		err := errors.New("could not get data")
		log.Println(err)
		return "", err
	}
	for _, row := range res.Values {
		if row[0] == clientIP {
			s, ok := row[1].(string)
			if !ok {
				return "", errors.New("failed to convert team name")
			}
			return s, nil
		}
	}
	return "", errors.New("client ip address and team name did not match")
}

func getError(err error) (events.LambdaFunctionURLResponse, error) {
	var message = Message{}
	message.Message = err.Error()
	msg, err := json.Marshal(message)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			Body:       err.Error(),
			StatusCode: 400,
		}, err
	}
	log.Println(string(msg))
	return events.LambdaFunctionURLResponse{
		Body:       string(msg),
		StatusCode: 400,
	}, err
}

func lambdaHandler(req events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	clientIP := req.RequestContext.HTTP.SourceIP
	teamName, err := getTeamNameFromSpreadSheets(clientIP)
	if err != nil {
		return getError(err)
	}

	targetURL := "http://" + clientIP
	cmd := exec.Command("./bin/benchmarker", "-t", targetURL, "-u", "./userdata")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Printf("Stdout: %s\n", stdout.String())
		log.Printf("Stderr: %s\n", stderr.String())
		return getError(err)
	}

	var output Output
	err = json.Unmarshal(stdout.Bytes(), &output)
	if err != nil {
		return getError(err)
	}

	// Mackerelにチーム名とスコアを渡す
	apiKey := os.Getenv("MACKEREL_API_KEY")
	serviceName := os.Getenv("MACKEREL_SERVICE_NAME")
	client := mackerel.NewClient(apiKey)
	err = client.PostServiceMetricValues(serviceName, []*mackerel.MetricValue{
		{
			Name:  "Score." + teamName,
			Time:  time.Now().Unix(),
			Value: output.Score,
		},
	})
	if err != nil {
		return getError(err)
	}

	return events.LambdaFunctionURLResponse{
		Body:       stdout.String(),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(lambdaHandler)
}
