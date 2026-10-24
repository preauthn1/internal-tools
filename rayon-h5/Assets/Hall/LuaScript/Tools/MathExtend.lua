require("Module")
mmodule(...)

function ParseInt(val)
    return math.floor(val + 0.000001)
end

function Pow(v, n)
    if (n == 1) then return v end
    local val = 1
    for i = 1, n, 1 do
        val = val * v
    end
    return val
end


function Decimal(v, num)
    if (num == nil) then num = 0 end
    v = ParseInt(v * Pow(10, num))
    v = v / Pow(10, num)
    return v
end


function MakeRandom(s,e)
    
    RndNum={}
    output={}
    newEnd=(e-s+1)
    for i = 1, (e-s+1), 1 do
          RndNum[i]=(s+i-1);
    end
    index=1;
    x=(function(rn)
        rnindex=1
        newRndNum={}
        for i,value in pairs(rn) do 
   
            if(rn[i] ~= nil ) then
                newRndNum[rnindex]=rn[i];                
                rnindex=rnindex+1
            end
        end     
        return newRndNum;
    end)
   breakindex=0;
  
    while #output < (e-s+1) do 
        rndn=ParseInt(CS.UnityEngine.Random.Range(1, (newEnd+0.99)))
      
        output[index]=RndNum[rndn];        
        newEnd=newEnd-1
        index=index+1
        RndNum[rndn]=nil;        
        RndNum=x(RndNum);        

        breakindex=breakindex+1;
        if(breakindex>10) then
            break;
        end
    end


    return output;
end

