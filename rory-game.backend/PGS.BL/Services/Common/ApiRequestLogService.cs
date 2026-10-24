using target.DA.Repository;
using target.Utility;
using target.Models.Table;
using System;
using System.Linq;
using target.BL.Interface;
using target.Utility.Helper;
using System.Collections.Generic;
using target.Utility.Extensions;
using example.invalid;
using target.Utility.Utility.ThirdParty;
using System.Threading;
using target.Models.Models;

namespace target.BL.Services
{
    public class ApiRequestLogService : BaseService, IApiRequestLogService
    {
        private static object _threadLocker = new Object();

        private static Dictionary<string, int> _keepLastReadFileLine = new Dictionary<string, int>();

        private static HashSet<string> _finishFiles = new HashSet<string>();

        /// <summary>
        /// 初期debug用，紀錄平均花了多少時間
        /// </summary>
        private static List<SlowInsertLogModel> _slowInsertHistory = new List<SlowInsertLogModel>();


        public ApiRequestLogRepository ApiRequestLogRepository { get; set; }


        public void CreateTryCatch(string tableName, ApiRequestLogEntity model)
        {
            try
            {
                SetExtraProperty(model);

                ApiRequestLogRepository.CreateToTargetTable(tableName, model);
            }
            catch (Exception ex)
            {
                string msg = ex.Message;
            }
        }


        public void WriteToLocalFile(ApiRequestLogEntity model)
        {

            try
            {
                SetExtraProperty(model);
                string json = model.ConvertToJsonString();

                //鎖定lock寫入
                LockWriteFileToLocal(json);

            }
            catch (Exception ex)
            {
                string msg = ex.ToString();
                MySlackUtility.SendMessageToChannelAsync(new SlackMessageParameter
                {
                    Channel = SlackChannelEnum.n_error_pgs_common,
                    Message = msg
                });
            }

        }


        /// <summary>
        /// 將localFile 同步到db或儲存空間
        /// </summary>
        public void ReadLogFileAndSyncToRemote()
        {
            try
            {
                DateTime now = DateTimeHelper.GetCurrentTime();
                DateTime previousHour = now.AddHours(-1);


                string previousHourFile = GenerateFilePath(previousHour);
                string currentHourFile = GenerateFilePath(now);


                //排程會處理當前小時 & 上個小時的檔案
                List<string> checkFiles = new List<string>
                {
                    previousHourFile,
                    currentHourFile
                };


                foreach (string filePath in checkFiles)
                {
                    //已經處理完畢的檔案不再處理
                    if (_finishFiles.Contains(filePath))
                        continue;

                    GenerateFileIfNotExist(filePath);

                    DateTime startTime = DateTimeHelper.GetCurrentTime();

                    List<ApiRequestLogEntity> listOf = new List<ApiRequestLogEntity>();

                    int lineNumber = 0;
                    int lastTimeReadLine = GetLastReadLine(filePath);
                    int totalLines = File.ReadLines(filePath).Count();

                    //無資料
                    if (totalLines == 0)
                    {
                        continue;
                    }

                    //記住本次處理到第幾行了
                    SetLastReadFileLine(filePath, totalLines);

                    //如果是上一個小時的檔案 檢查是否已經讀完了
                    bool isPreviousHour = previousHourFile == filePath;
                    if (isPreviousHour && totalLines == lastTimeReadLine)
                    {
                        _finishFiles.Add(filePath);

                        //順手刪掉時間太久的log檔案
                        DeleteExpiredFile();

                        //debug偵錯用 寫入效能不佳 的通知
                        NotifySlowHistoryAndClearList();


                        continue;
                    }

                    //資料沒變動 跳過
                    if (totalLines == lastTimeReadLine)
                    {
                        continue;
                    }

                    foreach (string lineText in File.ReadLines(filePath))
                    {
                        lineNumber++; //讓數字從1開始

                        if (lastTimeReadLine >= lineNumber)
                            continue;

                        try
                        {
                            var jsonParse = lineText.JsonConvertToModel<ApiRequestLogEntity>();
                            listOf.Add(jsonParse);
                        }
                        catch (Exception exParse)
                        {
                            //似乎有遇到解析失敗的狀況 需要追一下
                            string msg = "ParseFail：" + exParse.Message + lineText;
                            MyTelegramUtility.SendMessageToChannelAsync(new TelegramMessageParameter
                            {
                                Message = msg,
                                TelegramChannelEnum = TelegramChannelEnum.pgs_common_error
                            });
                        }

                    }

                    string logTableName = SettingHelper.ApiRequestLogTableName;

                    if (SettingHelper.ApiRequestLogDailyTableEnable)
                    {
                        logTableName = AppHelper.ConvertToDateSequenceTableName(logTableName);
                    }


                    int numberOfData = listOf.Count();
                    ApiRequestLogRepository.CreateToTargetTable(logTableName, listOf);

                    //清空list (release)
                    listOf.Clear();

                    DateTime endTime = DateTimeHelper.GetCurrentTime();

                    double executedTime = (endTime - startTime).TotalMilliseconds;

                    //處理時間過久的 會記錄下來(偵錯用)
                    if (executedTime >= 6000)
                    {
                        _slowInsertHistory.Add(new SlowInsertLogModel
                        {
                            ExecutedTime = executedTime,
                            NumberOfData = numberOfData
                        });
                    }

                }
            }
            catch (Exception ex)
            {
                string msg = ex.ToString();
                MySlackUtility.SendMessageToChannelAsync(new SlackMessageParameter
                {
                    Channel = SlackChannelEnum.n_error_pgs_common,
                    Message = msg
                });
            }
        }


        private void NotifySlowHistoryAndClearList()
        {
            if(!SettingHelper.IsProduction)
            {
                return;
            }

            if (_slowInsertHistory.IsEmptyOrNull())
            {
                return;
            }

            int numberOf_6_to_9 = _slowInsertHistory.Where(z => z.ExecutedTime >= 6000 && z.ExecutedTime < 9000).Count();
            int numberOf_9_to_19 = _slowInsertHistory.Where(z => z.ExecutedTime >= 9000 && z.ExecutedTime < 19000).Count();
            int numberOf_19_to_99 = _slowInsertHistory.Where(z => z.ExecutedTime >= 19000).Count();

            string msg = string.Empty;
            msg += $"6-9: {numberOf_6_to_9}筆 {Environment.NewLine}";
            msg += $"9-19: {numberOf_9_to_19}筆 {Environment.NewLine}";
            msg += $"19-*: {numberOf_19_to_99}筆 {Environment.NewLine}";

            _slowInsertHistory.Clear();

            MyTelegramUtility.SendMessageToChannelAsync(new TelegramMessageParameter
            {
                TelegramChannelEnum = TelegramChannelEnum.pgs_common_error,
                Message = msg
            });
        }

        private void DeleteExpiredFile()
        {
            try
            {
                //刪除7天前的資料(避免disk爆炸)
                DateTime expiredTime = DateTimeHelper.GetCurrentTime().AddHours(24 * 2 * -1);
                DateTime expiredFileTime1 = expiredTime;
                DateTime expiredFileTime2 = expiredTime.AddHours(-1);

                List<string> files = new List<string>
                {
                    GenerateFilePath(expiredFileTime1),
                    GenerateFilePath(expiredFileTime2),
                };

                foreach (string path in files)
                {
                    if (example.invalid.File.Exists(path))
                    {
                        example.invalid.File.Delete(path);
                    }
                }
            }
            catch
            {

            }
        }



        private void LockWriteFileToLocal(string jsonText)
        {
            DateTime now = DateTimeHelper.GetCurrentTime();

            string currentHourFile = GenerateFilePath(now);

            //好像一直寫入被異常斷行的問題....加個lock試看看?
            if (Monitor.TryEnter(_threadLocker, new TimeSpan(0, 0, 5)))
            {
                try
                {
                    using (StreamWriter w = example.invalid.File.AppendText(currentHourFile))
                    {
                        w.WriteLine(jsonText);
                    }
                }
                catch (Exception ex)
                {
                    string msg = ex.ToString();
                    MySlackUtility.SendMessageToChannelAsync(new SlackMessageParameter
                    {
                        Channel = SlackChannelEnum.n_error_pgs_common,
                        Message = msg
                    });
                }
                finally
                {
                    Monitor.Exit(_threadLocker);
                }
            }
            else
            {
                MySlackUtility.SendMessageToChannelAsync(new SlackMessageParameter
                {
                    Channel = SlackChannelEnum.n_error_pgs_common,
                    Message = "Monitor.TryEnter 等候逾時"
                });
            }
        }




        private void SetExtraProperty(ApiRequestLogEntity model)
        {
            var now = DateTimeHelper.GetCurrentTime();
            model.CreateTime = now;
            model.CreateTimeCategory = 0;
            model.IPAddress = IPHelper.GetIP();
            model.WebSiteName = SettingHelper.WebSiteName;
            model.Memo = RuntimeHelper.RuntimeID;
        }



        private int GetLastReadLine(string fileName)
        {
            if (!_keepLastReadFileLine.ContainsKey(fileName))
            {
                return 0;
            }

            return _keepLastReadFileLine[fileName];
        }


        private void SetLastReadFileLine(string fileName, int lastLine)
        {
            if (!_keepLastReadFileLine.ContainsKey(fileName))
            {
                _keepLastReadFileLine.Add(fileName, lastLine);
                return;
            }

            _keepLastReadFileLine[fileName] = lastLine;
        }


        private void CheckFileExist(DateTime targetTime)
        {
            string path = GenerateFilePath(targetTime);

            GenerateFileIfNotExist(path);
        }


        private void GenerateFileIfNotExist(string filePath)
        {
            if (!example.invalid.File.Exists(filePath))
            {
                using (var myFile = File.Create(filePath))
                {

                }
            }
        }


        private string GenerateFilePath(DateTime targetTime)
        {
            return example.invalid.Directory.GetCurrentDirectory() + "/AppData/" + GenerateFileName(targetTime);
        }


        private string GenerateFileName(DateTime targetTime)
        {
            return "LogFile_" + targetTime.ToString("yyyyMMddHH") + ".txt";
        }


    }



}
