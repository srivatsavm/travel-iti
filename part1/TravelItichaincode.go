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

var travelItiIndexStr = "_travelItiindex"				

type TravelIti struct{
	travelid int `json:"travelid"`
	balance int `json:"balance"`
	travelstate string `json:"travelstate"`				
	stateowner string `json:"stateowner"`
	}

// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *TravelItiChaincode) init(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var Aval int
	var err error

	return nil, nil
}

// ============================================================================================================================
// Run - Our entry point
// ============================================================================================================================
func (t *TravelItiChaincode) Run(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.init(stub, args)
	}
	fmt.Println("run did not find func: " + function)						//error

	return nil, errors.New("Received unknown function invocation")
}


// ============================================================================================================================
// Query - read a variable from chaincode state - (aka read)
// ============================================================================================================================
func (t *TravelItiChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	return nil, nil													//send it onward
}

func main() {
	err := shim.Start(new(TravelItiChaincode))
	if err != nil {
		fmt.Printf("Error starting Travel Iti chaincode: %s", err)
	}
}
