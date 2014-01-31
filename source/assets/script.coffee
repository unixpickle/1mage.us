onPopState = ->
  path = window.location.pathname
  window.onemage.load path

$ ->
  $.event.props.push 'dataTransfer'
  $(window).bind 'popstate', onPopState
  window.onemage.load window.location.pathname
