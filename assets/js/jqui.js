/////////////////////////////////////
// Drag and drop support
/////////////////////////////////////

goradd.draggable = function (parentId, draggableId) {
    // we are working around some jQuery UI bugs here..
    $j('#' + parentId).on("dragstart", function () {
        var c = $j(this);
        c.data ("originalPosition", c.position());
    }).on("dragstop", function () {
        var c = $j(this);
        gr.setControlValue(draggableId, "_DragData", {originalPosition: {left: c.data("originalPosition").left, top: c.data("originalPosition").top}, position: {left: c.position().left, top: c.position().top}});
    });
};

goradd.droppable = function (parentId, droppableId) {
    $j('#' + parentId).on("drop", function (event, ui) {
        gr.setControlValue(droppableId, "_DroppedId", ui.draggable.attr("id"));
    });
};

goradd.resizable = function (parentId, resizeableId) {
    $j('#' + parentId).on("resizestart", function () {
        var c = $j(this);
        c.data ("oW", c.width());
        c.data ("oH", c.height());
    })
        .on("resizestop", function () {
            var c = $j(this);
            gr.setControlValue(resizeableId, "_ResizeData", {originalSize: {width: c.data("oW"), height: c.data("oH")} , size:{width: c.width(), height: c.height()}});
        });
};

/////////////////////////////////////
// JQueryUI Support
/////////////////////////////////////

goradd.dialog = function(controlId) {
    $j('#' + controlId).on ("keydown", "input,select", function(event) {
        // makes sure a return key fires the default button if there is one
        if (event.which === 13) {
            var b = $j(this).closest("[role=\'dialog\']").find("button[type=\'submit\']");
            if (b && b[0]) {
                b[0].click();
            }
            event.preventDefault();
        }
    });
};

goradd.accordion = function(controlId) {
    $j('#' + controlId).on("accordionactivate", function(event, ui) {
        goradd.setControlValue(controlId, "_SelectedIndex", $j(this).accordion("option", "active"));
        $j(this).trigger("change");
    });
};

goradd.progressbar = function(controlId) {
    $j('#' + controlId).on("progressbarchange", function (event, ui) {
        goradd.setControlValue(controlId, "_Value", $j(this).progressbar ("value"));
    });
};

goradd.selectable = function(controlId) {
    $j('#' + controlId).on("selectablestop", function (event, ui) {
        var strItems;

        strItems = "";
        $j(".ui-selected", this).each(function() {
            strItems = strItems + "," + this.id;
        });

        if (strItems) {
            strItems = strItems.substring (1);
        }
        goradd.setControlValue(controlId, "_SelectedItems", strItems);

    });
};

goradd.slider = function(controlId) {
    $j('#' + controlId).on("slidechange", function (event, ui) {
        if (ui.values && ui.values.length) {
            gr.setControlValue(controlId, "_Values", ui.values[0] + ',' +  ui.values[1]);
        } else {
            gr.setControlValue(controlId, "_Value", ui.value);
        }
    });
};

goradd.tabs = function(controlId) {
    $j('#' + controlId).on("tabsactivate", function(event, ui) {
        var i = $j(this).tabs( "option", "active" );
        var id = ui.newPanel ? ui.newPanel.attr("id") : null;
        gr.setControlValue(controlId, "_active", [i,id]);
    });
};

goradd.datagrid2 = function(controlId) {
    $j('#' + controlId).on("click", "thead tr th a", function(event, ui) {
        var cellIndex = $j(this).parent()[0].cellIndex;
        $j(this).trigger('qdg2sort', cellIndex); // Triggers the QDataGrid_SortEvent
        event.stopPropagation();
    });
};

goradd.dialog = function(controlId) {
    $j('#' + controlId).on("tabsactivate", function(event, ui) {
        var i = $j(this).tabs( "option", "active" );
        var id = ui.newPanel ? ui.newPanel.attr("id") : null;
        gr.setControlValue(controlId, "_active", [i,id]);
    });
};
