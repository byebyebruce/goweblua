package web

const htmlStr = `<html>  
    <head>  
        <title>command</title>  
    </head>  
    <body>  
        <style type="text/css">
            #all {
                color: #000000;
                background: #ececec;
                width: 300px;
                height: 200px;
            }

            div {
                line-height: 40px;
                text-align: right;
            }

            input {
                width: 400px;
                height: 20px;
            }

            select {
                padding: 5px 82px;
            }
        </style>
       

    
        <label for="API">Parameter:</label>
        <p>
            <!-- <input type="text" name="API" id="api" value="help" rows="100"> -->
            <textarea id="api" name="API" rows="10" cols="80">Hello world</textarea>
            <button id="btn" style="height:50px;width:50px">click</button>
        </p>


        <!-- <label for="json">Param:</label>
        <p>
            <textarea id="json" name="API" rows="10" cols="80"></textarea>
        </p> -->

        <label for="Return">Return:</label>
        <p>
            <textarea id="return" name="summary" rows="50" cols="120"></textarea>
        </p>
        
    </body>  
    <script src="http://apps.bdimg.com/libs/jquery/2.1.4/jquery.min.js"></script>

    <script type="text/javascript">
        

        $("#btn").on("click", function () {
            var param = $("#api").val()


            var url1 = window.location.href 

            document.getElementById("return").value = "waiting result..."
            
            console.log(url1)
            $.ajax({
                url: url1,
                type: "POST",
				contentType: "application/json; charset=utf-8",
				//dataType: "json",

                data: param,

                success: function (res) {
                    document.getElementById("return").value = res;
                },

                error: function(xhr,textStatus){
                    document.getElementById("return").value = "error, state = " + textStatus
                }
            })
            
        })

    </script>
</html>`
