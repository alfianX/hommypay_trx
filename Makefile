buildAppLinux:
	go env -w GOOS=linux
	go build -o ./bin/hommypay_trx/app/hommypay-trx-app ./cmd/app/main.go
	go env -w GOOS=windows

buildAppWindows:
	go build -o ./bin/hommypay_trx/app/hommypay-trx-app.exe ./cmd/app/main.go

buildLogonLinux:
	go env -w GOOS=linux
	go build -o ./bin/hommypay_trx/logon/hommypay-trx-logon ./cmd/logon/main.go
	go env -w GOOS=windows

buildLogonWindows:
	go build -o ./bin/hommypay_trx/logon/hommypay-trx-logon.exe ./cmd/logon/main.go

buildSaleLinux:
	go env -w GOOS=linux
	go build -o ./bin/hommypay_trx/sale/hommypay-trx-sale ./cmd/sale/main.go
	go env -w GOOS=windows

buildSaleWindows:
	go build -o ./bin/hommypay_trx/sale/hommypay-trx-sale.exe ./cmd/sale/main.go

buildVoidLinux:
	go env -w GOOS=linux
	go build -o ./bin/hommypay_trx/void/hommypay-trx-void ./cmd/void/main.go
	go env -w GOOS=windows

buildVoidWindows:
	go build -o ./bin/hommypay_trx/void/hommypay-trx-void.exe ./cmd/void/main.go

buildReversalLinux:
	go env -w GOOS=linux
	go build -o ./bin/hommypay_trx/reversal/hommypay-trx-reversal ./cmd/reversal/main.go
	go env -w GOOS=windows

buildReversalWindows:
	go build -o ./bin/hommypay_trx/reversal/hommypay-trx-reversal.exe ./cmd/reversal/main.go

buildSettlementLinux:
	go env -w GOOS=linux
	go build -o ./bin/hommypay_trx/settlement/hommypay-trx-settlement ./cmd/settlement/main.go
	go env -w GOOS=windows

buildSettlementWindows:
	go build -o ./bin/hommypay_trx/settlement/hommypay-trx-settlement.exe ./cmd/settlement/main.go

buildBatchUploadLinux:
	go env -w GOOS=linux
	go build -o ./bin/hommypay_trx/batch-upload/hommypay-trx-batch-upload ./cmd/batch-upload/main.go
	go env -w GOOS=windows

buildBatchUploadWindows:
	go build -o ./bin/hommypay_trx/batch-upload/hommypay-trx-batch-upload.exe ./cmd/batch-upload/main.go

buildCronLinux:
	go env -w GOOS=linux
	go build -o ./bin/hommypay_trx/cron/hommypay-trx-cron ./cmd/cron/main.go
	go env -w GOOS=windows

buildCronWindows:
	go build -o ./bin/hommypay_trx/cron/hommypay-trx-cron.exe ./cmd/cron/main.go