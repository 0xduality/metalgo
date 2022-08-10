// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/MetalBlockchain/avalanchego/ids"
	"github.com/MetalBlockchain/avalanchego/snow"
	"github.com/MetalBlockchain/avalanchego/utils/crypto"
	"github.com/MetalBlockchain/avalanchego/utils/timer/mockable"
	"github.com/MetalBlockchain/avalanchego/vms/components/avax"
	"github.com/MetalBlockchain/avalanchego/vms/platformvm/reward"
	"github.com/MetalBlockchain/avalanchego/vms/platformvm/stakeable"
	"github.com/MetalBlockchain/avalanchego/vms/platformvm/validator"
	"github.com/MetalBlockchain/avalanchego/vms/secp256k1fx"
)

func TestAddValidatorTxSyntacticVerify(t *testing.T) {
	assert := assert.New(t)
	clk := mockable.Clock{}
	ctx := snow.DefaultContextTest()
	ctx.AVAXAssetID = ids.GenerateTestID()
	signers := [][]*crypto.PrivateKeySECP256K1R{preFundedKeys}

	var (
		stx            *Tx
		addValidatorTx *AddValidatorTx
		err            error
	)

	// Case : signed tx is nil
	assert.ErrorIs(stx.SyntacticVerify(ctx), errNilSignedTx)

	// Case : unsigned tx is nil
	assert.ErrorIs(addValidatorTx.SyntacticVerify(ctx), ErrNilTx)

	validatorWeight := uint64(2022)
	rewardAddress := preFundedKeys[0].PublicKey().Address()
	inputs := []*avax.TransferableInput{{
		UTXOID: avax.UTXOID{
			TxID:        ids.ID{'t', 'x', 'I', 'D'},
			OutputIndex: 2,
		},
		Asset: avax.Asset{ID: ctx.AVAXAssetID},
		In: &secp256k1fx.TransferInput{
			Amt:   uint64(5678),
			Input: secp256k1fx.Input{SigIndices: []uint32{0}},
		},
	}}
	outputs := []*avax.TransferableOutput{{
		Asset: avax.Asset{ID: ctx.AVAXAssetID},
		Out: &secp256k1fx.TransferOutput{
			Amt: uint64(1234),
			OutputOwners: secp256k1fx.OutputOwners{
				Threshold: 1,
				Addrs:     []ids.ShortID{preFundedKeys[0].PublicKey().Address()},
			},
		},
	}}
	stakes := []*avax.TransferableOutput{{
		Asset: avax.Asset{ID: ctx.AVAXAssetID},
		Out: &stakeable.LockOut{
			Locktime: uint64(clk.Time().Add(time.Second).Unix()),
			TransferableOut: &secp256k1fx.TransferOutput{
				Amt: validatorWeight,
				OutputOwners: secp256k1fx.OutputOwners{
					Threshold: 1,
					Addrs:     []ids.ShortID{preFundedKeys[0].PublicKey().Address()},
				},
			},
		},
	}}
	addValidatorTx = &AddValidatorTx{
		BaseTx: BaseTx{BaseTx: avax.BaseTx{
			NetworkID:    ctx.NetworkID,
			BlockchainID: ctx.ChainID,
			Ins:          inputs,
			Outs:         outputs,
		}},
		Validator: validator.Validator{
			NodeID: ctx.NodeID,
			Start:  uint64(clk.Time().Unix()),
			End:    uint64(clk.Time().Add(time.Hour).Unix()),
			Wght:   validatorWeight,
		},
		Stake: stakes,
		RewardsOwner: &secp256k1fx.OutputOwners{
			Locktime:  0,
			Threshold: 1,
			Addrs:     []ids.ShortID{rewardAddress},
		},
		Shares: reward.PercentDenominator,
	}

	// Case: valid tx
	stx, err = NewSigned(addValidatorTx, Codec, signers)
	assert.NoError(err)
	assert.NoError(stx.SyntacticVerify(ctx))

	// Case: Wrong network ID
	addValidatorTx.SyntacticallyVerified = false
	addValidatorTx.NetworkID++
	stx, err = NewSigned(addValidatorTx, Codec, signers)
	assert.NoError(err)
	err = stx.SyntacticVerify(ctx)
	assert.Error(err)
	addValidatorTx.NetworkID--

	// Case: Stake owner has no addresses
	addValidatorTx.SyntacticallyVerified = false
	addValidatorTx.Stake[0].
		Out.(*stakeable.LockOut).
		TransferableOut.(*secp256k1fx.TransferOutput).
		Addrs = nil
	stx, err = NewSigned(addValidatorTx, Codec, signers)
	assert.NoError(err)
	err = stx.SyntacticVerify(ctx)
	assert.Error(err)
	addValidatorTx.Stake = stakes

	// Case: Rewards owner has no addresses
	addValidatorTx.SyntacticallyVerified = false
	addValidatorTx.RewardsOwner.(*secp256k1fx.OutputOwners).Addrs = nil
	stx, err = NewSigned(addValidatorTx, Codec, signers)
	assert.NoError(err)
	err = stx.SyntacticVerify(ctx)
	assert.Error(err)
	addValidatorTx.RewardsOwner.(*secp256k1fx.OutputOwners).Addrs = []ids.ShortID{rewardAddress}

	// Case: Too many shares
	addValidatorTx.SyntacticallyVerified = false
	addValidatorTx.Shares++ // 1 more than max amount
	stx, err = NewSigned(addValidatorTx, Codec, signers)
	assert.NoError(err)
	err = stx.SyntacticVerify(ctx)
	assert.Error(err)
	addValidatorTx.Shares--
}

func TestAddValidatorTxSyntacticVerifyNotAVAX(t *testing.T) {
	assert := assert.New(t)
	clk := mockable.Clock{}
	ctx := snow.DefaultContextTest()
	ctx.AVAXAssetID = ids.GenerateTestID()
	signers := [][]*crypto.PrivateKeySECP256K1R{preFundedKeys}

	var (
		stx            *Tx
		addValidatorTx *AddValidatorTx
		err            error
	)

	assetID := ids.GenerateTestID()
	validatorWeight := uint64(2022)
	rewardAddress := preFundedKeys[0].PublicKey().Address()
	inputs := []*avax.TransferableInput{{
		UTXOID: avax.UTXOID{
			TxID:        ids.ID{'t', 'x', 'I', 'D'},
			OutputIndex: 2,
		},
		Asset: avax.Asset{ID: assetID},
		In: &secp256k1fx.TransferInput{
			Amt:   uint64(5678),
			Input: secp256k1fx.Input{SigIndices: []uint32{0}},
		},
	}}
	outputs := []*avax.TransferableOutput{{
		Asset: avax.Asset{ID: assetID},
		Out: &secp256k1fx.TransferOutput{
			Amt: uint64(1234),
			OutputOwners: secp256k1fx.OutputOwners{
				Threshold: 1,
				Addrs:     []ids.ShortID{preFundedKeys[0].PublicKey().Address()},
			},
		},
	}}
	stakes := []*avax.TransferableOutput{{
		Asset: avax.Asset{ID: assetID},
		Out: &stakeable.LockOut{
			Locktime: uint64(clk.Time().Add(time.Second).Unix()),
			TransferableOut: &secp256k1fx.TransferOutput{
				Amt: validatorWeight,
				OutputOwners: secp256k1fx.OutputOwners{
					Threshold: 1,
					Addrs:     []ids.ShortID{preFundedKeys[0].PublicKey().Address()},
				},
			},
		},
	}}
	addValidatorTx = &AddValidatorTx{
		BaseTx: BaseTx{BaseTx: avax.BaseTx{
			NetworkID:    ctx.NetworkID,
			BlockchainID: ctx.ChainID,
			Ins:          inputs,
			Outs:         outputs,
		}},
		Validator: validator.Validator{
			NodeID: ctx.NodeID,
			Start:  uint64(clk.Time().Unix()),
			End:    uint64(clk.Time().Add(time.Hour).Unix()),
			Wght:   validatorWeight,
		},
		Stake: stakes,
		RewardsOwner: &secp256k1fx.OutputOwners{
			Locktime:  0,
			Threshold: 1,
			Addrs:     []ids.ShortID{rewardAddress},
		},
		Shares: reward.PercentDenominator,
	}

	stx, err = NewSigned(addValidatorTx, Codec, signers)
	assert.NoError(err)
	assert.Error(stx.SyntacticVerify(ctx))
}
