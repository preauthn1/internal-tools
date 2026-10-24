using Newtonsoft.Json;
using System;
using System.Collections.Generic;
using System.Text;

namespace target.Models.Models
{
    public class TerproException: System.Exception
    {

        public ApiErrorCode ErrroCode { get; set; }

        public ErrorDetailModel ErrorDetail { get; set; }


        public TerproException(ApiErrorCode errorCode, ErrorDetailModel errorDetail = null)
        {
            this.ErrroCode = errorCode;
            this.ErrorDetail = errorDetail;
        }


        public TerproException(ApiErrorCode errorCode, string replaceMsg)
        {
            this.ErrroCode = errorCode;
            this.ErrorDetail = new ErrorDetailModel
            {
                ReplaceMessage = replaceMsg
            };
        }

        public TerproException(ApiErrorCode errorCode, bool isAlwaysNotify)
        {
            this.ErrroCode = errorCode;
            this.ErrorDetail = new ErrorDetailModel
            {
                IsAlwaysNotify = isAlwaysNotify
            };
        }

        public override string Message
        {
            get
            {
                return $"TerproException_{(int)ErrroCode}_{ErrroCode.ToString()}_{ErrorDetail?.LogExtraInfo}";
            }
        }

        public override string ToString()
        {
            return $"TerproException_{(int)ErrroCode}_{ErrroCode.ToString()}_{ErrorDetail?.LogExtraInfo}";
        }

    }


    public class ErrorDetailModel
    {

        public ErrorDetailModel()
        {

        }


        public ErrorDetailModel(string columnName)
        {
            this.ColumnName = columnName;
        }



        public string ColumnName { get; set; }


        [JsonProperty(NullValueHandling = NullValueHandling.Ignore)]
        public string LogExtraInfo { get; set; }

        /// <summary>
        /// 是否總是發slack通知
        /// </summary>
        [JsonIgnore]
        public bool IsAlwaysNotify { get; set; }


        /// <summary>
        /// 如果要強制覆蓋訊息
        /// </summary>
        [JsonProperty(NullValueHandling = NullValueHandling.Ignore)]
        public string ReplaceMessage { get; set; }

        [JsonProperty(NullValueHandling = NullValueHandling.Ignore)]
        public object ExtraData { get; set; }

    }


}
