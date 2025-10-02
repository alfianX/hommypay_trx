package main

import (
	"context"
	"log"

	// "github.com/alfianX/hommypay_trx/helper"

	"github.com/alfianX/hommypay_trx/internal/app"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run(ctx context.Context) error {
	server, err := app.NewServer()
	if err != nil {
		return err
	}

	// rawData, _ := hex.DecodeString("910ACC93F8F8F2F35FA0303071129F1804000000018609841E0000049D64204172129F180400000001860984240000047C6E0CBF")

	// lenPrefix := fmt.Sprintf("%03d", len(rawData))
	// rawData = append([]byte(lenPrefix), rawData...)

	// emvField := field.NewComposite(emvax.SpecEmv)
	// _, err = emvField.Unpack(rawData)
	// if err != nil {
	// 	return err
	// }

	// data := &emvax.Data{}

	// err = emvField.Unmarshal(data)
	// if err != nil {
	// 	return err
	// }

	// fmt.Println(data.IssuerScriptTemplate1.Value())

	// dataEncrytp := "00A8394398793C1159B7A94ED94F580BEAF53E5B024C0390B6C9F097C32A35E0534EF1FA4F29230EFE1957D311396727050A1F9AEFCFFDBBF0851664A85DC90B03EA876BF34C0B52B24C17019A9E2109E1077A4AA3313899B5C8"
	// fmt.Println("data encrypt: " + dataEncrytp + "\n")
	// dataDecrypt, _ := helper.HSMDecrypt("103.135.5.53:3501", "U28A99A9051BBD56A3FAA1919BB9428AA", dataEncrytp[4:])
	// fmt.Println("hasil decrypt: " + dataDecrypt)

	err = server.Run(ctx)
	return err
}
