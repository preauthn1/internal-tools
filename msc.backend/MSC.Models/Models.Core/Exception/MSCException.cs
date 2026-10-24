using Newtonsoft.Json;
using System;
using System.Collections.Generic;
using System.Text;

namespace MSC.Models.Models
{
    public class MSCException: System.Exception
    {

        public ErrorCode ErrroCode { get; set; }

        public ErrorDetailModel ErrorDetail { get; set; }


        public MSCException(ErrorCode errorCode, ErrorDetailModel errorDetail = null)
        {
            this.ErrroCode = errorCode;
            this.ErrorDetail = errorDetail;
        }


        public MSCException(ErrorCode errorCode, string replaceMsg)
        {
            this.ErrroCode = errorCode;
            this.ErrorDetail = new ErrorDetailModel
            {
                ReplaceMessage = replaceMsg
            };
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

        /// <summary>
        /// 如果要強制覆蓋訊息
        /// </summary>
        [JsonProperty(NullValueHandling = NullValueHandling.Ignore)]
        public string ReplaceMessage { get; set; }

        [JsonProperty(NullValueHandling = NullValueHandling.Ignore)]
        public object ExtraData { get; set; }

    }


}
