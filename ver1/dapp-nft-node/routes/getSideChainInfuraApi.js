const {createResponse} = require('./response')
const fetch = require("node-fetch");
global.Headers = fetch.Headers;
const {NODE} = require('./config')

const infuraUrl = NODE.infuraUrl
const infura_read_Url = NODE.infura_read_Url

const getSideChainEthCall = (address, data) => {
    return fetch(
        infura_read_Url,{
        method: 'POST',
        body: JSON.stringify({
            id: 1,
            jsonrpc: '2.0',
            method: 'eth_call',
            params: [{
                data: data,
                to:address,
            }, "latest"]
        }),
        headers: new Headers({
            'Content-Type': 'application/json'
        })
        }
    )
    .then((response) => response.json())
    .then((data)=>data.result)
    .then((data) => {
        // console.log('getEthCall '+data)
        return data
    });
};

const getSideChainEthLogs = (params) => {
    return fetch(
        infuraUrl,{
        method: 'POST',
        body: JSON.stringify({
            id: 1,
            jsonrpc: '2.0',
            method: 'eth_getLogs',
            params: params
        }),
        headers: new Headers({
            'Content-Type': 'application/json'
        })
        }
    )
    .then((response) => response.json())
    .then((data)=>data.result)
    .then((data) => {
        console.log('getEthLogs '+JSON.stringify(data))
        return data
    });
};

const getSideChainEthTx = (tx) => {
    return fetch(
        infuraUrl,{
        method: 'POST',
        body: JSON.stringify({
            id: 1,
            jsonrpc: '2.0',
            method: 'eth_getTransactionReceipt',
            params: [tx]
        }),
        headers: new Headers({
            'Content-Type': 'application/json'
        })
        }
    )
    .then((response) => response.json())
    .then((data)=>data.result)
    .then((data) => {
        console.log('getEthLogs '+JSON.stringify(data))

        if (data && data.status=="0x1") {
            return data
        } else {
            return false
        }
    });
};

const redirectUrl = (url) => {
    return fetch(
        url,{
        method: 'GET',
        headers: new Headers({
            'Content-Type': 'application/json'
        })
        }
    )
    .then((response) => response.json())
    .then((data) => {

        return data
    });
};

module.exports = {getSideChainEthCall, getSideChainEthLogs, getSideChainEthTx, redirectUrl}