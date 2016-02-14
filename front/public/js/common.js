$(document).ready(function(){
    $.get("/mbm/box", function(data){
        if (data.Error) {
            $("#error").text(data.Error)
        }
        else {
            $("#email").text(data.Box);
            setInterval(function(){
                $.get("/mbm/mails", {box: data.Box, sessid: data.Sessid})
                    .done(function(data){
                        if (data != "" && data.length != $("#inboxAmount").text()) {
                            $("#emptybox").hide();
                            $("#template-container").loadTemplate($("#template"), data);
                            $("#inboxAmount").text(data.length);
                            console.log(data, data.length)
                        }
                    })
            }, 5000);
        }
    })
});