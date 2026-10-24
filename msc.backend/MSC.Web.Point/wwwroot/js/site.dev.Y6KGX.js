
var initialHelper =
{
    initializeInputPage: function(requestID)
    {
        point_input.initialize(requestID);

        $("#submit").click(function () {

            point_input.submit();

        });
    }
}


var point_input =
{
    encryptKey: '',
    isLockBtn: false,
    errorTimes: 0,
    errorCardNoRecords: [],
    myLang: '',

    initialize: function (my_encryptKey) {

        this.encryptKey = my_encryptKey
        this.bindBtnAvailableListen();
    },

    setLang: function (lang) {

        this.myLang = lang
    },

    bindBtnAvailableListen: function() {
        
        $("input[type=text]").keyup(function () {

            
            var value1 = $("#card_1").val() + $("#card_2").val() + $("#card_3").val() + $("#card_4").val();
            var value2 = $("#cardPassword").val();

            var is_valid = value1.length >= 16 && value2.length >= 16;

            if (is_valid) {
                $("#submit").addClass("available");
            }
            else {
                $("#submit").removeClass("available");
            }
        });

    },


    submit: function () {

        var $this = this;
        
        if ($this.isLockBtn) {
            return;
        }

        var msg_cardNo = "#msg_cardNo";
        var msg_cardPassword = [REDACTED]";


        $this.resetErrorHint(msg_cardNo, msg_cardPassword );

               
        var cardNoValue = $("#card_1").val() + $("#card_2").val() + $("#card_3").val() + $("#card_4").val();        
        var cardPassword = $("#cardPassword").val();

        cardNoValue = cardNoValue ? cardNoValue.toUpperCase() : '';
        cardPassword = [REDACTED] ? cardPassword.toUpperCase() : '';


        if (cardNoValue.length < 16) {

            let text = $this.get_i18n("error_text_length");
            $this.setErrorHint(msg_cardNo , text);
            return;
        }

        if (cardPassword.length < 16) {

            let text = $this.get_i18n("error_password_text_length");
            $this.setErrorHint(msg_cardPassword , text);
            return;
        }

        //大寫檢查
        if (!this.isUpperCase(cardNoValue)) {
            let text = $this.get_i18n("error_text_lower_case");
            $this.setErrorHint(msg_cardNo ,text);
            return;
        }

        //大寫檢查
        if (!this.isUpperCase(cardPassword)) {
            let text = $this.get_i18n("error_password_text_lower_case");
            $this.setErrorHint(msg_cardPassword ,text);
            return;
        }
        
        //數字+英文檢查
        if (!this.is_onlyLettersAndNumber(cardNoValue)) {
            let text = $this.get_i18n("error_text_invalid_letters");
            $this.setErrorHint(msg_cardNo ,text);
            return;
        }

        //數字+英文檢查
        if (!this.is_onlyLettersAndNumber(cardPassword)) {
            let text = $this.get_i18n("error_password_text_invalid_letters");
            $this.setErrorHint(msg_cardPassword ,text);
            return;
        }

        //非法字母
        if (cardNoValue.includes("1") || cardNoValue.includes("0") ||
            cardNoValue.includes("I") || cardNoValue.includes("O")) {
            let text = $this.get_i18n("error_card_code_not_match");
            $this.setErrorHint(msg_cardNo ,text);
            return;
        }

        //若輸入相同的錯誤卡號，直接拋錯誤
        if ($this.errorCardNoRecords.includes(cardNoValue)) {
            let text = $this.get_i18n("error_card_code_not_match");

            console.log(cardNoValue);
            $this.setErrorHint(msg_cardNo , text);
            return;
        }

        var previous_15 = cardNoValue.substring(0, 15);
        var lastCode = cardNoValue.substring(15, 16);
        var expect_lastWord = this.getLastWord(previous_15);

        //最後一碼檢查
        if (expect_lastWord != lastCode) {
            let text = $this.get_i18n("error_card_code_not_match");
            $this.popupErrorMsg(text);
            return;
        }

        let ts = + new Date();

        var signKey = (cardNoValue + ts.toString() + this.encryptKey);

        var MD5 = CryptoJS.MD5(signKey).toString();

        var postBody =
        {
            MerchantToken: $("#Token").val(),
            CardNo: cardNoValue,
            cardPassword: [REDACTED],
            SignKey: MD5,
            TS: ts
        };

        var postJson = JSON.stringify(postBody);
        
        $this.lockBtn();


        $.ajax({
            type: "POST",
            url: window.API_DOMAIN + '/card/input',
            dataType: 'json',
            data: postJson,
            contentType: 'application/json; charset=utf-8',
            headers:
            {
                'x-device-id': $("#Device").val(),
                'locale': this.myLang
            },
            success: function (res) {

                console.log(res);

                let text = $this.get_i18n("success_msg_card_input")
                $this.popupErrorMsg(text);
                $this.unlockBtn();


                setTimeout(x => {

                    location.href = res.data.returnUrl;

                }, 4000);


            },
            error: function (xhr, status, error) {

                console.log('error');

                $this.errorTimes = $this.errorTimes + 1;

                let code = xhr.status;               
                
                console.log($this.errorTimes);


                if ($this.errorTimes < 10)
                {
                    if (code == 400) {
                        var parse = JSON.parse(xhr.responseText);
                        $this.popupErrorMsg(parse.message);
                        $this.addErrorCardNo(parse, cardNoValue);
                    }
                    
                    $this.unlockBtn();
                }
                else
                {
                    //錯誤太多次 故意延遲回應
                    setTimeout(function () {

                        if (code == 400) {
                            var parse = JSON.parse(xhr.responseText);
                            $this.popupErrorMsg(parse.message);
                            $this.addErrorCardNo(parse, cardNoValue);
                        }
                        $this.unlockBtn();

                    }, 3000);                    
                }

            }
        });

    },
    
    get_i18n: function (key) {

        return $t[key];
    },

    addErrorCardNo: function (responseJson, cardNo) {
        
        //如果是這兩種錯誤碼，前端會記錄起來 重複輸入直接拋錯誤 不打API
        if (responseJson.code == 550001 || 
            responseJson.code == 404056)
        {
            this.errorCardNoRecords.push(cardNo);
        }
    },


    lockBtn: function () {

        this.isLockBtn = true;
        $("#submit").attr("disabled", "disabled");
        $("#load_img").show();        

    },

    unlockBtn: function () {

        this.isLockBtn = false;
        $("#submit").removeAttr("disabled", "disabled");
        $("#load_img").hide();       

    },


    getLastWord: function (myText) {

        var numberOfValueTypes = 32;

        var dic =
        {
            "Z": 1,
            "2": 2,
            "3": 3,
            "4": 4,
            "5": 5,
            "6": 6,
            "7": 7,
            "8": 8,
            "9": 9,
            "Y": 10,
            "A": 11,
            "B": 12,
            "C": 13,
            "D": 14,
            "E": 15,
            "F": 16,
            "G": 17,
            "H": 18,
            "X": 19,
            "J": 20,
            "K": 21,
            "L": 22,
            "M": 23,
            "N": 24,
            "P": 25,
            "Q": 26,
            "R": 27,
            "S": 28,
            "T": 29,
            "U": 30,
            "V": 31,
            "W": 32,
        };

        var dicOfMultiply =
        {
            "1": 15,
            "2": 14,
            "3": 13,
            "4": 12,
            "5": 11,
            "6": 10,
            "7": 9,
            "8": 8,
            "9": 7,
            "10": 6,
            "11": 5,
            "12": 4,
            "13": 3,
            "14": 2,
            "15": 1,
        };

        var arrayOfAllLetters = myText.split('');
        var arrayOfAllMappingNumbers = [];

        for (var i = 0; i < 15; i++) {

            var this_letter = arrayOfAllLetters[i];
            var numberValue = dic[this_letter];
            arrayOfAllMappingNumbers.push(numberValue);
        }

        var totalSum = 0;

        for (var i = 0; i < 15; i++) {

            var mappingNumber = arrayOfAllMappingNumbers[i];
            var multiplyValue = dicOfMultiply[(i + 1).toString()];
            totalSum = totalSum + (mappingNumber * multiplyValue);
        }

        var modValue = totalSum % numberOfValueTypes;
        modValue = modValue == 0 ? numberOfValueTypes : modValue;

        var finalText = "";


        for (var key in dic) {
            var value = dic[key];

            if (value == modValue) {
                finalText = key;
            }
        }

        return finalText;

    },
    
    popupErrorMsg: function (msg) {

        $("#errorMsg").text(msg);
        $("#myModal").modal("show");
    },

    setErrorHint: function (element_id , msg) {

        $(element_id + " span").text(msg);
        $(element_id).removeClass("hide");
    },

    resetErrorHint: function (element_1, element_2) {

        $(element_1 + " span").text('');
        $(element_2 + " span").text('');
        $(element_1).addClass("hide");
        $(element_2).addClass("hide");
    },


    isUpperCase: function (str) {
        return str === str.toUpperCase();
    },

    is_onlyLettersAndNumber: function (str) {
        return str.match("^[A-Za-z0-9]+$");
    }


}