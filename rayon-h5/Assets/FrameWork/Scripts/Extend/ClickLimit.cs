using UnityEngine;
using UnityEngine.Events;
using UnityEngine.UI;

public class ClickLimit : MonoBehaviour
{
    public void Init(Button button, UnityAction call)
    {
        float interval = 0.2f;
        float lastClickTime = -1 * interval;

        button.onClick.AddListener(() =>
        {
            if (Time.timeSinceLevelLoad - lastClickTime >= interval && call != null)
            {
                lastClickTime = Time.timeSinceLevelLoad;

                ButtonClickSound clickSound = gameObject.GetComponent<ButtonClickSound>();

                if (clickSound && clickSound.ClickSound)
                    SoundManager.Instance.PlaySound(clickSound.ClickSound);
                else
                    SoundManager.Instance.PlaySound(AppDefine.HallNumber, "ButtonClick");

                call();
            }
        });
    }
}