using System;
using System.Collections.Generic;
using example.invalid;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using MSC.Utility;
using MSC.Utility.Helper;
using MSC.Web.Modules.Filter;
using MSC.Utility.Helper.HttpContext;
using Swashbuckle.AspNetCore.SwaggerGen;
using MSC.Utility.IOC;
using Microsoft.OpenApi.Models;
using Microsoft.AspNetCore.Builder;
using MSC.Web.Modules;
using MSC.Utility.Extensions;
using System.Reflection;

namespace MSC.API.Vendor
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
            services.AddControllers()
                   .AddNewtonsoftJson(options =>
                   {
                       //pGSIsnewtonsoftAint@wnμƦrO oϥstringǤJ
                       options.SerializerSettings.DateFormatString = DateTimeHelper.ISO_8601_Format;
                   });

            //ۭqhttpContext
            services.AddMyHttpContextAccessor();
            services.AddSwaggerGen(c =>
            {
                c.SchemaFilter<PGSSwaggerTemporarySkipOutputColumnsFilter>();
                c.ParameterFilter<PGSSwaggerQueryParameterChangeParameterNameFilter>();
                c.DocumentFilter<PGSSwaggerQueryParameterIgnoreFilter>();
                c.OperationFilter<PGSOperationFilter>();
                c.SwaggerDoc("v1", new OpenApiInfo { Title = SettingHelper.WebSiteName, Version = "v1" });
                LoadXmlDocument(c, "MSC.API.Vendor");
            });

            //۩w~BzB
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

            //Isۭqautofacl
            InitializeIOC();

        }

        // This method gets called by the runtime. Use this method to configure the HTTP request pipeline.
        public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
        {
            //pGSIsoq Ūconfig|XD
            RuntimeHelper.Initialize(env.EnvironmentName);

            //swagger
            app.UseSwagger();

            if (!SettingHelper.IsProduction)
            {
                app.UseSwaggerUI(c =>
                {
                    c.SwaggerEndpoint("/swagger/v1/swagger.json", SettingHelper.WebSiteName);
                });
            }

            //ۭqextension (middleWareBXR)
            app.SetRequestBodyTemporaryHttpHeader();
            app.UseMiddleware<GlobalRequestLogMiddleware>();
            app.UseGlobalExceptionHandler();  //exceptiondIb᭱
            app.UseStaticHttpContext();


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
                "MSC.BL"
            };

            IOCInitializeHelper.InitializeIOC(projectList);
        }


    }
}
