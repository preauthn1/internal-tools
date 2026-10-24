

var Global_default_no_filter = { 'text': '--ALL--', 'value': '' };


//vue js 類似繼承用的機制
var global_mixin = {

    data: {
        template: constant.template,               
        dropdown:
        {
            //如果沒有事先定義好空array 在redner有機率發生問題
            merchants: [],
            merchantsWithCurrency : [],
            currency: [],
            normalStatus: [],
            roles: []
        }
    },

    filters: {
        toMoneyDisplayFormat: function (value) {

            return format_helper.toMoneyDisplayFormat(value);
        },
        
        to_display_time: function (iso_8601) {
            return format_helper.to_display_time(iso_8601);
        },
    },

    methods: {

        isEmptyString(text) {

            if (text === undefined)
                return true;

            if (text == null)
                return true;

            if (text.length <= 0)
                return true;

            return false;
        },

        showMask() {
            let target = document.getElementById("global_overlay");

            if (target === undefined || target == null)
                return;

            target.style.display = "block";            
        },

        hideMask() {
            let target = document.getElementById("global_overlay");

            if (target === undefined || target == null)
                return;

            target.style.display = "none";
        },


        get_query_string: function (key) {

            const urlParams = new URLSearchParams(window.location.search);
            const myParam = urlParams.get(key);

            if (myParam === undefined || myParam === null || myParam === 'null')
                return "";

            return myParam;
        },


        to_display_time: function (iso_8601) {
            return format_helper.to_display_time(iso_8601);
        },

        to_display_time_or_empty: function (data, columnName) {

            if (!data || data == null) {
                return "-";
            }
            return format_helper.to_display_time(data[columnName]);
        },



        quickDatePick(begin, end) {
            
            app.dataTable.searchForm.createTimeBegin = begin;
            app.dataTable.searchForm.createTimeEnd = end;
            app.dataTable.searchForm.createTime_displayText = begin + ' ~ ' + end; 

            app.dataTable.searchForm.beginDate = begin;
            app.dataTable.searchForm.endDate = end;
            app.dataTable.searchForm.dateRange_displayText = begin.replace('00:00:00', '') + ' ~ ' + end.replace('00:00:00', '');
        },

        popupText(text) {

            $("#popupTextModal").modal("show");
            document.getElementById("popupTextModalContent").innerHTML = text;            

        },

        resetForm() {

            if (this.dataTable && this.dataTable.searchForm) {
                form_helper.resetForm(this.dataTable.searchForm);
            }
        },

        style_to_opcity50_when_null(item) {

            if (!item || item == null) {
                return " opcity50 ";
            };

            return "";
        },


        closeAllModal() {            
            $(".modal").modal("hide");
        },
                

        copyTextToClipboard(text) {

            clipboardHelper.copyText(text);

            let target = "#success-copy-alert";
            $(target).fadeIn();

            setTimeout(() => {

                $(target).fadeOut();

            }, 2000);

        },


        successCopyAlert() {

            let target = "#success-copy-alert";
            $(target).fadeIn();

            setTimeout(() => {

                $(target).fadeOut();

            }, 2000);
        },



        successAlert() {

            let target = "#success-alert";
            $(target).fadeIn();

            setTimeout(() => {

                $(target).fadeOut();

            }, 2000);
        },

        //寫操作紀錄 (某些行為要透過前端)
        writeOperateLog(type) {

            let url = window.API_DOMAIN + '/operate-logs';

            let postData =
            {
                operationType: type
            };

            ajaxHelper.post(url, postData)
                .then((response) => {

                })
                .catch((error) => { console.error(error) })
                .finally(() => { });

        },

        getQueryString() {
            var vars = [], hash;
            var hashes = window.location.href.slice(window.location.href.indexOf('?') + 1).split('&');
            for (var i = 0; i < hashes.length; i++) {
                hash = hashes[i].split('=');
                vars.push(hash[0]);
                vars[hash[0]] = hash[1];
            }
            return vars;
        },

        getDropdownCurrency: function () {

            return new Promise((resolve, reject) => {

                let url = window.API_DOMAIN + '/dropdown/currency';

                ajaxHelper.get(url)
                    .then((response) => {

                        //view那邊處理
                        this.dropdown.currency = response.data.data;
                        resolve(response);

                    })
                    .catch((error) => { console.error(error) })
                    .finally(() => { });
            });
        },


        getDropdownRoles: function () {

            return new Promise((resolve, reject) => {

                let url = window.API_DOMAIN + '/dropdown/roles';

                ajaxHelper.get(url)
                    .then((response) => {

                        //view那邊處理
                        this.dropdown.roles = response.data.data;
                        resolve(response);

                    })
                    .catch((error) => { console.error(error) })
                    .finally(() => { });
            });
        },


        getDropdownMerchants: function () {

            return new Promise((resolve, reject) => {

                let url = window.API_DOMAIN + '/dropdown/merchants';

                ajaxHelper.get(url)
                    .then((response) => {

                        //view那邊處理
                        this.dropdown.merchants = response.data.data;
                        resolve(response);

                    })
                    .catch((error) => { console.error(error) })
                    .finally(() => { });
            });
        },

        getDropdownMerchantsWithCurrency: function () {

            return new Promise((resolve, reject) => {

                let url = window.API_DOMAIN + '/dropdown/merchantsWithCurrency';

                ajaxHelper.get(url)
                    .then((response) => {

                        //view那邊處理
                        this.dropdown.merchantsWithCurrency = response.data.data;
                        resolve(response);

                    })
                    .catch((error) => { console.error(error) })
                    .finally(() => { });
            });
        },

        getDropdownNormalStatus: function () {

            return new Promise((resolve, reject) => {

                let url = window.API_DOMAIN + '/dropdown/status';

                ajaxHelper.get(url)
                    .then((response) => {

                        //view那邊處理
                        this.dropdown.normalStatus = response.data.data;
                        resolve(response);

                    })
                    .catch((error) => { console.error(error) })
                    .finally(() => { });
            });
        }   


    }
}
