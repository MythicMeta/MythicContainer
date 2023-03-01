package translationstructs

type CryptoKeys struct {
	EncKey   *[]byte `json:"enc_key,omitempty"`
	DecKey   *[]byte `json:"dec_key,omitempty"`
	Value    string  `json:"value"`
	Location string  `json:"location,omitempty"`
}
