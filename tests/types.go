package tests

type BroadcastTxSyncResp struct {
	ID      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result struct {
		Code int    `json:"code"`
		Log  string `json:"log"`
		Data string `json:"data"`
		Hash string `json:"hash"`
	} `json:"result"`
	Error struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
		Data    string `json:"data"`
	} `json:"error"`
}

type QueryTxResp struct {
	ID      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result struct {
		Height int64 `json:"height"`
		Index  int64 `json:"index"`
		TxResult struct {
			Data string      `json:"data"`
			Tags interface{} `json:"tags"`
		} `json:"tx_result"`

		Tx string `json:"tx"`
		Proof struct {
			Index    string `json:"Index"`
			Total    string `json:"Total"`
			RootHash string `json:"RootHash"`
			Data     string `json:"Data"`
			Proof struct {
				Aunts interface{} `json:"aunts"`
			} `json:"Proof"`
		} `json:"proof"`
	} `json:"result"`
	Error struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
		Data    string `json:"data"`
	} `json:"error"`
}

type BroadcastTxCommitResp struct {
	ID      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result struct {
		CheckTx struct {
			Code int    `json:"code"`
			Log  string `json:"log"`
		} `json:"check_tx"`
		DeliverTx struct {
			Code int    `json:"code"`
			Log  string `json:"log"`
		} `json:"deliver_tx"`
		Hash   string `json:"hash"`
		Height int    `json:"height"`
	} `json:"result"`
	Error struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
		Data    string `json:"data"`
	} `json:"error"`
}

type AbciQueryResp struct {
	ID      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result struct {
		Response struct {
			Value string `json:"value"`
			Log   string `json:"log"`
			Code  int64  `json:"code"`
		} `json:"response"`
	} `json:"result"`
}

type ChainStatus struct {
	ID      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result struct {
		LatestAppHash     string `json:"latest_app_hash"`
		LatestBlockHash   string `json:"latest_block_hash"`
		LatestBlockHeight int    `json:"latest_block_height"`
		LatestBlockTime   string `json:"latest_block_time"`
		NodeInfo struct {
			Channels   string   `json:"channels"`
			ListenAddr string   `json:"listen_addr"`
			Moniker    string   `json:"moniker"`
			Network    string   `json:"network"`
			Other      []string `json:"other"`
			PubKey struct {
				Data string `json:"data"`
				Type string `json:"type"`
			} `json:"pub_key"`
			Version string `json:"version"`
		} `json:"node_info"`
		PubKey struct {
			Data string `json:"data"`
			Type string `json:"type"`
		} `json:"pub_key"`
		Syncing bool `json:"syncing"`
	} `json:"result"`
}

type BlockHashHeight struct {
	Hash   string `json:"hash,omitempty" binding:"required"`
	Height int    `json:"height,omitempty" binding:"required"`
}

type BlockChain struct {
	ID      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result struct {
		BlockMetas []struct {
			BlockID struct {
				Hash string `json:"hash"`
				Parts struct {
					Hash  string `json:"hash"`
					Total int    `json:"total"`
				} `json:"parts"`
			} `json:"block_id"`
			Header struct {
				AppHash       string `json:"app_hash"`
				ChainID       string `json:"chain_id"`
				ConsensusHash string `json:"consensus_hash"`
				DataHash      string `json:"data_hash"`
				EvidenceHash  string `json:"evidence_hash"`
				Height        int    `json:"height"`
				LastBlockID struct {
					Hash string `json:"hash"`
					Parts struct {
						Hash  string `json:"hash"`
						Total int    `json:"total"`
					} `json:"parts"`
				} `json:"last_block_id"`
				LastCommitHash  string `json:"last_commit_hash"`
				LastResultsHash string `json:"last_results_hash"`
				NumTxs          int    `json:"num_txs"`
				Time            string `json:"time"`
				TotalTxs        int    `json:"total_txs"`
				ValidatorsHash  string `json:"validators_hash"`
			} `json:"header"`
		} `json:"block_metas"`
		LastHeight int `json:"last_height"`
	} `json:"result"`
}

type NumUnconfirmedTxsResp struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result struct {
		NTxs int         `json:"n_txs"`
		Txs  interface{} `json:"txs"`
	} `json:"result"`
}

type AccountsResp struct {
	PubKey string `json:"pub_key,omitempty" binding:"required"`
	Power  string `json:"power,omitempty" binding:"required"`
}

type ClientsGroup struct {
	Signature  string   `json:"signature,omitempty" binding:"required"`
	PubKey     string   `json:"pubkey,omitempty" binding:"required"`
	SubPubKeys []string `json:"subkeys,omitempty" binding:"required"`
}

type PrivFile struct {
	Address    string `json:"address"`
	LastHeight int    `json:"last_height"`
	LastRound  int    `json:"last_round"`
	LastSignature struct {
		Data string `json:"data"`
		Type string `json:"type"`
	} `json:"last_signature"`
	LastSignbytes string `json:"last_signbytes"`
	LastStep      int    `json:"last_step"`
	PrivKey struct {
		Data string `json:"data"`
		Type string `json:"type"`
	} `json:"priv_key"`
	PubKey struct {
		Data string `json:"data"`
		Type string `json:"type"`
	} `json:"pub_key"`
}
