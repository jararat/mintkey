package wordlist

import (
	"bytes"
	"testing"

	. "github.com/tendermint/go-common"
	"github.com/tendermint/go-crypto"
)

func TestEnglish(t *testing.T) {
	for i := 0; i < 5000; i++ {
		numBytes := RandInt()%31 + 1
		data := crypto.CRandBytes(numBytes)

		// Encode data to words
		words := BytesToWords("english", data)

		// Decode data from words
		data2 := make([]byte, numBytes)
		err := WordsToBytes("english", words, data2)
		if err != nil {
			t.Error(err)
			return
		}

		// Ensure that data matches
		if !bytes.Equal(data, data2) {
			t.Errorf("Expected data %X but got %X. Words: %v",
				data, data2, words)
			return
		}
	}
}
