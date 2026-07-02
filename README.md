# eth-contract-info-api

A minimal Go backend that retrieves Ethereum contract information using Infura and exposes it through REST endpoints documented with Swagger. This project is intended as a clean portfolio sample for Go backend, Web3 integration, and API design.

## Overview

The service provides:
- Contract code inspection via Ethereum JSON-RPC
- ERC-20 token metadata lookup
- REST API endpoints built with Gin
- Swagger/OpenAPI documentation
- Environment-based configuration

## Tech Stack

- Go
- Gin
- Infura Ethereum RPC
- Swagger / OpenAPI
- DDD-inspired project structure

## Project layout

```text
eth-contract-info-api/
├── cmd/
│   └── server/
│       └── main.go              # Entry point: wires config, Infura client, use case, HTTP handlers
├── internal/
│   ├── config/
│   │   └── config.go            # Loads INFURA_API_KEY, INFURA_NETWORK, PORT from environment/.env
│   ├── infrastructure/
│   │   └── infura/
│   │       └── client.go        # JSON-RPC client for Infura implementing the Reader interface
│   ├── domain/
│   │   └── contract/
│   │       ├── model.go         # Domain entities: Address, Info, ERC20Metadata
│   │       ├── validation.go    # Domain rules: address normalization, bytecode helpers, selectors
│   │       └── decode.go        # Minimal ABI decoding helpers for ERC-20 fields
│   ├── application/
│   │   └── contract/
│   │       ├── service.go       # Service orchestrates Reader calls and applies domain rules
│   │       └── types.go         # Input and output DTOs for contract info and ERC-20 metadata
│   └── interfaces/
│       └── http/
│           └── contract_handler.go  # Gin handlers, HTTP/JSON layer, Swagger annotations
├── docs/
│   └── swagger/                 # Generated OpenAPI/Swagger docs
├── go.mod                       # Go module definition and dependencies
└── README.md                    # Overview, setup, endpoints, and usage examples
```

## Setup

1. Copy `.env.example` to `.env`
2. Set your `INFURA_API_KEY`
3. Install dependencies
4. Generate Swagger docs
5. Start the server

```bash
go mod tidy
go install github.com/swaggo/swag/cmd/swag@v1.16.3
swag init -g cmd/server/main.go
go run ./cmd/server
```

## Available Endpoints

- `GET /health`
- `GET /api/v1/contracts/info?address=0x779877a7b0d9e8603169ddbd7836e478b4624789&blockTag=latest`
- `GET /api/v1/contracts/erc20?address=0x779877a7b0d9e8603169ddbd7836e478b4624789&blockTag=latest`
- `GET /swagger/index.html`

## Example Use Cases

### Check whether an address contains deployed contract code
Call:

```http
GET /api/v1/contracts/info?address=<contract_address>&blockTag=latest
```

### Read ERC-20 token details
Call:

```http
GET /api/v1/contracts/erc20?address=<erc20_contract_address>&blockTag=latest
```

## Notes

- Use a contract address for `eth_getCode`; wallet addresses typically return `0x`.
- Use an ERC-20 contract address for the `/contracts/erc20` endpoint.
- Example Sepolia contract address:
  `0x779877a7b0d9e8603169ddbd7836e478b4624789`

## Project Goal

This project is intentionally small and focused. The goal is to demonstrate clean API design, structured Go backend development, external RPC integration, and developer-friendly documentation rather than to build a full production blockchain platform.