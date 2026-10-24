using System;
using System.Collections.Generic;
using System.Text;
using System.Linq;

namespace target.Models.Models
{
    public class SystemConfigSettingParseModel
    {

        public int SettingKey { get; set; }

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

                Boolean.TryParse(SettingValue, out var tryValue);
                return tryValue;
            }

        }
        public List<string> SettingValueOfListString
        {
            get
            {
                if (string.IsNullOrWhiteSpace(SettingValue))
                    return new List<string>();

                try
                {
                    var split = SettingValue.Split(',').ToList();
                    return split;
                }
                catch
                {

                }

                return new List<string>();
            }

        }


        public List<int> SettingValueOfListInt
        {
            get
            {
                if (string.IsNullOrWhiteSpace(SettingValue))
                    return new List<int>();

                try
                {
                    var split = SettingValue.Split(',').Select(z=> Int32.Parse(z)).ToList();
                    return split;
                }
                catch
                {

                }

                return new List<int>();
            }

        }



    }
}
