var api = "http://localhost:8080"
var app = angular.module("app", ["ui.router", "ngResource", "ui.bootstrap"]);

app.config(function($stateProvider, $urlRouterProvider) {
  //
  // For any unmatched url, redirect to /state1
  $urlRouterProvider.otherwise("/teams");
  //
  // Now set up the states
  $stateProvider
    .state('team', {
      url: "/team/:teamId",
      templateUrl: "partials/team.html",
      "controller": "TeamController"
    })
    .state('teams', {
      url: "/teams",
      templateUrl: "partials/teams.html",
      controller: "TeamsController"
    })
    .state('login', {
      url: "/login",
      templateUrl: "partials/login.html",
      "controller": "LoginController"
    })
});

app.service("User", function($resource) {
  return $resource(api + "/api/me", {}, {login: {method: "POST", url: "/api/users/login"}})
})

app.service("Team", function($resource) {
  return $resource(api + "/api/team/:teamId", {teamId: "@ID"}, {update: {method: "PUT"}})
})

app.service("Member", function($resource) {
  return $resource(api + "/api/team/:teamId/members")
})

app.service("Chat", function($resource) {
  return $resource(api + "/api/team/:teamId/chats", {teamId: "@Team"})
})

app.service("Attachment", function($resource) {
  return $resource(api + "/api/team/:teamId/attachments", {teamId: "@Team"})
})

app.service("Event", function($resource) {
  return $resource(api + "/api/team/:teamId/events", {teamId: "@Team"})
})

app.controller('TeamsController', function($scope, Team) {
  $scope.teams = Team.query()
})

app.controller('LoginController', function($scope, User, $state) {
  $scope.user = new User()

  $scope.login = function() {
    $scope.user.$login(function() {
      $state.go('teams')
    })
  }
})

app.controller('TeamController', function($scope, Member, User, Team, Chat, Attachment, Event, $stateParams) {
  var teamId = $stateParams.teamId

  $scope.user = User.get();
  $scope.team = Team.get({teamId: teamId});
  $scope.chats = Chat.query({teamId: teamId});
  $scope.attachments = Attachment.query({teamId: teamId})
  $scope.events = Event.query({teamId: teamId})
  $scope.members = Member.query({teamId: teamId});

  $scope.saveChat = function() {
    var chat = new Chat({Message: $scope.newMessage, Team: teamId})
    chat.$save(function() {
      $scope.chats = Chat.query({teamId: teamId})
      $scope.newMessage = "";
    })
  }

  $scope.saveTeam = function() {
    console.log("SWAG")
    $scope.editTeam = false;
    $scope.team.$update();
  }
});

app.filter('fromNow', function() {
  return function(date) {
    return moment(date).fromNow();
  }
});

app.filter('to_trusted', ['$sce', function($sce){
  return function(text) {
    return $sce.trustAsHtml(text);
  };
}]);
