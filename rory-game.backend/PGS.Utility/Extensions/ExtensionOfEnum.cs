using System;
using System.Collections.Generic;
using System.Linq;
using System.ComponentModel;

namespace target.Utility.Extensions
{
    /// <summary>
    /// 列舉擴充
    /// </summary>
    public static class ExtensionOfEnum
    {
    
        /// <summary>
        /// 是否為定義好的enum編號
        /// </summary>
        public static bool IsDefinedValue(this System.Enum enumVal)
        {
            var type = enumVal.GetType();
            return System.Enum.IsDefined(type, enumVal);
        }


        /// <summary>
        /// 取得 Enum Attribute Description
        /// </summary>
        /// <param name="enumVal"></param>
        /// <returns></returns>
        public static string GetDescription(this Enum enumVal)
        {
            if (enumVal == null || enumVal.IsDefinedValue() == false)
            {
                return string.Empty;
            }

            return enumVal.GetAttributeFromEnum<DescriptionAttribute>()?.Description;
        }
        

        public static string GetDescriptionOrDefault(this Enum enumVal , string defaultValue)
        {
            if (enumVal == null || enumVal.IsDefinedValue() == false)            
                return defaultValue;            

            var attr = enumVal.GetAttributeFromEnum<DescriptionAttribute>();

            if (attr == null)
                return defaultValue;

            return attr.Description;
        }


        /// <summary>
        /// 取得Enums特定的attribute
        /// </summary>
        public static T GetAttributeFromEnum<T>(this System.Enum enumVal) where T : System.Attribute
        {
            if (!enumVal.IsDefinedValue())
                return null;

            var type = enumVal.GetType();
            var memInfo = type.GetMember(enumVal.ToString());
            var attributes = memInfo[0].GetCustomAttributes(typeof(T), false);
            return (attributes.Length > 0) ? (T)attributes[0] : null;
        }

        

        public static int ToInteger(this System.Enum value)
        {
            return value.GetHashCode();
        }

        public static byte ToByte(this System.Enum value)
        {
            return (byte)value.GetHashCode();
        }


        public static string ToIntegerString(this System.Enum value)
        {
            return value.GetHashCode().ToString();
        }



        /// <summary>
        /// 取得每一筆enum的值
        /// </summary>
        /// <typeparam name="T"></typeparam>
        /// <returns></returns>
        public static List<T> GetEnumAllValue<T>() where T : IComparable, IFormattable, IConvertible
        {
            List<T> result = new List<T>();

            foreach (T enumValue in System.Enum.GetValues(typeof(T)))
            {
                result.Add(enumValue);
            }

            return result;
        }


        public static T GetEnumByTextValue<T>(string textValue , T defaultValue) where T : IComparable, IFormattable, IConvertible
        {
            textValue = textValue.Trim();

            bool isNumeric = textValue.IsNumeric();

            if(isNumeric)
            {
                var intValue = Int32.Parse(textValue);
                return (T)Convert.ChangeType(intValue, typeof(T));
            }
            
            var allEnums = GetEnumAllValue<T>();

            foreach(var item in allEnums)
            {
                if (item.ToString().ToLower() == textValue.ToLower())
                    return item;
            }

            return defaultValue;
        }

    }
}
