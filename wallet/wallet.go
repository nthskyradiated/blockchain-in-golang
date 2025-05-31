package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"

	// "fmt"
	"log"

	"github.com/nthskyradiated/blockchain-in-golang/utils"
	"golang.org/x/crypto/ripemd160"
)

const (
	checksumLength = 4
	version = byte(0x00)
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey []byte
}

// * Added for the modified saveFile and loadFile methods
type SerializableWallet struct {
    PrivateKeyD []byte
    PublicKey   []byte
}

func (w Wallet) Address() []byte {
	publicKeyHash := PublicKeyHash(w.PublicKey)
	versionedPayload := append([]byte{version}, publicKeyHash...)
	checksum := Checksum(versionedPayload)
	fullPayload := append(versionedPayload, checksum...)
	address := utils.Base58Encode(fullPayload)
	
	// fmt.Printf("Pub Key: %x\n", w.PublicKey)
	// fmt.Printf("Pub hash: %x\n", publicKeyHash)
	// fmt.Printf("Address: %x\n", address)	
	
	return address
}

func NewKeyPair() (ecdsa.PrivateKey, []byte)  {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)

	if err != nil {
		log.Panic(err)
	}

	public := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, public
}

func CreateWallet() *Wallet {
	private, public := NewKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}

func PublicKeyHash(publicKey []byte) []byte {

	// Perform SHA-256 hashing
	hash := sha256.Sum256(publicKey)
	// Perform RIPEMD-160 hashing
	ripemd160Hasher := ripemd160.New()
	_, err := ripemd160Hasher.Write(hash[:])

	if err != nil {
		log.Panic(err)
	}

	publicRipMD := ripemd160Hasher.Sum(nil)
	return publicRipMD
}

func Checksum(payload []byte) []byte {

	hash := sha256.Sum256(payload)
	hash = sha256.Sum256(hash[:])

	return hash[:checksumLength]
}

func (w *Wallet) ToSerializable() SerializableWallet {
    return SerializableWallet{
        PrivateKeyD: w.PrivateKey.D.Bytes(),
        PublicKey:   w.PublicKey,
    }
}

func (sw *SerializableWallet) ToWallet() *Wallet {
    curve := elliptic.P256()
    x, y := elliptic.Unmarshal(curve, sw.PublicKey)

    priv := new(ecdsa.PrivateKey)
    priv.PublicKey.Curve = curve
    priv.PublicKey.X = x
    priv.PublicKey.Y = y
    priv.D = new(big.Int).SetBytes(sw.PrivateKeyD)

    return &Wallet{
        PrivateKey: *priv,
        PublicKey:  sw.PublicKey,
    }
}

// *Used by the modified saveFile()
// func (w *Wallet) Bytes() ([]byte, []byte) {
// 	return w.PrivateKey.D.Bytes(), w.PublicKey
// }

// // * Used by the modified loadFile()
// func (w *Wallet) LoadFromBytes(privKey, pubKey []byte) {
// 	curve := elliptic.P256()
// 	x, y := elliptic.Unmarshal(curve, pubKey)

// 	priv := new(ecdsa.PrivateKey)
// 	priv.PublicKey.Curve = curve
// 	priv.PublicKey.X = x
// 	priv.PublicKey.Y = y
// 	priv.D = new(big.Int).SetBytes(privKey)

// 	w.PrivateKey = *priv
// 	w.PublicKey = pubKey
// }