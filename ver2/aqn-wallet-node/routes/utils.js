const crypto = require('crypto');
const {NODE} = require('./config')

const iv = NODE.aesIV
const key = NODE.aesKey

const getCosignerAddress = () => {
  
  return NODE.coSignerAddress
};

const getCosignerPky = (coSignerAddress) => {

  return NODE.coSignerPky
};

const sha512=(text) => {
  var sha512 = crypto.createHash('sha512');
  sha512.update(text)
  var hash = sha512.digest('hex')

  return hash;
}

const paddingLeftZero = (digits, text) => {
  if (text.length>digits) {
      return text
  }

  const paddingDigits = digits - text.length

  var result = ""
  var i
  for (i = 0; i<paddingDigits; i++) {
      result += "0"
  }
  result = result + text

  return result
};

const sha256=(text) => {
  var sha256 = crypto.createHash('sha256');
  sha256.update(text)
  var hash = sha256.digest('hex')

  return hash;
}

const aesEncrypt = (text) => {
  const algorithm = 'aes128';

  const cipher = crypto.createCipheriv(algorithm, key, iv);  
  const encrypted = cipher.update(text, 'utf8', 'hex') + cipher.final('hex');
  return encrypted
}

const aesDecrypt = (text) => {
  const algorithm = 'aes128';

  const decipher = crypto.createDecipheriv(algorithm, key, iv);
  let str = decipher.update(text, 'hex', 'utf8');
  str += decipher.final('utf8');

  return str
}

module.exports = {getCosignerAddress, getCosignerPky, sha512
  , sha256, aesEncrypt, aesDecrypt, paddingLeftZero}