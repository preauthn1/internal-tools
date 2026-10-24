import * as BigNumber from 'bignumber.js';

export const showBalanceDigit = (balance, digit) => {
    if (balance == 0) {
        return '0.00'
    }
    
    var uint = new BigNumber(10).pow(digit)

    
    if (balance%(uint/100) == 0) {
        return new BigNumber(balance).div(uint).toFixed(2)
    }

    if (balance%(uint/1000) == 0) {
        return new BigNumber(balance).div(uint).toFixed(3)
    }

    return new BigNumber(balance).div(uint).toFixed(6)
};


export const paddingLeftZero = (digits, text) => {
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