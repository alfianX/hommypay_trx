buildAll:
	go env -w GOOS=linux
	go build -o ./bin/hommypay_trx/app/hommypay-trx-app ./cmd/app/main.go
	go build -o ./bin/hommypay_trx/logon/hommypay-trx-logon ./cmd/logon/main.go
	go build -o ./bin/hommypay_trx/sale/hommypay-trx-sale ./cmd/sale/main.go
	go build -o ./bin/hommypay_trx/void/hommypay-trx-void ./cmd/void/main.go
	go build -o ./bin/hommypay_trx/reversal/hommypay-trx-reversal ./cmd/reversal/main.go
	go build -o ./bin/hommypay_trx/settlement/hommypay-trx-settlement ./cmd/settlement/main.go
	go build -o ./bin/hommypay_trx/batch-upload/hommypay-trx-batch-upload ./cmd/batch-upload/main.go
	go env -w GOOS=windows
	go build -o ./bin/hommypay_trx/app/hommypay-trx-app.exe ./cmd/app/main.go
	go build -o ./bin/hommypay_trx/logon/hommypay-trx-logon.exe ./cmd/logon/main.go
	go build -o ./bin/hommypay_trx/sale/hommypay-trx-sale.exe ./cmd/sale/main.go
	go build -o ./bin/hommypay_trx/void/hommypay-trx-void.exe ./cmd/void/main.go
	go build -o ./bin/hommypay_trx/reversal/hommypay-trx-reversal.exe ./cmd/reversal/main.go
	go build -o ./bin/hommypay_trx/settlement/hommypay-trx-settlement.exe ./cmd/settlement/main.go
	go build -o ./bin/hommypay_trx/batch-upload/hommypay-trx-batch-upload.exe ./cmd/batch-upload/main.go

buildApp:
	go env -w GOOS=linux
	go build -o ./bin/hommypay_trx/app/hommypay-trx-app ./cmd/app/main.go
	go env -w GOOS=windows
	go build -o ./bin/hommypay_trx/app/hommypay-trx-app.exe ./cmd/app/main.go

buildLogon:
	go env -w GOOS=linux
	go build -o ./bin/hommypay_trx/logon/hommypay-trx-logon ./cmd/logon/main.go
	go env -w GOOS=windows
	go build -o ./bin/hommypay_trx/logon/hommypay-trx-logon.exe ./cmd/logon/main.go

buildSale:
	go env -w GOOS=linux
	go build -o ./bin/hommypay_trx/sale/hommypay-trx-sale ./cmd/sale/main.go
	go env -w GOOS=windows
	go build -o ./bin/hommypay_trx/sale/hommypay-trx-sale.exe ./cmd/sale/main.go

buildVoid:
	go env -w GOOS=linux
	go build -o ./bin/hommypay_trx/void/hommypay-trx-void ./cmd/void/main.go
	go env -w GOOS=windows
	go build -o ./bin/hommypay_trx/void/hommypay-trx-void.exe ./cmd/void/main.go

buildReversal:
	go env -w GOOS=linux
	go build -o ./bin/hommypay_trx/reversal/hommypay-trx-reversal ./cmd/reversal/main.go
	go env -w GOOS=windows
	go build -o ./bin/hommypay_trx/reversal/hommypay-trx-reversal.exe ./cmd/reversal/main.go

buildSettlement:
	go env -w GOOS=linux
	go build -o ./bin/hommypay_trx/settlement/hommypay-trx-settlement ./cmd/settlement/main.go
	go env -w GOOS=windows
	go build -o ./bin/hommypay_trx/settlement/hommypay-trx-settlement.exe ./cmd/settlement/main.go

buildBatchUpload:
	go env -w GOOS=linux
	go build -o ./bin/hommypay_trx/batch-upload/hommypay-trx-batch-upload ./cmd/batch-upload/main.go
	go env -w GOOS=windows
	go build -o ./bin/hommypay_trx/batch-upload/hommypay-trx-batch-upload.exe ./cmd/batch-upload/main.go