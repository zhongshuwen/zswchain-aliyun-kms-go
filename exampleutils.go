package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/zhongshuwen/zswchain-go/ecc"
	"github.com/zhongshuwen/zswchain-go/system"

	zsw "github.com/zhongshuwen/zswchain-go"
)

var version = "dev"

func Quit(message string, args ...interface{}) {
	fmt.Printf(message+"\n", args...)
	os.Exit(1)
}

func NoError(err error, message string, args ...interface{}) {
	if err != nil {
		Quit(message+": "+err.Error(), args...)
	}
}

func toJson(v interface{}) string {
	out, err := json.MarshalIndent(v, "", "  ")
	NoError(err, "unable to marshal json")

	return string(out)
}
func runTxBasic(ctx context.Context, api *zsw.API, actions []*zsw.Action) string {

	txOpts := &zsw.TxOptions{}
	if err := txOpts.FillFromChain(ctx, api); err != nil {
		panic(fmt.Errorf("filling tx opts: %w", err))
	}

	tx := zsw.NewTransaction(actions, txOpts)
	signedTx, packedTx, err := api.SignTransaction(ctx, tx, txOpts.ChainID, zsw.CompressionNone)
	if err != nil {
		panic(fmt.Errorf("sign transaction: %w", err))
	}
	hash, err := packedTx.ID()
	if err != nil {
		panic(fmt.Errorf("id err transaction: %w", err))
	}
	fmt.Printf("signature: %s\nhash: %s\n", toJson(signedTx.Signatures), hash)

	response, err := api.PushTransaction(context.Background(), packedTx)
	if err != nil {
		panic(fmt.Errorf("push transaction: %w", err))
	}

	fmt.Printf("Transaction [%s] submitted to the network succesfully.\n", hex.EncodeToString(response.Processed.ID))
	return hex.EncodeToString(response.Processed.ID)
}
func UuidToUint128OrQuit(uuidString string) zsw.Uint128 {
	var x zsw.Uint128
	NoError(x.FromUuidString(uuidString), "Invalid uuid: '%s'", uuidString)
	return x

}

func GetActionsCreateUserWithResources(creator zsw.AccountName, newAccount zsw.AccountName, newAccountPublicKey string, ramBytes uint32, cpuAmount zsw.Asset, netAmount zsw.Asset) []*zsw.Action {
	if cpuAmount.Amount == 0 && netAmount.Amount == 0 {
		return []*zsw.Action{
			//创建可信节点联盟链账号
			system.NewNewAccount(
				creator,    //中数文的内容审核管理账号
				newAccount, //联盟链账号名字
				ecc.MustNewPublicKey(newAccountPublicKey), //账号公钥
			),
			system.NewBuyRAMBytes(
				creator, //需要中数文签名
				newAccount,
				ramBytes,
			),
		}
	} else {

		return []*zsw.Action{
			system.NewNewAccount(
				creator,    //中数文的内容审核管理账号
				newAccount, //联盟链账号名字
				ecc.MustNewPublicKey(newAccountPublicKey), //账号公钥
			),
			system.NewBuyRAMBytes(
				creator, //需要中数文签名
				newAccount,
				ramBytes,
			),
			system.NewDelegateBW(
				creator, //需要中数文签名
				newAccount,
				cpuAmount,
				netAmount,
				true,
			),
		}
	}
}
