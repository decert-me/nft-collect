package utils

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/exp/constraints"
	"regexp"
)

// SliceIsExist 判断元素是否在slice
func SliceIsExist[T comparable](slice []T, val T) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func StructToMap(inter []any) map[string]interface{} {
	var m map[string]interface{}
	for _, v := range inter {
		if v == nil {
			continue
		}
		ja, _ := json.Marshal(v)
		json.Unmarshal(ja, &m)
	}
	return m
}

func MapPushStruct(m map[string]interface{}, inter []any) map[string]interface{} {
	for _, v := range inter {
		ja, _ := json.Marshal(v)
		json.Unmarshal(ja, &m)
	}
	return m
}

func SliceMax[T constraints.Ordered](slice []T) (index int, m T) {
	for i, e := range slice {
		if i == 0 || e > m {
			m = e
			index = i
		}
	}
	return
}

func SliceMin[T constraints.Ordered](slice []T) (index int, m T) {
	for i, e := range slice {
		if i == 0 || e < m {
			m = e
			index = i
		}
	}
	return
}

func RemoveDuplicate[T string | int](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func IsValidAddress(iaddress interface{}) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	switch v := iaddress.(type) {
	case string:
		return re.MatchString(v)
	case common.Address:
		return re.MatchString(v.Hex())
	default:
		return false
	}
}

// VerifySig 校验签名
func VerifySig(from, sigHex string, msg []byte) bool {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	sig := hexutil.MustDecode(sigHex)

	msg = accounts.TextHash(msg)
	if sig[crypto.RecoveryIDOffset] == 27 || sig[crypto.RecoveryIDOffset] == 28 {
		sig[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1
	}

	recovered, err := crypto.SigToPub(msg, sig)
	if err != nil {
		return false
	}

	recoveredAddr := crypto.PubkeyToAddress(*recovered)
	return from == recoveredAddr.Hex()
}
