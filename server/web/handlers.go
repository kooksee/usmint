package web

import (
	"github.com/labstack/echo"
	"net/http"
	"github.com/kooksee/usmint/server/core"
	"math/big"
	tp "github.com/tendermint/tendermint/types"
	"encoding/hex"
)

func index(c echo.Context) error {

	c.String(
		http.StatusOK,
		`ok`,
	)
	return nil
}

func abciInfo(c echo.Context) error {

	info, err := core.ABCIInfo()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, info)
}

func dump_consensus_state(c echo.Context) error {

	info, err := core.DumpConsensusState()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, info)
}

func genesis(c echo.Context) error {

	info, err := core.Genesis()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, info)
}

func net_info(c echo.Context) error {

	info, err := core.NetInfo()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, info)
}

func num_unconfirmed_txs(c echo.Context) error {

	info, err := core.NumUnconfirmedTxs()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, info)
}

func status(c echo.Context) error {

	info, err := core.Status()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, info)
}

func unconfirmed_txs(c echo.Context) error {

	info, err := core.UnconfirmedTxs()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, info)
}

func block(c echo.Context) error {
	a := big.NewInt(0)
	a.SetString(c.QueryParam("height"), 10)
	ret := a.Int64()
	info, err := core.Block(&ret)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, info)
}

func blockchain(c echo.Context) error {
	minh := big.NewInt(0)
	minh.SetString(c.QueryParam("height"), 10)
	min := minh.Int64()

	maxh := big.NewInt(0)
	maxh.SetString(c.QueryParam("height"), 10)
	max := maxh.Int64()

	info, err := core.BlockchainInfo(min, max)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, info)
}

func broadcast_tx_async(c echo.Context) error {

	ret, _ := hex.DecodeString(c.QueryParam("tx"))
	info, err := core.BroadcastTxAsync(tp.Tx(ret))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, info)
}

func tx(c echo.Context) error {

	ret, _ := hex.DecodeString(c.QueryParam("tx"))
	info, err := core.Tx(ret, false)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, info)
}
