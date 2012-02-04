package main

import (
	"strings"
)

var error_tmpl = `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
</head>
<body>
<h1>An error occured: {errormsg}</h1>
</body>
</html>
`


func ErrorPage(errmsg string) []byte {
	return []byte(strings.Replace(error_tmpl, "{errormsg}", errmsg, -1))
}

var upload_tmpl = `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script type="text/javascript">
window.uploader = window.uploader || {};
window.uploader = {
	id: "{upload_id}",
	start: function() {
		document.forms["upload"].submit();
	}
};
</script>
</head>
<body>
<h1>SuperUpload</h1>
<form action="/upload/{upload_id}" method="post" id="upload" name="upload" target="uploadiframe" enctype="multipart/form-data">
<input type="file" name="file" id="file" onchange="window.uploader.start();">
</form>
<div id="progress"></div>
<form action="/savetext/{upload_id}" method="post">
<textarea rows="4" cols="80" id="text"></textarea><br>
<input type="submit" value="Save">
</form>
<iframe name="uploadiframe" style="display: none"></iframe>
</body>
</html>
`

func UploadPage(upid string) []byte {
	return []byte(strings.Replace(upload_tmpl, "{upload_id}", upid, -1))
}
