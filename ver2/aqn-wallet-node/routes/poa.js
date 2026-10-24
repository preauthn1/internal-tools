var express = require('express');
var router = express.Router();
const Web3 = require('web3');
const BN = require('bn.js');

const {paddingLeftZero, sha256, sha512} = require('./utils')
const {NODE} = require('./config')

const emailContract = NODE.emailContract
const recoveryWalletContract = NODE.recoveryWalletContract

const {getSideChainEthCall, getSideChainEthNonce, sendSideChainEthRawTx, getSideChainEthTx, sendEmail} = require('./getSideChainInfuraApi')
const {createResponse} = require('./response')
const {getCosignerAddress, getCosignerPky, aesEncrypt} = require('./utils')

const poaChainId = NODE.CHAIN_ID

const gas = 1000000

const gasPrice = '0x9502F900'  // 2.5gwei

const coreWalletFactory = NODE.coreWalletFactory

const recoveryWalletAddress = NODE.recoveryWalletAddress
const recoveryWalletPrivatekey = NODE.recoveryWalletPrivatekey

var emailValidateTimesMap = new Map();

router.get('/', function(req, res, next) {
  res.send('Hello POA 1');
});

const generateValidateCode = (authorizedAddress, email) => {
  var now = Date.now();

  console.log('now '+now)
  //ts取到百秒
  var ts = now.toString().substr(0, 8)
  
  var code = paddingLeftZero(6, Web3.utils.toBN(
    "0x"+sha512(sha512(authorizedAddress)+ts).substr(0, 40)
  , 16).mod(new BN("1000000")).toString(10)) 
  
  console.log('code '+code)

  return new Promise((resolve, reject) => {
    console.log('Initial');
    const result = sendEmail(email, code)

    resolve(result);
  })
}

const validateEmailCode = (authorizedAddress, emailCode) => {
  var now = Date.now();

  //ts取到百秒 驗證碼有效時間300秒
  var ts_0 = now.toString().substr(0, 8)

  var code_0 = paddingLeftZero(6, Web3.utils.toBN(
    "0x"+sha512(sha512(authorizedAddress)+ts_0).substr(0, 40)
  , 16).mod(new BN("1000000")).toString(10)) 

  var ts_1 = (Number.parseInt(ts_0, 10) - 1 ).toString()
  var code_1 = paddingLeftZero(6, Web3.utils.toBN(
    "0x"+sha512(sha512(authorizedAddress)+ts_1).substr(0, 40)
  , 16).mod(new BN("1000000")).toString(10)) 

  var ts_2 = (Number.parseInt(ts_0, 10) - 2 ).toString()
  var code_2 = paddingLeftZero(6, Web3.utils.toBN(
    "0x"+sha512(sha512(authorizedAddress)+ts_2).substr(0, 40)
  , 16).mod(new BN("1000000")).toString(10)) 

  if(emailCode == code_0 || emailCode == code_1 || emailCode == code_2) {
    return true
  } else {
    return false
  }

}


router.post('/api/email/send', (req, res) => {

  const authorizedAddress = req.body.authorizedAddress
  const emailHex = req.body.emailHex

  if(!emailHex | !authorizedAddress) {
    res.send(createResponse(false, "!emailHex | !authorizedAddress"))

    return
  }

  const email =  Buffer.from(emailHex, "hex").toString("ascii").toLowerCase();
  
  generateValidateCode(authorizedAddress, email)
  .then((result) => {           
      if(result) {
        res.send(createResponse(true, ''))
      } else {
        res.send(createResponse(false, ''))
      }
  }).catch((err)=>{
      console.log("err "+err)
      res.send(createResponse(false, 'error '+err))

  });
});

router.post('/api/wallet/start', (req, res) => {
  const authorizedAddress = req.body.authorizedAddress
  const emailHex = req.body.emailHex
  const code = req.body.code

  if(!emailHex | !authorizedAddress | !code) {
    res.send(createResponse(false, "!emailHex | !authorizedAddress | !code"))

    return
  }

  const isValidate = validateEmailCode(authorizedAddress, code)

  if(!isValidate) {
    console.log('emailValidateTimes '+emailValidateTimesMap.get(emailHex))

    if(emailValidateTimesMap.get(emailHex)) {
      var times = emailValidateTimesMap.get(emailHex) + 1
      emailValidateTimesMap.set(emailHex, times)

    } else {
      emailValidateTimesMap.set(emailHex, 1)
    }

    res.send(createResponse(false, "!isValidate"))

    return
  }

  //max validate times:10
  if(emailValidateTimesMap.get(emailHex)>10) {
    res.send(createResponse(false, "validateOverflow"))

    return
  }

  const email =  Buffer.from(emailHex, "hex").toString("ascii").toLowerCase();
  const emailHash = sha256(email)

  new Promise((resolve, reject) => {
    const callData = "0x45b7b455"+emailHash
    const result = getSideChainEthCall(emailContract, callData)
    resolve(result)
    
  }).then((result) => {
    console.log('result '+result)


    var walletAddress = "0x"+result.substr(26, 40)
    if(walletAddress==0x0000000000000000000000000000000000000000) {
      poaRegister(req, res)
      return
    }

    const data = "0xd878efb0"
      + paddingLeftZero(64, walletAddress.toLowerCase().substr(2))
      + paddingLeftZero(64, authorizedAddress.toLowerCase().substr(2))
      + paddingLeftZero(64, "1")

    console.log('data '+data)

    recoveryWallet(data, walletAddress, res)

  }).catch(err=>{
    res.send(createResponse(false, err))
  })

});


const CoreWalletFactoryABI = [
  {
    "constant":false,
    "inputs":[
    {
    "internalType":"address",
    "name":"_authorizedAddress",
    "type":"address"
    },
    {
    "internalType":"bytes32",
    "name":"_emailHash",
    "type":"bytes32"
    },
    {
    "internalType":"bytes",
    "name":"_email",
    "type":"bytes"
    }
    ],
    "name":"deployCoreWallet",
    "outputs":[
    ],
    "payable":false,
    "stateMutability":"nonpayable",
    "type":"function"
    }
]

const poaRegister = (req, res) => {
  const authorizedAddress = req.body.authorizedAddress
  const emailHex = req.body.emailHex

  const email =  Buffer.from(emailHex, "hex").toString("ascii").toLowerCase();

  const emailHash = sha256(email)
  const emailAes = aesEncrypt(email)

  var web3 = new Web3(Web3.givenProvider)
  var contract = new web3.eth.Contract(CoreWalletFactoryABI, coreWalletFactory)
  var contractData = contract.methods.deployCoreWallet(
    authorizedAddress,
    "0x"+emailHash,
    "0x"+emailAes
  ).encodeABI()

  const params = contractData
  console.log('params '+params)

  var coSignerAddress = getCosignerAddress()

  new Promise((resolve, reject) => {
    console.log('Initial');

    const nonce = getSideChainEthNonce(coSignerAddress)
    resolve(nonce);
  })
  .then((nonce) => {              
      console.log('nonce '+nonce);
      
      var web3 = new Web3(Web3.givenProvider)
      var coSignerPky = getCosignerPky(coSignerAddress)

      return web3.eth.accounts.signTransaction({
          to: coreWalletFactory,
          value: 0x0,
          gas: gas,
          gasPrice: gasPrice,
          nonce: nonce,
          chainId: poaChainId,
          data: params,
      }, coSignerPky)
      .then((result)=>{
        console.log(result)
        return result.rawTransaction
      });

  })
  .then((param) => {
      console.log('broadcastParams '+param);

      return sendSideChainEthRawTx(param)
  }).then((txid) => {
    console.log('txid '+txid);

    const topics = "0x4459976dae5d173eda0ac09c39b840bed814c9d1e490e8b9ba23148e4b19a837"
    listenTx(res, txid, (logs)=>{
      logs.map(log=>{
        if(log.topics[0]==topics) {
          const poaAddress = "0x"+log.data.substr(26, 40)
          res.send(createResponse(true, poaAddress))
        }
      })
    })
  });
}

const recoveryWallet = (data, walletAddress, res) => {
  new Promise((resolve, reject) => {
    console.log('Initial');
    const nonce = getSideChainEthNonce(recoveryWalletAddress)
    resolve(nonce);
  })
  .then((nonce) => {              
      console.log('nonce '+nonce);
      
      var web3 = new Web3(Web3.givenProvider)

      return web3.eth.accounts.signTransaction({
          to: recoveryWalletContract,
          value: 0x0,
          gas: gas,
          gasPrice: gasPrice,
          nonce: nonce,
          chainId: poaChainId,
          data: data,
      }, recoveryWalletPrivatekey)
      .then((result)=>{
        console.log(result)
        return result.rawTransaction
      });

  })
  .then((param) => {
      console.log('broadcastParams '+param);

      return sendSideChainEthRawTx(param)
  }).then((txid) => {
    console.log('txid '+txid);

    listenTx(res, txid, (logs)=>{
      res.send(createResponse(true, walletAddress))
    })

  }).catch((err)=>{
      console.log("err "+err)
      res.send(createResponse(false, 'broadcast error '+err))

  });
}

const coSignerContractABI = [
  {
    "constant":false,
    "inputs":[
    {
    "internalType":"address",
    "name":"_walletAddress",
    "type":"address"
    },
    {
    "internalType":"uint8",
    "name":"v",
    "type":"uint8"
    },
    {
    "internalType":"bytes32",
    "name":"r",
    "type":"bytes32"
    },
    {
    "internalType":"bytes32",
    "name":"s",
    "type":"bytes32"
    },
    {
    "internalType":"uint256",
    "name":"nonce",
    "type":"uint256"
    },
    {
    "internalType":"address",
    "name":"toAddress",
    "type":"address"
    },
    {
    "internalType":"bytes",
    "name":"data",
    "type":"bytes"
    }
    ],
    "name":"callWalletContract",
    "outputs":[
    ],
    "payable":false,
    "stateMutability":"nonpayable",
    "type":"function"
  }
]

router.post('/api/invoke', (req, res) => {
  const walletAddress = req.body.walletAddress
  const v = req.body.v
  const r = req.body.r
  const s = req.body.s
  const walletNonce = req.body.nonce
  const data = req.body.data
  const authorizedAddress = req.body.authorizedAddress
  const toAddress = req.body.toAddress

  //check

  var coSignerAddress = getCosignerAddress()

  new Promise((resolve, reject) => {
      console.log('Initial');
        
      const nonce = getSideChainEthNonce(coSignerAddress)
      console.log('coSignerAddress '+coSignerAddress);

      resolve(nonce);
  })
  .then((nonce) => {              

      var coSignerPky = getCosignerPky(coSignerAddress)
      
      var web3 = new Web3(Web3.givenProvider)
      var contract = new web3.eth.Contract(coSignerContractABI, NODE.coSignerContractAddress)
      var contractData = contract.methods.callWalletContract(
          walletAddress,
          v, 
          r,
          s,
          walletNonce,
          toAddress,
          data
          ).encodeABI()

      console.log('contractData '+contractData);
      
      
      return web3.eth.accounts.signTransaction({
          to: NODE.coSignerContractAddress,
          value: 0x0,
          gas: gas,
          gasPrice: gasPrice,
          nonce: nonce,
          chainId: poaChainId,
          data: contractData,
      }, coSignerPky)
      .then((result)=>{
        console.log(result)
        return result.rawTransaction
      });
      
  })
  .then((param) => {
    console.log('broadcastParams '+param);

    return sendSideChainEthRawTx(param)
  }).then((txid) => {
    console.log('txid '+txid);

    // listenTx(res, txid, (logs)=>{
    //   res.send(createResponse(true, ""))
    // })

    res.send(createResponse(true, txid))

  }).catch((err)=>{
      console.log("err "+err)
      res.send(createResponse(false, 'broadcast error '+err))

  });
})

const listenTx = (res, txid, callback) => {
  var times = 0
  var myInterval = setInterval(()=>{
    console.log('listenTx '+times)

    
    new Promise((resolve, reject) => {
      const result = getSideChainEthTx(txid)
      resolve(result);
    }).then((result) => {

      if (!result) {
        return
      }

      if (result.status!="0x1") {
        clearTimeout(myInterval);
        return
      }
 
      callback(result.logs)
        
      clearTimeout(myInterval);
    })

    times += 1

    if (times == 30) {
      console.log('listenTxid stop')
      res.send(createResponse(false, ""))
      clearTimeout(myInterval);
      
    }
  },1000,"Interval");
}

router.post('/api/wallet/nonce', (req, res)=> {

  const walletAddress = req.body.walletAddress
  const authorizedAddress = req.body.authorizedAddress


  new Promise((resolve, reject) => {
    const result = getSideChainEthCall(walletAddress, "0x6e05dc11000000000000000000000000"+authorizedAddress.substr(2).toLowerCase())
    resolve(result);
  })
  .then((result) => {
    
    res.send(createResponse(true, result))
  })
})

router.post('/api/eth/nonce', (req, res) => {

  const ethAddress = req.body.ethAddress
  
  new Promise((resolve, reject) => {
    console.log('Initial');

    const nonce = getSideChainEthNonce(ethAddress)
    resolve(nonce);
  })
  .then((nonce) => {              
    console.log('nonce '+nonce);
      
    res.send(createResponse(true, nonce))

  }).catch((err)=>{
    console.log("err "+err)
    res.send(createResponse(false, 'broadcast error '+err))

});
})

router.post('/api/balance', (req, res)=> {

  const walletAddress = req.body.walletAddress
  const tokenAddress = req.body.tokenAddress

  new Promise((resolve, reject) => {
    const result = getSideChainEthCall(tokenAddress, "0x70a08231000000000000000000000000"+walletAddress.substr(2).toLowerCase())
    resolve(result);
  })
  .then((result) => {
      res.send(createResponse(true, Web3.utils.hexToNumberString(result)))
  })
})

module.exports = router;
