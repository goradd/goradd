/**
 * Widget script designed to be attached to a select table. Depends on ScrollIntoView.
 *
 * TODO: Capture focus and improve aria experience
 */

goradd.widget( "goradd.selectTable", {
    options: {
        selectedId: ""
    },
    _create: function() {
        this.on('click', 'tr', this._handleRowClick);
    },
    _handleRowClick: function(event) {
        this._selectRow(event.goradd.match);
        var $row = g$(event.goradd.match);
        var selId = $row.data("id");

        if ($row.hasClass ("selected") && selId !== this.selectedId) {
            this.selectedId = selId;
            goradd.setControlValue(this.element.id, "selectedId", selId);
            this.trigger('rowselected', selId);
        }
    },
    _selectRow: function(row) {
        if (g$(row).hasClass ("sel")) {
            var sel = g$(row.parentElement).qs('.selected');
            if (sel) {
                g$(sel).removeClass('selected');
            }
            g$(row).addClass('selected');
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
