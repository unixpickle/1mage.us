(function() {

  var DRAG_LEAVE_TIMEOUT = 10;

  function FileChooser() {
    this.onFilesChosen = null;

    this._enabled = true;
    this._input = document.getElementById('file-input');
    this._dragLeaveTimeout = null;
    this._registerEvents();
  }

  FileChooser.prototype.getEnabled = function() {
    return this._enabled;
  };

  FileChooser.prototype.setEnabled = function(e) {
    this._enabled = e;
    if (!e && this._dragLeaveTimeout !== null) {
      this._cancelDragLeaveTimeout();
    }
  };

  FileChooser.prototype._cancelDragLeaveTimeout = function() {
    clearTimeout(this._dragLeaveTimeout);
    this._dragLeaveTimeout = null;
  };

  FileChooser.prototype._dragLeave = function() {
    if (!this._enabled) {
      return;
    }
    if (this._dragLeaveTimeout !== null) {
      return;
    }
    this._dragLeaveTimeout = setTimeout(function() {
      this._dragLeaveTimeout = null;
      window.circle.switchScene(window.Circle.PROMPT_SCENE);
    }.bind(this), DRAG_LEAVE_TIMEOUT);
  };

  FileChooser.prototype._dragOver = function(e) {
    if (!this._enabled) {
      return;
    }

    if (e.dataTransfer) {
      e.dataTransfer.dropEffect = 'copy';
    }
    e.preventDefault();

    if (this._dragLeaveTimeout !== null) {
      this._cancelDragLeaveTimeout();
    } else {
      window.circle.switchScene(window.Circle.ANTS_SCENE);
    }
  };

  FileChooser.prototype._drop = function(e) {
    e.preventDefault();
    if (!this._enabled) {
      return;
    }
    if (e.dataTransfer.files.length === 0) {
      window.circle.switchScene(window.Circle.PROMPT_SCENE);
    } else {
      this._emitFiles(e.dataTransfer.files);
    }
  };

  FileChooser.prototype._emitFiles = function(files) {
    if ('function' !== typeof this.onFilesChosen) {
      throw new Error('no filesChosen listener');
    }
    this.onFilesChosen(files);
  };

  FileChooser.prototype._registerEvents = function() {
    document.body.addEventListener('dragover', this._dragOver.bind(this));
    document.body.addEventListener('dragleave', this._dragLeave.bind(this));
    document.body.addEventListener('drop', this._drop.bind(this));

    this._input.addEventListener('change', function() {
      this._emitFiles(this._input.files);
    }.bind(this));

    var circle = document.getElementById('circle');
    circle.addEventListener('click', function(e) {
      if (!this._enabled) {
        return;
      }
      var rect = circle.getBoundingClientRect();
      var centerX = rect.left + rect.width/2;
      var centerY = rect.top + rect.height/2;
      var distance = Math.sqrt(Math.pow(centerX-e.clientX, 2) + Math.pow(centerY-e.clientY, 2));
      if (distance <= rect.width/2) {
        this._input.click();
      }
    }.bind(this));
  };

  window.FileChooser = FileChooser;

})();
