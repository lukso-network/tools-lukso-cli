// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package depositMock

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

// DepositMockMetaData contains all meta data concerning the DepositMock contract.
var DepositMockMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"withdrawal_credentials\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"amount\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"index\",\"type\":\"bytes\"}],\"name\":\"DepositEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"withdrawal_credentials\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"deposit_data_root\",\"type\":\"bytes32\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"get_deposit_count\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"get_deposit_root\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// DepositMockABI is the input ABI used to generate the binding from.
// Deprecated: Use DepositMockMetaData.ABI instead.
var DepositMockABI = DepositMockMetaData.ABI

// DepositMock is an auto generated Go binding around an Ethereum contract.
type DepositMock struct {
	DepositMockCaller     // Read-only binding to the contract
	DepositMockTransactor // Write-only binding to the contract
	DepositMockFilterer   // Log filterer for contract events
}

// DepositMockCaller is an auto generated read-only Go binding around an Ethereum contract.
type DepositMockCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositMockTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DepositMockTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositMockFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DepositMockFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositMockSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DepositMockSession struct {
	Contract     *DepositMock      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DepositMockCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DepositMockCallerSession struct {
	Contract *DepositMockCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// DepositMockTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DepositMockTransactorSession struct {
	Contract     *DepositMockTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// DepositMockRaw is an auto generated low-level Go binding around an Ethereum contract.
type DepositMockRaw struct {
	Contract *DepositMock // Generic contract binding to access the raw methods on
}

// DepositMockCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DepositMockCallerRaw struct {
	Contract *DepositMockCaller // Generic read-only contract binding to access the raw methods on
}

// DepositMockTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DepositMockTransactorRaw struct {
	Contract *DepositMockTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDepositMock creates a new instance of DepositMock, bound to a specific deployed contract.
func NewDepositMock(address common.Address, backend bind.ContractBackend) (*DepositMock, error) {
	contract, err := bindDepositMock(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DepositMock{DepositMockCaller: DepositMockCaller{contract: contract}, DepositMockTransactor: DepositMockTransactor{contract: contract}, DepositMockFilterer: DepositMockFilterer{contract: contract}}, nil
}

// NewDepositMockCaller creates a new read-only instance of DepositMock, bound to a specific deployed contract.
func NewDepositMockCaller(address common.Address, caller bind.ContractCaller) (*DepositMockCaller, error) {
	contract, err := bindDepositMock(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DepositMockCaller{contract: contract}, nil
}

// NewDepositMockTransactor creates a new write-only instance of DepositMock, bound to a specific deployed contract.
func NewDepositMockTransactor(address common.Address, transactor bind.ContractTransactor) (*DepositMockTransactor, error) {
	contract, err := bindDepositMock(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DepositMockTransactor{contract: contract}, nil
}

// NewDepositMockFilterer creates a new log filterer instance of DepositMock, bound to a specific deployed contract.
func NewDepositMockFilterer(address common.Address, filterer bind.ContractFilterer) (*DepositMockFilterer, error) {
	contract, err := bindDepositMock(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DepositMockFilterer{contract: contract}, nil
}

// bindDepositMock binds a generic wrapper to an already deployed contract.
func bindDepositMock(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DepositMockMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DepositMock *DepositMockRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DepositMock.Contract.DepositMockCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DepositMock *DepositMockRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DepositMock.Contract.DepositMockTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DepositMock *DepositMockRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DepositMock.Contract.DepositMockTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DepositMock *DepositMockCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DepositMock.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DepositMock *DepositMockTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DepositMock.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DepositMock *DepositMockTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DepositMock.Contract.contract.Transact(opts, method, params...)
}

// GetDepositCount is a free data retrieval call binding the contract method 0x621fd130.
//
// Solidity: function get_deposit_count() view returns(bytes)
func (_DepositMock *DepositMockCaller) GetDepositCount(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _DepositMock.contract.Call(opts, &out, "get_deposit_count")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetDepositCount is a free data retrieval call binding the contract method 0x621fd130.
//
// Solidity: function get_deposit_count() view returns(bytes)
func (_DepositMock *DepositMockSession) GetDepositCount() ([]byte, error) {
	return _DepositMock.Contract.GetDepositCount(&_DepositMock.CallOpts)
}

// GetDepositCount is a free data retrieval call binding the contract method 0x621fd130.
//
// Solidity: function get_deposit_count() view returns(bytes)
func (_DepositMock *DepositMockCallerSession) GetDepositCount() ([]byte, error) {
	return _DepositMock.Contract.GetDepositCount(&_DepositMock.CallOpts)
}

// GetDepositRoot is a free data retrieval call binding the contract method 0xc5f2892f.
//
// Solidity: function get_deposit_root() view returns(bytes32)
func (_DepositMock *DepositMockCaller) GetDepositRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _DepositMock.contract.Call(opts, &out, "get_deposit_root")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetDepositRoot is a free data retrieval call binding the contract method 0xc5f2892f.
//
// Solidity: function get_deposit_root() view returns(bytes32)
func (_DepositMock *DepositMockSession) GetDepositRoot() ([32]byte, error) {
	return _DepositMock.Contract.GetDepositRoot(&_DepositMock.CallOpts)
}

// GetDepositRoot is a free data retrieval call binding the contract method 0xc5f2892f.
//
// Solidity: function get_deposit_root() view returns(bytes32)
func (_DepositMock *DepositMockCallerSession) GetDepositRoot() ([32]byte, error) {
	return _DepositMock.Contract.GetDepositRoot(&_DepositMock.CallOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0x22895118.
//
// Solidity: function deposit(bytes pubkey, bytes withdrawal_credentials, bytes signature, bytes32 deposit_data_root) payable returns()
func (_DepositMock *DepositMockTransactor) Deposit(opts *bind.TransactOpts, pubkey []byte, withdrawal_credentials []byte, signature []byte, deposit_data_root [32]byte) (*types.Transaction, error) {
	return _DepositMock.contract.Transact(opts, "deposit", pubkey, withdrawal_credentials, signature, deposit_data_root)
}

// Deposit is a paid mutator transaction binding the contract method 0x22895118.
//
// Solidity: function deposit(bytes pubkey, bytes withdrawal_credentials, bytes signature, bytes32 deposit_data_root) payable returns()
func (_DepositMock *DepositMockSession) Deposit(pubkey []byte, withdrawal_credentials []byte, signature []byte, deposit_data_root [32]byte) (*types.Transaction, error) {
	return _DepositMock.Contract.Deposit(&_DepositMock.TransactOpts, pubkey, withdrawal_credentials, signature, deposit_data_root)
}

// Deposit is a paid mutator transaction binding the contract method 0x22895118.
//
// Solidity: function deposit(bytes pubkey, bytes withdrawal_credentials, bytes signature, bytes32 deposit_data_root) payable returns()
func (_DepositMock *DepositMockTransactorSession) Deposit(pubkey []byte, withdrawal_credentials []byte, signature []byte, deposit_data_root [32]byte) (*types.Transaction, error) {
	return _DepositMock.Contract.Deposit(&_DepositMock.TransactOpts, pubkey, withdrawal_credentials, signature, deposit_data_root)
}

// DepositMockDepositEventIterator is returned from FilterDepositEvent and is used to iterate over the raw logs and unpacked data for DepositEvent events raised by the DepositMock contract.
type DepositMockDepositEventIterator struct {
	Event *DepositMockDepositEvent // Event containing the contract specifics and raw log

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
func (it *DepositMockDepositEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositMockDepositEvent)
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
		it.Event = new(DepositMockDepositEvent)
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
func (it *DepositMockDepositEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositMockDepositEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositMockDepositEvent represents a DepositEvent event raised by the DepositMock contract.
type DepositMockDepositEvent struct {
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
func (_DepositMock *DepositMockFilterer) FilterDepositEvent(opts *bind.FilterOpts) (*DepositMockDepositEventIterator, error) {

	logs, sub, err := _DepositMock.contract.FilterLogs(opts, "DepositEvent")
	if err != nil {
		return nil, err
	}
	return &DepositMockDepositEventIterator{contract: _DepositMock.contract, event: "DepositEvent", logs: logs, sub: sub}, nil
}

// WatchDepositEvent is a free log subscription operation binding the contract event 0x649bbc62d0e31342afea4e5cd82d4049e7e1ee912fc0889aa790803be39038c5.
//
// Solidity: event DepositEvent(bytes pubkey, bytes withdrawal_credentials, bytes amount, bytes signature, bytes index)
func (_DepositMock *DepositMockFilterer) WatchDepositEvent(opts *bind.WatchOpts, sink chan<- *DepositMockDepositEvent) (event.Subscription, error) {

	logs, sub, err := _DepositMock.contract.WatchLogs(opts, "DepositEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositMockDepositEvent)
				if err := _DepositMock.contract.UnpackLog(event, "DepositEvent", log); err != nil {
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
func (_DepositMock *DepositMockFilterer) ParseDepositEvent(log types.Log) (*DepositMockDepositEvent, error) {
	event := new(DepositMockDepositEvent)
	if err := _DepositMock.contract.UnpackLog(event, "DepositEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
