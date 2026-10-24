import { Eth } from "web3-eth";
import { paddingLeftZero, showBalanceDigit } from '../Utils'
import * as BigNumber from 'bignumber.js';
import * as Web3 from 'web3'

const apiUrl = "http://0.0.0.0:4000"


export const getTokenAllowance = (address, tokenAddress, spenderAddress) => {
  return fetch(
      `${apiUrl}/nft/api/tokenAllowance`,{
      method: 'POST',
      body: JSON.stringify({
          walletAddress:address,
          tokenAddress:tokenAddress,
          spenderAddress:spenderAddress
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

export const getTxid = (txid) => {
  return fetch(
      `${apiUrl}/nft/api/txid`,{
      method: 'POST',
      body: JSON.stringify({
          txid:txid
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

export const getNftInfo = (nftContract, nftId) => {
  return fetch(
      `${apiUrl}/nft/api/nftInfo`,{
      method: 'POST',
      body: JSON.stringify({
        nftContract:nftContract,
        nftId:nftId
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

export const postGetNftCollection = (address) => {
  return fetch(
      `${apiUrl}/nft/api/collection`,{
      method: 'POST',
      body: JSON.stringify({
        walletAddress:address
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

export const getTokenBalance = (address, tokenAddress) => {
    return fetch(
        `${apiUrl}/nft/api/balance`,{
        method: 'POST',
        body: JSON.stringify({
          walletAddress:address,
          tokenAddress:tokenAddress
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

export const postNftTransfer = (params) => {
  return fetch(
      `${apiUrl}/nft/api/transfer`,{
      method: 'POST',
      body: params,
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

export const postMarketplaceSell = (params) => {
  return fetch(
      `${apiUrl}/marketplace/api/sell`,{
      method: 'POST',
      body: params,
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

export const postGetMarketplaceOrders = (address) => {
  return fetch(
      `${apiUrl}/marketplace/api/orders`,{
      method: 'POST',
      body: JSON.stringify({
        walletAddress:address
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

export const postMarketplaceBuy = (params) => {
  return fetch(
      `${apiUrl}/marketplace/api/buy`,{
      method: 'POST',
      body: params,
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

export const postMarketplaceCancel = (params) => {
  return fetch(
      `${apiUrl}/marketplace/api/cancel`,{
      method: 'POST',
      body: params,
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