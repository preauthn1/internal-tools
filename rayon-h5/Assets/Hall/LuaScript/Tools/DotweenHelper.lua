-- region *.lua
-- Date
-- 此文件由[BabeLua]插件自动生成
require("Module")
mmodule(...)

local G = require "Global"

-- 显示完整金币数
function ToCashDetail(num, isConvert)
    num = tonumber(num)
    if (isConvert == nil or isConvert == true) then
        num = num / coinShowRate
    end
    if (coinShowRate > 1) then
        num = G.Tool.Math.Decimal(num, 2)
        return clearSuffix(string.format("%.2f", num))
    else
        return string.format("%s", G.Tool.Math.ParseInt(num * coinShowRate))
    end
end

function BubbleScaleS2B1( ani_obj, size,t,Completefunc)
       -- ani_obj.transform.localScale = CS.UnityEngine.Vector3.one
        local sequenceb = CS.DG.Tweening.DOTween.Sequence()
 
     
        sequenceb:Append(ani_obj.transform:DOScale(CS.UnityEngine.Vector3.one * size ,t))--縮小  (Size*0.5) 倍        
        sequenceb:Append(ani_obj.transform:DOScale(CS.UnityEngine.Vector3.one * 1 , t))--回復

        sequenceb:OnComplete( function()
              --結束要做的事情
              if(Completefunc) then
               Completefunc()
              end
        end)
end

function BubbleScaleS2B2( ani_obj, size,t,Completefunc)
       -- ani_obj.transform.localScale = CS.UnityEngine.Vector3.one
        local sequenceb = CS.DG.Tweening.DOTween.Sequence()
 
     
        sequenceb:Append(ani_obj.transform:DOScale(CS.UnityEngine.Vector3.one * size ,t))--縮小  (Size*0.5) 倍
        sequenceb:Append(ani_obj.transform:DOScale(CS.UnityEngine.Vector3.one * (size*0.5), 0.1))--縮小  (Size*0.5) 倍
        sequenceb:Append(ani_obj.transform:DOScale(CS.UnityEngine.Vector3.one * size, t))--放大 Size 倍
        sequenceb:Append(ani_obj.transform:DOScale(CS.UnityEngine.Vector3.one * (size*0.7) ,0.1))--縮小  (Size*0.7) 倍
        sequenceb:Append(ani_obj.transform:DOScale(CS.UnityEngine.Vector3.one * 1 , t))--回復

        sequenceb:OnComplete( function()
              --結束要做的事情
              if(Completefunc) then
               Completefunc()
              end
        end)
end

function LandingScale( ani_obj, size,t,Completefunc)
        ani_obj.transform.localScale = CS.UnityEngine.Vector3.one*size
        local sequenceb = CS.DG.Tweening.DOTween.Sequence()
        local sequenceb1 = CS.DG.Tweening.DOTween.Sequence()
      
        sequenceb:Append(ani_obj.transform:DOScale(CS.UnityEngine.Vector3.one * 1  ,t))
        sequenceb:Append(ani_obj.transform:DOScale(CS.UnityEngine.Vector3.one * 1.1  ,0.05))
        sequenceb:Append(ani_obj.transform:DOScale(CS.UnityEngine.Vector3.one  , 0.05))
        
        --  sequenceb:Append(ani_obj.transform:DORotate(CS.UnityEngine.Vector3.zero ,t))

        sequenceb:OnComplete( function()
              --結束要做的事情
              if(Completefunc) then
               Completefunc()
              end
        end)
end
--
function Shake(ani_obj,rootPos,size,t,Completefunc)
          local sequence = CS.DG.Tweening.DOTween.Sequence()

         -- local rootStartPos=ani_obj.transform.position;
           local rootStartPos= rootPos;
           sequence:Append(ani_obj.transform:DOMove(CS.UnityEngine.Vector3(rootStartPos.x + size,ani_obj.transform.position.y,0), 0.1))
           sequence:Append(ani_obj.transform:DOMove(CS.UnityEngine.Vector3(rootStartPos.x,ani_obj.transform.position.y,0), t))
           sequence:OnComplete( function() 
            
            ani_obj.transform.position = rootStartPos 

          end)
end










-- endregion
