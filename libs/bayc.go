package libs

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	common2 "github.com/nguyenbatam/example-balance/common"
	"github.com/nguyenbatam/example-balance/contracts"
	"math/big"
	"strconv"
	"strings"
)

func GetListBAYCNFTHolder(ctx context.Context, blockNumber uint64) ([]common.Address, error) {
	ethClient, err := ethclient.DialContext(ctx, common2.ETH_RPC)
	defer ethClient.Close()
	if err != nil {
		return nil, err
	}
	contract, err := contracts.NewContractsFilterer(common.HexToAddress(common2.BAYC_CONTRACT_ADDRESS), ethClient)
	if err != nil {
		return nil, err
	}
	start := common2.DEPLOY_CONTRACT_BLOCK_NUMBER
	step := uint64(10000)
	mapTokenOwner := make(map[int64]common.Address)
	end := start + step
	for start < blockNumber {
		if end > blockNumber {
			end = blockNumber
		}
		iter, err := contract.FilterTransfer(&bind.FilterOpts{
			Start: start,
			End:   &end,
		}, nil, nil, nil)
		if err != nil {
			errString := err.Error()
			if strings.Contains(errString, "Try with this block range") {
				endBlock := strings.Split(strings.TrimRight(errString, "Try with this block range"), ",")
				if len(endBlock) > 1 {
					endBlock = strings.Split(endBlock[1], "]")
					decimalValue, err := strconv.ParseUint(strings.TrimSpace(endBlock[0]), 0, 64)
					if err == nil {
						end = decimalValue
						continue
					}
				}
			}
		}
		if err != nil {
			return nil, err
		}
		for iter.Next() {
			if iter.Error() != nil {
				return nil, iter.Error()
			}
			mapTokenOwner[iter.Event.TokenId.Int64()] = iter.Event.To
		}
		start = end
		end = start + step
	}
	mapHolders := map[common.Address]bool{}
	for _, address := range mapTokenOwner {
		mapHolders[address] = true
	}
	var holders []common.Address
	for address, _ := range mapHolders {
		holders = append(holders, address)
	}
	return holders, nil
}

func GetTotalBalance(ctx context.Context, holders []common.Address, blockNumber uint64) (*big.Int, error) {
	ethClient, err := ethclient.DialContext(ctx, common2.ETH_RPC)
	defer ethClient.Close()
	if err != nil {
		return nil, err
	}
	totalETH := big.NewInt(0)
	for _, holder := range holders {
		ethBalance, err := ethClient.BalanceAt(ctx, holder, big.NewInt(0).SetUint64(blockNumber))
		if err != nil {
			continue
		}
		totalETH = new(big.Int).Add(totalETH, ethBalance)
	}
	return totalETH, nil
}

func GetBlockNumberByTimeStamp(ctx context.Context, timestamp uint64) (uint64, error) {
	ethClient, err := ethclient.DialContext(ctx, common2.ETH_RPC)
	defer ethClient.Close()
	if err != nil {
		return 0, err
	}
	currentBlock, err := ethClient.BlockByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}
	if currentBlock.Time() < timestamp {
		return currentBlock.NumberU64(), nil
	}
	startBlock := uint64(0)
	endBlock := currentBlock.NumberU64()

	for startBlock+1 < endBlock {
		middleBlock := (startBlock + endBlock) / 2
		block, err := ethClient.BlockByNumber(ctx, big.NewInt(0).SetUint64(middleBlock))
		if err != nil {
			return 0, err
		}
		if block.Time() == timestamp {
			return middleBlock, nil
		} else if block.Time() < timestamp {
			startBlock = middleBlock
		} else {
			endBlock = middleBlock
		}
	}
	return startBlock, nil
}
