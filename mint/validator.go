package mint

import (
	"github.com/kooksee/usmint/kts/consts"
	"github.com/kooksee/kdb"
	"github.com/kooksee/usmint/cmn"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/kooksee/usmint/kts"
	"errors"
)

func NewValidatorManager() *ValidatorManager {
	name := consts.ValidatorPrefix
	return &ValidatorManager{name: name, db: db.KHash([]byte(name))}
}

// Validator
type ValidatorManager struct {
	kts.IValidatorManager
	db     kdb.IKHash
	name   string
	pubkey crypto.PubKey
}

func (v *ValidatorManager) GetPubkey(pubk []byte) (crypto.PubKey, error) {
	pk, err := cryptoAmino.PubKeyFromBytes(pubk)
	if err != nil {
		return nil, cmn.ErrPipe("ValidatorManager GetPubkey error", err)
	}
	return pk, nil
}

// Check 检查Power值和Pubkey
func (v *ValidatorManager) CheckValidator(val *types.Validator) error {

	if val.Power > 9 || val.Power < -2 {
		return errors.New("power值超过限制")
	}

	// 检查pubkey的格式是否合格
	_, err := v.GetPubkey(val.PubKey.GetData())
	if err != nil {
		return err
	}
	return nil
}

// UpdateValidators 更新Validators
func (v *ValidatorManager) UpdateValidator(val *types.Validator) error {
	logger.Error("update node", "node", val.String())

	sss, err := val.PubKey.Marshal()
	if err != nil {
		return err
	}
	pk, err := v.GetPubkey(sss)
	if err = cmn.ErrPipe("ValidatorManager UpdateValidator error", err); err != nil {
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

func (v *ValidatorManager) DecodeValidator(val []byte) (*kts.Validator, error) {
	vt := &kts.Validator{}
	return vt, cmn.ErrPipe("validator decode error", cmn.ErrCurry(cmn.JsonUnmarshal, val, vt))
}
