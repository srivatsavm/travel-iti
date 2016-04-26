package main

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	"strings"

	"github.com/openblockchain/obc-peer/openchain/chaincode/shim"
)

type TravelItiChaincode struct {
}

type TravelIti struct{}

// ============================================================================================================================
// Run - Our entry point
// ============================================================================================================================
func (t *TravelItiChaincode) Run(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	return nil, nil
}

func main() {
	err := shim.Start(new(TravelItiChaincode))
	if err != nil {
		fmt.Printf("Error starting Travel Iti chaincode: %s", err)
	}
}
