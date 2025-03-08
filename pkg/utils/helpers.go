package utils

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"regexp"
)

func StringToTx(signedTx string) (*types.Transaction, error) {
	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(common.FromHex(signedTx)); err != nil {
		return nil, err
	}

	return tx, nil

}

func IsNotUUID(s string) bool {
	uuidRegex := `^[a-f0-9]{8}-[a-f0-9]{4}-[1-5][a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$`
	re := regexp.MustCompile(uuidRegex)
	return !re.MatchString(s)
}
