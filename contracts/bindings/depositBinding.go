// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// DepositMetaData contains all meta data concerning the Deposit contract.
var DepositMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"withdrawal_credentials\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"amount\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"index\",\"type\":\"bytes\"}],\"name\":\"DepositEvent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"depositCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"freezeContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDepositData\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"returnedArray\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getDepositDataByIndex\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"get_deposit_count\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"get_deposit_root\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getsVotesPerSupply\",\"outputs\":[{\"internalType\":\"uint256[101]\",\"name\":\"votesPerSupply\",\"type\":\"uint256[101]\"},{\"internalType\":\"uint256\",\"name\":\"totalVotes\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isContractFrozen\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"supplyVoteCounter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"depositData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"tokensReceived\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// DepositABI is the input ABI used to generate the binding from.
// Deprecated: Use DepositMetaData.ABI instead.
var DepositABI = DepositMetaData.ABI

// Deposit is an auto generated Go binding around an Ethereum contract.
type Deposit struct {
	DepositCaller     // Read-only binding to the contract
	DepositTransactor // Write-only binding to the contract
	DepositFilterer   // Log filterer for contract events
}

// DepositCaller is an auto generated read-only Go binding around an Ethereum contract.
type DepositCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DepositTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DepositFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DepositSession struct {
	Contract     *Deposit          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DepositCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DepositCallerSession struct {
	Contract *DepositCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// DepositTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DepositTransactorSession struct {
	Contract     *DepositTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// DepositRaw is an auto generated low-level Go binding around an Ethereum contract.
type DepositRaw struct {
	Contract *Deposit // Generic contract binding to access the raw methods on
}

// DepositCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DepositCallerRaw struct {
	Contract *DepositCaller // Generic read-only contract binding to access the raw methods on
}

// DepositTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DepositTransactorRaw struct {
	Contract *DepositTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDeposit creates a new instance of Deposit, bound to a specific deployed contract.
func NewDeposit(address common.Address, backend bind.ContractBackend) (*Deposit, error) {
	contract, err := bindDeposit(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Deposit{DepositCaller: DepositCaller{contract: contract}, DepositTransactor: DepositTransactor{contract: contract}, DepositFilterer: DepositFilterer{contract: contract}}, nil
}

// NewDepositCaller creates a new read-only instance of Deposit, bound to a specific deployed contract.
func NewDepositCaller(address common.Address, caller bind.ContractCaller) (*DepositCaller, error) {
	contract, err := bindDeposit(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DepositCaller{contract: contract}, nil
}

// NewDepositTransactor creates a new write-only instance of Deposit, bound to a specific deployed contract.
func NewDepositTransactor(address common.Address, transactor bind.ContractTransactor) (*DepositTransactor, error) {
	contract, err := bindDeposit(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DepositTransactor{contract: contract}, nil
}

// NewDepositFilterer creates a new log filterer instance of Deposit, bound to a specific deployed contract.
func NewDepositFilterer(address common.Address, filterer bind.ContractFilterer) (*DepositFilterer, error) {
	contract, err := bindDeposit(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DepositFilterer{contract: contract}, nil
}

// bindDeposit binds a generic wrapper to an already deployed contract.
func bindDeposit(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DepositMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Deposit *DepositRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Deposit.Contract.DepositCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Deposit *DepositRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deposit.Contract.DepositTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Deposit *DepositRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Deposit.Contract.DepositTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Deposit *DepositCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Deposit.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Deposit *DepositTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deposit.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Deposit *DepositTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Deposit.Contract.contract.Transact(opts, method, params...)
}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() view returns(uint256)
func (_Deposit *DepositCaller) DepositCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "depositCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() view returns(uint256)
func (_Deposit *DepositSession) DepositCount() (*big.Int, error) {
	return _Deposit.Contract.DepositCount(&_Deposit.CallOpts)
}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() view returns(uint256)
func (_Deposit *DepositCallerSession) DepositCount() (*big.Int, error) {
	return _Deposit.Contract.DepositCount(&_Deposit.CallOpts)
}

// GetDepositData is a free data retrieval call binding the contract method 0x38a42159.
//
// Solidity: function getDepositData() view returns(bytes[] returnedArray)
func (_Deposit *DepositCaller) GetDepositData(opts *bind.CallOpts) ([][]byte, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "getDepositData")

	if err != nil {
		return *new([][]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][]byte)).(*[][]byte)

	return out0, err

}

// GetDepositData is a free data retrieval call binding the contract method 0x38a42159.
//
// Solidity: function getDepositData() view returns(bytes[] returnedArray)
func (_Deposit *DepositSession) GetDepositData() ([][]byte, error) {
	return _Deposit.Contract.GetDepositData(&_Deposit.CallOpts)
}

// GetDepositData is a free data retrieval call binding the contract method 0x38a42159.
//
// Solidity: function getDepositData() view returns(bytes[] returnedArray)
func (_Deposit *DepositCallerSession) GetDepositData() ([][]byte, error) {
	return _Deposit.Contract.GetDepositData(&_Deposit.CallOpts)
}

// GetDepositDataByIndex is a free data retrieval call binding the contract method 0xf0cd185e.
//
// Solidity: function getDepositDataByIndex(uint256 index) view returns(bytes)
func (_Deposit *DepositCaller) GetDepositDataByIndex(opts *bind.CallOpts, index *big.Int) ([]byte, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "getDepositDataByIndex", index)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetDepositDataByIndex is a free data retrieval call binding the contract method 0xf0cd185e.
//
// Solidity: function getDepositDataByIndex(uint256 index) view returns(bytes)
func (_Deposit *DepositSession) GetDepositDataByIndex(index *big.Int) ([]byte, error) {
	return _Deposit.Contract.GetDepositDataByIndex(&_Deposit.CallOpts, index)
}

// GetDepositDataByIndex is a free data retrieval call binding the contract method 0xf0cd185e.
//
// Solidity: function getDepositDataByIndex(uint256 index) view returns(bytes)
func (_Deposit *DepositCallerSession) GetDepositDataByIndex(index *big.Int) ([]byte, error) {
	return _Deposit.Contract.GetDepositDataByIndex(&_Deposit.CallOpts, index)
}

// GetDepositCount is a free data retrieval call binding the contract method 0x621fd130.
//
// Solidity: function get_deposit_count() view returns(bytes)
func (_Deposit *DepositCaller) GetDepositCount(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "get_deposit_count")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetDepositCount is a free data retrieval call binding the contract method 0x621fd130.
//
// Solidity: function get_deposit_count() view returns(bytes)
func (_Deposit *DepositSession) GetDepositCount() ([]byte, error) {
	return _Deposit.Contract.GetDepositCount(&_Deposit.CallOpts)
}

// GetDepositCount is a free data retrieval call binding the contract method 0x621fd130.
//
// Solidity: function get_deposit_count() view returns(bytes)
func (_Deposit *DepositCallerSession) GetDepositCount() ([]byte, error) {
	return _Deposit.Contract.GetDepositCount(&_Deposit.CallOpts)
}

// GetDepositRoot is a free data retrieval call binding the contract method 0xc5f2892f.
//
// Solidity: function get_deposit_root() view returns(bytes32)
func (_Deposit *DepositCaller) GetDepositRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "get_deposit_root")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetDepositRoot is a free data retrieval call binding the contract method 0xc5f2892f.
//
// Solidity: function get_deposit_root() view returns(bytes32)
func (_Deposit *DepositSession) GetDepositRoot() ([32]byte, error) {
	return _Deposit.Contract.GetDepositRoot(&_Deposit.CallOpts)
}

// GetDepositRoot is a free data retrieval call binding the contract method 0xc5f2892f.
//
// Solidity: function get_deposit_root() view returns(bytes32)
func (_Deposit *DepositCallerSession) GetDepositRoot() ([32]byte, error) {
	return _Deposit.Contract.GetDepositRoot(&_Deposit.CallOpts)
}

// GetsVotesPerSupply is a free data retrieval call binding the contract method 0x7e84159a.
//
// Solidity: function getsVotesPerSupply() view returns(uint256[101] votesPerSupply, uint256 totalVotes)
func (_Deposit *DepositCaller) GetsVotesPerSupply(opts *bind.CallOpts) (struct {
	VotesPerSupply [101]*big.Int
	TotalVotes     *big.Int
}, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "getsVotesPerSupply")

	outstruct := new(struct {
		VotesPerSupply [101]*big.Int
		TotalVotes     *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.VotesPerSupply = *abi.ConvertType(out[0], new([101]*big.Int)).(*[101]*big.Int)
	outstruct.TotalVotes = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetsVotesPerSupply is a free data retrieval call binding the contract method 0x7e84159a.
//
// Solidity: function getsVotesPerSupply() view returns(uint256[101] votesPerSupply, uint256 totalVotes)
func (_Deposit *DepositSession) GetsVotesPerSupply() (struct {
	VotesPerSupply [101]*big.Int
	TotalVotes     *big.Int
}, error) {
	return _Deposit.Contract.GetsVotesPerSupply(&_Deposit.CallOpts)
}

// GetsVotesPerSupply is a free data retrieval call binding the contract method 0x7e84159a.
//
// Solidity: function getsVotesPerSupply() view returns(uint256[101] votesPerSupply, uint256 totalVotes)
func (_Deposit *DepositCallerSession) GetsVotesPerSupply() (struct {
	VotesPerSupply [101]*big.Int
	TotalVotes     *big.Int
}, error) {
	return _Deposit.Contract.GetsVotesPerSupply(&_Deposit.CallOpts)
}

// IsContractFrozen is a free data retrieval call binding the contract method 0x131d2873.
//
// Solidity: function isContractFrozen() view returns(bool)
func (_Deposit *DepositCaller) IsContractFrozen(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "isContractFrozen")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsContractFrozen is a free data retrieval call binding the contract method 0x131d2873.
//
// Solidity: function isContractFrozen() view returns(bool)
func (_Deposit *DepositSession) IsContractFrozen() (bool, error) {
	return _Deposit.Contract.IsContractFrozen(&_Deposit.CallOpts)
}

// IsContractFrozen is a free data retrieval call binding the contract method 0x131d2873.
//
// Solidity: function isContractFrozen() view returns(bool)
func (_Deposit *DepositCallerSession) IsContractFrozen() (bool, error) {
	return _Deposit.Contract.IsContractFrozen(&_Deposit.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Deposit *DepositCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Deposit *DepositSession) Owner() (common.Address, error) {
	return _Deposit.Contract.Owner(&_Deposit.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Deposit *DepositCallerSession) Owner() (common.Address, error) {
	return _Deposit.Contract.Owner(&_Deposit.CallOpts)
}

// SupplyVoteCounter is a free data retrieval call binding the contract method 0xe0197c7e.
//
// Solidity: function supplyVoteCounter(uint256 ) view returns(uint256)
func (_Deposit *DepositCaller) SupplyVoteCounter(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "supplyVoteCounter", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SupplyVoteCounter is a free data retrieval call binding the contract method 0xe0197c7e.
//
// Solidity: function supplyVoteCounter(uint256 ) view returns(uint256)
func (_Deposit *DepositSession) SupplyVoteCounter(arg0 *big.Int) (*big.Int, error) {
	return _Deposit.Contract.SupplyVoteCounter(&_Deposit.CallOpts, arg0)
}

// SupplyVoteCounter is a free data retrieval call binding the contract method 0xe0197c7e.
//
// Solidity: function supplyVoteCounter(uint256 ) view returns(uint256)
func (_Deposit *DepositCallerSession) SupplyVoteCounter(arg0 *big.Int) (*big.Int, error) {
	return _Deposit.Contract.SupplyVoteCounter(&_Deposit.CallOpts, arg0)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool)
func (_Deposit *DepositCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool)
func (_Deposit *DepositSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Deposit.Contract.SupportsInterface(&_Deposit.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool)
func (_Deposit *DepositCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Deposit.Contract.SupportsInterface(&_Deposit.CallOpts, interfaceId)
}

// FreezeContract is a paid mutator transaction binding the contract method 0xa584a9b5.
//
// Solidity: function freezeContract() returns()
func (_Deposit *DepositTransactor) FreezeContract(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "freezeContract")
}

// FreezeContract is a paid mutator transaction binding the contract method 0xa584a9b5.
//
// Solidity: function freezeContract() returns()
func (_Deposit *DepositSession) FreezeContract() (*types.Transaction, error) {
	return _Deposit.Contract.FreezeContract(&_Deposit.TransactOpts)
}

// FreezeContract is a paid mutator transaction binding the contract method 0xa584a9b5.
//
// Solidity: function freezeContract() returns()
func (_Deposit *DepositTransactorSession) FreezeContract() (*types.Transaction, error) {
	return _Deposit.Contract.FreezeContract(&_Deposit.TransactOpts)
}

// TokensReceived is a paid mutator transaction binding the contract method 0x0023de29.
//
// Solidity: function tokensReceived(address , address , address , uint256 amount, bytes depositData, bytes ) returns()
func (_Deposit *DepositTransactor) TokensReceived(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 common.Address, amount *big.Int, depositData []byte, arg5 []byte) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "tokensReceived", arg0, arg1, arg2, amount, depositData, arg5)
}

// TokensReceived is a paid mutator transaction binding the contract method 0x0023de29.
//
// Solidity: function tokensReceived(address , address , address , uint256 amount, bytes depositData, bytes ) returns()
func (_Deposit *DepositSession) TokensReceived(arg0 common.Address, arg1 common.Address, arg2 common.Address, amount *big.Int, depositData []byte, arg5 []byte) (*types.Transaction, error) {
	return _Deposit.Contract.TokensReceived(&_Deposit.TransactOpts, arg0, arg1, arg2, amount, depositData, arg5)
}

// TokensReceived is a paid mutator transaction binding the contract method 0x0023de29.
//
// Solidity: function tokensReceived(address , address , address , uint256 amount, bytes depositData, bytes ) returns()
func (_Deposit *DepositTransactorSession) TokensReceived(arg0 common.Address, arg1 common.Address, arg2 common.Address, amount *big.Int, depositData []byte, arg5 []byte) (*types.Transaction, error) {
	return _Deposit.Contract.TokensReceived(&_Deposit.TransactOpts, arg0, arg1, arg2, amount, depositData, arg5)
}

// DepositDepositEventIterator is returned from FilterDepositEvent and is used to iterate over the raw logs and unpacked data for DepositEvent events raised by the Deposit contract.
type DepositDepositEventIterator struct {
	Event *DepositDepositEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DepositDepositEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositDepositEvent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DepositDepositEvent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DepositDepositEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositDepositEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositDepositEvent represents a DepositEvent event raised by the Deposit contract.
type DepositDepositEvent struct {
	Pubkey                []byte
	WithdrawalCredentials []byte
	Amount                []byte
	Signature             []byte
	Index                 []byte
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterDepositEvent is a free log retrieval operation binding the contract event 0x649bbc62d0e31342afea4e5cd82d4049e7e1ee912fc0889aa790803be39038c5.
//
// Solidity: event DepositEvent(bytes pubkey, bytes withdrawal_credentials, bytes amount, bytes signature, bytes index)
func (_Deposit *DepositFilterer) FilterDepositEvent(opts *bind.FilterOpts) (*DepositDepositEventIterator, error) {

	logs, sub, err := _Deposit.contract.FilterLogs(opts, "DepositEvent")
	if err != nil {
		return nil, err
	}
	return &DepositDepositEventIterator{contract: _Deposit.contract, event: "DepositEvent", logs: logs, sub: sub}, nil
}

// WatchDepositEvent is a free log subscription operation binding the contract event 0x649bbc62d0e31342afea4e5cd82d4049e7e1ee912fc0889aa790803be39038c5.
//
// Solidity: event DepositEvent(bytes pubkey, bytes withdrawal_credentials, bytes amount, bytes signature, bytes index)
func (_Deposit *DepositFilterer) WatchDepositEvent(opts *bind.WatchOpts, sink chan<- *DepositDepositEvent) (event.Subscription, error) {

	logs, sub, err := _Deposit.contract.WatchLogs(opts, "DepositEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositDepositEvent)
				if err := _Deposit.contract.UnpackLog(event, "DepositEvent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDepositEvent is a log parse operation binding the contract event 0x649bbc62d0e31342afea4e5cd82d4049e7e1ee912fc0889aa790803be39038c5.
//
// Solidity: event DepositEvent(bytes pubkey, bytes withdrawal_credentials, bytes amount, bytes signature, bytes index)
func (_Deposit *DepositFilterer) ParseDepositEvent(log types.Log) (*DepositDepositEvent, error) {
	event := new(DepositDepositEvent)
	if err := _Deposit.contract.UnpackLog(event, "DepositEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
