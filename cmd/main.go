package main

import (
	"context"
	"fmt"
	"github.com/nguyenbatam/example-balance/common"
	"github.com/nguyenbatam/example-balance/libs"
	"math/big"
	"strconv"
	"time"
)

func GetTotalETHHolders(inputTime string) (*big.Int, error) {
	timestamp, err := strconv.ParseInt(inputTime, 10, 64)
	if err != nil {
		t, err := time.Parse(common.FormatDate, inputTime)
		if err != nil {
			return nil, err
		}
		timestamp = t.Unix()
	}
	ctx := context.Background()
	blockNumber, err := libs.GetBlockNumberByTimeStamp(ctx, uint64(timestamp))
	holders, err := libs.GetListBAYCNFTHolder(ctx, blockNumber)
	if err != nil {
		return nil, err
	}
	totalBalance, err := libs.GetTotalBalance(ctx, holders, blockNumber)
	if err != nil {
		return nil, err
	}
	return totalBalance, nil
}

func main() {
	inputTime := strconv.FormatInt(time.Now().Unix(), 10)
	totalBalance, err := GetTotalETHHolders(inputTime)
	if err != nil {
		panic(err)
	}
	fmt.Println(totalBalance)
}
