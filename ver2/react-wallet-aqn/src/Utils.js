import * as BigNumber from 'bignumber.js';
import crypto from 'crypto'

export const azureFunctionCode = "";

export const showBalance = (balance) => {
    if (balance == 0) {
        return '0.00'
    }
    
    if (balance%10000000000000000 == 0) {
        return new BigNumber(balance).div(1000000000000000000).toFixed(2)
    }

    var reminder = balance%10000000000
    return new BigNumber(balance-reminder).div(1000000000000000000).toFixed()
};

export const showBalanceBtc = (balance) => {
    if (balance == 0) {
        return '0.00'
    }
    
    if (balance%1000000 == 0) {
        return new BigNumber(balance).div(100000000).toFixed(2)
    }

    return new BigNumber(balance).div(100000000).toFixed()
};

export const showBalanceDigit = (balance, digit) => {
    if (balance == 0) {
        return '0.00'
    }

    // var uint = 1
    // for (var i=0; i<digit; i++) {
    //     uint = uint*10
    // }
    var uint = new BigNumber(10).pow(digit)

    
    if (balance%(uint/100) == 0) {
        return new BigNumber(balance).div(uint).toFixed(2)
    }

    if (balance%(uint/1000) == 0) {
        return new BigNumber(balance).div(uint).toFixed(3)
    }

    return new BigNumber(balance).div(uint).toFixed(6)

    // if (digit <= 10) {
    //     return new BigNumber(balance).div(uint).toFixed(digit)
    // } else {
    //     var reminder = balance%10000000000
    //     return new BigNumber(balance-reminder).div(uint).toFixed()
    // }
};

export const showSwapBalanceDigit = (balance, digit, outTokenDigit) => {
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

    if (outTokenDigit <= 10) {
        return new BigNumber(balance).div(uint).toFixed(outTokenDigit)
    } else {
        var reminder = new BigNumber(balance).mod(new BigNumber(10000000000))
        return (new BigNumber(balance).minus(reminder)).div(uint).toFixed(10)
    }

    // return new BigNumber(balance).div(uint).toFixed(outTokenDigit)
};

export const toBalance = (balance) => {
    return new BigNumber(balance).times(1000000000000000000).toFixed()
};

export const toBalanceBtc = (balance) => {
    return new BigNumber(balance).times(100000000).toFixed()
};

export const toBalanceDigit = (balance, digit) => {
    var uint = 1
    for (var i=0; i<digit; i++) {
        uint = uint*10
    }
    return new BigNumber(balance).times(uint).toFixed()
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

export const hashPassword = [REDACTED] => {

    return crypto.createHash('sha256').update(input).digest('hex')
};

// export const isValidPassword = [REDACTED] => {
    
//     return hashPassword(input) === localStorage.getItem('password')
// };

const paddingKey = (key) => {
    if (!key || key === '') {
        return ''
    }
    
    var interval = 16 - Buffer.from(key).length

    while (interval != 0) {
        if (interval > 0) {
            key += key
        } else {
            key = Buffer.from(key).slice(0, 16).toString()
        }
        interval = 16  - Buffer.from(key).length
    }

    return key
}

export const encryptMsg = (input, key) => {
    return input
    // if (key === '') {
    //     return ''
    // }

    // if (!window.encryptMsgWasm) {
    //     return ''
    // }

    // return window.encryptMsgWasm(input, paddingKey(key))
};

export const encryptPky = (input) => {
    return input
    // if (!window.encryptPkyWasm) {
    //     return ''
    // }

    // return window.encryptPkyWasm(input)
};

export const decryptMsg = (input, key) => {
    return input
    // if (key === '') {
    //     return ''
    // }

    // if (!window.decryptMsgWasm) {
    //     return ''
    // }

    // return window.decryptMsgWasm(input, paddingKey(key))
};

export const tsToDateString = (ts) => {
    var date = new Date(ts*1000)

    var formatted_string = date.getFullYear() + "-";

    if (date.getMonth() < 9) {
    formatted_string += "0";
    }
    formatted_string += (date.getMonth() + 1);
    formatted_string += "-";

    if(date.getDate() < 10) {
    formatted_string += "0";
    }
    formatted_string += date.getDate();
    formatted_string += " ";

    if(date.getHours() < 10) {
        formatted_string += "0";
    }
    formatted_string += date.getHours();
    formatted_string += ":";

    if(date.getMinutes() < 10) {
        formatted_string += "0";
    }
    formatted_string += date.getMinutes();
    formatted_string += ":";

    if(date.getSeconds() < 10) {
        formatted_string += "0";
    }
    formatted_string += date.getSeconds();

    return formatted_string
}

export const getNowTs = () => {
    return (new Date().getTime()/1000).toFixed(0)
}

export const showAmountDigit = (balance, digit) => {
    // if (balance == 0) {
    //     return '0.00'
    // }

    // var uint = 1
    // for (var i=0; i<digit; i++) {
    //     uint = uint*10
    // }
    
    // if (balance%(uint/1000000) == 0) {
    //     return new BigNumber(balance).div(uint).toFixed(6)
    // }

    // if (digit <= 10) {
    //     return new BigNumber(balance).div(uint).toFixed(2)
    // } else {
    //     var reminder = balance%10000000000
    //     return new BigNumber(balance-reminder).div(uint).toFixed()
    // }

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