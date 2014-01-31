class UploadScene extends window.onemage.Scene
  constructor: -> super()
  
  activate: (arg) ->
    if not arg?
      throw new Error 'must be called with a file argument'
    @drawProgress 0
    $('#upload-scene').css display: 'block'
    $('#upload-box').css top: -350
    $('#upload-box').animate top: 50
    $('#upload-shielder').animate {opacity: 1}, =>
      window.onemage.Scene::activate.call this
      @startLoading arg
  
  deactivate: ->
    @active = false
    $('#upload-box').stop()
    $('#upload-shielder').stop()
    $('#upload-box').animate {top: -350}, ->
      $('#upload-scene').css display: 'none'
    $('#upload-shielder').animate opacity: 0
  
  includesURL: (url) -> false
  
  pushURL: (url) ->
    history.pushState {}, 'onemage.us', '/prog'

  startLoading: (file) ->
    formData = new FormData()
    formData.append file.name, file
    xhr = new XMLHttpRequest()
    xhr.open 'POST', '/upload', true
    xhr.addEventListener 'load', (e) ->
      value = JSON.parse xhr.response
      if value.error?
        window.onemage.scenes.error.activate()
      else
        window.onemage.scenes.choose.activate()
    xhr.addEventListener 'progress', (e) =>
      return unless e.lengthComputable
      percent = e.loaded / e.total
      @drawProgress percent
    xhr.send formData
  
  drawProgress: (prog) ->
    value = Math.round prog * 100
    $('#progress-value').text value + '%'
    loader = document.getElementById 'loader'
    context = loader.getContext '2d'
    context.clearRect 0, 0, loader.width, loader.height
    
    # outer gray circle
    context.beginPath()
    context.moveTo loader.width / 2, loader.height / 2
    context.arc loader.width / 2, loader.height / 2,
      loader.width / 2, 0, Math.PI * 2
    context.closePath()
    context.fillStyle = 'rgb(141,147,152)'
    context.fill()
    
    # outer circle
    x = Math.PI / 2
    context.beginPath()
    context.moveTo loader.width / 2, loader.height / 2
    context.arc loader.width / 2, loader.height / 2,
      loader.width / 2, 0 - x, Math.PI * 2 * prog - x
    context.closePath()
    context.fillStyle = '#00a9e9'
    context.fill()
    
    # inner circle
    context.beginPath()
    context.arc loader.width / 2, loader.height / 2,
      loader.width / 2 - 30, 0, Math.PI * 2
    context.closePath()
    context.fillStyle = '#FFF'
    context.fill()


window.onemage.scenes.upload = new UploadScene()