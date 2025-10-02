buildAppLinux:
	go build -o ./bin/hommypay_trx/app/hommypay-trx-app ./cmd/app/main.go

buildAppWindows:
	go build -o ./bin/hommypay_trx/app/hommypay-trx-app.exe ./cmd/app/main.go

buildLogonLinux:
	go build -o ./bin/hommypay_trx/logon/hommypay-trx-logon ./cmd/logon/main.go

buildLogonWindows:
	go build -o ./bin/hommypay_trx/logon/hommypay-trx-logon.exe ./cmd/logon/main.go

buildSaleLinux:
	go build -o ./bin/hommypay_trx/sale/hommypay-trx-sale ./cmd/sale/main.go

buildSaleWindows:
	go build -o ./bin/hommypay_trx/sale/hommypay-trx-sale.exe ./cmd/sale/main.go

buildVoidLinux:
	go build -o ./bin/hommypay_trx/void/hommypay-trx-void ./cmd/void/main.go

buildVoidWindows:
	go build -o ./bin/hommypay_trx/void/hommypay-trx-void.exe ./cmd/void/main.go

buildReversalLinux:
	go build -o ./bin/hommypay_trx/reversal/hommypay-trx-reversal ./cmd/reversal/main.go

buildReversalWindows:
	go build -o ./bin/hommypay_trx/reversal/hommypay-trx-reversal.exe ./cmd/reversal/main.go

buildSettlementLinux:
	go build -o ./bin/hommypay_trx/settlement/hommypay-trx-settlement ./cmd/settlement/main.go

buildSettlementWindows:
	go build -o ./bin/hommypay_trx/settlement/hommypay-trx-settlement.exe ./cmd/settlement/main.go

buildBatchUploadLinux:
	go build -o ./bin/hommypay_trx/batch-upload/hommypay-trx-batch-upload ./cmd/batch-upload/main.go

buildBatchUploadWindows:
	go build -o ./bin/hommypay_trx/batch-upload/hommypay-trx-batch-upload.exe ./cmd/batch-upload/main.go

buildCronLinux:
	go build -o ./bin/hommypay_trx/cron/hommypay-trx-cron ./cmd/cron/main.go

buildCronWindows:
	go build -o ./bin/hommypay_trx/cron/hommypay-trx-cron.exe ./cmd/cron/main.go