using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Microsoft.OpenApi.Models;
using Newtonsoft.Json;
using target.Models.CoreLibs;
using target.Models.ViewModels;
using target.Models.ViewModels.Base;
using target.Utility;
using target.Utility.Helper;
using target.Utility.Helper.HttpContext;
using target.Utility.IOC;
using target.Web.Modules.Filter;
using Swashbuckle.AspNetCore.SwaggerGen;
using target.Utility.Extensions;
using example.invalid;
using Microsoft.AspNetCore.Routing;
using target.Utility.Extensions.Web;

namespace target.API.AdminSys
{
    public class Startup
    {

        public Startup(IConfiguration configuration)
        {
            Configuration = configuration;
        }

        public IConfiguration Configuration { get; }


        // This method gets called by the runtime. Use this method to add services to the container.
        public void ConfigureServices(IServiceCollection services)
        {
            services.AddControllers()
                    .AddNewtonsoftJson(options =>
                    { 
                        //如果沒有呼叫newtonsoft版的，int欄位一定要用數字型別 不得使用string傳入
                        options.SerializerSettings.DateFormatString = DateTimeHelper.ISO_8601_Format;
                    });

            services.AddCors(options =>
            {
                options.AddDefaultPolicy(builder =>
                {
                    builder.AllowAnyOrigin()
                           .AllowAnyMethod()
                           .AllowAnyHeader();
                });
            });

            //自訂的httpContext
            services.AddMyHttpContextAccessor();
            services.AddSwaggerGen(c =>
            {
                c.SchemaFilter<PGSSwaggerTemporarySkipOutputColumnsFilter>();
                c.ParameterFilter<PGSSwaggerQueryParameterChangeParameterNameFilter>();
                c.DocumentFilter<PGSSwaggerQueryParameterIgnoreFilter>();
                c.OperationFilter<PGSOperationFilter>();
                c.SwaggerDoc("v1", new OpenApiInfo { Title = SettingHelper.WebSiteName, Version = "v1" });
                LoadXmlDocument(c, "target.API.AdminSys");
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
            return;
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
            app.UseCors();

            app.UseRouting();
            app.UseEndpoints(endpoints =>
            {
                endpoints.MapControllers();
            });

            ApiConfigDefineHelper.AnalyzeAllApiConfigDefines(Assembly.GetExecutingAssembly());
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
