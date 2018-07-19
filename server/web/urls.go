package web

import "github.com/labstack/echo"

func initUrls(e *echo.Echo) {
	e.GET("/", index)
	e.GET("/abci_info", abciInfo)
	e.GET("/dump_consensus_state", dump_consensus_state)
	e.GET("/genesis", genesis)
	e.GET("/net_info", net_info)
	e.GET("/num_unconfirmed_txs", num_unconfirmed_txs)
	e.GET("/status", status)
	e.GET("/unconfirmed_txs", unconfirmed_txs)
	e.GET("/block", block)
	e.GET("/blockchain", blockchain)
	e.GET("/broadcast_tx", broadcast_tx_async)
	e.GET("/tx/:id", tx)
}
