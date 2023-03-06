// ╔╗   ╔╗ ╔╗╔╗╔═╗╔═══╗╔═══╗    ╔═══╗                          ╔╗  ╔╗     ╔╗     ╔╗      ╔╗
// ║║   ║║ ║║║║║╔╝║╔═╗║║╔═╗║    ║╔═╗║                          ║╚╗╔╝║     ║║     ║║     ╔╝╚╗
// ║║   ║║ ║║║╚╝╝ ║╚══╗║║ ║║    ║║ ╚╝╔══╗╔═╗ ╔══╗╔══╗╔╗╔══╗    ╚╗║║╔╝╔══╗ ║║ ╔╗╔═╝║╔══╗ ╚╗╔╝╔══╗╔═╗╔══╗
// ║║ ╔╗║║ ║║║╔╗║ ╚══╗║║║ ║║    ║║╔═╗║╔╗║║╔╗╗║╔╗║║══╣╠╣║══╣     ║╚╝║ ╚ ╗║ ║║ ╠╣║╔╗║╚ ╗║  ║║ ║╔╗║║╔╝║══╣
// ║╚═╝║║╚═╝║║║║╚╗║╚═╝║║╚═╝║    ║╚╩═║║║═╣║║║║║║═╣╠══║║║╠══║     ╚╗╔╝ ║╚╝╚╗║╚╗║║║╚╝║║╚╝╚╗ ║╚╗║╚╝║║║ ╠══║
// ╚═══╝╚═══╝╚╝╚═╝╚═══╝╚═══╝    ╚═══╝╚══╝╚╝╚╝╚══╝╚══╝╚╝╚══╝      ╚╝  ╚═══╝╚═╝╚╝╚══╝╚═══╝ ╚═╝╚══╝╚╝ ╚══╝

// SPDX-License-Identifier: CC0-1.0

pragma solidity 0.8.15;

import {IERC165} from "../interfaces/IERC165.sol";
import {IERC1820Registry} from "../interfaces/IERC1820Registry.sol";
import {IDepositContract} from "../interfaces/IDepositContract.sol";

contract DepositMock is IERC165 {
    // The address of the LYXe token contract.
    address constant LYXeAddress = 0x7A2AC110202ebFdBB5dB15Ea994ba6bFbFcFc215;

    // The address of the registry contract (ERC1820 Registry).
    address constant registryAddress = 0xa5594Cd0f68eDf204A49B62eaA19Acb6376FE8Ad;

    // The hash of the interface of the contract that receives tokens.
    bytes32 constant TOKENS_RECIPIENT_INTERFACE_HASH =
    0xb281fc8c12954d22544db45de3159a39272895b169a852b314f9cc762e44c53b;

    // The depth of the Merkle tree of deposits.
    uint256 constant DEPOSIT_CONTRACT_TREE_DEPTH = 32;

    // NOTE: this also ensures `deposit_count` will fit into 64-bits
    uint256 constant MAX_DEPOSIT_COUNT = 2**DEPOSIT_CONTRACT_TREE_DEPTH - 1;

    // _to_little_endian_64(uint64(32 ether / 1 gwei))
    bytes constant amount_to_little_endian_64 = hex"0040597307000000";

    // The current state of the Merkle tree of deposits.
    bytes32[DEPOSIT_CONTRACT_TREE_DEPTH] branch;

    // A pre-computed array of zero hashes for use in computing the Merkle root.
    bytes32[DEPOSIT_CONTRACT_TREE_DEPTH] zero_hashes;

    // The current number of deposits in the contract.
    uint256 internal deposit_count;

    event DepositEvent(
        bytes pubkey,
        bytes withdrawal_credentials,
        bytes amount,
        bytes signature,
        bytes index
    );

    /**
     * @dev Storing all the deposit data which should be sliced
     * in order to get the following parameters:
     * - pubkey - the first 48 bytes
     * - withdrawal_credentials - the following 32 bytes
     * - signature - the following 96 bytes
     * - deposit_data_root - last 32 bytes
     */
    mapping(uint256 => bytes) deposit_data;

    /**
     * @dev Storing the amount of votes for each supply where the index is the initial supply of LYX in million
     */
    mapping(uint256 => uint256) public supplyVoteCounter;

    /**
     * @dev Owner of the contract
     * Has access to `freezeContract()`
     */
    address public immutable owner;

    /**
     * @dev Default value is false which allows people to send 32 LYXe
     * to this contract with valid data in order to register as Genesis Validator
     */
    bool public isContractFrozen;

    /**
     * @dev Save the deployer as the owner of the contract
     */
    constructor(address owner_) {
        owner = owner_;

        isContractFrozen = false;

        // Set this contract as the implementer of the tokens recipient interface in the registry contract.
        IERC1820Registry(registryAddress).setInterfaceImplementer(
            address(this),
            TOKENS_RECIPIENT_INTERFACE_HASH,
            address(this)
        );

        // Compute hashes in empty sparse Merkle tree
        for (uint256 height = 0; height < DEPOSIT_CONTRACT_TREE_DEPTH - 1; height++)
            zero_hashes[height + 1] = sha256(
                abi.encodePacked(zero_hashes[height], zero_hashes[height])
            );
    }

    /**
     * @dev Whenever this contract receives LYXe tokens, it must be for the reason of
     * being a Genesis Validator.
     *
     * Requirements:
     * - `amount` MUST be exactly 32 LYXe
     * - `depositData` MUST be encoded properly
     * - `depositData` MUST contain:
     *   • pubkey - the first 48 bytes
     *   • withdrawal_credentials - the following 32 bytes
     *   • signature - the following 96 bytes
     *   • deposit_data_root - last 32 bytes
     */
    function tokensReceived(
        address, /* operator */
        address, /* from */
        address, /* to */
        uint256 amount,
        bytes calldata depositData,
        bytes calldata /* operatorData */
    ) external {
        require(!isContractFrozen, "LUKSOGenesisValidatorsDepositContract: Contract is frozen");
        require(
            msg.sender == LYXeAddress,
            "LUKSOGenesisValidatorsDepositContract: Not called on LYXe transfer"
        );
        require(
            amount == 32 ether,
            "LUKSOGenesisValidatorsDepositContract: Cannot send an amount different from 32 LYXe"
        );
        // 208 = 48 bytes pubkey + 32 bytes withdrawal_credentials + 96 bytes signature + 32 bytes deposit_data_root
        require(
            depositData.length == (209),
            "LUKSOGenesisValidatorsDepositContract: depositData not encoded properly"
        );

        uint8 supply = uint8(depositData[208]);
        require(supply <= 100, "LUKSOGenesisValidatorsDepositContract: Invalid supply vote");
        supplyVoteCounter[supply]++;

        // Store the deposit data in the contract state.
        deposit_data[deposit_count] = depositData;

        // Process the deposit and update the Merkle tree.
        _deposit(
            depositData[:48], // pubkey
            depositData[48:80], // withdrawal_credentials
            depositData[80:176], // signature
            _convertBytesToBytes32(depositData[176:208]) // deposit_data_root
        );
    }

    /**
     * @dev Freze the LUKSO Genesis Deposit Contract
     */
    function freezeContract() external {
        require(msg.sender == owner, "LUKSOGenesisValidatorsDepositContract: Caller not owner");
        isContractFrozen = true;
    }

    /**
     * @dev Returns the current number of deposits.
     *
     * @return The number of deposits.
     */
    function depositCount() external view returns (uint256) {
        return deposit_count;
    }

    /**
     * @dev Retrieves an array of votes per supply and the total number of votes
     */

    function getsVotesPerSupply()
    external
    view
    returns (uint256[101] memory votesPerSupply, uint256 totalVotes)
    {
        for (uint256 i = 0; i <= 100; i++) {
            votesPerSupply[i] = supplyVoteCounter[i];
        }
        return (votesPerSupply, deposit_count);
    }

    /**
     * @dev Get an array of all excoded deposit data
     */
    function getDepositData() external view returns (bytes[] memory returnedArray) {
        returnedArray = new bytes[](deposit_count);
        for (uint256 i = 0; i < deposit_count; i++) returnedArray[i] = deposit_data[i];
    }

    /**
     * @dev Get the encoded deposit data at the `index`
     */
    function getDepositDataByIndex(uint256 index) external view returns (bytes memory) {
        return deposit_data[index];
    }

    /**
     * @dev Determines whether the contract supports a given interface.
     *
     * @param interfaceId The interface ID to check.
     * @return True if the contract supports the interface, false otherwise.
     */
    function supportsInterface(bytes4 interfaceId) external pure override returns (bool) {
        return
        interfaceId == type(IERC165).interfaceId ||
        interfaceId == type(IDepositContract).interfaceId;
    }

    /**
     * @dev Returns the current root of the Merkle tree of deposits.
     *
     * @return The Merkle root of the deposit data.
     */
    function get_deposit_root() external view returns (bytes32) {
        bytes32 node;
        uint256 size = deposit_count;
        for (uint256 height = 0; height < DEPOSIT_CONTRACT_TREE_DEPTH; height++) {
            if ((size & 1) == 1) node = sha256(abi.encodePacked(branch[height], node));
            else node = sha256(abi.encodePacked(node, zero_hashes[height]));
            size /= 2;
        }

        return
        sha256(abi.encodePacked(node, _to_little_endian_64(uint64(deposit_count)), bytes24(0)));
    }

    /**
     * @dev Returns the current number of deposits in the contract.
     *
     * @return The number of deposits in little-endian order.
     */
    function get_deposit_count() external view returns (bytes memory) {
        return _to_little_endian_64(uint64(deposit_count));
    }

    /**
     * @dev Processes a deposit and updates the Merkle tree.
     *
     * @param pubkey The public key of the depositor.
     * @param withdrawal_credentials The withdrawal credentials of the depositor.
     * @param signature The deposit signature of the depositor.
     * @param deposit_data_root The root of the deposit data.
     */
    function _deposit(
        bytes calldata pubkey,
        bytes calldata withdrawal_credentials,
        bytes calldata signature,
        bytes32 deposit_data_root
    ) internal {
        // Emit `DepositEvent` log
        emit DepositEvent(
            pubkey,
            withdrawal_credentials,
            amount_to_little_endian_64,
            signature,
            _to_little_endian_64(uint64(deposit_count))
        );

        // Compute deposit data root (`DepositData` hash tree root)
        bytes32 pubkey_root = sha256(abi.encodePacked(pubkey, bytes16(0)));

        // Compute the root of the signature data.
        bytes32 signature_root = sha256(
            abi.encodePacked(
                sha256(abi.encodePacked(signature[:64])),
                sha256(abi.encodePacked(signature[64:], bytes32(0)))
            )
        );

        // Compute the root of the deposit data.
        bytes32 node = sha256(
            abi.encodePacked(
                sha256(abi.encodePacked(pubkey_root, withdrawal_credentials)),
                sha256(abi.encodePacked(amount_to_little_endian_64, bytes24(0), signature_root))
            )
        );

        // Verify computed and expected deposit data roots match
        require(
            node == deposit_data_root,
            "LUKSOGenesisValidatorsDepositContract: reconstructed DepositData does not match supplied deposit_data_root"
        );

        // Avoid overflowing the Merkle tree (and prevent edge case in computing `branch`)
        require(
            deposit_count < MAX_DEPOSIT_COUNT,
            "LUKSOGenesisValidatorsDepositContract: merkle tree full"
        );

        // Add deposit data root to Merkle tree (update a single `branch` node)
        deposit_count += 1;
        uint256 size = deposit_count;
        for (uint256 height = 0; height < DEPOSIT_CONTRACT_TREE_DEPTH; height++) {
            if ((size & 1) == 1) {
                branch[height] = node;
                return;
            }
            node = sha256(abi.encodePacked(branch[height], node));
            size /= 2;
        }
        // As the loop should always end prematurely with the `return` statement,
        // this code should be unreachable. We assert `false` just to be safe.
        assert(false);
    }

    /**
     * @dev Converts a uint64 value to a byte array in little-endian order.
     *
     * @param value The uint64 value to convert.
     * @return ret The byte array in little-endian order.
     */
    function _to_little_endian_64(uint64 value) internal pure returns (bytes memory ret) {
        ret = new bytes(8);
        bytes8 bytesValue = bytes8(value);
        // Byteswapping during copying to bytes.
        ret[0] = bytesValue[7];
        ret[1] = bytesValue[6];
        ret[2] = bytesValue[5];
        ret[3] = bytesValue[4];
        ret[4] = bytesValue[3];
        ret[5] = bytesValue[2];
        ret[6] = bytesValue[1];
        ret[7] = bytesValue[0];
    }

    /**
     * @dev Converts the first 32 bytes of a byte array to a bytes32 value.
     *
     * @param inBytes The byte array to convert.
     * @return outBytes32 The bytes32 value.
     */
    function _convertBytesToBytes32(bytes calldata inBytes)
    internal
    pure
    returns (bytes32 outBytes32)
    {
        bytes memory memoryInBytes = inBytes;
        assembly {
            outBytes32 := mload(add(memoryInBytes, 32))
        }
    }
}
