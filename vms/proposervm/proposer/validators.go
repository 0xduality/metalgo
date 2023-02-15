// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package proposer

import (
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/utils"
)

var _ utils.Sortable[validatorData] = validatorData{}

type validatorData struct {
	id     ids.NodeID
	weight uint64
}

func (d validatorData) Less(other validatorData) bool {
	return d.id.Less(other.id)
}
