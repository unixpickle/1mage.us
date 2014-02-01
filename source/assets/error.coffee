class ErrorScene extends window.onemage.Scene
  constructor: -> super()
  
  activate: (msg = 'Woops, looks like something went wrong') ->
    super()
    $('#error-scene').css display: 'block'
    $('#error-message').text msg
  
  deactivate: ->
    super()
    $('#error-scene').css display: 'none'
  
  includesURL: (url) -> url == '/error'
  
  pushURL: -> history.pushState {}, '1mage.us', '/error'

window.onemage.scenes.error = new ErrorScene()
