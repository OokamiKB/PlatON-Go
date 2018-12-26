package vm

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"bytes"
	"encoding/json"
	"errors"
	"math/big"
)

// error def
var (
	ErrOwnerNotOnly     = errors.New("Node ID cannot bind multiple owners")
	ErrPermissionDenied = errors.New("Transaction from address permission denied")
	ErrDepositEmpty     = errors.New("Deposit balance not zero")
	ErrWithdrawEmpty    = errors.New("No withdrawal amount")
	ErrCandidateEmpty   = errors.New("CandidatePool is null")
)

const (
	CandidateDepositEvent       = "CandidateDepositEvent"
	CandidateApplyWithdrawEvent = "CandidateApplyWithdrawEvent"
	CandidateWithdrawEvent      = "CandidateWithdrawEvent"
	SetCandidateExtraEvent      = "SetCandidateExtraEvent"
)

type candidatePool interface {
	SetCandidate(state StateDB, nodeId discover.NodeID, can *types.Candidate) error
	GetCandidate(state StateDB, nodeId discover.NodeID) (*types.Candidate, error)
	GetCandidateArr(state StateDB, nodeIds ...discover.NodeID) (types.CandidateQueue, error)
	WithdrawCandidate(state StateDB, nodeId discover.NodeID, price, blockNumber *big.Int) error
	GetChosens(state StateDB, flag int) types.CandidateQueue
	GetChairpersons(state StateDB) types.CandidateQueue
	GetDefeat(state StateDB, nodeId discover.NodeID) (types.CandidateQueue, error)
	IsDefeat(state StateDB, nodeId discover.NodeID) (bool, error)
	RefundBalance(state StateDB, nodeId discover.NodeID, blockNumber *big.Int) error
	GetOwner(state StateDB, nodeId discover.NodeID) common.Address
	SetCandidateExtra(state StateDB, nodeId discover.NodeID, extra string) error
	GetRefundInterval() uint64
}

type CandidateContract struct {
	Contract *Contract
	Evm      *EVM
}

func (c *CandidateContract) RequiredGas(input []byte) uint64 {
	return params.EcrecoverGas
}

func (c *CandidateContract) Run(input []byte) ([]byte, error) {
	if c.Evm.CandidatePool == nil {
		log.Error("Run==> ", "ErrCandidateEmpty", ErrCandidateEmpty.Error())
		return nil, ErrCandidateEmpty
	}
	var command = map[string]interface{}{
		"CandidateDetails":        c.CandidateDetails,
		"CandidateApplyWithdraw":  c.CandidateApplyWithdraw,
		"CandidateDeposit":        c.CandidateDeposit,
		"CandidateList":           c.CandidateList,
		"CandidateWithdraw":       c.CandidateWithdraw,
		"SetCandidateExtra":       c.SetCandidateExtra,
		"CandidateWithdrawInfos":  c.CandidateWithdrawInfos,
		"VerifiersList":           c.VerifiersList,
		"GetBatchCandidateDetail": c.GetBatchCandidateDetail,
	}
	return execute(input, command)
}

// Candidate Application && Increase Quality Deposit
func (c *CandidateContract) CandidateDeposit(nodeId discover.NodeID, owner common.Address, fee uint64, host, port, extra string) ([]byte, error) {
	// debug
	deposit := c.Contract.value
	txHash := c.Evm.StateDB.TxHash()
	txIdx := c.Evm.StateDB.TxIdx()
	height := c.Evm.Context.BlockNumber
	from := c.Contract.caller.Address()
	log.Info("CandidateDeposit==> ", "nodeId: ", nodeId.String(), " owner: ", owner.Hex(), " deposit: ", deposit,
		"  fee: ", fee, " txhash: ", txHash.Hex(), " txIdx: ", txIdx, " height: ", height, " from: ", from.Hex(),
		" host: ", host, " port: ", port, " extra: ", extra)
	//todo
	if deposit.Cmp(big.NewInt(0)) < 1 {
		r := ResultCommon{false, ErrDepositEmpty.Error()}
		event, _ := json.Marshal(r)
		c.addLog(CandidateDepositEvent, string(event))
		return nil, ErrDepositEmpty
	}
	can, err := c.Evm.CandidatePool.GetCandidate(c.Evm.StateDB, nodeId)
	if nil != err {
		log.Error("CandidateDeposit==> ", "err!=nill: ", err.Error())
		r := ResultCommon{false, err.Error()}
		event, _ := json.Marshal(r)
		c.addLog(CandidateDepositEvent, string(event))
		return nil, err
	}
	var alldeposit *big.Int
	if nil != can {
		if ok := bytes.Equal(can.Owner.Bytes(), owner.Bytes()); !ok {
			log.Error(ErrOwnerNotOnly.Error())
			r := ResultCommon{false, ErrOwnerNotOnly.Error()}
			event, _ := json.Marshal(r)
			c.addLog(CandidateDepositEvent, string(event))
			return nil, ErrOwnerNotOnly
		}
		alldeposit = new(big.Int).Add(can.Deposit, deposit)
		log.Info("CandidateDeposit==> ", "alldeposit: ", alldeposit, " can.Deposit: ", can.Deposit, " deposit: ", deposit)
	} else {
		alldeposit = deposit
	}
	canDeposit := types.Candidate{
		alldeposit,
		height,
		txIdx,
		nodeId,
		host,
		port,
		owner,
		from,
		extra,
		fee,
		//0,
		//new(big.Int).SetUint64(0),
		common.Hash{},
	}
	log.Info("CandidateDeposit==> ", "canDeposit: ", canDeposit)
	if err = c.Evm.CandidatePool.SetCandidate(c.Evm.StateDB, nodeId, &canDeposit); err != nil {
		// rollback transaction
		// ......
		log.Error("CandidateDeposit==> ", "SetCandidate return err: ", err.Error())
		r := ResultCommon{false, err.Error()}
		event, _ := json.Marshal(r)
		c.addLog(CandidateDepositEvent, string(event))
		return nil, err
	}
	r := ResultCommon{true, "success"}
	event, _ := json.Marshal(r)
	c.addLog(CandidateDepositEvent, string(event))
	log.Info("CandidateDeposit==> ", "json: ", string(event))
	return nil, nil
}

// Apply for a refund of the deposit
func (c *CandidateContract) CandidateApplyWithdraw(nodeId discover.NodeID, withdraw *big.Int) ([]byte, error) {
	// debug
	txHash := c.Evm.StateDB.TxHash()
	from := c.Contract.caller.Address()
	height := c.Evm.Context.BlockNumber
	log.Info("CandidateApplyWithdraw==> ", "nodeId: ", nodeId.String(), " from: ", from.Hex(), " txHash: ", txHash.Hex(), " withdraw: ", withdraw, " height: ", height)
	// todo
	can, err := c.Evm.CandidatePool.GetCandidate(c.Evm.StateDB, nodeId)
	if err != nil {
		log.Error("CandidateApplyWithdraw==> ", "err!=nill: ", err.Error())
		r := ResultCommon{false, err.Error()}
		event, _ := json.Marshal(r)
		c.addLog(CandidateApplyWithdrawEvent, string(event))
		return nil, err
	}
	if can.Deposit.Cmp(big.NewInt(0)) < 1 {
		r := ResultCommon{false, ErrWithdrawEmpty.Error()}
		event, _ := json.Marshal(r)
		c.addLog(CandidateApplyWithdrawEvent, string(event))
		return nil, ErrWithdrawEmpty
	}
	if ok := bytes.Equal(can.Owner.Bytes(), from.Bytes()); !ok {
		log.Error(ErrPermissionDenied.Error())
		r := ResultCommon{false, ErrPermissionDenied.Error()}
		event, _ := json.Marshal(r)
		c.addLog(CandidateApplyWithdrawEvent, string(event))
		return nil, ErrPermissionDenied
	}
	if withdraw.Cmp(can.Deposit) > 0 {
		withdraw = can.Deposit
	}
	if err := c.Evm.CandidatePool.WithdrawCandidate(c.Evm.StateDB, nodeId, withdraw, height); nil != err {
		log.Error(err.Error())
		r := ResultCommon{false, err.Error()}
		event, _ := json.Marshal(r)
		c.addLog(CandidateApplyWithdrawEvent, string(event))
		return nil, err
	}
	r := ResultCommon{true, "success"}
	event, _ := json.Marshal(r)
	c.addLog(CandidateApplyWithdrawEvent, string(event))
	log.Info("CandidateApplyWithdraw==> ", "json: ", string(event))
	return nil, nil
}

// Deposit withdrawal
func (c *CandidateContract) CandidateWithdraw(nodeId discover.NodeID) ([]byte, error) {
	// debug
	txHash := c.Evm.StateDB.TxHash()
	height := c.Evm.Context.BlockNumber
	log.Info("CandidateWithdraw==> ", "nodeId: ", nodeId.String(), " height: ", height, " txHash: ", txHash.Hex())
	// todo
	if err := c.Evm.CandidatePool.RefundBalance(c.Evm.StateDB, nodeId, height); nil != err {
		log.Error(err.Error())
		r := ResultCommon{false, err.Error()}
		event, _ := json.Marshal(r)
		c.addLog(CandidateWithdrawEvent, string(event))
		return nil, err
	}
	// return
	r := ResultCommon{true, "success"}
	event, _ := json.Marshal(r)
	c.addLog(CandidateWithdrawEvent, string(event))
	log.Info("CandidateWithdraw==> ", "json: ", string(event))
	return nil, nil
}

// Get the refund history you have applied for
func (c *CandidateContract) CandidateWithdrawInfos(nodeId discover.NodeID) ([]byte, error) {
	// debug
	log.Info("CandidateWithdrawInfos==> ", "nodeId: ", nodeId.String())
	// todo
	infos, err := c.Evm.CandidatePool.GetDefeat(c.Evm.StateDB, nodeId)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	// return
	type WithdrawInfo struct {
		Balance        *big.Int
		LockNumber     *big.Int
		LockBlockCycle uint64
	}
	type WithdrawInfos struct {
		Ret    bool
		ErrMsg string
		Infos  []WithdrawInfo
	}
	r := WithdrawInfos{true, "success", make([]WithdrawInfo, len(infos))}
	for i, v := range infos {
		r.Infos[i] = WithdrawInfo{v.Deposit, v.BlockNumber, c.Evm.CandidatePool.GetRefundInterval()}
	}
	data, _ := json.Marshal(r)
	sdata := DecodeResultStr(string(data))
	log.Info("CandidateWithdrawInfos==> ", "json: ", string(data))
	return sdata, nil
}

// Set up additional information
func (c *CandidateContract) SetCandidateExtra(nodeId discover.NodeID, extra string) ([]byte, error) {
	// debug
	txHash := c.Evm.StateDB.TxHash()
	from := c.Contract.caller.Address()
	log.Info("SetCandidate==> ", "nodeId: ", nodeId.String(), " extra: ", extra, " from: ", from.Hex(), " txHash: ", txHash.Hex())
	// todo
	owner := c.Evm.CandidatePool.GetOwner(c.Evm.StateDB, nodeId)
	if ok := bytes.Equal(owner.Bytes(), from.Bytes()); !ok {
		log.Error(ErrPermissionDenied.Error())
		r := ResultCommon{false, ErrPermissionDenied.Error()}
		event, _ := json.Marshal(r)
		c.addLog(SetCandidateExtraEvent, string(event))
		return nil, ErrPermissionDenied
	}
	if err := c.Evm.CandidatePool.SetCandidateExtra(c.Evm.StateDB, nodeId, extra); nil != err {
		log.Error(err.Error())
		r := ResultCommon{false, err.Error()}
		event, _ := json.Marshal(r)
		c.addLog(SetCandidateExtraEvent, string(event))
		return nil, err
	}
	r := ResultCommon{true, "success"}
	event, _ := json.Marshal(r)
	c.addLog(SetCandidateExtraEvent, string(event))
	log.Info("SetCandidate==> ", "json: ", string(event))
	return nil, nil
}

// Get candidate details
func (c *CandidateContract) CandidateDetails(nodeId discover.NodeID) ([]byte, error) {
	log.Info("CandidateDetails==> ", "nodeId: ", nodeId.String())
	candidate, err := c.Evm.CandidatePool.GetCandidate(c.Evm.StateDB, nodeId)
	if nil != err {
		log.Error("CandidateDetails==> ", "get CandidateDetails() occured error: ", err.Error())
		return nil, err
	}
	if nil == candidate {
		log.Error("CandidateDetails==> The candidate for the inquiry does not exist")
		return nil, nil
	}
	data, _ := json.Marshal(candidate)
	sdata := DecodeResultStr(string(data))
	log.Info("CandidateDetails==> ", "json: ", string(data), " []byte: ", sdata)
	return sdata, nil
}

// GetBatchCandidateDetail returns the batch of candidate info.
func (c *CandidateContract) GetBatchCandidateDetail(nodeIds []discover.NodeID) ([]byte, error) {
	input, _ := json.Marshal(nodeIds)
	log.Info("GetBatchCandidateDetail==>", "length: ", len(nodeIds), " nodeIds: ", string(input))
	candidates, err := c.Evm.CandidatePool.GetCandidateArr(c.Evm.StateDB, nodeIds...)
	if nil != err {
		log.Error("GetBatchCandidateDetail==> ", "get GetBatchCandidateDetail() occured error: ", err.Error())
		return nil, err
	}
	if nil == candidates {
		log.Error("GetBatchCandidateDetail==> The candidates for the inquiry does not exist")
		return nil, nil
	}
	data, _ := json.Marshal(candidates)
	sdata := DecodeResultStr(string(data))
	log.Info("GetBatchCandidateDetail==> ", "json: ", string(data), " []byte: ", sdata)
	return sdata, nil
}

// Get the current block candidate list
func (c *CandidateContract) CandidateList() ([]byte, error) {
	log.Info("CandidateList==> into func CandidateList... ")
	arr := c.Evm.CandidatePool.GetChosens(c.Evm.StateDB, 0)
	if nil == arr {
		log.Error("CandidateList==> The candidateList for the inquiry does not exist")
		return nil, nil
	}
	data, _ := json.Marshal(arr)
	sdata := DecodeResultStr(string(data))
	log.Info("CandidateList==>", "json: ", string(data), " []byte: ", sdata)
	return sdata, nil
}

// Get the current block round certifier list
func (c *CandidateContract) VerifiersList() ([]byte, error) {
	log.Info("VerifiersList==> into func VerifiersList... ")
	arr := c.Evm.CandidatePool.GetChairpersons(c.Evm.StateDB)
	if nil == arr {
		log.Error("VerifiersList==> The verifiersList for the inquiry does not exist")
		return nil, nil
	}
	data, _ := json.Marshal(arr)
	sdata := DecodeResultStr(string(data))
	log.Info("VerifiersList==> ", "json: ", string(data), " []byte: ", sdata)
	return sdata, nil
}

// transaction add event
func (c *CandidateContract) addLog(event, data string) {
	var logdata [][]byte
	logdata = make([][]byte, 0)
	logdata = append(logdata, []byte(data))
	buf := new(bytes.Buffer)
	if err := rlp.Encode(buf, logdata); nil != err {
		log.Error("addlog==> ", "rlp encode fail: ", err.Error())
	}
	c.Evm.StateDB.AddLog(&types.Log{
		Address:     common.CandidatePoolAddr,
		Topics:      []common.Hash{common.BytesToHash(crypto.Keccak256([]byte(event)))},
		Data:        buf.Bytes(),
		BlockNumber: c.Evm.Context.BlockNumber.Uint64(),
	})
}