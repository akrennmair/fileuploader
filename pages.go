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
var upload_id = "{upload_id}";

// TODO
</script>
</head>
<body>
<h1>File Upload</h1>
<form action="/upload/{upload_id}" method="post" target="upload" enctype="multipart/form-data">
<input type="file">
<input type="button" value="Upload File">
</form>
<div id="progress"></div>
<iframe name="upload" style="display: none"></iframe>
</body>
</html>
`

func UploadPage(upid string) []byte {
	return []byte(strings.Replace(upload_tmpl, "{upload_id}", upid, -1))
}
