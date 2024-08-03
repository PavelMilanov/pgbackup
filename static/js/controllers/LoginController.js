App.controller('LoginController', function ($scope) {
    $scope.user = {}
    $scope.save = function (user, loginForm) {
        if (loginForm.$valid) {
            // Handle form submission logic here
            console.log('User:', user);
        } else {
            console.log('Form is not valid');
        }
    };
})