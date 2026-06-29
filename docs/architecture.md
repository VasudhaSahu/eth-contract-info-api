# Architecture

## Overview

`eth-contract-info-api` is a small Go HTTP service that exposes Ethereum contract information and ERC‑20 metadata through a REST API. It is intentionally focused: one RPC backend (Infura), two main endpoints, and a simple layered design.

## Layers

- **Interfaces (HTTP)**  
  `internal/interfaces/http/contract_handler.go` defines Gin handlers for `/api/v1/contracts/info` and `/api/v1/contracts/erc20`. Handlers validate input, delegate to the use case layer, and translate domain responses into JSON.

- **Application / Use Case**  
  `internal/application/contract/usecase.go` contains `UseCase`, which orchestrates calls to the RPC client and applies domain rules (normalize addresses, classify bytecode, decode ERC‑20 metadata). It depends only on a small `Reader` interface.

- **Domain**  
  `internal/domain/contract/*` defines domain types (`Address`, `Info`, `ERC20Metadata`) and rules (`NormalizeAddress`, `HasContractCode`, `BytecodeSize`, `BytecodePreview`, ABI decoding helpers). These functions encapsulate Ethereum‑specific logic and keep it away from transport details.

- **Infrastructure**  
  `internal/infrastructure/infura/client.go` implements the `Reader` interface using Infura's JSON‑RPC API. It knows how to build and parse RPC requests and responses, but does not contain business logic.

- **Configuration**  
  `internal/config` loads environment variables such as `INFURA_API_KEY`, `INFURA_NETWORK`, and `PORT`. The `main.go` entrypoint wires configuration into the infrastructure and application layers.

## Request Flow

1. Client calls `/api/v1/contracts/info` or `/api/v1/contracts/erc20` with `address` and optional `blockTag`.
2. The Gin handler validates `address` and builds an input DTO.
3. The use case calls the Infura client to fetch contract bytecode or ERC‑20 fields via JSON‑RPC.
4. Domain helpers classify the code, decode ABI‑encoded fields, and build a response DTO.
5. The handler returns a JSON response to the caller and appropriate HTTP status codes.

## Error Handling

- Input validation errors (e.g., malformed address) return `400 Bad Request` with an `ErrorResponse`.
- RPC or network failures bubble up as `500 Internal Server Error` with a concise error message.
- The main entrypoint fails fast if configuration is invalid or the HTTP server cannot start.

## Limitations and Future Work

This project deliberately stays small and does not attempt to be production‑ready. Potential extensions include:

- Adding caching for common contract addresses.
- Using `context.Context` throughout for request timeouts and cancellation.
- Introducing structured logging and metrics.
- Supporting additional Ethereum networks and providers behind the same `Reader` abstraction.