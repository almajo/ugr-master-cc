#!/bin/bash
FILE=set_env.sh
az ad sp create-for-rbac > tmp.json

echo "export AZURE_TENANT_ID=$( cat tmp.json | jq '.tenant' | sed 's/"//g' )" > $FILE
echo "export AZURE_CLIENT_ID=$( cat tmp.json | jq '.appId' | sed 's/"//g' )" >> $FILE
echo "export AZURE_CLIENT_SECRET=$( cat tmp.json | jq '.password' | sed 's/"//g' )" >>$FILE

echo "export AZURE_SUBSCRIPTION_ID=$(az account list --query "[?isDefault].id" -o tsv)" >> $FILE

rm tmp.json
