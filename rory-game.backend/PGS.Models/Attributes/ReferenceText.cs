using System;
using System.Text.Json;
using System.Text.Json.Serialization;
using System.Text.RegularExpressions;
using target.Models.Extensions;

namespace target.Models.Attributes
{
    /// <summary>
    /// 用來標記文件用的attribute 實際無用途
    /// </summary>
    [AttributeUsage(AttributeTargets.Field, AllowMultiple = true)]
    public class ReferenceText : System.Attribute
    {
        public string Value { get; set; }

        public ReferenceText(string value)
        {
            this.Value = value;
        }
    }

    public class ReferenceType : System.Attribute
    {
        public Type Type { get; set; }

        public ReferenceType(Type value)
        {
            this.Type = value;
        }
    }


}
