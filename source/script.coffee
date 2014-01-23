dragInterfaceVisible = false
scene = 'upload'

configureDND = ->
  $(document).bind 'dragover', displayDragInterface
  $(document).bind 'dragleave', removeDragInterface
  $(document).bind 'drop', handleFileSelect

displayDragInterface = (evt) ->
  return unless scene is 'upload'
  evt.dataTransfer?.dropEffect = 'copy'
  evt.stopPropagation()
  evt.preventDefault()
  return if dragInterfaceVisible
  dragInterfaceVisible = true
  $('#drag-box').css display: 'block'
  $('#upload-button').css display: 'none'
  $('#upload-scene').css 'pointer-events': 'none'

removeDragInterface = (evt) ->
  return unless scene is 'upload'
  evt.stopPropagation()
  evt.preventDefault()
  return unless dragInterfaceVisible
  dragInterfaceVisible = false
  $('#drag-box').css display: 'none'
  $('#upload-button').css display: 'block'
  $('#upload-scene').css 'pointer-events': 'auto'

handleFileSelect = (evt) ->
  return unless scene is 'upload'
  removeDragInterface evt
  output = []
  uploadFile evt.dataTransfer.files[0]

uploadFile = (file) ->
  return unless scene is 'upload'
  formData = new FormData()
  formData.append file.name, file
  xhr = new XMLHttpRequest()
  xhr.open 'POST', '/upload', true
  xhr.onload = (e) ->
    value = JSON.parse xhr.response
    return window.displayError() if value.error?
    window.displayUploaded value.identifier
  xhr.send formData

window.displayError = (transition = true) ->
  scene = 'error'
  if transition
    history.pushState {}, 'Error', '/error'
  $('#upload-scene').css display: 'none'
  $('#gallery-scene').css display: 'none'
  $('#error-scene').css display: 'block'

window.displayUploaded = (imageId, transition = true) ->
  scene = 'uploaded'
  if transition
    history.pushState {}, 'Image ' + imageId, '/nav/' + imageId
  $('#upload-scene').css display: 'none'
  $('#gallery-scene').css display: 'block'
  $('#error-scene').css display: 'none'

window.displayHome = ->
  history.pushState {}, '1mage.us', '/'
  $('#upload-scene').css display: 'block'
  $('#gallery-scene').css display: 'none'
  $('#error-scene').css display: 'none'

onPopState = ->
  loadStateWithURL true

loadStateWithURL = (transition = true) ->
  path = window.location.pathname
  if (match = /\/nav\/([0-9]+)/.exec path)?
    window.displayUploaded parseInt(match[1]), transition
  else if path is '/error'
    window.displayError transition
  else
    window.displayHome transition

$ ->
  jQuery.event.props.push 'dataTransfer'
  if window.File? and window.FileReader? and window.FileList? and window.Blob?
    configureDND()
  $(window).bind 'popstate', onPopState
  loadStateWithURL false
