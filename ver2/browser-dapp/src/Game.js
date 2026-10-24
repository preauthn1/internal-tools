import React, { useEffect } from 'react'
import IframeComm from "react-iframe-comm";
import { createNewWallet, poaSign, poaNormalSign } from './Blockchain'
import { isValidAddress } from './RuleUtils'
import db from './model/db'

import Box from '@material-ui/core/Box';

var height = document.documentElement.clientHeight
var width  = window.screen.width


const Game = () => {

  var mainUrl = "http://0.0.0.0:3006"

  return (
    <React.Fragment>
      <Box style={{
        width: "100%",
        height: height,
        // backgroundColor:"#000000",
      }}>

        <IframeComm
              attributes={{
                id: "iframe",
                src: mainUrl,
                position: 'fixed',
                width:'100%',
                top: 'auto',
                height: height,
              }}
              postMessageData={"test99 hello iframe"}
              handleReady={()=>{
                console.log("test99 onReady");
              }}
              handleReceiveMessage={(event)=>{
                if (event.data.method) {
                  var iframe = document.getElementById('iframe').contentWindow

                  console.log("test99 onReceiveMessage "+JSON.stringify(event.data))

                  const request = event.data.data

                  switch (event.data.method) {
                    case "getWallet":
                      const poaAddress = localStorage.getItem('poaAddress_aqn')
                      if (poaAddress && isValidAddress(poaAddress)) {
                        iframe.postMessage({method:"returnWallet", data:poaAddress}, '*');
                      } else {
                        iframe.postMessage({method:"returnWallet", data:"EMPTY"}, '*');
                      }
                      break
                    case "getAuthAddress":

                      db.table_wallet.where('chainType').anyOf(1).first((wallet)=>{            
                        if(wallet) {
                          iframe.postMessage({method:"returnAuthAddress", data:wallet.address}, '*');
                        } else {
                          iframe.postMessage({method:"returnAuthAddress", data:""}, '*');
                        }
                      })
                      
                      break
                    case "generatePky":
                      createNewWallet().then(ethAddress=>{
                        console.log('createNewWallet '+ethAddress)

                        iframe.postMessage({method:"returnAuthAddress", data:ethAddress}, '*');
                      }).catch((err)=>{
                        console.log("err "+err)
                      })
                      break
                    case "setPoaAddress":
                      window.localStorage.setItem('poaAddress_aqn', request);
                      break
                    case "signWalletData":
                      poaSign(request).then(signature=>{
                        if(signature == "") {
                          return
                        }
                        console.log('ethSign '+JSON.stringify(signature))

                        iframe.postMessage({method:"returnWalletSignature", data:signature}, '*');
                      }).catch((err)=>{
                        console.log("err "+err)
                      })
                      break
                    case "signNormalData":
                      poaNormalSign(request.signData, request.toAddress).then(signature=>{
                        if(signature == "") {
                          return
                        }
                        // console.log('ethNormalSign '+JSON.stringify(signature))

                        iframe.postMessage({method:"returnNormalSignature", data:signature}, '*');
                      }).catch((err)=>{
                        console.log("err "+err)
                      })
                      break
                  }
                }
              }}
        />
      </Box>
    </React.Fragment>

  );
}

export default Game;
