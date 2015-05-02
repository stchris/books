
$(document).ready(function() {
    $('#addBookForm')
        .on('success.form.bv', function(e) {
            // Get the form instance
            var $form = $(e.target);

            // Get the BootstrapValidator instance
            var bv = $form.data('bootstrapValidator');

            // Use Ajax to submit form data
            $.post(
              $form.attr('action'),
              $form.serialize(),
              'json'
            );
        });
});
