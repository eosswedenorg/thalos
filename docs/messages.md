# Messages

This document describes the different messages that are sent

## Encoding

All messages are encoded in `json` format

## Types

### HeartBeat

Heartbeat messages are posted to the heartbeat channel periodically.

| Field                      | Datatype | Description                                 |
| -------------------------- | -------- | ------------------------------------------- |
| blocknum                   | int      | Current block number                        |
| head_blocknum              | int      | Head block number                           |
| last_irreversible_blocknum | int      | block number of the last irreversible block |

### Transaction


### ActionTrace

| Field          | Datatype          | Description                                                       |
| -------------- | ----------------- | ----------------------------------------------------------------- |
| tx_id          | string            | Transaction ID                                                    |
| blocknum       | int               | Block number where this action trace (and transaction) belongs to |
| blocktimestamp | time              | Block timestamp                                                   |
| receipt        | ActionReceipt     | Action receipt                                                    |
| receiver       | string            | Receiver account                                                  |
| contract       | string            | Contract account                                                  |
| action         | string            | What action was executed on the contract                          |
| data           | any               | Contract specific data (decoded using the contracts abi)          |
| authorization  | PermissionLevel[] | Authorization                                                     |

### ActionReceipt

| Field           | Datatype              | Description        |
| --------------- | --------------------- | ------------------ |
| receiver        | string                | Actor account name |
| act_digest      | string                | Action digest      |
| global_sequence | int                   | Global sequence    |
| recv_sequence   | int                   | Receive sequence   |
| auth_sequence   | AccountAuthSequence[] | Auth sequence      |
| code_sequence   | int                   | Code sequence      |
| abi_sequence    | int                   | ABI sequence       |

### PermisssionLevel

| Field      | Datatype | Description                      |
| ---------- | -------- | -------------------------------- |
| actor      | string   | Actor account name               |
| permission | string   | Permission (for example: active) |

### AccountAuthSequence

| Field    | Datatype | Description  |
| -------- | -------- | ------------ |
| account  | string   | Account name |
| sequence | int      | Sequence     |
