
// https://css-tricks.com/snippets/javascript/get-url-variables/
function getQueryVariable(variable) {
    var query = window.location.search.substring(1);
    var vars = query.split("&");
    for (var i=0; i<vars.length; i++) {
        var pair = vars[i].split("=");
        if (pair[0] == variable) {
            return pair[1];
        }
    }
    return false
}

function displayShortenedUrl(url) {
    $(".js-success-strip").removeClass("hidden");
    var url = "/urls/"+url;
    var html = "Your shortened url is <a href='"+url+"'>"+url+"</a>.";
    $(".js-success-strip>.message").html(html);
}

$(function () {
    var url = getQueryVariable("done");
    if (url) {
        displayShortenedUrl(url);
    }
});
