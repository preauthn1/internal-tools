using System;
using System.Collections.Generic;
using System.Text;
using System.Linq;

namespace MSC.Models.Models
{
    public class MerchantCustomizedSettingConvertModel
    {
        public int MerchantID { get; set; }

        public int CustomizedSettingKey { get; set; }

        public string SettingValue { get; set; }

        public List<decimal> SettingValueListDecimal
        {
            get
            {
                if (string.IsNullOrWhiteSpace(SettingValue))
                    return new List<decimal>();

                try
                {
                    var values = SettingValue.Split(',').Select(x => decimal.Parse(x)).ToList();
                    return values;
                }
                catch
                {
                    return new List<decimal>();
                }
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

                if (SettingValue == "1")
                    return true;

                if (SettingValue == "0")
                    return false;

                Boolean.TryParse(SettingValue, out var tryValue);
                return tryValue;
            }

        }
        

    }


    public class CustomizedSettingParseModel
    {
        public int Key { get; set; }

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

                decimal.TryParse(SettingValue, out decimal tryValue);
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
