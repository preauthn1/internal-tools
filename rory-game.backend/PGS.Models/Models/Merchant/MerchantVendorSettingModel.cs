using target.Models.Enums;
using System;
using System.Collections.Generic;
using System.Text;

namespace target.Models.Models
{
    public class MerchantVendorSettingModel
    {

        public int SettingKey { get; set; }

        public string SettingValue { get; set; }


        public VendorConfigSettingEnum SettingKeyEnum
        {
            get
            {
                return (VendorConfigSettingEnum)this.SettingKey;
            }

        }
               

        public int SettingValueOfInt
        {
            get
            {
                if (string.IsNullOrWhiteSpace(SettingValue))
                    return 0;

                Int32.TryParse(SettingValue, out int tryValue);
                return tryValue;
            }

        }

        public decimal SettingValueOfDecimal 
        { 
            get
            {
                if (string.IsNullOrWhiteSpace(SettingValue))
                    return 0;

                decimal.TryParse(SettingValue,  out decimal tryValue);
                return tryValue;
            }
        
        }

        public bool SettingValueOfBoolean
        {
            get
            {
                if (string.IsNullOrWhiteSpace(SettingValue))
                    return false;

                Boolean.TryParse(SettingValue, out var tryValue);
                return tryValue;
            }

        }
        

    }
}
