import { Eth } from "web3-eth";
import { paddingLeftZero, showBalanceDigit } from '../Utils'
import * as BigNumber from 'bignumber.js';
import * as Web3 from 'web3'

const apiUrl = "https://example.invalid"
// const apiUrl = "http://0.0.0.0:4000"

export const postGetPairList = () => {
  return fetch(
      `${apiUrl}/swap/api/pairlist`,{
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

export const getTokenBalance = (address, tokenAddress) => {
    return fetch(
        `${apiUrl}/swap/api/balance`,{
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

export const getSwapK = (inTokenAddress, outTokenAddress) => {
  return fetch(
      `${apiUrl}/swap/api/k`,{
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
        `${apiUrl}/swap/api/swap`,{
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

//pool
export const getInfo = () => {
  return fetch(
      `${apiUrl}/pool/api/info`,{
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
      `${apiUrl}/pool/api/poollist`,{
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

export const postPoolSubscribe = (params) => {
  return fetch(
      `${apiUrl}/pool/api/subscribe`,{
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

export const postPoolRedeem = (params) => {
  return fetch(
      `${apiUrl}/pool/api/redeem`,{
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

export const postPoolReward = (params) => {
  return fetch(
      `${apiUrl}/pool/api/reward`,{
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

export const getPoolDeposited = (walletAddress, nodeHash) => {
  return fetch(
      `${apiUrl}/pool/api/deposited`,{
      method: 'POST',
      body: JSON.stringify({
        walletAddress:walletAddress,
        nodeHash:nodeHash
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

export const getFarmList = () => {
  return fetch(
      `${apiUrl}/farm/api/farmlist`,{
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

export const postFarmSubscribe = (params) => {
  return fetch(
      `${apiUrl}/farm/api/subscribe`,{
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

export const getFarmDeposited = (walletAddress) => {
  return fetch(
      `${apiUrl}/farm/api/deposited`,{
      method: 'POST',
      body: JSON.stringify({
        walletAddress:walletAddress
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

export const postFarmRedeem = (params) => {
  return fetch(
      `${apiUrl}/farm/api/redeem`,{
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

export const postFarmReward = (params) => {
  return fetch(
      `${apiUrl}/farm/api/reward`,{
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

export const getLendList = () => {
  return fetch(
      `${apiUrl}/lend/api/lendlist`,{
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

export const postLendSubscribe = (params) => {
  return fetch(
      `${apiUrl}/lend/api/subscribe`,{
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

export const getLendOrder = (lend, address) => {
  return fetch(
      `${apiUrl}/lend/api/orders`,{
      method: 'POST',
      body: JSON.stringify({
        lend:lend,
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

export const postLendRedeem = (params) => {
  return fetch(
      `${apiUrl}/lend/api/redeem`,{
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

export const getLendingOrder = (lend, address) => {
  return fetch(
      `${apiUrl}/lend/api/loanOrders`,{
      method: 'POST',
      body: JSON.stringify({
        lend:lend,
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

export const postLendLiquidate = (params) => {
  return fetch(
      `${apiUrl}/lend/api/liquidate`,{
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

export const getLoanList = () => {
  return fetch(
      `${apiUrl}/loan/api/pairlist`,{
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

export const getLoanDealPrice = (symbol) => {
  return fetch(
      `${apiUrl}/loan/api/dealPrice`,{
      method: 'POST',
      body: JSON.stringify({
        symbol:symbol
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

export const postLoanToken = [REDACTED] => {
  return fetch(
      `${apiUrl}/loan/api/loan`,{
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

export const getLoanOrder = (lend, address) => {
  return fetch(
      `${apiUrl}/loan/api/orders`,{
      method: 'POST',
      body: JSON.stringify({
        lend:lend,
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

export const postRepayLoan = (params) => {
  return fetch(
      `${apiUrl}/loan/api/repay`,{
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