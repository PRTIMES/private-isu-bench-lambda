handler: main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o handler main.go

.PHONY: deploy_private_isu_bench_lambda
deploy_private_isu_bench_lambda: handler
	MACKEREL_API_KEY='' \
	MACKEREL_SERVICE_NAME='' \
	SPREADSHEETID='' \
	SPREADSHEET_CREDENTIALS_JSON='' \
	SPREADSHEET_RANGE='' \
	FUNCTION_NAME='' \
	FUNCTION_ROLE='' \
	S3_BUCKET='' \
	S3_KEY='' \
	lambroll deploy --function="function.json" --src='.'
