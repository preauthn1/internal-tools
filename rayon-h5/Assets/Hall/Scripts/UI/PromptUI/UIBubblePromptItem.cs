using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;

public class UIBubblePromptItem : MonoBehaviour
{
    [SerializeField]
    Text contentText = null;

    public void SetBubblePromptItemInfo(string content, float lastTime)
    {
        contentText.text = content;
        RectTransform itemRectTrans = transform.GetComponent<RectTransform>();
        itemRectTrans.sizeDelta = new Vector2(contentText.preferredWidth + 100f, itemRectTrans.sizeDelta.y);
        gameObject.SetActive(true);
        Invoke("DestroyBubble", lastTime);
    }

    public void DestroyBubble()
    {
        Destroy(gameObject);
    }
}
