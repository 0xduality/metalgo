// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package atomic

import (
	"github.com/MetalBlockchain/avalanchego/codec"
	"github.com/MetalBlockchain/avalanchego/codec/linearcodec"
)

const codecVersion = 0

// codecManager is used to marshal and unmarshal dbElements and chain IDs.
var codecManager codec.Manager

func init() {
	linearCodec := linearcodec.NewDefault()
	codecManager = codec.NewDefaultManager()
	if err := codecManager.RegisterCodec(codecVersion, linearCodec); err != nil {
		panic(err)
	}
}
