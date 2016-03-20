$(document).ready(function() {

    new Clipboard('.clipboard');

    $.get("/mbm/box", function(data){
        if (data.Error) {
            $("#error").text(data.Error)
        }
        else {
            $(".el-email").each(function() {
                $(this).text(data.Box);
                if($(this).is('a')) {
                    $(this).attr("href", "mailto:" + data.Box);
                }
            });

            $(".timer").data("seconds-left", data.ExpiresIn);
            $('.timer').startTimer();

            setInterval(function(){
                $.get("/mbm/mails", {box: data.Box, sessid: data.Sessid})
                    .done(function(data){
                        if (data != "" && data.length != $("#inboxAmount").text()) {
                            $("#emptybox").hide();
                            $("#template-container").loadTemplate($("#template"), data.reverse());
                            $("#inboxAmount").text(data.length);
                        }
                    })
            }, 5000);
        }
    })
});