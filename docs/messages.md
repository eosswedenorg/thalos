# Messages

This document describes the different messages that are sent

## Encoding

All messages are encoded in `json` format

## Types

### HearthBeat

Heartbeat messages are posted to the hearthbeat channel periodically.

| Field                      | Datatype | Description                                 |
| -------------------------- | -------- | ------------------------------------------- |
| blocknum                   | int      | Current block number                        |
| head_blocknum              | int      | Head block number                           |
| last_irreversible_blocknum | int      | block number of the last irreversible block |

### Transaction


### ActionTrace

| Field    | Datatype | Description                                                       |
| -------- | -------- | ----------------------------------------------------------------- |
| tx_id    | string   | Transaction ID                                                    |
| blocknum | int      | Block number where this action trace (and transaction) belongs to |
| receiver | string   | Receiver account                                                  |
| contract | string   | Contract account                                                  |
| action   | string   | What action was executed on the contract                          |
| data     | any      | Contract specific data (decoded using the contracts abi)          |
| hex_data | string   | Contract specific data (undecoded hex)                            |
