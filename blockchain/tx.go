package blockchain

type TXInput struct {
	ID       []byte
	OutIndex int
	Sig      string
}

type TXOutput struct {
	Value        int
	ScriptPubKey string
}

func (in *TXInput) CanUnlock(data string) bool {
	return in.Sig == data
}

func (out *TXOutput) CanBeUnlocked(data string) bool {
	return out.ScriptPubKey == data
}