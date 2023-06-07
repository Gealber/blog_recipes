#! /bin/bash

set -x
set -e
cd ../contracts
func -o ../bin/contract.fif -SPA stdlib.fc consts.fc utils.fc storage.fc user-storage.fc model.fc deposit.fc withdraw.fc borrow.fc requests.fc bank.fc get.fc
npx func-js stdlib.fc consts.fc utils.fc storage.fc user-storage.fc model.fc deposit.fc withdraw.fc borrow.fc requests.fc bank.fc get.fc --boc ../bin/contract.cell

# npx func-js stdlib.fc consts.fc bank.fc --boc ../bin/contract.cell
