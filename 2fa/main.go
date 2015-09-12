package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"github.com/dgryski/dgoogauth"
	"os"
	"strings"
)

func main() {

	// Get random secret
	sec := make([]byte, 6)
	_, err := rand.Read(sec)
	if err != nil {
		fmt.Printf("Error creating random secret key: %s", err)
	}

	// Encode secret to base32 string
	secret := base32.StdEncoding.EncodeToString(sec)

	// Give the user instructions
	fmt.Printf("Please enter token for secret: %s\n", secret)
	fmt.Printf("Hint - you can probably run: oathtool --totp -b %s\n", secret)

	// Read users' input, and strip newlines, etc.
	reader := bufio.NewReader(os.Stdin)
	token, _ := reader.ReadString('\n')
	token = strings.TrimSpace(token)

	// Configure the authentication checker.
	otpc := &dgoogauth.OTPConfig{
		Secret:      secret,
		WindowSize:  3,
		HotpCounter: 0,
	}

	// Validate the submitted token
	val, err := otpc.Authenticate(token)
	if err != nil {
		fmt.Printf("Error authenticating token: %s\n", err)
	}

	// Did it work?
	if val {
		fmt.Printf("Access granted - token was correct.\n")
	} else {
		fmt.Printf("Access denied - token was invalid.\n")
	}

}
