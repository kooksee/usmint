package mint

import (
	"github.com/tendermint/go-crypto"
	"encoding/hex"
	"errors"
	"github.com/kooksee/usmint/types/consts"
	"github.com/kooksee/kdb"
	"github.com/tendermint/abci/types"
	"github.com/kooksee/usmint/cmn"
)

func NewValidator(dbs ... *kdb.KDB) *Validator {
	db1 := db
	if len(dbs) > 0 {
		db1 = dbs[0]
	}

	name := consts.ValidatorPrefix
	return &Validator{name: name, db: db1.KHash([]byte(name))}
}

// Validator
type Validator struct {
	db     *kdb.KHash
	name   string
	pubkey crypto.PubKey

	PubKey string `json:"pubkey,omitempty" mapstructure:"pubkey"`
	Power  int64  `json:"power,omitempty" mapstructure:"power"`
}

func (v *Validator) GetPubkey() (crypto.PubKey, error) {
	if v.pubkey != nil {
		return v.pubkey, nil
	}

	d, err := hex.DecodeString(v.PubKey)
	if err != nil {
		return nil, err
	}

	pk, err := crypto.PubKeyFromBytes(d)
	if err != nil {
		return nil, err
	}

	v.pubkey = pk
	return pk, nil
}

func (v *Validator) Has() bool {
	pk, err := v.GetPubkey()
	if err != nil {
		return false
	}
	b, err := v.db.Exist(pk.Address().Bytes())
	cmn.ErrPipeLog("Validator Has error", err)
	return b
}

// Check 检查Power值和Pubkey
func (v *Validator) Check() error {
	if v.Power > 9 {
		return errors.New("Power值超过限制")
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

	// power等于-1的时候,开放节点的权限
	if val.Power == -1 {
		val.Power = 0

		val1, err := json.Marshal(val)
		if err != nil {
			return err
		}

		pk, err := crypto.PubKeyFromBytes(val.GetPubKey())
		if err != nil {
			return err
		}

		logger.Info("save node ok", "key", string(val.PubKey))
		return v.db.Set(pk.Address().Bytes(), val1)
	}

	// power等于-2的时候,删除节点
	if val.Power == -2 {
		val.Power = 0

		pk, err := crypto.PubKeyFromBytes(val.GetPubKey())
		if err != nil {
			return err
		}

		logger.Info("delete node ok", "key", string(val.PubKey))

		return v.db.Del(pk.Address().Bytes())
	}

	if v.Power > 9 || v.Power < -2 {
		return errors.New("power值超过限制")
	}

	val1, err := json.Marshal(val)
	if err != nil {
		return err
	}

	pk, err := crypto.PubKeyFromBytes(val.GetPubKey())
	if err != nil {
		return err
	}

	logger.Info("save node ok", "key", string(val.PubKey))
	return v.db.Set(pk.Address().Bytes(), val1)
}

func (v *Validator) Decode(val []byte) error {
	return json.Unmarshal(val, v)
}

func (v *Validator) Delete() error {
	pk, err := v.GetPubkey()
	if err != nil {
		return err
	}
	return v.db.Del(pk.Address().Bytes())
}

func (v *Validator) Save() error {
	pk, err := v.GetPubkey()
	if err != nil {
		return err
	}

	val, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return v.db.Set(pk.Address().Bytes(), val)
}
