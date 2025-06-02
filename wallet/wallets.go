package wallet

import (
	"fmt"
	"log"
	"os"
	"github.com/nthskyradiated/blockchain-in-golang/utils"
)

const walletFile = "./tmp/wallets_%s.data"

type Wallets struct {
	Wallets map[string]*Wallet
}

// * Added for the modified saveFile and loadFile methods
type SerializedWallets struct {
	Wallets map[string]SerializedWallet
}

func NewWallets(nodeId string) (*Wallets, error) {
	ws := Wallets{}
	ws.Wallets = make(map[string]*Wallet)
	err := ws.LoadFile(nodeId)
	return &ws, err
}

func (ws *Wallets) AddWallet(nodeId string) string {
	wallet := CreateWallet()
	address := string(wallet.Address())
	ws.Wallets[address] = wallet
	ws.SaveFile(nodeId)
	log.Printf("New wallet created with address: %s", address)
	return address
}

func (ws *Wallets) GetAllAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}
	return addresses
}

func (ws Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

func (ws *Wallets) LoadFile(nodeId string) error {
	walletFile := fmt.Sprintf(walletFile, nodeId)
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {	
		return err
	}

	fileContent, err := os.ReadFile(walletFile)
	if err != nil {
		return err
	}

	// var serialized SerializedWallets
	// decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	// err = decoder.Decode(&serialized)
	// if err != nil {
	// 	return err
	// }
	serialized := utils.Deserialize[SerializedWallets](fileContent)


	wallets := make(map[string]*Wallet)
	for addr, data := range serialized.Wallets {
		wallet := &Wallet{}
		wallet.LoadFromBytes(data.PrivateKey, data.PublicKey)
		wallets[addr] = wallet
	}

	ws.Wallets = wallets
	return nil
}

func (ws *Wallets) SaveFile(nodeId string) {
	walletFile := fmt.Sprintf(walletFile, nodeId)
	serialized := SerializedWallets{
		Wallets: make(map[string]SerializedWallet),
	}

	for addr, wallet := range ws.Wallets {
		privKey, pubKey := wallet.Bytes()
		serialized.Wallets[addr] = SerializedWallet{
			PrivateKey: privKey,
			PublicKey:  pubKey,
		}
	}

content := utils.Serialize(serialized)
	err := os.WriteFile(walletFile, content, 0644)
	utils.HandleError(err)
}