package usdb

import (
	"github.com/tendermint/tendermint/libs/db"
	"github.com/pingcap/tidb/kv"
	"github.com/kooksee/usmint/cmn"
	"github.com/pingcap/tidb/store/tikv"
	"context"
	"bytes"
	"fmt"
)

type TikvStore struct {
	db.DB
	name []byte
	c    kv.Storage
}

func NewTikvStore(name, url string) *TikvStore {
	tikv.MaxConnectionCount = 256

	// tikv://etcd-node1:port,etcd-node2:port?cluster=1&disableGC=false
	store, err := tikv.Driver{}.Open(fmt.Sprintf("tikv://%s/pd", url))
	cmn.MustNotErr("NewTikvStore Error", err)
	return &TikvStore{
		name: []byte(name),
		c:    store,
	}
}

func (db *TikvStore) withPrefix(key []byte) []byte {
	return append(db.name, key...)
}

func (db *TikvStore) withTxn(fn func(txn kv.Transaction) error) {
	txn, err := db.c.Begin()
	cmn.MustNotErr("tikv store open tx error", err)
	if err := fn(txn); err != nil && !kv.IsErrNotFound(fn(txn)) {
		cmn.MustNotErr("tikv store exec tx error", err)
	}

	cmn.MustNotErr("tikv store exec tx error", fn(txn))
	defer txn.Rollback()
	cmn.MustNotErr("tikv store commit tx error", txn.Commit(context.TODO()))
}

func (db *TikvStore) getSnapshot() kv.Snapshot {
	ss, err := db.c.GetSnapshot(kv.MaxVersion)
	cmn.MustNotErr("tikv store GetSnapshot error", err)
	return ss
}

// Implements DB.
func (db *TikvStore) Get(key []byte) []byte {
	ret, err := db.getSnapshot().Get(db.withPrefix(key))
	if !kv.IsErrNotFound(err) {
		cmn.MustNotErr("tikv store Get error", err)
	}
	return ret
}

// Implements DB.
func (db *TikvStore) Has(key []byte) bool {
	ret, err := db.getSnapshot().Get(db.withPrefix(key))
	return kv.IsErrNotFound(err) && len(ret) == 0
}

// Implements DB.
func (db *TikvStore) Set(key []byte, value []byte) {
	db.withTxn(func(txn kv.Transaction) (err error) {
		return txn.Set(db.withPrefix(key), value)
	})
}

// Implements DB.
func (db *TikvStore) SetSync(key []byte, value []byte) {
	db.Set(key, value)
}

// Implements DB.
func (db *TikvStore) Delete(key []byte) {
	db.withTxn(func(txn kv.Transaction) (err error) {
		return txn.Delete(db.withPrefix(key))
	})
}

// Implements DB.
func (db *TikvStore) DeleteSync(key []byte) {
	db.Delete(key)
}

// Implements DB.
func (db *TikvStore) Close() {
	cmn.MustNotErr("TikvStore Close Error", db.c.Close())
}

// Implements DB.
func (db *TikvStore) Print() {
}

// Implements DB.
func (db *TikvStore) Stats() map[string]string {
	//keys := []string{
	//	"leveldb.num-files-at-level{n}",
	//	"leveldb.stats",
	//	"leveldb.sstables",
	//	"leveldb.blockpool",
	//	"leveldb.cachedblock",
	//	"leveldb.openedtables",
	//	"leveldb.alivesnaps",
	//	"leveldb.aliveiters",
	//}

	return make(map[string]string)
}

//----------------------------------------
// Batch

// Implements DB.
func (db *TikvStore) NewBatch() db.Batch {
	return &tikvStoreBatch{data: make(map[string][]byte), db: db}
}

type tikvStoreBatch struct {
	db   *TikvStore
	data map[string][]byte
}

// Implements Batch.
func (m *tikvStoreBatch) Set(key, value []byte) {
	m.data[string(key)] = value
}

// Implements Batch.
func (m *tikvStoreBatch) Delete(key []byte) {
	delete(m.data, string(key))
}

// Implements Batch.
func (m *tikvStoreBatch) Write() {
	m.db.withTxn(func(txn kv.Transaction) error {
		for k, v := range m.data {
			if err := txn.Set([]byte(k), v); err != nil {
				return err
			}
		}
		return nil
	})
}

// Implements Batch.
func (m *tikvStoreBatch) WriteSync() {
	m.Write()
}

//----------------------------------------
// Iterator
// NOTE This is almost identical to db/c_level_db.Iterator
// Before creating a third version, refactor.

// Implements DB.
func (db *TikvStore) Iterator(start, end []byte) db.Iterator {
	it, err := db.getSnapshot().Seek(db.withPrefix(start))
	cmn.MustNotErr("TikvStore Iterator Error", err)
	return newTikvStoreIterator(db.name, false, it, db.withPrefix(start), db.withPrefix(end))
}

// Implements DB.
func (db *TikvStore) ReverseIterator(start, end []byte) db.Iterator {
	it, err := db.getSnapshot().SeekReverse(db.withPrefix(start))
	cmn.MustNotErr("TikvStore ReverseIterator Error", err)
	return newTikvStoreIterator(db.name, true, it, db.withPrefix(start), db.withPrefix(end))
}

type tikvStoreIterator struct {
	db.Iterator

	name    []byte
	r       kv.Iterator
	reverse bool
	start   []byte
	end     []byte
}

func newTikvStoreIterator(name []byte, reverse bool, r kv.Iterator, start, end []byte) *tikvStoreIterator {
	return &tikvStoreIterator{
		name:    name,
		r:       r,
		reverse: reverse,
		start:   start,
		end:     end,
	}
}

// Implements Iterator.
func (itr *tikvStoreIterator) Domain() ([]byte, []byte) {
	return itr.start, itr.end
}

// Implements Iterator.
func (itr *tikvStoreIterator) Valid() bool {
	if !itr.r.Valid() {
		return false
	}

	if !itr.reverse {
		if bytes.Compare(itr.r.Key(), itr.end) > 0 {
			return false
		}
	} else {
		if bytes.Compare(itr.r.Key(), itr.start) < 0 {
			return false
		}
	}

	return true
}

// Implements Iterator.
func (itr *tikvStoreIterator) Key() []byte {
	return bytes.TrimPrefix(itr.r.Key(), itr.name)
}

// Implements Iterator.
func (itr *tikvStoreIterator) Value() []byte {
	return itr.r.Value()
}

// Implements Iterator.
func (itr *tikvStoreIterator) Next() {
	cmn.MustNotErr("tikvStoreIterator next error", itr.r.Next())
}

// Implements Iterator.
func (itr *tikvStoreIterator) Close() {
	itr.r.Close()
}
