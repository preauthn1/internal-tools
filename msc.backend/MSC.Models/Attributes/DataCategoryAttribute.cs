using MSC.Models.Enums;

namespace MSC.Models.Attributes
{
    /// <summary>
    /// 設置某些資料分類
    /// </summary>
    public class DataCategoryAttribute : System.Attribute
    {
        public DataCategoryEnum Value { get; set; }

        public DataCategoryAttribute(DataCategoryEnum value)
        {
            this.Value = value;
        }
    }




}
