/**
 * Widget script designed to be attached to a select grid. Depends on ScrollIntoView.
 *
 * TODO: Capture focus and improve aria experience
 */
jQuery(function( $, undefined ) {

    $.widget( "goradd.selectTable", {
        options: {
            selectedId: ""
        },
        _create: function() {
            this._on({
                'click tr': this._handleRowClick
            });
        },
        _handleRowClick: function(event) {
            var $control = this.element,
                id = $control.attr('id');

            var $row = $(event.currentTarget);
            this._selectRow($row);
            var selId = $row.data("id");
            if ($row.hasClass ("selected") && selId != this.selectedId) {
                this.selectedId = selId;
                goradd.setControlValue(id, "selectedId", selId);
                $control.trigger('rowselected', selId);
            }
        },
        _selectRow: function($row) {
            if ($row.hasClass ("sel")) {
                $row.parent().find('.selected').toggleClass('selected');
                $row.toggleClass('selected');
            }
        },
        _setOption: function( key, value ) {
            this._super( key, value );
            if ( key === "selectedId" && value) {
                var $row = this.element.find("tr[data-id=" + value + "]");
                if ($row.length > 0) {
                    this._selectRow($row);
                    $row.scrollintoview();
                }
            }
        }

    });
});