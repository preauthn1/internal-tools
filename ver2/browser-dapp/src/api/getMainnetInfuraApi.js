const apiUrl = "http://0.0.0.0:4001"

export const getPoaWalletNonce = (walletAddress, authorizedAddress) => {
    return fetch(
        `${apiUrl}/poa/api/wallet/nonce`,{
            method: 'POST',
            body: JSON.stringify({
                walletAddress:walletAddress,
                authorizedAddress:authorizedAddress
            }),
            headers: new Headers({
            'Content-Type': 'application/json'
            })
        }
    )
        .then((response) => response.json())
        .then((data)=>data.data)
        .then((data) => {
            return data
        });
};

export const getSideChainEthNonce = (address) => {
    return fetch(
        `${apiUrl}/poa/api/eth/nonce`,{
            method: 'POST',
            body: JSON.stringify({
                ethAddress:address,
            }),
            headers: new Headers({
            'Content-Type': 'application/json'
            })
        }
    )
        .then((response) => response.json())
        .then((data)=>data.data)
        .then((data) => {
            return data
        });
};