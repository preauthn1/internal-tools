import { Eth } from "web3-eth";
import { paddingLeftZero, showBalanceDigit } from '../Utils'
import * as BigNumber from 'bignumber.js';
import * as Web3 from 'web3'

// const apiUrl = "https://example.invalid/"
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

export const getInfo = () => {
  return fetch(
      `${apiUrl}/api/info`,{
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

export const getPoolList = () => {
  return fetch(
      `${apiUrl}/api/poollist`,{
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

export const postTokenLending = (params) => {
    return fetch(
        `${apiUrl}/api/lending`,{
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

export const getLendingList = (walletAddress, tokenTab) => {
  return fetch(
      `${apiUrl}/api/lending/list`,{
      method: 'POST',
      body: JSON.stringify({
        walletAddress:walletAddress,
        tokenTab:tokenTab
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

export const postTokenRedeem = (params) => {
  return fetch(
      `${apiUrl}/api/redeem`,{
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