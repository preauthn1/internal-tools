const createResponse = (success, data) => {
    if(success) {
        return {
            "code":0,
            "data":data,
            "message":"success"
        }
    } else {
        return {
            "code":999,
            "data":data,
            "message":"success"
        }
    }
};

module.exports = {createResponse}