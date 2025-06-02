package wallet

import (
	"fmt"
	"log"
	"os"

	"google.golang.org/protobuf/proto"
)

const walletFile = "./tmp/wallets.data"

type Wallets struct {
	Wallets map[string]*Wallet
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

func (ws *Wallets) SaveFile(nodeId string) {
	walletFile := fmt.Sprintf(walletFile, nodeId)
    serialized := &SerializableWallets{
        Wallets: make(map[string]*SerializableWallet),
    }

    for addr, wallet := range ws.Wallets {
        serialized.Wallets[addr] = wallet.ToProtobuf()
    }

    data, err := proto.Marshal(serialized)
    if err != nil {
        log.Panic(err)
    }

    err = os.WriteFile(walletFile, data, 0644)
    if err != nil {
        log.Panic(err)
    }
}

func (ws *Wallets) LoadFile(nodeId string) error {
	walletFile := fmt.Sprintf(walletFile, nodeId)
    if _, err := os.Stat(walletFile); os.IsNotExist(err) {
        return err
    }

    data, err := os.ReadFile(walletFile)
    if err != nil {
        return err
    }

    serialized := &SerializableWallets{}
    err = proto.Unmarshal(data, serialized)
    if err != nil {
        return err
    }

    wallets := make(map[string]*Wallet)
    for addr, sw := range serialized.Wallets {
        wallets[addr] = FromProtobuf(sw)
    }

    ws.Wallets = wallets
    return nil
}
