import * as bitcoin from 'bitcoinjs-lib'
import * as Web3 from 'web3'
import * as bip39 from 'bip39';
import * as bip32 from 'bip32';
import * as secp256k1 from 'secp256k1'
import { v4 as uuidv4 } from 'uuid';

import db from './model/db'
import { getPoaWalletNonce, getSideChainEthNonce } from './api/getMainnetInfuraApi'

import { encryptMsg, decryptMsg } from './Utils'

const gas = 1000000
const gasPrice = 0  // 1 Gwei
const poaChainId = 55

export const createNewWallet = () => {

    return new Promise((resolve, reject) => {
        console.log('Initial ');

        if (typeof WebAssembly !== "object" && typeof WebAssembly.instantiate !== "function") {
          return
        }

        var mnemonics = bip39.generateMnemonic()

        if(mnemonics === "") {
          return
        }

        const seed = bip39.mnemonicToSeedSync(mnemonics);
        const root = bip32.fromSeed(seed);
        const path = "m/44'/60'/0'/0/0";
        const child = root.derivePath(path);
        
        var uuid = uuidv4()

        console.log('uuid '+uuid)
        uuid = uuid.replace("-", "")
        uuid = uuid.replace("-", "")
        uuid = uuid.replace("-", "")
        uuid = uuid.replace("-", "")

        const ecPair = bitcoin.ECPair.fromPrivateKey(child.privateKey, {compressed: false})
        const privateKey = ecPair.privateKey.toString('hex')
        const epkey = encryptMsg(privateKey, uuid);
        const publicKey = ecPair.publicKey.toString('hex').substr(2)
        console.log('publicKey '+ publicKey)

        var web3 = new Web3(Web3.givenProvider)
        const ethAddress = web3.eth.accounts.privateKeyToAccount('0x'+privateKey).address
        console.log('ethAddress '+ ethAddress)

        db.table_wallet.clear().then(()=>{
          console.log("table_wallet successfully deleted");

          return db.table_wallet.put({
            name:"ETH 錢包", 
            mainCoinId:"0x0000000000000000000000000000000000000000", 
            chainType:1, 
            epkey:epkey, 
            address:ethAddress, 
            isFromSystem:true, 
            uuid:uuid})
        }).then(()=>{          
          resolve(ethAddress);
        })

    })   
}

export const poaSign = (signData) => {
  const poaAddress = localStorage.getItem('poaAddress_aqn')

  return new Promise((resolve, reject) => {
          console.log('Initial');

          db.table_wallet.where('chainType').anyOf(1).first((wallet)=>{
            console.log(wallet.address)

            resolve(wallet);
          })
      })
      .then((wallet)=> {

        var pwd = [REDACTED]

        return getPoaWalletNonce(poaAddress, wallet.address).then(nonce=>{
          const epKey = decryptMsg(wallet.epkey, pwd)

          const message = "0x1900"+poaAddress.toLowerCase().substr(2)+nonce.substr(2)+signData
          console.log('message '+message)
          return {
            authorizedAddress:wallet.address,
            epKey:"0x"+epKey,
            nonce:Web3.utils.toBN(nonce).toString(10),
            message:message
          }
        })
        
      })
      .then((params) => {              
          
          var web3 = new Web3(Web3.givenProvider)


          var hash = web3.utils.soliditySha3(params.message)
          const sig = secp256k1.ecdsaSign(Buffer.from(web3.utils.hexToBytes(hash)), Buffer.from(web3.utils.hexToBytes(params.epKey)))

          return {
            walletAddress:poaAddress,
            authorizedAddress:params.authorizedAddress,
            v:sig.recid + 27,
            r:"0x"+Buffer.from(sig.signature.slice(0, 32)).toString('hex'),
            s:"0x"+Buffer.from(sig.signature.slice(32, 64)).toString('hex'),
            nonce:params.nonce,
            data:"0x"+signData
          }
          
      })
      .then((param) => {
          console.log('broadcastParams '+JSON.stringify(param));

          return param
      }).catch((err)=>{
        console.log("err "+err)
      });
}

export const poaNormalSign = (signData, toAddress) => {

  return new Promise((resolve, reject) => {
          console.log('Initial');

          db.table_wallet.where('chainType').anyOf(1).first((wallet)=>{
            console.log(wallet.address)

            resolve(wallet);
          })
      })
      .then((wallet)=> {

        var pwd = [REDACTED]

        return getSideChainEthNonce(wallet.address).then(nonce=>{
          const epKey = decryptMsg(wallet.epkey, pwd)
          return {
            authorizedAddress:wallet.address,
            epKey:"0x"+epKey,
            nonce:Web3.utils.toBN(nonce).toString(10),
            message:signData
          }
        })
        
      })
      .then((params) => {              
          
          var web3 = new Web3(Web3.givenProvider)
          
          return web3.eth.accounts.signTransaction({
            to: toAddress,
            value: 0x0,
            gas: gas,
            gasPrice: 0,
            nonce: params.nonce,
            chainId: poaChainId,
            data: params.message,
          }, params.epKey)
          .then((result)=>{
            console.log(result)
            return result.rawTransaction
          });
      })
      .then((param) => {
          console.log('broadcastParams '+JSON.stringify(param));

          return param
      }).catch((err)=>{
        console.log("err "+err)
      });
};