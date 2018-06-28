package mint

import (
	"github.com/tendermint/go-crypto"
	"encoding/hex"
	"errors"
	"github.com/kooksee/usmint/types/consts"
	"github.com/kooksee/kdb"
	"github.com/tendermint/abci/types"
	"bytes"
	"fmt"
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

	_, err := v.GetPubkey()
	if err != nil {
		return err
	}
	return nil
}

// UpdateValidators 更新Validators
func (v *Validator) UpdateValidator(val *types.Validator) error {
	key := []byte(v.name + hex.EncodeToString(val.PubKey))

	// power等于-1的时候,开放节点的权限
	if val.Power == -1 {
		value := bytes.NewBuffer(make([]byte, 0))
		if err := types.WriteMessage(val, value); err != nil {
			return errors.New(fmt.Sprintf("Error encoding validator: %v", err))
		}

		v.db.Set(key, value.Bytes())

		logger.Info("save node ok", "key", key)

		val.Power = 0
		return nil
	}

	// power等于-2的时候,删除节点
	if val.Power == -2 {
		v.db.Del(key)
		logger.Info("delete node ok", "key", key)

		val.Power = 0
		return nil
	}

	// power小于等于0的时候,删除验证节点
	if v.Power >= 0 {
		value := bytes.NewBuffer(make([]byte, 0))
		if err := types.WriteMessage(val, value); err != nil {
			return errors.New(fmt.Sprintf("Error encoding validator: %v", err))
		}

		v.db.Set(key, value.Bytes())

		logger.Info("save node ok", "key", key)
	}
	return nil
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
