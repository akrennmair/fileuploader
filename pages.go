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
				//console.log("responseText: " + xmlHttp.responseText);
				f(xmlHttp.responseText);
			}
		};
		xmlHttp.send(null);
	} else {
		//console.log("error: failed to create xmlHttp object");
	}
}

var upload_id = null;
var upload_started = false;
var percent = 0;
var save_pending = false;

function init() {
	alert('init called');
	upload_id = null;
	upload_started = false;
	percent = 0;
	save_pending = false;
}

function set_upload_id(id) {
	upload_id = id;
	document.forms["frm_upload"].action = "/upload/" + upload_id;
	document.forms["frm_save"].action = "/savedesc/" + upload_id;
}

function start_upload() {
	//console.log("enterting start_upload()");
	if (!upload_id) {
		//console.log("no upload ID yet");
		get('/requpid?nocache=' + Math.random(), function(id) {
			//console.log("setting upload ID");
			set_upload_id(id);
			//console.log("restarting start_upload()");
			window.setTimeout(start_upload, 1);
		});
		return;
	}
	if (!upload_started) {
		//console.log("submitting upload");
		document.forms["frm_upload"].submit();
		start_progress();
		upload_started = true;
	} else {
		alert("You already started a download!");
	}
}

function start_progress() {
	//console.log("starting progress");
	get("/progress/" + upload_id + "?nocache=" + Math.random(), function(text) {
		//console.log("got progress: " + text);
		var new_percent = parseInt(text);
		if (new_percent > percent) {
			percent = new_percent;
		}
		update_progress(percent);
		if (percent == 100) {
			//console.log("finished with upload");
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
	document.getElementById("progress").innerHTML = 'Upload finished. <a href="/files/' + upload_id + '">Uploaded to here.</a>';
	if (save_pending) {
		document.getElementById("input_desc").disabled = false;
		document.forms["frm_save"].submit();
	}
}

function save_desc() {
	if (upload_started) {
		if (percent == 100) {
			document.forms["frm_save"].submit();
		} else {
			document.getElementById("btn_save").disabled = true;
			document.getElementById("input_desc").disabled = true;
			save_pending = true;
			document.getElementById("div_savemsg").innerHTML = "Uploading file and saving description...";
		}
	} else {
		document.getElementById("div_savemsg").innerHTML = "Please choose a file to be uploaded.";
	}
}
</script>
</head>
<body onload="init();">
<h1>SuperUpload</h1>
<form action="/upload/INVALID" method="post" id="frm_upload" name="frm_upload" target="uploadiframe" enctype="multipart/form-data">
<input type="file" name="file" onchange="start_upload();">
</form>
<div id="progress">Please select file to upload</div>
<form action="/savedesc/INVALID" method="post" id="frm_save" name="frm_save">
<textarea rows="4" cols="80" id="input_desc" name="input_desc"></textarea><br>
<input type="submit" value="Save" onclick="save_desc(); return false;" id="btn_save">
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
