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