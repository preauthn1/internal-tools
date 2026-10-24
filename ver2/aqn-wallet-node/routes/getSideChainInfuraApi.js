const {createResponse} = require('./response')
const {NODE} = require('./config')

const fetch = require("node-fetch");
const FormData = require('form-data');

global.Headers = fetch.Headers;

const infuraUrl = NODE.infuraUrl
const infuraUrl_read = NODE.infuraUrl_read

const sendEmail = (to, content) => {

    var formData = {
        subject:"Wallet Verification",
        fromName:"Dice",
        fromAddress:"user@example.invalid",
        content:content,
        recipients:[
            {
                name:to,
                address:to
            }
        ]
    }
    
    return fetch(
            "https://example.invalid/v1/messages",{
            method: 'POST',
            headers: new Headers({
                'Content-Type': 'application/json',
                'x-api-key':NODE.surenotifyApiKey
            }),
            body: JSON.stringify(formData),
        }
    )
    .then((response) => response.json())
    .then((data)=>{
        console.log('sendEmail '+JSON.stringify(data))

        return data.id
    })
    .then((data) => {
        console.log('sendEmail '+data)
        return data
    });
}

const getSideChainEthCall = (address, data) => {
    return fetch(
        infuraUrl_read,{
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
        console.log('getEthCall '+data)
        return data
    });
};

const getSideChainEthLogs = (params) => {
    return fetch(
        infuraUrl_read,{
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

const getSideChainEthNonce = (address) => {
    return fetch(
        infuraUrl,{
        method: 'POST',
        body: JSON.stringify({
            id: 1,
            jsonrpc: '2.0',
            method: 'eth_getTransactionCount',
            params: [address, "latest"]
        }),
        headers: new Headers({
            'Content-Type': 'application/json'
        })
        }
    )
        .then((response) => response.json())
        .then((data)=>data.result)
        .then((data) => {
            console.log('getEthNonce '+data)
            return data
        });
};

const sendSideChainEthRawTxWithResponse = (rawTx, res) => {
    return fetch(
        infuraUrl,{
        method: 'POST',
        body: JSON.stringify({
            id: 1,
            jsonrpc: '2.0',
            method: 'eth_sendRawTransaction',
            params: [rawTx]
        }),
        headers: new Headers({
            'Content-Type': 'application/json'
        })
        }
    )
    .then((response) => response.json())
    // .then((data)=>data.result)
    .then((data) => {
        console.log('sendEthRawTx '+JSON.stringify(data))
        res.send(createResponse(true, data.result))

        return data.result
    });
};

const sendSideChainEthRawTx = (rawTx) => {
    return fetch(
        infuraUrl,{
        method: 'POST',
        body: JSON.stringify({
            id: 1,
            jsonrpc: '2.0',
            method: 'eth_sendRawTransaction',
            params: [rawTx]
        }),
        headers: new Headers({
            'Content-Type': 'application/json'
        })
        }
    )
    .then((response) => response.json())
    .then((data)=>data.result)
    .then((data) => {
        console.log('sendEthRawTx '+data)
        return data
    });
};

const getSideChainEthTx = (tx) => {
    return fetch(
        infuraUrl_read,{
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

        // if (data && data.status=="0x1") {
        //     return true
        // } else {
        //     return false
        // }
        return data
    });
};

const redirectUrl = (url, params) => {
    return fetch(
        url,{
        method: 'POST',
        body: JSON.stringify(params),
        headers: new Headers({
            'Content-Type': 'application/json'
        })
        }
    )
    .then((response) => response.json())
    .then((data) => {
        // console.log('redirectUrl '+JSON.stringify(data))

        return data
    });
};

module.exports = {getSideChainEthCall, getSideChainEthNonce, sendSideChainEthRawTxWithResponse, sendSideChainEthRawTx, getSideChainEthLogs, getSideChainEthTx, redirectUrl, sendEmail}