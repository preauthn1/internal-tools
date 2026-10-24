using System;
using System.Collections;
using System.Collections.Generic;
using System.Linq;
using UnityEngine;
using UnityEngine.SceneManagement;

public class AssetBundleOption
{

#if UNITY_EDITOR
    static int m_SimulateAssetBundleInEditor = -1;
    static string kSimulateAssetBundles = "SimulateAssetBundles";
    /// <summary>
    /// 编辑器模式下是否模拟资源包加载方式？（不需要真正打出资源包，避免每改动一个资源都打一次包的过程，抓取StreamingAssets中資料）
    /// </summary>
    public static bool SimulateAssetBundleInEditor
    {
        get
        {
            if (m_SimulateAssetBundleInEditor == -1)
            {
                m_SimulateAssetBundleInEditor = UnityEditor.EditorPrefs.GetBool(kSimulateAssetBundles, true) ? 1 : 0;
            }

            return m_SimulateAssetBundleInEditor != 0;
        }
        set
        {
            int newValue = value ? 1 : 0;
            if (newValue != m_SimulateAssetBundleInEditor)
            {
                m_SimulateAssetBundleInEditor = newValue;
                UnityEditor.EditorPrefs.SetBool(kSimulateAssetBundles, value);
            }
        }
    }

    // Add By Rayon
    static int m_SimulateAssetBundleHotUpdateInEditor = -1;
    static string kSimulateAssetBundlesHotUpdate = "SimulateAssetBundlesHotUpdate";
    /// <summary>
    /// 编辑器模式下是否模拟资源包加载方式？（網路上抓取熱更資料）
    /// </summary>
    public static bool SimulateAssetBundleHotUpdateInEditor
    {
        get
        {
            if (m_SimulateAssetBundleHotUpdateInEditor == -1)
            {
                m_SimulateAssetBundleHotUpdateInEditor = UnityEditor.EditorPrefs.GetBool(kSimulateAssetBundlesHotUpdate, true) ? 1 : 0;
            }

            return m_SimulateAssetBundleHotUpdateInEditor != 0;
        }
        set
        {
            int newValue = value ? 1 : 0;
            if (newValue != m_SimulateAssetBundleHotUpdateInEditor)
            {
                m_SimulateAssetBundleHotUpdateInEditor = newValue;
                UnityEditor.EditorPrefs.SetBool(kSimulateAssetBundlesHotUpdate, value);
            }
        }
    }

#endif
}

/// <summary>
/// AssetBundle 资源管理器
/// </summary>
public partial class AssetBundleManager : Singleton<AssetBundleManager>
{
    #region Manager Interface

    public bool IsReady { private set; get; }

    Dictionary<string, ProjectAssetBundleManager> m_ProjectAssetBundleManagerDic = new Dictionary<string, ProjectAssetBundleManager>();

    public void ReLaunch()
    {
        IsReady = false;

        foreach (var item in m_ProjectAssetBundleManagerDic)
            item.Value.Release();

        m_ProjectAssetBundleManagerDic.Clear();

        Launch();
    }

    public void Launch()
    {
        RegistProjectAssetBundleManager(AppDefine.HallNumber);
        IsReady = true;
    }

    public void RegistProjectAssetBundleManager(string projectNo)
    {
        if (m_ProjectAssetBundleManagerDic.ContainsKey(projectNo))
            unRegistProjectAssetBundleManager(projectNo);

        ProjectAssetBundleManager manager = gameObject.AddComponent<ProjectAssetBundleManager>();
        manager.Init(projectNo);
        m_ProjectAssetBundleManagerDic.Add(projectNo, manager);
    }

    public void unRegistProjectAssetBundleManager(string projectName)
    {
        m_ProjectAssetBundleManagerDic[projectName].Release();
        Destroy(m_ProjectAssetBundleManagerDic[projectName]);
        m_ProjectAssetBundleManagerDic.Remove(projectName);
    }

    public T LoadAsset<T>(string projectName, string asset, bool unloadAssetBundle = true, bool unloadAllLoadedObjects = false) where T : UnityEngine.Object
    {
        if (!IsReady || !m_ProjectAssetBundleManagerDic.ContainsKey(projectName))
            return null;

        return m_ProjectAssetBundleManagerDic[projectName].LoadAsset<T>(asset, unloadAssetBundle, unloadAllLoadedObjects);
    }

    public GameObject LoadGameObject(string projectName, string asset, bool unloadAssetBundle = true, bool unloadAllLoadedObjects = false)
    {
        return LoadAsset<GameObject>(projectName, asset, unloadAssetBundle, unloadAllLoadedObjects);
    }

    public byte[] LoadByte(string projectName, string asset, bool unloadAssetBundle = true, bool unloadAllLoadedObjects = false)
    {
        TextAsset textAsset = LoadAsset<TextAsset>(projectName, asset, unloadAssetBundle, unloadAllLoadedObjects);

        if (textAsset == null)
            return null;

        return textAsset.bytes;
    }

    public LoadAssetAsyncOperation LoadAssetAsync<T>(string projectName, string assetName, bool unloadAssetBundle = true, bool unloadAllLoadedObjects = false) where T : UnityEngine.Object
    {
        if (!IsReady || !m_ProjectAssetBundleManagerDic.ContainsKey(projectName))
            return null;

        return m_ProjectAssetBundleManagerDic[projectName].LoadAssetAsync<T>(assetName, unloadAssetBundle, unloadAllLoadedObjects);
    }

    public AsyncOperation LoadSceneAsync(string projectName, string scenceName)
    {
        // 以完整方式載入場景 By Rayon
        return m_ProjectAssetBundleManagerDic[projectName].LoadSceneAsync(scenceName, true, false, projectName);
    }

    #endregion
}

public class ProjectAssetBundleManager : MonoBehaviour
{
    private string m_ProjectName;       // 当前资源管理器对应的项目名称

    public void Release()
    {
        StopAllCoroutines();
        UnloadAllAssetBundles(true);
        m_AssetRecordDict.Clear();
        this.MainManifest = null;
    }

    public void Init(string projectName)
    {
        m_ProjectName = projectName;

        // 加载MainManifest            
        LoadLocalAssetBundleMainfest();
        // 加载资源记录文件
        LoadAssetRecord();
        // 加载预加载的ab 及相关
        if (this.MainManifest != null)
        {
            string[] assetBundleNames = this.MainManifest.GetAllAssetBundles();
            for (int i = 0; i < assetBundleNames.Length; i++)
            {
                string assetBundleName = assetBundleNames[i];
                if (assetBundleName.EndsWith("_pre"))
                {
                    LoadAssetBundleAndDependencies(assetBundleName, true);
                    DisposeTemporaryAssetBundles(false, false);
                }
            }
        }
    }

    /// <summary>
    /// Asset bundle 资源清单
    /// </summary>
    private AssetBundleManifest MainManifest { get; set; }

    /// <summary>
    /// 资源信息池（资源名，资源所在的资源包名）
    /// </summary>
    private Dictionary<string, AssetRecordInfo> m_AssetRecordDict = new Dictionary<string, AssetRecordInfo>();

    /// <summary>
    /// 常驻的AssetBundle
    /// </summary>
    private Dictionary<string, AssetBundle> m_PermanentAssetBundleDict = new Dictionary<string, AssetBundle>();

    /// <summary>
    /// 缓存的AssetBundle
    /// </summary>
    private Dictionary<string, AssetBundle> m_CacheAssetBundleDict = new Dictionary<string, AssetBundle>();

    /// <summary>
    /// 临时的AssetBundle
    /// </summary>
    private Dictionary<string, AssetBundle> m_TemporaryAssetBundleDict = new Dictionary<string, AssetBundle>();

    enum AssetBundleDictType
    {
        Permanent = 1,
        Cache,
        Temporary,
    }

    #region Load Asset

    /// <summary>
    /// 同步加载一个资源
    /// </summary>
    /// <typeparam name="T"></typeparam>
    /// <param name="sceneName">资源名称，不关注大小写</param>
    /// <param name="unloadAssetBundle"></param>
    /// <param name="unloadAllLoadedObjects"></param>
    /// <param name="projectName"></param>
    /// <returns></returns>
    public AsyncOperation LoadSceneAsync(string sceneName, bool unloadAssetBundle = true, bool unloadAllLoadedObjects = false, string projectName = "")
    {
        try
        {
            AssetRecordInfo assetRecord = GetAssetRecordInfoByAsset(sceneName);
            if (assetRecord == null)
                return null;
#if UNITY_EDITOR
            if (AssetBundleOption.SimulateAssetBundleInEditor)
            {
                // By Rayon
                // 以完整方式載入場景
                // 有多个名称相同但是路径不同的场景，你应该使用全部路径
                if (string.IsNullOrEmpty(projectName) == false)
                    return SceneManager.LoadSceneAsync("Game/" + projectName + "/" + sceneName, LoadSceneMode.Additive);
                else
                    return SceneManager.LoadSceneAsync(sceneName, LoadSceneMode.Additive);
            }
            else
#endif
            {
                AssetBundle assetBundle = LoadAssetBundleAndDependenciesByAsset(assetRecord);

                if (assetBundle != null)
                { 
                    Debug.Log("框架：场景资源加载完成，开始加载场景");

                    AsyncOperation async = SceneManager.LoadSceneAsync(sceneName, LoadSceneMode.Additive);
                    // DisposeTemporaryAssetBundles(unloadAssetBundle, unloadAllLoadedObjects);
                    return async;
                }
            }

            return null;
        }
        catch (System.Exception ex)
        {
            Debug.LogErrorFormat("Project {0} AssetBundleManager.LoadAsset is falid!\n{1}", m_ProjectName, ex.Message);
            return null;
        }
    }

    /// <summary>
    /// 同步加载一个资源
    /// </summary>
    /// <typeparam name="T"></typeparam>
    /// <param name="asset">资源名称，不关注大小写</param>
    /// <param name="unloadAssetBundle"></param>
    /// <param name="unloadAllLoadedObjects"></param>
    /// <returns></returns>
    public T LoadAsset<T>(string asset, bool unloadAssetBundle = true, bool unloadAllLoadedObjects = false) where T : UnityEngine.Object
    {
        try
        {
            T result = null;
            AssetRecordInfo assetRecord = GetAssetRecordInfoByAsset(asset);
            if (assetRecord == null)
                return null;
#if UNITY_EDITOR
            if (AssetBundleOption.SimulateAssetBundleInEditor)
            {
                result = UnityEditor.AssetDatabase.LoadAssetAtPath<T>(assetRecord.RelativePath);
            }
            else
#endif
            {
                AssetBundle assetBundle = LoadAssetBundleAndDependenciesByAsset(assetRecord);
                if (assetBundle != null)
                {
                    result = assetBundle.LoadAsset<T>(asset);
                    DisposeTemporaryAssetBundles(unloadAssetBundle, unloadAllLoadedObjects);
                }
            }
            return result;
        }
        catch (System.Exception ex)
        {
            Debug.LogErrorFormat("Project {0} AssetBundleManager.LoadAsset is falid!\n{1}", m_ProjectName, ex.Message);
        }
        return null;
    }

    /// <summary>
    /// 异步加载一个资源
    /// </summary>
    /// <typeparam name="T"></typeparam>
    /// <param name="assetName"></param>
    /// <param name="unloadAssetBundle"></param>
    /// <param name="unloadAllLoadedObjects"></param>
    /// <returns></returns>
    public LoadAssetAsyncOperation LoadAssetAsync<T>(string assetName, bool unloadAssetBundle = true, bool unloadAllLoadedObjects = false) where T : UnityEngine.Object
    {
        try
        {
            LoadAssetAsyncOperation operation = null;
            AssetRecordInfo assetRecord = GetAssetRecordInfoByAsset(assetName);
            if (assetRecord == null)
                return null;

#if UNITY_EDITOR
            if (AssetBundleOption.SimulateAssetBundleInEditor)
            {
                T asset = UnityEditor.AssetDatabase.LoadAssetAtPath<T>(assetRecord.RelativePath);
                operation = new LoadAssetAsyncOperation(asset);
            }
            else
#endif
            {
                AssetBundle assetBundle = LoadAssetBundleAndDependenciesByAsset(assetRecord);
                AssetBundleRequest request = null;
                if (assetBundle != null)
                {
                    request = assetBundle.LoadAssetAsync<T>(assetRecord.AssetName);
                    DisposeTemporaryAssetBundles(unloadAssetBundle, unloadAllLoadedObjects);
                }
                operation = new LoadAssetAsyncOperation(request);
            }
            return operation;
        }
        catch (System.Exception ex)
        {
            Debug.LogErrorFormat("Project {0} AssetBundleManager.LoadAsset is falid!\n{1}", m_ProjectName, ex.Message);
        }
        return null;
    }

    /// <summary>
    /// 加载本地的Mainfest 文件
    /// </summary>
    void LoadLocalAssetBundleMainfest()
    {
#if UNITY_EDITOR
        if (AssetBundleOption.SimulateAssetBundleInEditor)// 模拟运行时不需要加载本地Mainfest 文件
        {
            return;
        }
#endif
        AssetBundle assetBundle = LoadAssetBundle(m_ProjectName, false);
        if (assetBundle != null)
        {
            this.MainManifest = assetBundle.LoadAsset<AssetBundleManifest>("AssetBundleManifest");
            assetBundle.Unload(false);
        }
        if (this.MainManifest == null)
        {
            Debug.LogErrorFormat("Project {0} LocalAssetBundleMainfest: the MainManifest was null, please check and fix it!", m_ProjectName);
        }
    }

    /// <summary>
    /// 通过 资源名 获取AssetBundle名称
    /// </summary>
    /// <param name="assetName">资源名</param>
    /// <returns></returns>
    string GetAssetBundleNameByAsset(string assetName)
    {
        if (m_AssetRecordDict.ContainsKey(assetName))
        {
            return m_AssetRecordDict[assetName].AssetBundleName;
        }
        else
        {
            return m_AssetRecordDict[assetName].RelativePath;
        }
    }

    /// <summary>
    /// 获取资源记录信息
    /// </summary>
    /// <param name="assetRecordKey"></param>
    /// <returns></returns>
    AssetRecordInfo GetAssetRecordInfoByAsset(string assetRecordKey)
    {
        if (m_AssetRecordDict.ContainsKey(assetRecordKey))
        {
            return m_AssetRecordDict[assetRecordKey];
        }
        else
        {
            Debug.LogErrorFormat("Project {0} Can't find asset record key[{1}] in AssetRecords, please check it!", m_ProjectName, assetRecordKey);
            return null;
        }
    }

    /// <summary>
    /// 加载AB 包 并加载所有的依赖项
    /// </summary>
    /// <param name="assetName"></param>
    /// <returns></returns>
    AssetBundle LoadAssetBundleAndDependenciesByAsset(AssetRecordInfo assetRecord)
    {
        if (assetRecord == null)
            return null;
        // 加载所有的依赖项
        return LoadAssetBundleAndDependencies(assetRecord.AssetBundleName, false);
    }

    /// <summary>
    /// 加载AssetBundle和相关的依赖项
    /// </summary>
    /// <param name="assetBundleName"></param>
    /// <returns></returns>
    AssetBundle LoadAssetBundleAndDependencies(string assetBundleName, bool isPreLoadAssetBundle)
    {
        if (MainManifest == null)
            return null;

        string[] dependencies = MainManifest.GetAllDependencies(assetBundleName.ToLower());

        for (int index = 0; index < dependencies.Length; index++)
        {
            if (LoadAssetBundle(dependencies[index], isPreLoadAssetBundle) == null)
            {
                Debug.LogErrorFormat("Project {0} -> {1}'s Dependencie AssetBundle can't find. Name is {2}, please check it!", m_ProjectName, assetBundleName, dependencies[index]);
                return null;
            }
        }
        return LoadAssetBundle(assetBundleName, isPreLoadAssetBundle);
    }

    /// <summary>
    /// 加载 AssetBundle
    /// </summary>
    /// <param name="assetBundleName">AssetBundle 名称</param>
    /// <param name="isPreLoadAssetBundle">是否预先加载的assetBundle</param>
    /// <returns></returns>
    AssetBundle LoadAssetBundle(string assetBundleName, bool isPreLoadAssetBundle)
    {
        AssetBundle assetBundle = null;
        // 判断常驻资源字典中是否存在
        if (m_PermanentAssetBundleDict.ContainsKey(assetBundleName))
        {
            assetBundle = m_PermanentAssetBundleDict[assetBundleName];
            if (assetBundle != null)
                return assetBundle;
            else
                m_PermanentAssetBundleDict.Remove(assetBundleName);
        }
        // 判断缓存资源字典中是否存在
        if (m_CacheAssetBundleDict.ContainsKey(assetBundleName))
        {
            assetBundle = m_CacheAssetBundleDict[assetBundleName];
            if (assetBundle != null)
                return assetBundle;
            else
                m_CacheAssetBundleDict.Remove(assetBundleName);
        }
        //判断临时资源字典中是否存在
        if (m_TemporaryAssetBundleDict.ContainsKey(assetBundleName))
        {
            assetBundle = m_TemporaryAssetBundleDict[assetBundleName];
            if (assetBundle != null)
                return assetBundle;
            else
                m_TemporaryAssetBundleDict.Remove(assetBundleName);
        }

        assetBundle = InternalLoadAssetBundleFromFile(assetBundleName);

        if (assetBundle != null)
        {
            // 将资源添加到不同的字典中
            if (IsPermanentAssetBundle(assetBundleName) || isPreLoadAssetBundle)
            {
                m_PermanentAssetBundleDict.Add(assetBundleName, assetBundle);
            }
            else
            {
                m_TemporaryAssetBundleDict.Add(assetBundleName, assetBundle);
            }
        }

        return assetBundle;
    }

    /// <summary>
    /// 是否是不销毁资源
    /// </summary>
    /// <param name="assetBundleName"></param>
    /// <returns></returns>
    bool IsPermanentAssetBundle(string assetBundleName)
    {
        return false;
    }

    /// <summary>
    /// 从文件加载 AssetBundle 文件
    /// </summary>
    /// <param name="assetBundleName"></param>
    /// <returns></returns>
    AssetBundle InternalLoadAssetBundleFromFile(string assetBundleName)
    {
        AssetBundle assetBundle = null;
        string assetBundleFullPath = FileUtility.GetDataFileFullPath(m_ProjectName, assetBundleName);
        byte[] fileContent = FileUtility.ReadFileBytes(assetBundleFullPath);
        if (fileContent != null)
        {
            if (AppDefine.IsEncrypted)
            {
                fileContent = Utility.AESDecrypt(fileContent, AppDefine.AssetSecretKey);
            }

            assetBundle = AssetBundle.LoadFromMemory(fileContent);
        }

        return assetBundle;
    }

    #endregion

    #region Unload AssetBundle

    /// <summary>
    /// 释放所有的AssetBundle 资源
    /// </summary>
    /// <param name="unloadAllLoadedObjects"></param>
    void UnloadAllAssetBundles(bool unloadAllLoadedObjects)
    {
        UnloadAssetBundles(unloadAllLoadedObjects, AssetBundleDictType.Temporary);
        UnloadAssetBundles(unloadAllLoadedObjects, AssetBundleDictType.Cache);
        UnloadAssetBundles(unloadAllLoadedObjects, AssetBundleDictType.Permanent);
    }

    /// <summary>
    /// 释放对应字典中的 AssetBundle 资源
    /// </summary>
    /// <param name="unloadAllLoadedObjects"></param>
    /// <param name="dictType"></param>
    void UnloadAssetBundles(bool unloadAllLoadedObjects, AssetBundleDictType dictType)
    {
        Dictionary<string, AssetBundle> dict = null;
        switch (dictType)
        {
            case AssetBundleDictType.Permanent:
                dict = m_PermanentAssetBundleDict;
                break;
            case AssetBundleDictType.Cache:
                dict = m_CacheAssetBundleDict;
                break;
            case AssetBundleDictType.Temporary:
                dict = m_TemporaryAssetBundleDict;
                break;
            default:
                break;
        }
        if (dict == null)
            return;

        var enumerator = dict.Values.GetEnumerator();
        while (enumerator.MoveNext())
        {
            if (enumerator.Current != null)
            {
                enumerator.Current.Unload(unloadAllLoadedObjects);
            }
        }
        enumerator.Dispose();
        dict.Clear();
    }

    /// <summary>
    /// 释放掉AssetBundle
    /// </summary>
    /// <param name="assetBundleName"></param>
    /// <param name="unloadAllLoadedObjects"></param>
    void UnloadAssetBundle(string assetBundleName, bool unloadAllLoadedObjects)
    {
        AssetBundle ab = null;
        if (m_PermanentAssetBundleDict.TryGetValue(assetBundleName, out ab))
            m_PermanentAssetBundleDict.Remove(assetBundleName);
        else if (m_CacheAssetBundleDict.TryGetValue(assetBundleName, out ab))
            m_CacheAssetBundleDict.Remove(assetBundleName);
        else if (m_TemporaryAssetBundleDict.TryGetValue(assetBundleName, out ab))
            m_TemporaryAssetBundleDict.Remove(assetBundleName);

        if (ab != null)
        {
            ab.Unload(unloadAllLoadedObjects);
        }
    }

    /// <summary>
    /// 保存临时资源到缓存资源字典
    /// </summary>
    void SaveTemporaryAssetBundlesToCache()
    {
        var enumerator = m_TemporaryAssetBundleDict.GetEnumerator();
        while (enumerator.MoveNext())
        {
            m_CacheAssetBundleDict.Add(enumerator.Current.Key, enumerator.Current.Value);
        }
        enumerator.Dispose();
        m_TemporaryAssetBundleDict.Clear();
    }

    /// <summary>
    /// 释放临时资源
    /// </summary>
    /// <param name="unloadAssetBundle">是否释放资源，不释放，则放入缓存字典中</param>
    /// <param name="unloadAllLoadedObjects"></param>
    void DisposeTemporaryAssetBundles(bool unloadAssetBundle, bool unloadAllLoadedObjects)
    {
        if (unloadAssetBundle)
            UnloadAssetBundles(unloadAllLoadedObjects, AssetBundleDictType.Temporary);
        else
            SaveTemporaryAssetBundlesToCache();
    }

    #endregion

    #region Asset Record

    /// <summary>
    /// 加载资源记录文件
    /// </summary>
    void LoadAssetRecord()
    {
#if UNITY_EDITOR
        if (AssetBundleOption.SimulateAssetBundleInEditor)
        {
            m_AssetRecordDict = new Dictionary<string, AssetRecordInfo>();
            GetRecordList("Assets/" + (m_ProjectName != AppDefine.HallNumber ? ("Game/" + m_ProjectName) : m_ProjectName), ref m_AssetRecordDict);
        }
        else
#endif
        {
            byte[] contentBytes = FileUtility.ReadFileBytes(FileUtility.GetDataFileFullPath(m_ProjectName, AppDefine.AssetRecordsFileName));

            if (AppDefine.IsEncrypted)
            {
                contentBytes = Utility.AESDecrypt(contentBytes, AppDefine.AssetSecretKey);
            }

            string fileContent = Utility.BytesToUTF8String(contentBytes);
            if (string.IsNullOrEmpty(fileContent))
            {
                Debug.LogErrorFormat("Project {0} Load Asset Reocrds file error, please check and fix it!", m_ProjectName);
                return;
            }

            string[] assetRecordInfoArray = fileContent.Split(new string[] { "#$" }, StringSplitOptions.RemoveEmptyEntries);
            for (int index = 0; index < assetRecordInfoArray.Length; index++)
            {
                string[] recordItems = assetRecordInfoArray[index].Split(new string[] { "#@" }, StringSplitOptions.None);
                if (recordItems.Length > 1)
                {
                    string assetRecordKey = recordItems[0].Trim();
                    string assetBundleName = recordItems[1].Trim();
                    if (m_AssetRecordDict.ContainsKey(assetRecordKey))
                    {
                        Debug.LogErrorFormat("Project {0}  Some asset with the same assetrecordkey[{1}], please check and fix it! in assetBundleName:\n{2}\n{3}", m_ProjectName, assetRecordKey, m_AssetRecordDict[assetRecordKey].AssetBundleName, assetBundleName);
                    }
                    else
                    {
                        m_AssetRecordDict.Add(assetRecordKey, new AssetRecordInfo(assetRecordKey, assetBundleName));
                    }
                }
            }
        }
    }

    public static void GetRecordList(string path, ref Dictionary<string, AssetRecordInfo> recordList)
    {
        string[] files = example.invalid.Directory.GetFiles(path).Where(t => !t.EndsWith(".meta") && !t.EndsWith(".cs")).ToArray();

        if (files.Length != 0)
        {
            for (int i = 0; i < files.Length; i++)
                files[i] = files[i].Replace('\\', '/');

            for (int i = 0; i < files.Length; i++)
            {
                string filePath = files[i];

                if (filePath.EndsWith(".prefab")// 预制件
                    || filePath.EndsWith(".bytes") || filePath.EndsWith(".txt") || filePath.EndsWith(".xml") // 配置文件或者Lua文件
                    || filePath.EndsWith(".mp3") || filePath.EndsWith(".wav") || filePath.EndsWith(".ogg") // 音效文件
                    || filePath.EndsWith(".unity")
                    || filePath.EndsWith(".ttf")
                    || filePath.EndsWith(".lua")// 场景文件
                   )
                {
                    string assetRecordKey = AssetRecordInfo.ToAssetRecordKey(filePath) + (filePath.EndsWith(".lua") ? ".lua" : string.Empty);

                    if (recordList.ContainsKey(assetRecordKey))
                        Debug.LogErrorFormat("Some asset with the same assetrecordkey[{0}], please check and fix it! in assetBundleName:\n{1}\n{2}", assetRecordKey, recordList[assetRecordKey].AssetBundleName, filePath);
                    else
                        recordList.Add(assetRecordKey, new AssetRecordInfo(assetRecordKey, string.Empty, filePath));
                }
            }
        }

        string[] subFolders = example.invalid.Directory.GetDirectories(path);

        if (subFolders != null && subFolders.Length != 0)
        {
            for (int i = 0; i < subFolders.Length; i++)
                GetRecordList(subFolders[i], ref recordList);
        }
    }

    #endregion
}

public class AssetRecordInfo
{
    public AssetRecordInfo(string assetKey, string assetBundleName)
        : this(assetKey, assetBundleName, null)
    {
    }

    public AssetRecordInfo(string assetKey, string assetBundleName, string relativePath)
    {
        this.AssetName = assetKey;
        this.AssetBundleName = assetBundleName;
        this.RelativePath = relativePath;
    }

    /// <summary>
    /// 资源名称
    /// </summary>
    public string AssetName { get; private set; }

    /// <summary>
    /// 包名称
    /// </summary>
    public string AssetBundleName { get; private set; }

    /// <summary>
    /// 相对路径
    /// </summary>
    public string RelativePath { get; private set; }

    public static string ToAssetRecordKey(string fileName)
    {
        fileName = fileName.Replace('\\', '/');
        int lastIndex = fileName.LastIndexOf('/');
        string recordKey = string.Empty;
        if (lastIndex > -1)
        {
            recordKey = fileName.Substring(lastIndex + 1);
        }
        else
        {
            recordKey = fileName;
        }

        if (recordKey.Contains("."))
        {
            recordKey = recordKey.Substring(0, recordKey.LastIndexOf('.'));
        }

        return recordKey;
    }
}