package main

import (
        "fmt"
        "errors"
        "github.com/hyperledger/fabric-chaincode-go/shim"
        "github.com/hyperledger/fabric-protos-go/peer"
	"encoding/json"
	"strings"
)

type PaymentRequest struct{
	Collection string `json:"collection"`
	ID string `json:"ID"` 
	PayerID string `json:"payerID"`
	BeneficiaryID     string `json:"beneficiaryID"`
        Amount string `json:"amount"`
	Status string `json:"status"`
	WarningID string `json:"warningID"`	
	PaymentTransactionID string `json:"paymentTransactionID"`
	PaymentMode string `json:"paymentmode"`
}

const COLLECTION_PAYMENT_REQUEST="PR"

func(cc *TestChaincode)GeneratePaymentRequest(stub  shim.ChaincodeStubInterface, args []string)peer.Response{

 /**
	1. Check the argument length. 2. Map the arguments into proper variables
	3. Get the details of payer and beneficiary
	4. Create object of payment request and fill the details
	5. If the beneficiary account is associated with any warning create an object of warning
	6. Assign the warning ID to the payment request
	7. Save the warning if any
	8. Save the payment request 
  */
	
	if err:=ValidateArgLength(4,args);err!=nil{
                return shim.Error(err.Error())
        }
	request:=PaymentRequest{
		Collection:COLLECTION_PAYMENT_REQUEST,
		PayerID:args[0],
		BeneficiaryID:args[1],
		Amount:args[2],
		PaymentMode:args[3],
	}
	
	if request.PayerID==""{
                err:=errors.New("Payer ID cannot be blank")
                return shim.Error(err.Error())
        }
	if request.BeneficiaryID==""{
                err:=errors.New("Beneficiary ID cannot be blank")
                return shim.Error(err.Error())
        }
	if request.Amount==""{
                err:=errors.New("Amount cannot be blank")
                return shim.Error(err.Error())
        }
	if request.PaymentMode==""{
                err:=errors.New("Payment Mode Code cannot be blank")
                return shim.Error(err.Error())
        }

	// Get the payer
	splitted := strings.Split(request.PayerID, "-")
	channelId := "drchannel"
	chaincodeName := "payee"
	fcnArgs := [][]byte{
		[]byte("GetPayee"),
		[]byte(splitted[1]),
		[]byte(splitted[2]),
	}
	resp := stub.InvokeChaincode(chaincodeName, fcnArgs, channelId)
	payerBytes := resp.Payload

	if payerBytes==nil{
		err:=errors.New("Payer Doesn't Exist")
		return shim.Error(err.Error())
	}
	splitted = strings.Split(request.BeneficiaryID, "-")
	channelId = "drchannel"
	chaincodeName = "payee"
	fcnArgs = [][]byte{
		[]byte("GetPayee"),
		[]byte(splitted[1]),
		[]byte(splitted[2]),
	}
	resp = stub.InvokeChaincode(chaincodeName, fcnArgs, channelId)
	beneficiaryBytes := resp.Payload
	if beneficiaryBytes==nil{
		err:=errors.New("Beneficiary Doesn't Exist")
		return shim.Error(err.Error())
	}
	beneficiary:=Account{}
	if err:=json.Unmarshal(beneficiaryBytes,&beneficiary);err!=nil{
		return shim.Error(err.Error())
	}
	request.WarningID=beneficiary.WarningID
	request.ID=fmt.Sprintf("%v-%v",request.Collection,stub.GetTxID())
	
	Byts,err:=PutDataToLedger(stub,request.ID,request)/*key value pair for levelDB*/
	if err!=nil{
		return shim.Error(err.Error())
	}
	return shim.Success(Byts)
	

}

func (cc *TestChaincode) GetPaymentRequest(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if err := ValidateArgLength(1, args); err != nil {
		return shim.Error(err.Error())
	}
	ID := fmt.Sprintf("%v-%v", COLLECTION_PAYMENT_REQUEST, args[0])
	if ID == "" {
		err := errors.New("Invalid payment ID")
		return shim.Error(err.Error())
	}

	Bytes, err := stub.GetState(ID)
	if err != nil {
		return shim.Error(err.Error())
	}
	if Bytes == nil {
		err := errors.New("PaymentID Doesn't Exist")
		return shim.Error(err.Error())
	}
	return shim.Success(Bytes)

}

func(cc *TestChaincode)MakePayment(stub  shim.ChaincodeStubInterface, args []string)peer.Response{

/**
	1. Check the argument length. 2. Map the arguments into proper variables. 
	2a. Transaction ID is a mandatory parameter. 
	3. Get the payment request details
	4. Get the payer and beneficiary account details
	5. Transfer the balance 
	6. Save payer 7. Save beneficiary 8. Update payment request to paid 
*/
	if err:=ValidateArgLength(2,args);err!=nil{
                return shim.Error(err.Error())
        }
		
	requestID:=args[0]
	paymentConfirmationID:=args[1]
	if requestID==""{
                err:=errors.New("Payment Request ID cannot be blank")
                return shim.Error(err.Error())
        }
	if paymentConfirmationID==""{
                err:=errors.New("Payment Confirmation ID cannot be blank")
                return shim.Error(err.Error())
        }

	requestBytes,err:=stub.GetState(requestID)
	if err!=nil{
		return shim.Error(err.Error())
	}
	if requestBytes==nil{
		err:=errors.New("Payment Request Doesn't Exist")
		return shim.Error(err.Error())
	}
	request:=PaymentRequest{}
	if err:=json.Unmarshal(requestBytes,&request);err!=nil{
		return shim.Error(err.Error())
	}
	if request.PaymentTransactionID!=""{
		err:=errors.New("Payment Confirmation ID is already updated")
		return shim.Error(err.Error())
	}
	request.PaymentTransactionID=paymentConfirmationID
	Byts,err:=PutDataToLedger(stub,request.ID,request)/*key value pair for levelDB*/
	if err!=nil{
		return shim.Error(err.Error())
	}
	return shim.Success(Byts)
}
