/**
 * scrollIntoView extends the goradd widget to add a scrollIntoView function. This only scrolls vertically at this
 * point, not horizontally.
 */

(function() {
    /**
     * _newScrollTop figures out what the scroll top position should be given the parameters.
     * If the item is in view, it will return the viewTop to indicate there should be no change.
     * Otherwise it will return a new viewTop such that the item will just be in view.
     * This is tricky.
     * @param viewTop {number}
     * @param viewBottom {number}
     * @param scrollHeight {number}
     * @param itemTop {number}
     * @param itemBottom {number}
     * @private
     */
    function _newScrollTop(viewTop, viewBottom, scrollHeight, itemTop, itemBottom) {
        var viewHeight = viewBottom - viewTop;
        //var itemHeight = itemBottom - itemTop;
        var topPos = itemTop;
        var bottomPos = itemBottom - viewHeight;

        if (itemTop + viewHeight > scrollHeight) {
            // item is close to bottom, so find max top position
            topPos = scrollHeight - viewHeight;
            if (topPos < 0) {
                topPos = 0; // less items than would fill viewport
            }
        }

        if (bottomPos < 0) {
            bottomPos = 0; // less items than would fill viewport
        }

        if (bottomPos > topPos) {
            // item is bigger than view
            if (topPos <= viewTop && viewTop <= bottomPos) {
                // do nothing, item is fully visible
                return viewTop;
            } else {
                return topPos; // scroll the item so that top of item is visible
            }
        } else {
            if (bottomPos <= viewTop && viewTop <= topPos) {
                // do nothing, item is fully visible
                return viewTop;
            } else if (itemTop < viewTop) {
                return topPos;
            } else {
                return bottomPos;
            }
        }
    }

    goradd.extend(goradd.Control.prototype, {
        scrollIntoView: function () {
            var curEl = this.element;
            var scroller;
            var parents = this.parents();
            // Find wrapping scroller
            goradd.each(parents, function (i, el) {
                var o = g$(el).css("overflow-y");
                if (o === "auto" || o === "scroll") {
                    scroller = el;
                    return false;
                }
            });

            if (!!scroller) {
                var rEl = curEl.getBoundingClientRect();
                var rScroll = scroller.getBoundingClientRect();
                var viewTop = scroller.scrollTop;
                var viewBottom = viewTop + rScroll.bottom - rScroll.top;
                var scrollHeight = scroller.scrollHeight;
                var itemTop = rEl.top - rScroll.top + viewTop;
                var itemBottom = itemTop + rEl.bottom - rEl.top;
                var newTop = _newScrollTop(viewTop, viewBottom, scrollHeight, itemTop, itemBottom);

                // TODO: Create and use a simple animation library to do the scrolling.
                scroller.scrollTop = newTop;
            } else {
                // We went all the way up to the window and found no scrollable parent.
                // So, we will scroll the entire window.

                // TODO
            }
        }
    });
})();
