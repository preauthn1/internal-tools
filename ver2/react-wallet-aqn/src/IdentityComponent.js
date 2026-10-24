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

import { postPoaEmailSend, postPoaWalletStart } from './api/getMainnetInfuraApi'
import { paddingLeftZero, showBalanceDigit, toBalanceDigit } from './Utils'
import { isValidEmail, isValidEmailCode } from './RuleUtils'

import Paper from '@material-ui/core/Paper';

var height = document.documentElement.clientHeight
var pageWidth  = window.screen.width > 400 ? 400 : window.screen.width

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

const ConfirmButton = withStyles(theme => ({
  root: {
    height: 44,
    width: 180,
    backgroundColor: "#FFB41D",
    '&:hover': {
      backgroundColor: "#FFB41D",
    },
    borderRadius:18,
  },
  label: {
      color: '#ffffff',
      fontSize: 16,
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

export const EmailBindContainer = ({setEmailInputDialog, requestAppGeneratePky, requestAppSetPoaAddress})=>{

  const [step, setStep] = React.useState("0")
  const [authorizedAddress, setAuthorizedAddress] = React.useState("")
  const [email, setEmail] = React.useState("")
  const [validateCode, setValidateCode] = React.useState("")

  const [loadingDialog, setLoadingDialog] = React.useState(false);
  
  const bindEmail = ()=>{
    console.log("email "+email)

    if (!isValidEmail(email)) {
      return
    }

    var emailHex = Buffer(email, 'ascii').toString('hex');

    setLoadingDialog(true)

    requestAppGeneratePky((responseData) => { 
      console.log('responseData '+responseData)
      setAuthorizedAddress(responseData)
      postPoaEmailSend({
        authorizedAddress:responseData,
        emailHex:emailHex
      }).then(result=>{
        setLoadingDialog(false)

        if (result.code != 0) {
            return
        }

        setStep(1)
      }).catch(err=>{
        setLoadingDialog(false)
        console.error(err)
      })
    })


  }

      
  const validateEmail = ()=>{
    console.log("validateCode "+validateCode)

    if (!isValidEmailCode(validateCode)) {
      return
    }

    var emailHex = Buffer(email, 'ascii').toString('hex');

    setLoadingDialog(true)
    postPoaWalletStart({
      authorizedAddress:authorizedAddress,
      emailHex:emailHex,
      code:validateCode
    }).then(result=>{
      setLoadingDialog(false)
      if (result.code != 0) {
          return
      }
      console.log('postPoaEmailValidate '+result.data)
      requestAppSetPoaAddress(result.data, (status)=>{
        setStep(2)
      })
      // setEmailBindStatus(1)

    }).catch(err=>{
      setLoadingDialog(false)
      console.error(err)
    })
  }

  return <Box>
    {
      step=="0"?
      <EmailInputContainer {...{setEmail, bindEmail}}/>
      :step=="1"?
      <EmailValidateContainer {...{setValidateCode, validateEmail}}/>
      :step=="2"?
      <EmailBindSuccessContainer {...{setEmailInputDialog}}/>
      :
      <Box/>
    }
    <LoadingConainer {...{title:"", loadingDialog}}/>
  </Box>
}
const EmailInputContainer = ({setEmail, bindEmail})=>{

  const { t, i18n } = useTranslation();

  return <Box style={{
    width: 340,
    flexDirection: "column",
    display: 'flex',
    alignItems:"center"
  }}>
      <Box  style={{
        height:30
      }}>

      </Box>
      <Typography style={{
        color:"#FFB41D",
        fontSize:22,
      }}>
        {t('email_verification')}
      </Typography>
      <Typography style={{
        width:280,
        color:"#ffffff",
        fontSize:16,
        marginTop:8,
        textAlign:'center',
      }}>
        {t('description_email_verification')}
      </Typography>
      <Typography style={{
        width:280,
        color:"#ffffff",
        fontSize:16,
        marginTop:8,
        textAlign:'center',
      }}>
        {t('description_email_verification_2')}
      </Typography>
      <Divider  style={{ 
          background: '#ffffff', 
          width:282 ,
          marginTop:36,
          marginBottom:36,
      }}/>
      <Typography style={{
        width:280,
        color:"#ffffff",
        fontSize:14,
      }}>
        {t('please_enter_your_email')}
      </Typography>
      <InputContainer  style={{
          flexDirection: "row",
          display: 'flex',
          justifyContent:"space-between",
          alignItems:"center",
          marginTop:4,
      }}>
        <TextField
            margin="dense"
            type="email"
            fullWidth
            style={{
              marginLeft:16,
              marginRight:16,
            }}
            InputProps={{
              disableUnderline: true,
              style: {
                  color: "#ffffff",
              }
            }}
            onChange={event=>{
                var text = event.target.value
                setEmail(text)
            }}
        />
      </InputContainer>
      <ConfirmButton style={{
        alignSelf:'center',
        marginTop:38,
      }} onClick={()=>{
        bindEmail()
      }}>
        {t('next')}
      </ConfirmButton>
      <Box style={{
        height:36
      }}/>
  </Box>
}

const EmailValidateContainer = ({setValidateCode, validateEmail})=>{

  const { t, i18n } = useTranslation();

  return <Box style={{
    width: 340,
    flexDirection: "column",
    display: 'flex',
    alignItems:"center"
  }}>
      <Box  style={{
        height:30
      }}>

      </Box>
      <Typography style={{
        color:"#FFB41D",
        fontSize:22,
      }}>
        {t('email_verification')}
      </Typography>
      <Typography style={{
        width:280,
        color:"#ffffff",
        fontSize:16,
        marginTop:8,
        textAlign:'center',
      }}>
        {t('description_verification_code')}
      </Typography>
      <Divider  style={{ 
          background: '#ffffff', 
          width:282 ,
          marginTop:36,
          marginBottom:36,
      }}/>
      <Typography style={{
        width:280,
        color:"#ffffff",
        fontSize:14,
      }}>
        {t('please_enter_verification_code')}
      </Typography>
      <InputContainer  style={{
          flexDirection: "row",
          display: 'flex',
          justifyContent:"space-between",
          alignItems:"center",
          marginTop:4,
      }}>
        <TextField
            margin="dense"
            type="number"
            fullWidth
            style={{
              marginLeft:16,
              marginRight:16,
            }}
            InputProps={{
              disableUnderline: true,
              style: {
                  color: "#ffffff",
              }
            }}
            onChange={event=>{
                var text = event.target.value
                setValidateCode(text)
            }}
        />
      </InputContainer>
      <ConfirmButton style={{
        alignSelf:'center',
        marginTop:38,
      }} onClick={()=>{
        validateEmail()
      }}>
        {t('next')}
      </ConfirmButton>
      <Box style={{
        height:36
      }}/>
  </Box>
}
const EmailBindSuccessContainer = ({setEmailInputDialog}) => {
    
  const { t, i18n } = useTranslation();

  return <Box style={{
    width: 340,
    flexDirection: "column",
    display: 'flex',
    alignItems:"center"
  }}>
      <Box  style={{
        height:30
      }}>

      </Box>
      <Typography style={{
        color:"#ffffff",
        fontSize:20,
      }}>
        {t('login_succeeded')}
      </Typography>
      <ConfirmButton style={{
        alignSelf:'center',
        marginTop:38,
      }} onClick={()=>{
        setEmailInputDialog(false)
      }}>
        {t('ok')}
      </ConfirmButton>
      <Box style={{
        height:36
      }}/>
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