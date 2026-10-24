using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Http;
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
using Swashbuckle.AspNetCore.SwaggerGen;
using System.Collections.Generic;
using example.invalid;
using System.Reflection;

namespace target.API.Vendor
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

            //自訂的httpContext
            services.AddMyHttpContextAccessor();
            services.AddSwaggerGen(c =>
            {
                c.SchemaFilter<PGSSwaggerTemporarySkipOutputColumnsFilter>();
                c.ParameterFilter<PGSSwaggerQueryParameterChangeParameterNameFilter>();
                c.DocumentFilter<PGSSwaggerQueryParameterIgnoreFilter>();
                c.OperationFilter<PGSOperationFilter>();

                c.SwaggerDoc("tfMode", new OpenApiInfo { Title = SettingHelper.WebSiteName, Version = "tfMode" });
                c.SwaggerDoc("swMode", new OpenApiInfo { Title = SettingHelper.WebSiteName + "(Seamless)", Version = "swMode" });

                LoadXmlDocument(c , "target.API.Vendor");
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

        // This method gets called by the runtime. Use this method to configure the HTTP request pipeline.
        public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
        {

            //需要放前面一點 避免被其他middleWare覆蓋
            app.UseStaticFiles();

            //swagger相關
            app.UseSwagger();

            if(!SettingHelper.IsProduction)
            {
                app.UseSwaggerUI(c =>
                {
                    c.SwaggerEndpoint("/swagger/tfMode/swagger.json", "target.API.TransferMode");
                    c.SwaggerEndpoint("/swagger/swMode/swagger.json", "target.API.SeamlessMode");                    
                    c.InjectStylesheet("/css/swagger_customized.css");
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


            //app.UseStaticFiles(new StaticFileOptions()
            //{
            //    FileProvider = new PhysicalFileProvider(
            //    Path.Combine(env.ContentRootPath, @"node_modules")),
            //    RequestPath = new PathString("/third-party")
            //});

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
