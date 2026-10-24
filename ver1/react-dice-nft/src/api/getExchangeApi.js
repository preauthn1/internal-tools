import { Eth } from "web3-eth";
import { paddingLeftZero, showBalanceDigit } from '../Utils'
import * as BigNumber from 'bignumber.js';
import * as Web3 from 'web3'

const apiUrl = "https://example.invalid"

export const getPoolList = (symbol) => {
  return fetch(
      `${apiUrl}/index/poolList`,{
      method: 'GET',
      headers: new Headers({
        'Content-Type': 'application/json'
      })
    }
  )
  .then((response) => response.json())
  .then((data) => {

  return data.data
  });
};

export const getSwapList = (symbol) => {
  return fetch(
      `${apiUrl}/index/swapList`,{
      method: 'GET',
      headers: new Headers({
        'Content-Type': 'application/json'
      })
    }
  )
  .then((response) => response.json())
  .then((data) => {

  return data.data
  });
};

export const getExchangeList = (symbol) => {
  return fetch(
      `${apiUrl}/index/exchangeList`,{
      method: 'GET',
      headers: new Headers({
        'Content-Type': 'application/json'
      })
    }
  )
  .then((response) => response.json())
  .then((data) => {

  return data.data
  });
};

export const getOrderBook = (symbol) => {
  return fetch(
      `${apiUrl}/${symbol}/orderBook`,{
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

export const postOrderBuy = (symbol, params) => {
  return fetch(
      `${apiUrl}/${symbol}/api/order/buy`,{
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

export const postOrderSell = (symbol, params) => {
  return fetch(
      `${apiUrl}/${symbol}/api/order/sell`,{
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

export const getUserHistory = (symbol, address) => {
  return fetch(
      `${apiUrl}/${symbol}/api/history`,{
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

export const getUserOrder = (symbol, address) => {
  return fetch(
      `${apiUrl}/${symbol}/api/orders`,{
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

export const postOrderCancel = (symbol, params) => {
  return fetch(
      `${apiUrl}/${symbol}/api/order/cancel`,{
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

export const getUserStakeHistory = (symbol, address) => {
  return fetch(
      `${apiUrl}/${symbol}/api/stake/histroy`,{
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