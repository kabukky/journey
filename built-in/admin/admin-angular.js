//register the modules (don't forget ui bootstrap and bootstrap switch)
var adminApp = angular.module('adminApp', ['ngRoute', 'frapontillo.bootstrap-switch', 'ui.bootstrap', 'infinite-scroll']);
adminApp.config(function($routeProvider) {
	$routeProvider.
		when('/', {
			templateUrl: 'content.html',
			controller: 'ContentCtrl'
		}).
		when('/edit/:Id', {
			templateUrl: 'post.html',
			controller: 'EditCtrl'
		}).
		when('/create/', {
			templateUrl: 'post.html',
			controller: 'CreateCtrl'
		}).
    when('/settings/', {
      templateUrl: 'settings.html',
      controller: 'SettingsCtrl'
    }).
		otherwise({
        	redirectTo: '/'
		});
});

//service for sharing the markdown content across controllers
adminApp.factory('sharingService', function(){

  return { shared: { post: {}, blog: {}, user: {}, infiniteScrollFactory: null, selected: "" } }

});

//directive to handle visual selection of images
adminApp.directive('imgSelectionDirective', function() {
    return {
        restrict: 'A',
        link: function(scope, elem, attrs) {
            $(elem).click(function() {
              $('.imgselected').removeClass('imgselected');
              $(elem).children('img').addClass('imgselected');
        });
      }
    }
});

//filter to convert html to plain text
adminApp.filter('htmlToPlaintextPreview', function() {
    return function(text) {
      //remove html tags 
      var output = String(text).replace(/<[^>]+>/gm, '');
      //trim string (75 characters for now)
      if (output.length > 300) {
        output = output.substr(0, 300);
        //re-trim if in the middle of a word
        output = output.substr(0, Math.min(output.length, output.lastIndexOf(" ")));
      }
      return output;
    }
  }
);

//factory to load items in infinite-scroll
adminApp.factory('infiniteScrollFactory', function($http) {
  var infiniteScrollFactory = function(url) {
    this.url = url;
    this.items = [];
    this.busy = false;
    this.after = 1;
  };

  infiniteScrollFactory.prototype.nextPage = function() {
    if (this.busy) return;
    this.busy = true;
    var url = this.url + this.after;
    $http.get(url).success(function(data) {
      var items = data;
      for (var i = 0; i < items.length; i++) {
        this.items.push(items[i]);
      }
      this.after = this.after + 1;
      if (data.length > 0) this.busy = false;
    }.bind(this));
  };
  return infiniteScrollFactory;
});

adminApp.controller('ContentCtrl', function ($scope, $http, $sce, $location, infiniteScrollFactory, sharingService){
  //change the navbar according to controller
  $scope.navbarHtml = $sce.trustAsHtml('<ul class="nav navbar-nav"><li class="active"><a href="#/">Content<span class="sr-only">(current)</span></a></li><li><a href="#/create/">New Post</a></li><li><a href="#/settings/">Settings</a></li></ul>');
  $scope.infiniteScrollFactory = new infiniteScrollFactory('/admin/api/posts/');
  $scope.openPost = function(postId) {
    $location.url('/edit/' + postId);
  };
  $scope.deletePost = function(postId, postTitle) {
    if (confirm('Are you sure you want to delete post "' + postTitle + '"?')) {
      $http.delete('/admin/api/post/' + postId).success(function(data) {
        //delete post from array
        for (var i = 0; i < $scope.infiniteScrollFactory.items.length; i++) {
          if($scope.infiniteScrollFactory.items[i].Id == postId) {
            $scope.infiniteScrollFactory.items.splice(i, 1);
          }
        }
      });
    }
  };
});

adminApp.controller('SettingsCtrl', function ($scope, $http, $timeout, $sce, $location, sharingService){
  //change the navbar according to controller
  $scope.navbarHtml = $sce.trustAsHtml('<ul class="nav navbar-nav"><li><a href="#/">Content</a></li><li><a href="#/create/">New Post</a></li><li class="active"><a href="#/settings/">Settings<span class="sr-only">(current)</span></a></li></ul>');
  $scope.shared = sharingService.shared;
  $scope.loadData = function() {
    $http.get('/admin/api/blog/').success(function(data) {
      $scope.shared.blog = data;
      var themeIndex = $scope.shared.blog.Themes.indexOf($scope.shared.blog.ActiveTheme);
      $scope.shared.blog.ActiveTheme = $scope.shared.blog.Themes[themeIndex];
    });
    $http.get('/admin/api/userid/').success(function(data) {
      $scope.authenticatedUser = data;
      $http.get('/admin/api/user/' + $scope.authenticatedUser.Id).success(function(data) {
        $scope.shared.user = data;
      });
    });
  };
  $scope.loadData();
  $scope.save = function() {
    $http.patch('/admin/api/blog/', $scope.shared.blog);
    $http.patch('/admin/api/user/', $scope.shared.user).success(function(data) {
      $location.url('/');
    });
  };
});

adminApp.controller('CreateCtrl', function ($scope, $http, $sce, $location, sharingService){
  //create markdown converter
  var converter = new Showdown.converter();
  //change the navbar according to controller
  $scope.navbarHtml = $sce.trustAsHtml('<ul class="nav navbar-nav"><li><a href="#/">Content</a></li><li class="active"><a href="#/create/">New Post<span class="sr-only">(current)</span></a></li><li><a href="#/settings/">Settings</a></li></ul>');
	$scope.shared = sharingService.shared;
  $scope.shared.post = {Title: 'New Post', Slug: '', Markdown: 'Write something!', IsPublished: false, Image: '', Tags: ''}
	$scope.change = function() {
    document.getElementById('html-div').innerHTML = '<h1>' + $scope.shared.post.Title + '</h1><br>' + converter.makeHtml($scope.shared.post.Markdown);
    //resize the markdown textarea
    $('.textarea-autosize').val($scope.shared.post.Markdown).trigger('autosize.resize');
  };
  $scope.change();
  $scope.save = function() {
    $http.post('/admin/api/post/', $scope.shared.post).success(function(data) {
      $location.url('/');
    });
  };
});

adminApp.controller('EditCtrl', function ($scope, $routeParams, $http, $sce, $location, sharingService){
  //create markdown converter
  var converter = new Showdown.converter();
  //change the navbar according to controller
  $scope.navbarHtml = $sce.trustAsHtml('<ul class="nav navbar-nav"><li><a href="#/">Content</a></li><li><a href="#/create/">New Post</a></li><li><a href="#/settings/">Settings</a></li></ul>');
  $scope.shared = sharingService.shared;
  $scope.shared.post = {}
  $scope.change = function() {
    document.getElementById('html-div').innerHTML = '<h1>' + $scope.shared.post.Title + '</h1><br>' + converter.makeHtml($scope.shared.post.Markdown);
    //resize the markdown textarea
    $('.textarea-autosize').val($scope.shared.post.Markdown).trigger('autosize.resize');
  };
	$http.get('/admin/api/post/' + $routeParams.Id).success(function(data) {
		$scope.shared.post = data;
    $scope.change();
  });
  $scope.save = function() {
  	$http.patch('/admin/api/post/', $scope.shared.post).success(function(data) {
      $location.url('/');
    });
  };
});

//modal for post options and help
adminApp.controller('EmptyModalCtrl', function ($scope, $modal, $http, sharingService) {
  $scope.shared = sharingService.shared;
  $scope.deleteCover = function() {
    $scope.shared.post.Image = '';
  };
  $scope.open = function (size, callingFrom) {
    if (callingFrom == 'post-options') {
      var modalInstance = $modal.open({
        templateUrl: 'post-options-modal.tpl',
        controller: 'EmptyModalInstanceCtrl',
        size: size
      });
    } else if (callingFrom == 'post-help') {
      var modalInstance = $modal.open({
        templateUrl: 'post-help-modal.tpl',
        controller: 'EmptyModalInstanceCtrl',
        size: size
      });
    }
    modalInstance.result.then(function (selectedItem) {
      if (callingFrom == 'post-cover') {
        $scope.selected = selectedItem;
        $scope.shared.post.Image = $scope.selected;
      }
    });
  };
});

//modal for image selection and upload
adminApp.controller('ImageModalCtrl', function ($scope, $modal, $http, sharingService, infiniteScrollFactory) {
	$scope.shared = sharingService.shared;
  $scope.shared.infiniteScrollFactory = new infiniteScrollFactory('/admin/api/images/');
  $scope.shared.infiniteScrollFactory.nextPage();
  $scope.open = function (size, callingFrom) {
    //image modal
    var modalInstance = $modal.open({
      templateUrl: 'image-modal.tpl',
      controller: 'ImageModalInstanceCtrl',
      size: size
    });
    modalInstance.result.then(function (selectedItem) {
      if (callingFrom == 'post-image') {
        $scope.selected = selectedItem;
        $scope.shared.post.Markdown = $scope.shared.post.Markdown + '\n\n![](' + $scope.selected + ')\n\n';
        $scope.change(); //we invoke the change function of CreateCtrl or EditCtrl to update the html view.
      } else if (callingFrom == 'post-cover') {
        $scope.selected = selectedItem;
        $scope.shared.post.Image = $scope.selected;
      } else if (callingFrom == 'blog-logo') {
        $scope.selected = selectedItem;
        $scope.shared.blog.Logo = $scope.selected;
      } else if (callingFrom == 'blog-cover') {
        $scope.selected = selectedItem;
        $scope.shared.blog.Cover = $scope.selected;
      } else if (callingFrom == 'user-image') {
        $scope.selected = selectedItem;
        $scope.shared.user.Image = $scope.selected;
      } else if (callingFrom == 'user-cover') {
        $scope.selected = selectedItem;
        $scope.shared.user.Cover = $scope.selected;
      }
    });
  };
});

//empty modal window instance
adminApp.controller('EmptyModalInstanceCtrl', function ($scope, $http, $modalInstance, sharingService) {
  $scope.shared = sharingService.shared;
  $scope.ok = function () {
    $modalInstance.close();
  };
  $scope.cancel = function () {
    $modalInstance.dismiss('cancel');
  };
});

//image modal window instance
adminApp.controller('ImageModalInstanceCtrl', function ($scope, $http, $modalInstance, sharingService) {
  $scope.shared = sharingService.shared;
  $scope.shared.selected = $scope.shared.infiniteScrollFactory.items[0];
  $scope.ok = function () {
    $modalInstance.close($scope.shared.selected);
  };
  $scope.cancel = function () {
    $modalInstance.dismiss('cancel');
  };
});