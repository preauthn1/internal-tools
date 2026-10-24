import React, { useEffect } from 'react'
import * as BigNumber from 'bignumber.js';
import { useTranslation } from 'react-i18next';

import { withStyles } from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import Typography from '@material-ui/core/Typography';
import Box from '@material-ui/core/Box';
import CircularProgress from '@material-ui/core/CircularProgress';
import TextField from '@material-ui/core/TextField';

import BottomNavigation from '@material-ui/core/BottomNavigation';
import BottomNavigationAction from '@material-ui/core/BottomNavigationAction';
import Paper from '@material-ui/core/Paper';

import { paddingLeftZero, showBalanceDigit } from './Utils'
import { isValidNftId, isValidAmount, isValidAddress } from './RuleUtils'

import { getTokenBalance, getTxid, getNftInfo, getTokenAllowance }  from './api/getMainnetInfuraApi'

import * as Web3 from 'web3'

var height = document.documentElement.clientHeight
var width  = window.screen.width

const BottomNavigationActionButton = withStyles(theme => ({
  root: {
    color: "#757575"
  },
  selected: {
     color: "#FFB41D"
  }
}))(BottomNavigationAction);

const GreenButton = withStyles(theme => ({
  root: {
    height: 32,
    width: 120,
    backgroundColor: "rgb(2, 192, 118)",
    '&:hover': {
      backgroundColor: "rgb(2, 192, 118)",
    },
    borderRadius:16,
  },
  label: {
      color: '#ffffff',
      fontSize: '10px',
  },
}))(Button);

const RedButton = withStyles(theme => ({
  root: {
    height: 32,
    width: 120,
    backgroundColor: "#f84960",
    '&:hover': {
      backgroundColor: "#f84960",
    },
    borderRadius:16,
  },
  label: {
      color: '#ffffff',
      fontSize: '10px',
  },
}))(Button);

const InputContainer = withStyles(theme => ({
  root: {
      height: 44,
      width: 282,
      backgroundColor: "#333333",
      borderRadius:18,
  },
  label: {
      color: '#ffffff',
  },
}))(Paper);

const LoadingConainer = ({title, loadingDialog}) => {
  return  <Dialog  open={loadingDialog}>
      <Box style={{
          flexDirection: "column",
          display: 'flex',
          alignItems:"center",
          flexGrow: 1,
          paddingTop:16,
          paddingLeft:32,
          paddingRight:32,
          paddingBottom:32,
      }}>
          <Typography>
              {title}
          </Typography>
          <CircularProgress style={{
              marginTop:16,
          }}/>
      </Box>
  </Dialog>
}

const Game = () => {

  const { t, i18n } = useTranslation();

  const [walletAddress, setWalletAddress] = React.useState("");
  const [page, setPage] = React.useState(0);
  const [loadingDialog, setLoadingDialog] = React.useState(false);

  const [balance, setBalance] = React.useState(0);

  const [nftId, setNftId] = React.useState("0");
  const [nftInfo, setNftInfo] = React.useState({
    owner: "",
    face: "",
    diceNumber: "",
    payoutAddress: "",
    metaData: ""
});

  const [betTxid, setBetTxid] = React.useState("");
  const [betResult, setBetResult] = React.useState("");

  const [mintFace, setMintFace] = React.useState("");
  const [mintDiceNumber, setMintDiceNumber] = React.useState("");
  const [mintPayoutAddress, setMintPayoutAddress] = React.useState("");

  const [mintTxid, setMintTxid] = React.useState("");
  const [mintResult, setMintResult] = React.useState({
    nftId:"",
    owner: "",
    face: "",
    diceNumber: "",
    payoutAddress: ""
  });

  const [isApprove, setIsApprove] = React.useState(false);

  //Rinkeby
  const tokenAddress = "0x1afCa6ac39043171bb39BB1ED2f3699919a3134a"
  const tokenDigit = 18
  const nftContract = "0xD28a63Ca24d6217a3C378A1d1C62eB6B281D5302"

  const updateBalance = (address) => {
    getTokenBalance(address, tokenAddress).then(data=>{
      setBalance(data)
    }).catch(err=>{
      console.log('err '+err)
    })
  }

  const checkTokenAllowance = (address) => {
    getTokenAllowance(address, tokenAddress, nftContract).then(data=>{
      if(data) {
        setIsApprove(true)
      } else {
        setIsApprove(false)
      }
    }).catch(err=>{
      console.log('err '+err)
    })
  }

  const updateNftInfo = (nftId) => {
    getNftInfo(nftContract, nftId).then(data=>{
      
      if(data.owner=="0x0000000000000000000000000000000000000000") {
        alert("NFT Not Found")
      } else {
        setNftInfo(data)
      }
    }).catch(err=>{
      console.log('err '+err)
    })
  }

  const bet = (params) => {

    if(nftInfo.owner=="") {
      alert("NFT Not Found")
      return
    }

    setBetTxid("")
    setBetResult("")

    getTokenAllowance(walletAddress, tokenAddress, nftContract).then(isAllowance=>{
      if(isAllowance) {
        setLoadingDialog(true)
        window.ethereum.request({
          method: 'eth_sendTransaction',
          params,
        }).then((txid) => {
          getTxid(txid).then(result=>{
            setLoadingDialog(false)
            if(result){
              updateBalance(walletAddress)
    
              setBetTxid(result.transactionHash)
    
              const rollResult = result.logs.find(log=>log.topics[0]=="0xcda0c777d73e8b7526efa7ef0ce22a32bb922c9d790cd01a73754275fdf79e15").data
    
              const ts = new BigNumber(rollResult.substr(0, 66)).toNumber()
    
              const diceResult = rollResult.substr(194, rollResult.length)
              const diceNumber = Math.floor(diceResult.length/132)
    
              var index
              var diceResultString = ""
              for(index=0; index<diceNumber; index++) {
                var dicePoint = new BigNumber("0x"+diceResult.substr(132*index, 64)).toNumber()
                var basePoint = new BigNumber("0x"+diceResult.substr(132*index+66, 64)).toNumber()
    
                diceResultString += dicePoint + String.fromCharCode(parseInt(diceResult.substr(132*index+64, 2), 16))
                 + basePoint + String.fromCharCode(parseInt(diceResult.substr(132*index+130, 2), 16))
              }
              
              diceResultString = diceResultString.substr(0, diceResultString.length-1)
              
              setBetResult(diceResultString)
            }
          })
        }).catch((error) => {
    
        });
      } else {
        approveToken()
      }
    }).catch(err=>{
      console.log('err '+err)
    })




  }

  const mintNft = (params) => {
    setBetTxid("")
    setBetResult("")

    getTokenAllowance(mintPayoutAddress, tokenAddress, nftContract).then(isAllowance=>{
      if(isAllowance) {
        setLoadingDialog(true)
        window.ethereum.request({
          method: 'eth_sendTransaction',
          params,
        }).then((txid) => {
          getTxid(txid).then(result=>{
            setLoadingDialog(false)
            if(result){
              setMintTxid(result.transactionHash)
    
              const mintResultLog = result.logs.find(log=>log.topics[0]=="0x94c792774c59479f7bd68442f3af3691c02123a5aabee8b6f9116d8af8aa6669")
    
    
              const mintNftId = new BigNumber(mintResultLog.topics[2]).toString(10)
              
              getNftInfo(nftContract, mintNftId).then(data=>{
                setMintResult({
                  nftId:mintNftId,
                  owner: data.owner,
                  face: data.face,
                  diceNumber: data.diceNumber,
                  payoutAddress: data.payoutAddress
                })
              }).catch(err=>{
                console.log('err '+err)
              })
    
            }
          })
        }).catch((error) => {
    
        });
      } else {
        alert("PayoutAddress Need Approve")
      }
    }).catch(err=>{
      console.log('err '+err)
    })

    
  }

  const approveToken = () => {
    
    var params = [
      {
        from: walletAddress,
        to: tokenAddress,
        gas: '0xc350', // 50000
        gasPrice: '0x9502F900', // 2.5gwei
        value: '0x0',
        data:
          '0x095ea7b3'+
          paddingLeftZero(64, nftContract.substr(2))+
          "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
      },
    ];

    setLoadingDialog(true)
    window.ethereum.request({
      method: 'eth_sendTransaction',
      params,
    }).then((txid) => {
      getTxid(txid).then(result=>{
        setLoadingDialog(false)
        if(result){
          checkTokenAllowance(walletAddress)
        }
      })
    }).catch((error) => {

    });
  }

  useEffect(
    () => {

      if (window.ethereum) {
        window.web3 = new Web3(window.ethereum);

        window.ethereum.request({ method: 'eth_requestAccounts' }).then(accounts=>{
          if(accounts) {
            setWalletAddress(accounts[0])
            updateBalance(accounts[0])
            checkTokenAllowance(accounts[0])
          }
        }).catch(err=>{
  
        })
      }

      const params = (new URL(document.location)).searchParams;
      const id = params.get("id");
      
      if (id) {

        if(isValidNftId(id)) {
          setNftId(id)
          updateNftInfo(id)
        } else {
          setNftId("0")
          updateNftInfo("0")
        }
      } else {
        setNftId("0")
        updateNftInfo("0")
      }

    },[]
  )

  return (
    <React.Fragment>
      <Box style={{
        width:"100%",
        height: {height},
        backgroundColor:"#ffffff",
        backgroundRepeat:"no-repeat",
        backgroundSize:"100% 100%",
        flexDirection: "column",
        display: "flex",
      }}>

        <Box style={{
          position: 'fixed',
          width:'100%',
          height: height-36,
          background:"linear-gradient(#000000,#052F41)",
          backgroundRepeat:"no-repeat",
          backgroundSize:"100% 100%",
          overflow: 'auto', 
          visibility: page==0?"visible":"hidden",
          flexDirection: "column",
          display: 'flex',
        }}>
          <Typography style={{
              color:"#FFFFFF",
              fontSize: 18,
              alignSelf:'center',
          }}>
            Wallet Address
          </Typography>
          <Typography style={{
              color:"#FFFFFF",
              fontSize: 18,
              alignSelf:'center',
          }}>
            {walletAddress==""?'-':walletAddress}
          </Typography>
          
          <Typography style={{
              color:"#FFFFFF",
              fontSize: 18,
              alignSelf:'center',
          }}>
            Balance
          </Typography>
          <Typography style={{
              color:"#FFFFFF",
              fontSize: 18,
              alignSelf:'center',
          }}>
            {walletAddress==""?'-':showBalanceDigit(balance, tokenDigit)}
          </Typography>

          <Typography style={{
              color:"#FFFFFF",
              fontSize: 18,
              alignSelf:'center',
          }}>
            NFT Id
          </Typography>
          <Typography style={{
              color:"#FFFFFF",
              fontSize: 18,
              alignSelf:'center',
          }}>
            {nftId}
          </Typography>
          <Typography style={{
              color:"#FFFFFF",
              fontSize: 18,
              alignSelf:'center',
          }}>
            NFT Owner
          </Typography>
          <Typography style={{
              color:"#FFFFFF",
              fontSize: 18,
              alignSelf:'center',
          }}>
            {nftInfo.owner==""?'-':nftInfo.owner}
          </Typography>
          <Typography style={{
              color:"#FFFFFF",
              fontSize: 18,
              alignSelf:'center',
          }}>
            face:{nftInfo.face==""?'-':nftInfo.face} / dice number:{nftInfo.diceNumber==""?'-':nftInfo.diceNumber}
          </Typography>
          <Typography style={{
              color:"#FFFFFF",
              fontSize: 18,
              alignSelf:'center',
          }}>
            Payout Address
          </Typography>
          <Typography style={{
              color:"#FFFFFF",
              fontSize: 18,
              alignSelf:'center',
          }}>
            {nftInfo.payoutAddress==""?'-':nftInfo.payoutAddress}
          </Typography>

          <Box style={{
              width:300,
              flexDirection: "row",
              display: 'flex',
              justifyContent: 'space-between',
              alignSelf:'center',
          }}>
            <GreenButton onClick={()=>{
             
              var params = [
                {
                  from: walletAddress,
                  to: nftContract,
                  gas: '0x30D40', // 200000
                  gasPrice: '0x9502F900', // 2.5gwei
                  value: '0x0',
                  data:
                    '0xa1a2863b'+paddingLeftZero(64, parseInt(nftId).toString(16))+paddingLeftZero(64, "0")+paddingLeftZero(64, "DE0B6B3A7640000"),//1000000000000000000
                },
              ];

              bet(params)
            }}>
                <Typography style={{
                  fontSize: '14px',
                }}>
                  Small
                </Typography>
            </GreenButton> 
            <RedButton onClick={()=>{
              
              var params = [
                {
                  from: walletAddress,
                  to: nftContract,
                  gas: '0x30D40', // 200000
                  gasPrice: '0x9502F900', // 2.5gwei
                  value: '0x0',
                  data:
                    '0xa1a2863b'+paddingLeftZero(64, parseInt(nftId).toString(16))+paddingLeftZero(64, "1")+paddingLeftZero(64, "DE0B6B3A7640000"),//1000000000000000000
                },
              ];

              bet(params)

            }}>
                <Typography style={{
                  fontSize: '14px',
                }}>
                  Big
                </Typography>
            </RedButton> 
          </Box>
          <Box style={{
            height:30,
          }}/>
          <Typography style={{
              color:"#FFFFFF",
              fontSize: 18,
              alignSelf:'center',
          }}>
            TransactionHash
          </Typography>
          <Typography style={{
              color:"#FFFFFF",
              fontSize: 10,
              alignSelf:'center',
          }}>
            {betTxid==""?'-':betTxid}
          </Typography>
          <Typography style={{
              color:"#FFFFFF",
              fontSize: 18,
              alignSelf:'center',
          }}>
            Result
          </Typography>
          <Typography style={{
              color:"#FFFFFF",
              fontSize: 18,
              alignSelf:'center',
          }}>
            {betResult==""?'-':betResult}
          </Typography>

        </Box>

        <Box style={{
          position: 'fixed',
          width:'100%',
          height: height-36,
          background:"linear-gradient(#000000,#052F41)",
          backgroundRepeat:"no-repeat",
          backgroundSize:"100% 100%",
          overflow: 'auto', 
          visibility: page==1?"visible":"hidden",
          flexDirection: "column",
          display: 'flex',
        }}>
          <Typography style={{
              color:"#FFFFFF",
              fontSize: 18,
              alignSelf:'center',
          }}>
            Wallet Address
          </Typography>
          <Typography style={{
              color:"#FFFFFF",
              fontSize: 18,
              alignSelf:'center',
          }}>
            {walletAddress==""?'-':walletAddress}
          </Typography>
          <Typography style={{
              color:"#FFFFFF",
              fontSize: 18,
              alignSelf:'center',
          }}>
            Dice Face
          </Typography>
          <InputContainer style={{
            flexDirection: "row",
            display: 'flex',
            justifyContent:"space-between",
            alignItems:"center",
            alignSelf:'center',
            marginTop:10,
          }}>
            <TextField
                margin="dense"
                type="number"
                inputMode="decimal"
                fullWidth
                style={{
                    marginLeft:16,
                }}
                InputProps={{
                    disableUnderline: true,
                    style: {
                        color: "#ffffff",
                    }
                }}
                placeholder={"0"}
                onChange={event=>{
                    var amount = event.target.value
                    if (isValidAmount(amount)) {
                      setMintFace(amount)
                    }
                }}
            />
          </InputContainer>
          <Typography style={{
              color:"#FFFFFF",
              fontSize: 18,
              alignSelf:'center',
          }}>
            Dice Number
          </Typography>
          <InputContainer style={{
            flexDirection: "row",
            display: 'flex',
            justifyContent:"space-between",
            alignItems:"center",
            alignSelf:'center',
            marginTop:10,
          }}>
            <TextField
                margin="dense"
                type="number"
                inputMode="decimal"
                fullWidth
                style={{
                    marginLeft:16,
                }}
                InputProps={{
                    disableUnderline: true,
                    style: {
                        color: "#ffffff",
                    }
                }}
                placeholder={"0"}
                onChange={event=>{
                    var amount = event.target.value
                    if (isValidAmount(amount)) {
                      setMintDiceNumber(amount)
                    }
                }}
            />
          </InputContainer>
          <Typography style={{
              color:"#FFFFFF",
              fontSize: 18,
              alignSelf:'center',
          }}>
            Payout Address
          </Typography>
          <InputContainer style={{
            flexDirection: "row",
            display: 'flex',
            justifyContent:"space-between",
            alignItems:"center",
            alignSelf:'center',
            marginTop:10,
          }}>
            <TextField
                margin="dense"
                fullWidth
                style={{
                    marginLeft:16,
                }}
                InputProps={{
                    disableUnderline: true,
                    style: {
                        color: "#ffffff",
                    }
                }}
                onChange={event=>{
                    var address = event.target.value
                    if(isValidAddress) {
                      setMintPayoutAddress(address)
                    }
                }}
            />
          </InputContainer>
            <GreenButton style={{
                alignSelf:'center',
                marginTop:20,
            }} onClick={()=>{
              if(mintFace==0 | mintDiceNumber==0 | mintFace*mintDiceNumber%2==1) {
                alert("face*diceNumber Need Even")
                return
              }
              var metaData = ""

              var params = [
                {
                  from: walletAddress,
                  to: nftContract,
                  gas: '0x493E0', // 300000
                  gasPrice: '0x9502F900', // 2.5gwei
                  value: '0x0',
                  data:
                    '0xe7d3fe6b'+
                    paddingLeftZero(64, parseInt(mintFace).toString(16))+
                    paddingLeftZero(64, parseInt(mintDiceNumber).toString(16))+
                    paddingLeftZero(64, mintPayoutAddress.substr(2))
                },
              ];

              mintNft(params)
            }}>
                <Typography style={{
                  fontSize: '14px',
                }}>
                  mint
                </Typography>
            </GreenButton> 
            <Typography style={{
                color:"#FFFFFF",
                fontSize: 18,
                alignSelf:'center',
            }}>
              TransactionHash
            </Typography>
            <Typography style={{
                color:"#FFFFFF",
                fontSize: 10,
                alignSelf:'center',
            }}>
              {mintTxid==""?'-':mintTxid}
            </Typography>
            <Typography style={{
                color:"#FFFFFF",
                fontSize: 18,
                alignSelf:'center',
            }}>
              Result
            </Typography>
            <Typography style={{
                color:"#FFFFFF",
                fontSize: 18,
                alignSelf:'center',
            }}>
              {JSON.stringify(mintResult)}
            </Typography>

            <Box style={{
              position: 'fixed',
              width:'100%',
              height: height-36,
              background:"linear-gradient(#000000,#052F41)",
              backgroundRepeat:"no-repeat",
              backgroundSize:"100% 100%",
              overflow: 'auto', 
              visibility: page==2?"visible":"hidden",
              flexDirection: "column",
              display: 'flex',
            }}>
              <Typography style={{
                  color:"#FFFFFF",
                  fontSize: 18,
                  alignSelf:'center',
              }}>
                Wallet Address
              </Typography>
              <Typography style={{
                  color:"#FFFFFF",
                  fontSize: 18,
                  alignSelf:'center',
              }}>
                {walletAddress==""?'-':walletAddress}
              </Typography>
              <Typography style={{
                  color:"#FFFFFF",
                  fontSize: 18,
                  alignSelf:'center',
              }}>
                Token(ERC20) Address
              </Typography>
              <Typography style={{
                  color:"#FFFFFF",
                  fontSize: 18,
                  alignSelf:'center',
              }}>
                {tokenAddress}
              </Typography>
              <Typography style={{
                  color:"#FFFFFF",
                  fontSize: 18,
                  alignSelf:'center',
              }}>
                Approve NftContract
              </Typography>
              <Typography style={{
                  color:"#FFFFFF",
                  fontSize: 18,
                  alignSelf:'center',
              }}>
                {isApprove.toString()}
              </Typography>
              {
                isApprove==false?
                <GreenButton style={{
                    alignSelf:'center',
                    marginTop:20,
                }} onClick={()=>{
        
                  approveToken()
                }}>
                    <Typography style={{
                      fontSize: '14px',
                    }}>
                      approve
                    </Typography>
                </GreenButton>
                :
                <Box/>
              }
            </Box>
        </Box>

        <Box style={{
          position: 'fixed',
          width:'100%',
          top: 'auto',
          bottom: 0,
          backgroundColor:"#ffffff",
        }}>
          <BottomNavigation
            value={page}
            onChange={(event, newValue) => {
              setPage(newValue);
            }}
            showLabels
            style={{
              background: "#000000F0",
              height: 42
            }}
          >
            <BottomNavigationActionButton label="Dice" showLabel={true} />
            <BottomNavigationActionButton label="Workshop" showLabel={true} />
            <BottomNavigationActionButton label="Token" showLabel={true} />

          </BottomNavigation>
        </Box>

        <LoadingConainer {...{title:"", loadingDialog}}/>

      </Box>
    </React.Fragment>
  );
}

export default Game;
