package handlers

import (
	"context"
	"fmt"
)

func (cs *Service) Transactions() {
	data, err := cs.transactionDataService.GetTrxData(context.Background())

	if err != nil {
		fmt.Println(err)
	}

	if len(data) > 0 {
		for _, row := range data {
			switch row.TransactionType {
			case "01":
				cs.SaleTrx(row)
			case "31":
				cs.VoidTrx(row)
			case "41":
				cs.ReversalTrx(row)
			}
		}
	}
}