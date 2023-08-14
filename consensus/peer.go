package consensus

import (
	"CoralDB/Cluster"
	BcGrpc "CoralDB/Proto/blockchain"
	"CoralDB/blockchain/blockchain_data"
	"CoralDB/util"
	"bytes"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strconv"
	"sync"
)

var (
	Pros      map[uint64][]*Proposal
	Blcoks    map[uint64][]*blockchain_data.Block
	InVotes   map[string][]*VoteMessage
	LocalPeer *Peer
	vmu       sync.RWMutex // 投票锁
	bmu       sync.RWMutex // 区块列表锁
	pmu       sync.RWMutex // 提议列表锁
)

type Peer struct {
	incomingVotes map[string]*List                    // 投票的map
	blocks        map[uint64][]*blockchain_data.Block // 提议区块map
	maxProposals  map[uint64][]*Proposal              //最大提议区块
}

func NewPeer() *Peer {
	Pros = make(map[uint64][]*Proposal)
	Blcoks = make(map[uint64][]*blockchain_data.Block)
	InVotes = make(map[string][]*VoteMessage)
	return &Peer{
		incomingVotes: make(map[string]*List),
		blocks:        make(map[uint64][]*blockchain_data.Block),
		maxProposals:  make(map[uint64][]*Proposal),
	}
}

func (p *Peer) Start() {
	//GetPeerPool().add(p)
	LocalPeer = p
}

//func (p *Peer) stop() {
//	GetPeerPool().remove(p)
//}

func ID() string {
	return LocalAlg.id
}

func (p *Peer) Gossip(typ int, data []byte) {
	//log.Println("gossip is Run")
	// simulate gossiping
	var nodeIp []string
	for _, peer := range Cluster.LocalNode.Node {
		if peer.IP+port == ID() {
			continue
		}
		nodeIp = append(nodeIp, fmt.Sprintf("%s:%d", peer.IP, peer.Port))
		//log.Println("ip,port:", fmt.Sprintf("%s:%d", peer.IP, peer.Port))
		//conn, err := grpc.Dial(fmt.Sprintf("%s:%d", peer.IP, peer.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		//if err != nil {
		//	fmt.Printf("%s:%d 网络异常", peer.IP, peer.Port)
		//}
		//// 获得grpc句柄
		//client := BcGrpc.NewBlockChainServiceClient(conn)
		//// 通过句柄调用函数
		//_, err = client.Handle(context.Background(), &BcGrpc.TypAndData{
		//	Typ:  int32(typ),
		//	Data: data,
		//})
		//if err != nil {
		//	log.Println(err)
		//}
		//err = conn.Close()
		//if err != nil {
		//	log.Panic(err)
		//}
	}
	num := len(nodeIp)
	switch num {
	case 3:
		go handle(nodeIp[0], typ, data)
		go handle(nodeIp[1], typ, data)
		go handle(nodeIp[2], typ, data)
	case 7:
		go handle(nodeIp[0], typ, data)
		go handle(nodeIp[1], typ, data)
		go handle(nodeIp[2], typ, data)
		go handle(nodeIp[3], typ, data)
		go handle(nodeIp[4], typ, data)
		go handle(nodeIp[5], typ, data)
		go handle(nodeIp[6], typ, data)
	case 11:
		go handle(nodeIp[0], typ, data)
		go handle(nodeIp[1], typ, data)
		go handle(nodeIp[2], typ, data)
		go handle(nodeIp[3], typ, data)
		go handle(nodeIp[4], typ, data)
		go handle(nodeIp[5], typ, data)
		go handle(nodeIp[6], typ, data)
		go handle(nodeIp[7], typ, data)
		go handle(nodeIp[8], typ, data)
		go handle(nodeIp[9], typ, data)
		go handle(nodeIp[10], typ, data)
	case 15:
		go handle(nodeIp[0], typ, data)
		go handle(nodeIp[1], typ, data)
		go handle(nodeIp[2], typ, data)
		go handle(nodeIp[3], typ, data)
		go handle(nodeIp[4], typ, data)
		go handle(nodeIp[5], typ, data)
		go handle(nodeIp[6], typ, data)
		go handle(nodeIp[7], typ, data)
		go handle(nodeIp[8], typ, data)
		go handle(nodeIp[9], typ, data)
		go handle(nodeIp[10], typ, data)
		go handle(nodeIp[11], typ, data)
		go handle(nodeIp[12], typ, data)
		go handle(nodeIp[13], typ, data)
		go handle(nodeIp[14], typ, data)
	case 19:
		go handle(nodeIp[0], typ, data)
		go handle(nodeIp[1], typ, data)
		go handle(nodeIp[2], typ, data)
		go handle(nodeIp[3], typ, data)
		go handle(nodeIp[4], typ, data)
		go handle(nodeIp[5], typ, data)
		go handle(nodeIp[6], typ, data)
		go handle(nodeIp[7], typ, data)
		go handle(nodeIp[8], typ, data)
		go handle(nodeIp[9], typ, data)
		go handle(nodeIp[10], typ, data)
		go handle(nodeIp[11], typ, data)
		go handle(nodeIp[12], typ, data)
		go handle(nodeIp[13], typ, data)
		go handle(nodeIp[14], typ, data)
		go handle(nodeIp[15], typ, data)
		go handle(nodeIp[16], typ, data)
		go handle(nodeIp[17], typ, data)
		go handle(nodeIp[18], typ, data)
	case 24:
		go handle(nodeIp[0], typ, data)
		go handle(nodeIp[1], typ, data)
		go handle(nodeIp[2], typ, data)
		go handle(nodeIp[3], typ, data)
		go handle(nodeIp[4], typ, data)
		go handle(nodeIp[5], typ, data)
		go handle(nodeIp[6], typ, data)
		go handle(nodeIp[7], typ, data)
		go handle(nodeIp[8], typ, data)
		go handle(nodeIp[9], typ, data)
		go handle(nodeIp[10], typ, data)
		go handle(nodeIp[11], typ, data)
		go handle(nodeIp[12], typ, data)
		go handle(nodeIp[13], typ, data)
		go handle(nodeIp[14], typ, data)
		go handle(nodeIp[15], typ, data)
		go handle(nodeIp[16], typ, data)
		go handle(nodeIp[17], typ, data)
		go handle(nodeIp[18], typ, data)
		go handle(nodeIp[19], typ, data)
		go handle(nodeIp[20], typ, data)
		go handle(nodeIp[21], typ, data)
		go handle(nodeIp[22], typ, data)
		go handle(nodeIp[23], typ, data)
	}
}

func handle(addr string, typ int, data []byte) {
	//log.Println("ip,port:", fmt.Sprintf("%s:%d", peer.IP, peer.Port))
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("%s 网络异常", addr)
	}
	// 获得grpc句柄
	client := BcGrpc.NewBlockChainServiceClient(conn)
	// 通过句柄调用函数
	_, err = client.Handle(context.Background(), &BcGrpc.TypAndData{
		Typ:  int32(typ),
		Data: data,
	})
	if err != nil {
		log.Println(err)
	}
	err = conn.Close()
	if err != nil {
		log.Panic(err)
	}
}

func (p *Peer) Handle(typ int, data []byte) error {
	if typ == BLOCK {
		blk := blockchain_data.Deserialize(data)
		//log.Printf("区块的广播: %x\n", blk.CurrentBlockHash)
		p.AddBlock(&blk)
		return nil
	} else if typ == BlockProposal {
		bp := &Proposal{}
		if err := bp.Deserialize(data); err != nil {
			return err
		}
		//pmu.RLock()
		//maxProposal := p.maxProposals[bp.Round]
		//pmu.RUnlock()
		//if maxProposal != nil {
		//	if typ == BlockProposal && bytes.Compare(bp.Prior, maxProposal.Prior) <= 0 {
		//		return nil
		//	}
		//}
		if err := bp.Verify(LocalAlg.weight(bp.Address()), constructSeed(LocalAlg.SortitionSeed(bp.Round), Role(util.Proposer, bp.Round, util.PROPOSE))); err != nil {
			fmt.Printf("block proposal verification failed, %s\n", err)
			return err
		}
		//log.Printf("区块的提议: %x\n", bp.Hash)
		p.AddProposal(bp.Round, bp)
		return nil
	} else if typ == VOTE {
		vote := &VoteMessage{}
		if err := vote.Deserialize(data); err != nil {
			return err
		}
		key := ConstructVoteKey(vote.Round, vote.Step)
		//vmu.RLock()
		//list, ok := p.incomingVotes[key]
		//log.Println("list key ", key)
		//vmu.RUnlock()
		//if !ok {
		//	list = newList()
		//}
		//list.add(vote)
		//log.Printf("投票: %x\n", vote.Hash)
		vmu.Lock()
		//p.incomingVotes[key] = list
		InVotes[key] = append(InVotes[key], vote)
		vmu.Unlock()
		//log.Println("len vote: ", len(InVotes[key]))
	}
	return nil
}

// iterator 返回传入消息队列的迭代器。
func (p *Peer) voteIterator(round uint64, step int) *Iterator {
	key := ConstructVoteKey(round, step)
	vmu.RLock()
	list, ok := p.incomingVotes[key]
	vmu.RUnlock()
	if !ok {
		list = newList()
		vmu.Lock()
		p.incomingVotes[key] = list
		vmu.Unlock()
	}
	return &Iterator{
		list: list,
	}
}

func (p *Peer) getIncomingMsgs(round uint64, step int) []interface{} {
	vmu.RLock()
	defer vmu.RUnlock()
	l := p.incomingVotes[ConstructVoteKey(round, step)]
	if l == nil {
		return nil
	}
	return l.list
}

func (p *Peer) getBlock(round uint64, hash []byte) *blockchain_data.Block {
	bmu.RLock()
	defer bmu.RUnlock()
	//index := util.BytesToInt32(hash)
	//log.Println(round, hash)
	var block *blockchain_data.Block
	for i := 0; i < len(Blcoks[round]); i++ {
		if bytes.Equal(Blcoks[round][i].CurrentBlockHash, hash) {
			block = Blcoks[round][i]
		}
	}
	return block
}

func (p *Peer) AddBlock(blk *blockchain_data.Block) {
	bmu.Lock()
	defer bmu.Unlock()
	//index := util.BytesToInt32(hash)
	//log.Println("add block")
	p.blocks[blk.Round] = append(p.blocks[blk.Round], blk)
	//p.blocks[blk.Round][index] = blk
	Blcoks[blk.Round] = append(Blcoks[blk.Round], blk)
	//log.Println("AddBlock:", blk.Round, len(Blcoks[blk.Round]), blk.CurrentBlockHash)
}

func (p *Peer) AddProposal(round uint64, proposal *Proposal) {
	pmu.Lock()
	defer pmu.Unlock()
	//log.Println("add proposal")
	p.maxProposals[round] = append(p.maxProposals[round], proposal)
	//log.Printf("node %s set max proposal #%d,%x", ID(), proposal.Round, p.maxProposals[round].Hash)
	//maxProposal = proposal
	Pros[round] = append(Pros[round], proposal)
	// 确保maxBlock
	//maxBlockIndex := util.BytesToInt32(maxProposal.Hash)
	//maxBlock = p.blocks[round][maxBlockIndex]

}

func (p *Peer) getMaxProposal(round uint64) *Proposal {
	pmu.RLock()
	defer pmu.RUnlock()
	//log.Println("len proposal :", len(p.maxProposals[round]))
	proposal := Pros[round][0]
	for i := 0; i < len(Pros[round]); i++ {
		if bytes.Compare(proposal.Prior, Pros[round][i].Prior) <= 0 {
			proposal = Pros[round][i]
		}
	}
	//log.Printf("node %s get max proposal #%d,%x", ID(), proposal.Round, proposal.Hash)
	return proposal
}

func (p *Peer) clearProposal(round uint64) {
	pmu.Lock()
	defer pmu.Unlock()
	delete(p.maxProposals, round)
}

func ConstructVoteKey(round uint64, step int) string {
	key := strconv.Itoa(int(round)) + strconv.Itoa(step)
	return key
}

type List struct {
	mu   sync.RWMutex
	list []interface{}
}

func newList() *List {
	return &List{}
}

func (l *List) add(el interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.list = append(l.list, el)
}

func (l *List) get(index int) interface{} {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if index >= len(l.list) {
		return nil
	}
	return l.list[index]
}

type Iterator struct {
	list  *List
	index int
}

func (it *Iterator) next() interface{} {
	el := it.list.get(it.index)
	if el == nil {
		return nil
	}
	it.index++
	return el
}

//type PeerPool struct {
//	mu    sync.Mutex
//	peers map[string]*Peer
//}

//func GetPeerPool() *PeerPool {
//	if peerPool == nil {
//		peerPool = &PeerPool{
//			peers: make(map[string]*Peer),
//		}
//	}
//	return peerPool
//}
//
//func (pool *PeerPool) add(peer *Peer) {
//	pool.mu.Lock()
//	defer pool.mu.Unlock()
//	pool.peers[ID()] = peer
//}
//
//func (pool *PeerPool) remove(peer *Peer) {
//	pool.mu.Lock()
//	defer pool.mu.Unlock()
//	delete(pool.peers, ID())
//}
//
//func (pool *PeerPool) getPeers() []*Peer {
//	pool.mu.Lock()
//	defer pool.mu.Unlock()
//	var peers []*Peer
//	for _, peer := range pool.peers {
//		peers = append(peers, peer)
//	}
//	return peers
//}
