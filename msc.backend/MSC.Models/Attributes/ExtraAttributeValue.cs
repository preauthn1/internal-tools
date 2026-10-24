namespace MSC.Models.Attributes
{
    public class ExtraAttributeValue : System.Attribute
    {
        public string Value { get; set; }

        public ExtraAttributeValue (string value)
        {
            this.Value = value;
        }
    }


}
