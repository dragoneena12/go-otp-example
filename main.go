package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/pquerna/otp/totp"

	"image/png"
)

func main() {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Example.com",
		AccountName: "alice@example.com",
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
		err := ioutil.WriteFile("hoge.txt", []byte(key.Secret()), 0664)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("Access denied")
	}
}
