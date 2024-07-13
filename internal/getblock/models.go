package getblock

type BlockNumberResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  string `json:"result"`
}

type BlockResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  Block  `json:"result"`
}

type Block struct {
	Number       string        `json:"number"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Value    string `json:"value"`
	GasPrice string `json:"gasPrice"`
}
