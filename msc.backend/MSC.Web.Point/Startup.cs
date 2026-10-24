using System;
using System.Collections.Generic;
using example.invalid;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.HttpsPolicy;
using Microsoft.AspNetCore.StaticFiles;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.FileProviders;
using Microsoft.Extensions.Hosting;
using MSC.Utility;
using MSC.Utility.Helper.HttpContext;
using MSC.Utility.IOC;
using MSC.Web.Modules.Filter;

namespace MSC.Web.Point
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
            services.AddControllersWithViews();

            services.AddRazorPages().AddRazorRuntimeCompilation();

            //自訂的httpContext
            services.AddMyHttpContextAccessor();

            //呼叫自訂的autofac初始化
            InitializeIOC();
        }

        // This method gets called by the runtime. Use this method to configure the HTTP request pipeline.
        public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
        {              
            if (env.IsDevelopment())
            {
                app.UseDeveloperExceptionPage();
            }
            else
            {
                app.UseExceptionHandler("/Home/Error");
                // The default HSTS value is 30 days. You may want to change this for production scenarios, see https://aka.ms/aspnetcore-hsts.
                app.UseHsts();
            }
            app.UseHttpsRedirection();
            app.UseStaticFiles();

            var provider = new FileExtensionContentTypeProvider();
            provider.Mappings[".properties"] = "text/plain";

            app.UseStaticFiles(new StaticFileOptions
            {
                FileProvider = new PhysicalFileProvider(Path.Combine(Directory.GetCurrentDirectory(), "wwwroot", "bundle")),
                RequestPath = "/bundle",
                ContentTypeProvider = provider
            });

            app.UseRouting();

            //如果沒呼叫這段 讀config會出問題
            RuntimeHelper.Initialize(env.EnvironmentName);

            //自訂的extension (middleWare、或擴充)
            app.SetRequestBodyTemporaryHttpHeader();
            app.UseWebMvcExceptionHandler();  //exception攔截必須在比較後面
            app.UseStaticHttpContext();


            app.UseEndpoints(endpoints =>
            {
                endpoints.MapControllerRoute(
                    name: "default",
                    pattern: "{controller=Home}/{action=Index}/{id?}");
            });
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
