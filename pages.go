package main

import (
	"strings"
	"html"
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
function get(uri, f) {
	var xmlHttp = null;
	try {
		xmlHttp = new XMLHttpRequest();
	} catch(e) {
		try {
			xmlHttp = new ActiveXObject("Microsoft.XMLHTTP");
		} catch(e) {
			try {
				xmlHttp = new ActiveXObject("Msxml2.XMLHTTP");
			} catch(e) {
				xmlHttp = null;
			}
		}
	}

	if (xmlHttp) {
		xmlHttp.open('GET', uri, true);
		xmlHttp.onreadystatechange = function() {
			if (xmlHttp.readyState == 4) {
				f(xmlHttp.responseText);
			}
		};
		xmlHttp.send(null);
	} else {
		alert("error: failed to create xmlHttp object");
	}
}

var upload_started = false;
var percent = 0;

function start() {
	if (!upload_started) {
		document.forms["upload"].submit();
		start_progress();
		upload_started = true;
	} else {
		alert("You already started a download!");
	}
}

function start_progress() {
	get("/progress/{upload_id}", function(text) {
		var new_percent = parseInt(text);
		if (new_percent > percent) {
			percent = new_percent;
		}
		update_progress(percent);
		if (percent == 100) {
			finish_progress();
		} else {
			window.setTimeout(start_progress, 1000);
		}
	});
}

function update_progress(percent) {
	var div = document.getElementById("progress");
	div.innerHTML = "Uploading... " + percent + "%";
}

function finish_progress() {
	var div = document.getElementById("progress");
	div.innerHTML = 'Upload finished. <a href="/files/{upload_id}">Uploaded to here.</a>';
	var savebtn = document.getElementById("savebtn");
	savebtn.disabled = false;
}
</script>
</head>
<body>
<h1>SuperUpload</h1>
<form action="/upload/{upload_id}" method="post" id="upload" name="upload" target="uploadiframe" enctype="multipart/form-data">
<input type="file" name="file" id="file" onchange="start();">
</form>
<div id="progress">Please select file to upload</div>
<form action="/savetext/{upload_id}" method="post">
<textarea rows="4" cols="80" id="text" name="text"></textarea><br>
<input type="submit" value="Save" id="savebtn" disabled>
</form>
<iframe name="uploadiframe" style="display: none"></iframe>
</body>
</html>
`

func UploadPage(upload_id string) []byte {
	return []byte(strings.Replace(upload_tmpl, "{upload_id}", upload_id, -1))
}

var information_tmpl = `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
</head>
<body>
<h1>Uploaded File</h1>
<div>{description}</div>
<div><a href="/files/{upload_id}">Download here</a></div>
</html>
`

func InformationPage(upload_id, description string) []byte {
	tmpl := strings.Replace(information_tmpl, "{upload_id}", upload_id, -1)
	tmpl = strings.Replace(tmpl, "{description}", escape_text(description), -1)
	return []byte(tmpl)
}

func escape_text(text string) string {
	text = html.EscapeString(text)
	text = strings.Replace(text, "\n", "<br>", -1)
	return text
}
