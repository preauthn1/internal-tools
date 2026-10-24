using System.Collections;
using System.Collections.Generic;
using System.Diagnostics;
using UnityEditor;
using UnityEngine;
using example.invalid;

public class GitCommit
{
    public static void Log(string _targetPath)
    {
        string _commit = "git log";
        RunProcess(_commit, _targetPath);
    }

    public static void Reset(string _targetPath)
    {
        string _commit = "git reset --hard";
        RunProcess(_commit, _targetPath);
    }

    public static void Show(string _targetPath)
    {
        string _commit = "git show";
        RunProcess(_commit, _targetPath);
    }

    public static void Pull(string _targetPath)
    {
        string _commit = "git reset --hard";
        _commit += "\n" + "git pull";
        RunProcess(_commit, _targetPath);
    }

    public static void Checkout(string _targetPath,string _branch)
    {
        string _commit = "git reset --hard";
        _commit += "\n" + "git checkout " + _branch;
        RunProcess(_commit, _targetPath);
    }

    private static string RunProcess(string command,string _targetPath = "")
    {
        string _path = Application.dataPath + "/..";

        string _content = "";
        if (string.IsNullOrEmpty(_targetPath) == false)
            _content += "cd " + _targetPath;
        _content += "\n" + command;
        //_content += "\n" + "pause";
        _content += "\n" + "del /f process.bat";

        string _batPath = Path.Combine(_path, "process.bat");
        File.WriteAllText(_batPath, _content);

        Process process = new Process();
        process.StartInfo.FileName = _batPath;
        process.StartInfo.RedirectStandardOutput = true;
        process.StartInfo.UseShellExecute = false;
        process.Start();

        // Read the output stream first and then wait.
        string output = process.StandardOutput.ReadToEnd();
        process.WaitForExit();

        UnityEngine.Debug.Log(output);
        string[] _lines = output.Split(new char[] { '\n' }, System.StringSplitOptions.RemoveEmptyEntries);
        string _return = "";
        for (int i = 0; i < _lines.Length; ++i)
        {
            if (_lines[i].IndexOf("cd") >= 0)
                continue;

            if (string.IsNullOrEmpty(_lines[i]))
                continue;

            if (_lines[i].IndexOf("del") >= 0)
                break;

            if (_lines[i].IndexOf(command) >= 0)
                _return += command + "\n";
            else
                _return += _lines[i] + "\n";
        }

        return _return;
    }

    public static string GetLastCommit(string _targetPath)
    {
        string _gitlog = RunProcess("git log", _targetPath);
        string[] _lines = _gitlog.Split(new char[] { '\n' }, System.StringSplitOptions.RemoveEmptyEntries);
        string _content = "";
        bool _getLast = false;
        for (int i = 0; i < _lines.Length; ++i)
        {
            if (_lines[i].IndexOf("commit") >= 0 && _getLast == false)
                _getLast = true;
            else if (_lines[i].IndexOf("commit") >= 0)
                break;

            if (_lines[i].IndexOf("Author") >= 0)
                continue;
            
            _content +=  _lines[i] + "\n";
        }
        
        UnityEngine.Debug.Log(_content);
        return _content;
    }

    public static string GetCurrentBranch(string _targetPath)
    {
        string _content = RunProcess("git rev-parse --abbrev-ref HEAD", _targetPath);
        return _content;
    }

    public static string GetBranchList(string _targetPath)
    {
        string _content = RunProcess("git branch", _targetPath);
        return _content;
    }
}
