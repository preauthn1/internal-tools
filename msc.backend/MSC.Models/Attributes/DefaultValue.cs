namespace MSC.Models.Attributes
{
    public class SettingDefaultValue : System.Attribute
    {
        public string Value { get; set; }

        public SettingDefaultValue(string value)
        {
            this.Value = value;
        }
    }



}
