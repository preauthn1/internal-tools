using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.ComponentModel.DataAnnotations;
using example.invalid;
using System.Linq;
using System.Reflection;
using System.Runtime.Versioning;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using target.Utility;
using target.Utility.Extensions;

namespace target.API.Vendor.Controllers
{
    [ApiExplorerSettings(IgnoreApi = true)]
    public class InternalDocumentController : BaseController
    {
        
        /// <summary>
        /// 列出後端定義的所有enum
        /// </summary>
        [HttpGet]
        [Route("enums")]
        public object GetAllEnumsList()
        {
            Dictionary<string, string> dic = new Dictionary<string, string>();

            var allEnums = GetAllEnums();
            
            string hostUrl = HttpContext.Request.Host.ToString();

            foreach (var item in allEnums)
            {
                dic.Add(item.Name, string.Format("http://{0}/enums/{1}", hostUrl, item.Name.ToString()));
            }

            return dic
                .OrderBy(x => x.Key)
                .ToDictionary(x => x.Key, z => z.Value);
        }

     
        [HttpGet]
        [Route("enums/{enumName}")]
        public dynamic GetEnumByName(string enumName)
        {
            Dictionary<string, int> enumValues = new Dictionary<string, int>();
            Dictionary<string, string> descriptionValues = new Dictionary<string, string>();


            var allEnums = GetAllEnums();

            foreach (Type type in allEnums)
            {
                if (type.Name == enumName)
                {
                    foreach (var enumValue in System.Enum.GetValues(type))
                    {
                        AddNoneDuplicateKeys(enumValues, enumValue.ToString(), (int)enumValue);

                        var desAttr = ExtensionOfEnum.GetAttributeFromEnum<DescriptionAttribute>((Enum)enumValue);

                        if (desAttr != null)
                        {
                            AddNoneDuplicateKeys(descriptionValues, enumValue.ToString(), desAttr.Description);
                        }
                    }
                    break;
                }
            }

            return new
            {
                allValues = enumValues,
                allText = descriptionValues
            };
        }


        private void AddNoneDuplicateKeys<T>(Dictionary<string, T> dic, string key, T value)
        {
            if (dic.ContainsKey(key))
                return;

            dic.Add(key, value);
        }

        

        private List<Type> GetAllEnums()
        {
            List<Type> all = new List<Type>();
            var listOfEnums1 = GetTypesInNamespace(Assembly.Load("target.Models"), "target.Models.Enums");
            var listOfEnums2 = GetTypesInNamespace(Assembly.Load("target.Models"), "target.Models.Models");

            all.AddRange(listOfEnums1);
            all.AddRange(listOfEnums2);
            return all;
        }



        private List<Type> GetTypesInNamespace(Assembly assembly, string nameSpace)
        {
            return assembly
                .GetTypes()
                .Where(t => String.Equals(t.Namespace, nameSpace, StringComparison.Ordinal) &&  !t.IsClass)
                .ToList();
        }



    }
}
