#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#

# This is a collection of bash functions used by different scripts



# imports
. scripts/envVar.sh
. scripts/utils.sh

export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=${PWD}/configtx
export VERBOSE=false
export FABRIC_CFG_PATH=$PWD/../config/

CHANNEL_NAME=$1
CC_NAME=$2
FCN_NAME=$3
ARGS=$4
ORG=$5


setGlobals $ORG
res=$?

# while 'peer chaincode' command can get the orderer endpoint from the
# peer (if join was successful), let's supply it directly as we know
# it using the "-o" option
fcn_call="{\"function\":\"${FCN_NAME}\",\"Args\":${ARGS}}"
warnln "query fcn call:${fcn_call}"
set -x
peer chaincode query  -C $CHANNEL_NAME -n ${CC_NAME} -c "${fcn_call}" >&log.txt
# peer chaincode query -o localhost:7050 --ordererTLSHostnameOverride orderer1.example.com --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n ${CC_NAME} -c "${fcn_call}" >&log.txt
res=$?
{ set +x; } 2>/dev/null
cat log.txt
verifyResult $res "Invoke execution on $PEERS failed "
successln "Invoke transaction successful on $PEERS on channel '$CHANNEL_NAME'"
