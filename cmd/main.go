package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"igor/proto/contract"
	"math/big"
)

var (
	action     = flag.String("action", "", "set service config file")
	amount     = flag.String("amount", "", "set service config file")
	amountCoef = big.NewInt(1000000000000000000)
	amountCdai = big.NewInt(1e8)
	address    = common.HexToAddress("0x5105123b58f6313f21f357a9a3a885d1e3b934cd")
)

func main() {
	flag.Parse()

	client, err := ethclient.Dial("https://kovan.infura.io/v3/fd2ed624844d466a8d882ca668ceea07")
	if err != nil {
		fmt.Println(err)
		return
	}

	compound, err := contract.NewContract(common.HexToAddress("0xe7bc397DBd069fC7d0109C0636d06888bb50668c"), client)
	if err != nil {
		fmt.Println(err)
		return
	}

	tx, err := getTx(client)
	if err != nil {
		fmt.Println(err)
		return
	}

	var amountBig *big.Int
	var ok bool

	if *amount != "" {
		amountBig, ok = big.NewInt(0).SetString(*amount, 10)
		if !ok {
			fmt.Printf("Can't parse %s as big int", *amount)
			return
		}
	}

	switch *action {
	case "deposit":
		trans, err := compound.Mint(tx, big.NewInt(0).Mul(amountBig, amountCoef))
		if err != nil {
			fmt.Println(err)
			return
		}

		fmtTrans(trans)
		break
	case "withdraw":
		trans, err := compound.RedeemUnderlying(tx, big.NewInt(0).Mul(amountBig, amountCoef))
		if err != nil {
			fmt.Println(err)
			return
		}

		fmtTrans(trans)
		break
	case "info":
		co := &bind.CallOpts{}
		//
		//name, err := compound.Name(co)
		//if err != nil {
		//	fmt.Println(err)
		//}
		//
		//fmt.Println(`Name: `, name)
		//
		//symbol, err := compound.Symbol(co)
		//if err != nil {
		//	fmt.Println(err)
		//}
		//
		//fmt.Println(`Symbol: `, symbol)
		//
		//decimals, err := compound.Decimals(co)
		//if err != nil {
		//	fmt.Println(err)
		//}
		//
		//fmt.Println("Decimals:", decimals)
		//
		//storedBalance, err := compound.BorrowBalanceStored(co, address)
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
		//
		//fmt.Println("Stored balance:", storedBalance.String())
		//
		cDaiBalance, err := compound.BalanceOf(co, address)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Balance: ", cDaiBalance)

		//r, err := compound.GetCash(co)
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
		//
		//fmt.Println("Cash: ", r.Div(r, amountCdai).String())

		break
	}
}

func fmtTrans(trans *types.Transaction) {

	fmt.Println("Tx:", trans.Hash().Hex())
	fmt.Println("Gas limit:", trans.Gas())
	fmt.Println("Gas cost:", trans.Cost())
	fmt.Println("Gas price:", trans.GasPrice().String())
}

func getTx(client *ethclient.Client) (*bind.TransactOpts, error) {
	fmt.Println("compund inited")

	key := "8DE677CF3188A2F05C908535B32A65DAE32B3A442B67AF78DCEE701F0EE63F28"
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	tx := bind.NewKeyedTransactor(privateKey)
	tx.Nonce = big.NewInt(int64(nonce))
	tx.Value = big.NewInt(0)
	tx.GasPrice = gasPrice
	tx.GasLimit = 1 * 1000 * 1000
	return tx, nil
}
