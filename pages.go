package main

import (
	"html"
	"strings"
)

// this file contains various helper functions to render templates. Currently,
// this is as simple as possible, i.e. simple placeholders are replaced by text
// passed as arguments to functions.

var error_tmpl = `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>{errormsg}</title>
</head>
<body>
<h1>Error: {errormsg}</h1>
</body>
</html>
`

// this function renders the default error page
func ErrorPage(errmsg string) []byte {
	return []byte(strings.Replace(error_tmpl, "{errormsg}", errmsg, -1))
}

var upload_tmpl = `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script src="lib.js" type="text/javascript"></script>
<script src="upload.js" type="text/javascript"></script>
</head>
<body onload="ak.upload.reset_all();">
<h1>SuperUpload</h1>
<form action="/upload/INVALID" method="post" id="frm_upload" name="frm_upload" target="uploadiframe" enctype="multipart/form-data">
<input type="file" name="file" onchange="ak.upload.start_upload();">
</form>
<div id="progress">Please select file to upload</div>
<form action="/savedesc/INVALID" method="post" id="frm_save" name="frm_save">
<textarea rows="4" cols="80" id="input_desc" name="input_desc"></textarea><br>
<input type="submit" value="Save" onclick="ak.upload.save_desc(); return false;" id="btn_save">
</form>
<div id="div_savemsg"></div>
<iframe name="uploadiframe" style="display: none"></iframe>
</body>
</html>
`

// this function renders the upload page. All of the JavaScript client logic is contained
// in this template.
func UploadPage() []byte {
	return []byte(upload_tmpl)
}

var information_tmpl = `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
</head>
<body>
<h1>Uploaded File {filename}</h1>
<div>{description}</div>
<div><a href="/files/{upload_id}">Download {filename}</a></div>
</html>
`

// this function renders the description page.
func InformationPage(upload_id, description, filename string) []byte {
	tmpl := strings.Replace(information_tmpl, "{upload_id}", upload_id, -1)
	tmpl = strings.Replace(tmpl, "{description}", escape_text(description), -1)
	tmpl = strings.Replace(tmpl, "{filename}", escape_text(filename), -1)
	return []byte(tmpl)
}

// this helper function escapes a string so that it can be safely embedded into
// HTML code without the danger of cross-site scripting (XSS).
func escape_text(text string) string {
	text = html.EscapeString(text)
	text = strings.Replace(text, "\n", "<br>", -1)
	return text
}
