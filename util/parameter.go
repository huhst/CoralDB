package util

import "time"

var UserAmount uint64 = 100
var TokenPerUser uint64 = 10000
var Malicious uint64 = 0
var NetworkLatency = 0

func TotalTokenAmount() uint64 {
	return UserAmount * TokenPerUser
}

const (
	// Algorand 系统参数
	ExpectedBlockProposers        = 26 // 期望区块提议者数量
	ExpectedCommitteeMembers      = 10
	ThresholdOfBAStep             = 0.2
	ExpectedFinalCommitteeMembers = 20
	FinalThreshold                = 0.1
	MAXSTEPS                      = 12

	// timeout param
	LamdaPriority = 300 * time.Millisecond // time to gossip sortition proofs.
	//LamdaBlock    = 1 * time.Minute // timeout for receiving a block.
	LamdaBlock   = 300 * time.Millisecond // timeout for receiving a block.
	LamdaStep    = 250 * time.Millisecond // timeout for BA* step.
	LamdaStepvar = 300 * time.Millisecond // estimate of BA* completion time variance.

	// interval
	R = 1000 // seed 刷新间隔 (# of rounds)

	// helper const var
	Committee = "committee"
	Proposer  = "proposer"

	// step
	PROPOSE      = 1000
	ReductionOne = 1001
	ReductionTwo = 1002
	FINAL        = 1003

	FinalConsensus     = 0
	TentativeConsensus = 1
)
