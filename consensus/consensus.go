package consensus

import (
	"CoralDB/blockchain/blockchain_data"
	"CoralDB/common"
	"CoralDB/util"
	"bytes"
	"crypto/sha256"
	"errors"
	"log"
	"math/big"
	"math/rand"
	"time"
)

var LocalAlg *Algorand

var (
	ip                   = util.GetLocalIp()
	port                 = "3304"
	errCountVotesTimeout = errors.New("count votes timeout")

	// global metrics
	//metricsRound              uint64 = 1
	//proposerSelectedCounter          = metrics.NewRegisteredCounter("blockproposal/subusers/count", nil)
	//proposerSelectedHistogram        = metrics.NewRegisteredHistogram("blockproposal/subusers", nil, metrics.NewUniformSample(1028))
)

type Algorand struct {
	id      string           // 进程的id
	privkey *util.PrivateKey // 私钥
	Pubkey  *util.PublicKey  // 公钥

	chain *blockchain_data.BlockChain // 区块链
	peer  *Peer
	//quitCh      chan struct{}
	hangForever chan struct{}
}

func NewAlgorand() *Algorand {
	rand.Seed(time.Now().UnixNano())
	pub, priv, _ := util.NewKeyPair()
	alg := &Algorand{
		id:      ip + port,
		privkey: priv,
		Pubkey:  pub,
		chain:   blockchain_data.LocalDataBlockChain,
	}
	alg.peer = NewPeer()
	LocalAlg = alg
	return alg
}

func (alg *Algorand) Start() {
	//alg.quitCh = make(chan struct{})
	alg.hangForever = make(chan struct{})
	alg.peer.Start()
	//go alg.run()
}

func (alg *Algorand) Stop() {
	//close(alg.quitCh)
	close(alg.hangForever)
	//alg.peer.stop()
}

// Round 返回最新的整数。
func (alg *Algorand) Round() uint64 {
	return alg.lastBlock().Round
}

func (alg *Algorand) lastBlock() *blockchain_data.Block {
LOOP:
	block, err := alg.chain.GetBlockByHash(alg.chain.TailHash)
	if err != nil && err.Error() == "not the key" {
		time.Sleep(time.Millisecond)
		goto LOOP
	}
	return block
}

// weight 返回给定地址的权重。
func (alg *Algorand) weight(address common.Address) uint64 {
	//return TokenPerUser todo
	return uint64(rand.Intn(9950) + 50)
}

// TokenOwn 返回自身节点拥有的令牌数量（权重）。
func (alg *Algorand) TokenOwn() uint64 {
	return alg.weight(alg.Address())
}

// seed 返回块r的基于vrf的种子。
func (alg *Algorand) vrfSeed(round uint64) (seed, proof []byte, err error) {
	if round == 0 {
		emptySeed := sha256.Sum256([]byte{})
		return emptySeed[:], nil, nil
	}
	lastBlock := alg.chain.GetByRound(round - 1)
	// 最后一个区块不是genesis区块，验证种子r-1。
	if round != 1 {
		lastParentBlock, _ := alg.chain.GetBlockByHash(alg.chain.TailHash)
		if lastBlock.Proof != nil {
			// vrf-based seed
			pubkey := util.RecoverPubkey(lastBlock.Signature)
			m := bytes.Join([][]byte{lastParentBlock.Seed, common.Uint2Bytes(lastBlock.Round)}, nil)
			err = pubkey.VerifyVRF(lastBlock.Proof, m)
		} else if !bytes.Equal(lastBlock.Seed, common.Sha256(bytes.Join([][]byte{lastParentBlock.Seed, common.Uint2Bytes(lastBlock.Round)}, nil)).Bytes()) {
			// hash-based seed
			err = errors.New("hash seed invalid")
		}
		if err != nil {
			// seed r-1 invalid
			return common.Sha256(bytes.Join([][]byte{lastBlock.Seed, common.Uint2Bytes(lastBlock.Round + 1)}, nil)).Bytes(), nil, nil
		}
	}

	seed, proof, err = alg.privkey.Evaluate(bytes.Join([][]byte{lastBlock.Seed, common.Uint2Bytes(lastBlock.Round + 1)}, nil))
	return
}

func (alg *Algorand) emptyBlock(round uint64, prevHash []byte) *blockchain_data.Block {
	return &blockchain_data.Block{
		Round:             round,
		PreviousBlockHash: prevHash,
	}
}

// SortitionSeed 返回具有刷新间隔R的选择种子。
func (alg *Algorand) SortitionSeed(round uint64) []byte {
	realR := round - 1
	mod := round % util.R
	if realR < mod {
		seed := sha256.Sum256([]byte{})
		return seed[:]
	} else {
		realR -= mod
	}

	return alg.chain.GetByRound(realR).Seed
}

func (alg *Algorand) Address() common.Address {
	return common.BytesToAddress(alg.Pubkey.Bytes())
}

// ProcessMain 执行consensus算法的主要处理。
func (alg *Algorand) ProcessMain(transactions []*blockchain_data.Transaction, round uint64, vrf, proof []byte, subusers int) *blockchain_data.Block {
	currRound := alg.Round() + 1
	// 1. 区块提议
	//start := time.Now()
	block := alg.blockProposal(transactions, round, vrf, proof, subusers)
	//elapsed := time.Since(start)
	//fmt.Println("consensus  blockProposal have ", elapsed)
	//log.Printf("node %s init BA with block #%d %x, is empty? %v\n", alg.id, block.Round, block.CurrentBlockHash[:], block.Signature == nil)

	// 2. 用具有最高优先级的block初始化BA。
	//consensusType, block := alg.BA(currRound, block)
	//start1 := time.Now()
	_, block = alg.BA(currRound, block)
	//elapsed1 := time.Since(start1)
	//fmt.Println("consensus BA have ", elapsed1)
	//3. 就最终或暂定新区块达成共识。
	//log.Printf("node %s reach consensus %d at Round %d, block hash %x, is empty? %v\n", alg.id, consensusType, currRound, block.CurrentBlockHash, block.Signature == nil)
	return block
}

// proposeBlock 提出一个新的块
func (alg *Algorand) proposeBlock(transactions []*blockchain_data.Transaction) *blockchain_data.Block {
	currRound := alg.Round() + 1

	seed, proof, err := alg.vrfSeed(currRound)
	if err != nil {
		return alg.emptyBlock(currRound, alg.lastBlock().CurrentBlockHash)
	}
	block := blockchain_data.NewBlock()
	block.InitBlock(transactions, alg.chain.TailHash, alg.chain.LastID)
	block.Round = currRound
	block.Seed = seed
	block.Author = alg.Pubkey.Address()
	block.Proof = proof
	sign, _ := alg.privkey.Sign(block.CurrentBlockHash)
	block.Signature = sign
	//log.Printf("node %s propose a new block #%d %x\n", alg.id, block.Round, block.CurrentBlockHash)
	return &block
}

// blockProposal 执行块提议过程。
func (alg *Algorand) blockProposal(transactions []*blockchain_data.Transaction, round uint64, vrf, proof []byte, subusers int) *blockchain_data.Block {
	//round := alg.Round() + 1
	//vrf, proof, subusers := alg.Sortition(alg.SortitionSeed(round), Role(util.Proposer, round, util.PROPOSE), util.ExpectedBlockProposers, alg.TokenOwn())
	// have been selected.
	if subusers > 0 {
		var (
			newBlk       *blockchain_data.Block
			proposalType int
		)

		newBlk = alg.proposeBlock(transactions)
		proposalType = BlockProposal

		proposal := &Proposal{
			Round:  newBlk.Round,
			Hash:   newBlk.CurrentBlockHash,
			Prior:  maxPriority(vrf, subusers),
			VRF:    vrf,
			Proof:  proof,
			Pubkey: alg.Pubkey.Bytes(),
		}
		alg.peer.AddProposal(round, proposal)
		//log.Printf("自己提议区块的hash :%x\n", newBlk.CurrentBlockHash)
		alg.peer.AddBlock(newBlk)
		blkMsg := newBlk.Serialize()
		proposalMsg, _ := proposal.Serialize()

		//log.Println("gossip:", BLOCK, proposalType)

		go alg.peer.Gossip(BLOCK, blkMsg)
		go alg.peer.Gossip(proposalType, proposalMsg)

	}
	//log.Println("len proposal :", len(alg.peer.maxProposals[round]))

	//等待 λstepvar+λpriority 时间以确定最高优先级。
	timeoutForPriority := time.NewTimer(util.LamdaStepvar + util.LamdaPriority)
	<-timeoutForPriority.C

	// block提议超时
	timeoutForBlockFlying := time.NewTimer(util.LamdaBlock)
	ticker := time.NewTicker(200 * time.Millisecond)

	for {
		select {
		case <-timeoutForBlockFlying.C:
			// 空的block
			//log.Printf("maxBlock1:%x\n", alg.emptyBlock(round, alg.lastBlock().CurrentBlockHash).CurrentBlockHash)
			return alg.emptyBlock(round, alg.lastBlock().CurrentBlockHash)
		case <-ticker.C:
			// 获取具有最高优先级的块
			//pp := maxProposal
			pp := alg.peer.getMaxProposal(round)
			if pp == nil {
				continue
			}
			//log.Printf("pp.Hash: %x\n", pp.Prior)
			//blk := maxBlock
			blk := alg.peer.getBlock(round, pp.Hash)
			if blk != nil {
				//log.Printf("maxBlock2:%x\n", blk.CurrentBlockHash)
				return blk
			}
		}
	}
}

// Sortition 运行加密选择过程并返回vrf、证明和所选子用户的数量。
func (alg *Algorand) Sortition(seed, role []byte, expectedNum int, weight uint64) (vrf, proof []byte, selected int) {
	vrf, proof, _ = alg.privkey.Evaluate(constructSeed(seed, role))
	selected = subUsers(expectedNum, weight, vrf)
	return
}

// verifySort 验证vrf并返回所选子用户的数量。
func (alg *Algorand) verifySort(vrf, proof, seed, role []byte, expectedNum int) int {
	if err := alg.Pubkey.VerifyVRF(proof, constructSeed(seed, role)); err != nil {
		return 0
	}

	return subUsers(expectedNum, alg.TokenOwn(), vrf)
}

// committeeVote votes for `value`.
func (alg *Algorand) committeeVote(round uint64, step int, expectedNum int, hash []byte) error {

	vrf, proof, j := alg.Sortition(alg.SortitionSeed(round), Role(util.Committee, round, step), expectedNum, alg.TokenOwn())

	if j > 0 {
		// 广播投票信息
		voteMsg := &VoteMessage{
			Round:      round,
			Step:       step,
			VRF:        vrf,
			Proof:      proof,
			ParentHash: alg.chain.TailHash,
			Hash:       hash,
		}
		_, err := voteMsg.Sign(alg.privkey)
		if err != nil {
			return err
		}
		data, err := voteMsg.Serialize()
		if err != nil {
			return err
		}
		go alg.peer.Gossip(VOTE, data)
	}
	return nil
}

// BA 在下一轮次中，用一个提议的区块运行BA*共识
func (alg *Algorand) BA(round uint64, block *blockchain_data.Block) (int8, *blockchain_data.Block) {
	var (
		newBlk *blockchain_data.Block
		hash   []byte
	)
	hash = alg.reduction(round, block.CurrentBlockHash)
	//log.Println("reduction end!", hash)
	hash = alg.binaryBA(round, hash)
	//log.Println("binaryBA end!", hash)
	r, _ := alg.countVotes(round, util.FINAL, util.FinalThreshold, util.ExpectedFinalCommitteeMembers, util.LamdaStep)
	if prevHash := alg.lastBlock().CurrentBlockHash; bytes.Equal(hash, emptyHash(round, prevHash)) {
		// empty block
		newBlk = alg.emptyBlock(round, prevHash)
	} else {
		//LOOP:
		newBlk = alg.peer.getBlock(round, hash)
		//if newBlk == nil {
		//	time.Sleep(time.Microsecond)
		//	goto LOOP
		//}
		//newBlk = maxBlock
	}
	//log.Printf("%s newBlock is %d\n", ip, newBlk.Type)
	if bytes.Equal(r, hash) {
		newBlk.Type = util.FinalConsensus
		return util.FinalConsensus, newBlk
	} else {
		newBlk.Type = util.TentativeConsensus
		return util.TentativeConsensus, newBlk
	}
}

// 第二步 reduction.
func (alg *Algorand) reduction(round uint64, hash []byte) []byte {
	// step 1: 广播block的hash
	//log.Printf("reduction committeeVote begin: %x\n", hash)
	alg.committeeVote(round, util.ReductionOne, util.ExpectedCommitteeMembers, hash)
	//log.Printf("reduction committeeVote end: %x\n", hash)
	// 其他用户可能仍在等待块建议
	// 设置等待时间为 λblock + λstep
	empty := emptyHash(round, alg.chain.TailHash)
	// todo 投票为空
	//log.Printf("Round %d reduction countVotes begin: %x\n", round, hash)
	hash1, err := alg.countVotes(round, util.ReductionOne, util.ThresholdOfBAStep, util.ExpectedCommitteeMembers, util.LamdaBlock+util.LamdaStep)
	//log.Printf("Round %d reduction countVotes end: %x\n", round, hash1)
	// step 2: 重新广播block的hash
	//empty := emptyHash(round, alg.chain.TailHash)

	//if err == errCountVotesTimeout {
	//	//log.Printf("reduction 投票2: %x\n", empty)
	//	alg.committeeVote(round, util.ReductionTwo, util.ExpectedCommitteeMembers, empty)
	//} else {
	//	//log.Printf("reduction 投票3: %x\n", hash1)
	//	alg.committeeVote(round, util.ReductionTwo, util.ExpectedCommitteeMembers, hash1)
	//}

	//hash2, err := alg.countVotes(round, util.ReductionTwo, util.ThresholdOfBAStep, util.ExpectedCommitteeMembers, util.LamdaStep)
	if err == errCountVotesTimeout {
		return empty
	}
	//log.Printf("reduction 投票3: %x\n", hash1)
	return hash1
}

// binaryBA 执行，直到对给定的“hash”或“empty_hash”达成共识。
func (alg *Algorand) binaryBA(round uint64, hash []byte) []byte {
	var (
		step = 1
		r    = hash
		err  error
	)
	empty := emptyHash(round, alg.chain.TailHash)
	defer func() {
		//fmt.Printf("node %s complete binaryBA with %d steps\n", alg.id, step)
	}()
	for step < util.MAXSTEPS {
		// todo 检查投票的问题
		alg.committeeVote(round, step, util.ExpectedCommitteeMembers, r)
		r, err = alg.countVotes(round, step, util.ThresholdOfBAStep, util.ExpectedCommitteeMembers, util.LamdaStep)
		if err == errCountVotesTimeout {
			r = hash
		} else if !bytes.Equal(r, empty) {
			for s := step + 1; s <= step+3; s++ {
				alg.committeeVote(round, s, util.ExpectedCommitteeMembers, r)
			}
			if step == 1 {
				alg.committeeVote(round, util.FINAL, util.ExpectedFinalCommitteeMembers, r)
			}
			return r
		}
		step++

		alg.committeeVote(round, step, util.ExpectedCommitteeMembers, r)
		r, err = alg.countVotes(round, step, util.ThresholdOfBAStep, util.ExpectedCommitteeMembers, util.LamdaStep)
		if err == errCountVotesTimeout {
			r = empty
		} else if bytes.Equal(r, empty) {
			for s := step + 1; s <= step+3; s++ {
				alg.committeeVote(round, s, util.ExpectedCommitteeMembers, r)
			}
			return r
		}
		step++

		alg.committeeVote(round, step, util.ExpectedCommitteeMembers, r)
		r, err = alg.countVotes(round, step, util.ThresholdOfBAStep, util.ExpectedCommitteeMembers, util.LamdaStep)
		if err == errCountVotesTimeout {
			if alg.commonCoin(round, step, util.ExpectedCommitteeMembers) == 0 {
				r = hash
			} else {
				r = empty
			}
		}
	}

	log.Printf("reach the maxstep hang forever\n")
	// hang forever
	<-alg.hangForever
	return []byte{}
}

// countVotes 计算轮次和步数的选票
func (alg *Algorand) countVotes(round uint64, step int, threshold float64, expectedNum int, timeout time.Duration) (resHash []byte, err error) {
	expired := time.NewTimer(timeout)
	counts := make(map[int32]int)
	voters := make(map[string]struct{})
	//it := alg.peer.voteIterator(round, step)
	key := ConstructVoteKey(round, step)
	// 等待
	time.Sleep(timeout)
	// 获取投票的数量
	voteLens := len(InVotes[key])
	if voteLens == 0 {
		select {
		case <-expired.C:
			resHash = []byte{}
			err = errCountVotesTimeout
		default:
		}
	} else {
		//log.Println("统计投票")
		for i := 0; i < voteLens; i++ {
			voteMsg := InVotes[key][i]
			votes, hash, _ := alg.processMsg(voteMsg, expectedNum)
			//log.Printf("投票的hash %x\n", hash)
			pubkey := voteMsg.RecoverPubkey()
			if _, exist := voters[string(pubkey.Pk)]; exist || votes == 0 {
				continue
			}
			voters[string(pubkey.Pk)] = struct{}{}
			if bytes.Equal(hash, []byte{}) {
				continue
			}
			index := util.BytesToInt32(hash)
			counts[index] += votes
			// if we got enough votes, then output the target hash
			//log.Infof("node %d receive votes %v,threshold %v at step %d", alg.id, counts[hash], uint64(float64(expectedNum)*threshold), step)
			if uint64(counts[index]) >= uint64(float64(expectedNum)*threshold) {
				resHash = hash
				err = nil
				//return hash, nil
			}
		}
	}
	//log.Printf("投票的hash %x\n", resHash)
	return

	//for {
	//	msg := it.next()
	//	if msg == nil {
	//		select {
	//		case <-expired.C:
	//			// timeout
	//			return []byte{}, errCountVotesTimeout
	//		default:
	//		}
	//	} else {
	//		log.Println("统计投票")
	//		voteMsg := msg.(*VoteMessage)
	//		votes, hash, _ := alg.processMsg(msg.(*VoteMessage), expectedNum)
	//		pubkey := voteMsg.RecoverPubkey()
	//		if _, exist := voters[string(pubkey.Pk)]; exist || votes == 0 {
	//			continue
	//		}
	//		voters[string(pubkey.Pk)] = struct{}{}
	//		index := util.BytesToInt32(hash)
	//		counts[index] += votes
	//		// if we got enough votes, then output the target hash
	//		//log.Infof("node %d receive votes %v,threshold %v at step %d", alg.id, counts[hash], uint64(float64(expectedNum)*threshold), step)
	//		if uint64(counts[index]) >= uint64(float64(expectedNum)*threshold) {
	//			return hash, nil
	//		}
	//	}
	//}
}

// processMsg 验证传入的投票消息。
func (alg *Algorand) processMsg(message *VoteMessage, expectedNum int) (votes int, hash []byte, vrf []byte) {
	if err := message.VerifySign(); err != nil {
		return 0, []byte{}, nil
	}

	// discard messages that do not extend this chain
	prevHash := message.ParentHash
	if !bytes.Equal(prevHash, alg.chain.TailHash) {
		return 0, []byte{}, nil
	}

	votes = alg.verifySort(message.VRF, message.Proof, alg.SortitionSeed(message.Round), Role(util.Committee, message.Round, message.Step), expectedNum)
	hash = message.Hash
	vrf = message.VRF
	return
}

// commonCoin 计算所有用户共用的硬币。
// 如果对手向网络发送错误消息并阻止网络达成共识，这是一种帮助Algorand恢复的程序。
func (alg *Algorand) commonCoin(round uint64, step int, expectedNum int) int64 {
	minhash := new(big.Int).Exp(big.NewInt(2), big.NewInt(common.HashLength), big.NewInt(0))
	msgList := alg.peer.getIncomingMsgs(round, step)
	for _, m := range msgList {
		msg := m.(*VoteMessage)
		votes, _, vrf := alg.processMsg(msg, expectedNum)
		for j := 1; j < votes; j++ {
			h := new(big.Int).SetBytes(common.Sha256(bytes.Join([][]byte{vrf, common.Uint2Bytes(uint64(j))}, nil)).Bytes())
			if h.Cmp(minhash) < 0 {
				minhash = h
			}
		}
	}
	return minhash.Mod(minhash, big.NewInt(2)).Int64()
}

// Role 返回当前回合和步骤中的role的字节
func Role(iden string, round uint64, step int) []byte {
	return bytes.Join([][]byte{
		[]byte(iden),
		common.Uint2Bytes(round),
		common.Uint2Bytes(uint64(step)),
	}, nil)
}

// maxPriority 返回最高优先级的提议块
func maxPriority(vrf []byte, users int) []byte {
	var maxPrior []byte
	for i := 1; i <= users; i++ {
		prior := common.Sha256(bytes.Join([][]byte{vrf, common.Uint2Bytes(uint64(i))}, nil)).Bytes()
		if bytes.Compare(prior, maxPrior) > 0 {
			maxPrior = prior
		}
	}
	//log.Printf("maxPrior: %x\n", maxPrior)
	return maxPrior
}

// subUsers 返回根据数学协议确定的所选“子用户”数量
func subUsers(expectedNum int, weight uint64, vrf []byte) int {
	//binomial := NewBinomial(int64(weight), int64(expectedNum), int64(TotalTokenAmount()))
	binomial := util.NewApproxBinomial(int64(expectedNum), weight)
	// hash / 2^hashlen ∉ [ ∑0,j B(k;w,p), ∑0,j+1 B(k;w,p))
	hashBig := new(big.Int).SetBytes(vrf)
	maxHash := new(big.Int).Exp(big.NewInt(2), big.NewInt(common.HashLength*8), nil)
	hash := new(big.Rat).SetFrac(hashBig, maxHash)
	var lower, upper *big.Rat
	j := 0
	for uint64(j) <= weight {
		if upper != nil {
			lower = upper
		} else {
			lower = binomial.CDF(int64(j))
		}
		upper = binomial.CDF(int64(j + 1))
		//log.Infof("hash %v, lower %v , upper %v", hash.Sign(), lower.Sign(), upper.Sign())
		if hash.Cmp(lower) >= 0 && hash.Cmp(upper) < 0 {
			break
		}
		j++
	}
	//log.Infof("j %d", j)
	if uint64(j) > weight {
		j = 0
	}
	//j := parallelTrevels(runtime.NumCPU(), weight, hash, binomial)
	return j
}

//func parallelTrevels(core int, N uint64, hash *big.Rat, binomial util.Binomial) int {
//	var wg sync.WaitGroup
//	groups := N / uint64(core)
//	background, cancel := context.WithCancel(context.Background())
//	resChan := make(chan int)
//	notFound := make(chan struct{})
//	for i := 0; i < core; i++ {
//		go func(ctx context.Context, begin uint64) {
//			wg.Add(1)
//			defer wg.Done()
//			var (
//				end          uint64
//				upper, lower *big.Rat
//			)
//			if begin == uint64(core-2) {
//				end = N + 1
//			} else {
//				end = groups * (begin + 1)
//			}
//			for j := groups * begin; j < end; j++ {
//				select {
//				case <-ctx.Done():
//					return
//				default:
//				}
//				if upper != nil {
//					lower = upper
//				} else {
//					lower = binomial.CDF(int64(j))
//				}
//				upper = binomial.CDF(int64(j + 1))
//				//log.Infof("hash %v, lower %v , upper %v", hash.Sign(), lower.Sign(), upper.Sign())
//				if hash.Cmp(lower) >= 0 && hash.Cmp(upper) < 0 {
//					resChan <- int(j)
//					return
//				}
//				j++
//			}
//			return
//		}(background, uint64(i))
//	}
//
//	go func() {
//		wg.Wait()
//		close(notFound)
//	}()
//
//	select {
//	case j := <-resChan:
//		cancel()
//		return j
//	case <-notFound:
//		return 0
//	}
//}

// constructSeed 为vrf生成构造一个新的字节
func constructSeed(seed, role []byte) []byte {
	return bytes.Join([][]byte{seed, role}, nil)
}

func emptyHash(round uint64, prev []byte) []byte {
	hash := sha256.Sum256(bytes.Join([][]byte{
		common.Uint2Bytes(round),
		prev,
	}, nil))
	return hash[:]
}
