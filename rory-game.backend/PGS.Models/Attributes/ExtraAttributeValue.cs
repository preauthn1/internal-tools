using System;
using System.Text.Json;
using System.Text.Json.Serialization;
using System.Text.RegularExpressions;
using target.Models.Extensions;

namespace target.Models.Attributes
{
    public class ExtraAttributeValue : System.Attribute
    {
        public string Value { get; set; }

        public ExtraAttributeValue (string value)
        {
            this.Value = value;
        }
    }

    public class ExtraBooleanValue : System.Attribute
    {
        public bool BooleanValue { get; set; }

        public ExtraBooleanValue(string memo , bool value)
        {
            this.BooleanValue = value;
        }
    }


    public class SettingDefaultValue : System.Attribute
    {
        public string Value { get; set; }

        public SettingDefaultValue(string value)
        {
            this.Value = value;
        }
        public SettingDefaultValue(int value)
        {
            this.Value = value.ToString();
        }
    }

}
