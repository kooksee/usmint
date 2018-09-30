package db

type MintDb struct {
	data map[string][]byte
}

// Commit current transaction cache to block cache
func (t *MintDb) Set(key, value []byte) {
	t.data[string(key)] = value
}

func (t *MintDb) Get(key []byte) []byte {
	if t.data[string(key)] == nil {
		t.data[string(key)] = appDb.Get(key)
	}
	return t.data[string(key)]
}

func (t *MintDb) Exist(key []byte) bool {
	return t.Get(key) != nil
}

func (t *MintDb) Delete(key []byte) bool {
	_, b := t.data[string(key)]
	return b
}

func (t *MintDb) Commit() {
	for k, v := range t.data {
		appDb.Set([]byte(k), v)
	}
}

func (t *MintDb) RollBack() {
	t.data = make(map[string][]byte)
}
