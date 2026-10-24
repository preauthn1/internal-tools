using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;

public class LoadingStreamingAssetsTexture : MonoBehaviour {

    private RawImage m_RawImage;
    private Image m_Image;

    private void Awake()
    {
        m_RawImage = this.gameObject.GetComponentInChildren<RawImage>();
        if (m_RawImage != null)
            m_RawImage.color = Color.clear;

        m_Image = this.gameObject.GetComponentInChildren<Image>();
        if (m_Image != null)
            m_Image.color = Color.clear;
    }

    private void Start()
    {
        if (m_RawImage != null)
            StartCoroutine(UpdateAssetsYid());

        if (m_Image != null)
            StartCoroutine(UpdateAssetsYid());
    }

    private IEnumerator UpdateAssetsYid()
    {
        //图片路径 带后缀    
        /* 要注意  WWW是要带后缀才能找到从而加载到资源，Resources是不用带后缀的*/
        string url = Utility.GetFullPathOfStreamingAsset("Logo.png");

        WWW www = new WWW(url);

        //等待资源加载完成
        yield return www;
        if (string.IsNullOrEmpty(www.error) == false)
        {
            Debug.Log("*********************** 资源加载失敗\n" + www.error);
        }
        else
        {
            if (m_RawImage != null)
            {
                m_RawImage.texture = www.texture;
            }

            if (m_Image != null)
            {
                // assign the downloaded image to sprite
                Texture2D texture = new Texture2D(www.texture.width, www.texture.height, TextureFormat.DXT1, false);
                www.LoadImageIntoTexture(texture);
                Rect rec = new Rect(0, 0, texture.width, texture.height);
                Sprite spriteToUse = Sprite.Create(texture, rec, new Vector2(0.5f, 0.5f), 100);
                m_Image.sprite = spriteToUse;
            }
        }

        if (m_RawImage != null)
        {
            m_RawImage.SetNativeSize();
            m_RawImage.color = Color.white;
        }

        if (m_Image != null)
        {
            m_Image.SetNativeSize();
            m_Image.color = Color.white;
        }
    }
}
