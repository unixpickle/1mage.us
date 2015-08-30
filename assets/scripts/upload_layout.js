(function() {

  var CIRCLE_SIZE = 0.75;
  var CIRCLE_TOP = 0.05;

  function recomputeFontSize() {
    var height = window.innerHeight;
    var width = window.innerWidth;

    var circleSizeByHeight = height * CIRCLE_SIZE;
    var leftRightPadding = height * CIRCLE_TOP;
    var circleSizeByWidth = width - leftRightPadding*2;

    if (circleSizeByHeight < circleSizeByWidth) {
      document.body.style.fontSize = height.toFixed(3) + 'px';
    } else {
      var fontSize = height * (circleSizeByWidth / circleSizeByHeight);
      document.body.style.fontSize = fontSize.toFixed(3) + 'px';
    }
  }
  window.addEventListener('resize', recomputeFontSize);
  window.addEventListener('load', recomputeFontSize);

})();
