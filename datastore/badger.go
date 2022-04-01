package datastore

import (
	"context"
	"errors"
	"github.com/dgraph-io/badger/v3"
	"github.com/quantosnetwork/dev-0.1.0/logger"
	"go.uber.org/zap"
	"log"
	"sync"
	"time"
)

type Repo struct {
	closeLk        sync.RWMutex
	closed         bool
	closeOnce      sync.Once
	closing        chan struct{}
	gcDiscardRatio float64
	gcSleep        time.Duration
	gcInterval     time.Duration

	syncWrites bool

	instance *badger.DB
	Logger   *zap.Logger
}

var ErrClosed = errors.New("closed")
var ErrNotFound = errors.New("not found")

type txn struct {
	ds       *Repo
	txn      *badger.Txn
	implicit bool
}

// Options are the badger datastore options, reexported here for convenience.
type Options struct {
	// Please refer to the Badger docs to see what this is for
	GcDiscardRatio float64

	// Interval between GC cycles
	//
	// If zero, the datastore will perform no automatic garbage collection.
	GcInterval time.Duration

	// Sleep time between rounds of a single GC cycle.
	//
	// If zero, the datastore will only perform one round of GC per
	// GcInterval.
	GcSleep time.Duration

	badger.Options
}

// DefaultOptions are the default options for the badger datastore.
var DefaultOptions Options

func init() {
	DefaultOptions = Options{
		GcDiscardRatio: 0.2,
		GcInterval:     15 * time.Minute,
		GcSleep:        10 * time.Second,
		Options:        badger.LSMOnlyOptions(""),
	}
	// This is to optimize the database on close so it can be opened
	// read-only and efficiently queried. We don't do that and hanging on
	// stop isn't nice.
	DefaultOptions.Options.CompactL0OnClose = false

	DefaultOptions.Options.BypassLockGuard = true

	// Explicitly set this to mmap. This doesn't use much memory anyways.
	DefaultOptions.Options.InMemory = true

	// Reduce this from 64MiB to 16MiB. That means badger will hold on to
	// 20MiB by default instead of 80MiB.
	//
	// This does not appear to have a significant performance hit.
	DefaultOptions.Options.BaseTableSize = 16 << 20
}

type badgerLog struct {
	*zap.Logger
}

func (b badgerLog) Errorf(s string, i ...interface{}) {
	log.Fatalf(s, i...)
}

func (b badgerLog) Warningf(s string, args ...interface{}) {
	log.Printf(s, args...)
}

func (b badgerLog) Infof(s string, i ...interface{}) {
	log.Printf(s, i...)

}

func (b badgerLog) Debugf(s string, i ...interface{}) {
	log.Printf(s, i...)
}

func (r *Repo) Create() (err error) {
	opt := badger.DefaultOptions("")

	opt.Dir = string("./data/db")
	opt.ValueDir = string("./data/db")
	opt.Logger = &badgerLog{logger.Logger}

	kv, err := badger.Open(opt)
	if err != nil {
		panic(err)
	}
	r.instance = kv
	r.closing = make(chan struct{})

	return

}

func (r *Repo) NewTransaction(ctx context.Context, readOnly bool) (*txn, error) {
	r.closeLk.RLock()
	defer r.closeLk.RUnlock()
	if r.closed {
		return &txn{}, errors.New("closed")
	}

	return &txn{r, r.instance.NewTransaction(!readOnly), false}, nil
}

func (r *Repo) NewImplicitTransaction(readOnly bool) *txn {
	return &txn{r, r.instance.NewTransaction(!readOnly), true}
}

type RepoKey []byte

func (r *Repo) Put(ctx context.Context, key RepoKey, value []byte) error {
	r.closeLk.RLock()
	defer r.closeLk.RUnlock()
	if r.closed {
		return errors.New("closed")
	}
	txn := r.NewImplicitTransaction(false)
	defer txn.discard()
	if err := txn.put(key, value); err != nil {
		return err
	}
	return txn.commit()

}

func (t *txn) Put(ctx context.Context, key RepoKey, value []byte) error {
	t.ds.closeLk.RLock()
	defer t.ds.closeLk.RUnlock()
	if t.ds.closed {
		return errors.New("closed")
	}
	return t.put(key, value)
}

func (t *txn) put(key RepoKey, value []byte) error {
	return t.txn.Set(key, value)
}

func (t *txn) Sync(ctx context.Context, prefix RepoKey) error {
	t.ds.closeLk.RLock()
	defer t.ds.closeLk.RUnlock()
	if t.ds.closed {
		return ErrClosed
	}

	return nil
}

func (t *txn) PutWithTTL(ctx context.Context, key RepoKey, value []byte, ttl time.Duration) error {
	t.ds.closeLk.RLock()
	defer t.ds.closeLk.RUnlock()
	if t.ds.closed {
		return ErrClosed
	}
	return t.putWithTTL(key, value, ttl)
}

func (t *txn) putWithTTL(key RepoKey, value []byte, ttl time.Duration) error {
	return t.txn.SetEntry(badger.NewEntry(key, value).WithTTL(ttl))
}

func (t *txn) GetExpiration(ctx context.Context, key RepoKey) (time.Time, error) {
	t.ds.closeLk.RLock()
	defer t.ds.closeLk.RUnlock()
	if t.ds.closed {
		return time.Time{}, ErrClosed
	}

	return t.getExpiration(key)
}

func (t *txn) getExpiration(key RepoKey) (time.Time, error) {
	item, err := t.txn.Get(key)
	if err == badger.ErrKeyNotFound {
		return time.Time{}, ErrNotFound
	} else if err != nil {
		return time.Time{}, err
	}
	return time.Unix(int64(item.ExpiresAt()), 0), nil
}

func (t *txn) SetTTL(ctx context.Context, key RepoKey, ttl time.Duration) error {
	t.ds.closeLk.RLock()
	defer t.ds.closeLk.RUnlock()
	if t.ds.closed {
		return ErrClosed
	}

	return t.setTTL(key, ttl)
}

func (t *txn) setTTL(key RepoKey, ttl time.Duration) error {
	item, err := t.txn.Get(key)
	if err != nil {
		return err
	}
	return item.Value(func(data []byte) error {
		return t.putWithTTL(key, data, ttl)
	})

}

func (t *txn) Get(ctx context.Context, key RepoKey) ([]byte, error) {
	t.ds.closeLk.RLock()
	defer t.ds.closeLk.RUnlock()
	if t.ds.closed {
		return nil, ErrClosed
	}

	return t.get(key)
}

func (t *txn) get(key RepoKey) ([]byte, error) {
	item, err := t.txn.Get(key)
	if err == badger.ErrKeyNotFound {
		err = ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return item.ValueCopy(nil)
}

func (t *txn) Has(ctx context.Context, key RepoKey) (bool, error) {
	t.ds.closeLk.RLock()
	defer t.ds.closeLk.RUnlock()
	if t.ds.closed {
		return false, ErrClosed
	}

	return t.has(key)
}

func (t *txn) has(key RepoKey) (bool, error) {
	_, err := t.txn.Get(key)
	switch err {
	case badger.ErrKeyNotFound:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}
}

// Alias to commit
func (t *txn) Close() error {
	t.ds.closeLk.RLock()
	defer t.ds.closeLk.RUnlock()
	if t.ds.closed {
		return ErrClosed
	}
	return t.close()
}

func (t *txn) close() error {
	return t.txn.Commit()
}

func (t *txn) Discard(ctx context.Context) {
	t.ds.closeLk.RLock()
	defer t.ds.closeLk.RUnlock()
	if t.ds.closed {
		return
	}

	t.discard()
}

func (t *txn) discard() {
	t.txn.Discard()
}

func (t *txn) Commit(ctx context.Context) error {
	t.ds.closeLk.RLock()
	defer t.ds.closeLk.RUnlock()
	if t.ds.closed {
		return ErrClosed
	}

	return t.commit()
}

func (t *txn) commit() error {
	return t.txn.Commit()
}

func expires(item *badger.Item) time.Time {
	return time.Unix(int64(item.ExpiresAt()), 0)
}

func NewDatastore() *Repo {
	r := &Repo{}
	err := r.Create()
	if err != nil {
		r.Logger.Fatal(err.Error())
	}
	defer r.instance.Close()
	return r
}
