package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/howeyc/gopass"
	. "github.com/tendermint/go-common"
	"github.com/tendermint/go-crypto"
	"github.com/tendermint/go-wire"
	"github.com/tendermint/mintkey/wordlist"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	// Mix user-supplied random bytes
	//fmt.Print("\033[?1049h\033[H")
	fmt.Println("We need to get some random bits of entropy from you.\nPlease smash the keyboard, and press Enter when done.")
	randStr := readlineKeyboardPass()
	//fmt.Print("\033[?1049l")
	crypto.MixEntropy([]byte(randStr))
}

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

func main() {
	// Generage a privKey and a secret (seed)
	privKey, secret := generatePrivKey()

	// Print words for secret
	words := wordlist.BytesToWords("english", secret)
	fmt.Printf("Your seed phrase is:\n\n > %v\n", strings.Join(words, " "))

	fmt.Println("\nWrite the seed phrase and store it in a secure location.")
	fmt.Print("Enter the seed phrase again:\n\n > ")
	words2 := readlineKeyboard()
	words2 = strings.TrimSpace(words2)
	fmt.Println("")
	if strings.Join(words, " ") != words2 {
		Exit("Seed phrase did not match!")
	}

	// Encrypt key
	fmt.Println("Enter a passphrase to encrypt your secret key:")
	passStr := readlineKeyboardPass()
	fmt.Println("Re-enter the passphrase:")
	passStr2 := readlineKeyboardPass()
	if passStr != passStr2 {
		Exit("Passphrase didn't match!")
	}
	encBytes := encryptPrivKey(privKey, passStr)
	header := map[string]string{"Encryption": "NACLv0"}
	armorStr := crypto.EncodeArmor("TENDERMINT PRIVATE KEY", header, encBytes)

	// Save armored & encrypted key to file
	fmt.Printf("Armored:\n%v\n", armorStr)

}

// Generates a 128bit secret, and returns the generated PrivKey
func generatePrivKey() (privKey crypto.PrivKey, secret []byte) {
	secret = crypto.CRandBytes(16)
	privKey = crypto.GenPrivKeyEd25519FromSecret(secret)
	return
}

func encryptPrivKey(privKey crypto.PrivKey, passphrase string) []byte {
	key, err := bcrypt.GenerateFromPassword([]byte(passphrase), 12) // TODO parameterize.  12 is good today (2016)
	if err != nil {
		Exit("Error generating bcrypt key from passphrase: " + err.Error())
	}
	key = crypto.Sha256(key) // Get 32 bytes
	privKeyBytes := wire.BinaryBytes(privKey)
	return crypto.EncryptSymmetric(privKeyBytes, key)
}

func decryptPrivKey(encBytes []byte, passphrase string) (privKey crypto.PrivKey, err error) {
	key, err := bcrypt.GenerateFromPassword([]byte(passphrase), 12) // TODO parameterize.  12 is good today (2016)
	if err != nil {
		Exit("Error generating bcrypt key from passphrase: " + err.Error())
	}
	key = crypto.Sha256(key) // Get 32 bytes
	privKeyBytes, err := crypto.DecryptSymmetric(encBytes, key)
	if err != nil {
		return nil, err
	}
	err = wire.ReadBinaryBytes(privKeyBytes, &privKey)
	return privKey, err
}
