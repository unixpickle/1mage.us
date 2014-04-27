db = require './db'

module.exports = (req, res) ->
  db.latest (err, seq) ->
    return res.riderct '/error' if err?
    res.redirect '/nav/' + (seq - 1)

