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

// EthereumDepositMetaData contains all meta data concerning the EthereumDeposit contract.
var EthereumDepositMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"withdrawal_credentials\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"amount\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"index\",\"type\":\"bytes\"}],\"name\":\"DepositEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"pubkey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"withdrawal_credentials\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"deposit_data_root\",\"type\":\"bytes32\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"get_deposit_count\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"get_deposit_root\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
}

// EthereumDepositABI is the input ABI used to generate the binding from.
// Deprecated: Use EthereumDepositMetaData.ABI instead.
var EthereumDepositABI = EthereumDepositMetaData.ABI

// EthereumDeposit is an auto generated Go binding around an Ethereum contract.
type EthereumDeposit struct {
	EthereumDepositCaller     // Read-only binding to the contract
	EthereumDepositTransactor // Write-only binding to the contract
	EthereumDepositFilterer   // Log filterer for contract events
}

// EthereumDepositCaller is an auto generated read-only Go binding around an Ethereum contract.
type EthereumDepositCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthereumDepositTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EthereumDepositTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthereumDepositFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EthereumDepositFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthereumDepositSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EthereumDepositSession struct {
	Contract     *EthereumDeposit  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EthereumDepositCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EthereumDepositCallerSession struct {
	Contract *EthereumDepositCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// EthereumDepositTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EthereumDepositTransactorSession struct {
	Contract     *EthereumDepositTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// EthereumDepositRaw is an auto generated low-level Go binding around an Ethereum contract.
type EthereumDepositRaw struct {
	Contract *EthereumDeposit // Generic contract binding to access the raw methods on
}

// EthereumDepositCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EthereumDepositCallerRaw struct {
	Contract *EthereumDepositCaller // Generic read-only contract binding to access the raw methods on
}

// EthereumDepositTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EthereumDepositTransactorRaw struct {
	Contract *EthereumDepositTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEthereumDeposit creates a new instance of EthereumDeposit, bound to a specific deployed contract.
func NewEthereumDeposit(address common.Address, backend bind.ContractBackend) (*EthereumDeposit, error) {
	contract, err := bindEthereumDeposit(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EthereumDeposit{EthereumDepositCaller: EthereumDepositCaller{contract: contract}, EthereumDepositTransactor: EthereumDepositTransactor{contract: contract}, EthereumDepositFilterer: EthereumDepositFilterer{contract: contract}}, nil
}

// NewEthereumDepositCaller creates a new read-only instance of EthereumDeposit, bound to a specific deployed contract.
func NewEthereumDepositCaller(address common.Address, caller bind.ContractCaller) (*EthereumDepositCaller, error) {
	contract, err := bindEthereumDeposit(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EthereumDepositCaller{contract: contract}, nil
}

// NewEthereumDepositTransactor creates a new write-only instance of EthereumDeposit, bound to a specific deployed contract.
func NewEthereumDepositTransactor(address common.Address, transactor bind.ContractTransactor) (*EthereumDepositTransactor, error) {
	contract, err := bindEthereumDeposit(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EthereumDepositTransactor{contract: contract}, nil
}

// NewEthereumDepositFilterer creates a new log filterer instance of EthereumDeposit, bound to a specific deployed contract.
func NewEthereumDepositFilterer(address common.Address, filterer bind.ContractFilterer) (*EthereumDepositFilterer, error) {
	contract, err := bindEthereumDeposit(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EthereumDepositFilterer{contract: contract}, nil
}

// bindEthereumDeposit binds a generic wrapper to an already deployed contract.
func bindEthereumDeposit(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EthereumDepositMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EthereumDeposit *EthereumDepositRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EthereumDeposit.Contract.EthereumDepositCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EthereumDeposit *EthereumDepositRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthereumDeposit.Contract.EthereumDepositTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EthereumDeposit *EthereumDepositRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EthereumDeposit.Contract.EthereumDepositTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EthereumDeposit *EthereumDepositCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EthereumDeposit.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EthereumDeposit *EthereumDepositTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthereumDeposit.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EthereumDeposit *EthereumDepositTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EthereumDeposit.Contract.contract.Transact(opts, method, params...)
}

// GetDepositCount is a free data retrieval call binding the contract method 0x621fd130.
//
// Solidity: function get_deposit_count() view returns(bytes)
func (_EthereumDeposit *EthereumDepositCaller) GetDepositCount(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _EthereumDeposit.contract.Call(opts, &out, "get_deposit_count")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetDepositCount is a free data retrieval call binding the contract method 0x621fd130.
//
// Solidity: function get_deposit_count() view returns(bytes)
func (_EthereumDeposit *EthereumDepositSession) GetDepositCount() ([]byte, error) {
	return _EthereumDeposit.Contract.GetDepositCount(&_EthereumDeposit.CallOpts)
}

// GetDepositCount is a free data retrieval call binding the contract method 0x621fd130.
//
// Solidity: function get_deposit_count() view returns(bytes)
func (_EthereumDeposit *EthereumDepositCallerSession) GetDepositCount() ([]byte, error) {
	return _EthereumDeposit.Contract.GetDepositCount(&_EthereumDeposit.CallOpts)
}

// GetDepositRoot is a free data retrieval call binding the contract method 0xc5f2892f.
//
// Solidity: function get_deposit_root() view returns(bytes32)
func (_EthereumDeposit *EthereumDepositCaller) GetDepositRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _EthereumDeposit.contract.Call(opts, &out, "get_deposit_root")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetDepositRoot is a free data retrieval call binding the contract method 0xc5f2892f.
//
// Solidity: function get_deposit_root() view returns(bytes32)
func (_EthereumDeposit *EthereumDepositSession) GetDepositRoot() ([32]byte, error) {
	return _EthereumDeposit.Contract.GetDepositRoot(&_EthereumDeposit.CallOpts)
}

// GetDepositRoot is a free data retrieval call binding the contract method 0xc5f2892f.
//
// Solidity: function get_deposit_root() view returns(bytes32)
func (_EthereumDeposit *EthereumDepositCallerSession) GetDepositRoot() ([32]byte, error) {
	return _EthereumDeposit.Contract.GetDepositRoot(&_EthereumDeposit.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool)
func (_EthereumDeposit *EthereumDepositCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _EthereumDeposit.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool)
func (_EthereumDeposit *EthereumDepositSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _EthereumDeposit.Contract.SupportsInterface(&_EthereumDeposit.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) pure returns(bool)
func (_EthereumDeposit *EthereumDepositCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _EthereumDeposit.Contract.SupportsInterface(&_EthereumDeposit.CallOpts, interfaceId)
}

// Deposit is a paid mutator transaction binding the contract method 0x22895118.
//
// Solidity: function deposit(bytes pubkey, bytes withdrawal_credentials, bytes signature, bytes32 deposit_data_root) payable returns()
func (_EthereumDeposit *EthereumDepositTransactor) Deposit(opts *bind.TransactOpts, pubkey []byte, withdrawal_credentials []byte, signature []byte, deposit_data_root [32]byte) (*types.Transaction, error) {
	return _EthereumDeposit.contract.Transact(opts, "deposit", pubkey, withdrawal_credentials, signature, deposit_data_root)
}

// Deposit is a paid mutator transaction binding the contract method 0x22895118.
//
// Solidity: function deposit(bytes pubkey, bytes withdrawal_credentials, bytes signature, bytes32 deposit_data_root) payable returns()
func (_EthereumDeposit *EthereumDepositSession) Deposit(pubkey []byte, withdrawal_credentials []byte, signature []byte, deposit_data_root [32]byte) (*types.Transaction, error) {
	return _EthereumDeposit.Contract.Deposit(&_EthereumDeposit.TransactOpts, pubkey, withdrawal_credentials, signature, deposit_data_root)
}

// Deposit is a paid mutator transaction binding the contract method 0x22895118.
//
// Solidity: function deposit(bytes pubkey, bytes withdrawal_credentials, bytes signature, bytes32 deposit_data_root) payable returns()
func (_EthereumDeposit *EthereumDepositTransactorSession) Deposit(pubkey []byte, withdrawal_credentials []byte, signature []byte, deposit_data_root [32]byte) (*types.Transaction, error) {
	return _EthereumDeposit.Contract.Deposit(&_EthereumDeposit.TransactOpts, pubkey, withdrawal_credentials, signature, deposit_data_root)
}

// EthereumDepositDepositEventIterator is returned from FilterDepositEvent and is used to iterate over the raw logs and unpacked data for DepositEvent events raised by the EthereumDeposit contract.
type EthereumDepositDepositEventIterator struct {
	Event *EthereumDepositDepositEvent // Event containing the contract specifics and raw log

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
func (it *EthereumDepositDepositEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthereumDepositDepositEvent)
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
		it.Event = new(EthereumDepositDepositEvent)
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
func (it *EthereumDepositDepositEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthereumDepositDepositEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthereumDepositDepositEvent represents a DepositEvent event raised by the EthereumDeposit contract.
type EthereumDepositDepositEvent struct {
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
func (_EthereumDeposit *EthereumDepositFilterer) FilterDepositEvent(opts *bind.FilterOpts) (*EthereumDepositDepositEventIterator, error) {

	logs, sub, err := _EthereumDeposit.contract.FilterLogs(opts, "DepositEvent")
	if err != nil {
		return nil, err
	}
	return &EthereumDepositDepositEventIterator{contract: _EthereumDeposit.contract, event: "DepositEvent", logs: logs, sub: sub}, nil
}

// WatchDepositEvent is a free log subscription operation binding the contract event 0x649bbc62d0e31342afea4e5cd82d4049e7e1ee912fc0889aa790803be39038c5.
//
// Solidity: event DepositEvent(bytes pubkey, bytes withdrawal_credentials, bytes amount, bytes signature, bytes index)
func (_EthereumDeposit *EthereumDepositFilterer) WatchDepositEvent(opts *bind.WatchOpts, sink chan<- *EthereumDepositDepositEvent) (event.Subscription, error) {

	logs, sub, err := _EthereumDeposit.contract.WatchLogs(opts, "DepositEvent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthereumDepositDepositEvent)
				if err := _EthereumDeposit.contract.UnpackLog(event, "DepositEvent", log); err != nil {
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
func (_EthereumDeposit *EthereumDepositFilterer) ParseDepositEvent(log types.Log) (*EthereumDepositDepositEvent, error) {
	event := new(EthereumDepositDepositEvent)
	if err := _EthereumDeposit.contract.UnpackLog(event, "DepositEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
