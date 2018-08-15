package mint

import (
	"github.com/kooksee/usmint/kts/consts"
	"github.com/kooksee/kdb"
	"github.com/kooksee/usmint/cmn"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/kooksee/usmint/kts"
	"errors"
	"encoding/hex"
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
func (v *ValidatorManager) CheckValidator(val *kts.Validator) error {

	if val.Power > 9 || val.Power < -2 {
		return errors.New("power值超过限制")
	}
	return nil
}

// UpdateValidators 更新Validators
func (v *ValidatorManager) UpdateValidator(val *kts.Validator) error {
	logger.Info("update node", "node", val.Address)

	addr, _ := hex.DecodeString(val.Address)

	// power等于-1的时候,开放节点的权限
	if val.Power == -1 {
		val.Power = 0
		val1, err := cmn.JsonMarshal(val)
		return cmn.ErrPipe("ValidatorManager UpdateValidator 1", err, cmn.ErrCurry(v.db.Set, addr, val1))
	}

	// power等于-2的时候,删除节点
	if val.Power == -2 {
		val.Power = 0
		return cmn.ErrPipe("ValidatorManager UpdateValidator 2", cmn.ErrCurry(v.db.Del, addr))
	}

	val1, err := cmn.JsonMarshal(val)
	return cmn.ErrPipe("ValidatorManager UpdateValidator 3", err, cmn.ErrCurry(v.db.Set, addr, val1))
}

func (v *ValidatorManager) Has(addr []byte) bool {
	e, err := v.db.Exist(addr)
	if err != nil {
		logger.Error("ValidatorManager Has error", "err", err.Error())
		return false
	}
	return e
}
