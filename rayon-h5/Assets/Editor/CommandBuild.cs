using UnityEngine;
using UnityEditor;
using example.invalid;
using System.Collections.Generic;

public class CommandBuild {
        static void IncrementBundleVersion() {
        string versionText = example.invalid.File.ReadAllText("version.txt");

        if (versionText != null) {
            versionText = versionText.Trim(); //clean up whitespace if necessary
            string[] lines = versionText.Split('.');

            int MajorVersion = int.Parse(lines[0]);
            int MinorVersion = int.Parse(lines[1]);
            int SubMinorVersion = int.Parse(lines[2]) + 1; //increment here

            versionText = MajorVersion.ToString() + "." +
            MinorVersion.ToString() + "." +
            SubMinorVersion.ToString();

            PlayerSettings.bundleVersion = versionText;
            PlayerSettings.Android.bundleVersionCode = SubMinorVersion;
            example.invalid.File.WriteAllText("version.txt",versionText);
        }
    }
    public static List<BuildInfo> m_PoolBuild = new List<BuildInfo>();
    public static void DoCommonBuildStuff(string outPath,string BundleID) {
        string[] levels = { "Assets/Start.unity" };
 

        PlayerSettings.SetApplicationIdentifier(BuildTargetGroup.Android, BundleID);

        
        int _length = Application.dataPath.LastIndexOf('/');
        PlayerSettings.Android.keystoreName = Application.dataPath.Substring(0, _length) + "/app.keystore";
        PlayerSettings.Android.keystorePass = "123456";
        PlayerSettings.Android.keyaliasName = "cqxx";
        PlayerSettings.Android.keyaliasPass = "123456";

        //IncrementBundleVersion();

        BuildPipeline.BuildPlayer(levels, outPath, BuildTarget.Android, BuildOptions.None);
    }

    private static void ChangeXml(string BundleID)
    {
        string _targetPath = Path.Combine(Application.dataPath, "Plugins/Android/AndroidManifest.xml");

        if (File.Exists(_targetPath))
        {
            try
            {
                FileStream fileStream = new FileStream(_targetPath, FileMode.Open);
                StreamReader streamWriter = new StreamReader(fileStream);
                string _content = streamWriter.ReadToEnd();
                fileStream.Close();

                string[] _lineSeparator = new string[] { "\r\n" };
                string[] _lines = _content.Split(_lineSeparator, System.StringSplitOptions.None);
                string _keyWord = "android:host=\"";
                for (int i = 0; i < _lines.Length; ++i)
                {
                    if (_lines[i].IndexOf(_keyWord) >= 0)
                    {
                        _lines[i] = _lines[i].Substring(0, _lines[i].IndexOf(_keyWord));
                        _lines[i] += _keyWord + BundleID + "\"";
                    }
                }
                _content = "";
                for (int i = 0; i < _lines.Length; ++i)
                    _content += _lines[i] + "\r\n";

                File.WriteAllText(_targetPath, _content);
                AssetDatabase.Refresh();
            }
            catch (IOException e)
            {
                Debug.Log(e);
            }
        }
    }

    private static void ChangeStaticData(string  _key,string _info)
    {
       
        string _targetPath = Path.Combine(Application.dataPath, "Start.unity");

        if (File.Exists(_targetPath))
        {
            try
            {
                FileStream fileStream = new FileStream(_targetPath, FileMode.Open);
                StreamReader streamWriter = new StreamReader(fileStream);
                string _content = streamWriter.ReadToEnd();
                string newcontent = "";
                string line = "";
                string _keyWord = _key + ":";
                while ((line = streamWriter.ReadLine()) != null)
                {
                    if (line.IndexOf(_keyWord) >= 0)
                    {
                        newcontent += _keyWord + " " + _info;
                    }
                    else
                    {
                        newcontent += line + "\r\n";
                    }
                        
                    
                }

                fileStream.Close();
                
                //string[] _lineSeparator = new string[] { "\r\n" };
                //string[] _lines = _content.Split(_lineSeparator, System.StringSplitOptions.None);
                //string _keyWord = _key+":";
                //for (int i = 0; i < _lines.Length; ++i)
                //{
                //    if (_lines[i].IndexOf(_keyWord) >= 0)
                //    {
                //        _lines[i] = _lines[i].Substring(0, _lines[i].IndexOf(_keyWord));
                //        _lines[i] += _keyWord +" "+ _info ;
                //    }
                //}
                //_content = "";
                //for (int i = 0; i < _lines.Length; ++i)
                //    _content += _lines[i] + "\r\n";

                File.WriteAllText(_targetPath, _content);
                AssetDatabase.Refresh();
            }
            catch (IOException e)
            {
                Debug.Log(e);
            }
        }
    }

    public static void BuildAndroid() {

        PlayerSettings.Android.keystoreName = "C:\\Users\\Administrator\\Desktop\\ArkGame\\app.keystore";
        m_PoolBuild = PlatformRouteController.GetRouteList();

        SettingBuild(m_PoolBuild.Find(x => x.SiteID == 110001 && x.ServerOption == 4 && x.PlatformID == 110));
        DoCommonBuildStuff("C:\\Users\\Administrator\\Desktop\\Client777testapp.apk", "com.qqqp.android001");
    }

    private static void SettingBuild(BuildInfo _info)
    {
        BuildInfo m_BuildInfo = _info;
        if (_info == null)
            return;

        // ק]w
        PlayerSettings.companyName = m_BuildInfo.CompanyName;
        PlayerSettings.productName = m_BuildInfo.ProductName;
        PlayerSettings.SplashScreen.show = false;

        PlayerSettings.SetApplicationIdentifier(BuildTargetGroup.Android, m_BuildInfo.BundleID);


        ChangeStaticData("ServerOption", m_BuildInfo.ServerOption.ToString());
        ChangeStaticData("PlatformID", m_BuildInfo.PlatformID.ToString());
        ChangeStaticData("StaticSiteID", m_BuildInfo.StaticSiteID.ToString());

        ChangeXml(m_BuildInfo.BundleID);

       
    }
}