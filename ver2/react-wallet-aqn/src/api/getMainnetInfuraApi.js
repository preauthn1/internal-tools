// import { Eth } from "web3-eth";
// import { paddingLeftZero, showBalanceDigit } from '../Utils'
// import * as BigNumber from 'bignumber.js';
// import * as Web3 from 'web3'

const apiUrl = "http://0.0.0.0:4001"
const apiReadUrl = "http://0.0.0.0:4001"

// const apiUrl = "http://0.0.0.0:4001"
// const apiReadUrl = "http://0.0.0.0:4001"

const pbsUrl = "https://example.invalid"
// const pbsUrl = "http://0.0.0.0:4000"

export const postWalletInvoke = (params) => {
  return fetch(
      `${apiUrl}/poa/api/invoke`,{
      method: 'POST',
      body: params,
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
}

export const getTokenBalance = (address, tokenAddress) => {
    return fetch(
        `${apiReadUrl}/poa/api/balance`,{
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

export const getTokenBalanceWithUsd = (address, tokenAddress) => {
  return fetch(
      `${apiReadUrl}/poa/api/balanceWithUsd`,{
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

export const getTokenBalanceAllWithUsd = (address, tokenAddress) => {
  return fetch(
      `${apiReadUrl}/poa/api/balanceAll`,{
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

export const postWalletRegister = (address) => {
    return fetch(
        `${apiUrl}/poa/api/wallet/register`,{
        method: 'POST',
        body: JSON.stringify({
            authorizedAddress:address,
            referenceCode:"admin"
        }),
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

export const postPbsGetChainAddress = (address, tokenAddress) => {
  return fetch(
      `${pbsUrl}/api/chainAddress`,{
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
  .then((data) => {

  return data
  });
};

export const postPbsRegisterChainAddress = (params) => {
    return fetch(
        `${pbsUrl}/api/register`,{
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

export const postPbsGetTokenTransferList = (address, tokenAddress) => {
  return fetch(
      `${pbsUrl}/api/transfer/list`,{
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
  .then((data) => {

  return data
  });
};

export const postPbsGetTokenWithdrawList = (address, tokenAddress) => {
  return fetch(
      `${pbsUrl}/api/withdraw/list`,{
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
  .then((data) => {

  return data
  });
};

export const postPbsChainAddressSearch = (params) => {
  return fetch(
      `${pbsUrl}/api/chainAddressSearch`,{
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

export const postPbsTokenTransfer = (params) => {
  return fetch(
      `${pbsUrl}/api/transfer`,{
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

export const postPbsGetUsdtTxList = (params) => {
    return fetch(
        `${pbsUrl}/usdt/api/tx/usdt`,{
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

export const postPbsGetUsdtDepositList = (address) => {
    return fetch(
        `${pbsUrl}/usdt/api/deposit/list`,{
        method: 'POST',
        body: JSON.stringify({
            walletAddress:address,
        }),
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

export const postPbsApplyUsdtDeposit = (params) => {
    return fetch(
        `${pbsUrl}/usdt/api/apply/deposit`,{
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

export const postPbsUsdtWithdrawContract = () => {
  return fetch(
      `${pbsUrl}/usdt/api/withdrawContract`,{
      method: 'POST',
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

export const postPbsWithdraw = (params) => {
  return fetch(
      `${pbsUrl}/api/withdraw`,{
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

export const postPbsWithdrawFee = (token) => {
  return fetch(
      `${pbsUrl}/api/withdrawFee`,{
      method: 'POST',
      body: JSON.stringify({
        token:[REDACTED],
    }),
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

export const postPbsUsdtWithdraw = (params) => {
  return fetch(
      `${pbsUrl}/usdt/api/withdrawV2`,{
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

export const postPbsUsdtWithdrawFee = () => {
  return fetch(
      `${pbsUrl}/usdt/api/withdrawFee`,{
      method: 'POST',
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

export const postPbsGetBtcTxList = (params) => {
  return fetch(
      `${pbsUrl}/btc/api/tx/btc`,{
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

export const postPbsApplyBtcDeposit = (params) => {
  return fetch(
      `${pbsUrl}/btc/api/apply/deposit`,{
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

export const postPbsGetBtcDepositList = (address) => {
  return fetch(
      `${pbsUrl}/btc/api/deposit/list`,{
      method: 'POST',
      body: JSON.stringify({
          walletAddress:address,
      }),
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

export const postPbsBtcWithdraw = (params) => {
  return fetch(
      `${pbsUrl}/btc/api/withdraw`,{
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

export const postPbsBtcWithdrawFee = () => {
  return fetch(
      `${pbsUrl}/btc/api/withdrawFee`,{
      method: 'POST',
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

export const postPbsBtcWithdrawContract = () => {
  return fetch(
      `${pbsUrl}/btc/api/withdrawContract`,{
      method: 'POST',
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

export const postPbsGetEthTxList = (params) => {
  return fetch(
      `${pbsUrl}/eth/api/tx/eth`,{
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

export const postPbsApplyEthDeposit = (params) => {
  return fetch(
      `${pbsUrl}/eth/api/apply/deposit`,{
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

export const postPbsGetEthDepositList = (address) => {
  return fetch(
      `${pbsUrl}/eth/api/deposit/list`,{
      method: 'POST',
      body: JSON.stringify({
          walletAddress:address,
      }),
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

export const postPbsEthWithdrawContract = () => {
  return fetch(
      `${pbsUrl}/eth/api/withdrawContract`,{
      method: 'POST',
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

export const postPbsEthWithdraw = (params) => {
  return fetch(
      `${pbsUrl}/eth/api/withdrawV2`,{
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

export const postPbsEthWithdrawFee = () => {
  return fetch(
      `${pbsUrl}/eth/api/withdrawFee`,{
      method: 'POST',
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

export const postPbsAqnWithdrawContract = () => {
  return fetch(
      `${pbsUrl}/aqn/api/withdrawContract`,{
      method: 'POST',
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

export const postPbsAqnWithdraw = (params) => {
  return fetch(
      `${pbsUrl}/aqn/api/withdrawByUsdtV2`,{
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

export const postPbsAqnWithdrawFee = () => {
  return fetch(
      `${pbsUrl}/aqn/api/withdrawFeeByUsdt`,{
      method: 'POST',
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

export const postPoaTokenTxList = (address, tokenAddress) => {
  return fetch(
      `${apiUrl}/poa/api/txlist`,{
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
  .then((data) => {

  return data
  });
};

export const postPoaPbkCheck = (params) => {
  return fetch(
      `${apiUrl}/id/api/pbk/check`,{
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

export const postPoaEmailCheck = (address) => {
  return fetch(
      `${apiUrl}/id/api/email/check`,{
      method: 'POST',
      body: JSON.stringify({
          walletAddress:address,
      }),
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

export const postPoaEmailSend = (params) => {
  return fetch(
      `${apiUrl}/poa/api/email/send`,{
      method: 'POST',
      body: JSON.stringify(params),
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

export const postPoaWalletStart = (params) => {
  return fetch(
      `${apiUrl}/poa/api/wallet/start`,{
      method: 'POST',
      body: JSON.stringify(params),
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


export const postPoaEmailRecoverySend = (params) => {
  return fetch(
      `${apiUrl}/id/api/email/recovery/send`,{
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

export const postPoaEmailRecovery = (params) => {
  return fetch(
      `${apiUrl}/id/api/email/recovery`,{
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

export const postGetTokenList = () => {
  return fetch(
      `${apiReadUrl}/poa/api/wallet/tokenlist`,{
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

export const postTokenTransfer = (params) => {
  return fetch(
      `${apiUrl}/poa/api/token/transfer`,{
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