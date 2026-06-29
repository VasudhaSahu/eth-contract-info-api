package http

import (
	stdhttp "net/http"
	"regexp"

	app "github.com/VasudhaSahu/eth-contract-info-api.git/internal/application/contract"

	"github.com/gin-gonic/gin"
)

type ContractHandler struct {
	UseCase *app.UseCase
}

func NewContractHandler(useCase *app.UseCase) *ContractHandler {
	return &ContractHandler{UseCase: useCase}
}

// we accept only canonical 0x + 40‑hex addresses here
var ethAddressRegex = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)

// GetContractInfo godoc
// @Summary 		Get smart contract info
// @Description 	Fetch deployed bytecode information for an Ethereum address using Infura. Contract address recommended; wallet address returns no code.
// @Tags 			contracts
// @Accept 			json
// @Produce 		json
// @Param 			address query string true "Ethereum address (contract address recommended)"
// @Param 			blockTag query string false "Block tag" default(latest)
// @Success 		200	 {object} 	contract.ContractInfoResponse
// @Failure 		400 {object} 	contract.ErrorResponse
// @Failure 		500 {object} 	contract.ErrorResponse
// @Router			/api/v1/contracts/info [get]
func (h *ContractHandler) GetContractInfo(c *gin.Context) {
	address := c.Query("address")
	blockTag := c.DefaultQuery("blockTag", "latest")

	// Be strict on input rather than trying to guess what the user meant
	if !ethAddressRegex.MatchString(address) {
		c.JSON(stdhttp.StatusBadRequest, app.ErrorResponse{Error: "invalid Ethereum address; expected 0x + 40 hex chars"})
		return
	}

	resp, err := h.UseCase.GetInfo(app.GetContractInfoInput{Address: address, BlockTag: blockTag})
	if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, app.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(stdhttp.StatusOK, resp)
}

// GetERC20Metadata godoc
// @Summary 		Get ERC20 token metadata
// @Description 	Fetch ERC20 name, symbol, decimals, and totalSupply using eth_call via Infura. Requires an ERC20 contract address.
// @Tags 			contracts
// @Accept 			json
// @Produce 		json
// @Param 			address query string true "ERC20 contract address"
// @Param 			blockTag query string false "Block tag" default(latest)
// @Success 		200		 {object} 	contract.ERC20MetadataResponse
// @Failure 		400		 {object} 	contract.ErrorResponse
// @Failure 		500		 {object} 	contract.ErrorResponse
// @Router			/api/v1/contracts/erc20 [get]
func (h *ContractHandler) GetERC20Metadata(c *gin.Context) {
	address := c.Query("address")
	blockTag := c.DefaultQuery("blockTag", "latest")

	// strict input check
	if !ethAddressRegex.MatchString(address) {
		c.JSON(stdhttp.StatusBadRequest, app.ErrorResponse{Error: "invalid Ethereum address; expected 0x + 40 hex chars"})
		return
	}

	// for production we should cache these calls, for now we keep it simple and stateless
	resp, err := h.UseCase.GetERC20Metadata(app.GetERC20MetadataInput{Address: address, BlockTag: blockTag})
	if err != nil {
		c.JSON(stdhttp.StatusInternalServerError, app.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(stdhttp.StatusOK, resp)
}
