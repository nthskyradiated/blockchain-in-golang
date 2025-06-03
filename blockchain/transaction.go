package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/nthskyradiated/blockchain-in-golang/utils"
	"github.com/nthskyradiated/blockchain-in-golang/wallet"
	"log"
	"math/big"
	"strings"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

func (tx Transaction) Serialize() []byte {
	return utils.Serialize(tx)
}

func (tx *Transaction) Hash() []byte {
	var hash [32]byte
	txCopy := *tx
	txCopy.ID = nil // Clear the ID to avoid using it in the hash
	hash = sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].OutIndex == -1
}

// func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
// 	if tx.IsCoinbase() {
// 		return
// 	}

// 	for _, input := range tx.Inputs {
// 		if prevTXs[hex.EncodeToString(input.ID)].ID == nil {
// 			log.Panicf("Previous transaction not found: %x", input.ID)
// 		}
// 	}
// 	txCopy := tx.TrimmedCopy()

// 	for i, input := range txCopy.Inputs {
// 		prevTX := prevTXs[hex.EncodeToString(input.ID)]
// 		txCopy.Inputs[i].Sig = nil // Clear the signature for signing
// 		txCopy.Inputs[i].PubKey = prevTX.Outputs[input.OutIndex].ScriptPubKey
// 		txCopy.ID = txCopy.Hash()

// 		dataToSign := fmt.Sprintf("%x\n", txCopy)

// 		r, s, err := ecdsa.Sign(rand.Reader, &privKey, []byte(dataToSign))

// 		log.Printf("Signing Transaction: ID=%x", tx.ID)
// 		log.Printf("Public Key: %x", privKey.PublicKey.X.Bytes())
// 		log.Printf("Signature: R=%x, S=%x", r.Bytes(), s.Bytes())
// 		utils.HandleError(err)
// 		signature := append(r.Bytes(), s.Bytes()...)
// 		tx.Inputs[i].Sig = signature
// 		txCopy.Inputs[i].PubKey = nil
// 	}
// }

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}

	for _, input := range tx.Inputs {
		if prevTXs[hex.EncodeToString(input.ID)].ID == nil {
			log.Panicf("Previous transaction not found: %x", input.ID)
		}
	}

	txCopy := tx.TrimmedCopy()

	for i, input := range txCopy.Inputs {
		prevTX := prevTXs[hex.EncodeToString(input.ID)]
		txCopy.Inputs[i].Sig = nil // Clear the signature for signing
		txCopy.Inputs[i].PubKey = prevTX.Outputs[input.OutIndex].ScriptPubKey
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[i].PubKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		utils.HandleError(err)

		// Combine R and S into a single signature
		signature := append(r.Bytes(), s.Bytes()...)

		// Populate the transaction input fields
		tx.Inputs[i].Sig = signature
		tx.Inputs[i].PubKey = append(privKey.PublicKey.X.Bytes(), privKey.PublicKey.Y.Bytes()...)

		log.Printf("Signing Transaction: ID=%x", tx.ID)
		log.Printf("Public Key: %x", tx.Inputs[i].PubKey)
		log.Printf("Signature: R=%x, S=%x", r.Bytes(), s.Bytes())
	}
}

func (tx *Transaction) TrimmedCopy() Transaction {
	var outputs []TxOutput
	var inputs []TxInput
	for _, in := range tx.Inputs {
		inputs = append(inputs, TxInput{in.ID, in.OutIndex, nil, nil})
	}
	for _, out := range tx.Outputs {
		outputs = append(outputs, TxOutput{out.Value, out.ScriptPubKey})
	}
	txCopy := Transaction{tx.ID, inputs, outputs}
	return txCopy
}

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	for _, input := range tx.Inputs {
		if prevTXs[hex.EncodeToString(input.ID)].ID == nil {
			log.Panicf("Previous transaction not found: %x", input.ID)
		}
	}

	txCopy := tx.TrimmedCopy()
	fmt.Println(txCopy)
	curve := elliptic.P256()
	for i, input := range tx.Inputs {
		prevTX := prevTXs[hex.EncodeToString(input.ID)]
		txCopy.Inputs[i].Sig = nil // Clear the signature for verification
		txCopy.Inputs[i].PubKey = prevTX.Outputs[input.OutIndex].ScriptPubKey
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[i].PubKey = nil
		log.Printf("Transaction Hash: %x", txCopy.ID)
		r := new(big.Int).SetBytes(input.Sig[:len(input.Sig)/2])
		s := new(big.Int).SetBytes(input.Sig[len(input.Sig)/2:])
		log.Printf("Signature Components: R=%x, S=%x", r, s)

		x := new(big.Int).SetBytes(input.PubKey[:len(input.PubKey)/2])
		y := new(big.Int).SetBytes(input.PubKey[len(input.PubKey)/2:])

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: x, Y: y}
		log.Printf("Reconstructed Public Key: X=%x, Y=%x", x, y)
		if !ecdsa.Verify(&rawPubKey, txCopy.ID, r, s) {
			return false
		}
	}
	return true
}
func NewTransaction(w *wallet.Wallet, to string, amount int, UTXO *UTXOSet) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	pubKeyHash := wallet.PublicKeyHash(w.PublicKey)

	acc, validOutputs := UTXO.FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		log.Panicf("Not enough funds: %d < %d", acc, amount)
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		utils.HandleError(err)
		for _, out := range outs {
			input := TxInput{txID, out, nil, w.PublicKey}
			inputs = append(inputs, input)
		}
	}
	from := string(w.Address())
	outputs = append(outputs, *NewTXOutput(amount, to))
	if acc > amount {
		outputs = append(outputs, *NewTXOutput(acc-amount, from))
	}
	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	UTXO.Blockchain.SignTransaction(&tx, w.PrivateKey)
	return &tx
}
func CoinbaseTx(to, data string) *Transaction {
	if data == "" {
		randData := make([]byte, 24)
		_, err := rand.Read(randData)
		utils.HandleError(err)
		data = fmt.Sprintf("%x", randData)
	}
	txin := TxInput{[]byte{}, -1, nil, []byte(data)}
	txout := NewTXOutput(100, to)

	tx := Transaction{nil, []TxInput{txin}, []TxOutput{*txout}}
	tx.ID = tx.Hash()
	return &tx
}

func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.ID))
	for i, input := range tx.Inputs {
		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:     %x", input.ID))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.OutIndex))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Sig))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))
	}

	for i, output := range tx.Outputs {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.ScriptPubKey))
	}

	return strings.Join(lines, "\n")
}
