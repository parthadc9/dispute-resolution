package main

import (
        "fmt"
        "errors"
        "github.com/hyperledger/fabric-chaincode-go/shim"
        "github.com/hyperledger/fabric-protos-go/peer"
	"encoding/json"
	"strings"
)

const (

	CHALLENGE_TYPE_WARNING = "warning"
	CHALLENGE_TYPE_PAYMENT = "payment"
	CHALLENGE_TYPE_MESSAGE = "message"
)

type Verdict struct{
	ID string `json: "ID"`
	Verdict string `json: "verdict"`
}

type Complaint struct{
	Collection string `json:"collection"`
	ID string `json:"ID"` 
	PaymentRequestID string `json:"payementRequestID"`
	BeneficiaryID string `json:"beneficiaryID"`
        WarningID string `json:"warningID"`
	PayerID string `json:"payerID"`
	PaymentConfirmationID string `json:"paymentID"`
	ChallengeType string `json:"challengetype"`
	Evidence string `json:"evidence"`
	Verdict map[string]Verdict `json:"verdict"`
}

const (
	CHALLENGE_MESSAGE="e"
	CHALLENGE_WARNING="x"
	CHALLENGE_PAYMENT="y"

)
const COLLECTION_COMPLAINT="CMP"

func(complaint *Complaint)IsResolved()bool{

	if complaint.Verdict==nil{
		return false
	}
	return len(complaint.Verdict)>=3

}


func(cc *TestChaincode)GenerateComplaint(stub  shim.ChaincodeStubInterface, args []string)peer.Response{

/**

	1. Check the argument length. 2.Map the arguments into proper variables.
	3. Validate the arguments. 4. Check the paymentID, warningID, beneficiaryID, payerID if they are valid
	5. Only payer can raise the dispute.
	6. Challenge type should be valid.
	7. If challenge type is warning then evidence type should not be empty.
	8. Save the complaint.  

*/
	if err:=ValidateArgLength(2,args);err!=nil{
                return shim.Error(err.Error())
        }
	
	PaymentRequestID:=args[0]
	ChallengeType:=args[1]	
	
	if PaymentRequestID==""{
                err:=errors.New("Payment Request ID cannot be blank")
                return shim.Error(err.Error())
        }
	switch ChallengeType{
	case CHALLENGE_MESSAGE:
	case CHALLENGE_WARNING:
	case CHALLENGE_PAYMENT:
	default:
		err:=errors.New("Invalid Challenge Type")
		return shim.Error(err.Error())
	}
	// PaymentRequestID=fmt.Sprintf("%v-%v",COLLECTION_PAYMENT_REQUEST,PaymentRequestID)
	
	splitted := strings.Split(PaymentRequestID, "-")
	channelId := "drchannel"
	chaincodeName := "payment"
	fcnArgs := [][]byte{
		[]byte("GetPaymentRequest"),
		[]byte(splitted[1]),
	}
	resp := stub.InvokeChaincode(chaincodeName, fcnArgs, channelId)
	requestByts := resp.Payload
	if requestByts==nil{
		err:=errors.New("Payment Request Doesn't Exist")
		return shim.Error(err.Error())
	}
	request:=PaymentRequest{}
	if err:=json.Unmarshal(requestByts,&request);err!=nil{
		return shim.Error(err.Error())
	}
	if request.PaymentTransactionID==""{
		err:=errors.New("Payment Not Processed")
		return shim.Error(err.Error())
	}


	complaintID:=fmt.Sprintf("%v-%v", COLLECTION_COMPLAINT,request.PaymentTransactionID)
	ExistingBytes,err:=stub.GetState(complaintID)
	if err!=nil{
		return shim.Error(err.Error())
	}
	if ExistingBytes!=nil{
		err:=errors.New("Complaint Already Exists")
		return shim.Error(err.Error())
	}
	/* Creating the Compliant */
	
	complaint:=Complaint{
		Collection:COLLECTION_COMPLAINT,
		PaymentRequestID:request.ID,
		BeneficiaryID:request.BeneficiaryID,
		PaymentConfirmationID:request.PaymentTransactionID,
        	WarningID:request.WarningID,
		PayerID:request.PayerID,
		ChallengeType:ChallengeType,
		Evidence:"", 
	}
	complaint.ID=fmt.Sprintf("%v-%v",complaint.Collection,complaint.PaymentRequestID)
	Byts,err:=PutDataToLedger(stub,complaint.ID,complaint)/*key value pair for levelDB*/
	if err!=nil{
		return shim.Error(err.Error())
	}
	return shim.Success(Byts)
}
