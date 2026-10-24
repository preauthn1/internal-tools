using System;
using System.Collections.Generic;
using System.Linq;
using System.Text.Encodings.Web;
using System.Text.Unicode;
using System.Threading.Tasks;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Mvc.Razor;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.WebEncoders;
using MSC.Utility;
using MSC.Utility.Helper.HttpContext;
using MSC.Utility.IOC;
using MSC.Web.AdminSys.Models.Core;
using MSC.Web.Modules.Filter;

namespace MSC.Web.AdminSys
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

            services.AddMyHttpContextAccessor();

            services.Configure<RazorViewEngineOptions>(options => {
                options.ViewLocationExpanders.Add(new ViewLocationExpander());
            });

            services.AddSingleton(HtmlEncoder.Create(UnicodeRanges.All));
            services.Configure<WebEncoderOptions>(op => op.TextEncoderSettings = new TextEncoderSettings(UnicodeRanges.All));

            InitializeIOC();

            return;
        }

        // This method gets called by the runtime. Use this method to configure the HTTP request pipeline.
        public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
        {   
            //pGSIsoq Ūconfig|XD
            RuntimeHelper.Initialize(env.EnvironmentName);

            if (env.IsDevelopment())
            {
                app.UseDeveloperExceptionPage();
            }
            else
            {
                app.UseExceptionHandler("/Home/Error");
            }
            app.UseStaticFiles();
            app.UseRouting();
            app.UseStaticHttpContext();
            app.UseWebMvcExceptionHandler();  //exceptiondIb᭱

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
