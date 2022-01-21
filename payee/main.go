package main

import (
	"fmt"
	

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)


type TestChaincode struct {
}

//var logger *shim.ChaincodeLogger

// main function starts up the chaincode in the container during instantiate
func main() {
	
	//logger = shim.NewLogger("payment")
	if err := shim.Start(new(TestChaincode)); err != nil {
		fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}
}

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
func (t *TestChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// Get the args from the transaction proposal
	// args := stub.GetStringArgs()
	// if len(args) != 3 {
	// 	return shim.Error("Incorrect arguments. Expecting a name, version and sequence number")
	// }
	// data := map[string]interface{}{}

	// data["name"] = args[0]
	// data["version"] = args[1]
	// data["sequence"] = args[2]

	// dataBytes, err := json.Marshal(data)
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }
	// // Set up any variables or assets here by calling stub.PutState()

	// // We store the key and the value on the ledger
	// err = stub.PutState("meta", dataBytes)
	// if err != nil {
	// 	return shim.Error(fmt.Sprintf("Failed to update meta: %s", data))
	// }
	return shim.Success(nil)
}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The Set
// method may create a new asset by specifying a new key-value pair.
func (t *TestChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	fnc, args := stub.GetFunctionAndParameters()

	switch fnc {
	case "NewPayee":
		return t.NewPayee(stub, args)
	case "GetPayee":
		return t.GetPayee(stub, args)
	}

	// Return the result as success payload
	return shim.Success(nil)
}

/**
peer chaincode invoke -c paymentchannel -C payment [ "GetBalance", "name" ]
*/

