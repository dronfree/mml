$(document).ready(function() {

    new Clipboard('.clipboard');

    $.get("/mbm/box", function(data){
        boxData = data;
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
                            $.each(data, function(index, value) {
                                value.Raw = '/mbm/mail?' + $.param({
                                    box:    boxData.Box,
                                    sessid: boxData.Sessid,
                                    id:     value.Id
                                })
                            })
                            $("#template-container").loadTemplate($("#template"), data.reverse());
                            $("#inboxAmount").text(data.length);
                        }
                    })
            }, 5000);
        }
    })
});