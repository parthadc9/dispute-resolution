package main

import (
	"strings"
	"crypto/x509"
	"encoding/pem"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	
)

// ValidateArgLength checks the argument length against expected length of argument
func ValidateArgLength(expectedArgLength int, args []string) error {
	length := len(args)
	if length != expectedArgLength {
		err := fmt.Sprintf("expected %v arguments, received %v", expectedArgLength, length)
		return errors.New(err)
	}
	return nil
}

// convert Query Result into a JSON
func generateJSONFromQueryResult(resultsIterator shim.StateQueryIteratorInterface) (bytes.Buffer, error) {
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			//logger.Error(err.Error())
			return buffer, err
		}

		// Add a comma before array members, suppress it for the first array member

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	return buffer, nil
}

// PutDataToLedger updates the ledger
func PutDataToLedger(stub shim.ChaincodeStubInterface, key string, data interface{}) ([]byte, error) {
	dataAsBytes, err := json.Marshal(data)
	// raise error if any
	if err != nil {
		return nil, err
	}

	// update ledger
	if err := stub.PutState(key, dataAsBytes); err != nil {
		return nil, err
	}

	return dataAsBytes, nil
}

// QueryLedger executes a query and returns result
func QueryLedger(stub shim.ChaincodeStubInterface, queryString string) (shim.StateQueryIteratorInterface, error) {
	// raise error if query string is blank
	if queryString == "" {
		return nil, errors.New("queryString is required")
	}
	//logger.Info(queryString)
	return stub.GetQueryResult(queryString)
}

/* get userNAme and orgName from certificate */
func GetCreator(certificate []byte) (string, string) {
	data := certificate[strings.Index(string(certificate), "-----") : strings.LastIndex(string(certificate), "-----")+5]
	block, _ := pem.Decode([]byte(data))
	cert, _ := x509.ParseCertificate(block.Bytes)
	organization := cert.Issuer.Organization[0]
	commonName := cert.Subject.CommonName
	// logger.Debug("commonName: " + commonName + ", organization: " + organization)
	organizationShort := strings.Split(organization, ".")[0]
	return commonName, organizationShort
}
