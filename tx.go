// Copyright (c) 2013 Conformal Systems LLC.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package btcutil

import (
	"bytes"
	"github.com/conformal/btcwire"
)

// TxIndexUnknown is the value returned for a transaction index that is unknown.
// This is typically because the transaction has not been inserted into a block
// yet.
const TxIndexUnknown = -1

// Tx defines a bitcoin transaction that provides easier and more efficient
// manipulation of raw transactions.  It also memoizes the hash for the
// transaction on its first access so subsequent accesses don't have to repeat
// the relatively expensive hashing operations.
type Tx struct {
	msgTx        *btcwire.MsgTx   // Underlying MsgTx
	serializedTx []byte           // Serialized bytes for the transaction
	txSha        *btcwire.ShaHash // Cached transaction hash
	txIndex      int              // Position within a block or TxIndexUnknown
}

// MsgTx returns the underlying btcwire.MsgTx for the transaction.
func (t *Tx) MsgTx() *btcwire.MsgTx {
	// Return the cached transaction.
	return t.msgTx
}

// Sha returns the hash of the transaction.  This is equivalent to
// calling TxSha on the underlying btcwire.MsgTx, however it caches the
// result so subsequent calls are more efficient.
func (t *Tx) Sha() *btcwire.ShaHash {
	// Return the cached hash if it has already been generated.
	if t.txSha != nil {
		return t.txSha
	}

	// Generate the transaction hash.  Ignore the error since TxSha can't
	// currently fail.
	sha, _ := t.msgTx.TxSha()

	// Cache the hash and return it.
	t.txSha = &sha
	return &sha
}

// Index returns the saved index of the transaction within a block.  This value
// will be TxIndexUnknown if it hasn't already explicitly been set.
func (t *Tx) Index() int {
	return t.txIndex
}

// SetIndex sets the index of the transaction in within a block.
func (t *Tx) SetIndex(index int) {
	t.txIndex = index
}

// NewTx returns a new instance of a bitcoin transaction given an underlying
// btcwire.MsgTx.  See Tx.
func NewTx(msgTx *btcwire.MsgTx) *Tx {
	return &Tx{
		msgTx:   msgTx,
		txIndex: TxIndexUnknown,
	}
}

// NewTxFromBytes returns a new instance of a bitcoin transaction given the
// serialized bytes.  See Tx.
func NewTxFromBytes(serializedTx []byte) (*Tx, error) {
	// Deserialize the bytes into a MsgTx.
	var msgTx btcwire.MsgTx
	br := bytes.NewBuffer(serializedTx)
	err := msgTx.Deserialize(br)
	if err != nil {
		return nil, err
	}

	t := Tx{
		msgTx:        &msgTx,
		serializedTx: serializedTx,
		txIndex:      TxIndexUnknown,
	}
	return &t, nil
}