package main

import (
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

// func main() {
// 	// Generate Alice's key pair
// 	alicePrivKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	alicePrivKeyString, err := privateKeyToString(alicePrivKey)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	fmt.Println("INI PRIVATE KEY ALICE: ", alicePrivKeyString)

// 	alicePubKeyString, err := publicKeyToString(&alicePrivKey.PublicKey)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	fmt.Println("INI PUBLICK KEY ALICE: ", alicePubKeyString)

// 	// Generate Bob's key pair
// 	bobPrivKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	bobPrivKeyString, err := privateKeyToString(bobPrivKey)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	log.Println("INI PRIVATE KEY BOB: ", bobPrivKeyString)

// 	bobPubKeyString, err := publicKeyToString(&bobPrivKey.PublicKey)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	log.Println("INI PUBLIC KEY BOB: ", bobPubKeyString)

// 	// Compute Alice's shared secret
// 	aliceSharedSecret, err := ecdh(alicePrivKey, bobPubKeyString)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	// Compute Bob's shared secret
// 	bobSharedSecret, err := ecdh(bobPrivKey, alicePubKeyString)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	// Print the shared secrets
// 	fmt.Printf("Alice's Shared Secret: %x\n", aliceSharedSecret)
// 	fmt.Printf("Bob's Shared Secret: %x\n", bobSharedSecret)
// 	if fmt.Sprintf("%x", aliceSharedSecret) == fmt.Sprintf("%x", bobSharedSecret) {
// 		fmt.Println("Shared secrets match!")
// 	}
// }

// func privateKeyToString(privKey *ecdsa.PrivateKey) (string, error) {
// 	keyBytes, err := json.Marshal(privKey)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(keyBytes), nil
// }

// func publicKeyToString(pubKey *ecdsa.PublicKey) (string, error) {
// 	keyBytes, err := json.Marshal(pubKey)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(keyBytes), nil
// }

// func ecdh(privKey *ecdsa.PrivateKey, pubKeyString string) ([]byte, error) {
// 	var pubKey ecdsa.PublicKey
// 	err := json.Unmarshal([]byte(pubKeyString), &pubKey)
// 	if err != nil {
// 		return nil, err
// 	}

// 	x, _ := pubKey.Curve.ScalarMult(pubKey.X, pubKey.Y, privKey.D.Bytes())

// 	return x.Bytes(), nil
// }

func main() {
	// Generate Alice's key pair
	alicePrivKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("INI PRIVATE KEY ALICE: ", alicePrivKey)
	fmt.Println("INI PUBLICK KEY ALICE: ", alicePrivKey.PublicKey())

	// Generate Bob's key pair
	bobPrivKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Println("INI PRIVATE KEY BOB: ", bobPrivKey)
	log.Println("INI PUBLIC EY BOB: ", bobPrivKey.PublicKey())

	// Compute Alice's shared secret
	aliceSharedSecret, err := alicePrivKey.ECDH(bobPrivKey.Public().(*ecdh.PublicKey))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Compute Bob's shared secret
	bobSharedSecret, err := bobPrivKey.ECDH(alicePrivKey.Public().(*ecdh.PublicKey))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the shared secrets
	fmt.Printf("Alice's Shared Secret: %x\n", aliceSharedSecret)
	fmt.Printf("Bob's Shared Secret: %x\n", bobSharedSecret)
	if fmt.Sprintf("%x", aliceSharedSecret) == fmt.Sprintf("%x", bobSharedSecret) {
		fmt.Println("Shared secrets match!")
	}
}

func GenerateECWithP256(ApiOrEnc string) ([]byte, error) {

	// Generate ECDSA private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	// Encode private key to PEM format
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	privateKeyPEM := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	// Encode public key to PEM format
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}
	publicKeyPEM := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	// Encode PEM blocks to []byte
	privateKeyPEMBytes := pem.EncodeToMemory(privateKeyPEM)
	publicKeyPEMBytes := pem.EncodeToMemory(publicKeyPEM)

	err = os.WriteFile(ApiOrEnc, privateKeyPEMBytes, 0644)
	if err != nil {
		return nil, err
	}

	return publicKeyPEMBytes, nil
}

// func main() {
// 	publicKeyForAPI, err := GenerateECWithP256("API")
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}

// 	publicKeyForEnc, err := GenerateECWithP256("Enc")
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}

// 	fmt.Println("key for API: ", string(publicKeyForAPI))
// 	fmt.Println("key for Enc: ", string(publicKeyForEnc))

// 	privateKey, err := os.ReadFile("API")
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}

// 	fmt.Println("private key: ", string(privateKey))
// }
