package wordlist

import (
	"errors"
	"math/big"
	"sort"
	"strings"

	. "github.com/tendermint/go-common"
)

// Encodes data into words
func BytesToWords(bank string, data []byte) (words []string) {
	wordsAll := GetWords("english")

	// 2048 words per bank, which is 2^11.
	numWords := (8*len(data) + 10) / 11

	n2048 := big.NewInt(2048)
	nData := big.NewInt(0).SetBytes(data)
	nRem := big.NewInt(0)
	// Alternative, use condition "nData.BitLen() > 0"
	// to allow for shorter words when data has leading 0's
	for i := 0; i < numWords; i++ {
		nData.DivMod(nData, n2048, nRem)
		rem := nRem.Int64()
		words = append(words, wordsAll[rem])
	}
	return words
}

// Decodes words into bytes
func WordsToBytes(bank string, words []string, dest []byte) error {
	wordsAll := GetWords("english")

	// 2048 words per bank, which is 2^11.
	numWords := (8*len(dest) + 10) / 11
	if numWords != len(words) {
		return errors.New(Fmt("Expected %v words for %v dest bytes", numWords, len(dest)))
	}

	n2048 := big.NewInt(2048)
	nData := big.NewInt(0)
	for i := 0; i < numWords; i++ {
		rem := GetWordIndex(wordsAll, words[numWords-i-1])
		if rem < 0 {
			return errors.New(Fmt("Unrecognized word %v for bank %v", words[i], bank))
		}
		nRem := big.NewInt(int64(rem))
		nData.Mul(nData, n2048)
		nData.Add(nData, nRem)
	}
	nDataBytes := nData.Bytes()
	if len(nDataBytes) > len(dest) {
		return errors.New(Fmt("Value %v (len=%v) overflows dest len %v",
			nData, len(nDataBytes), len(dest)))
	}
	copy(dest[len(dest)-len(nDataBytes):], nDataBytes)
	return nil
}

func GetWords(bank string) (words []string) {
	wordsAllStr, err := Asset("wordlist/" + bank + ".txt")
	if err != nil {
		Exit("Error loading wordlist: " + err.Error())
	}
	wordsAll := strings.Split(strings.TrimSpace(string(wordsAllStr)), "\n")
	return wordsAll
}

func GetWordIndex(words []string, word string) int {
	idx := sort.SearchStrings(words, word)
	if idx >= len(words) {
		return -1
	}
	if words[idx] != word {
		return -1
	}
	return idx
}
