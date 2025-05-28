package cronhandlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alfianX/hommypay_trx/databases/trx/reversals"
	h "github.com/alfianX/hommypay_trx/helper"
	"github.com/alfianX/hommypay_trx/pkg/iso"
)

func (cs *CronService) SafReversal() {
	data, err := cs.reversalService.GetDataSafReversal(context.Background())
	if err != nil {
		h.ErrorLog("Cron AR -> Get data auto reversal: " + err.Error())
		return
	}

	if len(data) > 0 {
		for _, row := range data {
			var isoReqEnc string
			if row.IsoRequest != "" {
				isoReqEnc = row.IsoRequest

				_, issuerConnType, _, issuerService, err := cs.issuerService.GetUrlByIssuerID(context.Background(), row.IssuerID)
				if err != nil {
					h.ErrorLog("Cron AR -> Get url service: " + err.Error())
					continue
				}

				if issuerConnType == 0 {
					h.ErrorLog("Cron AR -> Issuer connection type not found!")
					continue
				}

				if issuerService == "" {
					h.ErrorLog("Cron AR -> Issuer service not found!")
					continue
				}

				ip, port, err := cs.hsmConfigService.GetHSMIpPort(context.Background())
				if err != nil {
					h.ErrorLog("Cron AR -> Get IP PORT HSM: " + err.Error())
					continue
				}

				zek, err := cs.keyConfigService.GetZEK(context.Background())
				if err != nil {
					h.ErrorLog("Cron AR -> Get ZEK: " + err.Error())
					continue
				}

				var isoResEnc string
				var responseCode string
				if issuerConnType == 1 {
					lenIsoReq, err := strconv.ParseInt(isoReqEnc[:4], 16, 0)
					if err != nil {
						h.ErrorLog("Cron AR -> Parse Int: " + err.Error())
						continue
					}

					isoReq, err := h.HSMDecrypt(ip+":"+port, zek, isoReqEnc[4:])
					if err != nil {
						h.ErrorLog("Cron AR -> ISO decrypt: " + err.Error())
						continue
					}
					isoReq = isoReq[:lenIsoReq]

					isoReversal := isoReq

					ret, res := h.SendHost(issuerService, isoReversal, cs.config.TimeoutTrx)
					if ret != 0 {
						if ret == -4 {
							if row.RepeatCount < 3 {
								err := cs.transactionService.UpdateTOReversalFlag(context.Background(), row.TransactionID)
								if err != nil {
									h.ErrorLog("Cron AR -> Update reversal flag TO: " + err.Error())
									continue
								}

								repeatCount := row.RepeatCount + 1

								err = cs.reversalService.UpdateBackFlagReversal(context.Background(), row.TransactionID, repeatCount)
								if err != nil {
									h.ErrorLog("Cron AR -> Update back reversal flag: " + err.Error())
									continue
								}
							}
						}
						h.ErrorLog("Cron AR -> " + res)
						continue
					}

					DERes := iso.Parse(res[14:])
					if DERes[39] != "" {
						responseCode = DERes[39]
					}

					isoResEnc, err = h.HSMEncrypt(ip+":"+port, zek, res)
					if err != nil {
						h.ErrorLog("Cron AR -> ISO encrypt: " + err.Error())
						continue
					}
				} else if issuerConnType == 2 {
					var resp *http.Response

					type send struct {
						TransactionID string `json:"transactionID"`
					}

					dataToSend := send{}
					dataToSend.TransactionID = row.TransactionID

					payload, err := json.Marshal(dataToSend)
					if err != nil {
						h.ErrorLog("Cron AR -> Marshal data send: " + err.Error())
						continue
					}

					exReq, err := http.NewRequest("POST", issuerService+"/reversal", bytes.NewReader(payload))
					if err != nil {
						h.ErrorLog("Cron AR -> Prepare to send : " + err.Error())
						continue
					}

					exReq.Header = map[string][]string{
						"Content-Type": {"application/json"},
					}

					client := &http.Client{
						Timeout: time.Duration(cs.config.TimeoutTrx) * time.Second,
					}
					resp, err = client.Do(exReq)

					if err != nil {
						if strings.Contains(err.Error(), "Timeout") || strings.Contains(err.Error(), "timeout") {
							if row.RepeatCount < 3 {
								err := cs.transactionService.UpdateTOReversalFlag(context.Background(), row.TransactionID)
								if err != nil {
									h.ErrorLog("Cron AR -> Update reversal flag TO: " + err.Error())
									continue
								}

								repeatCount := row.RepeatCount + 1

								err = cs.reversalService.UpdateBackFlagReversal(context.Background(), row.TransactionID, repeatCount)
								if err != nil {
									h.ErrorLog("Cron AR -> Update back reversal flag: " + err.Error())
									continue
								}
							}
						}
						h.ErrorLog("Cron AR -> " + err.Error())
						continue
					}

					defer resp.Body.Close()

					var extResp map[string]interface{}
					err = json.NewDecoder(resp.Body).Decode(&extResp)
					if err != nil {
						h.ErrorLog("Cron AR -> Decode response: " + err.Error())
						continue
					}

					if resp.StatusCode != 200 {
						h.ErrorLog("Cron AR -> " + resp.Status)
						continue
					}

					responseCode = extResp["responseCode"].(string)
				}

				if responseCode == "00" {
					err := cs.transactionService.UpdateReversalFlag(context.Background(), row.TransactionID)
					if err != nil {
						h.ErrorLog("Cron AR -> Update reversal flag: " + err.Error())
						continue
					}
				} else if responseCode == "1B" {
					err = cs.reversalService.UpdateDataReversal(context.Background(), reversals.UpdateDataReversalParams{
						TransactionID: row.TransactionID,
						ResponseCode:  responseCode,
						IsoResponse:   isoResEnc,
					})
					if err != nil {
						h.ErrorLog("Cron AR -> Update reversal: " + err.Error())
						continue
					}

					repeatCount := row.RepeatCount

					err = cs.reversalService.UpdateBackFlagReversal(context.Background(), row.TransactionID, repeatCount)
					if err != nil {
						h.ErrorLog("Cron AR -> Update back reversal flag: " + err.Error())
						continue
					}

					continue
				}

				err = cs.reversalService.UpdateDataReversal(context.Background(), reversals.UpdateDataReversalParams{
					TransactionID: row.TransactionID,
					ResponseCode:  responseCode,
					IsoResponse:   isoResEnc,
				})
				if err != nil {
					h.ErrorLog("Cron AR -> Update reversal: " + err.Error())
					continue
				}

				err = cs.reversalService.CreateAutoReversalLog(context.Background(), row.TransactionID)
				if err != nil {
					h.ErrorLog("Cron AR -> Save reversal log: " + err.Error())
					continue
				}

				err = cs.reversalService.DeleteReversal(context.Background(), row.TransactionID)
				if err != nil {
					h.ErrorLog("Cron AR -> Delete reversal data: " + err.Error())
					continue
				}
			}
		}
	}
}
