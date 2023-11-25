package commons

import (
	"errors"
	"math"
	"strconv"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type MoneyAmount struct {
	rawNum int
}

func NewMoneyAmount(strAmnt string) *MoneyAmount {
	amnt := &MoneyAmount{}
	amnt.unmarshalByteArr([]byte(strAmnt))
	return amnt
}

func (amnt *MoneyAmount) Add(amnt2 *MoneyAmount) *MoneyAmount {
	return &MoneyAmount{rawNum: amnt.rawNum + amnt2.rawNum}
}

func (amnt *MoneyAmount) String() string {
	floNum := float64(amnt.rawNum) / 100.0
	return strconv.FormatFloat(floNum, 'f', 2, 64)
}

func (amnt *MoneyAmount) UnmarshalJSON(bytArrNum []byte) error {
	return amnt.unmarshalByteArr(bytArrNum)
}

func (amnt *MoneyAmount) MarshalJSON() ([]byte, error) {
	return []byte(amnt.String()), nil
}

func (amnt *MoneyAmount) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	strNum := amnt.String()
	av.N = &strNum
	return nil
}

func (amnt *MoneyAmount) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	if av == nil {
		return errors.New("Attribute value must not be null.")
	}

	if av.N == nil {
		return errors.New("Attribute value must be a number.")
	}

	return amnt.unmarshalByteArr([]byte(*av.N))
}

func (amnt *MoneyAmount) unmarshalByteArr(bytArr []byte) error {
	if bytArr == nil {
		return errors.New("Value must not be null.")
	}

	dotIdx := -1
	for i, byt := range bytArr {
		if (byt < '0' || '9' < byt) && byt != '.' {
			return errors.New("Value is not a number.")
		}

		if byt == '.' {
			dotIdx = i
		}
	}

	if dotIdx == -1 {
		intNum, err := strconv.Atoi(string(bytArr))
		if err != nil {
			return err
		}
		amnt.rawNum = intNum * 100

		return nil
	}

	scale := len(bytArr) - dotIdx - 1
	strNum := string(bytArr[:dotIdx]) + string(bytArr[dotIdx+1:])
	floNum, err := strconv.ParseFloat(strNum, 64)
	if err != nil {
		return err
	}
	amnt.rawNum = int(math.Round(floNum / math.Pow(10, float64(scale-2))))

	return nil
}
