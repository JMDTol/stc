
package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type StellarNet struct {
	NetworkId string
	Horizon string
}
var Networks = map[string]StellarNet{
	"main": { "Public Global Stellar Network ; September 2015",
		"https://horizon.stellar.org/"},
	"test": { "Test SDF Network ; September 2015",
		"https://horizon-testnet.stellar.org/"},
}

func get(net *StellarNet, query string) []byte {
	resp, err := http.Get(net.Horizon + query)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}
	return body
}

type HorizonAccountEntry struct {
	Sequence json.Number
	Thresholds struct {
		Low_threshold uint8
		Med_threshold uint8
		High_threshold uint8
	}
	Signers []struct {
		Key string
		Weight uint32
	}
}

func GetAccountEntry(net *StellarNet, acct string) *HorizonAccountEntry {
	if body := get(net, "accounts/" + acct); body != nil {
		var ae HorizonAccountEntry
		if err := json.Unmarshal(body, &ae); err != nil {
			return nil
		}
		return &ae
	}
	return nil
}

func GetLedgerHeader(net *StellarNet) (ret *LedgerHeader) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			ret = nil
		}
	}()

	body := get(net, "ledgers?limit=1&order=desc")
	if body == nil {
		return nil
	}

	var lhx struct {
		Embedded struct {
			Records []struct {
				Header_xdr string
			}
		} `json:"_embedded"`
	}
	if err := json.Unmarshal(body, &lhx);
	err != nil || len(lhx.Embedded.Records) == 0 {
		panic(err)
	}

	ret = &LedgerHeader{}
	b64i := base64.NewDecoder(base64.StdEncoding,
		strings.NewReader(lhx.Embedded.Records[0].Header_xdr))
	ret.XdrMarshal(&XdrIn{b64i}, "")
	return
}

func PostTransaction(
	net *StellarNet, e *TransactionEnvelope) *TransactionResult {
	tx := txOut(e)
	resp, err := http.PostForm(net.Horizon + "/transactions",
		url.Values{"tx": {tx}})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil
	}
	defer resp.Body.Close()

	js := json.NewDecoder(resp.Body)
	var res struct {
		Result_xdr string
		Extras struct {
			Result_xdr string
		}
	}
	if err = js.Decode(&res); err != nil {
		fmt.Fprintf(os.Stderr, "PostTransaction: %s\n", err.Error())
		return nil
	}
	if res.Result_xdr == "" { res.Result_xdr = res.Extras.Result_xdr }

	var ret TransactionResult
	if err = txIn(&ret, res.Result_xdr); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid result_xdr\n")
		return nil
	}
	return &ret
}