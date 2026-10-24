using UnityEngine;
using UnityEngine.UI;

/// <summary>
/// 打包的时候，去掉引用大厅字体，通过资源加载方式读取
/// </summary>
[RequireComponent(typeof(Text))]
public class BindingFont : MonoBehaviour
{
    [SerializeField]
    private string m_FontName;

    public void SetFontName(string fontName)
    {
        m_FontName = fontName;
    }

    public void Start()
    {
        Font font = AssetBundleManager.Instance.LoadAsset<Font>(AppDefine.HallNumber, m_FontName);
        gameObject.GetComponent<Text>().font = font == null ? Resources.GetBuiltinResource<Font>("Arial.ttf") : font;
    }
}
