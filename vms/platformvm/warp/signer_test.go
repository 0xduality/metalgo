// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package warp

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/utils/crypto/bls"
)

func TestSigner(t *testing.T) {
	for _, test := range SignerTests {
		sk, err := bls.NewSecretKey()
		require.NoError(t, err)

		chainID := ids.GenerateTestID()
		s := NewSigner(sk, chainID)

		test(t, s, sk, chainID)
	}
}
