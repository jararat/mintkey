package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/howeyc/gopass"
	. "github.com/tendermint/go-common"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-crypto/bcrypt"
)

const (
	blockTypePrivKey = "TENDERMINT PRIVATE KEY"
)

func readlineKeyboardPass() string {
	str, err := gopass.GetPasswdMasked()
	if err != nil {
		Exit("Error reading from keyboard: " + err.Error())
	}
	return string(str)
}

func readlineKeyboard() string {
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		Exit("Error reading from keyboard: " + err.Error())
	}
	return text
}

func encryptArmorPrivKey(privKey crypto.PrivKey, passphrase string) string {
	saltBytes, encBytes := encryptPrivKey(privKey, passphrase)
	header := map[string]string{
		"kdf":  "bcrypt",
		"salt": Fmt("%X", saltBytes),
	}
	armorStr := crypto.EncodeArmor(blockTypePrivKey, header, encBytes)
	return armorStr
}

func unarmorDecryptPrivKey(armorStr string, passphrase string) (crypto.PrivKey, error) {
	blockType, header, encBytes, err := crypto.DecodeArmor(armorStr)
	if err != nil {
		return nil, err
	}
	if blockType != blockTypePrivKey {
		return nil, fmt.Errorf("Unrecognized armor type: %v", blockType)
	}
	if header["kdf"] != "bcrypt" {
		return nil, fmt.Errorf("Unrecognized KDF type: %v", header["KDF"])
	}
	if header["salt"] == "" {
		return nil, fmt.Errorf("Missing salt bytes")
	}
	saltBytes, err := hex.DecodeString(header["salt"])
	if err != nil {
		return nil, fmt.Errorf("Error decoding salt: %v", err.Error())
	}
	privKey, err := decryptPrivKey(saltBytes, encBytes, passphrase)
	return privKey, err
}

func encryptPrivKey(privKey crypto.PrivKey, passphrase string) (saltBytes []byte, encBytes []byte) {
	saltBytes = crypto.CRandBytes(16)
	key, err := bcrypt.GenerateFromPassword(saltBytes, []byte(passphrase), 12) // TODO parameterize.  12 is good today (2016)
	if err != nil {
		Exit("Error generating bcrypt key from passphrase: " + err.Error())
	}
	key = crypto.Sha256(key) // Get 32 bytes
	privKeyBytes := privKey.Bytes()
	return saltBytes, crypto.EncryptSymmetric(privKeyBytes, key)
}

func decryptPrivKey(saltBytes []byte, encBytes []byte, passphrase string) (privKey crypto.PrivKey, err error) {
	key, err := bcrypt.GenerateFromPassword(saltBytes, []byte(passphrase), 12) // TODO parameterize.  12 is good today (2016)
	if err != nil {
		Exit("Error generating bcrypt key from passphrase: " + err.Error())
	}
	key = crypto.Sha256(key) // Get 32 bytes
	privKeyBytes, err := crypto.DecryptSymmetric(encBytes, key)
	if err != nil {
		return nil, err
	}
	privKey, err = crypto.PrivKeyFromBytes(privKeyBytes)
	return privKey, err
}

func defaultPath(file string) string {
	return os.Getenv("HOME") + "/.mintkey/" + file
}
