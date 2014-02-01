class GalleryScene extends window.onemage.Scene
  constructor: -> super()
  
  activate: (arg) ->
    if not arg?
      throw new Error 'must be called with sequence argument'
    super()
    $('#gallery-scene').css display: 'block'
  
  deactivate: ->
    super()
    $('#gallery-scene').css display: 'none'
  
  includesURL: (url) -> /^\/nav\/[0-9]+$/.exec(url)?
  
  pushURL: (url) ->
    history.pushState {}, '1mage.us', '/nav'


window.onemage.scenes.gallery = new GalleryScene()
