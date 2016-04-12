var ENDPOINT = "http://localhost:8080";

var me = {};
$.getJSON(ENDPOINT + "/api/me", function(results) {
  me = results;
  console.log(me)
});

var data = {User: 0, Author: "Stuart", Type: "visit", Data: document.URL};
$.post("http://localhost:8080/api/team/events", JSON.stringify(data), function() {}, "json")

jQuery.fn.getPath = function () {
    if (this.length != 1) throw 'Requires one element.';

    var path, node = this;
    while (node.length) {
        var realNode = node[0], name = realNode.localName;
        if (!name) break;
        name = name.toLowerCase();

        var parent = node.parent();

        var siblings = parent.children(name);
        if (siblings.length > 1) { 
            name += ':eq(' + siblings.index(realNode) + ')';
        }

        path = name + (path ? '>' + path : '');
        node = parent;
    }

    return path;
};

function addComment(parent, body, path) {
  console.log("[+] addComment", parent, body, path)

  var url = document.URL;
  var data = {Url: url, Body: body, Path: path, Parent: parent};
  $.post(ENDPOINT + "/api/team/comments", JSON.stringify(data), function() {}, "json")
}

$( "*" ).on( "dblclick", function(event) {
  event.stopImmediatePropagation()

  var path = $( this ).getPath();
  var comment = prompt("Enter your comment", "");
  if (comment == null || comment == "") {
    return
  }

  addComment(null, comment, path);
});


function buildComments(x, y, comments, me) {
  out ='<ul id="markit-comment" class="list-group" style=" z-index: 999; width: 400px; position: absolute; top: '+x+'px; left: '+y+'px">'

  for (var i = 0; i < comments.length; i++) {
    var comment = comments[i];
    out += '<li class="list-group-item">'
    if (i == 0) {
      out += '<button id="markit-close" type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>' 
    }
    out += '<img class="pull-left" src="'+ comment.ProfilePic + '" style="width: 54px; height: 54px; margin-right: 10px">' +
    '<p>'+comment.Author+'<small> '+comment.TS+'</small></p>'+
    '<p>'+comment.Body+'</p>'+
    '</li>'
  }
  return out +
  // <li class="list-group-item">
  //   <img class="pull-left" src="https://static.wixstatic.com/media/a905fd_76d88cea66d64e56a47221b59b381f21.jpg/v1/fill/w_870,h_836,al_c,q_85,usm_0.66_1.00_0.01/a905fd_76d88cea66d64e56a47221b59b381f21.jpg" style="width: 54px; height:54 px; margin-right: 10px">
  //   <p>Katie Honadle<small> 30 minutes ago</small></p>
  //   <p>That looks awesome! Could we get the title in purple?</p>  
  // </li>
  '<li class="list-group-item" style="min-height: 64px">'+
      '<img class="pull-left" src="'+ me.ProfilePic +'" style="width:54px; height:54 px; margin-right: 10px">'+
      '<div class="input-group">' +
        '<input id="markit-input" type="text" class="form-control" placeholder="Message">' +
        '<span class="input-group-btn">' +
        '<button id="markit-add" class="btn btn-default" type="button">Reply</button>' +
        '</span>' +
      '</div><!-- /input-group -->' +
     '<span class="clearfix"></span>' +
  '</li>' +
'</ul>'
}

$.getJSON("http://localhost:8080/api/team/comments", function( comments ) {
  console.log("Get comments", comments)
  for (var i = 0; i < comments.length; i++) {
    var comment = comments[i];
    console.log(comment.Parent)
    if (comment.Parent != 0) {
      continue
    }

    bindComment(comment, comments);
  }
});

function bindComment(comment, comments) {
  $(comment.Path).css("border", "dotted")
  $(comment.Path).css("border-color", "red")


  $(comment.Path).hover(function(e) {
    if ($('#markit-comment').size() != 0) {
      return
    }

    commentTree = buildCommentTree(comment, comments);
    popupComment(e.pageX, e.pageY, commentTree);
  });
}

function popupComment(x, y, commentTree) {
  var popup = buildComments(y, x, commentTree, me)
  $('body').append(popup)

  $('#markit-close').click(function() {
    $('#markit-comment').remove()
  })

  $('#markit-comment').hover(function(){}, function() {
    $('#markit-comment').remove()
  })

  $('#markit-add').click(function() {
    var body = $('#markit-input').val();
    if (body == "") {
      return
    }
    var last = commentTree[commentTree.length - 1]
    addComment(last.ID, body, last.Path);
    $('#markit-input').val("");

    var newComment = {Parent: last.ID, Body: body, Path: last.Path, ProfilePic: me.ProfilePic, Author: me.Name, TS: new Date()}
    commentTree.push(newComment)

    $('#markit-comment').remove()
    popupComment(x, y, commentTree)

  })
}

function buildCommentTree(parent, comments) {
  out = [parent]
  var i = 0;
  while (i < comments.length) {
    var comment = comments[i];
    if (comment.Parent == parent.ID) {
      out.push(comment)
      parent = comment
      i = 0
    } else {
      i += 1
    }
  }
  return out
}

