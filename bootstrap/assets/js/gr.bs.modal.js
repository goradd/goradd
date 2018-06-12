/**
 * Helper to implement a modal via bootstrap.
 *
 * Uses the jQuery UI widget mechanism to provide object oriented functionality.
 *
 * Assumes that a wrapper is being used on the control and is set up with a modal class. This is required for correct functioning.
 */

jQuery(function( $, undefined ) {
    $.widget( "qcubed.bsModal", {
        options: {
            fade:    true,          // Modal transition
            hasCloseButton: true,   // Has a close button
            title: "",              // Title text
            headerClasses: "",      // Classes to put in the header. A good one is bg-primary, bg-success, etc.
            buttons: null,          // Array of button options
            size: null,             // The modal-lg or modal-sm option
            backdrop: null,         // The backdrop option that can be specified in the initialization options.
                                    // Boolean, or the string "static", which means do not allow closing by clicking outside of dialog.
            keyboard: true,         // Boolean, whether to allow ESC key to close
            show: false             // Boolean, whether to show immediately upon initialization
        },
        _create: function() {
            var self = this,
                $control = this.element,
                id = $control.attr('id');

            $control.addClass("modal");
            $control.attr('tabindex', -1);



            var $md = $('<div class="modal-dialog" role="document"></div>');
            if (this.options.size) {
                $md.addClass(this.options.size);
            }

            // allows capturing of enter key events
            $md.attr('tabindex', -1);

            var $mc = $('<div class="modal-content"></div>');
            var $mh = $('<div class="modal-header"></div>');
            var hasHeader = false;
            if (this.options.hasCloseButton) {
                $mh.append('<button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>');
                hasHeader = true;
            }
            if (this.options.title) {
                $mh.append('<h4 class="modal-title" id="gridSystemModalLabel">' + this.options.title + '</h4>');
                hasHeader = true;
            }
            if (this.options.headerClasses) {
                $mh.addClass(this.options.headerClasses);
            }
            $control.addClass('modal-body');

            var $mf = $('<div class="modal-footer"></div>');
            if (this.options.buttons) {
                $.each(this.options.buttons, function(i, objButton) {
                    var buttonStyle = 'btn-default';
                    var $b = $('<button type="button">' + objButton.label + '</button>');
                    if (objButton.primary) {
                        $b.attr('data-primary', 1);
                    }

                    if (objButton.style) {
                        buttonStyle = objButton.style;
                    }
                    $b.addClass(buttonStyle);

                    if (objButton.size) {
                        $b.addClass(objButton.size);
                    }

                    $b.attr('data-btnid', objButton.id);

                    $b.addClass('btn');

                    // create an action for the button
                    if (objButton.close) {
                        // button click closes the dialog
                        $b.attr("data-dismiss", "modal");
                    }

                    objButton.instance = self;

                    $b.click(objButton, self._buttonClick);

                    $mf.append($b);
                });
            }

            // makes sure a return key fires the default button if there is one
            $control.parent().on ("keydown", function(event) {
                if (event.which === 13 && !$(event.target).is('textarea')) {
                    var b = $(this).closest("[role=\'dialog\']").find("button[data-primary=1]");
                    if (b && b[0]) {
                        b[0].click();
                    }
                    event.preventDefault();
                    return false;
                }
            });


            // put it all together
            $md.append($mc);
            if (hasHeader) {
                $mc.append($mh);
            }
            $mc.append($control);
            if (this.options.buttons) {
                $mc.append($mf);
            }
            $wrapper.append($md);

            // initialize the modal
            var options = {};
            if (this.options.keyboard) {
                options.keyboard = this.options.keyboard;
            }
            if (this.options.backdrop) {
                options.backdrop = this.options.backdrop;
            }
            options.show = this.options.show;
            $wrapper.modal(options);

            $wrapper.on('shown.bs.modal', function (e) {
                qcubed.recordControlModification(id, "_IsOpen", true);
                // focus first element if possible, or whole dialog to capture return event
                var obj = $control.find(':input:enabled:visible:first');
                if (obj.length) {
                    obj.focus();
                } else {
                    $md.focus(); // allows capturing of enter key events if no control is selected
                }

            });
            $wrapper.on('hidden.bs.modal', function (e) {
                qcubed.recordControlModification(id, "_IsOpen", false);
            });
        },
        open: function() {
            var $control = this.element,
                id = $control.attr('id'),
                wrapperId = id + "_ctl",
                $wrapper = $('#' + wrapperId);

            $wrapper.modal("show");
        },
        close: function() {
            var $control = this.element,
                id = $control.attr('id'),
                wrapperId = id + "_ctl",
                $wrapper = $('#' + wrapperId);

            $wrapper.modal("hide");

        },
        showButton: function(btnId, visible) {
            var $control = $(this.element).parent(),
                $button = $control.find("button[data-btnid=" + btnId + "]");

            if ($button) {
                if (visible) {
                    $button.show();
                } else {
                    $button.hide();
                }
            }
        },
        setButtonCss: function(btnId, css) {
            var $control = this.element,
                $button = $control.find("button[data-btnid=" + btnId + "]");

            if ($button) {
                $button.css(css);
            }
        },
        /**
         * Uses bootstrap to perform a confirm message.
         * @param string message
         * @param function success The function to execute if the user confirms the message. Otherwise, no function is executed.
         */
        confirm: function(message, success) {
            var $form = $(this.element).closest('form');
            var $w = $('<div id="bsConfirm_ctl" class="modal"></div>');
            var $m = $('<div id="bsConfirm">' + message + '</div>');
            $w.append($m);
            $form.append($w);
            $m.bsModal({
                hasCloseButton: false,
                backdrop: "static",
                title: " ",
                headerClasses: "bg-warning",
                buttons: [
                    {label: "No", close: true, click: false},
                    {label: "Yes", close: true, click: success}
                ]
            });
            $m.bsModal("open");
            $w.on('hidden.bs.modal', function () {
                $w.remove();
            });
        },
        _buttonClick: function(event) {
            var objButton = event.data;
            var self = objButton.instance;

            if (objButton.confirm) {
                self.confirm(objButton.confirm, function () {
                    self._recordButtonClick(event);
                });
            } else {
                self._recordButtonClick(event);
            }
        },
        _recordButtonClick: function(event) {
            var objButton = event.data;
            var self = objButton.instance;
            var controlId = self.element.attr('id');

            qcubed.recordControlModification(controlId, "_ClickedButton", objButton.id);
            self.element.trigger("QDialog_Button", objButton.id);
        },
        _destroy: function () {
            this.close(); // there currently is no tear down function for modals. v4 of bootstrap is supposed to have one.
            this._super();
        }
    });
});