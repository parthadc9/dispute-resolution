package main

import (
        "fmt"
        "errors"
        "github.com/hyperledger/fabric-chaincode-go/shim"
        "github.com/hyperledger/fabric-protos-go/peer"
)

type Account struct{
	ID string `json:"ID"` 
	Collection string `json:"collection"`
	SortCode     string `json:"sortcode"`
        AccountNumber string `json:"accnumber"`
	IsActive     bool   `json:"isActive"`
	WarningID string `json:"warningID"`
}

const COLLECTION_ACCOUNT="ACCOUNT"

func(cc *TestChaincode)NewPayee(stub  shim.ChaincodeStubInterface, args []string)peer.Response{
	/**
  	1. Check the argument length and if they are blank 2. Map the arguments into proper variables 
	3. Check if any payee is existing with the given id
	4. If payee exists then response error id is taken
	5. Create the payee
	6. Payee confirmation

	*/

	if err:=ValidateArgLength(3,args);err!=nil{
		return shim.Error(err.Error())
	}

	account:= Account{
		Collection:COLLECTION_ACCOUNT,
		SortCode:args[0],
		AccountNumber:args[1],
		WarningID:args[2],
	}
	if account.SortCode==""{
		err:=errors.New("Account Sort Code cannot be blank")
                return shim.Error(err.Error())
	}
	 if account.AccountNumber==""{
                err:=errors.New(" Account Number cannot be blank")
                return shim.Error(err.Error())
        }

	account.ID=fmt.Sprintf("%v-%v-%v",account.Collection,account.SortCode,account.AccountNumber)
	ExistingBytes,err:=stub.GetState(account.ID)
	if err!=nil{
		return shim.Error(err.Error())
	}
	if ExistingBytes!=nil{
		err:=errors.New("Account Already Exists")
		return shim.Error(err.Error())
	}
	Byts,err:=PutDataToLedger(stub,account.ID,account)/*key value pair for levelDB*/
	if err!=nil{
		return shim.Error(err.Error())
	}
	return shim.Success(Byts)
}

/** This is to retrive a payee details based on sort code and account number*/
func(cc *TestChaincode) GetPayee(stub shim.ChaincodeStubInterface, args []string)peer.Response{
	/**
	1. Check the argument length and if they are blank 2. Map the arguments into proper variables
	3. Get state from the ledger
	4. Send the response
	*/

	if err:=ValidateArgLength(2,args);err!=nil{
                return shim.Error(err.Error())
        }
	SortCode:=args[0]
	AccountNumber:=args[1]
	 if SortCode==""{
                err:=errors.New("Account Sort Code cannot be blank")
                return shim.Error(err.Error())
        }
         if AccountNumber==""{
                err:=errors.New("Account Number cannot be blank")
                return shim.Error(err.Error())
        }
	ID:=fmt.Sprintf("%v-%v-%v",COLLECTION_ACCOUNT,SortCode,AccountNumber)
	Byts,err:=stub.GetState(ID)
	if err!=nil{
		return shim.Error(err.Error())
	}
	if Byts==nil{
                err:=errors.New("Account Does Not Exist")
                return shim.Error(err.Error())
        }
	return shim.Success(Byts)
}
/**

peer chaincode query -c drchannel -C payment -args '["GetPayee","12345", "56789"]'

*/

