package ethereum

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/nocturna-ta/election/config"
	"github.com/nocturna-ta/golib/ethereum"
)

func GetEthereumClient(cfg *config.BlockchainConfig) (ethereum.Client, error) {
	return ethereum.New(&ethereum.Options{
		URL: cfg.GanacheURL,
	})
}

func GetContractAddress(cfg *config.BlockchainConfig) common.Address {
	return common.HexToAddress(cfg.ContractAddress)

}
