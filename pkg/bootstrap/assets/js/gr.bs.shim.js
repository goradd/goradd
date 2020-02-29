/**
 * This file will make it so bootstrap javascript plays nicely with jquery ui. They have some name conflicts,
 * but bootstrap has a name resolution mechanism.
 *
 * After this file, to call the bootstrap versions, you would use:
 * .bootstrapButton()
 * .bootstrapTooltip()
 *
 * Rather than the built-in versions.
 *
 * The load sequence should be:
 * jquery ui
 * bootstrap
 * this file
 */

jQuery.fn.bootstrapButton = jQuery.fn.button.noConflict();
jQuery.fn.bootstrapTooltip = jQuery.fn.tooltip.noConflict();

// Shorthand for $( document ).ready()
$(function() {
    // Bootstrap fires events using jquery, so we have to capture it using jquery to shunt it to us.
    var ctrls = $('[data-grctl="bs-RadioListGroup"]');
    ctrls.on("change", "input", function(event) {
        g$(this).trigger("change");
    });
});
