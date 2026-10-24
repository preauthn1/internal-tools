

function callVueFunctionByName(functionName , data) {
    app[functionName](data);
}

var date_picker =
{

    //請注意 datetimepicker系列的參數 config 包含哪些prop要傳入

    //包含 小時、分鐘的選擇器 預設空值
    datetime_range_empty: function (config) {

        let vue_instance = config.vue_instance;

        if (config.format === undefined || config.format === null) {
            config["format"] = "YYYY-MM-DD HH:mm";
        }

        $(config.element_ID).daterangepicker({
            autoUpdateInput: false,
            timePicker: true,
            locale: {
                format: config.format
            }
        }, function (start, end, label) {

            var myTimzoneOffset = new Date().getTimezoneOffset();

            //傳給API 需要強制轉成utc 0
            let start_utc_0 = moment(start).add(myTimzoneOffset, 'minutes').format(config.format);
            let end_utc_0 = moment(end).add(myTimzoneOffset, 'minutes').format(config.format);

            vue_instance[config.vue_prop_text] = start.format(config.format) + " - " + end.format(config.format);
            vue_instance[config.vue_prop_begin] = start_utc_0;
            vue_instance[config.vue_prop_end] = end_utc_0;

        });

    },

    //包含 小時、分鐘的選擇器 預設今日-明日
    datetime_range_today: function (config) {

        let vue_instance = config.vue_instance;

        if (config.format === undefined || config.format === null) {
            config["format"] = "YYYY-MM-DD HH:mm";
        }

        $('#createtime').daterangepicker({
            timePicker: true,
            startDate: moment().startOf('day'),
            endDate: moment().startOf('day').add(24, 'hour'),
            locale: {
                format: config.format
            }
        }, function (start, end, label) {

            var myTimzoneOffset = new Date().getTimezoneOffset();

            //傳給API 需要強制轉成utc 0
            let start_utc_0 = moment(start).add(myTimzoneOffset, 'minutes').format(config.format);
            let end_utc_0 = moment(end).add(myTimzoneOffset, 'minutes').format(config.format);

            vue_instance[start_time_prop_name] = start_utc_0;
            vue_instance[end_time_prop_name] = end_utc_0;
        });

    }

}


var htmlHelper =
{
    to_anchor(url, text) {
        return String.format("<a target='_blank' href='{0}'>{1}</a>", url, text);
    }
}



var format_helper =
{
    to_display_time(iso_8601) {

        var myTimeZone = 8;

        var parseDate = this.iso_8601_add_timeZone(iso_8601, myTimeZone);

        return parseDate.format('YYYY-MM-DD HH:mm:ss');
    },

    iso_8601_add_timeZone(iso_8601, timeZone) {

        var parseDate = moment.utc(iso_8601);
        parseDate.add(timeZone, 'hours');

        return parseDate;

    },


    toMoneyDisplayFormat: function (value) {

        if (value == null || typeof value == 'undefined')
            return "";

        // 加上千分位符號
        return parseFloat(value)
            .toFixed(2)
            .toString()
            .replace(/^(-?\d+?)((?:\d{3})+)(?=\.\d+$|$)/, function (all, pre, groupOf3Digital) {
                return pre + groupOf3Digital.replace(/\d{3}/g, ',$&');
            });

    }
}


var arrayHelper =
{
    removeLastArrayItem(array) {

        array.splice(-1, 1)
    }
}


var clipboardHelper = {

    copyElement: function ($element) {
        this.copyText($element.text())
    },

    copyText: function (text) // Linebreaks with \n
    {
        var $tempInput = $("<textarea>");
        $("body").append($tempInput);
        $tempInput.val(text).select();
        document.execCommand("copy");
        $tempInput.remove();
    }
};


//ajax 請求API 請統一在此呼叫
var ajaxHelper =
{
    isLockAlertAgain: false,

    get: function (url) {

        var headers = this.getHeaders();

        return new Promise((resolve, reject) => {
            axios.get(url, { 'headers': headers })
                .then((response) => {

                    //view那邊處理
                    resolve(response);

                })
                .catch((error) => {

                    console.error('---ajax error---');
                    console.log(error.response);

                    ajaxHelper.authTokenErrorProcess(error.response);
                    reject(error.response);

                });
        });
    },

    post: function (url, data) {

        var headers = this.getHeaders();

        return new Promise((resolve, reject) => {

            var config =
            {
                headers: headers,
            };

            axios.post(url, data, config)
                .then((response) => {

                    //view那邊處理
                    resolve(response);

                })
                .catch((error) => {

                    console.log('---ajax error---');
                    console.log(error.response);

                    ajaxHelper.authTokenErrorProcess(error.response);
                    reject(error.response);
                });
        });
    },

    put: function (url, data) {

        var headers = this.getHeaders();

        return new Promise((resolve, reject) => {

            var config =
            {
                headers: headers,
            };

            axios.put(url, data, config)
                .then((response) => {

                    //view那邊處理
                    resolve(response);

                })
                .catch((error) => {

                    console.log('---ajax error---');
                    console.log(error.response);

                    ajaxHelper.authTokenErrorProcess(error.response);
                    reject(error.response);
                });
        });
    },



    authTokenErrorProcess(res) {

        let my = this;

        if (res.data.code >= 401001 && res.data.code <= 401099) {
            if (!my.isLockAlertAgain) {

                my.isLockAlertAgain = true;
                Cookies.set(constant.KEY_COOKIE_TOKEN, '');

                setTimeout(() => {
                    my.isLockAlertAgain = false;
                }, 1000);

                alert(res.data.message);
                location.href = "/Login";

            }

        }
    },


    getHeaders: function () {

        let token = [REDACTED];
        let device = localStorage.getItem("X-Device-ID");

        var alwaysHeader =
        {
            'Content-Type': 'application/json;charset=UTF-8',
            deviceID: device,
            accessToken: [REDACTED],
            locale: 'zh-cn'
        };

        return alwaysHeader;
    }
}


var form_helper =
{
    //使用方式請參考 Card/List那頁
    getSearchFormQueryString(vue_dataTable) {

        let paging =
        {
            pageSize: vue_dataTable.pagination.pageSize,
            pageIndex: vue_dataTable.pagination.pageIndex,
            orderByColumn: vue_dataTable.pagination.orderByColumn,
        };

        return this.parseObjectToQueryString(vue_dataTable, paging);
    },

    //使用方式請參考 Card/List那頁
    getSearchFormAndSpecificPageSize(vue_dataTable, pageSize, orderByColumn) {

        let paging =
        {
            pageSize: pageSize,
            pageIndex: 1,
            orderByColumn: orderByColumn,
        };

        return this.parseObjectToQueryString(vue_dataTable, paging);
    },


    parseObjectToQueryString(vue_dataTable, pagingConfig) {

        let queryString = "";

        let pageSize = pagingConfig.pageSize;
        let pageIndex = pagingConfig.pageIndex ? pagingConfig.pageIndex : 1;
        let orderByColumn = pagingConfig.orderByColumn ? pagingConfig.orderByColumn : '';

        queryString = "?pageSize=" + pageSize + "&pageIndex=" + pageIndex + "&orderByColumn=" + orderByColumn;

        var keys = Object.keys(vue_dataTable.searchForm);

        for (var i = 0; i < keys.length; i++) {
            var keyValue = keys[i];
            var value = vue_dataTable.searchForm[keyValue];

            queryString += "&" + keyValue + "=" + value;
        }

        return queryString;
    },


    resetForm(vue_searchForm) {

        var keys = Object.keys(vue_searchForm);

        for (var i = 0; i < keys.length; i++) {
            var keyValue = keys[i];
            vue_searchForm[keyValue] = "";
        }
    }


}




var common =
{
    exportExcel: function (config) {

        //config參數 請傳入類似以下這樣的物件
        var config_sample =
        {
            headers: ["MemberID", "MemberName", "CreateTime"],
            rows: [
                { "MemberID": 1234, "MemberName": "test001", "CreateTime": "2020-10-26" },
                { "MemberID": 1001, "MemberName": "test002", "CreateTime": "2020-10-26" }
            ],
            fileName: "export_members_20200101.xlsx"
        }

        var createXLSLFormatObj = [];

        createXLSLFormatObj.push(config.headers);
        $.each(config.rows, function (index, value) {
            var innerRowData = [];
            $.each(value, function (ind, val) {
                innerRowData.push(val);
            });
            createXLSLFormatObj.push(innerRowData);
        });


        /* File Name */
        var filename = config.fileName;

        /* Sheet Name */
        var ws_name = "Sheet1";

        if (typeof console !== 'undefined') console.log(new Date());
        var wb = XLSX.utils.book_new(),
            ws = XLSX.utils.aoa_to_sheet(createXLSLFormatObj);


        let columnWidthSetting = [];

        $.each(config.headers, function (index, value) {

            //預設欄位寬度
            columnWidthSetting.push({ width: 25 });
        });

        ws['!cols'] = columnWidthSetting;

        /* Add worksheet to workbook */
        XLSX.utils.book_append_sheet(wb, ws, ws_name);

        /* Write workbook and Download */
        if (typeof console !== 'undefined') console.log(new Date());

        XLSX.writeFile(wb, filename);
        if (typeof console !== 'undefined') console.log(new Date());
    }

}


//全域擴充方法
Date.prototype.addHours = function (h) {
    this.setTime(this.getTime() + (h * 60 * 60 * 1000));
    return this;
}