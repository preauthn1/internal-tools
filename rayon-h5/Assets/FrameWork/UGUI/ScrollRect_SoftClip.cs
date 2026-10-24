using UnityEngine;
using UnityEngine.UI;

public class ScrollRect_SoftClip : MonoBehaviour
{
    public GameObject m_MaskObjectRoot;

    public int m_ClipSoftX = 80;    // X方向 Range(1,200)
    public int m_ClipSoftY = 1;     // Y方向 Range(1,200)

    public void Start()
    {
        // check if a regular mask is used, if so, replace with RectMask2D
        if (this.gameObject.GetComponentInChildren<Mask>() != null)
        {
            Destroy(this.gameObject.GetComponentInChildren<Mask>());
            Image vpImage = this.gameObject.GetComponentInChildren<Image>();
            vpImage.color = Color.clear;
            vpImage.sprite = null;
        }

        // apply rectMask2D
        if (this.gameObject.GetComponentInChildren<RectMask2D>() == null)
            this.gameObject.AddComponent<RectMask2D>();

        Material mat = new Material(Shader.Find("ZZL/UGUI/UI Soft Mask"));
        mat.SetFloat("_ClipSoftX", m_ClipSoftX);
        mat.SetFloat("_ClipSoftY", m_ClipSoftY);

        if (m_MaskObjectRoot != null)
            setMaterialOnChildren(m_MaskObjectRoot.transform, mat);
        else
            setMaterialOnChildren(this.gameObject.transform, mat);
    }

    /// <summary>
    /// Recursive function that sets the material on all children of the content
    /// </summary>
    private void setMaterialOnChildren(Transform T, Material mat)
    {
        Graphic[] graphics = T.gameObject.GetComponentsInChildren<Graphic>(true);
        foreach (Graphic graphic in graphics)
            graphic.material = mat;

        Text[] texts = T.gameObject.GetComponentsInChildren<Text>(true);
        foreach (Text text in texts)
            text.material = mat;
    }
}
