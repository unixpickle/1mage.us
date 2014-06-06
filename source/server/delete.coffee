db = require './db'
args = require './args'
path = require 'path'

get = (req, res) ->
  res.sendfile path.resolve __dirname + '/../../static/delete.html'

post = (req, res) ->
  if req.body.password isnt args.password
    return res.redirect '/error'
  db.delete parseInt(req.body.seq), (err) ->
    return res.redirect '/error' if err?
    res.redirect '/'

module.exports =
  post: post
  get: get

