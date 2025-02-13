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

	// dataEncrytp := "00F8E281E50BB1754564A94ED94F580BEAF520171ED258A94EA5716AE387AEF276EAD9D3040ECD11208C714667806431C41346E4160A7B514ACA3ABCCF587F5CD9BA275D845854BBF80EC5B1DD1D4E3ED89C97AAC7FD65F7FE5D27269979687F95879AFCB541A681ECCC69BC1AADCD7DE825EE28796EBC8D391A06D29FC00CF40ED9"
	// fmt.Println("data encrypt: " + dataEncrytp + "\n")
	// dataDecrypt, _ := helper.HSMDecrypt("103.135.5.53:3501", "U28A99A9051BBD56A3FAA1919BB9428AA", dataEncrytp[4:])
	// fmt.Println("hasil decrypt: " + dataDecrypt)

	err = server.Run(ctx)
	return err
}