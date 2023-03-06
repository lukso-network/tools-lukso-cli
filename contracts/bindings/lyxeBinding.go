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

// LYXeMetaData contains all meta data concerning the LYXe contract.
var LYXeMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"tokenHolder\",\"type\":\"address\"}],\"name\":\"AuthorizedOperator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"operatorData\",\"type\":\"bytes\"}],\"name\":\"Burned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"operatorData\",\"type\":\"bytes\"}],\"name\":\"Minted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"tokenHolder\",\"type\":\"address\"}],\"name\":\"RevokedOperator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"operatorData\",\"type\":\"bytes\"}],\"name\":\"Sent\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"authorizeOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"burn\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"defaultOperators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"granularity\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenHolder\",\"type\":\"address\"}],\"name\":\"isOperatorFor\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"operatorData\",\"type\":\"bytes\"}],\"name\":\"operatorBurn\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"operatorData\",\"type\":\"bytes\"}],\"name\":\"operatorSend\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"revokeOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"send\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// LYXeABI is the input ABI used to generate the binding from.
// Deprecated: Use LYXeMetaData.ABI instead.
var LYXeABI = LYXeMetaData.ABI

// LYXe is an auto generated Go binding around an Ethereum contract.
type LYXe struct {
	LYXeCaller     // Read-only binding to the contract
	LYXeTransactor // Write-only binding to the contract
	LYXeFilterer   // Log filterer for contract events
}

// LYXeCaller is an auto generated read-only Go binding around an Ethereum contract.
type LYXeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LYXeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LYXeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LYXeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LYXeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LYXeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LYXeSession struct {
	Contract     *LYXe             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LYXeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LYXeCallerSession struct {
	Contract *LYXeCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// LYXeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LYXeTransactorSession struct {
	Contract     *LYXeTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LYXeRaw is an auto generated low-level Go binding around an Ethereum contract.
type LYXeRaw struct {
	Contract *LYXe // Generic contract binding to access the raw methods on
}

// LYXeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LYXeCallerRaw struct {
	Contract *LYXeCaller // Generic read-only contract binding to access the raw methods on
}

// LYXeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LYXeTransactorRaw struct {
	Contract *LYXeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLYXe creates a new instance of LYXe, bound to a specific deployed contract.
func NewLYXe(address common.Address, backend bind.ContractBackend) (*LYXe, error) {
	contract, err := bindLYXe(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LYXe{LYXeCaller: LYXeCaller{contract: contract}, LYXeTransactor: LYXeTransactor{contract: contract}, LYXeFilterer: LYXeFilterer{contract: contract}}, nil
}

// NewLYXeCaller creates a new read-only instance of LYXe, bound to a specific deployed contract.
func NewLYXeCaller(address common.Address, caller bind.ContractCaller) (*LYXeCaller, error) {
	contract, err := bindLYXe(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LYXeCaller{contract: contract}, nil
}

// NewLYXeTransactor creates a new write-only instance of LYXe, bound to a specific deployed contract.
func NewLYXeTransactor(address common.Address, transactor bind.ContractTransactor) (*LYXeTransactor, error) {
	contract, err := bindLYXe(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LYXeTransactor{contract: contract}, nil
}

// NewLYXeFilterer creates a new log filterer instance of LYXe, bound to a specific deployed contract.
func NewLYXeFilterer(address common.Address, filterer bind.ContractFilterer) (*LYXeFilterer, error) {
	contract, err := bindLYXe(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LYXeFilterer{contract: contract}, nil
}

// bindLYXe binds a generic wrapper to an already deployed contract.
func bindLYXe(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LYXeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LYXe *LYXeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LYXe.Contract.LYXeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LYXe *LYXeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LYXe.Contract.LYXeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LYXe *LYXeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LYXe.Contract.LYXeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LYXe *LYXeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LYXe.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LYXe *LYXeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LYXe.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LYXe *LYXeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LYXe.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_LYXe *LYXeCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LYXe.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_LYXe *LYXeSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _LYXe.Contract.BalanceOf(&_LYXe.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_LYXe *LYXeCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _LYXe.Contract.BalanceOf(&_LYXe.CallOpts, owner)
}

// DefaultOperators is a free data retrieval call binding the contract method 0x06e48538.
//
// Solidity: function defaultOperators() view returns(address[])
func (_LYXe *LYXeCaller) DefaultOperators(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _LYXe.contract.Call(opts, &out, "defaultOperators")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// DefaultOperators is a free data retrieval call binding the contract method 0x06e48538.
//
// Solidity: function defaultOperators() view returns(address[])
func (_LYXe *LYXeSession) DefaultOperators() ([]common.Address, error) {
	return _LYXe.Contract.DefaultOperators(&_LYXe.CallOpts)
}

// DefaultOperators is a free data retrieval call binding the contract method 0x06e48538.
//
// Solidity: function defaultOperators() view returns(address[])
func (_LYXe *LYXeCallerSession) DefaultOperators() ([]common.Address, error) {
	return _LYXe.Contract.DefaultOperators(&_LYXe.CallOpts)
}

// Granularity is a free data retrieval call binding the contract method 0x556f0dc7.
//
// Solidity: function granularity() view returns(uint256)
func (_LYXe *LYXeCaller) Granularity(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LYXe.contract.Call(opts, &out, "granularity")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Granularity is a free data retrieval call binding the contract method 0x556f0dc7.
//
// Solidity: function granularity() view returns(uint256)
func (_LYXe *LYXeSession) Granularity() (*big.Int, error) {
	return _LYXe.Contract.Granularity(&_LYXe.CallOpts)
}

// Granularity is a free data retrieval call binding the contract method 0x556f0dc7.
//
// Solidity: function granularity() view returns(uint256)
func (_LYXe *LYXeCallerSession) Granularity() (*big.Int, error) {
	return _LYXe.Contract.Granularity(&_LYXe.CallOpts)
}

// IsOperatorFor is a free data retrieval call binding the contract method 0xd95b6371.
//
// Solidity: function isOperatorFor(address operator, address tokenHolder) view returns(bool)
func (_LYXe *LYXeCaller) IsOperatorFor(opts *bind.CallOpts, operator common.Address, tokenHolder common.Address) (bool, error) {
	var out []interface{}
	err := _LYXe.contract.Call(opts, &out, "isOperatorFor", operator, tokenHolder)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOperatorFor is a free data retrieval call binding the contract method 0xd95b6371.
//
// Solidity: function isOperatorFor(address operator, address tokenHolder) view returns(bool)
func (_LYXe *LYXeSession) IsOperatorFor(operator common.Address, tokenHolder common.Address) (bool, error) {
	return _LYXe.Contract.IsOperatorFor(&_LYXe.CallOpts, operator, tokenHolder)
}

// IsOperatorFor is a free data retrieval call binding the contract method 0xd95b6371.
//
// Solidity: function isOperatorFor(address operator, address tokenHolder) view returns(bool)
func (_LYXe *LYXeCallerSession) IsOperatorFor(operator common.Address, tokenHolder common.Address) (bool, error) {
	return _LYXe.Contract.IsOperatorFor(&_LYXe.CallOpts, operator, tokenHolder)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_LYXe *LYXeCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LYXe.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_LYXe *LYXeSession) Name() (string, error) {
	return _LYXe.Contract.Name(&_LYXe.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_LYXe *LYXeCallerSession) Name() (string, error) {
	return _LYXe.Contract.Name(&_LYXe.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_LYXe *LYXeCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LYXe.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_LYXe *LYXeSession) Symbol() (string, error) {
	return _LYXe.Contract.Symbol(&_LYXe.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_LYXe *LYXeCallerSession) Symbol() (string, error) {
	return _LYXe.Contract.Symbol(&_LYXe.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_LYXe *LYXeCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LYXe.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_LYXe *LYXeSession) TotalSupply() (*big.Int, error) {
	return _LYXe.Contract.TotalSupply(&_LYXe.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_LYXe *LYXeCallerSession) TotalSupply() (*big.Int, error) {
	return _LYXe.Contract.TotalSupply(&_LYXe.CallOpts)
}

// AuthorizeOperator is a paid mutator transaction binding the contract method 0x959b8c3f.
//
// Solidity: function authorizeOperator(address operator) returns()
func (_LYXe *LYXeTransactor) AuthorizeOperator(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _LYXe.contract.Transact(opts, "authorizeOperator", operator)
}

// AuthorizeOperator is a paid mutator transaction binding the contract method 0x959b8c3f.
//
// Solidity: function authorizeOperator(address operator) returns()
func (_LYXe *LYXeSession) AuthorizeOperator(operator common.Address) (*types.Transaction, error) {
	return _LYXe.Contract.AuthorizeOperator(&_LYXe.TransactOpts, operator)
}

// AuthorizeOperator is a paid mutator transaction binding the contract method 0x959b8c3f.
//
// Solidity: function authorizeOperator(address operator) returns()
func (_LYXe *LYXeTransactorSession) AuthorizeOperator(operator common.Address) (*types.Transaction, error) {
	return _LYXe.Contract.AuthorizeOperator(&_LYXe.TransactOpts, operator)
}

// Burn is a paid mutator transaction binding the contract method 0xfe9d9303.
//
// Solidity: function burn(uint256 amount, bytes data) returns()
func (_LYXe *LYXeTransactor) Burn(opts *bind.TransactOpts, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _LYXe.contract.Transact(opts, "burn", amount, data)
}

// Burn is a paid mutator transaction binding the contract method 0xfe9d9303.
//
// Solidity: function burn(uint256 amount, bytes data) returns()
func (_LYXe *LYXeSession) Burn(amount *big.Int, data []byte) (*types.Transaction, error) {
	return _LYXe.Contract.Burn(&_LYXe.TransactOpts, amount, data)
}

// Burn is a paid mutator transaction binding the contract method 0xfe9d9303.
//
// Solidity: function burn(uint256 amount, bytes data) returns()
func (_LYXe *LYXeTransactorSession) Burn(amount *big.Int, data []byte) (*types.Transaction, error) {
	return _LYXe.Contract.Burn(&_LYXe.TransactOpts, amount, data)
}

// OperatorBurn is a paid mutator transaction binding the contract method 0xfc673c4f.
//
// Solidity: function operatorBurn(address account, uint256 amount, bytes data, bytes operatorData) returns()
func (_LYXe *LYXeTransactor) OperatorBurn(opts *bind.TransactOpts, account common.Address, amount *big.Int, data []byte, operatorData []byte) (*types.Transaction, error) {
	return _LYXe.contract.Transact(opts, "operatorBurn", account, amount, data, operatorData)
}

// OperatorBurn is a paid mutator transaction binding the contract method 0xfc673c4f.
//
// Solidity: function operatorBurn(address account, uint256 amount, bytes data, bytes operatorData) returns()
func (_LYXe *LYXeSession) OperatorBurn(account common.Address, amount *big.Int, data []byte, operatorData []byte) (*types.Transaction, error) {
	return _LYXe.Contract.OperatorBurn(&_LYXe.TransactOpts, account, amount, data, operatorData)
}

// OperatorBurn is a paid mutator transaction binding the contract method 0xfc673c4f.
//
// Solidity: function operatorBurn(address account, uint256 amount, bytes data, bytes operatorData) returns()
func (_LYXe *LYXeTransactorSession) OperatorBurn(account common.Address, amount *big.Int, data []byte, operatorData []byte) (*types.Transaction, error) {
	return _LYXe.Contract.OperatorBurn(&_LYXe.TransactOpts, account, amount, data, operatorData)
}

// OperatorSend is a paid mutator transaction binding the contract method 0x62ad1b83.
//
// Solidity: function operatorSend(address sender, address recipient, uint256 amount, bytes data, bytes operatorData) returns()
func (_LYXe *LYXeTransactor) OperatorSend(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int, data []byte, operatorData []byte) (*types.Transaction, error) {
	return _LYXe.contract.Transact(opts, "operatorSend", sender, recipient, amount, data, operatorData)
}

// OperatorSend is a paid mutator transaction binding the contract method 0x62ad1b83.
//
// Solidity: function operatorSend(address sender, address recipient, uint256 amount, bytes data, bytes operatorData) returns()
func (_LYXe *LYXeSession) OperatorSend(sender common.Address, recipient common.Address, amount *big.Int, data []byte, operatorData []byte) (*types.Transaction, error) {
	return _LYXe.Contract.OperatorSend(&_LYXe.TransactOpts, sender, recipient, amount, data, operatorData)
}

// OperatorSend is a paid mutator transaction binding the contract method 0x62ad1b83.
//
// Solidity: function operatorSend(address sender, address recipient, uint256 amount, bytes data, bytes operatorData) returns()
func (_LYXe *LYXeTransactorSession) OperatorSend(sender common.Address, recipient common.Address, amount *big.Int, data []byte, operatorData []byte) (*types.Transaction, error) {
	return _LYXe.Contract.OperatorSend(&_LYXe.TransactOpts, sender, recipient, amount, data, operatorData)
}

// RevokeOperator is a paid mutator transaction binding the contract method 0xfad8b32a.
//
// Solidity: function revokeOperator(address operator) returns()
func (_LYXe *LYXeTransactor) RevokeOperator(opts *bind.TransactOpts, operator common.Address) (*types.Transaction, error) {
	return _LYXe.contract.Transact(opts, "revokeOperator", operator)
}

// RevokeOperator is a paid mutator transaction binding the contract method 0xfad8b32a.
//
// Solidity: function revokeOperator(address operator) returns()
func (_LYXe *LYXeSession) RevokeOperator(operator common.Address) (*types.Transaction, error) {
	return _LYXe.Contract.RevokeOperator(&_LYXe.TransactOpts, operator)
}

// RevokeOperator is a paid mutator transaction binding the contract method 0xfad8b32a.
//
// Solidity: function revokeOperator(address operator) returns()
func (_LYXe *LYXeTransactorSession) RevokeOperator(operator common.Address) (*types.Transaction, error) {
	return _LYXe.Contract.RevokeOperator(&_LYXe.TransactOpts, operator)
}

// Send is a paid mutator transaction binding the contract method 0x9bd9bbc6.
//
// Solidity: function send(address recipient, uint256 amount, bytes data) returns()
func (_LYXe *LYXeTransactor) Send(opts *bind.TransactOpts, recipient common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _LYXe.contract.Transact(opts, "send", recipient, amount, data)
}

// Send is a paid mutator transaction binding the contract method 0x9bd9bbc6.
//
// Solidity: function send(address recipient, uint256 amount, bytes data) returns()
func (_LYXe *LYXeSession) Send(recipient common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _LYXe.Contract.Send(&_LYXe.TransactOpts, recipient, amount, data)
}

// Send is a paid mutator transaction binding the contract method 0x9bd9bbc6.
//
// Solidity: function send(address recipient, uint256 amount, bytes data) returns()
func (_LYXe *LYXeTransactorSession) Send(recipient common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _LYXe.Contract.Send(&_LYXe.TransactOpts, recipient, amount, data)
}

// LYXeAuthorizedOperatorIterator is returned from FilterAuthorizedOperator and is used to iterate over the raw logs and unpacked data for AuthorizedOperator events raised by the LYXe contract.
type LYXeAuthorizedOperatorIterator struct {
	Event *LYXeAuthorizedOperator // Event containing the contract specifics and raw log

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
func (it *LYXeAuthorizedOperatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LYXeAuthorizedOperator)
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
		it.Event = new(LYXeAuthorizedOperator)
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
func (it *LYXeAuthorizedOperatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LYXeAuthorizedOperatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LYXeAuthorizedOperator represents a AuthorizedOperator event raised by the LYXe contract.
type LYXeAuthorizedOperator struct {
	Operator    common.Address
	TokenHolder common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterAuthorizedOperator is a free log retrieval operation binding the contract event 0xf4caeb2d6ca8932a215a353d0703c326ec2d81fc68170f320eb2ab49e9df61f9.
//
// Solidity: event AuthorizedOperator(address indexed operator, address indexed tokenHolder)
func (_LYXe *LYXeFilterer) FilterAuthorizedOperator(opts *bind.FilterOpts, operator []common.Address, tokenHolder []common.Address) (*LYXeAuthorizedOperatorIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var tokenHolderRule []interface{}
	for _, tokenHolderItem := range tokenHolder {
		tokenHolderRule = append(tokenHolderRule, tokenHolderItem)
	}

	logs, sub, err := _LYXe.contract.FilterLogs(opts, "AuthorizedOperator", operatorRule, tokenHolderRule)
	if err != nil {
		return nil, err
	}
	return &LYXeAuthorizedOperatorIterator{contract: _LYXe.contract, event: "AuthorizedOperator", logs: logs, sub: sub}, nil
}

// WatchAuthorizedOperator is a free log subscription operation binding the contract event 0xf4caeb2d6ca8932a215a353d0703c326ec2d81fc68170f320eb2ab49e9df61f9.
//
// Solidity: event AuthorizedOperator(address indexed operator, address indexed tokenHolder)
func (_LYXe *LYXeFilterer) WatchAuthorizedOperator(opts *bind.WatchOpts, sink chan<- *LYXeAuthorizedOperator, operator []common.Address, tokenHolder []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var tokenHolderRule []interface{}
	for _, tokenHolderItem := range tokenHolder {
		tokenHolderRule = append(tokenHolderRule, tokenHolderItem)
	}

	logs, sub, err := _LYXe.contract.WatchLogs(opts, "AuthorizedOperator", operatorRule, tokenHolderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LYXeAuthorizedOperator)
				if err := _LYXe.contract.UnpackLog(event, "AuthorizedOperator", log); err != nil {
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

// ParseAuthorizedOperator is a log parse operation binding the contract event 0xf4caeb2d6ca8932a215a353d0703c326ec2d81fc68170f320eb2ab49e9df61f9.
//
// Solidity: event AuthorizedOperator(address indexed operator, address indexed tokenHolder)
func (_LYXe *LYXeFilterer) ParseAuthorizedOperator(log types.Log) (*LYXeAuthorizedOperator, error) {
	event := new(LYXeAuthorizedOperator)
	if err := _LYXe.contract.UnpackLog(event, "AuthorizedOperator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LYXeBurnedIterator is returned from FilterBurned and is used to iterate over the raw logs and unpacked data for Burned events raised by the LYXe contract.
type LYXeBurnedIterator struct {
	Event *LYXeBurned // Event containing the contract specifics and raw log

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
func (it *LYXeBurnedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LYXeBurned)
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
		it.Event = new(LYXeBurned)
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
func (it *LYXeBurnedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LYXeBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LYXeBurned represents a Burned event raised by the LYXe contract.
type LYXeBurned struct {
	Operator     common.Address
	From         common.Address
	Amount       *big.Int
	Data         []byte
	OperatorData []byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterBurned is a free log retrieval operation binding the contract event 0xa78a9be3a7b862d26933ad85fb11d80ef66b8f972d7cbba06621d583943a4098.
//
// Solidity: event Burned(address indexed operator, address indexed from, uint256 amount, bytes data, bytes operatorData)
func (_LYXe *LYXeFilterer) FilterBurned(opts *bind.FilterOpts, operator []common.Address, from []common.Address) (*LYXeBurnedIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _LYXe.contract.FilterLogs(opts, "Burned", operatorRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &LYXeBurnedIterator{contract: _LYXe.contract, event: "Burned", logs: logs, sub: sub}, nil
}

// WatchBurned is a free log subscription operation binding the contract event 0xa78a9be3a7b862d26933ad85fb11d80ef66b8f972d7cbba06621d583943a4098.
//
// Solidity: event Burned(address indexed operator, address indexed from, uint256 amount, bytes data, bytes operatorData)
func (_LYXe *LYXeFilterer) WatchBurned(opts *bind.WatchOpts, sink chan<- *LYXeBurned, operator []common.Address, from []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _LYXe.contract.WatchLogs(opts, "Burned", operatorRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LYXeBurned)
				if err := _LYXe.contract.UnpackLog(event, "Burned", log); err != nil {
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

// ParseBurned is a log parse operation binding the contract event 0xa78a9be3a7b862d26933ad85fb11d80ef66b8f972d7cbba06621d583943a4098.
//
// Solidity: event Burned(address indexed operator, address indexed from, uint256 amount, bytes data, bytes operatorData)
func (_LYXe *LYXeFilterer) ParseBurned(log types.Log) (*LYXeBurned, error) {
	event := new(LYXeBurned)
	if err := _LYXe.contract.UnpackLog(event, "Burned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LYXeMintedIterator is returned from FilterMinted and is used to iterate over the raw logs and unpacked data for Minted events raised by the LYXe contract.
type LYXeMintedIterator struct {
	Event *LYXeMinted // Event containing the contract specifics and raw log

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
func (it *LYXeMintedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LYXeMinted)
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
		it.Event = new(LYXeMinted)
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
func (it *LYXeMintedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LYXeMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LYXeMinted represents a Minted event raised by the LYXe contract.
type LYXeMinted struct {
	Operator     common.Address
	To           common.Address
	Amount       *big.Int
	Data         []byte
	OperatorData []byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterMinted is a free log retrieval operation binding the contract event 0x2fe5be0146f74c5bce36c0b80911af6c7d86ff27e89d5cfa61fc681327954e5d.
//
// Solidity: event Minted(address indexed operator, address indexed to, uint256 amount, bytes data, bytes operatorData)
func (_LYXe *LYXeFilterer) FilterMinted(opts *bind.FilterOpts, operator []common.Address, to []common.Address) (*LYXeMintedIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LYXe.contract.FilterLogs(opts, "Minted", operatorRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LYXeMintedIterator{contract: _LYXe.contract, event: "Minted", logs: logs, sub: sub}, nil
}

// WatchMinted is a free log subscription operation binding the contract event 0x2fe5be0146f74c5bce36c0b80911af6c7d86ff27e89d5cfa61fc681327954e5d.
//
// Solidity: event Minted(address indexed operator, address indexed to, uint256 amount, bytes data, bytes operatorData)
func (_LYXe *LYXeFilterer) WatchMinted(opts *bind.WatchOpts, sink chan<- *LYXeMinted, operator []common.Address, to []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LYXe.contract.WatchLogs(opts, "Minted", operatorRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LYXeMinted)
				if err := _LYXe.contract.UnpackLog(event, "Minted", log); err != nil {
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

// ParseMinted is a log parse operation binding the contract event 0x2fe5be0146f74c5bce36c0b80911af6c7d86ff27e89d5cfa61fc681327954e5d.
//
// Solidity: event Minted(address indexed operator, address indexed to, uint256 amount, bytes data, bytes operatorData)
func (_LYXe *LYXeFilterer) ParseMinted(log types.Log) (*LYXeMinted, error) {
	event := new(LYXeMinted)
	if err := _LYXe.contract.UnpackLog(event, "Minted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LYXeRevokedOperatorIterator is returned from FilterRevokedOperator and is used to iterate over the raw logs and unpacked data for RevokedOperator events raised by the LYXe contract.
type LYXeRevokedOperatorIterator struct {
	Event *LYXeRevokedOperator // Event containing the contract specifics and raw log

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
func (it *LYXeRevokedOperatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LYXeRevokedOperator)
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
		it.Event = new(LYXeRevokedOperator)
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
func (it *LYXeRevokedOperatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LYXeRevokedOperatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LYXeRevokedOperator represents a RevokedOperator event raised by the LYXe contract.
type LYXeRevokedOperator struct {
	Operator    common.Address
	TokenHolder common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRevokedOperator is a free log retrieval operation binding the contract event 0x50546e66e5f44d728365dc3908c63bc5cfeeab470722c1677e3073a6ac294aa1.
//
// Solidity: event RevokedOperator(address indexed operator, address indexed tokenHolder)
func (_LYXe *LYXeFilterer) FilterRevokedOperator(opts *bind.FilterOpts, operator []common.Address, tokenHolder []common.Address) (*LYXeRevokedOperatorIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var tokenHolderRule []interface{}
	for _, tokenHolderItem := range tokenHolder {
		tokenHolderRule = append(tokenHolderRule, tokenHolderItem)
	}

	logs, sub, err := _LYXe.contract.FilterLogs(opts, "RevokedOperator", operatorRule, tokenHolderRule)
	if err != nil {
		return nil, err
	}
	return &LYXeRevokedOperatorIterator{contract: _LYXe.contract, event: "RevokedOperator", logs: logs, sub: sub}, nil
}

// WatchRevokedOperator is a free log subscription operation binding the contract event 0x50546e66e5f44d728365dc3908c63bc5cfeeab470722c1677e3073a6ac294aa1.
//
// Solidity: event RevokedOperator(address indexed operator, address indexed tokenHolder)
func (_LYXe *LYXeFilterer) WatchRevokedOperator(opts *bind.WatchOpts, sink chan<- *LYXeRevokedOperator, operator []common.Address, tokenHolder []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var tokenHolderRule []interface{}
	for _, tokenHolderItem := range tokenHolder {
		tokenHolderRule = append(tokenHolderRule, tokenHolderItem)
	}

	logs, sub, err := _LYXe.contract.WatchLogs(opts, "RevokedOperator", operatorRule, tokenHolderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LYXeRevokedOperator)
				if err := _LYXe.contract.UnpackLog(event, "RevokedOperator", log); err != nil {
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

// ParseRevokedOperator is a log parse operation binding the contract event 0x50546e66e5f44d728365dc3908c63bc5cfeeab470722c1677e3073a6ac294aa1.
//
// Solidity: event RevokedOperator(address indexed operator, address indexed tokenHolder)
func (_LYXe *LYXeFilterer) ParseRevokedOperator(log types.Log) (*LYXeRevokedOperator, error) {
	event := new(LYXeRevokedOperator)
	if err := _LYXe.contract.UnpackLog(event, "RevokedOperator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LYXeSentIterator is returned from FilterSent and is used to iterate over the raw logs and unpacked data for Sent events raised by the LYXe contract.
type LYXeSentIterator struct {
	Event *LYXeSent // Event containing the contract specifics and raw log

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
func (it *LYXeSentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LYXeSent)
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
		it.Event = new(LYXeSent)
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
func (it *LYXeSentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LYXeSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LYXeSent represents a Sent event raised by the LYXe contract.
type LYXeSent struct {
	Operator     common.Address
	From         common.Address
	To           common.Address
	Amount       *big.Int
	Data         []byte
	OperatorData []byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterSent is a free log retrieval operation binding the contract event 0x06b541ddaa720db2b10a4d0cdac39b8d360425fc073085fac19bc82614677987.
//
// Solidity: event Sent(address indexed operator, address indexed from, address indexed to, uint256 amount, bytes data, bytes operatorData)
func (_LYXe *LYXeFilterer) FilterSent(opts *bind.FilterOpts, operator []common.Address, from []common.Address, to []common.Address) (*LYXeSentIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LYXe.contract.FilterLogs(opts, "Sent", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &LYXeSentIterator{contract: _LYXe.contract, event: "Sent", logs: logs, sub: sub}, nil
}

// WatchSent is a free log subscription operation binding the contract event 0x06b541ddaa720db2b10a4d0cdac39b8d360425fc073085fac19bc82614677987.
//
// Solidity: event Sent(address indexed operator, address indexed from, address indexed to, uint256 amount, bytes data, bytes operatorData)
func (_LYXe *LYXeFilterer) WatchSent(opts *bind.WatchOpts, sink chan<- *LYXeSent, operator []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LYXe.contract.WatchLogs(opts, "Sent", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LYXeSent)
				if err := _LYXe.contract.UnpackLog(event, "Sent", log); err != nil {
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

// ParseSent is a log parse operation binding the contract event 0x06b541ddaa720db2b10a4d0cdac39b8d360425fc073085fac19bc82614677987.
//
// Solidity: event Sent(address indexed operator, address indexed from, address indexed to, uint256 amount, bytes data, bytes operatorData)
func (_LYXe *LYXeFilterer) ParseSent(log types.Log) (*LYXeSent, error) {
	event := new(LYXeSent)
	if err := _LYXe.contract.UnpackLog(event, "Sent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
