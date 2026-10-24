using Microsoft.OpenApi.Any;
using Microsoft.OpenApi.Models;
using target.Models;
using target.Models.Attributes;
using target.Models.CoreLibs;
using target.Models.Enums;
using target.Models.ViewModels;
using target.Utility;
using target.Utility.Extensions;
using target.Utility.Helper.HttpContext;
using target.Web.Modules.Filter;
using Swashbuckle.AspNetCore.SwaggerGen;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;


namespace target.API.AdminSys
{
    public class PGSSwaggerGameAPICertificateSampleValueFilter : IOperationFilter
    {
        public void Apply(OpenApiOperation operation, OperationFilterContext context)
        {
            string path = context.ApiDescription.RelativePath;

            if (operation.Parameters.IsEmptyOrNull())
                operation.Parameters = new List<OpenApiParameter>();

            operation.Parameters.Add(new OpenApiParameter
            {
                Name = "X-API-Certificate",
                In = ParameterLocation.Header,
                Description = "每隻API都要帶此header",
                Schema = new OpenApiSchema
                {
                    Type = "String",
                    Default = new OpenApiString(SettingHelper.GameApiSetting.ApiCertificateValue)
                }
            });

        }
    }




    public class PGSOperationFilter : IOperationFilter
    {

      
        public void Apply(OpenApiOperation operation, OperationFilterContext context)
        {
            string path = context.ApiDescription.RelativePath;

            if (operation.Parameters.IsEmptyOrNull())
                operation.Parameters = new List<OpenApiParameter>();

            operation.Description = operation.Description + GetApiDescription(context);

            var firstResponseContent = operation.Responses.FirstOrDefault().Value;
            var responseContent = operation.Responses.FirstOrDefault().Value.Content;

            bool isPaged = false;

            var classFullName = responseContent?.Values?.FirstOrDefault()?.Schema?.Reference?.Id ?? string.Empty;


            dynamic actionDesc = context.ApiDescription.ActionDescriptor;

            string methodName = (actionDesc.ActionName).ToString();

            bool isPost = context.ApiDescription.HttpMethod == "POST";

            if (classFullName.Contains("PagedList"))
            {
                isPaged = true;
            }

            bool isDebugMatch = path.Contains("merchants");
            if (isDebugMatch)
            {
                string testDebug = "斷點請設在這";
            }

            if(methodName.Contains("GetGameRecords"))
            {
                string testDebug = "斷點請設在這";
            }

            if (methodName.Contains("GetAllGames"))
            {
                string testDebug = "斷點請設在這";
            }

            var requestBody = operation.RequestBody;
            if(requestBody != null)
            {
                int loopTimes = 0; //避免物件無窮迴圈
                var match = requestBody.Content?.FirstOrDefault().Value;

                if (match != null)
                {
                    var newModels = GenerateColumnDocument(context.SchemaRepository.Schemas, "data", match.Schema, isPaged , loopTimes);

                    dynamic listOfModels = newModels;

                    //因應我方API 會多包裝一層data的結構 修改swagger格式
                    match.Schema.Example = newModels;
                    match.Schema.Reference = null;

                    if(isPost)
                    {
                        if(operation.Parameters.IsEmptyOrNull())
                        {
                            operation.Parameters = new List<OpenApiParameter>();    
                        }

                        Dictionary<string, ApiDescriptionModel> dic = PGSSwaggerSpecialCaseSetting.GetParameterDescriptionByWebsite();

                        foreach (var item in listOfModels)
                        {
                            try
                            {
                                //有些結構是物件包物件 先不處理了
                                string key = item.Key.ToString();

                                var dataType = item.Value.PrimitiveType.ToString();

                                operation.Parameters.Add(new OpenApiParameter
                                {
                                    Name = item.Key,
                                    In = ParameterLocation.Query,
                                    Description = dic.ContainsKey(key) ? dic[key].FullDescription : "",
                                    Schema = new OpenApiSchema
                                    {
                                        Type = dataType,
                                    }
                                });
                            }
                            catch
                            {

                            }
                        }
                    }

                }
            }

           
            if (responseContent.IsNotEmptyNull())
            {
                int loopTimes = 0; //避免物件無窮迴圈
                var responseSchema = responseContent.FirstOrDefault().Value;

                var newModels = GenerateColumnDocument(context.SchemaRepository.Schemas, "data", responseSchema.Schema , isPaged , loopTimes);

                //因應我方API 會多包裝一層data的結構 修改swagger格式
                responseSchema.Schema.Example = new OpenApiObject
                {
                    ["code"] = new OpenApiLong(0),
                    ["message"] = new OpenApiString(null),
                    ["data"] = newModels
                };

                dynamic listOfModels = newModels;

                Dictionary<string, ApiDescriptionModel> dic = PGSSwaggerSpecialCaseSetting.GetParameterDescriptionByWebsite();

                HashSet<string> columns = new HashSet<string>();

                foreach (var item in listOfModels)
                {
                    try
                    {
                        if(IsList(item))
                        {
                            foreach(var child in item)
                            {
                                string myKey = child.Key.ToString();
                                columns.Add(myKey);
                            }
                        }
                        else if(item.Key.ToString() == "items")
                        {
                            foreach (var child in item.Value[0])
                            {
                                string myKey = child.Key.ToString();
                                columns.Add(myKey);
                            }
                        }

                        string key = item.Key.ToString(); 
                        columns.Add(key);
                    }
                    catch
                    {

                    }
                }

                firstResponseContent.Description = "<table class='responseColmns'>";
                firstResponseContent.Description += "<tr> <th>Column</th> <th>Description</th> </tr>";

                RemoveUselessResponseColumns(columns);

                foreach (var item in columns)
                {
                    string desc = dic.ContainsKey(item) ? dic[item].Description : "";
                    firstResponseContent.Description += $"<tr> <td class='columnName'>{item}</td> <td>{desc}</td> </tr>";
                }

                firstResponseContent.Description += "</table>";
                responseSchema.Schema.Reference = null;
            }

            operation.Parameters = operation.Parameters.OrderByDescending(z => z.Name == "version ").ToList();
            return;
        }


        /// <summary>
        /// 移除无用的栏位
        /// </summary>
        private void RemoveUselessResponseColumns(HashSet<string> list)
        {
            list.Remove("items");
            list.Remove("extraObject");
            list.Remove("gameSummary");            
        }


        public bool IsList(object o)
        {
            return (o is System.Collections.IEnumerable);
        }

        private IOpenApiAny GenerateColumnDocument(Dictionary<string, OpenApiSchema> schemas, string columnName, OpenApiSchema target , bool isPaged , int loopTimes)
        {
            target.Format = target.Format ?? string.Empty;
            if (target.Reference == null && target.Items?.Reference == null)
            {
                return GenerateBaseTypeColumnOpenApiObject(schemas , columnName , target , isPaged);
            }
            else
            {
                if (target.Type == "array")
                {
                    if (target.Items.Reference != null)
                    {
                        loopTimes++;
                        return new OpenApiArray
                        {
                            GenerateFromReferenceClass(schemas , target.Items.Reference.Id , isPaged , loopTimes) 
                        };
                    }
                }
                else
                {
                    if (target.Reference != null)
                    {
                        loopTimes++;
                        return GenerateFromReferenceClass(schemas, target.Reference.Id , isPaged , loopTimes);
                    }
                }
            }

            return null;
        }


        private IOpenApiAny GenerateBaseTypeColumnOpenApiObject(Dictionary<string, OpenApiSchema> schemas, string columnName, OpenApiSchema target, bool isPaged)
        {
            if (PGSSwaggerSpecialCaseSetting.IsOutputKeyExist(columnName))
                return new OpenApiString(PGSSwaggerSpecialCaseSetting.OutputParameterSampleValueDefined[columnName]);

            if (PGSSwaggerSpecialCaseSetting.IsOutputNumberColumnsExist(columnName))
                return new OpenApiLong(PGSSwaggerSpecialCaseSetting.OutputParameterSampleNumberValueDefined[columnName]);

            if (PGSSwaggerSpecialCaseSetting.IsOutputDecimalColumnsExist(columnName))
                return new OpenApiDouble((double)PGSSwaggerSpecialCaseSetting.OutputParameterSampleDecimalValueDefined[columnName]);

            string type = target.Type.ToLower();
            switch (type)
            {
                case "string":

                    if (target.Format.Contains("date-time"))
                    {
                        return new OpenApiString("2020-01-01T12:00:00Z");
                    }

                    return new OpenApiString("string");

                case "long":
                case "number":
                case "integer":

                    if(target.Format.Contains("double"))
                    {
                        return new OpenApiDouble(0.0);
                    }

                    return new OpenApiInteger(0);


                case "boolean":
                    return new OpenApiBoolean(false);

                case "array":
                    return new OpenApiArray
                    {

                    };

                default:
                    return new OpenApiString(columnName);
            }
        }

        private string GetApiDescription(OperationFilterContext context)
        {
            string resultText = string.Empty;
            string changeLine = Environment.NewLine + "<BR>";

            if (context.MethodInfo.CustomAttributes != null)
            {
                var attr_DocumentMemo = (DocumentMemo)context.MethodInfo.GetCustomAttribute(typeof(DocumentMemo), false);
                var desc_EN = (SwaggerApiDescriptionEN)context.MethodInfo.GetCustomAttribute(typeof(SwaggerApiDescriptionEN), false);
                var desc_CN = (SwaggerApiDescriptionCN)context.MethodInfo.GetCustomAttribute(typeof(SwaggerApiDescriptionCN), false);


                if (attr_DocumentMemo != null)
                {
                    resultText += changeLine + attr_DocumentMemo.Memo;
                }

                if(desc_EN != null && PGSSwaggerSpecialCaseSetting.IsLanguageEN())
                {
                    resultText += changeLine + desc_EN.Description;
                }
                else if (desc_CN != null && !PGSSwaggerSpecialCaseSetting.IsLanguageEN())
                {
                    resultText += changeLine + desc_CN.Description;
                }

                if (IsAttributeExist<SwaggerDocumentAccessTokenRequired>(context) ||
                    IsAttributeExist<PlayerTokenAuthorizeFilter>(context) ||
                    IsAttributeExist<AdminTokenrAuthorizeFilter>(context))
                {
                    resultText += changeLine + " 此API需帶header `access_token` ";
                }

                if (IsAttributeExist<SignKeyVerifyFilter>(context))
                {
                    resultText += changeLine + " 此API需帶header `X-SignKey` (加密驗證簽) ";
                }

            }

            return resultText;
        }


        private bool IsAttributeExist<T>(OperationFilterContext context) where T : class
        {
            dynamic attr = context.MethodInfo.GetCustomAttribute(typeof(T), false);
            var convert = (T)attr;
            return convert != null;
        }




        private IOpenApiAny GenerateFromReferenceClass(Dictionary<string, OpenApiSchema> schemas, string refID , bool isPaged ,int loopTimes)
        {
            //避免無窮迴圈
            if (loopTimes >= 6)
                return new OpenApiString("object");

            var match = schemas.Where(x => x.Key == refID).FirstOrDefault();

            var newModel = new OpenApiObject();

            foreach (var item in match.Value.Properties)
            {
                var isIgnoreWhenPaged = PGSSwaggerSpecialCaseSetting.IsExistInSkipOutputResponseColumnWhenPagedList(refID , item.Key);

                //若有標記attribute 特定的欄位 不會輸出到swagger文件 (1)
                if (isPaged && isIgnoreWhenPaged)
                    continue;

                var isIgnore = PGSSwaggerSpecialCaseSetting.IsExistInSkipOutputRequestBodyList(refID, item.Key);

                //若有標記attribute 特定的欄位 不會輸出到swagger文件 (2)
                if (isIgnore)
                    continue;

                //if (PGSSwaggerSpecialCaseSetting.OutputParameterIgnoreOutputDocumentColumns.Contains(item.Key))
                //{
                //    newModel.Add(item.Key, new OpenApiString("無窮延伸的巢狀物件"));
                //    continue;
                //}

                newModel.Add(item.Key, GenerateColumnDocument(schemas, item.Key, item.Value , isPaged , loopTimes));
            }

            return newModel;
        }

    }



    /// <summary>
    /// 第一次讀取的時候 會先把不要的欄位記住、稍後輸出的時候會進行隱蔽
    /// </summary>
    public class PGSSwaggerTemporarySkipOutputColumnsFilter : ISchemaFilter
    {
        public void Apply(OpenApiSchema schema, SchemaFilterContext context)
        {
            var type = context.ApiModel.Type;
          
            foreach (var item in context.ApiModel.Type.GetProperties())
            {
                var name = item.Name;
                var ignoreAttr = item.GetCustomAttribute(typeof(SwaggerIgnoreOutputWhenPagedList), false);
                var ignoreAttr2 = item.GetCustomAttribute(typeof(SwaggerIgnoreAttribute), false);

                if (ignoreAttr != null)
                {
                    string concateKeyName = $"{context.ApiModel.Type.Name}_{name}";

                    //滿鳥的寫法，總之就是先在這邊記住 model + 欄位名稱，後面輸出的時候會拿來判斷
                    PGSSwaggerSpecialCaseSetting.AddToPagedIgnoreDictionary(concateKeyName);
                }

                if (ignoreAttr2 != null)
                {
                    string concateKeyName = $"{context.ApiModel.Type.Name}_{name}";
                    PGSSwaggerSpecialCaseSetting.AddToRequestBodyIgnoreDictionary(concateKeyName);
                }
            }            
        }
    }


    /// <summary>
    /// 第一次讀取的時候 會先把不要的欄位改名、隱蔽
    /// </summary>
    public class PGSSwaggerQueryParameterChangeParameterNameFilter : IParameterFilter
    {
      
        public void Apply(OpenApiParameter parameter, ParameterFilterContext context)
        {
            var prop = context.PropertyInfo;

            parameter.Name = parameter.Name.ToCamelCase();

            bool isEN = false;

            if(MyHttpContext.Current.Request.Cookies.ContainsKey(SystemConstant.HeaderKey_Lang))
            {
                var cookieValue = MyHttpContext.Current.Request.Cookies[SystemConstant.HeaderKey_Lang];
                isEN = cookieValue == LanguageEnum.en_us.ToString();
            }

            //這個檔案裏面管理某些命名的欄位會顯示自訂的sampleValue        
            Dictionary<string, ApiDescriptionModel> dic = PGSSwaggerSpecialCaseSetting.GetParameterDescriptionByWebsite();

            //記得移除
            //記得移除
            //記得移除
            //記得移除
            //記得移除
            //記得移除
            //記得移除
            //記得移除
            //記得移除
            //accessJwtToken

            if (parameter.Name.Contains("accessJwtToken"))
            {
                string qqqqqq = "";
            }


            if (parameter.Name.Contains("version"))
            {
                parameter.Description = "Api Version (v1)";
                parameter.Required = false;
            }

            if (dic.ContainsKey(parameter.Name))
            {
                parameter.Description = dic[parameter.Name].FullDescription;
            }

            //搞不太懂swagger的生命週期 要第二次進入頁面 才會把標記 SwaggerIgnoreAttribute 的欄位隱藏
            //無可奈何之下只好在第一次的時候
            if (prop != null)
            {
                var ignoreAttr = context.PropertyInfo.GetCustomAttribute(typeof(SwaggerIgnoreAttribute), false);

                if (ignoreAttr != null)
                {
                    parameter.Name = "XXXXXXXXXX";
                    parameter.Description = "--------------此欄位已隱蔽，請無視或F5重新刷新頁面--------------";
                }
            }
        }
    }



    /// <summary>
    /// 產生swagger文件的時候，移除掉有被標記特定attr的欄位 (沒有研究的很好，有些bug...)
    /// </summary>
    public class PGSSwaggerQueryParameterIgnoreFilter : IDocumentFilter
    {
        public void Apply(OpenApiDocument swaggerDoc, DocumentFilterContext context)
        {
            var model = context.ApiDescriptions;

            List<string> removePropertyNames = new List<string>();

            foreach (var description in context.ApiDescriptions)
            {
                var debug = description.ActionDescriptor;
                var actiob = debug.DisplayName;

                foreach (var each in description.ParameterDescriptions)
                {
                    var myModelType = each.ModelMetadata?.ContainerType;

                    if (myModelType != null)
                    {
                        //檢查是否有特殊標記的attribute
                        var prop = myModelType.GetProperty(each.Name);

                        if (prop == null)
                            continue;

                        var ignoreAttr = prop.GetCustomAttribute(typeof(SwaggerIgnoreAttribute), false);

                        if (ignoreAttr != null)
                        {
                            removePropertyNames.Add(each.Name);
                        }
                    }
                }

                //如果有被標記要ignore的欄位，移除之
                if (removePropertyNames.IsNotEmptyNull())
                {
                    foreach (string removeName in removePropertyNames)
                    {
                        var match = description.ParameterDescriptions.Where(x => x.Name == removeName).FirstOrDefault();
                        description.ParameterDescriptions.Remove(match);
                    }
                }
            }




            bool isVendorAPI = PGSSwaggerSpecialCaseSetting.IsVendorAPI();

            //如果是vendorAPI 才有的特殊處理
            if(isVendorAPI)
            {
                if(PGSSwaggerSpecialCaseSetting.IsLanguageEN())
                {
                    example.invalid.Description = "expand to following items to see more detail.";
                }
                else
                {
                    example.invalid.Description = "展开各个项目，查看更多API说明";
                }

                var sortItems = new OpenApiPaths();
                var sort = swaggerDoc.Paths.OrderBy(z => z.Value?.Operations?.FirstOrDefault().Value?.Summary);

                foreach (var item in sort)
                {
                    sortItems.Add(item.Key, item.Value);
                }

                swaggerDoc.Paths = sortItems;
            }

            return;

        }
    }


}
