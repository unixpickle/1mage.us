(function() {

  function Uploader() {
    this._chooser = new window.FileChooser();
    this._chooser.onFilesChosen = this._handleFiles.bind(this);
  }

  Uploader.prototype._handleFiles = function(files) {
    this._chooser.setEnabled(false);
    window.circle.switchScene(window.Circle.UPLOAD_SCENE);
    window.circle.setProgress(0);

    var upload = new Upload(files);
    upload.onDone = function() {
      // TODO: this.
      window.alert('TODO: this. (onDone)');
    }.bind(this);
    upload.onError = function(msg) {
      window.alert('error: ' + msg);
      this._chooser.setEnabled(true);
      window.circle.switchScene(window.Circle.PROMPT_SCENE);
    }.bind(this);
    upload.onProgress = function(p) {
      window.circle.setProgress(p);
    }.bind(this);

    // TODO: here, handle an authentication prompt.

    upload.start();
  };

  function Upload(files) {
    this._formData = new FormData();
    for (var i = 0, len = files.length; i < len; ++i) {
      this._formData.append(files[i].name, files[i]);
    }
    this._xhr = new XMLHttpRequest();

    this.onDone = null;
    this.onError = null;
    this.onProgress = null;
  }

  Upload.prototype.start = function() {
    this._registerEvents();
    this._xhr.open('POST', '/upload', true);
    this._xhr.send(this._formData);
  };

  Upload.prototype._handleError = function(errorStr) {
    this.onError(errorStr);
  };

  Upload.prototype._handleLoad = function() {
    var value;
    try {
      value = JSON.parse(this._xhr.response);
    } catch (e) {
      this._handleError('invalid response');
      return;
    }
    this._handleError('TODO: process the response to the upload call. ' + this._xhr.response);
  };

  Upload.prototype._handleProgress = function(e) {
    if (e.lengthComputable) {
      this.onProgress(e.loaded / e.total);
    }
  };

  Upload.prototype._registerEvents = function() {
    this._xhr.addEventListener('load', this._handleLoad.bind(this));
    this._xhr.addEventListener('error', this._handleError.bind(this, 'request failed'));
    this._xhr.addEventListener('progress', this._handleProgress.bind(this));
  };

  window.addEventListener('load', function() {
    window.uploader = new Uploader();
  });

})();
