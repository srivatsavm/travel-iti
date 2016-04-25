/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	"strings"

	"github.com/openblockchain/obc-peer/openchain/chaincode/shim"
)

// TravelItiChaincode example simple Chaincode implementation
type TravelItiChaincode struct {
}

var travelItiIndexStr = "_travelItiindex"				//name for the key/value that will store a list of all known travel itineraries
var openTradesStr = "_opentrades"				//name for the key/value that will store all open trades

type TravelIti struct{
	Name string `json:"name"`					//the fieldtags are needed to keep case from bouncing around
	Color string `json:"color"`
	Size int `json:"size"`
	User string `json:"user"`
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
	err = stub.PutState("abc", []byte(strconv.Itoa(Aval)))				//making a test var "abc", I find it handy to read/write to it right away to test the
	network
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
	} else if function == "delete" {										//deletes an entity from its state
		return t.Delete(stub, args)
	} else if function == "write" {											//writes a value to the chaincode state
		return t.Write(stub, args)
	} else if function == "init_travelIti" {									//create a new travelIti
		return t.init_travelIti(stub, args)
	} else if function == "set_user" {										//change owner of a travelIti
		return t.set_user(stub, args)
	}
	fmt.Println("run did not find func: " + function)						//error

	return nil, errors.New("Received unknown function invocation")
}

// ============================================================================================================================
// Delete - remove a key/value pair from state
// ============================================================================================================================
func (t *TravelItiChaincode) Delete(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	
	name := args[0]
	err := stub.DelState(name)													//remove the key from chaincode state
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	//srivatsav
	//get the travelIti index
	travelItiAsBytes, err := stub.GetState(travelItiIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get travelIti index")
	}
	var travelItiIndex []string
	json.Unmarshal(travelItiAsBytes, &travelItiIndex)								//un stringify it aka JSON.parse()
	
	//remove travelIti from index
	for i,val := range travelItiIndex{
		fmt.Println(strconv.Itoa(i) + " - looking at " + val + " for " + name)
		if val == name{															//find the correct travelIti
			fmt.Println("found Travel Itinary")
			travelItiIndex = append(travelItiIndex[:i], travelItiIndex[i+1:]...)			//remove it
			for x:= range travelItiIndex{											//debug prints...
				fmt.Println(string(x) + " - " + travelItiIndex[x])
			}
			break
		}
	}
	jsonAsBytes, _ := json.Marshal(travelItiIndex)									//save new index
	err = stub.PutState(travelItiIndexStr, jsonAsBytes)
	return nil, nil
}

// ============================================================================================================================
// Query - read a variable from chaincode state - (aka read)
// ============================================================================================================================
func (t *TravelItiChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function != "query" {
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetState(name)									//get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil													//send it onward
}

func main() {
	err := shim.Start(new(TravelItiChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// ============================================================================================================================
// Write - write variable into chaincode state
// ============================================================================================================================
func (t *TravelItiChaincode) Write(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var name, value string // Entities
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the variable and value to set")
	}

	name = args[0]															//rename for funsies
	value = args[1]
	err = stub.PutState(name, []byte(value))								//write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// ============================================================================================================================
// Init TravelIti - create a new Travel Itinerary, store into chaincode state
// ============================================================================================================================
func (t *TravelItiChaincode) init_travelIti(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var err error

	//   0       1       2     3
	// "asdf", "blue", "35", "bob"
	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	fmt.Println("- start init Travel Itinerary")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("4th argument must be a non-empty string")
	}
	
	size, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("3rd argument must be a numeric string")
	}
	
	color := strings.ToLower(args[1])
	user := strings.ToLower(args[3])

	str := `{"name": "` + args[0] + `", "color": "` + color + `", "size": ` + strconv.Itoa(size) + `, "user": "` + user + `"}`
	err = stub.PutState(args[0], []byte(str))								//store travel Itinerary with id as key
	if err != nil {
		return nil, err
	}
		
	//get the Travel Itinerary index
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

// ============================================================================================================================
// Set User Permission on travelIti
// ============================================================================================================================
func (t *TravelItiChaincode) set_user(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var err error
	
	//   0       1
	// "name", "bob"
	if len(args) < 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	
	fmt.Println("- start set user")
	fmt.Println(args[0] + " - " + args[1])
	travelItiAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get thing")
	}
	res := travelIti{}
	json.Unmarshal(travelItiAsBytes, &res)										//un stringify it aka JSON.parse()
	res.User = args[1]														//change the user
	
	jsonAsBytes, _ := json.Marshal(res)
	err = stub.PutState(args[0], jsonAsBytes)								//rewrite the travelIti with id as key
	if err != nil {
		return nil, err
	}
	
	fmt.Println("- end set user")
	return nil, nil
}
