using System;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.EventSystems;

public class RuntimeDebug : MonoBehaviour
{
    public static RuntimeDebug s_instance;

    public void Awake()
    {
        s_instance = this;
    }

    public bool m_IsShow = false;

    private const int LOGLENGTH = 300;
    private static bool m_ShowLog = false;
    private static bool m_IsLock = false;

    private static List<string> logList = new List<string>();


    int m_TakeStackLineNo = 5;

    public void Start()
    {
        Application.logMessageReceived += Application_logMessageReceived;

        this.enabled = false;
    }

    public void SetShow(bool IsShow)
    {
#if UNITY_EDITOR
        this.enabled = m_IsShow;
#else
        this.enabled = IsShow;
#endif
    }

    private void Application_logMessageReceived(string condition, string stackTrace, LogType type)
    {
        if (type == LogType.Warning)
            return;

        switch (type)
        {
            case LogType.Error:
            case LogType.Exception:
                AppendToLogList(string.Format("<color=#ff0000ff>{0}  {1}</color>", DateTime.Now.ToString("hh:mm:ss"), condition));

                string[] subList = stackTrace.Split(new String[] { "\n" }, StringSplitOptions.None);

                if (subList.Length < 5)
                    m_TakeStackLineNo = subList.Length;

                for (int i = 0; i < m_TakeStackLineNo; i++)
                    AppendToLogList(string.Format("<color=#ff0000ff>{0}</color>", subList[i]));
                break;
            default:
                AppendToLogList(string.Format("{0}  {1}", DateTime.Now.ToString("hh:mm:ss"), condition));
                break;
        }
    }

    static string logString = string.Empty;
    static Vector2 scrollPosition = new Vector2(0, 0);
    static bool m_IsLogListChanged = false;             //  用于判断Log是否需要刷新，如果不需要不必每帧都去重新拼字符串，重新拼字符串效率比较低

    private static void AppendToLogList(string log)
    {
        string[] subList = log.Split(new String[] { "\n" }, StringSplitOptions.None);

        for (int i = 0; i < subList.Length; i++)
        {
            if (logList.Count >= LOGLENGTH)
                logList.RemoveAt(0);

            logList.Add(subList[i] + "\n");
        }

        if (!m_IsLock)
            scrollPosition = new Vector2(0, int.MaxValue);

        m_IsLogListChanged = true;
    }

    void OnGUI()
    {
        GUILayout.BeginHorizontal();

        GUILayout.Space(Screen.width / 2 - 60);

        if (GUILayout.Button(m_ShowLog ? "隐藏Log" : "显示Log", GUILayout.Width(120), GUILayout.Height(60)))
            m_ShowLog = !m_ShowLog;

        if (m_ShowLog)
        {
            if (GUILayout.Button("清空Log", GUILayout.Width(120), GUILayout.Height(60)))
            {
                logList.Clear();
                logString = string.Empty;
            }

            if (!m_IsLock && m_IsLogListChanged)
            {
                m_IsLogListChanged = false;
                logString = string.Empty;
                logList.ForEach((o) => { logString += o; });
            }

            if (GUILayout.Button("复制Log", GUILayout.Width(120), GUILayout.Height(60)))
                Utility.CopyText(logString);

            if (GUILayout.Button(m_IsLock ? "√ 锁定Log " : "锁定Log", GUILayout.Width(120), GUILayout.Height(60)))
            {
                if (!m_IsLock)
                    m_IsLogListChanged = true;

                m_IsLock = !m_IsLock;
            }
            GUILayout.EndHorizontal();

            if (!string.IsNullOrEmpty(logString))
            {
                GUILayout.BeginHorizontal();
                GUILayout.Space(10);
                scrollPosition = GUILayout.BeginScrollView(scrollPosition, GUILayout.Width(Screen.width - 10), GUILayout.Height(Screen.height - 80));

                GUIStyle guiStyle = new GUIStyle();
                guiStyle.normal.background = null;                          //这是设置背景填充的 
                guiStyle.fontSize = 20;                                     //当然，这是字体大小
                guiStyle.normal.textColor = new Color(1, 1, 1);

                GUILayout.Label(logString, guiStyle);
                GUILayout.EndScrollView();
                GUILayout.EndHorizontal();
            }
        }
        else
        {
            GUILayout.EndHorizontal();
        }
    }

    private float m_StartMouseY = 0;
    private Vector2 m_StartScrollPosition = Vector2.zero;
    private bool m_IsMouseDown = false;
    private bool m_isMouseMove = false;

    public void Update()
    {
        if (!m_ShowLog)
            return;

        if (Input.mousePosition.x > Screen.width / 2 - 200f)
        {
            m_IsMouseDown = false;
            m_isMouseMove = false;
            return;
        }

        if (Input.GetMouseButtonDown(0))
        {
            m_StartMouseY = Input.mousePosition.y;
            m_StartScrollPosition = scrollPosition;
            m_IsMouseDown = true;
        }

        if (Input.GetMouseButtonUp(0))
        {
            m_IsMouseDown = false;
            m_isMouseMove = false;
        }

        if (m_IsMouseDown && Math.Abs(Input.mousePosition.y - m_StartMouseY) > 20)
        {
            m_isMouseMove = true;
        }

        if (m_isMouseMove)
        {
            float moveLength = Input.mousePosition.y - m_StartMouseY;
            scrollPosition = Vector2.Lerp(scrollPosition, new Vector2(scrollPosition.x, m_StartScrollPosition.y + moveLength * 1f), 50f * Time.deltaTime);
        }

    }
}
