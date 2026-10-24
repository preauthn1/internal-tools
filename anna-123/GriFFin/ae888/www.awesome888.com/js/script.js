// scroll
$(document).ready(function () {
    $('.pop-banner').delay( 3000 ).fadeOut( 400 );
    $('.icon-close').click(function () {
        $('.pop-banner').hide();
    });

    $('#fullpage').fullpage({
        menu: '#menu',
        anchors: ['', 'service', 'advantage', 'products', 'partners', 'contact'],
        autoScrolling: false,
        fitToSection: false,
        afterLoad: function (anchorLink, index) {
            var loadedSection = $(this);
            if (index >= 2) {
                $('.sicon, .gotop').addClass('showicon');
            }
            else {
                $('.gotop').removeClass('showicon');
            }
        }
    });
});

//lang
function setLang(e) {
    var title = Object.keys(e);
    for (var i = 0; i < title.length; i++) {
        $('[text=' + title[i] + ']').html(e[title[i]]);
    };
}

$('.lang').click(function () {
    $(this).toggleClass('zh');
    if ($('.lang').hasClass('zh')) {
        setLang(i18n.zh);
    } else {
        setLang(i18n.en);
    }
})

function ipLookUp() {
    $.ajax('https://example.invalid/json')
        .then(
            function success(response) {
                console.log(response.country_name);
                if (response.country_name == 'Taiwan' | 'China') {
                    $('.lang').addClass('zh');
                    setLang(i18n.zh);
                } else {
                    $('.lang').removeClass('zh');
                    setLang(i18n.en);
                };
            },
            function fail(data, status) {
                console.log('Request failed.  Returned status of',
                    status);
            }
        );
}
ipLookUp()


$('.menu-close').click(function () {
    $('header').toggleClass('menu-open');
    var menuOpen = $('header').hasClass('menu-open');
    if (menuOpen == true) {
        $('header.menu-open nav').click(function () {
            $('header').removeClass('menu-open');
        })
    }
});



$('.logo, .gotop').click(function () {
    $("html, body").animate({ scrollTop: 0 }, 300);
});

$('.chess-case .chess-awd-box-desktop').click(function () {
    $('.chess-awd-box div, section.slider').removeClass('active');
    $('.chess-case .chess-awd-box-desktop, .s-desktop').addClass('active');
});
$('.chess-case .chess-awd-box-mobile').click(function () {
    $('.chess-awd-box div, section.slider').removeClass('active');
    $('.chess-case .chess-awd-box-mobile, .s-mobile').addClass('active');
})