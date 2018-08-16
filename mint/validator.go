package mint

import (
	"github.com/kooksee/usmint/kts/consts"
	"github.com/kooksee/kdb"
	"github.com/kooksee/usmint/cmn"
	"github.com/tendermint/tendermint/crypto"
	"github.com/kooksee/usmint/kts"
	"encoding/hex"
	"github.com/tendermint/tendermint/abci/types"
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

// Check 检查Power值和Pubkey
func (v *ValidatorManager) CheckValidatorWithTx(tx *kts.Transaction) error {
	val, err := kts.DecodeValidator(tx.Data)
	if err != nil {
		return cmn.ErrPipe("ValidatorManager CheckValidatorWithTx", err)
	}

	if val.Power > 9 || val.Power < -2 {
		return cmn.Err("ValidatorManager CheckValidatorWithTx: power值超过限制")
	}
	
	return nil
}

func (v *ValidatorManager) UpdateValidatorWithTx(tx *kts.Transaction) (types.Validator, error) {
	val, err := kts.DecodeValidator(tx.Data)
	return types.Validator{
		Address: val.GetAddress(),
		PubKey:  val.PubKey,
		Power:   val.Power,
	}, cmn.ErrPipe("ValidatorManager UpdateValidatorWithTx", err, cmn.ErrCurry(v.UpdateValidator, val))
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
