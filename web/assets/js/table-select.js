/**
 * Widget script designed to be attached to a select table.
 */

(function() {
    goradd.SelectTable = goradd.extendWidget({
        constructor: function(element, options) {
            var optionDefaults = {
                selectedId: "",
                scrollable: false,
                reselect: false
            };
            options = goradd.extendOptions(optionDefaults, options);
            this._super(element, options);
            this.on('click', 'tr', this._handleRowClick, {bubbles: true});
            this._initSelectedId();
            if (this.options.scrollable) {
                var scroller = goradd.tagBuilder("div")
                    .attr("id", this.id + "_scroller")
                    .attr("tabindex", 0)
                    .attr("style", "overflow-y:auto;padding-right:10px")
                    .wrap(this.element);
                g$(scroller).on('keydown', [this, this._handleKeyDown]);
                this.attr("tabindex", false);

            } else {
                this.on('keydown', this._handleKeyDown);
            }
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

            if ((this.options.reselect || selId !== prevSelId) && $row.hasClass ("selected")) {
                goradd.setControlValue(this.element.id, "selectedId", selId);
                this.trigger('gr-rowselected', selId);
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
                    this.trigger('gr-rowselected', this.options.selectedId);
                    this.showSelectedItem();
                    e.preventDefault();
                }
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
                this.attr("aria-activedescendant", false);
            }
            if ($r && !$r.hasClass ("nosel")) {
                $r.class('+selected');
                $r.attr("aria-selected", true);
                this.attr("aria-activedescendant", row.id);
                this.options.selectedId = $r.data("id");
            } else {
                // provide for keyboard interaction to select first item
                row = this.qs("tr"); // get first row
                $r = g$(row);
            }
        },
        _setOption: function(key, value) {
            this._super._setOption(key, value);
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

    goradd.registerWidget("goradd.SelectTable", goradd.SelectTable);
})();
