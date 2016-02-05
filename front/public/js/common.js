$(document).ready(function(){
    $.get("/mbm/box", function(data){
        if (data.Error) {
            $("#error").text(data.Error)
        }
        else {
            $("#email").text(data.Box);
            $("#sessid").text(data.Sessid);
            setInterval(function(){
                $.get("/mbm/mails", {box: data.Box, sessid: data.Sessid})
                    .done(function(content){
                        if (content != "") {
                            $("#content").text(content)
                        }
                    })
            }, 5000);
        }
    })
});