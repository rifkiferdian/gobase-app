/*
Template Name: Skote - Admin & Dashboard Template
Author: Themesbrand
Website: https://themesbrand.com/
Contact: themesbrand@gmail.com
File: Material design Init Js File
*/

var stockOutCaseModal = document.getElementById('stockOutCaseModal');
if (stockOutCaseModal) {
    stockOutCaseModal.addEventListener('show.bs.modal', function () {
        var form = stockOutCaseModal.querySelector('#stockOutCaseForm');
        if (form) {
            form.reset();
            form.classList.remove('was-validated');
        }
    });

    var form = stockOutCaseModal.querySelector('#stockOutCaseForm');
    if (form) {
        form.addEventListener('submit', function (event) {
            if (!form.checkValidity()) {
                event.preventDefault();
                event.stopPropagation();
            }
            form.classList.add('was-validated');
        });
    }
}
