configureDND = ->
  document.addEventListener 'dragenter', handleDragOver, false
  document.addEventListener 'dragleave', handleDragLeave, false
  document.addEventListener 'drop', handleFileSelect, false

handleDragOver = (evt) ->
  evt.stopPropagation()
  evt.preventDefault()
  evt.dataTransfer.dropEffect = 'copy'
  $('#upload-button').attr 'class', 'upload-button-drag'

handleDragLeave = (evt) ->
  $('#upload-button').attr 'class', 'upload-button-regular'

handleFileSelect = (evt) ->
  evt.stopPropagation()
  evt.preventDefault()
  output = []
  for f in evt.dataTransfer.files
    console.log f.name

if window.File? and window.FileReader? and window.FileList? and window.Blob?
  configureDND()
