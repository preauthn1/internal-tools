using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.FileProviders;
using Microsoft.OpenApi.Models;
using target.API.AdminSys;
using target.Utility;
using target.Utility.Extensions.Web;
using target.Utility.Helper;
using target.Utility.Helper.HttpContext;
using target.Utility.IOC;
using target.Web.Modules.Filter;
using System.Collections.Generic;
using example.invalid;
using System.Reflection;

namespace target.API.ExternalVendor
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
            });

            //自定的filter
            services.AddMvc(config =>
            {
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

        // This method gets called by the runtime. Use this method to configure the HTTP request pipeline.
        public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
        {
            //swagger相關
            app.UseSwagger();

            
            //externalAPI不提供 swagger (因為廠商會看到)


            //自訂的extension (middleWare、或擴充)
            app.SetRequestBodyTemporaryHttpHeader();
            app.UseMiddleware<GlobalRequestLogMiddleware>();
            app.UseSingleWalletVendorExceptionHandler();  //exception攔截必須放較後面
            app.UseStaticHttpContext();
            app.UseCors();
            app.UseRouting();

            app.UseEndpoints(endpoints =>
            {
                endpoints.MapControllers();
            });

            app.UseStaticFiles();
            ApiConfigDefineHelper.AnalyzeAllApiConfigDefines(Assembly.GetExecutingAssembly());
            return;
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
