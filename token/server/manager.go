/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package server

import (
	"github.com/hyperledger/fabric/token/identity"
	"github.com/hyperledger/fabric/token/ledger"
	"github.com/hyperledger/fabric/token/tms/plain"
	"github.com/pkg/errors"
)

// Manager implements  token/server/TMSManager interface
// TODO: it will be updated after lscc-baased tms configuration is available
type Manager struct {
	LedgerManager              ledger.LedgerManager
	TokenOwnerValidatorManager identity.TokenOwnerValidatorManager
}

// For now it returns a plain issuer.
// After lscc-based tms configuration is available, it will be updated
// to return an issuer configured for the specific channel
func (m *Manager) GetIssuer(channel string, privateCredential, publicCredential []byte) (Issuer, error) {
	tokenOwnerValidator, err := m.TokenOwnerValidatorManager.Get(channel)
	if err != nil {
		return nil, errors.Wrapf(err, "failed getting token owner validator for channel: %s", channel)
	}

	return &plain.Issuer{TokenOwnerValidator: tokenOwnerValidator}, nil
}

// GetTransactor returns a Transactor bound to the passed channel and whose credential
// is the tuple (privateCredential, publicCredential).
func (m *Manager) GetTransactor(channel string, privateCredential, publicCredential []byte) (Transactor, error) {
	ledger, err := m.LedgerManager.GetLedgerReader(channel)
	if err != nil {
		return nil, errors.Wrapf(err, "failed getting ledger for channel: %s", channel)
	}

	tokenOwnerValidator, err := m.TokenOwnerValidatorManager.Get(channel)
	if err != nil {
		return nil, errors.Wrapf(err, "failed getting token owner validator for channel: %s", channel)
	}

	return &plain.Transactor{
		Ledger:              ledger,
		PublicCredential:    publicCredential,
		TokenOwnerValidator: tokenOwnerValidator}, nil
}
