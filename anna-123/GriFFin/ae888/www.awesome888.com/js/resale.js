var provider_list = {
    chess: [
        { id: 'kingmaker', name: 'King Maker', whitelabel: true },
        { id: '365', name: '365棋牌', whitelabel: true },
        { id: 'cq9', name: 'CQ9', whitelabel: true },
        { id: 'leg', name: 'LEG 乐游', whitelabel: true },
        { id: 'ky', name: 'KY 開元', whitelabel: true },
        { id: 'vg', name: 'VG 财神棋牌', whitelabel: true },
        { id: 'jdb', name: 'JDB', whitelabel: true },
        { id: 'mt', name: 'mt 美天棋牌', whitelabel: false },
        { id: 'ds', name: 'DS 德胜棋牌', whitelabel: true },
        { id: 'lc', name: '龙城棋牌', whitelabel: true },
        { id: 'lucky', name: '幸运棋牌', whitelabel: false },
        { id: 'tm', name: 'tm 天美棋牌', whitelabel: false },
        { id: false, name: '人人棋牌', whitelabel: false },
        { id: 'ae', name: 'AE 棋牌', whitelabel: true },
        { id: false, name: '761 棋乐游', whitelabel: false },
        { id: 'hl', name: '欢乐棋牌', whitelabel: false },
        { id: 'wm', name: 'WM', whitelabel: false },
        { id: false, name: 'SBG', whitelabel: false },
        { id: false, name: '新世界', whitelabel: false },
    ],
    sports: [
        { id: 'ag', name: 'AG', whitelabel: false },
        { id: 'bbin', name: 'BBIN', whitelabel: false },
        { id: 'cmdbet', name: 'CMDBET', whitelabel: true },
        { id: false, name: 'IBC', whitelabel: false },
        { id: 'im', name: 'InplayMatrix', whitelabel: false },
        { id: false, name: 'M8', whitelabel: false },
        { id: 'sbobet', name: 'SBOBET', whitelabel: true },
        { id: 'ugaming', name: 'United Gaming', whitelabel: true },
        { id: 'bti', name: 'BTi', whitelabel: true },
        { id: 'mg', name: 'MG Plus', whitelabel: false },
        { id: false, name: 'Majasports', whitelabel: false },
        { id: false, name: 'UPBET', whitelabel: true },
    ],
    live: [
        { id: 'ag', name: 'AG', whitelabel: true },
        { id: 'allbet', name: 'Allbet 歐博', whitelabel: true },
        { id: 'bbin', name: 'BBIN', whitelabel: false },
        { id: 'bg', name: 'BG 大游', whitelabel: false },
        { id: 'dg', name: 'Dream Gaming', whitelabel: true },
        { id: 'ebet', name: 'Ebet', whitelabel: false },
        { id: 'evo', name: 'Evolution Gaming', whitelabel: true },
        { id: 'gd', name: 'Gold Deluxe', whitelabel: true },
        { id: 'gpi', name: 'Gameplay Interactive', whitelabel: false },
        { id: false, name: 'IBO', whitelabel: false },
        { id: 'ig', name: 'IGaming', whitelabel: false },
        { id: 'venus', name: 'Venus Gaming', whitelabel: true },
        { id: 'mg', name: 'MG Plus', whitelabel: true },
        { id: false, name: 'Oriental Game', whitelabel: false },
        { id: 'pt', name: 'Playtech', whitelabel: false },
        { id: 'sa', name: 'SA Gaming 沙龙', whitelabel: true },
        { id: 'sbobet', name: 'SBOBET', whitelabel: false },
        { id: 'sexybcrt', name: 'AE Sexy Baccrat', whitelabel: true },
        { id: 'ean2', name: 'N2-Live', whitelabel: false },
        { id: 'sb', name: 'Suncity-Group 申博', whitelabel: true },
        { id: 'wm', name: 'WM 完美', whitelabel: true },
        { id: false, name: 'XPG', whitelabel: false },
        { id: false, name: 'IDN Live', whitelabel: false },
        { id: 'cq9', name: 'CQ9', whitelabel: true },
    ],
    mpg: [
        { id: 'ag', name: 'AG', whitelabel: true },
        { id: 'bbin', name: 'BBIN', whitelabel: false },
        { id: 'bg', name: 'BG 大游', whitelabel: false },
        { id: 'cq9', name: 'CQ9', whitelabel: true },
        { id: 'jdb', name: 'JDB', whitelabel: true },
        { id: 'mw', name: 'MW', whitelabel: true },
        { id: 'mg', name: 'MG Plus', whitelabel: true },
        { id: 'pt', name: 'Playtech', whitelabel: false },
        { id: 'simpleplay', name: 'SimplePlay', whitelabel: true },
        { id: 'spadegaming', name: 'Spade Gaming', whitelabel: false },
        { id: 'bowjin', name: '爆金捕鱼', whitelabel: false },
        { id: false, name: 'YL 捕鱼', whitelabel: false },
        { id: 'joker', name: 'Joker', whitelabel: false },
        { id: 'gamatron', name: 'Gamatron', whitelabel: false },
    ],
    
    lottery: [
        { id: false, name: '马来 4d彩', whitelabel: false },
        { id: 'tcg', name: '天成 TCG', whitelabel: false },
        { id: 'bbin', name: 'BBIN', whitelabel: false },
        { id: 'ae', name: 'AE', whitelabel: true },
        { id: 'vr', name: 'VR', whitelabel: true },
        { id: 'bg', name: 'BG 大游', whitelabel: false },
        { id: 'sgwin', name: 'SG 双赢彩票', whitelabel: false },
        { id: 'vsl', name: 'VSL', whitelabel: true },
        { id: 'gdlotto', name: 'GD Lotto 豪龙彩票', whitelabel: false },
        { id: false, name: 'KK 彩票', whitelabel: false },
        { id: false, name: '亚博彩票', whitelabel: false },
    ],
    egame: [
        { id: 'ag', name: 'AG', whitelabel: true },
        { id: 'xin', name: 'XIN电子', whitelabel: true },
        { id: 'ameba', name: 'Ameba', whitelabel: true },
        { id: 'bbin', name: 'BBIN', whitelabel: false },
        { id: 'bg', name: 'BG 大游', whitelabel: false },
        { id: 'betsoft', name: 'BetSoft', whitelabel: false },
        { id: 'cq9', name: 'CQ9', whitelabel: true },
        { id: 'dt', name: 'Dream Tech', whitelabel: false },
        { id: 'ggaming', name: 'GG', whitelabel: true },
        { id: 'gns', name: 'GNS', whitelabel: false },
        { id: 'hb', name: 'HB', whitelabel: false },
        { id: 'jdb', name: 'JDB', whitelabel: true },
        { id: 'laxino', name: 'LAXINO', whitelabel: true },
        { id: 'mg', name: 'MG Plus', whitelabel: true },
        { id: 'mw', name: 'MW', whitelabel: true },
        { id: 'netent', name: 'Netent', whitelabel: false },
        { id: 'nyx', name: 'NYX', whitelabel: false },
        { id: false, name: 'PGSOFT', whitelabel: false },
        { id: 'playngo', name: 'Play n GO', whitelabel: false },
        { id: 'ps', name: 'Playstar', whitelabel: true },
        { id: 'pp', name: 'Pragmatic Play', whitelabel: false },
        { id: 'pt', name: 'Playtech', whitelabel: false },
        { id: 'qtech', name: 'Qtech', whitelabel: false },
        { id: 'rtg', name: 'RealTimeGaming', whitelabel: false },
        { id: 'simpleplay', name: 'SimplePlay', whitelabel: true },
        { id: 'sbobet', name: 'SBOBET', whitelabel: false },
        { id: 'skywind', name: 'skywind', whitelabel: false },
        { id: 'spadegaming', name: 'Spade Gaming', whitelabel: false },
        { id: 'pg', name: 'Pocket Game', whitelabel: false },
        { id: 'rt', name: 'Red Tiger Gaming', whitelabel: true },
        { id: 'ttg', name: 'Top Trend Gaming', whitelabel: false },
        { id: 'uc8', name: 'UC8', whitelabel: false },
        { id: 'vt', name: 'VT 路虎', whitelabel: true },
        { id: 'yoplay', name: 'Yoplay', whitelabel: true },
        { id: '918kiss', name: '918kiss', whitelabel: false },
        { id: 'ig', name: 'IG电子', whitelabel: false },
        { id: 'ganapati', name: 'ganapati', whitelabel: false },
        { id: 'pussy888', name: 'Pussy 888', whitelabel: true },
        { id: 'onegame', name: 'One Game', whitelabel: false },
        { id: 'laxino', name: 'Laxino Premium', whitelabel: false },
        { id: false, name: 'GA', whitelabel: false },
        { id: 'joker', name: 'Joker', whitelabel: false },
        { id: 'sea', name: 'SEA Gaming', whitelabel: true },
        { id: 'gamatron', name: 'Gamatron', whitelabel: false },
    ],
    esport: [
        { id: 'avia', name: '泛亚电竞', whitelabel: false },
        { id: 'esportsbull', name: 'EsportBull 电竞牛', whitelabel: true },
        { id: 'tf', name: 'TFGaming 雷火电竞', whitelabel: true },
    ],
    cock: [
        { id: 's128', name: 'S128 龍鳳鬥雞', whitelabel: true },
        { id: 'sv388', name: 'SV388 斗鸡', whitelabel: true },
    ],
    horse: [
        { id: 'trc', name: 'TRC', whitelabel: true },
    ],
    bull: [
        { id: 'sm558', name: 'SM558', whitelabel: false },
    ],
};

var type_title = {
    horse: '赛马 Horse Racing',
    mpg: '捕鱼 Fishing',
    lottery: '彩票 Lottery',
    egame: '电子 Slot',
    esport: '電競 Esport',
    cock: '鬥雞 Cockfighting',
    bull: '鬥牛 Bullfighting',
    chess: '棋牌 Card Game',
    sports: '体育 Sport',
    live: '真人 Live Casino',
}

function setLang(e) {
    var title = Object.keys(e);
    for (var i = 0, max = title.length; i < max; i++) {
        $('[text=' + title[i] + ']').text(e[title[i]]);
    };
    for (var i = 0, max = $('[data-word]').length; i < max; i++) {
        $('[data-word]').eq(i).attr('data-word', $('[data-word]').eq(i).text())
    };
}

$(document).ready(function() {
    $.each( provider_list, function(type) {
        var providerHtml = '';
        $.each(provider_list[type], function(provideridx, provider) {
            var cover_exist = provider.id && ' ' || 'no-img';
            var background_img  = provider.id && './img/api-logo/'+provider.id+'.png' || '';
            var tags = '';
            if(provider.whitelabel){ tags += '<span class="tag tag-whitelabel">包网</span>' };

            var item = 
            '<li class="'+provider.id+'">'+
                '<div class="provider-img '+cover_exist+'" data-name="'+provider.name+'" style="background-image: url('+background_img+');"></div>'+
                '<h4 class="provider-title">'+provider.name+'</h4>'+tags
            '</li>';
            providerHtml += item;
        });

        var typeHtml = 
        '<div class="type-block type'+type+'">'+
            '<h3>'+type_title[type]+'</h3>'+
            '<ul class="provider-list">'+providerHtml+'</ul>'
        '</div>';

        $('#resale_list').append(typeHtml);
    });

    $("#year").html(new Date().getFullYear());

    function afterSetLang(){
        // come in, computer language. window.navigator.language ： zh-TW 繁中 、 zh-CN 簡中 、 th-TH 泰文 、 vi-VN 越南 、 en-US 英文
        switch(window.navigator.language) {
          case 'zh-TW' :
          case 'zh-CN' :
            setLang(i18n.zh);
            break;
          case 'th-TH' :
            setLang(i18n.th);
            break;
          case 'vi-VN' :
            setLang(i18n.vn);
            break;
          default :
            setLang(i18n.en);
            break;
        };
    }
    $.ajax({
        url:'./load/contact.html',
        dataType:'html'
    }).done(function(data) {
        $('#contact').append(data);
        afterSetLang();
    })
    $.ajax({
        url:'./load/foot.html',
        dataType:'html'
    }).done(function(data) {
        $('#foot').append(data);
        afterSetLang();
    })

    $('#goToTop').click(function () {
        $("html,body").animate({ scrollTop: 0 }, 600);
    });
});