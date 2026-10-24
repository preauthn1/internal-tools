export const isValidAddress = (address) => {
    var pattern = new RegExp('^0x[a-fA-F0-9]{40}$')
    return address.match(pattern)
};

export const isValidFilAddress = (address) => {
    var pattern = new RegExp('^f1[a-z0-9]{39}$')
    return address.match(pattern)
};

export const isValidBtcAddress = (address) => {
    var pattern = new RegExp('^[13][1-9A-HJ-NP-Za-km-z]{26,33}$')
    return address.match(pattern)
};

export const isValidPrivateKey = (pky) => {
    var pattern = new RegExp('^[a-fA-F0-9]{64}$')
    return pky.match(pattern)
};

export const isValidPassword = [REDACTED] => {
    var pattern = new RegExp('^(?=.*[A-Za-z])(?=.*\\d)[A-za-z\\d]{8,20}$')
    return password.match(pattern)
};

export const isValidMnemonic = (mnemonic) => {
    var pattern = new RegExp('[a-z]{1,}\\s[a-z]{1,}\\s[a-z]{1,}\\s[a-z]{1,}\\s[a-z]{1,}\\s[a-z]{1,}\\s[a-z]{1,}\\s[a-z]{1,}\\s[a-z]{1,}\\s[a-z]{1,}\\s[a-z]{1,}\\s[a-z]{1,}$')
    return mnemonic.match(pattern)
};

export const isValidAmount = (amount) => {
    var pattern = new RegExp('^[0-9]{0,}(\.[0-9]{1,6})?$')
    return amount.match(pattern)
};

export const isValidReferenceCode = (code) => {
    return  0<code.length && code.length <= 10
};

export const isValidTxid = (txid) => {
    var pattern = new RegExp('^0x[a-fA-F0-9]{64}$')
    return txid.match(pattern)
};

export const isValidEmail = (email) => {
    var pattern = new RegExp('^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$')
    return email.match(pattern)
};

export const isValidEmailCode = (email) => {
    var pattern = new RegExp('^[0-9]{6}$')
    return email.match(pattern)
};

export const isValidNftId = (nftId) => {
    var pattern = new RegExp('^[0-9]{0,6}$')
    return nftId.match(pattern)
};