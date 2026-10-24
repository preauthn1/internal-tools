using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.HttpsPolicy;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Microsoft.OpenApi.Models;
using target.Utility;
using target.Utility.Helper;
using target.Utility.IOC;
using target.Web.Modules.Filter;
using target.Utility.Helper.HttpContext;
using target.API.AdminSys;
using example.invalid;
using Swashbuckle.AspNetCore.SwaggerGen;
using Microsoft.AspNetCore.Routing;
using System.Reflection;
using target.Utility.Extensions.Web;

namespace target.API.Game
{
    public class Startup
    {
        public Startup(IConfiguration configuration)
        {
            Configuration = configuration;
        }

        public IConfiguration Configuration { get; }


        public void ConfigureServices(IServiceCollection services)
        {
            services.AddControllers(x =>
            {
                x.Filters.Add(typeof(ApiCertificateAuthorizeFilter));
            })
            .AddNewtonsoftJson(options =>
            {
                options.SerializerSettings.DateFormatString = DateTimeHelper.ISO_8601_Format;
            })

            ; //如果沒有呼叫newtonsoft版的，int欄位一定要用數字型別 不得使用string傳入

            //自訂的httpContext
            services.AddMyHttpContextAccessor();
            services.AddSwaggerGen(c =>
            {
                c.SchemaFilter<PGSSwaggerTemporarySkipOutputColumnsFilter>();
                c.ParameterFilter<PGSSwaggerQueryParameterChangeParameterNameFilter>();
                c.DocumentFilter<PGSSwaggerQueryParameterIgnoreFilter>();
                c.OperationFilter<PGSOperationFilter>();
                c.OperationFilter<PGSSwaggerGameAPICertificateSampleValueFilter>();


                c.SwaggerDoc("v1", new OpenApiInfo
                {
                    Title = SettingHelper.WebSiteName,
                    Version = "v1",
                    Description = $"所有API請求都必須帶header `X-API-Certificate` </br>若有標記 access_Token* 則需要帶header `access_Token` "
                });

                LoadXmlDocument(c, "target.API.Game");
            });

            //自定的錯誤處理、
            services.AddMvc(config =>
            {
                config.Filters.Add(new ApiResponseFormatterFilter());
                config.Filters.Add(new GlobalApiHeaderTemporaryInfoFilter());
            })
            .ConfigureApiBehaviorOptions(options =>
            {
                options.InvalidModelStateResponseFactory = context =>
                {
                    return ModalValidCustomizedResponse.ResponseCustomizedErrorResponse(context);
                };
            });

            //呼叫自訂的autofac初始化
            InitializeIOC();

        }


        public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
        {
            //swagger相關
            app.UseSwagger();

            if (!SettingHelper.IsProduction)
            {
                app.UseSwaggerUI(c =>
                {
                    c.SwaggerEndpoint("/swagger/v1/swagger.json", SettingHelper.WebSiteName);
                });
            }


            //自訂的extension (middleWare、或擴充)
            app.SetRequestBodyTemporaryHttpHeader();
            app.UseMiddleware<GlobalRequestLogMiddleware>();
            app.UseGlobalExceptionHandler();  //exception攔截必須在比較後面
            app.UseStaticHttpContext();


            app.UseRouting();
            app.UseEndpoints(endpoints =>
            {
                endpoints.MapControllers();
            });

            ApiConfigDefineHelper.AnalyzeAllApiConfigDefines(Assembly.GetExecutingAssembly());
            return;

        }


        private void LoadXmlDocument(SwaggerGenOptions c, string fileName)
        {
            try
            {
                var filePath = Path.Combine(System.AppContext.BaseDirectory, $"AppData/{fileName}.xml");
                c.IncludeXmlComments(filePath);
            }
            catch
            {

            }
        }


        private void InitializeIOC()
        {
            List<string> projectList = new List<string>
            {
                "target.BL"
            };

            IOCInitializeHelper.InitializeIOC(projectList);
        }



    }
}
