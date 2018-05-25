package tx

import (
	"github.com/kooksee/kchain/types/cnst"
)


// Db
type Db struct {
	Key   string `json:"key,omitempty" mapstructure:"key"`
	Value string `json:"value,omitempty" mapstructure:"value"`
}

func NewDb() *Db {
	return &Db{}
}

func (db *Db)ToSortBytes() []byte {
	return []byte(db.Key + db.Value)
}
func (db *Db)ToBytes() []byte {
	d, _ := json.Marshal(db)
	return d
}
func (db *Db) FromBytes(d []byte) error {
	return json.Unmarshal(d, db)
}
func (db *Db)ToKv() ([]byte, []byte) {
	return db.K(), db.ToBytes()
}
func (db *Db)K() []byte {
	return []byte(db.GetPrefix() + db.Key)
}
func (db *Db)GetPrefix() string {
	return cnst.DbPrefix
}
