/*
Package storage provides an in-memory key-value database engine.

Storage supports transactions and transaction trees.

Use New() to instantiate the storage.

Basic usage

Storage supports two basic operations - Set and Unset.

Set adds or modifies the variable by its key. Existing variable is overwritten.

Unset removes the variable.

Storage supports count-index by variables' values - use NumEqualTo to count variables
with given values.

Transactions

Storage supports transactions with unlimited nesting.

When the transaction is committed it traverses its parents all the way down to the root storage,
effectively committing the whole Tx tree.

Rollback rolls back only one transaction, returning its parent.

See DB interface for the API and usage examples.
*/
package storage
