using System;
using System.Collections.Generic;
using System.Text;

namespace target.Models.Models
{
    public class MerchantCustomizedSettingModel
    {
        public int MerchantID { get; set; }

        public int CustomizedSettingKey { get; set; }

        public string SettingValue { get; set; }


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

                if (SettingValue == "1")
                    return true;

                if (SettingValue == "0")
                    return false;

                Boolean.TryParse(SettingValue, out var tryValue);
                return tryValue;
            }

        }
        

    }
}
