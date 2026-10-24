using Microsoft.VisualStudio.TestTools.UnitTesting;
using MSC.Models.Models;
using MSC.Models.ViewModels;
using MSC.Utility.Extensions;
using System.Collections.Generic;
using System.Linq;

namespace MSC.Testing
{
    [TestClass]
    public class QueueNotifyModelValidateTests
    {

        [TestMethod]
        public void CheckDeuplicateEnumValue()
        {
            //如果亂加新欄位 可能導致商戶notify錯誤，加個MSTest確保          
            List<string> validateDefineColumns = new List<string>
            {
                nameof(CardSuccessNotifyRequestInfo_V1_Model.CardNo),
                nameof(CardSuccessNotifyRequestInfo_V1_Model.CreateTime),
                nameof(CardSuccessNotifyRequestInfo_V1_Model.IPAddress),
                nameof(CardSuccessNotifyRequestInfo_V1_Model.MemberCode),
                nameof(CardSuccessNotifyRequestInfo_V1_Model.MerchantID),
                nameof(CardSuccessNotifyRequestInfo_V1_Model.MerchantRequestID),
                nameof(CardSuccessNotifyRequestInfo_V1_Model.MSCOrderID),
                nameof(CardSuccessNotifyRequestInfo_V1_Model.Point),
                nameof(CardSuccessNotifyRequestInfo_V1_Model.RandomID),
                nameof(CardSuccessNotifyRequestInfo_V1_Model.TS),
            };

            foreach (var item in typeof(CardSuccessNotifyRequestInfo_V1_Model).GetProperties())
            {
                var propName = item.Name;
                bool isDefined = validateDefineColumns.Any(x => x == propName);

                if(!isDefined)
                {
                    Assert.Fail($"undefined columns: {propName}");
                }               
            }
        }
    }
}
