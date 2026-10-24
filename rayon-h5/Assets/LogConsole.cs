using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using example.invalid;

public class LogConsole : MonoBehaviour
{
    private struct LogData
    {
        public string m_log;
        public string m_stackTrace;
        public LogType m_type;
    }

    private static List<LogData> s_logs = new List<LogData>();

    private void OnEnable()
    {
        Application.RegisterLogCallback(HandleLog);
    }

    private void OnDisable()
    {
        Application.RegisterLogCallback(null);
    }

    private void HandleLog(string log, string stackTrace, LogType type)
    {
        s_logs.Add(new LogData() { m_log = log, m_stackTrace = stackTrace, m_type = type });

        if (s_logs.Count > 500) // 避免 log 數量太多
        {
            s_logs.RemoveAt(0);
        }
    }

    public static void StartLog()
    {
        s_logs.Clear();
    }

    public static void StopLog()
    {
        string _folderPath = System.Environment.GetFolderPath(System.Environment.SpecialFolder.Desktop) + "/LogConsole";

        // 創新資料夾
        if (Directory.Exists(_folderPath) == false)
            Directory.CreateDirectory(_folderPath);
        
        string _content = "";
        foreach (LogData logdata in s_logs)
            _content += logdata.m_log + "\n\n";

        string _fileName = "LogConsole_"+System.DateTime.Now.ToString("yyMMdd_HHmmss")+".txt";
        string _targetPath = Path.Combine(_folderPath, _fileName);
        File.WriteAllText(_targetPath, _content);
    }
}
