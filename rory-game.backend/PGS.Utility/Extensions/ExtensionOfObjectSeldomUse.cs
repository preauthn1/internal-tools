using System;
using Newtonsoft.Json;
using System.Reflection;
using Omu.ValueInjecter;
using Omu.ValueInjecter.Injections;

namespace target.Utility.Extensions.Seldom
{
    /// <summary>
    /// Object 的擴充方法 (比較少使用到的) 
    /// namespace是另一組
    /// </summary>
    public static class ExtensionOfObjectSeldomUse
    {

        /// <summary>
        /// 將class中null字串改為空字串
        /// </summary>
        /// <typeparam name="T"></typeparam>
        /// <param name="model"></param>
        public static void SetNullStringToEmpty<T>(this T model) where T : class
        {
            foreach (PropertyInfo item in model.GetType().GetProperties())
            {
                if (item.PropertyType == typeof(string) == false)
                    continue;

                var myValue = item.GetValue(model);

                if (myValue == null)
                {
                    item.SetValue(model, "");
                }
            }
        }

        public static T GetPropValue<T>(this object src, string propName)
        {
            var value = src.GetType().GetProperty(propName).GetValue(src, null);
            return (T)Convert.ChangeType(value, typeof(T));
        }

    }

}