using System.Collections;
using System.Collections.Generic;
using System.Diagnostics;
using example.invalid;
using UnityEngine;
using UnityEditor;

public class ReleaseScript : MonoBehaviour
{
    private static string m_sourceFolder;
    private static string m_releasePath;
    private static string m_gitFolder = "ClientRelease";


    public static void AutoRunRelease(string _sourceFolder,string Id)
    {
        m_sourceFolder = _sourceFolder;
        m_releasePath = Application.dataPath + "/..";
        // 創新資料夾
        if (Directory.Exists(m_releasePath) == false)
            Directory.CreateDirectory(m_releasePath);

        UnityEngine.Debug.Log("_sourceFolder__" + _sourceFolder);
        UnityEngine.Debug.Log("_releasePath__" + m_releasePath);

        string _content = "echo Start Wait Copy";
        _content += "\n" + "TIMEOUT /T 10";
        _content += "\n" + "echo Start Copy";

        //step1. copy folder
        _content += CopyFolderToStorage();
       
        //step6. delete
        _content += "\n" + "cd " + m_releasePath;
        _content += "\n" + "wget http://0.0.0.0/api/gv/CompleteBuild/" + Id;
        _content += "\n" + "TIMEOUT /T 2";
        _content += DeleteFolder();       
        _content += "\n" + "pause";

        string _batPath = Path.Combine(m_releasePath, "release.bat");
        File.WriteAllText(_batPath, _content);

        Process process = new Process();
        process.StartInfo.FileName = _batPath;
        process.Start();
    }

    public static void RunRelease(string _sourceFolder)
    {
        m_sourceFolder = _sourceFolder;
        m_releasePath = Application.dataPath + "/..";
        // 創新資料夾
        if (Directory.Exists(m_releasePath) == false)
            Directory.CreateDirectory(m_releasePath);

        UnityEngine.Debug.Log("_sourceFolder__" + _sourceFolder);
        //UnityEngine.Debug.Log("_releasePath__" + m_releasePath);

        string _content = "";
        //step1. copy folder
        _content += CopyFolderToQA();
        //step2. git clone
        _content += GitClone();
        //step3. git reset
        _content += "\n" + "cd " + m_gitFolder;
        _content += GitReset();
        //step4. copy folder
        _content += CopyFolder();
        //step5. git commit push
        _content += GitUpload();
        //step6. delete
        _content += "\n" + "cd..";
        _content += DeleteFolder();
        _content += "\n" + "pause";

        string _batPath = Path.Combine(m_releasePath, "release.bat");
        File.WriteAllText(_batPath, _content);

        Process process = new Process();
        process.StartInfo.FileName = _batPath;
        process.Start();
    }

    public static void RunSynchronize(string _sourceFolder)
    {
        m_sourceFolder = _sourceFolder;
        m_releasePath = Application.dataPath + "/..";
        // 創新資料夾
        if (Directory.Exists(m_releasePath) == false)
            Directory.CreateDirectory(m_releasePath);

        UnityEngine.Debug.Log("_sourceFolder__" + _sourceFolder);
        //UnityEngine.Debug.Log("_releasePath__" + m_releasePath);

        GetRedmineLog();

        string _content = "";
        //step1. copy folder
        _content += CopyFolderToQA();
        //step2. delete
        _content += DeleteFolder();
        _content += "\n" + "pause";

        string _batPath = Path.Combine(m_releasePath, "release.bat");
        File.WriteAllText(_batPath, _content);

        Process process = new Process();
        process.StartInfo.FileName = _batPath;
        process.Start();
    }

    private static string GitClone()
    {
        string _content = "\n" + "git clone http://0.0.0.0/ClientTeam/" + m_gitFolder  + ".git";

        return _content;
    }

    private static string CopyFolder()
    {
        DirectoryInfo _info = new DirectoryInfo(m_sourceFolder);
        string _destinationFolder = m_releasePath + "/" + m_gitFolder + "/" + _info.Name;

        string _content = "\n" + "Xcopy /E /I";
        _content += " \"" + m_sourceFolder + "\"";
        _content += " \"" + _destinationFolder + "\"";

        return _content;
    }

    private static string GitReset()
    {
        string _content = "\n" + "git reset --hard";
        _content += "\n" + "git fetch";
        _content += "\n" + "git checkout Deploy";
        _content += "\n" + "git pull --rebas";

        return _content;
    }

    private static string GitUpload()
    {
        DirectoryInfo _info = new DirectoryInfo(m_sourceFolder);
        string _content = "";
        _content += "\n" + "git config --global core.autocrlf false";
        _content += "\n" + "git add " + _info.Name;
        _content += "\n" + "git commit -m \'"+ GetRedmineLog() +"\'";
        _content += "\n" + "git push";
        _content += "\n" + "git tag " + _info.Name;
        _content += "\n" + "git push origin --tags";
        _content += "\n pause";
               
        return _content;
    }

    private static string DeleteFolder()
    {
        string _content = "\n" + "del /f release.bat";

        return _content;
    }

    private static string CopyFolderToQA()
    {
        string _destinationFolder = "\\\\0.0.0.0\\update\\UpdateFiles";

        string _content = "\n" + "Xcopy /E /I";
        _content += " \"" + m_sourceFolder + "\\CN\"";
        _content += " \"" + _destinationFolder + "\"";
        _content += "\n" + "pause";

        return _content;
    }
    private static string CopyFolderToStorage()
    {
        string _destinationFolder = "Z:\\UpdateFiles\\"+ m_sourceFolder.Split('/')[m_sourceFolder.Split('/').Length-1];

        string _content ="\n" + "net use Z: \\\\0.0.0.0\\98.clientrelease /user:domain\\bob example.invalid /persistent:no";
        _content +="\n" + "Xcopy /E /I";
        _content += " \"" + m_sourceFolder + "\"";
        _content += " \"" + _destinationFolder + "\"";
      

        return _content;
    }

    private static string GetRedmineLog()
    {
        DirectoryInfo _info = new DirectoryInfo(m_sourceFolder);
        string _platform = AppDefine.PlatformPath;
        string _projectList = "";
        string _filePath = Application.dataPath + "/StreamingAssets/" + _platform;
        string[] filenames = Directory.GetFiles(_filePath);
        foreach (var e in filenames)
        {
            FileInfo _file = new FileInfo(e);
            if (_file.Extension != ".txt")
                continue;

            if (_file.Name.IndexOf("version") < 0)
                continue;

            _projectList += "_";
            _projectList += _file.Name.Replace("_version.txt", "");
            FileStream _fs = new FileStream(_file.FullName, FileMode.Open);
            StreamReader _sr = new StreamReader(_fs); // 以串流讀取
            _projectList += "(" + _sr.ReadToEnd() + ")";
            _sr.Close();
            _fs.Close();
        }

        string _content = "[" + _info.Name + "][" + _platform + "]" + _projectList;
        UnityEngine.Debug.Log("to redmine log : " + _content);

        return _content;
    }

    private static Dictionary<string, List<string>> m_HallFolder = new Dictionary<string, List<string>>()
    {        
#if UNITY_IOS
        {"3",new List<string>(){ "CN"} },
#else
        {"3",new List<string>(){ "CN" } },
#endif
        {"4",new List<string>(){ "CN/UI_4_1" } },
        {"5",new List<string>(){ "CN/UI_5_1" } },
        {"6",new List<string>(){ "CN/UI_6_1" } },
        {"7",new List<string>(){ "CN/UI_7_1" } },
        {"8",new List<string>(){ "CN/UI_8_1" } },
        {"6_986",new List<string>(){ "CN/UI_6_1_986" } },
        {"6_Blue",new List<string>(){ "CN/UI_6_1_Blue" } },
    };
    private static Dictionary<string, List<string>> m_GameFolder = new Dictionary<string, List<string>>()
    {    
#if UNITY_IOS
        {"3",new List<string>(){ "CN"} },
#else
        {"3",new List<string>(){ "CN" } },
#endif
        {"4",new List<string>(){ "CN/UI_4_1" } },
        {"5",new List<string>(){ "CN/UI_5_1" } },
        {"6",new List<string>(){ "CN/UI_6_1" } },
        {"7",new List<string>(){ "CN/UI_7_1" } },
        {"8",new List<string>(){ "CN/UI_8_1" } },
    };

    public static string CopyFilesToHotFixFolder(bool showmsg=true)
    {
        // 本地上传更新根路径
        string updateRootPath = Path.GetFullPath(AppDefine.LOCAL_INIT_DATA_PATH);
        if (Directory.Exists(updateRootPath) == false)
        {
            EditorUtility.DisplayDialog("通知", "發布失敗!!!沒有要發布的資料", "確定");
            return null; 
        }

        DirectoryInfo rootFolder = new DirectoryInfo(AppDefine.LOCAL_INIT_DATA_PATH);
        DirectoryInfo[] directoryInfo = rootFolder.GetDirectories();
        
        // UpdateFiles\20180524_190233\bjl\1.0.0\Android\
        string baseDataPath = Application.dataPath + "/../UpdateFiles/" + System.DateTime.Now.ToString("yyyyMMdd_HHmmss");

        int Acc = 0;
        foreach (DirectoryInfo data in directoryInfo)
        {
            // 取得設定檔版本號
            string version = FileUtility.ReadFileText(data.Name + "_version.txt", PathType.InitData);

            // 取得設定檔样板
            string skin = FileUtility.ReadFileText(data.Name + "_skin.txt", PathType.InitData);
            string[] skinList = skin.Split(new char[] { ',' }, System.StringSplitOptions.RemoveEmptyEntries);

            if (skinList.Length == 0)
            {
                EditorUtility.DisplayDialog("通知", data.Name + "_skin.txt 查無發布樣板", "確定");
            }
            else
            {
                for (int i = 0; i < skinList.Length; ++i)
                {
                    List<string> OutPaths = null;
                    if (data.Name == "Hall" && m_HallFolder.ContainsKey(skinList[i]))
                        OutPaths = m_HallFolder[skinList[i]];
                    else if (m_GameFolder.ContainsKey(skinList[i]))
                        OutPaths = m_GameFolder[skinList[i]];

                    if (OutPaths == null)
                        EditorUtility.DisplayDialog("通知", "查無發布樣板", "確定");

                    foreach (string OutPath in OutPaths)
                    {
                        EditorUtility.DisplayProgressBar("Copy update assets", "Copy Update assets...", Acc / directoryInfo.Length);

                        // UpdateFiles\20180524_190233\bjl\1.0.0\Android\
                        string dataPath = baseDataPath + "/" + OutPath + "/" + data.Name + "/" + version + "/" + AppDefine.PlatformPath;

                        string localUpdateRootPath = Path.GetFullPath(dataPath);

                        if (Directory.Exists(localUpdateRootPath))
                            Directory.Delete(localUpdateRootPath, true);

                        // 创建文件夹
                        Directory.CreateDirectory(localUpdateRootPath);

                        CopyFolderToUpdate(data, localUpdateRootPath);
                        Acc++;
                    }
                }
            }
        }
        EditorUtility.ClearProgressBar();
        if(showmsg)
            EditorUtility.DisplayDialog("通知", "發布完成!!!", "確定");

        EditorUtility.RevealInFinder(baseDataPath);

        return baseDataPath;
    }

    private static void CopyFolderToUpdate(DirectoryInfo data, string localUpdateRootPath)
    {
        // 上層要複製的資料
        List<string> copyFile = new List<string>();
        copyFile.Add("_version.txt");
        copyFile.Add("_BuildLog.txt");
        copyFile.Add("_FileList.xml");

        DirectoryInfo rootFolder = new DirectoryInfo(AppDefine.LOCAL_INIT_DATA_PATH);

        // 直接複製上層資料過去
        foreach (string _copyFileData in copyFile)
        {
            string _copyFileDataName = data.Name + _copyFileData;
            FileUtility.CopyFile(rootFolder + "/" + _copyFileDataName, localUpdateRootPath + "/" + _copyFileDataName);
        }

        rootFolder = new DirectoryInfo(AppDefine.LOCAL_INIT_DATA_PATH + "/" + data.Name);
        if (rootFolder.Exists)
        {
            FileInfo[] fileInfos = rootFolder.GetFiles("*.*", SearchOption.AllDirectories);
            if (fileInfos.Length > 0)
            {
                for (int i = 0; i < fileInfos.Length; i++)
                {
                    FileInfo fileInfo = fileInfos[i];
                    if (fileInfo.Extension == ".meta")
                        continue;

                    string fileName = fileInfo.FullName.Replace("\\", "/").Replace(AppDefine.LOCAL_INIT_DATA_PATH.Replace("\\", "/"), "").TrimStart('/');
                    string s = AppDefine.LOCAL_INIT_DATA_PATH + "/" + data.Name + "/" + fileInfo.Name;
                    string d = localUpdateRootPath + "/" + fileName;
                    FileUtility.CopyFile(s, d);
                }
                EditorUtility.ClearProgressBar();
            }
        }
    }
    
    public static void RefreshGameVersion(System.Action _callback)
    {
        string _url = "http://0.0.0.0/api/gv/g/" + AppDefine.CurrentProjectFolder;

        UnityEngine.Debug.Log(_url);
        WWW www = new WWW(_url);

        bool _check = false;
        while (www.isDone == false)
        {
            if (www.isDone == true)
            {
                _check = true;
                ReceiveGameVersion(www, _callback);
                break;
            }
        }
        if (www.isDone == true && _check == false)
            ReceiveGameVersion(www, _callback);
    }

    public static void GetCurrentGame(System.Action<Dictionary<string, object>> _callback)
    {
        string _url = "http://0.0.0.0/api/gv/curpack/A/";

        UnityEngine.Debug.Log(_url);
        WWW www = new WWW(_url);

        bool _check = false;
        while (www.isDone == false)
        {
            if (www.isDone == true)
            {
                _check = true;
                SetCurrentGame(www, _callback);
                break;
            }
        }
        if (www.isDone == true && _check == false)
            SetCurrentGame(www, _callback);
    }

    public static void CallProcess(string Id)
    {
        string _url = "http://0.0.0.0/api/gv/ProcessBuild/"+Id;

        UnityEngine.Debug.Log(_url);
        WWW www = new WWW(_url);

        bool _check = false;
        while (www.isDone == false)
        {
            if (www.isDone == true)
            {
                break;
            }
        }
       
    }

    public class UpdateGameVersion
    {
        public string Id;
        public string Game;
        public string Platform;
        public string Ver;
        public string UpdateAt;
    }

    private static void ReceiveGameVersion(WWW www, System.Action _callback)
    {
        if (string.IsNullOrEmpty(www.error) == false)
        {
            UnityEngine.Debug.Log(www.error);
        }
        else
        {
            UnityEngine.Debug.Log(www.text);

            string[] _lines = www.text.Split(new string[] { "[{", "},{", "}]" }, System.StringSplitOptions.RemoveEmptyEntries);
            for (int i = 0; i < _lines.Length; ++i)
            {
                UpdateGameVersion content = JsonUtility.FromJson<UpdateGameVersion>("{" + _lines[i] + "}");

                GameConfig config = ProjectEditUtility.GetCurrentProjectConfig();

                if (config != null && content.Platform == "A" && AppDefine.PlatformPath == "Android")
                {
                    config.Version = content.Ver;
                    EditorUtility.SetDirty(config);
                }
                else if (config != null && content.Platform == "I" && AppDefine.PlatformPath == "iOS")
                {
                    config.Version = content.Ver;
                    EditorUtility.SetDirty(config);
                }
            }

            if (_callback != null)
                _callback();
        }
    }

    private static void SetCurrentGame(WWW www, System.Action<Dictionary<string, object>> _callback)
    {
        if (string.IsNullOrEmpty(www.error) == false)
        {
            UnityEngine.Debug.Log(www.error);
        }
        else
        {
            UnityEngine.Debug.Log(www.text);

            Dictionary<string, object> g = MiniJSON.Json.Deserialize(www.text) as Dictionary<string, object>;



            if (_callback != null)
                _callback(g);
        }
    }
}