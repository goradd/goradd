
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
        var itemHeight = itemBottom - itemTop;
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

    goradd.extend(goradd.g.prototype, {
        scrollIntoView: function () {
            var self = this;
            var curEl = this.element;
            var objTop = curEl.offsetTop; // This is the top from the non-static parent, which we will search for below.
            var ps = this.parents();
            var found = false;
            goradd.each(ps, function (i, el) {
                if (el.position !== "static") {
                    var o = el.css("overflow-y");
                    if (o === "auto" || o === "scroll") {
                        // this is the scrollable parent. It might not be currently scrollable, but for our purposes, its what we want
                        found = true;
                        var objBottom = objTop + self.element.clientHeight;
                        var s = window.getComputedStyle(el, null);
                        var h = el.clientHeight - s.paddingTop - s.paddingBottom;
                        var t = el.scrollTop;
                        el.scrollTop = _newScrollTop(t, t+h, el.scrollHeight, objTop, objBottom);
                        return false;
                    } else {
                        // the non-static parent is not the scrollable parent, so we need to add to the next non-static parent and keep looking
                        curEl = el;
                        objTop += curEl.offsetTop;
                    }
                }
            });
            if (!found) {
                // We went all the way up to the window and found no scrollable parent that was not styled
                // so that it would scroll (non-static). So, we will scroll the entire window.

                // TODO
            }
        }
    });
})();
