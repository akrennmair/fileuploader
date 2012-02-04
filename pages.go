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
		console.log("error: failed to create xmlHttp object");
	}
}

var upload_started = false;
var percent = 0;
var save_pending = false;

function start_upload() {
	if (!upload_started) {
		document.forms["frm_upload"].submit();
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
	document.getElementById("progress").innerHTML = "Uploading... " + percent + "%";
}

function finish_progress() {
	document.getElementById("progress").innerHTML = 'Upload finished. <a href="/files/{upload_id}">Uploaded to here.</a>';
	if (save_pending) {
		document.forms["frm_save"].submit();
	}
}

function save_desc() {
	if (upload_started) {
		document.getElementById("btn_save").disabled = true;
		document.getElementById("input_desc").disabled = true;
		if (percent == 100) {
			document.forms["frm_save"].submit();
		} else {
			save_pending = true;
			document.getElementById("div_savemsg").innerHTML = "Uploading file and saving description...";
		}
	} else {
		document.getElementById("div_savemsg").innerHTML = "Please choose a file to be uploaded.";
	}
}
</script>
</head>
<body>
<h1>SuperUpload</h1>
<form action="/upload/{upload_id}" method="post" id="frm_upload" target="uploadiframe" enctype="multipart/form-data">
<input type="file" name="file" onchange="start_upload();">
</form>
<div id="progress">Please select file to upload</div>
<form action="/savetext/{upload_id}" method="post" id="frm_save">
<textarea rows="4" cols="80" id="input_desc"></textarea><br>
<input type="submit" value="Save" onclick="save_desc(); return false;" id="btn_save">
</form>
<div id="div_savemsg"></div>
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
<h1>Uploaded File {filename}</h1>
<div>{description}</div>
<div><a href="/files/{upload_id}">Download {filename}</a></div>
</html>
`

func InformationPage(upload_id, description, filename string) []byte {
	tmpl := strings.Replace(information_tmpl, "{upload_id}", upload_id, -1)
	tmpl = strings.Replace(tmpl, "{description}", escape_text(description), -1)
	tmpl = strings.Replace(tmpl, "{filename}", escape_text(filename), -1)
	return []byte(tmpl)
}

func escape_text(text string) string {
	text = html.EscapeString(text)
	text = strings.Replace(text, "\n", "<br>", -1)
	return text
}
