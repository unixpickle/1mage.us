class ChooseScene extends window.onemage.Scene
  constructor: ->
    super()
    @_dispInt = @displayDragInterface.bind this
    @_hideInt = @hideDragInterface.bind this
    @_drop = @handleDrop.bind this
    @interfaceVisible = false
  
  activate: ->
    super()
    {File: a, FileReader: b, FileList: c, Blob: d} = window
    return unless a? and b? and c? and d?
    $(window).bind 'dragover', @_dispInt
    $('#dropzone').bind 'dragleave', @_hideInt
    $('#dropzone').bind 'drop', @_drop
    $('#choose-scene').css display: 'block'
  
  deactivate: ->
    super()
    $(window).unbind 'dragover', @_dispInt
    $('#dropzone').unbind 'dragleave', @_hideInt
    $('#dropzone').unbind 'drop', @_drop
    $('#choose-scene').css display: 'none'
  
  pushURL: -> history.pushState {}, '1mage.us', '/'
  
  includesURL: (url) ->
    url is '/' or url is ''
  
  displayDragInterface: (evt) ->
    evt.dataTransfer?.dropEffect = 'copy'
    evt.stopPropagation()
    evt.preventDefault()
    return if @interfaceVisible
    @interfaceVisible = true
    $('#drag-box').css display: 'block'
    $('#choose-button').css display: 'none'
    $('#dropzone').css display: 'block'
  
  hideDragInterface: (evt) ->
    evt.stopPropagation()
    evt.preventDefault()
    return unless @interfaceVisible
    @interfaceVisible = false
    $('#drag-box').css display: 'none'
    $('#choose-button').css display: 'block'
    $('#dropzone').css display: 'none'
  
  handleDrop: (evt) ->
    @hideDragInterface evt
    file = evt.dataTransfer.files[0]
    return if not file?
    window.onemage.scenes.upload.go file


window.onemage.scenes.choose = new ChooseScene()
