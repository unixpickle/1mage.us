db = require './db'

module.exports = (req, res) ->
  db.latest (err, seq) ->
    return res.redirect '/error' if err?
    res.redirect '/nav/' + seq

