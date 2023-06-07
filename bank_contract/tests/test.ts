import { Address, Cell, beginCell, Dictionary } from "ton-core";
import { SmartContract, internal } from "ton-contract-executor";
import { TonClient } from "ton";

const contractAddress = Address.parse('EQCJyC8kTRnKZIx4K1f4fepy417p3isZhCVaZvxjJ8fcVEDr')

let client = new TonClient({
    endpoint: 'https://toncenter.com/api/v2/jsonRPC'
})

async function main() {
    let state = await client.getContractState(contractAddress)

    let code = Cell.fromBoc(state.code!)[0]
    let data = Cell.fromBoc(state.data!)[0]
    
    let contract = await SmartContract.fromCell(
        code,
        data,
        // {debug: true},
    )

    const depositOPPayload = beginCell()
        .storeUint(100, 32)
        .endCell()        

    const depRes = await contract.sendInternalMessage(internal({
            dest: contractAddress,
            value: BigInt(1000000000), // 1 ton
            bounce: false,
            body: depositOPPayload,
        }))

    console.log("DEPOSIT RES: ", depRes)    

    const withdrawOPPayload = beginCell()
        .storeUint(119, 32)
        .storeCoins(600000000)
        .endCell()

    const res = await contract.sendInternalMessage(internal({
        dest: contractAddress,
        value: BigInt(0), // 1 ton
        bounce: false,
        body: withdrawOPPayload,
    }))

    console.log("WITHDRAW: ", res)

    const borrowOPPayload = beginCell()
        .storeUint(98, 32)        
        .endCell()

    var borowRes = await contract.sendInternalMessage(internal({
        dest: contractAddress,
        value: BigInt(0),
        bounce: false,
        body: borrowOPPayload,
    }))

    console.log("BORROW: ", borowRes)    

    const getRes = await contract.invokeGetMethod("get_users", [])

    console.log("RES: ", getRes.result)
}

main()