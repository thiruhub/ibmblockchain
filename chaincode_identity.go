/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	
)

// IdentityChainCode example simple Chaincode implementation
type IdentityChainCode struct {
}


// Init callback representing the invocation of a chaincode
// This chaincode will manage two accounts A and B and will transfer X units from A to B upon invoke
func (t *IdentityChainCode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Assigning Asset...")
	
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}
	
	//In PolicyName,In ResourceName,In ResourceIP, In ResourceURL,In UserName,In UserType,In UserRole,In EnrollID, In IsAuthorized
	
	// Create IAManager Table
	err := stub.CreateTable("IAManagerTable", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "UserName", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "IsAuthorized", Type: shim.ColumnDefinition_BOOL, Key: true},
		&shim.ColumnDefinition{Name: "PolicyName", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ResourceName", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ResourceIP", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ResourceURL", Type: shim.ColumnDefinition_STRING, Key: false},		
		&shim.ColumnDefinition{Name: "UserType", Type: shim.ColumnDefinition_BYTES, Key: false},
		&shim.ColumnDefinition{Name: "UserRole", Type: shim.ColumnDefinition_BYTES, Key: false},
		&shim.ColumnDefinition{Name: "EnrollID", Type: shim.ColumnDefinition_STRING, Key: false},		
	})
	if err != nil {
		return nil, fmt.Errorf("Failed creating IAManager table, [%v]", err)
	}
	
	fmt.Println("Created IAManagerTable...")
	
	// Create Log Table
	errLog := stub.CreateTable("AccessLogTable", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "ResourceName", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ResourceIP", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "ResourceURL", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "UserName", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "EnrollID", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "IsAuthorized", Type: shim.ColumnDefinition_BOOL, Key: false},
		&shim.ColumnDefinition{Name: "DateTime", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if errLog != nil {
		return nil, fmt.Errorf("Failed creating AccessLog table, [%v]", errLog)
	}
	
	fmt.Println("Created AccessLog Table...")

	// Initialize the chaincode
	fmt.Println("Initialize chaincode")
	return nil, nil
}

func (t *IdentityChainCode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("In invoke with arg " + function)

	// Handle different functions
	if function == "Init" {					//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "resourcecreate" {		//Create Resource
		return t.resourcecreate(stub, args)
	}else if function == "policycreate" {			//Create policy
		return t.policycreate(stub, args)
	} else if function == "policydelete" {			//Delete policy
		return t.policydelete(stub, args)
	} else if function == "policymodify" {			//Modify policy
		return t.policymodify(stub, args)
	} else if function == "fetchlogs" {				//Fetch Access Logs
		return t.fetchlogs(stub, args)
	}
	
	
	fmt.Println("invoke did not find func: " + function)	//error

	return nil, errors.New("Received unknown function invocation")
}

// ============================================================================================================================
// Create Policy - create a key/value pair from state
// ============================================================================================================================
func (t *IdentityChainCode) policycreate(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("In policycreate...")
	
	var err error
		
	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	
	/*policyname  := args[0]
	resourceurl := args[1] 
	username    := args[2]
	permission  := args[3]*/
	
	str := `{"PolicyName": "` + args[0] + `", "ResourceURL": "` + args[1] + `", "UserName": ` + args[2] + `, "IsAuthorized": "` + args[3] + `"}`
	
	err = stub.PutState(args[0], []byte(str))								
	if err != nil {
		return nil, err
	}

	fmt.Println("created the policy successfully")
	return nil, nil
}

// ============================================================================================================================
// policydelete - remove a key/value pair from state
// ============================================================================================================================
func (t *IdentityChainCode) policydelete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("In policydelete...")
	
    var policyname string
	
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	
	policyname = args[0]
	err := stub.DelState(policyname)													//remove the key from chaincode state
	if err != nil {
		return nil, errors.New("Failed to delete policy")
	}
	
	fmt.Println("Removed the policy successfully")
	return nil, nil
}

// ============================================================================================================================
// Policyremove - remove a key/value pair from state
// ============================================================================================================================
func (t *IdentityChainCode) policymodify(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("in policymodify...")
	
	var policyoldname string
	
	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	
	policyoldname = args[1] 
		
	err := stub.DelState(policyoldname)													//remove the key from chaincode state
	if err != nil {
		return nil, errors.New("Failed to delete policy")
	}
		
	str := `{"PolicyName": "` + args[0] + `", "ResourceURL": "` + args[2] + `", "UserName": ` + args[3] + `, "IsAuthorized": "` + args[4] + `"}`
	
	err = stub.PutState(args[0], []byte(str))								//store policyname with id as key
	if err != nil {
		return nil, err
	}
	
	fmt.Println("Modified the policy successfully")
	return nil, nil
}

// ============================================================================================================================
// resourcecreate - create a key/value pair from state
// ============================================================================================================================
func (t *IdentityChainCode) resourcecreate(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("In resourcecreate...")
		
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}
	
	//(In ResourceName,In ResourceIP, In ResourceURL)
	str := `{"ResourceName": "` + args[0] + `", "ResourceIP": "` + args[1] + `", "ResourceURL": ` + args[2] + `"}`
	
	err := stub.PutState(args[0], []byte(str))								//store policyname with id as key
	if err != nil {
		return nil, err
	}
	
	fmt.Println("Modified the policy successfully")
	return nil, nil
}

// ============================================================================================================================
// fetchlogs - create a key/value pair from state
// ============================================================================================================================
func (t *IdentityChainCode) fetchlogs(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("In fetchlogs...")
	
	var indexfrom string
	var indexto string
	
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	
	//(In ResourceName,In ResourceIP, In ResourceURL)
	str := `{"ResourceName": "` + args[0] + `", "ResourceIP": "` + args[1] + `", "ResourceURL": ` + args[2] + `"}`
	
	err := stub.PutState(args[0], []byte(str))								//store policyname with id as key
	if err != nil {
		return nil, err
	}
	
	fmt.Println("Modified the policy successfully")
	return nil, nil
}


// Query callback representing the query of a chaincode
func (t *IdentityChainCode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return nil, nil
}

func main() {
	err := shim.Start(new(IdentityChainCode))
	if err != nil {
		fmt.Printf("Error starting Identity chaincode: %s", err)
	}	
}
