using Microsoft.OpenApi.Any;
using Microsoft.OpenApi.Models;
using MSC.Models.CoreLibs;
using MSC.Models.ViewModels;
using MSC.Utility;
using MSC.Utility.Extensions;
using Swashbuckle.AspNetCore.SwaggerGen;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;


namespace MSC.Web.Modules
{
    public class PGSSwaggerGameAPICertificateSampleValueFilter : IOperationFilter
    {
        public void Apply(OpenApiOperation operation, OperationFilterContext context)
        {
            string path = context.ApiDescription.RelativePath;

            if (operation.Parameters.IsEmptyOrNull())
                operation.Parameters = new List<OpenApiParameter>();

        }
    }


    public class PGSOperationFilter : IOperationFilter
    {
        public void Apply(OpenApiOperation operation, OperationFilterContext context)
        {
            string path = context.ApiDescription.RelativePath;

            bool isDebugMatch = path.Contains("merchants");

            if (operation.Parameters.IsEmptyOrNull())
                operation.Parameters = new List<OpenApiParameter>();          

            var responseContent = operation.Responses.FirstOrDefault().Value.Content;

            bool isPaged = false;

            var classFullName = responseContent?.Values?.FirstOrDefault()?.Schema?.Reference?.Id ?? string.Empty;

            if(classFullName.Contains("PagedList"))
            {
                isPaged = true;
            }
            
            if (isDebugMatch)
            {
                string testDebug = "斷點請設在這";
            }

            var requestBody = operation.RequestBody;
            if(requestBody != null)
            {
                int loopTimes = 0; //避免物件無窮迴圈
                var match = requestBody.Content?.FirstOrDefault().Value;

                if(match != null)
                {
                    var newModel = GenerateColumnDocument(context.SchemaRepository.Schemas, "data", match.Schema, isPaged , loopTimes);

                    //因應我方API 會多包裝一層data的結構 修改swagger格式
                    match.Schema.Example = newModel;
                    match.Schema.Reference = null;
                }
            }
           
            if (responseContent.IsNotEmptyNull())
            {
                int loopTimes = 0; //避免物件無窮迴圈
                var responseSchema = responseContent.FirstOrDefault().Value;

                var newModel = GenerateColumnDocument(context.SchemaRepository.Schemas, "data", responseSchema.Schema , isPaged , loopTimes);

                //因應我方API 會多包裝一層data的結構 修改swagger格式
                responseSchema.Schema.Example = new OpenApiObject
                {
                    ["data"] = newModel
                };

                responseSchema.Schema.Reference = null;
            }
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


        private IOpenApiAny GenerateFromReferenceClass(Dictionary<string, OpenApiSchema> schemas, string refID , bool isPaged ,int loopTimes)
        {
            //避免無窮迴圈
            if (loopTimes >= 6)
                return new OpenApiString("object");

            var match = schemas.Where(x => x.Key == refID).FirstOrDefault();

            var newModel = new OpenApiObject();

            foreach (var item in match.Value.Properties)
            {
                var isExist = PGSSwaggerSpecialCaseSetting.IsKeyExist(refID , item.Key);

                if (isPaged && isExist)
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
            foreach (var item in context.Type.GetProperties())
            {
                var name = item.Name;
                var ignoreAttr = item.GetCustomAttribute(typeof(SwaggerIgnoreOutputWhenPagedList), false);

                if(ignoreAttr != null)
                {
                    string concateKeyName = $"{context.Type.Name}_{name}";

                    //滿鳥的寫法，總之就是先在這邊記住 model + 欄位名稱，後面輸出的時候會拿來判斷
                    PGSSwaggerSpecialCaseSetting.AddWhenNotExist(concateKeyName);
                }
            }            
        }
    }


    /// <summary>
    /// 第一次讀取的時候 會先把不要的欄位改名、隱蔽
    /// </summary>
    public class PGSSwaggerQueryParameterChangeParameterNameFilter : IParameterFilter
    {
        //這個檔案裏面管理某些命名的欄位會顯示自訂的sampleValue
        private static Dictionary<string, string> dic = PGSSwaggerSpecialCaseSetting.QueryParameterSampleValueDefined;

        public void Apply(OpenApiParameter parameter, ParameterFilterContext context)
        {
            var prop = context.PropertyInfo;

            parameter.Name = parameter.Name.ToCamelCase();

            if (dic.ContainsKey(parameter.Name))
            {
                parameter.Description = dic[parameter.Name];
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
        }
    }


}
