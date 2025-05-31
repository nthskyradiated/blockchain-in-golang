package wallet

import (
	// "crypto/elliptic"
	"log"
	"os"
	"google.golang.org/protobuf/proto"
)

const walletFile = "./tmp/wallets.data"

type Wallets struct {
	Wallets map[string]*Wallet
}

func NewWallets() (*Wallets, error) {
	ws := Wallets{}
	ws.Wallets = make(map[string]*Wallet)
	err := ws.LoadFile()
	return &ws, err
}

func (ws *Wallets) AddWallet() string {
	wallet := CreateWallet()
	address := string(wallet.Address())
	ws.Wallets[address] = wallet
	ws.SaveFile()
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

// ! doesn't work. kept for reference
// func (ws *Wallets) LoadFile() error {

// 	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
// 		return err
// 	}

// 	var wallets Wallets
// 	fileContent, err := os.ReadFile(walletFile)

// 	if err != nil {
// 		return err
// 	}

// 	gob.Register(elliptic.P256())
// 	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
// 	err = decoder.Decode(&wallets)

// 	if err != nil {
// 		return err
// 	}

// 	ws.Wallets = wallets.Wallets
// 	return nil
// }

func (ws *Wallets) SaveFile() {
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

func (ws *Wallets) LoadFile() error {
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

// ! doesn't work. kept for reference
// func (ws *Wallets) SaveFile() {
// 	var content bytes.Buffer

// 	gob.Register(elliptic.P256())
// 	encoder := gob.NewEncoder(&content)
// 	err := encoder.Encode(ws)

// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	err = os.WriteFile(walletFile, content.Bytes(), 0644)

// 	if err != nil {
// 		log.Panic(err)
// 	}
// }
