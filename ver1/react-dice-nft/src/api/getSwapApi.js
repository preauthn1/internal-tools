import { Eth } from "web3-eth";
import { paddingLeftZero, showBalanceDigit } from '../Utils'
import * as BigNumber from 'bignumber.js';
import * as Web3 from 'web3'

// const apiUrl = "https://example.invalid"
const apiUrl = "http://0.0.0.0:4000"

export const getTokenBalance = (address, tokenAddress) => {
    return fetch(
        `${apiUrl}/api/balance`,{
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

export const getTokenAllowance = (address) => {
    return fetch(
        `${apiUrl}/api/wallet/token/allowance`,{
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

export const postTokenApprove = (params) => {
    return fetch(
        `${apiUrl}/api/token/approve`,{
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

export const getTokenRate = (inTokenAddress, outTokenAddress) => {
  return fetch(
      `${apiUrl}/api/rate`,{
      method: 'POST',
      body: JSON.stringify({
        inTokenAddress:inTokenAddress,
        outTokenAddress:outTokenAddress
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

export const postTokenSwap = (params) => {
    return fetch(
        `${apiUrl}/api/swap`,{
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

export const postGetPairList = () => {
  return fetch(
      `${apiUrl}/api/pairlist`,{
      method: 'POST',
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

export const postGetSwapTxList = (address) => {
  return fetch(
      `${apiUrl}/api/swapTxList`,{
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