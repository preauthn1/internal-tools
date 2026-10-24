using System.Collections;
using System.Collections.Generic;
using UnityEditor;
using UnityEditor.SceneManagement;
using UnityEngine;
using UnityEngine.SceneManagement;
using System;
using example.invalid;
using System.Text;

public class EditorCoroutine
{
    public static EditorCoroutine start(IEnumerator _routine)
    {
        EditorCoroutine coroutine = new EditorCoroutine(_routine);
        coroutine.start();
        return coroutine;
    }

    readonly IEnumerator routine;
    EditorCoroutine(IEnumerator _routine)
    {
        routine = _routine;
    }

    void start()
    {
        //Debug.Log("start");
        EditorApplication.update += update;
    }
    public void stop()
    {
        //Debug.Log("stop");
        EditorApplication.update -= update;
    }

    void update()
    {
        /* NOTE: no need to try/catch MoveNext,
         * if an IEnumerator throws its next iteration returns false.
         * Also, Unity probably catches when calling EditorApplication.update.
         */

        //Debug.Log("update");
        if (!routine.MoveNext())
        {
            stop();
        }
    }
}

public class BuildWindow : EditorWindow
{
    private static int s_Branch = 0;

    [MenuItem("Publish/Build Setting")]
    public static void ShowWindow()
    {
        GetWindow<BuildWindow>(false, "Build Setting", true);
    }

    private static List<BuildInfo> m_RouteList = new List<BuildInfo>();
    private static List<BuildInfo> m_FilterList = new List<BuildInfo>();
    private static List<BuildInfo> m_PublishList = new List<BuildInfo>();
    private static CustomDefine m_CustomDefine;
    private static BuildInfo m_BuildInfo;
    private static bool m_FoldoutMore;
    private static int m_FilterServerOption;
    private static string m_FilterContent;
    private static string m_FilterCompanyName;
    private static string m_FilterProductName;
    private static string m_FilterBundleID;
    private static string m_FilterAppName;
    private static string m_FilterPlatformID;
    private static string m_FilterStaticSiteID;

    private static bool m_FoldoutOpen = true;
    private static bool m_FoldoutCurrent;
    private static bool m_FoldoutPublish;
    private static List<bool> m_FoldoutFilterList = new List<bool>();
    private static List<bool> m_FoldoutPublishList = new List<bool>();
    private Vector2 m_ScrollVeiwRouteList;
    private Vector2 m_ScrollVeiwPublishList;

    private void OnFocus()
    {
        m_FoldoutCurrent = true;
        m_FoldoutPublish = true;
        m_RouteList = PlatformRouteController.GetRouteList();
        //ChangeFilter();

        if (m_CustomDefine != null)
            SettingBuild(GetBuildInfoByID(PlayerSettings.applicationIdentifier, m_CustomDefine.StaticSiteID));
    }

    void OnGUI()
    {
        // filter
        EditorGUI.BeginChangeCheck();
        EditorGUILayout.BeginVertical(GUI.skin.box);
        EditorGUILayout.Foldout(m_FoldoutOpen, "Filter");
        if (m_FoldoutOpen == true)
        {
            EditorGUILayout.BeginHorizontal();
            List<string> filterList = new List<string>() { "全部", "0.開發", "1.QA", "5.Stage", "3.美服", "4.國服" };
            for (int i = 0; i < filterList.Count; ++i)
            {
                string _content = filterList[i];
                if (m_FilterServerOption == i)
                    _content = "> " + _content + " <";
                if (GUILayout.Button(_content))
                    m_FilterServerOption = i;
            }
            EditorGUILayout.EndHorizontal();

            m_FilterContent = EditorGUILayout.TextField("FilterContent", m_FilterContent);

            m_FoldoutMore = EditorGUILayout.Foldout(m_FoldoutMore, "FilterMore");
            if (m_FoldoutMore == true)
            {
                m_FilterCompanyName = EditorGUILayout.TextField("CompanyName", m_FilterCompanyName);
                m_FilterProductName = EditorGUILayout.TextField("ProductName", m_FilterProductName);
                m_FilterBundleID = EditorGUILayout.TextField("BundleID", m_FilterBundleID);
                m_FilterAppName = EditorGUILayout.TextField("AppName", m_FilterAppName);
                m_FilterPlatformID = EditorGUILayout.TextField("PlatformID", m_FilterPlatformID);
                m_FilterStaticSiteID = EditorGUILayout.TextField("StaticSiteID", m_FilterStaticSiteID);
            }
        }
        EditorGUILayout.EndVertical();
        EditorGUILayout.Space();
        // 搜尋條件異動
        if (EditorGUI.EndChangeCheck() == true)
            ChangeFilter();

        //----
        if (m_CustomDefine == null)
        {
            foreach (CustomDefine _define in Resources.FindObjectsOfTypeAll(typeof(CustomDefine)) as CustomDefine[])
            {
                m_CustomDefine = _define;

                m_BuildInfo = GetBuildInfoByID(PlayerSettings.applicationIdentifier, m_CustomDefine.StaticSiteID);
            }
        }


        // current
        EditorGUILayout.BeginHorizontal();
        EditorGUILayout.BeginVertical(GUI.skin.window, GUILayout.MinWidth(200), GUILayout.MaxWidth(400));
        if (GUILayout.Button("Open Player Settings", GUILayout.Height(30)))
            EditorApplication.ExecuteMenuItem("Edit/Project Settings/Player");

        m_FoldoutCurrent = EditorGUILayout.Foldout(m_FoldoutCurrent, "Current Route");
        if (m_FoldoutCurrent == false)
            m_FoldoutPublish = true;
        else if (m_FoldoutCurrent == true)
        {
            if (m_CustomDefine == null)
            {
                EditorGUILayout.HelpBox("Please change scene to [Start] then you can use me.", MessageType.Warning);
            }
            else if (m_BuildInfo == null)
                EditorGUILayout.HelpBox("Current BundleID and StaticSiteID does not match assign platformID.Please choose Route List of right side", MessageType.Warning);
            else
            {
                EditorGUILayout.LabelField("CompanyName", m_BuildInfo.CompanyName);
                EditorGUILayout.LabelField("ProductName", m_BuildInfo.ProductName);
                EditorGUILayout.LabelField("BundleID", m_BuildInfo.BundleID);
                EditorGUILayout.LabelField("AppName", m_BuildInfo.AppName);
                EditorGUILayout.LabelField("AppleID", m_BuildInfo.AppleID);
                EditorGUILayout.LabelField("IconName", m_BuildInfo.IconName);
                EditorGUILayout.LabelField("ServerOption", m_BuildInfo.ServerOption.ToString());
                EditorGUILayout.LabelField("Branch", m_BuildInfo.Branch.ToString());
                EditorGUILayout.LabelField("PlatformID", m_BuildInfo.PlatformID.ToString());
                EditorGUILayout.LabelField("SiteID", m_BuildInfo.SiteID.ToString());
                EditorGUILayout.LabelField("StaticSiteID", m_BuildInfo.StaticSiteID.ToString());
                EditorGUILayout.LabelField("FlagShare", m_BuildInfo.FlagShare ? "V" : "X");
                EditorGUILayout.LabelField("NotExport", (m_BuildInfo.NotExport == 0) ? "V" : "X");

                if (GUILayout.Button("Change SplashImage and Icon", GUILayout.Height(30)))
                {
                    CopyImage();
                    AssetDatabase.Refresh();
                }

               
                if (Directory.Exists(Application.dataPath + "/FrameWork/XLua/Gen") == false)
                    EditorGUILayout.HelpBox("Publish > 生成Wrap腳本", MessageType.Warning);
                else if (GUILayout.Button("Build", GUILayout.Height(30)))
                    Build();
            }
        }

#if !UNITY_IOS
        //publish list
        if (m_PublishList.Count > 0)
            EditorGUILayout.BeginVertical(GUI.skin.box, GUILayout.MinWidth(300), GUILayout.MaxWidth(500));
        else
            EditorGUILayout.BeginVertical(GUI.skin.box);
        m_FoldoutPublish = EditorGUILayout.Foldout(m_FoldoutPublish, "Publish List");
        if (m_FoldoutPublish == false)
            m_FoldoutCurrent = true;
        else if (m_FoldoutPublish == true)
        {
            EditorGUILayout.BeginVertical(GUI.skin.box);
            m_ScrollVeiwPublishList = EditorGUILayout.BeginScrollView(m_ScrollVeiwPublishList);

            for (int i = 0; i < m_PublishList.Count; ++i)
            {
                EditorGUILayout.BeginHorizontal(GUI.skin.box);
                EditorGUILayout.BeginVertical();
                m_FoldoutPublishList[i] = EditorGUILayout.Foldout(m_FoldoutPublishList[i], m_PublishList[i].GetKeyWord());
                if (m_FoldoutPublishList[i] == true)
                {
                    EditorGUILayout.LabelField("CompanyName", m_PublishList[i].CompanyName);
                    EditorGUILayout.LabelField("ProductName", m_PublishList[i].ProductName);
                    EditorGUILayout.LabelField("BundleID", m_PublishList[i].BundleID);
                    EditorGUILayout.LabelField("AppName", m_PublishList[i].AppName);
                    EditorGUILayout.LabelField("AppleID", m_PublishList[i].AppleID);
                    EditorGUILayout.LabelField("IconName", m_PublishList[i].IconName);
                    EditorGUILayout.LabelField("ServerOption", m_PublishList[i].ServerOption.ToString());
                    EditorGUILayout.LabelField("Branch", m_PublishList[i].Branch.ToString());
                    EditorGUILayout.LabelField("PlatformID", m_PublishList[i].PlatformID.ToString());
                    EditorGUILayout.LabelField("SiteID", m_PublishList[i].SiteID.ToString());
                    EditorGUILayout.LabelField("StaticSiteID", m_PublishList[i].StaticSiteID.ToString());
                    EditorGUILayout.LabelField("FlagShare", m_PublishList[i].FlagShare ? "V" : "X");
                    EditorGUILayout.LabelField("NotExport", (m_PublishList[i].NotExport == 0) ? "V" : "X");
                }
                EditorGUILayout.EndVertical();
                if (GUILayout.Button("X", GUILayout.Width(50)))
                {
                    m_PublishList.RemoveAt(i);
                    m_FoldoutPublishList.RemoveAt(i);

                    if (m_PublishList.Count == 0)
                    {
                        m_FoldoutCurrent = true;
                        m_FoldoutPublish = false;
                    }
                }
                EditorGUILayout.EndHorizontal();
            }

            EditorGUILayout.EndScrollView();

            if (m_CustomDefine != null)
            {
                if (Directory.Exists(Application.dataPath + "/FrameWork/XLua/Gen") == false)
                    EditorGUILayout.HelpBox("Publish > 生成Wrap腳本", MessageType.Warning);
                else if (GUILayout.Button("Build List", GUILayout.Height(30)))
                {
                    EditorCoroutine.start(BuildList());
                }
            }

            EditorGUILayout.EndVertical();
        }
        EditorGUILayout.EndVertical();
#endif
        EditorGUILayout.EndVertical();
        EditorGUILayout.Space();

        // route list
        EditorGUILayout.BeginVertical(GUI.skin.box, GUILayout.MinWidth(300), GUILayout.MaxWidth(500));
        EditorGUILayout.BeginHorizontal();
        if (GUILayout.Button("Refresh Route List ( 刷新 )", GUILayout.Height(30)))
        {
            PlatformRouteController.Refresh(() => {
                m_RouteList = PlatformRouteController.GetRouteList();
                ChangeFilter();
            });
        }
        if (GUILayout.Button("Create New Route ( 新增 )", GUILayout.Height(30)))
            Application.OpenURL(PlatformRouteController.GetCreateIP());
        EditorGUILayout.EndHorizontal();

        EditorGUILayout.BeginHorizontal();
        if (GUILayout.Button("輸出渠道資料", GUILayout.Height(30)))
            EditorApplication.delayCall += Export;
        if (GUILayout.Button("輸出所有渠道資料", GUILayout.Height(30)))
            EditorApplication.delayCall += ExportAll;

        EditorGUILayout.EndHorizontal();
        EditorGUILayout.Foldout(m_FoldoutOpen, "Route List - count : " + m_FilterList.Count);
        if (m_FoldoutOpen == true)
        {
            m_ScrollVeiwRouteList = EditorGUILayout.BeginScrollView(m_ScrollVeiwRouteList);
            EditorGUILayout.BeginVertical(GUI.skin.box);
            if (m_FilterList.Count == 0)
                EditorGUILayout.HelpBox("Not found.", MessageType.Warning);
            else
            {
                for (int i = 0; i < m_FilterList.Count; ++i)
                {
                    EditorGUILayout.BeginHorizontal(GUI.skin.box);
                    if (m_PublishList.Contains(m_FilterList[i]))
                        EditorGUILayout.TextField("Publish", GUI.skin.box, GUILayout.Width(50), GUILayout.Height(20));
                    else if (m_BuildInfo != null && m_BuildInfo.GetAppKey() == m_FilterList[i].GetAppKey())
                        EditorGUILayout.TextField("Current", GUI.skin.box, GUILayout.Width(50), GUILayout.Height(20));
                    EditorGUILayout.BeginVertical();
                    m_FoldoutFilterList[i] = EditorGUILayout.Foldout(m_FoldoutFilterList[i], m_FilterList[i].GetKeyWord());
                    if (m_FoldoutFilterList[i] == true)
                    {
                        EditorGUILayout.LabelField("CompanyName", m_FilterList[i].CompanyName);
                        EditorGUILayout.LabelField("ProductName", m_FilterList[i].ProductName);
                        EditorGUILayout.LabelField("BundleID", m_FilterList[i].BundleID);
                        EditorGUILayout.LabelField("AppName", m_FilterList[i].AppName);
                        EditorGUILayout.LabelField("AppleID", m_FilterList[i].AppleID);
                        EditorGUILayout.LabelField("ServerOption", m_FilterList[i].ServerOption.ToString());
                        EditorGUILayout.LabelField("Branch", m_FilterList[i].Branch.ToString());
                        EditorGUILayout.LabelField("PlatformID", m_FilterList[i].PlatformID.ToString());
                        EditorGUILayout.LabelField("SiteID", m_FilterList[i].SiteID.ToString());
                        EditorGUILayout.LabelField("StaticSiteID", m_FilterList[i].StaticSiteID.ToString());
                        EditorGUILayout.LabelField("FlagShare", m_FilterList[i].FlagShare ? "V" : "X");
                        EditorGUILayout.LabelField("IconName", m_FilterList[i].IconName);

                        EditorGUILayout.BeginHorizontal();
                        if (GUILayout.Button("Edit ( 編輯 )", GUILayout.Width(150)))
                            Application.OpenURL(PlatformRouteController.GetEditorIP(m_FilterList[i].GetAppKey()));
                        if (m_CustomDefine != null)
                        {
                            if (GUILayout.Button("Set Current", GUILayout.Width(150)))
                            {
                                SettingBuild(m_FilterList[i]);
                                
                                m_FoldoutCurrent = true;
                                m_FoldoutFilterList[i] = false;
                            }

                            #if !UNITY_IOS
                            if (GUILayout.Button("Add Publish List", GUILayout.Width(150)))
                            {
                                m_PublishList.Add(m_FilterList[i]);
                                m_FoldoutPublishList.Add(false);

                                m_FoldoutPublish = true;
                                m_FoldoutFilterList[i] = false;
                            }
                            #endif
                        }
                        EditorGUILayout.EndHorizontal();
                    }
                    EditorGUILayout.EndVertical();
                    EditorGUILayout.EndHorizontal();
                }
            }

            EditorGUILayout.EndVertical();
            EditorGUILayout.EndScrollView();
        }
        EditorGUILayout.EndVertical();

        EditorGUILayout.EndHorizontal();
    }

    IEnumerator BuildList()
    {
        m_PublishList.Sort((x, y) => { return y.IconName.CompareTo(x.IconName); });
        for (int i = 0; i < m_PublishList.Count; ++i)
        {
            SettingBuild(m_PublishList[i]);
            CopyImage();
            AssetDatabase.Refresh();
            Build();
            yield return new WaitForSeconds(1);
        }
    }

    private void ChangeFilter()
    {
        m_FoldoutFilterList.Clear();
        m_FilterList.Clear();
        for (int i = 0; i < m_RouteList.Count; ++i)
        {
            // 開發
            if (m_FilterServerOption == 1 && m_RouteList[i].ServerOption != 0)
                continue;
            // QA
            if (m_FilterServerOption == 2 && m_RouteList[i].ServerOption != 1)
                continue;
            // Stage
            if (m_FilterServerOption == 3 && m_RouteList[i].ServerOption != 5)
                continue;
            // 美服
            if (m_FilterServerOption == 4 && m_RouteList[i].ServerOption != 3)
                continue;
            // 國服
            if (m_FilterServerOption == 5 && m_RouteList[i].ServerOption != 4)
                continue;

            if (string.IsNullOrEmpty(m_FilterCompanyName) == false)
                if (m_RouteList[i].CompanyName.IndexOf(m_FilterCompanyName) < 0)
                    continue;

            if (string.IsNullOrEmpty(m_FilterProductName) == false)
                if (m_RouteList[i].ProductName.IndexOf(m_FilterProductName) < 0)
                    continue;

            if (string.IsNullOrEmpty(m_FilterBundleID) == false)
                if (m_RouteList[i].BundleID.IndexOf(m_FilterBundleID) < 0)
                    continue;

            if (string.IsNullOrEmpty(m_FilterAppName) == false)
                if (m_RouteList[i].AppName.IndexOf(m_FilterAppName) < 0)
                    continue;

            if (string.IsNullOrEmpty(m_FilterPlatformID) == false)
                if (m_RouteList[i].PlatformID.ToString().IndexOf(m_FilterPlatformID) < 0)
                    continue;

            if (string.IsNullOrEmpty(m_FilterStaticSiteID) == false)
                if (m_RouteList[i].StaticSiteID.ToString().IndexOf(m_FilterStaticSiteID) < 0)
                    continue;

            if (CheckFilterContent(m_RouteList[i]) == true)
            {
                m_FilterList.Add(m_RouteList[i]);
                m_FoldoutFilterList.Add(false);
            }
        }
    }

    private bool CheckFilterContent(BuildInfo _data)
    {
        if (string.IsNullOrEmpty(m_FilterContent) == true)
            return true;

        if (_data.Index.ToString().IndexOf(m_FilterContent) >= 0)
            return true;

        if (_data.CompanyName.IndexOf(m_FilterContent) >= 0)
            return true;

        if (_data.ProductName.IndexOf(m_FilterContent) >= 0)
            return true;

        if (_data.BundleID.IndexOf(m_FilterContent) >= 0)
            return true;

        if (_data.AppName.IndexOf(m_FilterContent) >= 0)
            return true;

        if (_data.AppleID.IndexOf(m_FilterContent) >= 0)
            return true;

        if (_data.PlatformID.ToString().IndexOf(m_FilterContent) >= 0)
            return true;

        if (_data.SiteID.ToString().IndexOf(m_FilterContent) >= 0)
            return true;

        if (_data.StaticSiteID.ToString().IndexOf(m_FilterContent) >= 0)
            return true;

        if (_data.IconName != null && _data.IconName.IndexOf(m_FilterContent) >= 0)
            return true;

        return false;
    }

    private BuildInfo GetBuildInfoByID(string _bundleID, int _staticSiteID)
    {
        for (int i = 0; i < m_RouteList.Count; ++i)
        {
            if (_bundleID == m_RouteList[i].BundleID)
            {
                if (_staticSiteID == m_RouteList[i].StaticSiteID)
                {
                    return m_RouteList[i];
                }
            }
        }
        return null;
    }

    private static string GetBranchName(int Branch)
    {
        if (Branch == 0)
            return "方舟";
        else if (Branch == 1)
            return "986";
        else if (Branch == 2)
            return "蓝钻";
        else
            return Branch.ToString();
    }

    private static int CheckArkData()
    {
        if (m_BuildInfo.Branch != s_Branch)
        {
            string _Message = "ArkGame和渠道資料不符\n";
            _Message += "SiteID " + m_BuildInfo.SiteID + "\n";
            _Message += "ArkGame " + GetBranchName(s_Branch) + "\n";
            _Message += "渠道 " + GetBranchName(m_BuildInfo.Branch);

            Debug.LogError(_Message);
            EditorUtility.DisplayDialog("通知", _Message, "確定");
            return -1;
        }
        else
        {
            return 0;
        }
    }

    private static int CheckHallData()
    {
        string filePath;
#if UNITY_IOS
        filePath = Application.dataPath + "/StreamingAssets/iOS/Hall_BuildLog.txt";
#else
        filePath = Application.dataPath + "/StreamingAssets/Android/Hall_BuildLog.txt";
#endif

        if (example.invalid.File.Exists(filePath) == false)
        {
            EditorUtility.DisplayDialog("通知", "大廳熱更資料不存在", "確定");
            return -1;
        }

        FileStream fileStream_Open = new FileStream(Path.GetFullPath(filePath), FileMode.Open);
        StreamReader streamWriter = new StreamReader(fileStream_Open);
        string oldScriptContent = streamWriter.ReadToEnd();
        fileStream_Open.Close();
        streamWriter.Close();

        string _NoteBegin = "Branch ";
        string _NoteEnd = "Version";
        int begin = oldScriptContent.IndexOf(_NoteBegin);
        int end = oldScriptContent.IndexOf(_NoteEnd);
        if (begin == -1 || end == -1)
        {
            Debug.LogWarning("找不到Branch資料");
            return 0;
        }
        int Length = end - (begin + _NoteBegin.Length);

        oldScriptContent = oldScriptContent.Substring(begin + _NoteBegin.Length, Length);
        int result = Int32.Parse(oldScriptContent);
        if (result != s_Branch)
        {
            string _Message = "ArkGame和大廳熱更資料不符";
            if (m_BuildInfo != null)
                _Message += " SiteID " + m_BuildInfo.SiteID;

            _Message += "\n";
            _Message += "ArkGame " + GetBranchName(s_Branch) + "\n";
            _Message += "大廳熱更 " + GetBranchName(result);

            Debug.LogError(_Message);
            EditorUtility.DisplayDialog("通知", _Message, "確定");
            return -1;
        }
        else
        {
            return 1;
        }
    }

    private void Build()
    {
        // 檢查Ark和渠道
        if (CheckArkData() == -1)
            return;

        // 檢查Ark和大廳
        // 0:找不到 -1:錯誤 1正確
        if (CheckHallData() == -1)
            return;

        if (!CheckFishTag())
        {
            EditorUtility.DisplayDialog("通知", "捕魚標籤掉了請重新設定標籤再打包", "確定");
            return;
        }

        if (!CheckFishWarp())
        {
            EditorUtility.DisplayDialog("通知", "捕魚Warp掉了，請重新生成Warp再打包", "確定");
            return;
        }

#if UNITY_IOS
        SettingIOS();
        BuildIOS();
#else
        SettingAndroid();
        ChangeXml();
        BuildAndroid();
#endif
        Debug.Log("finished! custom build");
    }
    private bool CheckFishTag()
    {
        string search_tag="fish";
       
         for (int i = 0; i < UnityEditorInternal.InternalEditorUtility.tags.Length; i++) {
             if (UnityEditorInternal.InternalEditorUtility.tags[i].Contains(search_tag)) {
                 Debug.Log ("At Position " + i + " is the Tag " + search_tag + " found :) ");
                return true;
             }
         }
        Debug.Log("No tag be found :) ");
        return false;
    }
    private bool CheckFishWarp()
    {
        if (File.Exists(Application.dataPath + "/FrameWork/XLua/Gen/UnityEngineTextMeshWrap.cs") == false)
            return false;

        if (File.Exists(Application.dataPath + "/FrameWork/XLua/Gen/mogenFishWrap.cs") == false)
            return false;

        if (File.Exists(Application.dataPath + "/FrameWork/XLua/Gen/mogenBulletWrap.cs") == false)
            return false;

        if (File.Exists(Application.dataPath + "/FrameWork/XLua/Gen/mogenFishMsgWrap.cs") == false)
            return false;

        return true;        
    }
    private void SettingBuild(BuildInfo _info)
    {
        m_BuildInfo = _info;
        if (_info == null)
            return;

        // 修改設定
        PlayerSettings.companyName = m_BuildInfo.CompanyName;
        PlayerSettings.productName = m_BuildInfo.ProductName;
        PlayerSettings.SplashScreen.show = false;

        PlayerSettings.SetApplicationIdentifier(BuildTargetGroup.Android, m_BuildInfo.BundleID);
        PlayerSettings.SetApplicationIdentifier(BuildTargetGroup.iOS, m_BuildInfo.BundleID);
        PlayerSettings.SetApplicationIdentifier(BuildTargetGroup.Standalone, m_BuildInfo.BundleID);

        m_CustomDefine.ServerOption = m_BuildInfo.ServerOption;
        m_CustomDefine.PlatformID = m_BuildInfo.PlatformID;
        m_CustomDefine.StaticSiteID = m_BuildInfo.StaticSiteID;

        EditorSceneManager.SaveScene(EditorSceneManager.GetActiveScene());        
    }

    private void SettingAndroid()
    {
        PlayerSettings.SetApplicationIdentifier(BuildTargetGroup.Android, m_BuildInfo.BundleID);
        PlayerSettings.SetApplicationIdentifier(BuildTargetGroup.Standalone, m_BuildInfo.BundleID);

        bool _checkLog = (m_CustomDefine.ServerOption == 3) || (m_CustomDefine.ServerOption == 4);
        string _defines = _checkLog ? "NOLOG" : "";
        PlayerSettings.SetScriptingDefineSymbolsForGroup(BuildTargetGroup.Android, _defines);

        int _length = Application.dataPath.LastIndexOf('/');
        PlayerSettings.Android.keystoreName = Application.dataPath.Substring(0, _length) + "/app.keystore";
        PlayerSettings.Android.keystorePass = "123456";
        PlayerSettings.Android.keyaliasName = "cqxx";
        PlayerSettings.Android.keyaliasPass = "123456";
    }

    private void SettingIOS()
    {
        PlayerSettings.SetApplicationIdentifier(BuildTargetGroup.iOS, m_BuildInfo.BundleID);

        //string _defines = (m_CustomSetting.m_Define.ServerOption > 2) ? "NOLOG" : "";
        //PlayerSettings.SetScriptingDefineSymbolsForGroup(BuildTargetGroup.iOS, _defines);
    }

    private void BuildAndroid()
    {
        string _name = m_BuildInfo.AppName;
        _name += "_" + GetRandomString(4);//加入四碼亂數字串
        _name += "_p" + m_CustomDefine.PlatformID;
        _name += "_s" + m_CustomDefine.StaticSiteID.ToString("00");
        _name += "_op" + m_CustomDefine.ServerOption;
        _name += "_v" + System.DateTime.Now.ToString("yyMMdd");

        string _path = Path.Combine(System.Environment.GetFolderPath(System.Environment.SpecialFolder.Desktop), System.DateTime.Now.ToString("yyMMdd"));
        _path = Path.Combine(_path, _name);

        int _index = 1;
        while (File.Exists(_path + _index.ToString("00") + ".apk"))
            ++_index;
        _path += _index.ToString("00") + ".apk";
        BuildPlayerOptions buildPlayerOptions = new BuildPlayerOptions();
        buildPlayerOptions.locationPathName = _path;
        buildPlayerOptions.target = BuildTarget.Android;
        buildPlayerOptions.options = BuildOptions.None;
        buildPlayerOptions.scenes = GetSceneList();
        BuildPipeline.BuildPlayer(buildPlayerOptions);
        Debug.Log(_path);
    }
    
    //產生包含英文大小寫的亂數字串
    public static string GetRandomString(int length)
    {
        var str = "0123456789abcdefghijklmnopqrstuvwxyz";
        System.Random next = new System.Random();
        var builder = new StringBuilder();
        for (var i = 0; i < length; i++)
        {
            builder.Append(str[next.Next(0, str.Length)]);
        }
        return builder.ToString();
    }

    private void BuildIOS()
    {
        BuildPlayerOptions buildPlayerOptions = new BuildPlayerOptions();
        buildPlayerOptions.locationPathName = m_BuildInfo.AppName;
        buildPlayerOptions.target = BuildTarget.iOS;
        buildPlayerOptions.options = BuildOptions.None;
        buildPlayerOptions.scenes = GetSceneList();
        BuildPipeline.BuildPlayer(buildPlayerOptions);
    }

    private string[] GetSceneList()
    {
        if (m_BuildInfo.BundleID == "com.cqxx.myGameserver")
            return new[] { "Assets/ServerIP.unity", "Assets/Start.unity" };

        return new[] { "Assets/Start.unity" };
    }

    public void CopyImage()
    {
        string _folder = m_BuildInfo.IconName;
        if (string.IsNullOrEmpty(_folder) == true)
            _folder = m_BuildInfo.CompanyName.ToString();

        // splash
        string _pathSource = Application.dataPath + "/../Icon/" + _folder + "/";
        string _pathDest = Application.dataPath + "/Image/";
        string _fileName = "SplashImage.png";
        FileCopy(_pathSource + _fileName, _pathDest + _fileName);

        // logo
        _pathSource = Application.dataPath + "/../Icon/" + _folder + "/";
        _pathDest = Application.streamingAssetsPath + "/";
        _fileName = "Logo.png";
        FileCopy(_pathSource + _fileName, _pathDest + _fileName);

        // hall logo
        _pathSource = Application.dataPath + "/../Icon/" + _folder + "/";
        _pathDest = Application.streamingAssetsPath + "/";
        _fileName = "HallLogo.png";

        if (File.Exists(_pathDest + _fileName))
            File.Delete(_pathDest + _fileName);
        else
            Debug.LogWarning("hall logo Dest is null");

        if (File.Exists(_pathSource + _fileName))
            FileCopy(_pathSource + _fileName, _pathDest + _fileName);
        else
            Debug.LogWarning("hall logo Source is null");

#if UNITY_IOS
        // ios icon
        _pathSource = Application.dataPath + "/../Icon/" + _folder + "/ios/";
        _pathDest = Application.dataPath + "/Image/Icon/";
        string[] picList = Directory.GetFiles(_pathSource, "*.png");
        foreach (string f in picList)
        {
            _fileName = f.Substring(_pathSource.Length);
            FileCopy(f, _pathDest + _fileName);
        }
#else
        // android icon
        _pathSource = Application.dataPath + "/../Icon/" + _folder + "/android/";
        _pathDest = Application.dataPath + "/Plugins/Android/res/";
        _fileName = "app_icon.png";
        FileCopy(_pathSource + "Icon.png", _pathDest + "drawable/" + _fileName);
        FileCopy(_pathSource + "Icon-hdpi.png", _pathDest + "drawable-hdpi/" + _fileName);
        FileCopy(_pathSource + "Icon-ldpi.png", _pathDest + "drawable-ldpi/" + _fileName);
        FileCopy(_pathSource + "Icon-mdpi.png", _pathDest + "drawable-mdpi/" + _fileName);
        FileCopy(_pathSource + "Icon-xhdpi.png", _pathDest + "drawable-xhdpi/" + _fileName);
        FileCopy(_pathSource + "Icon-xxhdpi.png", _pathDest + "drawable-xxhdpi/" + _fileName);
        FileCopy(_pathSource + "Icon-xxxhdpi.png", _pathDest + "drawable-xxxhdpi/" + _fileName);
#endif
    }
    private void FileCopy(string _pathSource, string _pathDest)
    {
        if (File.Exists(_pathSource))
        {
            File.Delete(_pathDest);
            File.Copy(_pathSource, _pathDest, true);
        }
        else
            Debug.Log("FileNotFoundException: " + _pathSource + " does not exist");
    }

    private void ChangeXml()
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

                string[] _lineSeparator = new string[] { "\n" };
                string[] _lines = _content.Split(_lineSeparator, System.StringSplitOptions.None);
                string _keyWord = "android:host=\"";
                for (int i = 0; i < _lines.Length; ++i)
                {
                    _lines[i] = _lines[i].Replace("\r", "");
                    if (_lines[i].IndexOf(_keyWord) >= 0)
                    {
                        _lines[i] = _lines[i].Substring(0, _lines[i].IndexOf(_keyWord));
                        _lines[i] += _keyWord + m_BuildInfo.BundleID + "\"";
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


    private string ExportData(BuildInfo data)
    {
        Dictionary<string, object> _data = new Dictionary<string, object>();

        Dictionary<string, object> _systemConfig = new Dictionary<string, object>();
        _data.Add("systemConfig", _systemConfig);

        Dictionary<string, object> _Info = new Dictionary<string, object>();
        _Info.Add("BundleID", data.BundleID);
        _Info.Add("AppName", data.AppName);
        if (data.AppleID != "nil")
            _Info.Add("AppleID", data.AppleID);
        else
            _Info.Add("AppleID", "");
        _Info.Add("ServerOption", data.ServerOption);
        _Info.Add("Branch", data.Branch);
        _Info.Add("PlatformID", data.PlatformID);
        _Info.Add("SiteID", data.SiteID);
        _Info.Add("StaticSiteID", data.StaticSiteID);
        _Info.Add("Agent", data.Agent);
        _Info.Add("Template", "0");
        _Info.Add("FlagExchage", 0);
        _Info.Add("ExchangeFontType", data.Version);
        _systemConfig.Add("Info", _Info);

        Dictionary<string, JiuLiaoInfo> m_JiuLiaoList = PlatformRouteController.GetJiuLiaoList();
        Dictionary<string, object> _JiuLiao = new Dictionary<string, object>();
        string Key = data.BundleID.Replace('.', '_');
        if (m_JiuLiaoList.ContainsKey(Key) == true)
        {
            JiuLiaoInfo _JiuLiaoData = m_JiuLiaoList[Key];
            _JiuLiao.Add("JiuLiaoId", _JiuLiaoData.JiuLiaoId);
            _JiuLiao.Add("JiuLiaoSecrety", _JiuLiaoData.JiuLiaoSecrety);
            _JiuLiao.Add("JiuLiaoName", _JiuLiaoData.JiuLiaoName);
            _JiuLiao.Add("DownloadUrl", _JiuLiaoData.DownloadUrl);
        }
        else
        {
            _JiuLiao.Add("JiuLiaoId", "");
            _JiuLiao.Add("JiuLiaoSecrety", "");
            _JiuLiao.Add("JiuLiaoName", "");
            _JiuLiao.Add("DownloadUrl", "");
        }
        _systemConfig.Add("JiuLiao", _JiuLiao);

        Dictionary<string, object> _payConfig = new Dictionary<string, object>();
        _payConfig.Add("UIFlagPay", true);
        _payConfig.Add("UIFlagAgent", true);
        _payConfig.Add("UIFlagApplyAgent", false);
        if (_JiuLiao["JiuLiaoId"].ToString() != "")
            _payConfig.Add("UIFlagJiuLiao", true);
        else
            _payConfig.Add("UIFlagJiuLiao", false);
        _systemConfig.Add("payConfig", _payConfig);

        string config = MiniJSON.Json.Serialize(_data);
        return config;
    }

    private void Export()
    {
        Dictionary<string, object> Data = new Dictionary<string, object>();
        BuildInfo _info = m_BuildInfo.Copy();
        string config = ExportData(_info);
        Debug.Log(_info.SiteID + "," + config);
    }

    private void ExportAll()
    {
        string _content = "";

        Dictionary<int, BuildInfo> test = new Dictionary<int, BuildInfo>();

        List<BuildInfo> _Pool = new List<BuildInfo>();
        for (int index = 0; index < m_RouteList.Count; index++)
        {
            if (m_RouteList[index].NotExport == 0)
                _Pool.Add(m_RouteList[index]);
            else
                continue;

            if (test.ContainsKey(m_RouteList[index].SiteID) == false)
                test.Add(m_RouteList[index].SiteID, m_RouteList[index]);
            else
                Debug.LogWarning("@@@@@@@ " + m_RouteList[index].SiteID);
        }

        _Pool.Sort((x, y) => {
            int compare = x.SiteID.CompareTo(y.SiteID);
            if (compare != 0)
                return compare;

            return x.AppName.CompareTo(y.AppName);
        });


        Dictionary<string, object> Data = new Dictionary<string, object>();
        for (int index = 0; index < _Pool.Count; index++)
        {
            _content += _Pool[index].SiteID + ",";
            _content += ExportData(_Pool[index]);
            _content += "\n";
        }

        string _targetPath = Application.dataPath;
        File.WriteAllText(_targetPath + "JsonInfo.txt", _content);
        EditorUtility.DisplayDialog("通知", "資料JsonInfo.txt輸出完成!!!", "確定");

        EditorUtility.RevealInFinder(_targetPath);
    }
}
