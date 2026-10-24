using UnityEngine;

public class RuntimeStats : MonoBehaviour
{
    public static RuntimeStats s_instance;

    public bool m_IsShow = false;

    public void Awake()
    {
        s_instance = this;
    }

    private float fpsMeasuringDelta = 2.0f;

    private float timePassed;
    private int m_FrameCount = 0;
    private float m_FPS = 0.0f;

    private void Start()
    {
        timePassed = 0.0f;

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

    private void Update()
    {
        m_FrameCount = m_FrameCount + 1;
        timePassed = timePassed + Time.deltaTime;

        if (timePassed > fpsMeasuringDelta)
        {
            m_FPS = m_FrameCount / timePassed;

            timePassed = 0.0f;
            m_FrameCount = 0;
        }
    }
     
    private void OnGUI()
    {
        GUIStyle guiStyle = new GUIStyle();
        guiStyle.normal.background = null;                          //这是设置背景填充的 
        guiStyle.fontSize = 30;                                     //当然，这是字体大小

        string labText = string.Empty;
        if (m_FPS > 25)
        {
            labText = string.Format("<color=#00FF00FF>{0}</color>", "FPS: " + m_FPS.ToString("F1"));
        }
        else if (m_FPS > 20)
        {
            labText = string.Format("<color=#ffff00ff>{0}</color>", "FPS: " + m_FPS.ToString("F1"));
        }
        else
        {
            labText = string.Format("<color=#ff0000ff>{0}</color>", "FPS: " + m_FPS.ToString("F1"));
        }

        //居中显示FPS
        GUI.Label(new Rect(20, 10, 100, 100), labText, guiStyle);
    }
}