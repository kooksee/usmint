package mint

import (
	"errors"
	"github.com/kooksee/usmint/kts/consts"
	"github.com/kooksee/kdb"
	"github.com/kooksee/usmint/cmn"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
)

func NewValidator(dbs ... kdb.IKDB) *Validator {
	db1 := db
	if len(dbs) > 0 {
		db1 = dbs[0]
	}

	name := consts.ValidatorPrefix
	return &Validator{name: name, db: db1.KHash([]byte(name))}
}

// Validator
type Validator struct {
	db     kdb.IKHash
	name   string
	pubkey crypto.PubKey

	PubKey []byte `json:"pubkey,omitempty" mapstructure:"pubkey"`
	Power  int64  `json:"power,omitempty" mapstructure:"power"`
}

func (v *Validator) GetPubkey() (crypto.PubKey, error) {
	if v.pubkey != nil {
		return v.pubkey, nil
	}

	pk, err := cryptoAmino.PubKeyFromBytes(v.PubKey)
	if err != nil {
		return nil, cmn.ErrPipe("validator parse pubkey", err)
	}

	v.pubkey = pk
	return pk, nil
}

func (v *Validator) Has(pk crypto.PubKey) bool {
	b, err := v.db.Exist(pk.Address().Bytes())
	cmn.ErrPipe("validator has error", err)
	return b
}

// Check 检查Power值和Pubkey
func (v *Validator) Check() error {

	if v.Power > 9 || v.Power < -2 {
		return errors.New("power值超过限制")
	}

	// 检查pubkey的格式是否合格
	_, err := v.GetPubkey()
	if err != nil {
		return err
	}
	return nil
}

// UpdateValidators 更新Validators
func (v *Validator) UpdateValidator(val *types.Validator) error {
	logger.Error("update node", "node", val.String())

	ppk := val.GetPubKey()
	pk, err := cryptoAmino.PubKeyFromBytes((&ppk).GetData())
	if err = cmn.ErrPipe("pubkey parse error", err); err != nil {
		return err
	}

	// power等于-1的时候,开放节点的权限
	if val.Power == -1 {
		val.Power = 0
		val1, err := cmn.JsonMarshal(val)
		return cmn.ErrPipe("save node error", err, cmn.ErrCurry(v.db.Set, pk.Address().Bytes(), val1))
	}

	// power等于-2的时候,删除节点
	if val.Power == -2 {
		val.Power = 0
		return cmn.ErrPipe("delete node error", cmn.ErrCurry(v.db.Del, pk.Address().Bytes()))
	}

	val1, err := cmn.JsonMarshal(val)
	return cmn.ErrPipe("save node error", err, cmn.ErrCurry(v.db.Set, pk.Address().Bytes(), val1))
}

func (v *Validator) Decode(val []byte) error {
	return cmn.ErrPipe("validator decode error", cmn.ErrCurry(cmn.JsonUnmarshal, val, v))
}

func (v *Validator) Delete() error {
	pk, err := v.GetPubkey()
	return cmn.ErrPipe("Validator.Delete error", err, cmn.ErrCurry(v.db.Del, pk.Address().Bytes()))
}

func (v *Validator) Save() error {
	pk, err := v.GetPubkey()
	if err != nil {
		return err
	}

	val, err := cmn.JsonMarshal(v)
	return cmn.ErrPipe("Validator.Save Error", err, cmn.ErrCurry(v.db.Set, pk.Address().Bytes(), val))
}
