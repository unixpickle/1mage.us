db = require './db'

module.exports = (req, res) ->
  unless (match = /^\/nextlast\/([0-9]+)$/.exec req.url)?
    return res.redirect '/error'
  sequence = parseInt match[1]
  
  # make sure they aren't trying to pull something
  if isNaN sequence or sequence.toString() isnt match[1]
    return res.redirect '/error'
  
  # now, get our results
  db.findNextLast sequence, (err, next, last) ->
    return res.json error: err.toString() if err?
    obj = {}
    obj.next = next if next?
    obj.last = last if last?
    res.json obj
