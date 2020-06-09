package platonstats

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/syndtr/goleveldb/leveldb/errors"

	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"

	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"

	"github.com/PlatONnetwork/PlatON-Go/core/state"

	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"

	"github.com/PlatONnetwork/PlatON-Go/common"

	"github.com/PlatONnetwork/PlatON-Go/eth"

	"github.com/PlatONnetwork/PlatON-Go/trie"

	"github.com/PlatONnetwork/PlatON-Go/rlp"

	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/PlatONnetwork/PlatON-Go/rpc"

	"github.com/PlatONnetwork/PlatON-Go/core/types"

	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/Shopify/sarama"
)

const (
	// historyUpdateRange is the number of blocks a node should report upon login or
	// history request.
	sampleEventChanSize = 50

	kafkaBlockTopic = "platon-block"
)

var (
	statsLogPath = "./"
	statsLogFile = "platonstats.log"
)

func SetPlatonStatsLogPath(path string) {
	statsLogPath = path
	statsLogFile = filepath.Join(statsLogPath, statsLogFile)
	log.Info("PlatON stats service log file", "path", statsLogFile)
}

type platonStats interface {
	SubscribeSampleEvent(ch chan<- SampleEvent) event.Subscription
}

type StatsBlockExt struct {
	BlockType    uint8
	EpochNo      uint64
	NodeID       discover.NodeID      `rlp:"nil"`
	Block        *types.Block         `rlp:"nil"`
	Receipts     []*types.Receipt     `rlp:"nil"`
	ExeBlockData *common.ExeBlockData `rlp:"nil"`
	GenesisData  *common.GenesisData  `rlp:"nil"`
}

type PlatonStatsService struct {
	server *p2p.Server // Peer-to-peer server to retrieve networking infos

	kafkaUrl string
	eth      *eth.Ethereum // Full Ethereum service if monitoring a full node

	blockProducer sarama.SyncProducer
	msgProducer   sarama.AsyncProducer

	stopSampleMsg chan struct{}
	stopBlockMsg  chan struct{}
	stopOnce      sync.Once
}

var (
	//platonStatsServiceOnce sync.Once
	platonStatsService *PlatonStatsService
)

func New(kafkaUrl string, ethServ *eth.Ethereum) (*PlatonStatsService, error) {
	platonStatsService = &PlatonStatsService{
		kafkaUrl: kafkaUrl,
		eth:      ethServ,
	}
	return platonStatsService, nil
}

func GetPlatonStatsService() *PlatonStatsService {
	return platonStatsService
}

func (s *PlatonStatsService) BlockChain() *core.BlockChain {
	return s.eth.BlockChain()
}

func (s *PlatonStatsService) ChainDb() ethdb.Database {
	return s.eth.ChainDb()
}

// Protocols implements node.Service, returning the P2P network protocols used
// by the stats service (nil as it doesn't use the devp2p overlay network).
func (s *PlatonStatsService) Protocols() []p2p.Protocol { return nil }

// APIs implements node.Service, returning the RPC API endpoints provided by the
// stats service (nil as it doesn't provide any user callable APIs).
func (s *PlatonStatsService) APIs() []rpc.API { return nil }

// Start implements node.Service, starting up the monitoring and reporting daemon.
func (s *PlatonStatsService) Start(server *p2p.Server) error {
	s.server = server
	urls := []string{s.kafkaUrl}

	if msgProducer, err := sarama.NewAsyncProducer(urls, msgProducerConfig()); err != nil {
		return err
	} else {
		s.msgProducer = msgProducer
	}

	if blockProducer, err := sarama.NewSyncProducer(urls, blockProducerConfig()); err != nil {
		return err
	} else {
		s.blockProducer = blockProducer
	}

	defer func() {
		if s.msgProducer != nil {
			s.msgProducer.AsyncClose()
		}
		if s.blockProducer != nil {
			s.blockProducer.Close()
		}
	}()

	go s.blockMsgLoop()
	go s.sampleMsgLoop()

	log.Info("PlatON stats daemon started")
	return nil
}
func blockProducerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // 发送完数据需要leader和follow都确认
	config.Producer.Return.Successes = true
	config.Producer.Compression = sarama.CompressionGZIP
	return config
}

func msgProducerConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true
	config.Producer.Compression = sarama.CompressionGZIP
	//config.Producer.Retry
	return config
}

// Stop implements node.Service, terminating the monitoring and reporting daemon.
func (s *PlatonStatsService) Stop() error {
	s.stopOnce.Do(func() {
		close(s.stopSampleMsg)
		close(s.stopBlockMsg)
	})

	log.Info("PlatON stats daemon stopped")
	return nil
}

func (s *PlatonStatsService) blockMsgLoop() {
	var nextBlockNumber uint64
	nextBlockNumber = 0

	if loggedBlockNumber, err := readBlockNumber(statsLogFile); err == nil {
		nextBlockNumber = loggedBlockNumber + 1
	}

	for {
		nextBlock := s.BlockChain().GetBlockByNumber(nextBlockNumber)
		if nextBlock != nil {
			if err := s.reportBlockMsg(nextBlock); err == nil {
				if err := writeBlockNumber(statsLogFile, nextBlockNumber); err == nil {
					nextBlockNumber = nextBlockNumber + 1
				}
			}
		} else {
			time.Sleep(time.Microsecond * 100)
		}
	}
}

func (s *PlatonStatsService) reportBlockMsg(block *types.Block) error {
	var genesisData *common.GenesisData
	var receipts []*types.Receipt
	var exeBlockData *common.ExeBlockData

	var err error
	if block.NumberU64() == 0 {
		if genesisData, err = s.scanGenesis(block); err != nil {
			log.Error("cannot read genesis block", err)
			return err
		}
	} else {
		receipts = s.BlockChain().GetReceiptsByHash(block.Hash())
		exeBlockData = rawdb.ReadExeBlockData(s.ChainDb(), block.Number())
	}

	statsBlockExt := &StatsBlockExt{
		Block:        block,
		Receipts:     receipts,
		ExeBlockData: exeBlockData,
		GenesisData:  genesisData,
	}

	data, err := rlp.EncodeToBytes(statsBlockExt)
	if err != nil {
		log.Error("encode platon stats block ext message error")
		return err
	}
	// send message
	msg := &sarama.ProducerMessage{
		Topic:     kafkaBlockTopic,
		Partition: 0,
		Key:       sarama.StringEncoder(strconv.FormatUint(block.NumberU64(), 10)),
		Value:     sarama.ByteEncoder(data),
		Timestamp: time.Now(),
	}

	partition, offset, err := s.blockProducer.SendMessage(msg)

	if err != nil {
		log.Error("send block message error.", "blockNumber=", block.NumberU64(), "error", err)
	} else {
		log.Info("send block message success.", "blockNumber=", block.NumberU64(), "partition", partition, "offset", offset)
	}

	rawdb.DeleteExeBlockData(s.ChainDb(), block.Number())
	return nil
}

func readBlockNumber(file string) (uint64, error) {
	if bytes, err := ioutil.ReadFile(file); err != nil || len(bytes) == 0 {
		return 0, errors.New("Failed to read PlatON stats service log")
	} else {
		if blockNumber, err := strconv.ParseUint(string(bytes), 10, 64); err != nil {
			log.Warn("Failed to read PlatON stats service log", "error", err)
			return 0, errors.New("Failed to read PlatON stats service log")
		} else {
			return blockNumber, nil
		}
	}
}

func writeBlockNumber(file string, blockNumber uint64) error {
	return ioutil.WriteFile(file, []byte(strconv.FormatUint(blockNumber, 10)), 666)
}

func (s *PlatonStatsService) sampleMsgLoop() {
	var sampleEventProducer SampleEventProducer
	sampleEventCh := make(chan SampleEvent, sampleEventChanSize)
	sampleEventSub := sampleEventProducer.SubscribeSampleEvent(sampleEventCh)
	defer sampleEventSub.Unsubscribe()

	for {
		select {
		case sampleEvent := <-sampleEventCh:
			log.Debug("received a sample event", sampleEvent)
		case <-sampleEventSub.Err():
			return
		case <-s.stopSampleMsg:
			return
		}
	}
}

func (s *PlatonStatsService) scanGenesis(genesisBlock *types.Block) (*common.GenesisData, error) {
	genesisData := &common.GenesisData{}
	/*hash := rawdb.ReadCanonicalHash(s.eth.ChainDb(), 0)
	println("genesis block hash:", hash.String())
	block := rawdb.ReadBlock(s.eth.ChainDb(), hash, 0)
	if block == nil {
		return nil, fmt.Errorf("cannot read genesis block")
	}
	*/
	root := genesisBlock.Root()
	tr, err := trie.NewSecure(root, trie.NewDatabase(s.ChainDb()), 0)
	if err != nil {
		return nil, err
	}

	iter := tr.NodeIterator(nil)
	for iter.Next(true) {
		if iter.Leaf() {
			var obj state.Account
			err := rlp.DecodeBytes(iter.LeafBlob(), &obj)
			if err != nil {
				return nil, fmt.Errorf("parse account error:%s", err.Error())
			}
			key := iter.LeafKey()
			address := common.BytesToAddress(key)
			balance := obj.Balance.Uint64()
			genesisData.AddAllocItem(address, balance)

			log.Debug("alloc account", "address", address, "balance", balance)
		}
	}
	return genesisData, nil
}
