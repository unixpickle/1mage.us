class Scene
  constructor: -> @active = false
  
  deactivate: -> @active = false
  
  activate: ->
    @active = true
    for _, obj of window.onemage.scenes
      continue if obj is this
      if obj.active
        obj.deactivate()
        break

  go: (args...) ->
    @activate args...
    @pushURL()

  urlTransition: ->

  pushURL: -> throw new Error 'Implement in subclass'
  
  includesURL: (url) -> false

load = (path) ->
  for own x, scene of window.onemage.scenes
    if scene.includesURL path
      if scene.active
        return scene.urlTransition()
      else return scene.activate()
  window.onemage.scenes.error.go()

window.onemage = Scene: Scene, scenes: {}, load: load
