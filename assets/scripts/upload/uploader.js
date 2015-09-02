(function() {

  function Uploader() {
    this._chooser = new window.FileChooser();
    this._chooser.onFilesChosen = this._handleFiles.bind(this);
  }

  Uploader.prototype._handleFiles = function(files) {
    this._chooser.setEnabled(false);
    window.circle.switchScene(window.Circle.PROMPT_SCENE);
  };

  function Upload(files) {
    var formData = new FormData();
    for (var i = 0, len = files.length; i < len; ++i) {
      formData.append(files[i].name, files[i]);
    }
    this._xhr = new XMLHttpRequest();

    this.onDone = null;
    this.onError = null;
    this.onProgress = null;
  }

  Upload.prototype._handleError = function(errorStr) {
    this.onError();
  };

  Upload.prototype._handleLoad = function() {
    var value;
    try {
      value = JSON.parse(this._xhr.response);
    } catch (e) {
      this._handleError('invalid response');
      return;
    }
    if (value.error) {
      this._handleError(value.error);
    } else {
      this.onDone();
    }
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
