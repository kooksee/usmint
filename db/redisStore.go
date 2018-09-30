package db

import (
	"bytes"
	"github.com/tendermint/tendermint/libs/db"
	"github.com/go-redis/redis"
)

type RedisDb struct {
	db.DB
	name string

	data map[string][]byte
	c    *redis.Client
}

func NewRedisDB(name, url string) (*RedisDb, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)
	if err := client.Ping().Err(); err != nil {
		return nil, err
	}

	return &RedisDb{c: client, name: name, data: make(map[string][]byte)}, nil
}

// Implements DB.
func (db *RedisDb) Get(key []byte) []byte {
	if dt, err := db.c.Get(string(key)).Bytes(); err != nil {
		panic(err.Error())
	} else {
		db.data[string(key)] = dt
	}
	return db.data[string(key)]
}

// Implements DB.
func (db *RedisDb) Has(key []byte) bool {
	return db.Get(key) != nil
}

// Implements DB.
func (db *RedisDb) Set(key []byte, value []byte) {
	db.data[string(key)] = value
}

// Implements DB.
func (db *RedisDb) SetSync(key []byte, value []byte) {
	db.data[string(key)] = value
}

// Implements DB.
func (db *RedisDb) Delete(key []byte) {
	db.data[string(key)] = nil
}

// Implements DB.
func (db *RedisDb) DeleteSync(key []byte) {
	db.data[string(key)] = nil
}

// Implements DB.
func (db *RedisDb) Close() {
	db.c.Close()
}

// Implements DB.
func (db *RedisDb) Print() {
}

// Implements DB.
func (db *RedisDb) Stats() map[string]string {
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
func (db *RedisDb) NewBatch() db.Batch {
	return &goLevelDBBatch{db: db, data: make(map[string][]byte)}
}

type goLevelDBBatch struct {
	db   *RedisDb
	data map[string][]byte
}

// Implements Batch.
func (mBatch *goLevelDBBatch) Set(key, value []byte) {
	mBatch.data[string(key)] = value
}

// Implements Batch.
func (mBatch *goLevelDBBatch) Delete(key []byte) {
	delete(mBatch.data, string(key))
}

// Implements Batch.
func (mBatch *goLevelDBBatch) Write() {
	for k, v := range mBatch.data {
		mBatch.db.Set([]byte(k), v)
	}
}

// Implements Batch.
func (mBatch *goLevelDBBatch) WriteSync() {
	mBatch.Write()
}

//----------------------------------------
// Iterator
// NOTE This is almost identical to db/c_level_db.Iterator
// Before creating a third version, refactor.

// Implements DB.
func (db *RedisDb) Iterator(start, end []byte) db.Iterator {
	return newRedisDBIterator(db.c, start, end, false)
}

// Implements DB.
func (db *RedisDb) ReverseIterator(start, end []byte) db.Iterator {
	return newRedisDBIterator(db.c, start, end, true)
}

type redisDBIterator struct {
	db.Iterator

	r *redis.Client

	start     []byte
	end       []byte
	isReverse bool
	isInvalid bool
}

func newRedisDBIterator(r *redis.Client, start, end []byte, isReverse bool) *redisDBIterator {
	return &redisDBIterator{
		r:         r,
		start:     start,
		end:       end,
		isReverse: isReverse,
		isInvalid: false,
	}
}

// Implements Iterator.
func (itr *redisDBIterator) Domain() ([]byte, []byte) {
	return itr.start, itr.end
}

// Implements Iterator.
func (itr *redisDBIterator) Valid() bool {

	// Once invalid, forever invalid.
	if itr.isInvalid {
		return false
	}

	// Panic on DB error.  No way to recover.
	itr.assertNoError()

	// If source is invalid, invalid.
	if !itr.source.Valid() {
		itr.isInvalid = true
		return false
	}

	// If key is end or past it, invalid.
	var end = itr.end
	var key = itr.source.Key()

	if itr.isReverse {
		if end != nil && bytes.Compare(key, end) <= 0 {
			itr.isInvalid = true
			return false
		}
	} else {
		if end != nil && bytes.Compare(end, key) <= 0 {
			itr.isInvalid = true
			return false
		}
	}

	// Valid
	return true
}

// Implements Iterator.
func (itr *redisDBIterator) Key() []byte {
	// Key returns a copy of the current key.
	// See https://github.com/syndtr/goleveldb/blob/52c212e6c196a1404ea59592d3f1c227c9f034b2/leveldb/iterator/iter.go#L88
	itr.assertNoError()
	itr.assertIsValid()
	return cp(itr.source.Key())
}

// Implements Iterator.
func (itr *redisDBIterator) Value() []byte {
	// Value returns a copy of the current value.
	// See https://github.com/syndtr/goleveldb/blob/52c212e6c196a1404ea59592d3f1c227c9f034b2/leveldb/iterator/iter.go#L88
	itr.assertNoError()
	itr.assertIsValid()
	return cp(itr.source.Value())
}

// Implements Iterator.
func (itr *redisDBIterator) Next() {
	itr.assertNoError()
	itr.assertIsValid()
	if itr.isReverse {
		itr.source.Prev()
	} else {
		itr.source.Next()
	}
}

// Implements Iterator.
func (itr *redisDBIterator) Close() {
	itr.source.Release()
}

func (itr *redisDBIterator) assertNoError() {
	if err := itr.source.Error(); err != nil {
		panic(err)
	}
}

func (itr redisDBIterator) assertIsValid() {
	if !itr.Valid() {
		panic("goLevelDBIterator is invalid")
	}
}
