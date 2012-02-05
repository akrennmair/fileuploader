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
// wrapper function to do XMLHttpRequests more easily
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

var upload_id = null;
var upload_started = false;
var percent = 0;
var save_pending = false;

// this function resets all state variables to their initial values and enables
// the description textarea and the save button.
function reset_all() {
	upload_id = null;
	upload_started = false;
	percent = 0;
	save_pending = false;
	document.getElementById("input_desc").disabled = false;
	document.getElementById("btn_save").disabled = false;
}

// this function saves the upload ID and resets the action URIs for
// the two forms.
function set_upload_id(id) {
	upload_id = id;
	document.forms["frm_upload"].action = "/upload/" + upload_id;
	document.forms["frm_save"].action = "/savedesc/" + upload_id;
}

// this function is called when the user selects a file.
function start_upload() {
	// if we have no upload ID yet, we fetch one and rerun start_upload().
	if (!upload_id) {
		// the nocache?=random is there to randomize the URL. This was added
		// to remedy browser caching that even wouldn't go away when the usual
		// HTTP response headers to prevent caching were set.
		get('/requpid?nocache=' + Math.random(), function(id) {
			set_upload_id(id);
			start_upload();
		});
		return;
	}
	// the user hasn't started an upload yet -> submit the form, and start
	// polling the upload progress.
	if (!upload_started) {
		document.forms["frm_upload"].submit();
		start_progress();
		upload_started = true;
	} else {
		// the user selected another file while an upload was already in progress.
		// we warn him that he has already started an upload.
		alert("You already started a download!");
	}
}

// this function fetches the current upload progress, displays it, and if the
// upload isn't complete yet, it schedules another upload progress poll.
function start_progress() {
	// see above. the ?nocache=random is there to prevent browser caching.
	get("/progress/" + upload_id + "?nocache=" + Math.random(), function(text) {
		var new_percent = parseInt(text);
		if (new_percent > percent) {
			percent = new_percent;
		}
		update_progress(percent);
		// if the upload is finished, we completely finish it by showing
		// the user additional information
		if (percent == 100) {
			finish_progress();
		} else {
			// upload not yet finished -> schedule another poll
			window.setTimeout(start_progress, 1000);
		}
	});
}

// this function displays the upload progress percents to the user while the upload
// is going on.
function update_progress(percent) {
	document.getElementById("progress").innerHTML = "Uploading... " + percent + "%";
}

// this function completes the upload progress.
function finish_progress() {
	// we show the user that the upload has finished plus a download link to the uploaded file.
	document.getElementById("progress").innerHTML = 'Upload finished. <a href="/files/' + upload_id + '">Uploaded to here.</a>';
	// if the user pressed the save button in the mean while, we submit the description save form.
	if (save_pending) {
		// the previously disabled (see below) input textarea needs to be reenabled, otherwise
		// its content will not the sent as part of the form.
		document.getElementById("input_desc").disabled = false;
		submit_save_form();
	}
}

// this function is called when the user presses the save button
function save_desc() {
	if (upload_started) {
		// if an upload has already been started that is finished already, we simply submit the description save form.
		if (percent == 100) {
			submit_save_form();
		} else {
			// otherwise, if the upload isn't finished yet, we temporarily disable the input textarea
			// (so that the user can't modify the description text anymore) and the save button (so that
			// the user can't press the button twice), and mark that a save operation is pending.
			// this will be honored by the finish_progress() function (see above).
			document.getElementById("btn_save").disabled = true;
			document.getElementById("input_desc").disabled = true;
			save_pending = true;
			document.getElementById("div_savemsg").innerHTML = "Uploading file and saving description...";
		}
	} else {
		// if no upload has been started yet, a message is shown instead that the user shall first
		// select a file to be uploaded.
		document.getElementById("div_savemsg").innerHTML = "Please choose a file to be uploaded.";
	}
}

// this function submit the description save form.
function submit_save_form() {
	document.forms["frm_save"].submit();
	// we need to reset the state variables just in case the user uses the back button to start another upload.
	// Chrome and IE8 fire a onload event, while Firefox and Safari don't. That's why we first submit the form,
	// and then clean up the previous state so that the page is ready for another upload if the user decides
	// to use the back button.
	reset_all();
}
</script>
</head>
<body onload="reset_all();">
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
