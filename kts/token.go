package kts

import (
	"math/big"
	"github.com/ethereum/go-ethereum/common"
	"github.com/kooksee/kdb"
	"github.com/kooksee/usmint/kts/consts"
	"github.com/kooksee/usmint/cmn"
	"errors"
)

func NewToken(address common.Address, db kdb.IKDB) *Token {
	return &Token{address: address, name: consts.TokenPrefix, db: db.KHash([]byte(consts.TokenPrefix))}
}

// 设计token机制
// 参考erc20标准
type Token struct {
	IToken
	address common.Address
	name    string
	db      kdb.IKHash
}

// InitToken 初始化token
func (t *Token) InitToken(amount *big.Int) error {
	// 只在启动的时候初始化一次,在init chain的时候调用,然后把token放到一个总地址下面

	b, err := t.db.Exist([]byte(consts.TokenAddress))
	if err != nil {
		return cmn.ErrPipe("Token.InitToken",err)
	}
	if b {
		return errors.New("")
	}

	return cmn.ErrPipe(
		"Token InitToken Error",
		cmn.ErrCurry(t.db.Set, []byte(consts.TotalSupply), amount.Bytes()),
		cmn.ErrCurry(t.db.Set, []byte(consts.TokenAddress), amount.Bytes()),
	)
}

// TransferTo 转账
func (t *Token) TransferTo(toAddress common.Address, value *big.Int) error {

	fromAmount, err := t.db.Get(t.address.Bytes())
	if err != nil {
		return err
	}
	fa := big.NewInt(0).SetBytes(fromAmount)

	toAmount, err := t.db.Get(toAddress.Bytes())
	if err != nil {
		return err
	}
	ta := big.NewInt(0).SetBytes(toAmount)

	if fa.Cmp(value) < 0 {
		return errors.New("账户余额不足")
	}

	fa = fa.And(fa, value)
	ta = ta.Sub(ta, value)

	return cmn.ErrPipe(
		"TransferTo Error",
		t.db.Set(t.address.Bytes(), fa.Bytes()),
		t.db.Set(toAddress.Bytes(), ta.Bytes()),
	)
}

func (t *Token) TransferFrom(fromAddress common.Address, toAddress common.Address, value *big.Int) error {
	fromAddr := t.address.Bytes()
	approveAddr := append(fromAddress.Bytes(), t.address.Bytes()...)
	toAddr := toAddress.Bytes()

	fd, err := t.db.Get(fromAddr)
	if err != nil {
		return err
	}
	fa := big.NewInt(0).SetBytes(fd)

	ad, err := t.db.Get(approveAddr)
	if err != nil {
		return err
	}
	aa := big.NewInt(0).SetBytes(ad)

	td, err := t.db.Get(toAddr)
	if err != nil {
		return err
	}
	ta := big.NewInt(0).SetBytes(td)

	if fa.Cmp(value) < 0 || aa.Cmp(value) < 0 {
		return errors.New("账户余额不足")
	}

	aa = aa.Sub(aa, value)
	ta = ta.Add(ta, value)

	return cmn.ErrPipe(
		"Token TransferFrom Error",
		t.db.Set(approveAddr, aa.Bytes()),
		t.db.Set(toAddr, ta.Bytes()),
	)
}

func (t *Token) BalanceOf(address common.Address) *big.Int {
	d, err := t.db.Get(address.Bytes())
	if err != nil {
		return big.NewInt(0)
	}
	if len(d) == 0 {
		return big.NewInt(0)
	}

	return big.NewInt(0).SetBytes(d)
}

func (t *Token) TotalSupply() *big.Int {
	d, err := t.db.Get([]byte(consts.TotalSupply))
	if err != nil {
		return big.NewInt(0)
	}
	if len(d) == 0 {
		return big.NewInt(0)
	}

	return big.NewInt(0).SetBytes(d)
}

func (t *Token) Approve(spenderAddress common.Address, value *big.Int) error {
	fromAddr := t.address
	toAddr := append(t.address.Bytes(), spenderAddress.Bytes()...)

	fd, err := t.db.Get(fromAddr.Bytes())
	if err != nil {
		return err
	}
	fa := big.NewInt(0).SetBytes(fd)

	td, err := t.db.Get(toAddr)
	if err != nil {
		return err
	}
	ta := big.NewInt(0).SetBytes(td)

	if fa.Cmp(value) < 0 {
		return errors.New("账户余额不足")
	}

	ta = ta.And(ta, value)

	return cmn.ErrPipe(
		"Token Approve Error",
		t.db.Set(toAddr, ta.Bytes()),
	)
}

func (t *Token) Allowance(ownerAddress common.Address, spenderAddress common.Address) *big.Int {
	addr := append(ownerAddress.Bytes(), spenderAddress.Bytes()...)
	ad, err := t.db.Get(addr)
	if err != nil {
		cmn.ErrPipe("Token Allowance Error", err)
		return big.NewInt(0)
	}
	return big.NewInt(0).SetBytes(ad)
}

func (t *Token) Balances(adds chan common.Address) error {
	return t.db.Range(func(key, _ []byte) error {
		if len(key) == len(t.address.Bytes()) {
			adds <- common.BytesToAddress(key)
		}
		return nil
	})
}

type IToken interface {
	//初始化token
	InitToken(amount *big.Int) error

	//转账
	TransferTo(toAddress common.Address, value *big.Int) error

	//从授权地址转账给其他地址
	TransferFrom(fromAddress common.Address, toAddress common.Address, value *big.Int) error

	//代币总量
	TotalSupply() *big.Int

	//获得所有持有代币的用户的地址
	Balances(adds chan common.Address) error

	//查询账户余额
	BalanceOf(address common.Address) *big.Int

	//授权给某个账户
	Approve(spenderAddress common.Address, value *big.Int) error

	//查询授权余额
	Allowance(ownerAddress common.Address, spenderAddress common.Address) *big.Int
}
