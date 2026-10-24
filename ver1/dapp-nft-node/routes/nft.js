var express = require('express');
var router = express.Router();
const Web3 = require('web3');
const BN = require('bn.js');

const {redirectUrl, getSideChainEthCall, getSideChainEthTx} = require('./getSideChainInfuraApi')
const {createResponse} = require('./response')
const {paddingLeftZero} = require('./utils')
const {NODE} = require('./config');


/* GET home page. */
router.get('/', function(req, res, next) {
  res.send('Hello NFT 1');
});

router.post('/api/balance', (req, res)=> {

  const walletAddress = req.body.walletAddress
  const tokenAddress = req.body.tokenAddress

  new Promise((resolve, reject) => {
    const result = getSideChainEthCall(tokenAddress, "0x70a08231"+paddingLeftZero(64, walletAddress.substr(2).toLowerCase()))
    resolve(result);
  })
  .then((result) => {
      res.send(createResponse(true, Web3.utils.hexToNumberString(result)))
  })
})

router.post('/api/tokenAllowance', (req, res)=> {

  const walletAddress = req.body.walletAddress
  const tokenAddress = req.body.tokenAddress
  const spenderAddress = req.body.spenderAddress

  new Promise((resolve, reject) => {
    const result = getSideChainEthCall(tokenAddress, "0xdd62ed3e"+paddingLeftZero(64, walletAddress.substr(2).toLowerCase())+paddingLeftZero(64, spenderAddress.substr(2).toLowerCase()))
    resolve(result);
  })
  .then((result) => {
      res.send(createResponse(true, Web3.utils.hexToNumberString(result)=="0"?false:true))
  })
})

router.post('/api/txid', (req, res)=> {

  const txid = req.body.txid

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

      res.send(createResponse(true, result))

      clearTimeout(myInterval);
    })

    times += 1

    if (times == 60) {
      console.log('listenTxid stop')
      res.send(createResponse(false, false))
      clearTimeout(myInterval);
      
    }
  },500,"Interval");
})

router.post('/api/nftInfo', (req, res)=> {

  const nftContract = req.body.nftContract
  const nftId = req.body.nftId

  const nftIdHex = parseInt(nftId).toString(16);

  console.log('nftIdHex '+nftIdHex)

  var callList = []
  callList.push(getSideChainEthCall(nftContract, "0x6352211e"+paddingLeftZero(64, nftIdHex))) //owner
  callList.push(getSideChainEthCall(nftContract, "0xd46727e1"+paddingLeftZero(64, nftIdHex))) //face
  callList.push(getSideChainEthCall(nftContract, "0x2f96f1ad"+paddingLeftZero(64, nftIdHex))) //dice number
  callList.push(getSideChainEthCall(nftContract, "0xbe5e9418"+paddingLeftZero(64, nftIdHex))) //payout address

  Promise.all(callList).then(values=>{
    var owner = "0x"+values[0].substr(26, 40)
    var face = Web3.utils.toBN(values[1]).toString(10)
    var diceNumber = Web3.utils.toBN(values[2]).toString(10)
    var payoutAddress = "0x"+values[3].substr(26, 40)

    res.send(createResponse(true, {
      owner:owner,
      face:face,
      diceNumber:diceNumber,
      payoutAddress:payoutAddress
    }))

  }).catch(err=>{
    console.log(err)
    res.send(createResponse(false, err))
  })
})

router.get('/metadata/:pair', function(req, res, next) {

  const path = req.path.split("/")
  const pair = path[2]

  redirectUrl("https://example.invalid/dice/"+pair).then(data=>{
    res.send(data);
  }).catch(err=>{
    res.send("");
  })
});

module.exports = router;
