var booksApp = angular.module('booksApp', []);

booksApp.controller('BookListCtrl', function($scope, $http) {
  $http.get('api/book').success(function(data) {
      $scope.books = data;
    });
});
