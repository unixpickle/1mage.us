class GalleryNode
  constructor: (@seq) ->
    @nextLast = null
    @img = new Image()
    @img.onload = => @_imageLoaded()
    @img.src = '/' + @seq
    @size = null
    @callback = null
    @error = null
    @req = $.ajax
      dataType: 'json'
      url: '/nextlast/' + @seq
      success: (obj) => @_nextLastLoaded obj
      error: => @_errorEncountered 'failed to load AJAX'
  
  load: (handler) ->
    if @_isCompleted()
      if not @size?
        @size = width: @img.width, height: @img.height
      handler()
    else @callback = handler
  
  cancel: ->
    @callback = null
    delete @img.onload
    @img.src = null
    @req?.abort?()
    @req = null
  
  _imageLoaded: ->
    @size = width: @img.width, height: @img.height
    if @_isCompleted()
      @callback?()
      @callback = null
  
  _nextLastLoaded: (nl) ->
    @req = null
    if nl.error?
      @error = nl.error
    else @nextLast = nl
    if @_isCompleted()
      @callback?()
      @callback = null
  
  _errorEncountered: (e) ->
    @req = null
    @error = e
    if @_isCompleted()
      @callback?()
      @callback = null
  
  _isCompleted: ->
    val = @error? or (@img.complete and @nextLast?)
    if @img.complete and not @size?
      @size = width: @img.width, height: @img.height
    return val


class GalleryScene extends window.onemage.Scene
  constructor: ->
    super()
    @current = null
    @next = null
    @last = null
    
    @imageTag = null
    @_handleResize = @resizeImage.bind this
    @_handleKeyPress = @handleKeyPress.bind this
  
  activate: (arg) ->
    if not arg?
      url = window.location.pathname
      unless (match = /^\/nav\/([0-9]+)$/.exec url)?
        throw new Error 'must be called with sequence argument'
      arg = parseInt match[1]
    
    super()
    $('#gallery-scene').css display: 'block'
    @destroyContext()
    @loadNode new GalleryNode arg
    $(window).bind 'resize', @_handleResize
    $(window).bind 'keydown', @_handleKeyPress
  
  deactivate: ->
    super()
    @destroyContext()
    $('#image-well').html ''
    $('#gallery-scene').css display: 'none'
    $(window).unbind 'resize', @_handleResize
    $(window).unbind 'keydown', @_handleKeyPress
  
  urlTransition: ->
    url = window.location.pathname
    unless (match = /^\/nav\/([0-9]+)$/.exec url)?
      throw new Error 'must be called with sequence argument'
    arg = parseInt match[1]
    
    $('#image-well').stop()
    @destroyContext()
    @loadNode new GalleryNode arg
  
  includesURL: (url) -> /^\/nav\/[0-9]+$/.exec(url)?
  
  pushURL: ->
    history.pushState {}, '1mage.us', '/nav/' + @current.seq
  
  destroyContext: ->
    @current?.cancel?()
    @next?.cancel?()
    @last?.cancel?()
    [@current, @next, @last] = [null, null, null]
  
  switchToNode: (node) ->
    @pushURL()
    $('#image-well').stop()
    $('#image-well').animate {opacity: 0}, 'fast', =>
      @loadNode node
  
  loadNode: (node) ->
    @current = node
    @current.load =>
      if @current.nextLast.next?
        @next = new GalleryNode @current.nextLast.next
      if @current.nextLast.last?
        @last = new GalleryNode @current.nextLast.last
      # display the image
      @_displayCurrentImage()
  
  _displayCurrentImage: ->
    return if not @current
    $('#image-well').html ''
    $('#image-well').append @current.img
    $('#image-well').animate {opacity: 1}, 'fast'
    @resizeImage()

  resizeImage: ->
    return if not @current?.size?
    width = $(window).width()
    height = $(window).height()
    if width < 200 or height < 200
      return $('#image-well').css width: 0, height: 0
    width -= 200
    height -= 200
    
    widthRat = width / @current.size.width
    heightRat = height / @current.size.height
    scalar = if widthRat > heightRat then heightRat else widthRat
    newWidth = scalar * @current.size.width
    newHeight = scalar * @current.size.height
    
    sel = $ '#image-well img'
    sel.css width: newWidth, height: newHeight
  
  loadNext: ->
    return if not @next?
    [@last, @current, @next] = [@current, @next, null]
    @switchToNode @current
  
  loadLast: ->
    return if not @last?
    [@last, @current, @next] = [null, @last, @current]
    @switchToNode @current
  
  handleKeyPress: (evt) ->
    evt.preventDefault()
    if evt.keyCode is 37
      @loadNext()
    else if evt.keyCode is 39
      @loadLast()


window.onemage.scenes.gallery = new GalleryScene()
