using System;
using System.Text.Json;
using System.Text.Json.Serialization;
using System.Text.RegularExpressions;
using target.Models.Extensions;

namespace target.Models.CoreLibs
{

    /// <summary>
    /// 影響SQLHelper，用來標記此table非自動增長auto Identity
    /// </summary>
    [AttributeUsage(AttributeTargets.Class, Inherited = false)]
    public class IsNotAutoKeyTable : System.Attribute
    {

    }

}
