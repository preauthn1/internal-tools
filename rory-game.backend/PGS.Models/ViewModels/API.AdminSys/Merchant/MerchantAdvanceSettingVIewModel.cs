using Newtonsoft.Json;
using target.Models.Models;
using target.Models.Table;
using target.Models.ViewModels.API.AdminSys;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class MerchantAdvanceSettingViewModel
    {
        public int MerchantID { get; set; }

        public List<MerchantAdvanceSettingItemViewModel> Items { get; set; }

    }


    public class MerchantAdvanceSettingItemViewModel
    {
        public int CustomizedSettingKey { get; set; }

        public string CustomizedSettingDisplayName { get; set; }

        /// <summary>
        /// 基本上是給前端的一些判斷
        /// </summary>
        public string DataTypeName { get; set; }

        [JsonIgnore]
        public DefineTypeCheckResultModel DefineType { get; set; }

        /// <summary>
        /// 部分欄位僅允許透過DB修改數據 不開放普通人員編輯
        /// </summary>
        public bool IsDisableEdit { get; set; }

        public List<DropdownValueViewModel> DataDropdowns { get; set; }


        public bool IsDropdownDataType 
        { 
            get
            {
                return DataDropdowns != null && DataDropdowns.Count() > 0;
            }
        }


        public string Value { get; set; }
    }


    public class MerchantAdvanceSettingUpdateRequestModel
    {
        public int MerchantID { get; set; }

        public int CustomizedSettingKey { get; set; }

        public string Value { get; set; }
    }


}
