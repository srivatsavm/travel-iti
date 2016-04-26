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

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	// Initialize the chaincode
	Aval, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}

	// Write the state to the ledger
	err = stub.PutState("abc", []byte(strconv.Itoa(Aval)))				//making a test var "abc", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}
	
	var empty []string
	jsonAsBytes, _ := json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	err = stub.PutState(travelItiIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	
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
	} else if function == "init_travelIti" {									//create a new travel Iti
		return t.init_travelIti(stub, args)
	}
	fmt.Println("run did not find func: " + function)						//error

	return nil, errors.New("Received unknown function invocation")
}


// ============================================================================================================================
// Query - read a variable from chaincode state - (aka read)
// ============================================================================================================================
func (t *TravelItiChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function != "query" {
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}
	var travelid, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting travlid to query")
	}

	travelid = args[0]
	valAsbytes, err := stub.GetState(travelid)									//get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + travelid + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil													//send it onward
}

func main() {
	err := shim.Start(new(TravelItiChaincode))
	if err != nil {
		fmt.Printf("Error starting Travel Iti chaincode: %s", err)
	}
}


// ============================================================================================================================
// Init TravelIti - create a new travelIti, store into chaincode state
// ============================================================================================================================
func (t *TravelItiChaincode) init_travelIti(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var err error

	
	balance, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("3rd argument must be a numeric string")
	}
	

	str := `{"travelid": "` + args[0] + `", "balance": ` + strconv.Itoa(balance) + `,"travelstate": ` + args[0] + `, "stateowner": "` + args[0] + `"}`

	err = stub.PutState(args[0], []byte(str))								//store travelIti with id as key
	if err != nil {
		return nil, err
	}
		
	//get the travelIti index
	travelItiAsBytes, err := stub.GetState(travelItiIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get travelIti index")
	}
	var travelItiIndex []string
	json.Unmarshal(travelItiAsBytes, &travelItiIndex)							//un stringify it aka JSON.parse()
	
	//append
	travelItiIndex = append(travelItiIndex, args[0])								//add travelIti name to index list
	fmt.Println("! travelIti index: ", travelItiIndex)
	jsonAsBytes, _ := json.Marshal(travelItiIndex)
	err = stub.PutState(travelItiIndexStr, jsonAsBytes)						//store name of travelIti

	fmt.Println("- end init travelIti")
	return nil, nil
}
