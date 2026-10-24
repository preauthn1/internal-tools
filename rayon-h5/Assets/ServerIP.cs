using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using UnityEngine.SceneManagement;

public class ServerIP : MonoBehaviour
{
    [SerializeField] Text m_TextServerIP;
	[SerializeField] InputField m_InputServerIP;
    [SerializeField] InputField m_InputServerPort;
    [SerializeField] Text m_TextUpdateDomainIP;
    [SerializeField] InputField m_InputUpdateDomainIP;
    [SerializeField] InputField m_InputUpdateDomainPort;
    [SerializeField] Text m_TextSiteID;
    [SerializeField] InputField m_InputSiteID;
    [SerializeField] Text m_TextPlatformID;
    [SerializeField] InputField m_InputPlatformID;
    [SerializeField] Text m_TextServerOption;
    [SerializeField] InputField m_InputServerOption;
    [SerializeField] Text m_TextChargeDomainIP;
    [SerializeField] InputField m_InputChargeDomainIP;
    [SerializeField] Text m_TextIPReuestIP;
    [SerializeField] InputField m_InputIPReuestIP;
    [SerializeField] GameObject m_AttachConfig;
    [SerializeField] GameObject m_PrefabConfig;
    [SerializeField] RectTransform m_ChooseServer;
    [SerializeField] Button m_ButtonEnter;
    [SerializeField] Button m_ButtonUnique;
    [SerializeField] GameObject m_ChooseUnique;
    [SerializeField] Text m_TextMessage;
    [SerializeField] Text m_TextServer;
    [SerializeField] Dropdown m_Dropdown;
    public List<Config> m_ConfigList = new List<Config>();
    private List<GameObject> m_PoolConfig = new List<GameObject>();

    /// <summary> 登入流程類型 </summary>
    private int m_LoginProcessType;
    private int m_ServerIndex;

    [System.Serializable]
    public class Config
    {
        public string m_Name;
        public string m_ServerListDomain;
        public string m_IP;
        public string m_WWWPort;
        public string m_ListPort;
        public string m_ServerPort;
        //public string m_UpdateDomainIP;
        //public string m_UpdateDomainPort;
        public string m_SiteID;
        public string m_PlatformID;
        public string m_ServerOption;
        //public string m_ChargeDomainIP;
        //public string m_IPReuestIP;
    }

    private void Start()
    {
        AppDefine.UseChangeServerIP = true;
        AppDefine.useDeviceUniqueIdentifier = (PlayerPrefs.GetInt("mUniqueIdentifier", 0) == 1);

        m_LoginProcessType = AppDefine.LoginProcessType;
        m_Dropdown.value = m_LoginProcessType;

        m_ButtonEnter.onClick.AddListener(ClickEnter);
        m_ButtonUnique.onClick.AddListener(ClickUnique);
        m_InputServerIP.onValueChanged.AddListener(ChangeConfig);
        m_InputServerPort.onValueChanged.AddListener(ChangeConfig);
        m_InputUpdateDomainIP.onValueChanged.AddListener(ChangeConfig);
        m_InputUpdateDomainPort.onValueChanged.AddListener(ChangeConfig);
        m_InputSiteID.onValueChanged.AddListener(ChangeConfig);
        m_InputPlatformID.onValueChanged.AddListener(ChangeConfig);
        m_InputServerOption.onValueChanged.AddListener(ChangeConfig);
        m_InputIPReuestIP.onValueChanged.AddListener(ChangeConfig);
        m_InputChargeDomainIP.onValueChanged.AddListener(ChangeConfig);
        m_Dropdown.onValueChanged.AddListener(ClickDropdown);

        if (Application.platform != RuntimePlatform.WindowsEditor)
        {
            if (Application.platform != RuntimePlatform.WindowsPlayer)
            {
                AppDefine.useDeviceUniqueIdentifier = false;
                m_ButtonUnique.gameObject.SetActive(false);
            }
        }
        ShowUnique();
        ShowDropdown();
        ShowConfig();
        ChangeConfig();
        ClickDropdown(m_LoginProcessType);

        if (m_ConfigList.Count == 1)
        {
            m_LoginProcessType = 2;
            ClickServer(0);
        }
    }

    /// <summary> 複製物件 </summary>
    private void CreateConfig()
    {
        GameObject clone = Instantiate(m_PrefabConfig) as GameObject;
        clone.transform.SetParent(m_AttachConfig.transform);
        clone.transform.localScale = Vector3.one;
        clone.gameObject.SetActive(true);
        m_PoolConfig.Add(clone);
    }

    private void ShowConfig()
    {
        for (int i = 0; i < m_ConfigList.Count; ++i)
        {
            if (i >= m_PoolConfig.Count)
                CreateConfig();

            GameObject clone = m_PoolConfig[i];
            Config _config = m_ConfigList[i];

            Text _name = clone.transform.Find("TextName").GetComponent<Text>();
            _name.text = _config.m_Name;

            Text _ip = clone.transform.Find("TextIP").GetComponent<Text>();
            _ip.text = _config.m_IP;
            
            Text _port = clone.transform.Find("TextPort").GetComponent<Text>();
            if (m_LoginProcessType == 0)
                _port.text = _config.m_WWWPort;
            else if (m_LoginProcessType == 1)
                _port.text = _config.m_ListPort;
            else if (m_LoginProcessType == 2)
                _port.text = _config.m_ServerPort;

            m_InputServerPort.gameObject.SetActive(m_LoginProcessType != 0);

            int _index = i;
            Button _button = clone.GetComponent<Button>();
            _button.onClick.AddListener(()=> { ClickServer(_index); });
        }
    }

    private void ClickServer(int _index)
    {
        Config _config = new Config();
        _config.m_ServerListDomain = AppDefine.DefaultServerIP;
        _config.m_IP = AppDefine.DefaultServerIP;
        _config.m_WWWPort = AppDefine.DefaultServerPort.ToString();
        _config.m_ListPort = AppDefine.DefaultServerPort.ToString();
        _config.m_ServerPort = AppDefine.DefaultServerPort.ToString();
        _config.m_SiteID = AppDefine.GameSiteID.ToString();
        _config.m_PlatformID = AppDefine.GamePlatformID.ToString();
        _config.m_ServerOption = AppDefine.ServerOption.ToString();

        if (_index >= 0)
            _config = m_ConfigList[_index];

        if (m_LoginProcessType == 0)
        {
            m_InputServerIP.text = _config.m_ServerListDomain;
            m_InputServerPort.text = _config.m_WWWPort;
            m_InputServerPort.gameObject.SetActive(false);
        }
        else if (m_LoginProcessType == 1)
        {
            m_InputServerIP.text = _config.m_IP;
            m_InputServerPort.text = _config.m_ListPort;
            m_InputServerPort.gameObject.SetActive(true);
        }
        else if (m_LoginProcessType == 2)
        {
            m_InputServerIP.text = _config.m_IP;
            m_InputServerPort.text = _config.m_ServerPort;
            m_InputServerPort.gameObject.SetActive(true);
        }
        //m_InputUpdateDomainIP.text = _config.m_UpdateDomainIP;
        //m_InputUpdateDomainPort.text = _config.m_UpdateDomainPort;
        m_InputSiteID.text = _config.m_SiteID;
        m_InputPlatformID.text = _config.m_PlatformID;
        m_InputServerOption.text = _config.m_ServerOption;
        //m_InputChargeDomainIP.text = _config.m_ChargeDomainIP;
        //m_InputIPReuestIP.text = _config.m_IPReuestIP;

        ChangeConfig();
    }
    
    public void ClickEnter()
    {
        m_TextMessage.text = m_TextServerIP.text +  "必須填值";
        if (string.IsNullOrEmpty(m_InputServerIP.text))
            return;

        m_TextMessage.text = m_TextServerIP.text + "的port必須填值";
        if (string.IsNullOrEmpty(m_InputServerPort.text))
            return;

        //m_TextMessage.text = m_TextUpdateDomainIP.text + "必須填值";
        //if (string.IsNullOrEmpty(m_InputUpdateDomainIP.text))
           // return;

        //m_TextMessage.text = m_TextUpdateDomainIP.text + "必須填值";
        //if (string.IsNullOrEmpty(m_InputUpdateDomainPort.text))
            //return;

        m_TextMessage.text = m_TextSiteID.text + "必須填值";
        if (string.IsNullOrEmpty(m_InputSiteID.text))
            return;

        m_TextMessage.text = m_TextPlatformID.text + "必須填值";
        if (string.IsNullOrEmpty(m_InputPlatformID.text))
            return;

        //m_TextMessage.text = m_TextChargeDomainIP.text + "必須填值";
        //if (string.IsNullOrEmpty(m_InputChargeDomainIP.text))
            //return;

        //m_TextMessage.text = m_TextIPReuestIP.text + "必須填值";
        //if (string.IsNullOrEmpty(m_InputIPReuestIP.text))
            //return;

        m_TextMessage.text = "";
        
        AppDefine.DefaultServerIP = m_InputServerIP.text;
        AppDefine.DefaultServerPort = int.Parse(m_InputServerPort.text);
        AppDefine.UpdateDomainIP = m_InputUpdateDomainIP.text;
        AppDefine.UpdateDomainPort = m_InputUpdateDomainPort.text;
        AppDefine.GameSiteID = int.Parse(m_InputSiteID.text);
        AppDefine.GamePlatformID = int.Parse(m_InputPlatformID.text);
        AppDefine.LoginProcessType = m_LoginProcessType;
        AppDefine.ChargeDomainIP = m_InputChargeDomainIP.text;
        AppDefine.IPReuestIP = m_InputIPReuestIP.text;
        AppDefine.ServerOption= int.Parse(m_InputServerOption.text);
        if (m_LoginProcessType == 0)
            AppDefine.ServerListDomain = AppDefine.DefaultServerIP;
        else if (m_LoginProcessType == 1)
            AppDefine.ServerListDomain = "http://" + AppDefine.DefaultServerIP + ":" + AppDefine.DefaultServerPort;
        else if (m_LoginProcessType == 2)
            AppDefine.ServerListDomain = "http://" + AppDefine.DefaultServerIP + ":" + AppDefine.DefaultServerPort;

        SceneManager.LoadScene("Start");
    }
    
    public void ShowDropdown()
    {
        m_Dropdown.value = m_LoginProcessType;
    }

    private void ClickDropdown(int _index)
    {
        m_LoginProcessType = _index;
        ShowConfig();
        ClickServer(m_ServerIndex);

        if (m_LoginProcessType == 0)
            m_TextServerIP.text = "伺服器列表文檔";
        else if (m_LoginProcessType == 1)
            m_TextServerIP.text = "伺服器列表網址";
        else if (m_LoginProcessType == 2)
            m_TextServerIP.text = "大廳伺服器網址";
    }

    private void ClickUnique()
    {
        AppDefine.useDeviceUniqueIdentifier = (AppDefine.useDeviceUniqueIdentifier == false);
        PlayerPrefs.SetInt("mUniqueIdentifier", (AppDefine.useDeviceUniqueIdentifier == true) ? 1 : 0);
        ShowUnique();
    }

    private void ShowUnique()
    {
        m_ChooseUnique.SetActive(AppDefine.useDeviceUniqueIdentifier);
    }
    
    private void ChangeConfig(string _result = "")
    {
        m_ChooseServer.gameObject.SetActive(false);
        for (int i = 0; i < m_ConfigList.Count; ++i)
        {
            GameObject clone = m_PoolConfig[i];
            Config _config = m_ConfigList[i];

            if (m_LoginProcessType == 0 && m_InputServerIP.text != _config.m_ServerListDomain)
                continue;
            if (m_LoginProcessType == 1 && m_InputServerIP.text != _config.m_IP)
                continue;
            if (m_LoginProcessType == 2 && m_InputServerIP.text != _config.m_IP)
                continue;
            //if (m_InputUpdateDomainIP.text != _config.m_UpdateDomainIP)
            //continue;
            //if (m_InputUpdateDomainPort.text != _config.m_UpdateDomainPort)
            //continue;
            if (m_InputSiteID.text != _config.m_SiteID)
                continue;
            if (m_InputPlatformID.text != _config.m_PlatformID)
                continue;
            if (m_InputServerOption.text != _config.m_ServerOption)
                continue;
            //if (m_InputChargeDomainIP.text != _config.m_ChargeDomainIP)
            //continue;
            //if (m_InputIPReuestIP.text != _config.m_IPReuestIP)
            //continue;
            if (m_LoginProcessType == 0 && m_InputServerPort.text != _config.m_WWWPort)
                continue;
            if (m_LoginProcessType == 1 && m_InputServerPort.text != _config.m_ListPort)
                continue;
            if (m_LoginProcessType == 2 && m_InputServerPort.text != _config.m_ServerPort)
                continue;

            m_ChooseServer.gameObject.SetActive(true);
            m_ChooseServer.transform.SetParent(clone.transform);
            m_ChooseServer.transform.SetAsFirstSibling();
            m_ChooseServer.anchoredPosition = Vector2.zero;
            m_ServerIndex = i;
            return;
        }
        m_ServerIndex = -1;
    }
}
