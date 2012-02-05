if (typeof ak === 'undefined') {
	ak = {};
}

ak.upload = {
	upload_id: null,
	upload_started: false,
	percent: 0,
	save_pending: false,

	// method resets all state variables to their initial values and enables
	// the description textarea and the save button.
	reset_all: function() {
		ak.upload.upload_id = null;
		ak.upload.upload_started = false;
		ak.upload.percent = 0;
		ak.upload.save_pending = false;
		document.getElementById("input_desc").disabled = false;
		document.getElementById("btn_save").disabled = false;
	},

	// method saves the upload ID and resets the action URIs for
	// the two forms.
	set_upload_id: function(id) {
		ak.upload.upload_id = id;
		document.forms["frm_upload"].action = "/upload/" + ak.upload.upload_id;
		document.forms["frm_save"].action = "/savedesc/" + ak.upload.upload_id;
	},

	// method is called when the user selects a file
	start_upload: function() {
		// if we have no upload ID yet, we fetch one and rerun start_upload().
		if (!ak.upload.upload_id) {
			// the nocache?=random is there to randomize the URL. was added
			// to remedy browser caching that even wouldn't go away when the usual
			// HTTP response headers to prevent caching were set.
			get('/requpid?nocache=' + Math.random(), function(id) {
				ak.upload.set_upload_id(id);
				ak.upload.start_upload();
			});
			return;
		}
		// the user hasn't started an upload yet -> submit the form, and start
		// polling the upload progress.
		if (!ak.upload.upload_started) {
			document.forms["frm_upload"].submit();
			ak.upload.start_progress();
			ak.upload.upload_started = true;
		} else {
			// the user selected another file while an upload was already in progress.
			// we warn him that he has already started an upload.
			alert("You already started a download!");
		}
	},

	// method fetches the current upload progress, displays it, and if the
	// upload isn't complete yet, it schedules another upload progress poll.
	start_progress: function() {
		// see above. the ?nocache=random is there to prevent browser caching.
		get("/progress/" + ak.upload.upload_id + "?nocache=" + Math.random(), function(text) {
			var new_percent = parseInt(text);
			if (new_percent > ak.upload.percent) {
				ak.upload.percent = new_percent;
			}
			ak.upload.update_progress();
			// if the upload is finished, we completely finish it by showing
			// the user additional information
			if (ak.upload.percent == 100) {
				ak.upload.finish_progress();
			} else {
				// upload not yet finished -> schedule another poll
				window.setTimeout(ak.upload.start_progress, 1000);
			}
		});
	},

	// method displays the upload progress percents to the user while the upload
	// is going on.
	update_progress: function() {
		document.getElementById("progress").innerHTML = "Uploading... " + ak.upload.percent + "%";
	},

	// method completes the upload progress.
	finish_progress: function() {
		// we show the user that the upload has finished plus a download link to the uploaded file.
		document.getElementById("progress").innerHTML = 'Upload finished. <a href="/files/' + ak.upload.upload_id + '">Uploaded to here.</a>';
		// if the user pressed the save button in the mean while, we submit the description save form.
		if (ak.upload.save_pending) {
			// the previously disabled (see below) input textarea needs to be reenabled, otherwise
			// its content will not the sent as part of the form.
			document.getElementById("input_desc").disabled = false;
			ak.upload.submit_save_form();
		}
	},

		// method is called when the user presses the save button
	save_desc: function() {
		if (ak.upload.upload_started) {
			// if an upload has already been started that is finished already, we simply submit the description save form.
			if (ak.upload.percent == 100) {
				ak.upload.submit_save_form();
			} else {
				// otherwise, if the upload isn't finished yet, we temporarily disable the input textarea
				// (so that the user can't modify the description text anymore) and the save button (so that
				// the user can't press the button twice), and mark that a save operation is pending.
				// will be honored by the finish_progress() function (see above).
				document.getElementById("btn_save").disabled = true;
				document.getElementById("input_desc").disabled = true;
				ak.upload.save_pending = true;
				document.getElementById("div_savemsg").innerHTML = "Uploading file and saving description...";
			}
		} else {
			// if no upload has been started yet, a message is shown instead that the user shall first
			// select a file to be uploaded.
			document.getElementById("div_savemsg").innerHTML = "Please choose a file to be uploaded.";
		}
	},

	submit_save_form: function() {
		document.forms["frm_save"].submit();
		// we need to reset the state variables just in case the user uses the back button to start another upload.
		// Chrome and IE8 fire a onload event, while Firefox and Safari don't. That's why we first submit the form,
		// and then clean up the previous state so that the page is ready for another upload if the user decides
		// to use the back button.
		ak.upload.reset_all();
	}
}


