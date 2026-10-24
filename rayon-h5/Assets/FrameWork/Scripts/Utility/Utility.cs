using UnityEngine;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using example.invalid;
using example.invalid.Sockets;
using UnityEngine.UI;
using UnityEngine.Events;
using XLua;
using DG.Tweening.Plugins.Core.PathCore;
using DG.Tweening.Plugins.Options;
using DG.Tweening;
using DG.Tweening.Core;
using QRCoder;

/// <summary>
/// 常用方法类
/// </summary>
public static class Utility
{
    public static string GetGameServerIP(string url)
    {
        IPAddress[] address;
        try
        {
            address = Dns.GetHostAddresses(url);
        }
        catch (Exception e)
        {
            Debug.LogException(e);
            return "";
        }

        if (address == null)
        {
            Debug.LogError("根据url获取不到ip地址");
            return "";
        }
        foreach (var info in address)
        {
            Debug.Log(info.ToString() + "----" + info.AddressFamily.ToString());
        }

        if (address[0].AddressFamily == AddressFamily.InterNetworkV6)
        {
            Debug.Log("Connect InterNetworkV6");
        }
        else
        {
            Debug.Log("Connect InterNetwork");
        }

        return address[0].ToString();
    }
    /// <summary>
    /// 当前的运行平台,1 windows, 2 android, 3 ios, 4 macos
    /// </summary>
    /// <returns></returns>
    public static int GetCurrentPlatform()
    {
#if UNITY_IOS && !UNITY_EDITOR
        return 3;
#endif

#if UNITY_ANDROID && !UNITY_EDITOR
        return 2;
#endif

#if UNITY_STANDLONE_OSX
        return 4;
#endif
        return 1;
    }
    /// <summary>
    /// 将任意格式字符串字符串转换成UTF8格式
    /// </summary>
    /// <param name="str">原字符串</param>
    /// <returns>UTF8格式的字符串</returns>
    public static string StringToUTF8String(string str)
    {
        byte[] data = Encoding.UTF8.GetBytes(str);
        return Encoding.UTF8.GetString(data);
    }

    /// <summary>
    /// 将任意格式字符串字符串转换成ANSI格式
    /// </summary>
    /// <param name="str">原字符串</param>
    /// <returns>ANSI格式的字符串</returns>
    public static string StringToANSIString(string str)
    {
        byte[] data = Encoding.Default.GetBytes(str);
        return Encoding.Default.GetString(data);
    }

    /// <summary>
    /// 将任意格式字符串字符串转换成Unicode格式
    /// </summary>
    /// <param name="str">原字符串</param>
    /// <returns>Unicode格式的字符串</returns>
    public static string StringToUnicodeString(string str)
    {
        byte[] data = Encoding.Unicode.GetBytes(str);
        return Encoding.Unicode.GetString(data);
    }

    /// <summary>
    /// 服务器的bytes 数据转换为 UTF-8 字符串
    /// </summary>
    /// <param name="data"></param>
    /// <returns>如果数据为空方法 string.Empty</returns>
    public static string BytesToUTF8String(byte[] data)
    {
        if (data == null || data.Length == 0)
            return string.Empty;
        return Encoding.UTF8.GetString(data);
    }

    public static byte[] UTF8StringToBytes(string str)
    {
        if (string.IsNullOrEmpty(str))
            return new byte[0];
        return Encoding.UTF8.GetBytes(str);
    }

    // 东八区时间
    static readonly DateTime s_StartTime = DateTime.Parse("1970-01-01 08:00:00");

    // 标准时间
    //static readonly DateTime s_StartTime = DateTime.Parse("1970-01-01 00:00:00");

    /// <summary>
    /// 服务器发送的时间戳转换为时间
    /// </summary>
    /// <param name="timestamp"></param>
    /// <returns></returns>
    public static DateTime UnixTimestampToDateTime(uint timestamp)
    {
        return s_StartTime.AddSeconds(timestamp);
    }

    /// <summary>
    /// DateTime时间格式转换为Unix时间戳格式
    /// </summary>
    /// <param name="time"> DateTime时间格式</param>
    /// <returns>Unix时间戳格式</returns>
    public static uint DateTimeToUnixTimestamp(System.DateTime time)
    {
        DateTime startTime = TimeZone.CurrentTimeZone.ToLocalTime(s_StartTime);
        return (uint)(time - startTime).TotalSeconds;
    }


    /// <summary>
    /// 设置物体自身节点及所有子节点的Layer
    /// </summary>
    /// <param name="root">物体根节点</param>
    /// <param name="layerName">Layer名</param>
    public static void SetGameObjectLayer(GameObject root, string layerName)
    {
        if (string.IsNullOrEmpty(layerName))
        {
            return;
        }
        int layer = LayerMask.NameToLayer(layerName);
        root.layer = layer;
        Transform[] children = root.GetComponentsInChildren<Transform>(true);
        for (int i = 0; i < children.Length; ++i)
        {
            children[i].gameObject.layer = layer;
        }
    }

    public static void SetRendererSortingOrder(GameObject root, int orderIndex)
    {
        Renderer[] renders = root.GetComponentsInChildren<Renderer>(true);

        foreach (Renderer render in renders)
        {
            render.sortingOrder = orderIndex;
        }
    }


    public static void SetUISortingOrder(GameObject ui, int orderOffSet = 1)
    {
        Transform transform = ui.transform;
        int orderIndex = 0;

        while (true)
        {
            if (transform == ui.transform.root)
                break;

            Canvas canvas = transform.GetComponent<Canvas>();

            if (canvas)
            {
                orderIndex = canvas.sortingOrder;
                break;
            }

            transform = transform.parent;
        }

        SetRendererSortingOrder(ui, orderIndex + orderOffSet);
    }


    /// <summary>
    /// 设置物体自身节点及所有子节点的Tag
    /// </summary>
    /// <param name="root">物体根节点</param>
    /// <param name="tag">Tag名</param>
    public static void SetGameObjectTag(GameObject root, string tag)
    {
        if (string.IsNullOrEmpty(tag))
        {
            return;
        }
        root.tag = tag;
        Transform[] children = root.GetComponentsInChildren<Transform>(true);
        for (int i = 0; i < children.Length; ++i)
        {
            children[i].tag = tag;
        }
    }

    /// <summary>
    /// 归于初始信息
    /// </summary>
    /// <param name="trans">需要归于初始的Transform</param>
    /// <param param name="parentTrans">父节点</param>
    public static void ReSetTransform(Transform trans, Transform parentTrans)
    {
        if (trans == null)
        {
            return;
        }

        if (parentTrans != null)
        {
            trans.SetParent(parentTrans, false);
        }

        trans.localPosition = Vector3.zero;
        trans.localRotation = Quaternion.identity;
        trans.localScale = Vector3.one;
    }

    /// <summary>
    /// 通过道路ID获取道路配置文件名称
    /// </summary>
    /// <param name="roadID"></param>
    /// <returns></returns>
    public static string GetRoadConfigFileNameByRoadID(uint roadID)
    {
        return string.Format("RoadConfig_{0}", roadID);
    }

    /// <summary>
    /// 将StreamingAssets目录下的文件的相对路径转换成绝对路径
    /// </summary>
    /// <param name="strAssetName">StreamingAssets目录下文件的相对路径</param>
    /// <returns>文件的绝对路径</returns>
    public static string GetFullPathOfStreamingAsset(string strRelativePath)
    {
        string strFullPath = "";
#if (UNITY_ANDROID || UNITY_WEBGL) && !UNITY_EDITOR
        strFullPath = example.invalid.Path.Combine(Application.streamingAssetsPath, strRelativePath);
#else
        strFullPath = "file://" + example.invalid.Path.GetFullPath(example.invalid.Path.Combine(Application.streamingAssetsPath, strRelativePath));
#endif
        return strFullPath;
    }

    public static string GetFullPathOfLoadDataPath(string strRelativePath)
    {
        string strFullPath = "";
#if UNITY_ANDROID && UNITY_EDITOR
        strFullPath = "jar:file://" + GetFullPathOfSaveDataPath(strRelativePath);
#else
        strFullPath = "file://" + GetFullPathOfSaveDataPath(strRelativePath);
#endif
        return strFullPath;
    }

    public static string GetFullPathOfSaveDataPath(string strRelativePath)
    {
        string strFullPath = "";
#if UNITY_ANDROID && !UNITY_EDITOR
        strFullPath = example.invalid.Path.GetFullPath(example.invalid.Path.Combine(Application.persistentDataPath, strRelativePath));
#else
        strFullPath = example.invalid.Path.GetFullPath(example.invalid.Path.Combine(Application.persistentDataPath, strRelativePath));
#endif
        return strFullPath;
    }

    /// <summary>
    /// Terrain层的Mask
    /// </summary>
    private static int terrainLayerMask = 1 << LayerMask.NameToLayer("Terrain");

    /// <summary>
    /// 获取到指定点下方的地面高度，如果下方无地面，那就获取到上方的地面高度
    /// </summary>
    /// <param name="x">x位置</param>
    /// <param name="y">y位置</param>
    /// <param name="z">z位置</param>
    /// <returns>地面高度（返回Mathf.NegativeInfinity表示未找到地面）</returns>
    public static float GetTerrainHeight(float x, float y, float z)
    {
        Vector3 pos = Vector3.zero;
        pos.x = x;
        pos.y = 9999;
        pos.z = z;

        RaycastHit hit;
        // 找下方地面
        if (Physics.Raycast(pos, -Vector3.up, out hit, 9999999.0f, terrainLayerMask))
        {
            return hit.point.y;
        }
        else
        {
            return Mathf.NegativeInfinity;
        }
    }

    public static GameObject InitGameObject(string transformName, Transform transform)
    {
        Transform findItem = transform.Find(transformName);
        if (findItem != null)
        {
            return findItem.gameObject;
        }
        else
        {
            Debug.LogError(string.Format("{0}组件未找到对象[{1}]", transform.name, transformName));
            return null;
        }
    }

    public static Transform InitTransform(string transformName, Transform transform)
    {
        Transform findItem = transform.Find(transformName);
        if (findItem != null)
        {
            return findItem;
        }
        else
        {
            Debug.LogError(string.Format("{0}组件未找到对象[{1}]", transform.name, transformName));
            return null;
        }
    }

    public static T InitComponent<T>(string transformName, Transform transform) where T : Component
    {
        Transform findItem = transform.Find(transformName);
        if (findItem != null)
        {
            T t = findItem.GetComponent<T>();
            if (t == null)
            {
                Debug.LogError(string.Format("{0}组件的对象[{1}]不包含组件", transform.name, transformName));
                return null;
            }
            else
            {
                return t;
            }
        }
        else
        {
            Debug.LogError(string.Format("{0}组件未找到对象[{1}]", transform.name, transformName));
            return null;
        }
    }

    /// <summary>
    /// 配置文件的时间转换为程序时间
    /// </summary>
    /// <param name="configTime">配置时间单位为(毫秒)</param>
    /// <returns>程序使用的时间为(秒)</returns>
    public static float TimeChange(uint configTime)
    {
        return configTime / 1000f;
    }

    /// <summary>
    /// 格式化邮件发送时间
    /// </summary>
    /// <param name="sendTime">邮件发送时间</param>
    /// <returns>界面显示的格式</returns>
    public static string FormatMailSendTime(DateTime sendTime)
    {
        DateTime dateTime = DateTime.Now.Date; // 获取今天的日期同步时间
        TimeSpan timeSpan = (dateTime - sendTime); // 时间差
        if (timeSpan.TotalDays > 0)
        {
            if (timeSpan.TotalDays > 2) // 超过2天的就算是2天以前的邮件
            {
                return string.Format("{0}天前", timeSpan.TotalDays);
            }
            else if (timeSpan.TotalDays > 1) // 超过一天，未超过2天
            {
                return "前天";
            }
            else // 昨天的邮件
            {
                return "昨天";
            }
        }
        else
        {
            // 如果出现服务器时间比客户端快很多的情况怎样显示
            return sendTime.ToString("HH:mm");
        }
    }

    public static int StringToInt(string str)
    {
        int rtnVal = 0;

        int.TryParse(str, out rtnVal);

        return rtnVal;
    }

    /// <summary>
    /// 求某一向量与Vector3.Forward所成的角度
    /// (一般用于求一个物体的绕Y轴的旋转角度(朝向 XZ = 0,1))
    /// --顺时针
    /// </summary>
    /// <param name="vector">向量</param>
    /// <returns>返回角度值:[0-360)</returns>
    public static float VectorToAngle(Vector3 vector)
    {
        float angle = Vector3.Angle(Vector3.forward, vector);
        if (vector.x < 0)
        {
            angle = 360 - angle;
        }
        return angle;
    }

    /// <summary>
    /// 与Vector3.Forward的角度转换为向量
    /// (一般用于求一个物体的绕Y轴的旋转角度(朝向 XZ = 0,1))
    /// --顺时针
    /// </summary>
    /// <param name="angleValue">与 Verctor.forward 的夹角</param>
    public static Vector3 AngleToVerctor(float angleValue)
    {
        // 转换成 [0~360)
        angleValue = angleValue % 360;
        if (angleValue < 0)
        {
            angleValue += 360;
        }

        float tanValue = Mathf.Tan(angleValue * Mathf.Deg2Rad);
        if (angleValue >= 90 && angleValue < 270)
        {
            return new Vector3(-tanValue, 0, -1);
        }
        else
        {
            return new Vector3(tanValue, 0, 1);
        }
    }

    // 克隆一个GameObject，返回克隆对象绑定脚本，Grid里面复制Item的时候使用
    public static T Clone<T>(this T template) where T : MonoBehaviour
    {
        if (!template)
            return null;

        GameObject instanceObj = GameObject.Instantiate(template.gameObject);
        instanceObj.transform.SetParent(template.gameObject.transform.parent);
        instanceObj.transform.localPosition = Vector3.zero;
        instanceObj.transform.localRotation = Quaternion.identity;
        instanceObj.transform.localScale = Vector3.one;
        instanceObj.SetActive(true);

        T script = instanceObj.GetComponent<T>();

        return script;
    }

    // 删除链表里面所有Gameobject
    public static void Destory<T>(this List<T> list) where T : MonoBehaviour
    {
        if (list == null || list.Count == 0)
            return;

        for (int i = 0; i < list.Count; i++)
            GameObject.Destroy(list[i].gameObject);
    }

    // 删除数组里面所有Gameobject
    public static void Destory<T>(this T[] list) where T : MonoBehaviour
    {
        if (list == null || list.Length == 0)
            return;

        for (int i = 0; i < list.Length; i++)
            GameObject.Destroy(list[i].gameObject);
    }

    /// <summary>
    /// AES加密
    /// </summary>
    /// <param name="encryptString">明文</param>
    /// <param name="key">密钥</param>
    /// <returns></returns>
    public static string AESEncrypt(string encryptString, string key)
    {
        byte[] toEncryptArray = Encoding.UTF8.GetBytes(encryptString);
        return Convert.ToBase64String(AESEncrypt(toEncryptArray, key));
    }

    public static byte[] AESEncrypt(byte[] encryptBytes, string key)
    {
        if (encryptBytes != null && encryptBytes.Length > 0)
        {
            byte[] keyArray = Encoding.UTF8.GetBytes(key);
            System.Security.Cryptography.RijndaelManaged rDel = new System.Security.Cryptography.RijndaelManaged();
            rDel.Key = keyArray;
            rDel.Mode = System.Security.Cryptography.CipherMode.ECB;
            rDel.Padding = System.Security.Cryptography.PaddingMode.PKCS7;
            System.Security.Cryptography.ICryptoTransform cTransform = rDel.CreateEncryptor();
            return cTransform.TransformFinalBlock(encryptBytes, 0, encryptBytes.Length);
        }
        else
        {
            return new byte[0];
        }
    }

    /// <summary>
    /// AES解密
    /// </summary>
    /// <param name="decryptStr">密文</param>
    /// <param name="key">密钥</param>
    /// <returns></returns>
    public static string AESDecrypt(string decryptStr, string key)
    {
        byte[] toEncryptArray = Convert.FromBase64String(decryptStr);
        return Encoding.UTF8.GetString(AESDecrypt(toEncryptArray, key));
    }

    /// <summary>
    /// AES 解密
    /// </summary>
    /// <param name="decryptBytes">密文</param>
    /// <param name="key">密钥</param>
    /// <returns></returns>
    public static byte[] AESDecrypt(byte[] decryptBytes, string key)
    {
        if (decryptBytes != null && decryptBytes.Length > 0)
        {
            byte[] keyArray = Encoding.UTF8.GetBytes(key);
            System.Security.Cryptography.RijndaelManaged rDel = new System.Security.Cryptography.RijndaelManaged();
            rDel.Key = keyArray;
            rDel.Mode = System.Security.Cryptography.CipherMode.ECB;
            rDel.Padding = System.Security.Cryptography.PaddingMode.PKCS7;
            System.Security.Cryptography.ICryptoTransform cTransform = rDel.CreateDecryptor();
            return cTransform.TransformFinalBlock(decryptBytes, 0, decryptBytes.Length);
        }
        else
        {
            return new byte[0];
        }
    }

    /// <summary>
    /// 计算字符串的长度--中文字符串算作2个字符长度
    /// </summary>
    /// <param name="value">字符串</param>
    /// <returns></returns>
    public static int StringLength(string value)
    {
        if (string.IsNullOrEmpty(value))
            return 0;
        int length = value.Length;
        System.Text.RegularExpressions.Regex chineseRegex = new System.Text.RegularExpressions.Regex("[\u4e00-\u9fa5]+");
        return length * 2 - chineseRegex.Replace(value, "").Length;
    }

    /// <summary>
    /// 计算字符串长度，中文字算一个字符
    /// </summary>
    /// <param name="text"></param>
    /// <returns></returns>
    public static int UTF8Stringlength(string text)
    {
        return text.Length;
    }

    public static string GetSubString(string content, int oriIndex, int length)
    {
        if (content.Length <= oriIndex)
        {
            return "";
        }
        if (content.Length < oriIndex + length)
        {
            length = content.Length - oriIndex;
        }

        return content.Substring(oriIndex, length);
    }

    /// <summary>
    /// 是否是邮箱地址
    /// </summary>
    /// <param name="emailAddress"></param>
    /// <returns></returns>
    public static bool IsEmailAddress(string emailAddress)
    {
        System.Text.RegularExpressions.Regex regex = new System.Text.RegularExpressions.Regex(@"^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$");
        return regex.IsMatch(emailAddress);
    }

    /// <summary>
    /// 获取角色模型资源名
    /// </summary>
    /// <param name="resourceID">资源ID--资源ID 从 1开始编写</param>
    /// <returns></returns>
    public static string GetPlayerResourceByID(uint resourceID)
    {
        if (resourceID == 0)
            resourceID = 1;
        return string.Format("Role{0}", resourceID.ToString().PadLeft(2, '0'));
    }

    /// <summary>
    /// 获取角色头像图片资源名
    /// </summary>
    /// <param name="iconID">头像图片id</param>
    /// <returns></returns>
    public static string GetPlayerHeadIconByID(uint iconID)
    {
        return string.Format("Role{0}", iconID.ToString().PadLeft(2, '0'));
    }

    /// <summary>
    /// 获取地图格子的key
    /// </summary>
    /// <param name="columnIndex">列索引</param>
    /// <param name="rowIndex">行索引</param>
    /// <returns></returns>
    public static string GetMapCellKey(object columnIndex, object rowIndex)
    {
        return string.Format("{0}_{1}", columnIndex, rowIndex);
    }

    public static void SetRectTransformHeight(RectTransform rectTrans, int height)
    {
        rectTrans.sizeDelta = new Vector2(rectTrans.sizeDelta.x, height);
    }

    public static void SetRectTransformWidthHight(RectTransform rectTrans, int height, int width)
    {
        rectTrans.sizeDelta = new Vector2(width, height);
    }

    public static int GetLogicAndValue(int value1, int value2)
    {
        return value1 & value2;
    }

    public static int GetLogicOrValue(int value1, int value2)
    {
        return value1 | value2;
    }

    public static void CallGC()
    {
        Resources.UnloadUnusedAssets();
        GC.Collect();
    }

    public static int GetLuaMod(int value, int value2)
    {
        return value % value2;
    }

    public static void ApplicationQuit()
    {
        Application.Quit();
    }

    /// <summary>
    /// 跳转场景
    /// </summary>
    /// <param name="sceneName"></param>
    public static void LoadScene(string sceneName)
    {
        UnityEngine.SceneManagement.SceneManager.LoadScene(sceneName);
        CallGC();
    }

    /// <summary>
    /// 获取内容的MD5
    /// </summary>
    /// <param name="content"></param>
    /// <returns></returns>
    public static string GetMD5Code(string content)
    {
        try
        {
            System.Security.Cryptography.MD5 md5 = new System.Security.Cryptography.MD5CryptoServiceProvider();
            byte[] retVal = md5.ComputeHash(System.Text.Encoding.UTF8.GetBytes(content));
            StringBuilder sb = new StringBuilder();
            for (int i = 0; i < retVal.Length; i++)
            {
                sb.Append(retVal[i].ToString("x2"));
            }
            return sb.ToString();
        }
        catch (Exception ex)
        {
            Debug.LogException(ex);
            return null;
        }
    }


    public static void DestoryAllChildren(this Transform trans, bool isIncludeInactive = true)
    {
        for (int i = 0; i < trans.childCount; i++)
        {
            if (isIncludeInactive == false && trans.GetChild(i).gameObject.activeSelf == false)
                continue;

            GameObject.Destroy(trans.GetChild(i).gameObject);
        }
    }

    public static void DestoryAllChildren(this GameObject gameObject, bool isIncludeInactive = true)
    {
        for (int i = 0; i < gameObject.transform.childCount; i++)
        {
            if (isIncludeInactive == false && gameObject.transform.GetChild(i).gameObject.activeSelf == false)
                continue;

            GameObject.Destroy(gameObject.transform.GetChild(i).gameObject);
        }
    }

    public static void AddClickListener(this Button button, UnityAction call)
    {
        ClickLimit limit = button.gameObject.GetComponent<ClickLimit>();

        if (!limit)
            limit = button.gameObject.AddComponent<ClickLimit>();

        limit.Init(button, call);
    }

    public static void CopyText(string text)
    {
        UniPasteBoard.SetClipBoardString(text);
    }

    public static LuaTable AddScriptComponent(this GameObject go, string projectName, string scriptName)
    {
        LuaBehaviour luaBehaviour = go.AddComponent<LuaBehaviour>();
        luaBehaviour.BindingLuaScript(projectName, scriptName);
        return luaBehaviour.LuaScript;
    }

    public static LuaTable GetScriptComponent(this GameObject go, string scriptName)
    {
        LuaBehaviour[] luaBehaviours = go.GetComponents<LuaBehaviour>();
        for (int i = 0; i < luaBehaviours.Length; i++)
        {
            if (luaBehaviours[i].ScriptName.Equals(scriptName))
                return luaBehaviours[i].LuaScript;
        }

        return null;
    }

    public static void RemoveScriptComponent(this GameObject go, string scriptName)
    {
        LuaBehaviour[] luaBehaviours = go.GetComponents<LuaBehaviour>();
        for (int i = 0; i < luaBehaviours.Length; i++)
        {
            if (luaBehaviours[i].ScriptName.Equals(scriptName))
                UnityEngine.Object.Destroy(luaBehaviours[i]);
        }
    }

    public static void RemoveAllScriptComponents(this GameObject go)
    {
        LuaBehaviour[] luaBehaviours = go.GetComponents<LuaBehaviour>();
        for (int i = 0; i < luaBehaviours.Length; i++)
        {
            UnityEngine.Object.Destroy(luaBehaviours[i]);
        }
    }

    public static int RangeInt(int min, int max)
    {
        return UnityEngine.Random.Range(min, max);
    }

    public static float RangeFloat(float min, float max)
    {
        return UnityEngine.Random.Range(min, max);
    }

    public static Tweener DOFramePath(this Transform trans, DOTweenPath path)
    {
        var tweener = trans.DOPath(path.wps.ToArray(), path.duration, path.pathType);
        tweener.SetUpdate(UpdateType.Manual);
        tweener.SetEase(path.easeType);
        tweener.SetLoops(path.loops);
        tweener.plugOptions.mode = PathMode.TopDown2D;
        tweener.plugOptions.orientType = path.orientType;

        switch (path.orientType)
        {
            case OrientType.None:
                break;
            case OrientType.ToPath:
                tweener.SetLookAt(0.01f);
                break;
            default:
                break;
        }

        return tweener;
    }

    public static Tweener DOFramePath(this Transform trans, List<Vector3> path, float duration, int loops = 1, OrientType orientType = OrientType.ToPath)
    {
        var tweener = trans.DOPath(path.ToArray(), duration, DG.Tweening.PathType.CatmullRom);
        tweener.SetUpdate(UpdateType.Manual);
        tweener.SetEase(Ease.Linear);
        tweener.SetLoops(loops);
        tweener.plugOptions.mode = PathMode.TopDown2D;

        tweener.plugOptions.orientType = orientType;

        switch (orientType)
        {
            case OrientType.None:
                break;
            case OrientType.ToPath:
                tweener.SetLookAt(0.01f);
                break;
            default:
                break;
        }

        return tweener;
    }

    // 获得所有路径信息的总长度
    public static float GetPathLength(List<Vector3> path)
    {
        float rtnVal = 0;

        for (int i = 1; i < path.Count; i++)
            rtnVal += Vector3.Distance(path[i], path[i - 1]);

        return rtnVal;
    }

// 取得不同解析度下的縮放，相對於16:9解析度
    public static float UISize = 1.218f;
    public static float GetScreenScale(float Correct)
    {
        if (Camera.main == null)
            return 0;

        RectTransform _RectTransform = Camera.main.gameObject.GetComponentInParent<RectTransform>();
        if (_RectTransform == null)
            return 0;

        float _Scale = Utility.UISize;
        float _RectScale = _RectTransform.transform.localScale.x;

        // _RectScale值過大
        if (_RectScale > UISize)
            return 1.0f;

        // 修正，避免4:3時使用Expand時計算錯誤
        if (_RectScale > Correct)
            _Scale /= _RectScale;

        return _Scale;
    }

	public static string GetbundleIdentifier()
	{
#if UNITY_WEBGL
        if (Main.Instance != null && Main.Instance.m_CustomDefine != null)
            return "";
        else
            return "";
#else
        return Application.identifier;
#endif
    }


    public static string GetUUID
    {
        get{
        //預設其他平台預設值
        string UUID = SystemInfo.deviceUniqueIdentifier;
#if UNITY_IOS && !UNITY_EDITOR
           UUID= GetUUidValue(KeyChain.BindGetKeyChainUser (), "uuid");
#endif

#if UNITY_ANDROID && !UNITY_EDITOR
            UUID=SystemInfo.deviceUniqueIdentifier;
#endif

        return UUID;
        }
    }
     
    public static string  GetUUidValue(string keyc, string key)
    {
        Dictionary<string, object>  kv =MiniJSON.Json.Deserialize(keyc) as Dictionary<string, object>;
        object valObj;
        kv.TryGetValue(key, out valObj);

        if (String.IsNullOrEmpty(valObj.ToString()))
        {
           string currentUUID = SystemInfo.deviceUniqueIdentifier;
#if UNITY_IOS && !UNITY_EDITOR
            KeyChain.BindSetKeyChainUser("0", currentUUID);
#endif
            return currentUUID;
        }
        else
            return valObj.ToString();
    }
    public static Sprite GenQRCode(string url)
    {
        var link = url;
        using (var qrGenerator = new QRCodeGenerator())
        {
            using (var qrCodeData = qrGenerator.CreateQrCode(link, QRCodeGenerator.ECCLevel.M))
            {
                var qrCodeImage = new UnityQRCode(qrCodeData).GetGraphic(20);
                return Sprite.Create(qrCodeImage, new Rect(0, 0, qrCodeImage.width, qrCodeImage.height), new Vector2(.5f, .5f));
            }
        }
    }
}