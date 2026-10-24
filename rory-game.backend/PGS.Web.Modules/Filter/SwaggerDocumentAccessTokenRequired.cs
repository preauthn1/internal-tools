using System;

namespace target.Web.Modules.Filter
{
    /// <summary>
    /// 純粹標記 作為Swagger API文件識別標記使用
    /// </summary>
    [AttributeUsage(AttributeTargets.Method)]
    public class SwaggerDocumentAccessTokenRequired : System.Attribute
    {
       
    }

}
