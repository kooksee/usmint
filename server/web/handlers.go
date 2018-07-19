package web

import (
	"github.com/labstack/echo"
	"net/http"
	"github.com/kooksee/usmint/types"
)

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

	switch tx.Event {
	// 添加验证节点
	case "":
		//	调用交易合约函数
	case "1":
		//	调用查询合约函数
	case "2":
		//	得到block hash
	case "3":
		//	得到区块信息
		//  q abci_info
		//  q dump_consensus_state
		//  q genesis
		//  q net_info
		//  q num_unconfirmed_txs
		//  q status
		//  q unconfirmed_txs

		// ﻿abci_query?path=&data=&height=&prove=
		// ﻿block?height=_
		// ﻿block_results?height=
		//  blockchain?minHeight=&maxHeight=_
		// ﻿broadcast_tx_async?tx=_
		// ﻿broadcast_tx_commit?tx=_
		// ﻿commit?height=
		//	tx?hash=&prove=_
		// ﻿tx_search?query=&prove=
		// ﻿validators?height=_

	default:
		// 方法不存在
	}

	return c.JSON(http.StatusOK, m{
		"ok": "ok",
		"oo": m{
			"rid": c.Response().Header().Get(echo.HeaderXRequestID),
		},
		"d": tx,
	})
}

func txGet(c echo.Context) error {
	txId := c.Param("id")
	return c.JSON(http.StatusOK, txId)
}
