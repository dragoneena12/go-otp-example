package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"image/png"
)

func main() {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Example.com",       // OTPを発行する組織・企業名。必須。
		AccountName: "alice@example.com", // OTP発行先のアカウント名。必須。
		Period:      30,                  // TOTPハッシュが有効な時間（秒）。デフォルトでは30秒。
		SecretSize:  20,                  // 生成されるSecretのバイト長。デフォルトでは20バイト。
		Secret:      []byte{},            // Secretとして用いるバイト列。デフォルトではempty。
		Digits:      otp.DigitsSix,       // 生成されるTOTPハッシュの桁数。デフォルトでは6桁。
		Algorithm:   otp.AlgorithmSHA1,   // HMACに用いるハッシュアルゴリズム。デフォルトではSHA1。
		Rand:        rand.Reader,         // Secret生成に用いる乱数生成io.Reader。デフォルトではrand.Reader。
	})
	if err != nil {
		log.Fatal(err)
	}

	// Convert TOTP key into a QR code encoded as a PNG image.
	img, err := key.Image(200, 200)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create("qrcode.png")
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	// Now Validate that the user's successfully added the passcode.
	var passcode string
	fmt.Println("input one time password:")
	fmt.Scanf("%s", &passcode)
	valid := totp.Validate(passcode, key.Secret())

	if valid {
		// User successfully used their TOTP, save it to your backend!
		fmt.Println("Authorization succeed!")
		err := ioutil.WriteFile("secret.key", []byte(key.Secret()), 0600)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("Access denied")
	}
}
