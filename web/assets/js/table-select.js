/**
 * Widget script designed to be attached to a select table.
 *
 * TODO: Capture focus and improve aria experience
 */

goradd.widget( "goradd.selectTable", {
    options: {
        selectedId: ""
    },
    _create: function() {
        this._super();
        this.on('click', 'tr', this._handleRowClick, {bubbles: true});
        this._initSelectedId();
        this.on('keydown', 'tr', this._handleKeyDown);
        this.on('focus', 'tr', this._handleFocus);
    },
    _initSelectedId: function() {
        var row;
        if (!!this.options.selectedId) {
            row = this.qs("tr[data-id='" + this.options.selectedId + "']");
        }
        this._selectRow(row);
        this.showSelectedItem();
    },
    _handleRowClick: function(event) {
        var prevSelId = this.options.selectedId;
        this._selectRow(event.goradd.match);
        var $row = g$(event.goradd.match);
        var selId = this.options.selectedId;

        if (selId !== prevSelId && $row.hasClass ("selected")) {
            goradd.setControlValue(this.element.id, "selectedId", selId);
            this.trigger('rowselected', selId);
            $row.trigger("focus");
        }
        event.preventDefault();
    },
    _handleKeyDown: function(e) {
        var row;

        switch (e.keyCode) {
            case 38: // up
            case 40: // down
            case 35: // end
            case 36: // home
                if (!!this.options.selectedId) {
                    row = this.qs("tr[data-id='" + this.options.selectedId + "']");
                }
                if (!row) {
                    return;
                }
        }
        var prevId = this.options.selectedId;
        var newRow;
        switch (e.keyCode) {
            case 38: // up
                newRow = row.previousElementSibling;
                break;
            case 40: // down
                newRow = row.nextElementSibling;
                break;
            case 35: // end
                newRow = this.qs("tr:last-child");
                break;
            case 36: // home
                newRow = this.qs("tr");
                break;
        }

        if (newRow) {
            var newId = g$(newRow).data("id");
            if (newId !== prevId) {
                this._selectRow(newRow);
                goradd.setControlValue(this.element.id, "selectedId", this.options.selectedId);
                this.trigger('rowselected', this.options.selectedId);
                g$(newRow).trigger("focus");
                this.showSelectedItem();
            }
        }

        e.preventDefault();
    },
    _handleFocus: function(event) {
        if (!this.options.selectedId) {
            var t = event.target;
            this._selectRow(t);
        }
    },
    _selectRow: function(row) {
        var $r = g$(row);
        var sel;
        if ($r) {
            sel = g$(row.parentElement).qs('.selected');
        }
        if (sel) { // should we make sure we are not selecting same item?
            g$(sel).class('-selected');
            g$(sel).attr("aria-selected", false);
            if ($r){$r.attr("tabindex", false);}
            this.attr("aria-activedescendant", false);
        }
        if ($r && !$r.hasClass ("nosel")) {
            $r.class('+selected');
            $r.attr("aria-selected", true);
            $r.attr("tabindex", 0);
            this.attr("aria-activedescendant", row.id);
            this.options.selectedId = $r.data("id");
        } else {
            // provide for keyboard interaction to select first item
            row = this.qs("tr"); // get first row
            $r = g$(row);
            if ($r) {$r.attr("tabindex", "0");}
        }
    },
    _setOption: function(key, value) {
        this._super(key, value);
        if (key === "selectedId") {
            this._initSelectedId();
        }
    },
    showSelectedItem: function() {
        if (!this.options.selectedId) {
            return;
        }
        var row = this.qs("tr[data-id='" + this.options.selectedId + "']");
        if (row) {
            g$(row).scrollIntoView();
        }
    }

});
