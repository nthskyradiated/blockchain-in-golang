package wallet

import (
	"bytes"
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

func ValidateAddress(address string) bool {

	pubKeyHash := utils.Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-checksumLength:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-checksumLength]
	targetChecksum := Checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Equal(actualChecksum, targetChecksum)
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

	public := elliptic.Marshal(curve, private.PublicKey.X, private.PublicKey.Y)
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

func (w *Wallet) ToProtobuf() *SerializableWallet {
    return &SerializableWallet{
        PrivateKeyD: w.PrivateKey.D.Bytes(),
        PublicKey:   w.PublicKey,
    }
}

func FromProtobuf(sw *SerializableWallet) *Wallet {
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