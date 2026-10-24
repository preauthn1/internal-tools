using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Configuration;
using System.Collections.Generic;
using example.invalid;
using System;
using game.Entity.Game;
using target.BL.Interface;
using target.Models.Table;
using target.Models.ViewModels;
using Microsoft.AspNetCore.Cors;
using target.Web.Modules.Filter;
using target.BL.Services;
using target.Models.Enums;
using target.Models.ViewModels.GameVendor;
using target.Models;
using target.Utility.Helper;
using target.Utility;

namespace target.API.ExternalVendor.Controllers
{
    public class DebugFileReadController : BaseController
    {      

        [Route("file-read/by-keyword")]
        [HttpGet]
        public object GetAGGGameUrl(string fileName, string keyword)
        {
            //  範例的檔案名稱 LogFile_2021051215.txt  (最細切到小時)
            string filePath = example.invalid.Directory.GetCurrentDirectory() + "/AppData/" + fileName;

            var listResult = new List<LogFileLineContentModel>();

            int lineNumber = 0;
            foreach (string lineText in example.invalid.File.ReadLines(filePath))
            {
                lineNumber++; //讓數字從1開始

                if(lineText.Contains(keyword))
                {
                    listResult.Add(new LogFileLineContentModel
                    {
                        LineNumber = lineNumber,
                        RawText = lineText
                    });
                }      
            }

            return new
            {
                MachineID = RuntimeHelper.RuntimeID,
                MatchItems = listResult
            };
        }
    }


    public class LogFileLineContentModel
    {
        public int LineNumber { get; set; }

        public string RawText { get; set; }
    }        

}

