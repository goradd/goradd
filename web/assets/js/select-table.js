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
        this._super();
        this.on('click', 'tr', this._handleRowClick);
        if (!!this.options.selectedId) {
            this._initSelectedId();
        }
    },
    _initSelectedId: function() {
        var row = this.qs("tr[data-id=" + this.options.selectedId + "]");
        this._selectRow(row);
        if (row) {
            //$g(row)._scrollIntoView();
        }
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
        var sel = g$(row.parentElement).qs('.selected');
        if (sel) {
            g$(sel).removeClass('selected');
        }
        if (!!row && g$(row).hasClass ("sel")) {
            g$(row).addClass('selected');
        }
    },
    _setOption: function(key, value) {
        this._super(key, value);
        if (key === "selectedId") {
            this._initSelectedId();
        }
    }

});
