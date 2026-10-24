const paddingKey = (key) => {
    // if (!key || key === '') {
    //     return ''
    // }
    
    // var interval = 16 - Buffer.from(key).length

    // while (interval != 0) {
    //     if (interval > 0) {
    //         key += key
    //     } else {
    //         key = Buffer.from(key).slice(0, 16).toString()
    //     }
    //     interval = 16  - Buffer.from(key).length
    // }

    return key
}

export const encryptMsg = (input, key) => {
    // return input
    if (key === '') {
        return input
    }

    if (!window.encryptMsgWasm) {
        return input
    }

    return window.encryptMsgWasm(input, paddingKey(key))
};

export const decryptMsg = (input, key) => {
    // return input
    if (key === '') {
        return ''
    }

    if (!window.decryptMsgWasm) {
        return ''
    }

    return window.decryptMsgWasm(input, paddingKey(key))
};