package main

import (
        
        "errors"
        "github.com/hyperledger/fabric-chaincode-go/shim"
        "github.com/hyperledger/fabric-protos-go/peer"
	"encoding/json"
	"fmt"
)

const (
	VERDICT_MESSAGE="e"
	VERDICT_WARNING="x"
	VERDICT_PAYMENT="y"
	VERDICT_DISSENT="d"
)

func(cc *TestChaincode)ViewComplaint(stub  shim.ChaincodeStubInterface, args []string)peer.Response{

	if err:=ValidateArgLength(1,args);err!=nil{
                return shim.Error(err.Error())
        }
	complaintID:=args[0]
	if complaintID==""{
		err:=errors.New("Invalid Complaint ID")
                return shim.Error(err.Error())
        }

	Bytes,err:=stub.GetState(complaintID)
	if err!=nil{
		return shim.Error(err.Error())
	}
	if Bytes==nil{
		err:=errors.New("Complaint Doesn't Exist")
		return shim.Error(err.Error())
	}
	return shim.Success(Bytes)

}

func(cc *TestChaincode)ResolveComplaint(stub  shim.ChaincodeStubInterface, args []string)peer.Response{

	if err:=ValidateArgLength(2,args);err!=nil{
                return shim.Error(err.Error())
        }
	complaintID:=args[0]
	VerdictType:=args[1]
	if complaintID==""{
		err:=errors.New("Invalid Complaint ID")
                return shim.Error(err.Error())
        }
	switch VerdictType{
	case VERDICT_MESSAGE:
	case VERDICT_WARNING:
	case VERDICT_PAYMENT:
	case VERDICT_DISSENT:
	default:
		err:=errors.New("Invalid Verdict")
		return shim.Error(err.Error())
	}
	complaintBytes,err:=stub.GetState(complaintID)
	if err!=nil{
		return shim.Error(err.Error())
	}
	if complaintBytes==nil{
		err:=errors.New("Complaint Doesn't Exist")
		return shim.Error(err.Error())
	}
	complaint:=Complaint{}
	if err:=json.Unmarshal(complaintBytes,&complaint);err!=nil{
		return shim.Error(err.Error())
	}

	if complaint.Verdict != nil && len(complaint.Verdict) >=5{
		err:=errors.New("Voting is Closed")
		return shim.Error(err.Error())
	}

	verdict:=Verdict{}

	verdict.Verdict=VerdictType
	creatorBytes, err := stub.GetCreator()
	if err != nil{
		err:=errors.New("Cannot Find Arbitrer Name From Certificate")
		return shim.Error(err.Error())
	}	
	verdict.ID,_=GetCreator(creatorBytes)
	
	if complaint.Verdict==nil{
		complaint.Verdict=map[string]Verdict{}
	}
	complaint.Verdict[verdict.ID]=verdict

	Bytes,err:=PutDataToLedger(stub,complaint.ID,complaint)/*key value pair for levelDB*/
	if err!=nil{
		return shim.Error(err.Error())
	}

	return shim.Success(Bytes)
	
}

func(cc *TestChaincode)GetComplaintStatus(stub  shim.ChaincodeStubInterface, args []string)peer.Response{

/**

	1. Check the argument length. 2.Map the arguments into proper variables.
	3. Validate the arguments. 4. Count verdict 5. Display majority verdict

*/

	if err:=ValidateArgLength(1,args);err!=nil{
                return shim.Error(err.Error())
        }
	complaintID:=args[0]
	if complaintID==""{
		err:=errors.New("Invalid Complaint ID")
                return shim.Error(err.Error())
        }
	complaintBytes,err:=stub.GetState(complaintID)
	if err!=nil{
		return shim.Error(err.Error())
	}
	if complaintBytes==nil{
		err:=errors.New("Complaint Doesn't Exist")
		return shim.Error(err.Error())
	}
	complaint:=Complaint{}
	if err:=json.Unmarshal(complaintBytes,&complaint);err!=nil{
		return shim.Error(err.Error())
	}
	if complaint.Verdict==nil{
		message:="No Verdict Received"
		return shim.Success([]byte(message))
	}
	count:=map[string]int{}
	for _,verdict:=range complaint.Verdict{
		value,ok:=count[verdict.Verdict]
		if ok{
			count[verdict.Verdict]=value+1
		}else{
			count[verdict.Verdict]=1
		}

	}
	message:=""
	if complaint.IsResolved(){
		message="The Complaint Is Resolved As \n"
	
		for verdict,value:=range count{
		message=fmt.Sprintf("%v%v: %v\t",message,verdict,value)
		}
	}else{
		message="The Complaint is Not Yet Resolved"
	}
	return shim.Success([]byte(message))
}
