(function() {

  function Circle() {
    this._promptScene = document.getElementById('prompt-scene');
    this._antsScene = document.getElementById('ants-scene');
    this._uploadScene = document.getElementById('upload-scene');

    this._currentScene = Circle.PROMPT_SCENE;

    this._ants = new Ants();
  }

  Circle.PROMPT_SCENE = 0;
  Circle.ANTS_SCENE = 1;
  Circle.UPLOAD_SCENE = 2;

  Circle.prototype.switchScene = function(newScene) {
    if (newScene === this._currentScene) {
      return;
    }

    if (this._currentScene === Circle.ANTS_SCENE) {
      this._ants.stop();
    } else if (newScene === Circle.ANTS_SCENE) {
      this._ants.start();
    }

    this._elementForScene(this._currentScene).setAttribute('visibility', 'hidden');
    this._elementForScene(newScene).setAttribute('visibility', 'visible');
    this._currentScene = newScene;
  };

  Circle.prototype._elementForScene = function(scene) {
    return [this._promptScene, this._antsScene, this._uploadScene][scene];
  };

  function Ants() {
    this._frameRequest = null;
    this._startTime = null;
    this._element = document.getElementById('ants-scene');
  }

  Ants.prototype.start = function() {
    this._setAngle(0);
    this._frameRequest = window.requestAnimationFrame(this._frame.bind(this));
  };

  Ants.prototype.stop = function() {
    this._startTime = null;
    window.cancelAnimationFrame(this._frameRequest);
  };

  Ants.prototype._frame = function(time) {
    if (this._startTime === null) {
      this._startTime = time;
    }

    var angle = ((time - this._startTime) / 20) % 360;
    this._setAngle(angle);

    this._frameRequest = window.requestAnimationFrame(this._frame.bind(this));
  };

  Ants.prototype._setAngle = function(a) {
    this._element.setAttribute('transform', 'rotate(' + a.toFixed(3) + ', 0.5, 0.5)');
  };

  window.addEventListener('load', function() {
    window.circle = new Circle();
  });

})();
