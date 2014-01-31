path = require 'path'

db = require './db'
args = require './args'

###
Should be used as the handler for /*
###
module.exports = (req, res) ->
  unless (match = /^\/([0-9]+)$/.exec req.url)?
    return res.sendHome()
  
  sequence = parseInt match[1]
  if isNaN sequence or sequence.toString() isnt match[1]
    return res.redirect '/error'
  
  db.fetch sequence, (err, doc) ->
    return res.redirect '/error' if err? or not doc?
    thePath = path.join args.directory, sequence.toString()
    res.contentType doc.mime
    res.sendfile thePath
