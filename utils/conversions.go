package utils

import (
	"fmt"
	"math/big"
)

// ConvertFloatAmountToBigInt - converts a given float64 amount to a bigint with the correct base
func ConvertFloatAmountToBigInt(amount float64) *big.Int {
	bigAmount := new(big.Float).SetFloat64(amount)
	base := new(big.Float).SetInt(big.NewInt(1000000000000000000))
	bigAmount.Mul(bigAmount, base)
	realAmount := new(big.Int)
	bigAmount.Int(realAmount)

	return realAmount
}

// ConvertNumeralStringToBigFloat - converts a numeral string back to a big float with the correct base set
func ConvertNumeralStringToBigFloat(balance string) (*big.Float, error) {
	floatBalance := new(big.Float)
	floatBalance, ok := floatBalance.SetString(balance)

	if !ok {
		return nil, fmt.Errorf("can't convert balance string %s to a float balance", balance)
	}

	base := new(big.Float).SetInt(big.NewInt(1000000000000000000))
	value := new(big.Float).Quo(floatBalance, base)
	return value, nil
}

// CalculateTotalGasCost - calculates the total gas cost for a given transaction
func CalculateTotalGasCost(gasPrice uint64, gasLimit uint64) *big.Int {
	bigGasPrice := new(big.Int).SetUint64(gasPrice)
	bigGasLimit := new(big.Int).SetUint64(gasLimit)
	bigGasCost := bigGasPrice.Mul(bigGasPrice, bigGasLimit)

	return bigGasCost
}

// CalculateAmountWithoutGasCost - calculates the amount to send for a tx excluding its gas cost
func CalculateAmountWithoutGasCost(amount *big.Int, gasCost *big.Int) *big.Int {
	amountAfterGas := new(big.Int).Sub(amount, gasCost)

	return amountAfterGas
}
