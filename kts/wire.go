package kts

func init() {
	wire.RegisterInterface((*DataHandler)(nil))

	wire.Register("m", &M{})
	wire.Register("tx", &Transaction{})
	wire.Register("baseHandler", &BaseDataHandler{})
}
