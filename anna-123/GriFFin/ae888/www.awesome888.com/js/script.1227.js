//lang
function setLang(e) {
    var title = Object.keys(e);
    for (var i = 0; i < title.length; i++) {
        $('[text=' + title[i] + ']').html(e[title[i]]);
    };
}

(function() {
    $.ajax('https://example.invalid/json').then(
        function success(response) {
            // console.log(response.country_name);
            if (response.country_name == 'Taiwan' | 'China') {
                $('.lang').addClass('zh');
                setLang(i18n.zh);
            } else {
                $('.lang').removeClass('zh');
                setLang(i18n.en);
            };
        },
        function fail(data, status) {
            console.log('Request failed.  Returned status of', status);
        }
    );
}());

// scroll
$(function(){
    $('.pop-banner').delay( 3000 ).fadeOut( 400 );
    $('.icon-close').click(function () {
        $('.pop-banner').hide();
    });
    // $('#fullpage').fullpage({
    //     menu: '#menu',
    //     anchors: ['', 'service', 'advantage', 'products', 'partners', 'contact'],
    //     autoScrolling: false,
    //     fitToSection: false,
    //     afterLoad: function (anchorLink, index) {
    //         var loadedSection = $(this);
    //         if (index >= 2) {
    //             $('.sicon, .gotop').addClass('showicon');
    //         }
    //         else {
    //             $('.gotop').removeClass('showicon');
    //         }
    //     }
    // });
    $('.lang').click(function () {
        $(this).toggleClass('zh');
        $('.lang').hasClass('zh') ? setLang(i18n.zh) : setLang(i18n.en);
    })
    $('.menu-close').click(function () {
        $('header').toggleClass('menu-open');
        if ($('header').hasClass('menu-open')) {
            $('header ul').click(function () {
                $('header').removeClass('menu-open');
            })
        }
    });
    $(window).on('scroll', function(){
        if ($(window).scrollTop() > $('.chess-feature').offset().top) {
            $('.gotop').addClass('showicon');
        } else {
            $('.gotop').removeClass('showicon');
        }
    }).scroll();
    $('.gotop').click(function () {
        $("html, body").animate({ scrollTop: 0 }, 300);
    });
    // $('.chess-case .chess-awd-box-desktop').click(function () {
    //     $('.chess-awd-box div, section.slider').removeClass('active');
    //     $('.chess-case .chess-awd-box-desktop, .s-desktop').addClass('active');
    // });
    // $('.chess-case .chess-awd-box-mobile').click(function () {
    //     $('.chess-awd-box div, section.slider').removeClass('active');
    //     $('.chess-case .chess-awd-box-mobile, .s-mobile').addClass('active');
    // })
});
