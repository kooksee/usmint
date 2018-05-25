package web

import (
	"github.com/labstack/echo"
	"net/http"
	"github.com/kooksee/kchain/types"
	// cmn "github.com/tendermint/tmlibs/common"

	"fmt"
)

func txPost1(c echo.Context) error {

	a := <-cfg().Node.Switch().Broadcast(0, []byte("hello"))

	cfg().Node.Switch().Broadcast(0x60, []byte("hello123456"))

	for _, p := range cfg().Node.Switch().Peers().List() {
		d, _ := json.Marshal(p.NodeInfo().String())
		fmt.Println(string(d))
		fmt.Println(p.String())
		fmt.Println(p.NodeInfo().Channels)
	}

	c.Logger().Error(cfg().Node.Switch().String())
	c.Logger().Error(cfg().Node.Switch().NodeInfo().String())
	c.Logger().Error(cfg().Node.Switch().NumPeers())
	c.Logger().Error(cfg().Node.Switch().Reactors())

	for _, p := range cfg().Node.Switch().Peers().List() {
		c.Logger().Error(p.String())
		c.Logger().Error(p.NodeInfo().String())
	}

	c.Logger().Error(a)

	return c.JSON(http.StatusOK, m{
		"ok": "ok",
		"oo": m{
			"rid": c.Response().Header().Get(echo.HeaderXRequestID),
		},
	})
	return nil

}

func txPost(c echo.Context) error {
	tx := types.NewTransaction()

	if err := c.Bind(tx); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(tx); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	//abci := kcfg.Abci()
	////cmtRes, err := abci.BroadcastTxCommit(ctypes.Tx([]byte("")))
	//if err != nil {
	//	c.Logger().Error(err.Error())
	//	return c.JSON(http.StatusBadRequest, err.Error())
	//}

	// 2. 把txid发送tx
	// 1. 得到txID

	return c.JSON(http.StatusOK, m{
		"ok": "ok",
		"oo": m{
			"rid": c.Response().Header().Get(echo.HeaderXRequestID),
		},
		"d": tx,
	})
}

func txGet(c echo.Context) error {
	txId := c.Param("txId")
	return c.JSON(http.StatusOK, txId)
}
