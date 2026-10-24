import React, { useEffect } from 'react'
import { useSelector, useDispatch } from 'react-redux';
import * as BigNumber from 'bignumber.js';
import date from 'date-and-time'
import QRCode from 'qrcode.react'
import { useTranslation } from 'react-i18next';

import { withStyles } from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogContent from '@material-ui/core/DialogContent';
import Typography from '@material-ui/core/Typography';
import Box from '@material-ui/core/Box';
import CircularProgress from '@material-ui/core/CircularProgress';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import Slide from '@material-ui/core/Slide';
import AppBar from '@material-ui/core/AppBar';
import Card from '@material-ui/core/Card';
import CardActionArea from '@material-ui/core/CardActionArea';
import Divider from '@material-ui/core/Divider';

import CloseIcon from '@material-ui/icons/Close';
import AddIcon from '@material-ui/icons/Add';
import ArrowBackIcon from '@material-ui/icons/ArrowBack';
import SwapHorizIcon from '@material-ui/icons/SwapHoriz';
import FilterListIcon from '@material-ui/icons/FilterList';
import InputIcon from '@material-ui/icons/Input';
import AccountBalanceIcon from '@material-ui/icons/AccountBalance';
import LocalAtmIcon from '@material-ui/icons/LocalAtm';
import LibraryBooksIcon from '@material-ui/icons/LibraryBooks';
import AutorenewIcon from '@material-ui/icons/Autorenew';
import FiberNewIcon from '@material-ui/icons/FiberNew';
import DeleteIcon from '@material-ui/icons/Delete';
import ChatIcon from '@material-ui/icons/Chat';

import IconButton from '@material-ui/core/IconButton';
import Toolbar from '@material-ui/core/Toolbar';
import TextField from '@material-ui/core/TextField';
import ListAltIcon from '@material-ui/icons/ListAlt';
import InputLabel from '@material-ui/core/InputLabel';
import Select from '@material-ui/core/Select';
import FormControl from '@material-ui/core/FormControl';

import { postPoaEmailRecoverySend, postPoaEmailRecovery } from './api/getMainnetInfuraApi'
import { paddingLeftZero, showBalanceDigit, toBalanceDigit } from './Utils'
import { isValidEmail, isValidEmailCode } from './RuleUtils'


var height = document.documentElement.clientHeight
var pageWidth  = window.screen.width > 400 ? 400 : window.screen.width

const ConfirmButton = withStyles(theme => ({
    root: {
      height: 32,
      width: 120,
      backgroundColor: "rgb(2, 192, 118)",
      '&:hover': {
        backgroundColor: "rgb(2, 192, 118)",
      },
      borderRadius:4,
    },
    label: {
        color: '#ffffff',
        fontSize: '10px',
    },
  }))(Button);

const CancelButton = withStyles(theme => ({
    root: {
        height: 32,
        width: 120,
        backgroundColor: "#f84960",
        '&:hover': {
        backgroundColor: "#f84960",
        },
        borderRadius:4,
    },
    label: {
        color: '#ffffff',
        fontSize: '10px',
    },
}))(Button);

const Transition = React.forwardRef(function Transition(props, ref) {
    return <Slide direction="up" ref={ref} {...props} />;
  });

export const EmailRecoveryContainer = ({emailRecoveryDialog, setEmailRecoveryDialog, requestAppGeneratePky, requestAppSignData, setWalletAddress, requestAppSetPoaAddress, setOpenWalletInitDialog})=>{
    
  const { t, i18n } = useTranslation();

  const [step, setStep] = React.useState("0")
  const [email, setEmail] = React.useState("")
  const [authorizedAddress, setAuthorizedAddress] = React.useState("")
  const [loadingDialog, setLoadingDialog] = React.useState(false)
  
  useEffect(
    () => {

      return () => {
        setStep("0")
        setEmail("")
      }
    }, [emailRecoveryDialog]
  )

  return  <Box>
    <Dialog fullScreen style={{
      height:height, // height*0.96,
      position: 'fixed',
      width:'100%',
      top: 'auto',
      bottom: 0,
    }} open={emailRecoveryDialog} onClose={()=>{
      setEmailRecoveryDialog(false)
    }} TransitionComponent={Transition}>
      <Box style={{
        padding: 8,
      }}>
        <AppBar elevation={0} style={{
          position: 'fixed',
          width:'100%',
          height: 20,
          background: "linear-gradient(#2295ff,#77c1fe)",
        }}></AppBar>
        <Toolbar style={{
          top: 11
        }}>
          {
            step!="2"?
              <Box style={{
                flexDirection: "row",
                display: 'flex',
                width:"100%",
                alignItems:'center',
              }}>
                <IconButton edge="start" color="inherit" onClick={()=>{
                    setEmailRecoveryDialog(false)
                }} aria-label="close">
                    <ArrowBackIcon />
                </IconButton>
                <Typography>
                  {t('restore_the_wallet')}
                </Typography>
              </Box>
            :
              <Box/>
          }
          
        </Toolbar>
      </Box>
      <Box style={{
        padding: 8,
        height: '100%',
        backgroundColor: '#c7e5ff'
      }}>
      {
        step=="0"?
        <EmailInputContainer {...{setStep, setLoadingDialog, email, setEmail, setAuthorizedAddress, requestAppGeneratePky, requestAppSignData}}/>
        :step=="1"?
        <EmailValidateContainer {...{setStep, setLoadingDialog, authorizedAddress, email, requestAppSignData, setWalletAddress, requestAppSetPoaAddress}}/>
        :step=="2"?
        <EmailRecoverySuccessContainer {...{setEmailRecoveryDialog, setOpenWalletInitDialog}}/>
        :
        <Box/>
      }
      </Box>
    </Dialog>
    <LoadingConainer {...{title:t('processing'), loadingDialog}}/>
  </Box>
}

const EmailInputContainer = ({setStep, setLoadingDialog, email, setEmail, setAuthorizedAddress, requestAppGeneratePky, requestAppSignData}) => {
    
  const { t, i18n } = useTranslation();

  const bindEmail = ()=>{
    console.log("email "+email)

    if (!isValidEmail(email)) {
      return
    }

    var emailHex = Buffer(email, 'ascii').toString('hex');

    var data = "0x00000000"+emailHex

    console.log("data "+data)

    setLoadingDialog(true)
    requestAppGeneratePky((authorizedAddress) => { 
      console.log('authorizedAddress '+authorizedAddress)
      setAuthorizedAddress(authorizedAddress)

      var params = {
        authorizedAddress:authorizedAddress,
        data:data
      }

      postPoaEmailRecoverySend(JSON.stringify(params)).then(result=>{
          setLoadingDialog(false)
          if (result.code != 0) {
              return
          }
          console.log('postPoaEmailRecoverySend '+result.data)

          setStep(1)
      }).catch(err=>{
          setLoadingDialog(false)
          console.error(err)
      })

    })

    
  }

  return  <Box style={{
    width:'auto', //pageWidth-32,
    flexDirection: "column",
    display: 'flex',
    marginLeft:16,
    marginRight:16,
    backgroundColor: "#c7e5ff"
  }}>
      {/* <Typography style={{
          marginTop:12,
      }}>
          Your Email
      </Typography> */}
      <Typography style={{
          marginTop:12,
      }}>
        {t('description_email_verification')}
      </Typography>
      <TextField
          margin="dense"
          type="email"
          placeholder={t('please_enter_your_email')}
          variant="outlined"
          fullWidth
          onChange={event=>{
              var text = event.target.value
              setEmail(text)
          }}
      />
      <ConfirmButton style={{
        alignSelf:'center',
        marginTop:12,
        background: "linear-gradient(#2295ff,#77c1fe)"
      }} onClick={()=>{
        bindEmail()
      }}>
        {t('next')}
      </ConfirmButton>
      
  </Box>
}

const EmailValidateContainer = ({setStep, setLoadingDialog, authorizedAddress, email, requestAppSignData, setWalletAddress, requestAppSetPoaAddress}) => {

  const { t, i18n } = useTranslation();

  const [validateCode, setValidateCode] = React.useState("")
    
  const validateEmail = ()=>{
    console.log("validateCode "+validateCode)

    if (!isValidEmailCode(validateCode)) {
      return
    }

    var emailHex = Buffer(email, 'ascii').toString('hex');

    var data = "0x00000000"+validateCode+emailHex

    console.log("data "+data)

    var params = {
      authorizedAddress:authorizedAddress,
      data:data
    }

    setLoadingDialog(true)
    postPoaEmailRecovery(JSON.stringify(params)).then(result=>{
        setLoadingDialog(false)
        if (result.code != 0) {
            return
        }
        console.log('postPoaEmailRecovery '+result.data)
        setWalletAddress(result.data)
        requestAppSetPoaAddress(result.data)
        setStep(2)
    }).catch(err=>{
        console.error(err)
        setLoadingDialog(false)
    })
  }

  return  <Box style={{
    width: 'auto',// pageWidth-32,
    flexDirection: "column",
    display: 'flex',
    marginLeft:16,
    marginRight:16,
  }}>
      <Typography>
        {t('email_verification')}
      </Typography>
      <Typography style={{
          marginTop:12,
      }}>
        {t('description_verification_code')}
      </Typography>
      <TextField
          margin="dense"
          type="number"
          variant="outlined"
          fullWidth
          onChange={event=>{
              var text = event.target.value
              setValidateCode(text)
          }}
      />
      <ConfirmButton style={{
        alignSelf:'center',
        marginTop:12,
        background: "linear-gradient(#2295ff,#77c1fe)"
      }} onClick={()=>{
        validateEmail()
      }}>
        {t('confirm')}
      </ConfirmButton>
  </Box>
}

const EmailRecoverySuccessContainer = ({setEmailRecoveryDialog, setOpenWalletInitDialog}) => {
    
  const { t, i18n } = useTranslation();

  return  <Box style={{
    width:'auto',
    flexDirection: "column",
    display: 'flex',
    marginLeft:16,
    marginRight:16,
  }}>
      <Typography style={{
        alignSelf:'center',
      }}>
        {t('wallet_restore_success')}
      </Typography>
      
      <ConfirmButton style={{
        alignSelf:'center',
        marginTop:12,
        background: "linear-gradient(#2295ff,#77c1fe)"
      }} onClick={()=>{
        setEmailRecoveryDialog(false)
        setOpenWalletInitDialog(false)
      }}>
        {t('get_started')}
      </ConfirmButton>
  </Box>
}

export const LoadingConainer = ({title, loadingDialog}) => {
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